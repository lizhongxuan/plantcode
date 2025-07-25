package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"ai-dev-platform/internal/ai"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/repository"
	"github.com/google/uuid"
)

type SpecService struct {
	db         *sql.DB
	aiManager  *ai.AIManager
	repo       *repository.MySQLRepository
}

func NewSpecService(db *sql.DB, aiManager *ai.AIManager, repo *repository.MySQLRepository) *SpecService {
	return &SpecService{
		db:        db,
		aiManager: aiManager,
		repo:      repo,
	}
}

// InitProjectSpec 初始化项目 Spec
func (s *SpecService) InitProjectSpec(ctx context.Context, projectID uuid.UUID) (*model.ProjectSpec, error) {
	specID := uuid.New()
	
	spec := &model.ProjectSpec{
		SpecID:       specID,
		ProjectID:    projectID,
		CurrentStage: model.SpecStageRequirements,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO project_specs (spec_id, project_id, current_stage, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query, spec.SpecID, spec.ProjectID, spec.CurrentStage, spec.CreatedAt, spec.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create project spec: %w", err)
	}

	return spec, nil
}

// GetProjectSpec 获取项目 Spec
func (s *SpecService) GetProjectSpec(ctx context.Context, projectID uuid.UUID) (*model.ProjectSpec, error) {
	query := `
		SELECT spec_id, project_id, current_stage, created_at, updated_at
		FROM project_specs
		WHERE project_id = ?
	`
	
	spec := &model.ProjectSpec{}
	err := s.db.QueryRowContext(ctx, query, projectID).Scan(
		&spec.SpecID, &spec.ProjectID, &spec.CurrentStage, &spec.CreatedAt, &spec.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有 spec，自动创建一个
			return s.InitProjectSpec(ctx, projectID)
		}
		return nil, fmt.Errorf("failed to get project spec: %w", err)
	}

	return spec, nil
}

// GenerateRequirements 生成需求文档
func (s *SpecService) GenerateRequirements(ctx context.Context, userID uuid.UUID, req *model.GenerateRequirementsRequest) (*model.RequirementsDoc, error) {
	// 获取用户AI配置（这里简化处理，使用默认配置）
	// 在实际项目中，应该从用户配置中获取AI提供商信息

	// 构建需求分析的提示词
	prompt := s.buildRequirementsPrompt(req)
	
	// 调用AI生成需求文档（使用ProjectChat作为通用接口）
	response, err := s.aiManager.ProjectChat(ctx, prompt, "", ai.ProviderOpenAI)
	if err != nil {
		return nil, fmt.Errorf("failed to generate requirements: %w", err)
	}

	// 解析AI响应
	reqDoc, err := s.parseRequirementsResponse(response.Message, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requirements response: %w", err)
	}

	// 保存到数据库
	if err := s.saveRequirementsDoc(ctx, reqDoc); err != nil {
		return nil, fmt.Errorf("failed to save requirements: %w", err)
	}

	// 更新项目 spec 状态
	if err := s.updateSpecStage(ctx, req.ProjectID, model.SpecStageRequirements); err != nil {
		log.Printf("Warning: failed to update spec stage: %v", err)
	}

	return reqDoc, nil
}

// GenerateDesign 生成设计文档
func (s *SpecService) GenerateDesign(ctx context.Context, userID uuid.UUID, req *model.GenerateDesignRequest) (*model.DesignDoc, error) {
	// 获取需求文档
	reqDoc, err := s.getRequirementsDoc(ctx, req.RequirementsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requirements doc: %w", err)
	}

	// 构建设计分析的提示词
	prompt := s.buildDesignPrompt(reqDoc, req)
	
	// 调用AI生成设计文档
	response, err := s.aiManager.ProjectChat(ctx, prompt, "", ai.ProviderOpenAI)
	if err != nil {
		return nil, fmt.Errorf("failed to generate design: %w", err)
	}

	// 解析AI响应
	designDoc, err := s.parseDesignResponse(response.Message, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse design response: %w", err)
	}

	// 保存到数据库
	if err := s.saveDesignDoc(ctx, designDoc); err != nil {
		return nil, fmt.Errorf("failed to save design: %w", err)
	}

	// 更新项目 spec 状态
	if err := s.updateSpecStage(ctx, req.ProjectID, model.SpecStageDesign); err != nil {
		log.Printf("Warning: failed to update spec stage: %v", err)
	}

	return designDoc, nil
}

// GenerateTasks 生成任务列表
func (s *SpecService) GenerateTasks(ctx context.Context, userID uuid.UUID, req *model.GenerateTasksRequest) (*model.TaskListDoc, error) {
	// 获取需求和设计文档
	reqDoc, err := s.getRequirementsDoc(ctx, req.RequirementsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requirements doc: %w", err)
	}

	designDoc, err := s.getDesignDoc(ctx, req.DesignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get design doc: %w", err)
	}

	// 构建任务分析的提示词
	prompt := s.buildTasksPrompt(reqDoc, designDoc, req)
	
	// 调用AI生成任务文档
	response, err := s.aiManager.ProjectChat(ctx, prompt, "", ai.ProviderOpenAI)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks: %w", err)
	}

	// 解析AI响应
	taskDoc, err := s.parseTasksResponse(response.Message, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks response: %w", err)
	}

	// 保存到数据库
	if err := s.saveTaskListDoc(ctx, taskDoc); err != nil {
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	// 更新项目 spec 状态
	if err := s.updateSpecStage(ctx, req.ProjectID, model.SpecStageTasks); err != nil {
		log.Printf("Warning: failed to update spec stage: %v", err)
	}

	return taskDoc, nil
}

// buildRequirementsPrompt 构建需求分析提示词
func (s *SpecService) buildRequirementsPrompt(req *model.GenerateRequirementsRequest) string {
	prompt := fmt.Sprintf(`
你是一个专业的产品经理和需求分析师。请基于以下信息生成详细的需求文档：

**项目信息：**
- 项目类型：%s
- 初始需求：%s
`, req.ProjectType, req.InitialPrompt)

	if req.TargetAudience != "" {
		prompt += fmt.Sprintf("- 目标用户：%s\n", req.TargetAudience)
	}

	if req.BusinessGoals != "" {
		prompt += fmt.Sprintf("- 业务目标：%s\n", req.BusinessGoals)
	}

	prompt += `
请按照以下格式生成需求文档，使用 EARS (Easy Approach to Requirements Syntax) 语法：

**输出格式（JSON）：**
` + "```json" + `
{
  "content": "完整的需求文档内容（Markdown格式）",
  "user_stories": [
    {
      "title": "用户故事标题",
      "description": "作为[角色]，我希望[功能]，以便[价值]",
      "acceptance_criteria": ["验收标准1", "验收标准2"],
      "priority": "high/medium/low",
      "story_points": 5
    }
  ],
  "functional_requirements": ["功能需求1", "功能需求2"],
  "non_functional_requirements": ["性能要求", "安全要求", "可用性要求"],
  "assumptions": ["假设条件1", "假设条件2"],
  "edge_cases": ["边界情况1", "边界情况2"]
}
` + "```" + `

**要求：**
1. 使用EARS语法编写需求
2. 包含详细的用户故事和验收标准
3. 考虑边界情况和异常处理
4. 明确功能和非功能需求
5. 提供清晰的假设条件
`

	return prompt
}

// buildDesignPrompt 构建设计分析提示词
func (s *SpecService) buildDesignPrompt(reqDoc *model.RequirementsDoc, req *model.GenerateDesignRequest) string {
	prompt := fmt.Sprintf(`
你是一个专业的系统架构师和技术设计师。请基于以下需求文档生成技术设计方案：

**需求文档：**
%s

**架构偏好：**
- 架构风格：%s
`, reqDoc.Content, req.ArchitectureStyle)

	if req.FocusAreas != "" {
		prompt += fmt.Sprintf("- 重点关注：%s\n", req.FocusAreas)
	}

	prompt += `
请按照以下格式生成设计文档：

**输出格式（JSON）：**
` + "```json" + `
{
  "content": "完整的设计文档内容（Markdown格式）",
  "puml_diagrams": [
    {
      "title": "系统架构图",
      "type": "component",
      "code": "@startuml\n!define RECTANGLE class\n...\n@enduml",
      "description": "图表说明"
    }
  ],
  "interfaces": [
    {
      "name": "接口名称",
      "code": "interface UserInterface {\n  id: string;\n  name: string;\n}",
      "description": "接口描述"
    }
  ],
  "api_endpoints": [
    {
      "path": "/api/users",
      "method": "GET",
      "description": "获取用户列表",
      "request_body": {},
      "response_body": {"users": []},
      "headers": {"Authorization": "Bearer token"}
    }
  ],
  "database_schema": "数据库设计说明",
  "architecture_notes": ["架构说明1", "架构说明2"]
}
` + "```" + `

**要求：**
1. 生成多种类型的PUML图表（架构图、时序图、活动图等）
2. 定义清晰的数据接口
3. 设计RESTful API端点
4. 考虑数据库设计
5. 提供架构决策说明
`

	return prompt
}

// buildTasksPrompt 构建任务分析提示词
func (s *SpecService) buildTasksPrompt(reqDoc *model.RequirementsDoc, designDoc *model.DesignDoc, req *model.GenerateTasksRequest) string {
	teamSize := 3
	sprintDuration := 2
	
	if req.TeamSize != nil {
		teamSize = *req.TeamSize
	}
	if req.SprintDuration != nil {
		sprintDuration = *req.SprintDuration
	}

	prompt := fmt.Sprintf(`
你是一个专业的项目管理师和敏捷教练。请基于需求文档和设计文档生成开发任务列表：

**需求文档：**
%s

**设计文档：**
%s

**团队信息：**
- 团队大小：%d 人
- Sprint 周期：%d 周
`, reqDoc.Content, designDoc.Content, teamSize, sprintDuration)

	prompt += `
请按照以下格式生成任务文档：

**输出格式（JSON）：**
` + "```json" + `
{
  "content": "完整的任务规划文档（Markdown格式）",
  "tasks": [
    {
      "title": "任务标题",
      "description": "详细描述",
      "type": "feature/bug/refactor/test/docs",
      "priority": "high/medium/low",
      "status": "todo",
      "estimated_hours": 8,
      "dependencies": [],
      "user_story_id": "关联的用户故事ID"
    }
  ],
  "test_cases": [
    {
      "title": "测试用例标题",
      "description": "测试描述",
      "type": "unit/integration/e2e/api",
      "steps": ["步骤1", "步骤2"],
      "expected_result": "预期结果"
    }
  ],
  "estimated_total_hours": 120,
  "milestones": ["里程碑1", "里程碑2"]
}
` + "```" + `

**要求：**
1. 任务分解要详细且可执行
2. 估算工作量要合理
3. 包含完整的测试策略
4. 考虑任务依赖关系
5. 设定清晰的里程碑
`

	return prompt
}

// parseRequirementsResponse 解析需求文档响应
func (s *SpecService) parseRequirementsResponse(content string, projectID uuid.UUID) (*model.RequirementsDoc, error) {
	// 提取JSON部分
	jsonStr := s.extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	// 解析JSON
	var parsed struct {
		Content                   string   `json:"content"`
		UserStories              []struct {
			Title              string   `json:"title"`
			Description        string   `json:"description"`
			AcceptanceCriteria []string `json:"acceptance_criteria"`
			Priority           string   `json:"priority"`
			StoryPoints        *int     `json:"story_points"`
		} `json:"user_stories"`
		FunctionalRequirements    []string `json:"functional_requirements"`
		NonFunctionalRequirements []string `json:"non_functional_requirements"`
		Assumptions               []string `json:"assumptions"`
		EdgeCases                 []string `json:"edge_cases"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 创建需求文档
	reqDoc := &model.RequirementsDoc{
		ID:        uuid.New(),
		ProjectID: projectID,
		Content:   parsed.Content,
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 转换为JSON字符串存储
	if assumptions, err := json.Marshal(parsed.Assumptions); err == nil {
		reqDoc.Assumptions = string(assumptions)
	}
	if edgeCases, err := json.Marshal(parsed.EdgeCases); err == nil {
		reqDoc.EdgeCases = string(edgeCases)
	}
	if funcReq, err := json.Marshal(parsed.FunctionalRequirements); err == nil {
		reqDoc.FunctionalRequirements = string(funcReq)
	}
	if nonFuncReq, err := json.Marshal(parsed.NonFunctionalRequirements); err == nil {
		reqDoc.NonFunctionalRequirements = string(nonFuncReq)
	}

	return reqDoc, nil
}

// parseDesignResponse 解析设计文档响应
func (s *SpecService) parseDesignResponse(content string, projectID uuid.UUID) (*model.DesignDoc, error) {
	// 提取JSON部分
	jsonStr := s.extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	// 解析JSON
	var parsed struct {
		Content         string `json:"content"`
		PUMLDiagrams   []struct {
			Title       string `json:"title"`
			Type        string `json:"type"`
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"puml_diagrams"`
		Interfaces     []struct {
			Name        string `json:"name"`
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"interfaces"`
		APIEndpoints   []struct {
			Path         string      `json:"path"`
			Method       string      `json:"method"`
			Description  string      `json:"description"`
			RequestBody  interface{} `json:"request_body"`
			ResponseBody interface{} `json:"response_body"`
			Headers      map[string]string `json:"headers"`
		} `json:"api_endpoints"`
		DatabaseSchema    string   `json:"database_schema"`
		ArchitectureNotes []string `json:"architecture_notes"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 创建设计文档
	designDoc := &model.DesignDoc{
		ID:               uuid.New(),
		ProjectID:        projectID,
		Content:          parsed.Content,
		DatabaseSchema:   parsed.DatabaseSchema,
		Version:          1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 转换为JSON字符串存储
	if notes, err := json.Marshal(parsed.ArchitectureNotes); err == nil {
		designDoc.ArchitectureNotes = string(notes)
	}

	return designDoc, nil
}

// parseTasksResponse 解析任务文档响应
func (s *SpecService) parseTasksResponse(content string, projectID uuid.UUID) (*model.TaskListDoc, error) {
	// 提取JSON部分
	jsonStr := s.extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	// 解析JSON
	var parsed struct {
		Content             string `json:"content"`
		Tasks              []struct {
			Title          string   `json:"title"`
			Description    string   `json:"description"`
			Type           string   `json:"type"`
			Priority       string   `json:"priority"`
			Status         string   `json:"status"`
			EstimatedHours int      `json:"estimated_hours"`
			Dependencies   []string `json:"dependencies"`
			UserStoryID    string   `json:"user_story_id"`
		} `json:"tasks"`
		TestCases          []struct {
			Title          string   `json:"title"`
			Description    string   `json:"description"`
			Type           string   `json:"type"`
			Steps          []string `json:"steps"`
			ExpectedResult string   `json:"expected_result"`
		} `json:"test_cases"`
		EstimatedTotalHours int      `json:"estimated_total_hours"`
		Milestones          []string `json:"milestones"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 创建任务文档
	taskDoc := &model.TaskListDoc{
		ID:                  uuid.New(),
		ProjectID:           projectID,
		Content:             parsed.Content,
		EstimatedTotalHours: parsed.EstimatedTotalHours,
		Version:             1,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// 转换为JSON字符串存储
	if milestones, err := json.Marshal(parsed.Milestones); err == nil {
		taskDoc.Milestones = string(milestones)
	}

	return taskDoc, nil
}

// extractJSON 从响应中提取JSON部分
func (s *SpecService) extractJSON(content string) string {
	// 寻找JSON代码块
	if start := strings.Index(content, "```json"); start != -1 {
		start += 7 // 跳过 ```json
		if end := strings.Index(content[start:], "```"); end != -1 {
			return strings.TrimSpace(content[start : start+end])
		}
	}
	
	// 寻找普通的JSON对象
	if start := strings.Index(content, "{"); start != -1 {
		if end := strings.LastIndex(content, "}"); end != -1 && end > start {
			return strings.TrimSpace(content[start : end+1])
		}
	}

	return ""
}

// 数据库操作方法
func (s *SpecService) saveRequirementsDoc(ctx context.Context, doc *model.RequirementsDoc) error {
	query := `
		INSERT INTO requirements_docs (id, project_id, content, assumptions, edge_cases, 
			functional_requirements, non_functional_requirements, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query,
		doc.ID, doc.ProjectID, doc.Content, doc.Assumptions, doc.EdgeCases,
		doc.FunctionalRequirements, doc.NonFunctionalRequirements,
		doc.Version, doc.CreatedAt, doc.UpdatedAt,
	)
	return err
}

func (s *SpecService) saveDesignDoc(ctx context.Context, doc *model.DesignDoc) error {
	query := `
		INSERT INTO design_docs (id, project_id, content, database_schema, architecture_notes, 
			version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query,
		doc.ID, doc.ProjectID, doc.Content, doc.DatabaseSchema, doc.ArchitectureNotes,
		doc.Version, doc.CreatedAt, doc.UpdatedAt,
	)
	return err
}

func (s *SpecService) saveTaskListDoc(ctx context.Context, doc *model.TaskListDoc) error {
	query := `
		INSERT INTO task_list_docs (id, project_id, content, estimated_total_hours, milestones,
			version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := s.db.ExecContext(ctx, query,
		doc.ID, doc.ProjectID, doc.Content, doc.EstimatedTotalHours, doc.Milestones,
		doc.Version, doc.CreatedAt, doc.UpdatedAt,
	)
	return err
}

func (s *SpecService) getRequirementsDoc(ctx context.Context, reqID uuid.UUID) (*model.RequirementsDoc, error) {
	query := `
		SELECT id, project_id, content, assumptions, edge_cases, 
			functional_requirements, non_functional_requirements, version, created_at, updated_at
		FROM requirements_docs
		WHERE id = ?
	`
	
	doc := &model.RequirementsDoc{}
	err := s.db.QueryRowContext(ctx, query, reqID).Scan(
		&doc.ID, &doc.ProjectID, &doc.Content, &doc.Assumptions, &doc.EdgeCases,
		&doc.FunctionalRequirements, &doc.NonFunctionalRequirements,
		&doc.Version, &doc.CreatedAt, &doc.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *SpecService) getDesignDoc(ctx context.Context, designID uuid.UUID) (*model.DesignDoc, error) {
	query := `
		SELECT id, project_id, content, database_schema, architecture_notes,
			version, created_at, updated_at
		FROM design_docs
		WHERE id = ?
	`
	
	doc := &model.DesignDoc{}
	err := s.db.QueryRowContext(ctx, query, designID).Scan(
		&doc.ID, &doc.ProjectID, &doc.Content, &doc.DatabaseSchema, &doc.ArchitectureNotes,
		&doc.Version, &doc.CreatedAt, &doc.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *SpecService) updateSpecStage(ctx context.Context, projectID uuid.UUID, stage string) error {
	query := `
		UPDATE project_specs 
		SET current_stage = ?, updated_at = ?
		WHERE project_id = ?
	`
	
	_, err := s.db.ExecContext(ctx, query, stage, time.Now(), projectID)
	return err
}