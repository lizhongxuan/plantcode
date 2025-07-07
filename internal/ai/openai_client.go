package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// OpenAIClient OpenAI客户端实现
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey  string
	BaseURL string // 可选，默认为OpenAI官方API
	Model   string // 可选，默认为gpt-4
}

// NewOpenAIClient 创建OpenAI客户端
func NewOpenAIClient(config OpenAIConfig) *OpenAIClient {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	
	model := config.Model
	if model == "" {
		model = "gpt-4"
	}

	return &OpenAIClient{
		apiKey:  config.APIKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GetProvider 返回AI服务提供商类型
func (c *OpenAIClient) GetProvider() AIProvider {
	return ProviderOpenAI
}

// AnalyzeRequirement 分析业务需求
func (c *OpenAIClient) AnalyzeRequirement(ctx context.Context, requirement string) (*RequirementAnalysis, error) {
	prompt := c.buildAnalysisPrompt(requirement)
	
	response, err := c.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用OpenAI API失败: %w", err)
	}

	analysis, err := c.parseAnalysisResponse(response.Content, requirement)
	if err != nil {
		return nil, fmt.Errorf("解析需求分析结果失败: %w", err)
	}

	return analysis, nil
}

// GenerateQuestions 基于分析结果生成补充问题
func (c *OpenAIClient) GenerateQuestions(ctx context.Context, analysis *RequirementAnalysis) ([]Question, error) {
	prompt := c.buildQuestionsPrompt(analysis)
	
	response, err := c.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用OpenAI API失败: %w", err)
	}

	questions, err := c.parseQuestionsResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析问题生成结果失败: %w", err)
	}

	return questions, nil
}

// GeneratePUML 生成PUML图表代码
func (c *OpenAIClient) GeneratePUML(ctx context.Context, analysis *RequirementAnalysis, diagramType PUMLType) (*PUMLDiagram, error) {
	prompt := c.buildPUMLPrompt(analysis, diagramType)
	
	response, err := c.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用OpenAI API失败: %w", err)
	}

	diagram, err := c.parsePUMLResponse(response.Content, analysis.ProjectID, diagramType)
	if err != nil {
		return nil, fmt.Errorf("解析PUML生成结果失败: %w", err)
	}

	return diagram, nil
}

// GenerateDocument 生成开发文档
func (c *OpenAIClient) GenerateDocument(ctx context.Context, analysis *RequirementAnalysis) (*DevelopmentDocument, error) {
	prompt := c.buildDocumentPrompt(analysis)
	
	response, err := c.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用OpenAI API失败: %w", err)
	}

	document, err := c.parseDocumentResponse(response.Content, analysis.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("解析文档生成结果失败: %w", err)
	}

	return document, nil
}

// ProjectChat 项目上下文AI对话
func (c *OpenAIClient) ProjectChat(ctx context.Context, message, context string) (*ProjectChatResponse, error) {
	prompt := c.buildProjectChatPrompt(message, context)
	
	response, err := c.callOpenAI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用OpenAI API失败: %w", err)
	}

	chatResponse, err := c.parseProjectChatResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析项目对话结果失败: %w", err)
	}

	return chatResponse, nil
}

// callOpenAI 调用OpenAI API
func (c *OpenAIClient) callOpenAI(ctx context.Context, prompt string) (*AIResponse, error) {
	req := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的业务分析师和软件架构师。请严格按照要求的JSON格式回复，不要添加额外的说明文字。",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":   2000,
		"temperature":  0.3,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("构建请求数据失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp struct {
		ID      string `json:"id"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("解析OpenAI响应失败: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI响应中没有生成内容")
	}

	return &AIResponse{
		ID:      openAIResp.ID,
		Content: openAIResp.Choices[0].Message.Content,
		Usage: AIUsage{
			PromptTokens:     openAIResp.Usage.PromptTokens,
			CompletionTokens: openAIResp.Usage.CompletionTokens,
			TotalTokens:      openAIResp.Usage.TotalTokens,
		},
		Model: openAIResp.Model,
	}, nil
}

// buildAnalysisPrompt 构建需求分析的提示语
func (c *OpenAIClient) buildAnalysisPrompt(requirement string) string {
	return fmt.Sprintf(`请分析以下业务需求，提取关键信息并识别缺失的信息。

业务需求：
%s

请按照以下JSON格式返回分析结果：
{
  "core_functions": ["功能1", "功能2"],
  "roles": ["角色1", "角色2"],
  "business_processes": [
    {
      "name": "流程名称",
      "description": "流程描述", 
      "steps": ["步骤1", "步骤2"],
      "actors": ["参与者1", "参与者2"]
    }
  ],
  "data_entities": [
    {
      "name": "实体名称",
      "description": "实体描述",
      "attributes": [
        {
          "name": "属性名",
          "type": "数据类型",
          "required": true,
          "description": "属性描述"
        }
      ],
      "relations": [
        {
          "target_entity": "目标实体",
          "relation_type": "one-to-many",
          "description": "关系描述"
        }
      ]
    }
  ],
  "missing_info": ["缺失信息1", "缺失信息2"],
  "completion_score": 0.7
}

注意：
1. 仔细分析业务逻辑，识别所有可能的功能点
2. 数据实体要包含完整的属性定义
3. 缺失信息要具体指出哪些业务细节不明确
4. 完整度评分基于需求描述的详细程度(0-1之间)`, requirement)
}

// buildQuestionsPrompt 构建问题生成的提示语
func (c *OpenAIClient) buildQuestionsPrompt(analysis *RequirementAnalysis) string {
	missingInfo := strings.Join(analysis.MissingInfo, "\n- ")
	
	return fmt.Sprintf(`基于以下需求分析中的缺失信息，生成具体的补充问题。

缺失信息：
- %s

请按照以下JSON格式返回问题列表：
{
  "questions": [
    {
      "category": "business_rule",
      "content": "具体的问题内容",
      "options": ["选项1", "选项2"],
      "priority": 3,
      "target_info": "目标获取的信息类型"
    }
  ]
}

问题分类包括：
- business_rule: 业务规则
- exception_handling: 异常处理
- data_structure: 数据结构
- external_interface: 外部接口
- performance_requirement: 性能需求
- security_requirement: 安全需求

优先级1-5，5最高。`, missingInfo)
}

// buildPUMLPrompt 构建PUML生成的提示语
func (c *OpenAIClient) buildPUMLPrompt(analysis *RequirementAnalysis, diagramType PUMLType) string {
	var diagramDescription string
	var example string
	
	switch diagramType {
	case PUMLTypeBusinessFlow:
		diagramDescription = "业务流程图（活动图）"
		example = `@startuml 业务流程图
start
:用户登录;
if (验证成功?) then (是)
  :进入系统;
  :执行操作;
else (否)
  :显示错误信息;
endif
stop
@enduml`
	case PUMLTypeArchitecture:
		diagramDescription = "系统架构图（组件图）"
		example = `@startuml 系统架构图
package "前端" {
  [用户界面]
}
package "后端" {
  [API服务]
  [业务逻辑]
}
database "数据库" {
  [用户数据]
}
[用户界面] --> [API服务]
[API服务] --> [业务逻辑]
[业务逻辑] --> [用户数据]
@enduml`
	case PUMLTypeSequence:
		diagramDescription = "序列图"
		example = `@startuml 序列图
actor 用户
participant 前端
participant 后端
participant 数据库
用户 -> 前端: 发起请求
前端 -> 后端: API调用
后端 -> 数据库: 查询数据
数据库 -> 后端: 返回数据
后端 -> 前端: 返回结果
前端 -> 用户: 显示结果
@enduml`
	case PUMLTypeDataModel:
		diagramDescription = "数据模型图（ER图）"
		example = `@startuml 数据模型图
!define table(x) class x << (T,#FFAAAA) >>
!define primary_key(x) <u>x</u>
!define foreign_key(x) <i>x</i>

table(用户) {
  primary_key(用户ID) : bigint
  用户名 : varchar(50)
  邮箱 : varchar(100)
  密码哈希 : varchar(255)
  创建时间 : datetime
}

table(项目) {
  primary_key(项目ID) : bigint
  foreign_key(用户ID) : bigint
  项目名称 : varchar(100)
  描述 : text
  状态 : varchar(20)
  创建时间 : datetime
}

用户 ||--o{ 项目 : 拥有
@enduml`
	case PUMLTypeClass:
		diagramDescription = "类图"
		example = `@startuml 类图
class 用户服务 {
  +登录(邮箱, 密码) : 用户
  +注册(用户信息) : 用户
  +获取用户信息(用户ID) : 用户
}

class 项目服务 {
  +创建项目(项目信息) : 项目
  +获取项目列表(用户ID) : 项目[]
  +更新项目(项目ID, 项目信息) : 项目
}

用户服务 --> 项目服务 : 使用
@enduml`
	}

	coreFunc := strings.Join(analysis.CoreFunctions, "\n- ")
	
	return fmt.Sprintf(`基于以下需求分析结果，生成%s的PlantUML代码。

核心功能：
- %s

请按照以下JSON格式返回：
{
  "title": "图表标题",
  "content": "完整的PlantUML代码",
  "description": "图表说明"
}

示例格式：
%s

要求：
1. 代码要完整可执行
2. 包含所有主要功能模块
3. 体现业务流程逻辑关系
4. 使用中文标注`, diagramDescription, coreFunc, example)
}

// buildDocumentPrompt 构建文档生成的提示语
func (c *OpenAIClient) buildDocumentPrompt(analysis *RequirementAnalysis) string {
	entities := make([]string, len(analysis.DataEntities))
	for i, entity := range analysis.DataEntities {
		entities[i] = entity.Name
	}
	entitiesStr := strings.Join(entities, ", ")
	
	return fmt.Sprintf(`基于以下需求分析结果，生成详细的开发文档。

核心功能：%s
数据实体：%s

请按照以下JSON格式返回完整的开发文档：
{
  "function_modules": [
    {
      "name": "模块名称",
      "description": "模块描述",
      "sub_modules": ["子模块1", "子模块2"],
      "dependencies": ["依赖模块1"],
      "priority": 1,
      "complexity": "medium",
      "estimated_hours": 40
    }
  ],
  "development_plan": {
    "phases": [
      {
        "name": "阶段名称",
        "description": "阶段描述",
        "tasks": ["任务1", "任务2"],
        "duration": "2周",
        "dependencies": []
      }
    ],
    "duration": "总开发周期",
    "resources": "资源需求描述"
  },
  "tech_stack": {
    "backend": {
      "recommended": "Go",
      "alternatives": ["Java", "Python"],
      "reason": "选择理由"
    },
    "frontend": {
      "recommended": "React",
      "alternatives": ["Vue", "Angular"],
      "reason": "选择理由"
    },
    "database": {
      "recommended": "MySQL",
      "alternatives": ["PostgreSQL"],
      "reason": "选择理由"
    }
  },
  "api_design": [
    {
      "path": "/api/endpoint",
      "method": "POST",
      "summary": "API简介",
      "description": "详细描述",
      "parameters": [],
      "request_body": {
        "description": "请求体描述",
        "schema": {"type": "object"},
        "required": true
      },
      "responses": {
        "200": {
          "description": "成功响应",
          "schema": {"type": "object"}
        }
      }
    }
  ]
}

要求：
1. 功能模块要完整覆盖所有核心功能
2. 开发计划要有明确的时间线
3. 技术选型要有合理的理由
4. API设计要包含主要的业务接口`, strings.Join(analysis.CoreFunctions, ", "), entitiesStr)
}

// parseAnalysisResponse 解析需求分析响应
func (c *OpenAIClient) parseAnalysisResponse(content, originalText string) (*RequirementAnalysis, error) {
	// 提取JSON部分
	jsonContent := extractJSON(content)
	
	var result struct {
		CoreFunctions     []string          `json:"core_functions"`
		Roles             []string          `json:"roles"`
		BusinessProcesses []BusinessProcess `json:"business_processes"`
		DataEntities      []DataEntity      `json:"data_entities"`
		MissingInfo       []string          `json:"missing_info"`
		CompletionScore   float64           `json:"completion_score"`
	}
	
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	
	return &RequirementAnalysis{
		ID:                uuid.New().String(),
		OriginalText:      originalText,
		CoreFunctions:     result.CoreFunctions,
		Roles:             result.Roles,
		BusinessProcesses: result.BusinessProcesses,
		DataEntities:      result.DataEntities,
		MissingInfo:       result.MissingInfo,
		CompletionScore:   result.CompletionScore,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

// parseQuestionsResponse 解析问题生成响应
func (c *OpenAIClient) parseQuestionsResponse(content string) ([]Question, error) {
	jsonContent := extractJSON(content)
	
	var result struct {
		Questions []struct {
			Category   string   `json:"category"`
			Content    string   `json:"content"`
			Options    []string `json:"options"`
			Priority   int      `json:"priority"`
			TargetInfo string   `json:"target_info"`
		} `json:"questions"`
	}
	
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	
	questions := make([]Question, len(result.Questions))
	for i, q := range result.Questions {
		questions[i] = Question{
			ID:         uuid.New().String(),
			Category:   q.Category,
			Content:    q.Content,
			Options:    q.Options,
			Priority:   q.Priority,
			TargetInfo: q.TargetInfo,
		}
	}
	
	return questions, nil
}

// parsePUMLResponse 解析PUML生成响应
func (c *OpenAIClient) parsePUMLResponse(content, projectID string, diagramType PUMLType) (*PUMLDiagram, error) {
	jsonContent := extractJSON(content)
	
	var result struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}
	
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	
	return &PUMLDiagram{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		Type:        diagramType,
		Title:       result.Title,
		Content:     result.Content,
		Description: result.Description,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// parseDocumentResponse 解析文档生成响应
func (c *OpenAIClient) parseDocumentResponse(content, projectID string) (*DevelopmentDocument, error) {
	jsonContent := extractJSON(content)
	
	var result DevelopmentDocument
	
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	
	result.ID = uuid.New().String()
	result.ProjectID = projectID
	result.Version = 1
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()
	
	return &result, nil
}

// buildProjectChatPrompt 构建项目对话的提示语
func (c *OpenAIClient) buildProjectChatPrompt(message, context string) string {
	return fmt.Sprintf(`你是一个专业的AI项目助手，专门帮助用户优化项目需求分析和开发细节。

项目上下文信息：
%s

用户问题：
%s

请根据项目上下文和用户问题，提供专业的回答和建议。如果用户的问题涉及需求分析的优化，请同时提供相关的建议。

请按照以下JSON格式回复：
{
  "message": "回答用户问题的详细内容",
  "should_update_analysis": false,
  "related_questions": ["相关问题1", "相关问题2"],
  "suggestions": ["建议1", "建议2"],
  "analysis_updates": {}
}

注意：
1. 回答要专业、准确、有针对性
2. 如果建议更新需求分析，将should_update_analysis设为true
3. 提供相关的后续问题帮助用户深入思考
4. 给出实用的建议来改进项目`, context, message)
}

// parseProjectChatResponse 解析项目对话响应
func (c *OpenAIClient) parseProjectChatResponse(content string) (*ProjectChatResponse, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var response ProjectChatResponse
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &response, nil
}

// extractJSON 从文本中提取JSON部分
func extractJSON(content string) string {
	// 寻找JSON开始和结束标记
	start := strings.Index(content, "{")
	if start == -1 {
		return content
	}
	
	// 从后往前找最后一个}
	end := strings.LastIndex(content, "}")
	if end == -1 || end <= start {
		return content
	}
	
	return content[start : end+1]
} 