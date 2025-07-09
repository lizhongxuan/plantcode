package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GeminiClient Google Gemini客户端实现
type GeminiClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// GeminiConfig Gemini配置
type GeminiConfig struct {
	APIKey  string
	BaseURL string // 可选，默认为Google AI Studio API
	Model   string // 可选，默认为gemini-1.5-flash
}

// NewGeminiClient 创建Gemini客户端
func NewGeminiClient(config GeminiConfig) *GeminiClient {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	
	model := config.Model
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &GeminiClient{
		apiKey:  config.APIKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GetProvider 返回AI服务提供商类型
func (c *GeminiClient) GetProvider() AIProvider {
	return ProviderGemini
}

// AnalyzeRequirement 分析业务需求
func (c *GeminiClient) AnalyzeRequirement(ctx context.Context, requirement string) (*RequirementAnalysis, error) {
	prompt := c.buildAnalysisPrompt(requirement)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	analysis, err := c.parseAnalysisResponse(response.Content, requirement)
	if err != nil {
		return nil, fmt.Errorf("解析需求分析结果失败: %w", err)
	}

	return analysis, nil
}

// GenerateQuestions 基于分析结果生成补充问题
func (c *GeminiClient) GenerateQuestions(ctx context.Context, analysis *RequirementAnalysis) ([]Question, error) {
	prompt := c.buildQuestionsPrompt(analysis)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	questions, err := c.parseQuestionsResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析问题生成结果失败: %w", err)
	}

	return questions, nil
}

// GeneratePUML 生成PUML图表代码
func (c *GeminiClient) GeneratePUML(ctx context.Context, analysis *RequirementAnalysis, diagramType PUMLType) (*PUMLDiagram, error) {
	prompt := c.buildPUMLPrompt(analysis, diagramType)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	diagram, err := c.parsePUMLResponse(response.Content, analysis.ProjectID, diagramType)
	if err != nil {
		return nil, fmt.Errorf("解析PUML生成结果失败: %w", err)
	}

	return diagram, nil
}

// GenerateDocument 生成开发文档
func (c *GeminiClient) GenerateDocument(ctx context.Context, analysis *RequirementAnalysis) (*DevelopmentDocument, error) {
	prompt := c.buildDocumentPrompt(analysis)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	document, err := c.parseDocumentResponse(response.Content, analysis.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("解析文档生成结果失败: %w", err)
	}

	return document, nil
}

// ProjectChat 项目上下文AI对话
func (c *GeminiClient) ProjectChat(ctx context.Context, message, context string) (*ProjectChatResponse, error) {
	prompt := c.buildProjectChatPrompt(message, context)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	chatResponse, err := c.parseProjectChatResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析项目对话结果失败: %w", err)
	}

	return chatResponse, nil
}

// GenerateStageSpecificDocument 生成特定阶段的文档
func (c *GeminiClient) GenerateStageSpecificDocument(ctx context.Context, analysis *RequirementAnalysis, documentType string) (*DevelopmentDocument, error) {
	prompt := c.buildStageDocumentPrompt(analysis, documentType)
	
	response, err := c.callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用Gemini API失败: %w", err)
	}

	document, err := c.parseStageDocumentResponse(response.Content, analysis.ProjectID, documentType)
	if err != nil {
		return nil, fmt.Errorf("解析文档生成响应失败: %w", err)
	}

	return document, nil
}

// CallGemini 公开的Gemini调用方法
func (c *GeminiClient) CallGemini(ctx context.Context, prompt string) (*AIResponse, error) {
	return c.callGemini(ctx, prompt)
}

// callGemini 调用Gemini API
func (c *GeminiClient) callGemini(ctx context.Context, prompt string) (*AIResponse, error) {
	req := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.3,
			"maxOutputTokens": 2000,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("构建请求数据失败: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("Gemini API返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("解析Gemini响应失败: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("Gemini响应中没有生成内容")
	}

	return &AIResponse{
		ID:      uuid.New().String(),
		Content: geminiResp.Candidates[0].Content.Parts[0].Text,
		Usage: AIUsage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
		Model: c.model,
	}, nil
}

// buildAnalysisPrompt 构建需求分析的提示语
func (c *GeminiClient) buildAnalysisPrompt(requirement string) string {
	return fmt.Sprintf(`请将以下需求分析作为一个完整的软件项目，进行全面的项目架构和功能分析。

业务需求：
%s

请按照以下JSON格式返回完整的项目分析结果：
{
  "project_overview": {
    "project_name": "推荐的项目名称",
    "project_type": "web_application|mobile_app|desktop_app|api_service|other",
    "target_users": ["目标用户群体1", "目标用户群体2"],
    "core_value": "项目核心价值主张"
  },
  "system_architecture": {
    "architecture_pattern": "MVC|MVP|MVVM|微服务|单体应用|其他",
    "frontend_tech": ["推荐的前端技术栈"],
    "backend_tech": ["推荐的后端技术栈"],
    "database_tech": ["推荐的数据库技术"],
    "deployment_env": ["推荐的部署环境"],
    "external_services": ["需要集成的外部服务"]
  },
  "core_functions": [
    {
      "name": "功能模块名称",
      "description": "功能详细描述",
      "priority": "高|中|低",
      "complexity": "简单|中等|复杂",
      "sub_functions": ["子功能1", "子功能2"],
      "dependencies": ["依赖的其他功能模块"]
    }
  ],
  "user_roles": [
    {
      "name": "角色名称",
      "description": "角色描述",
      "permissions": ["权限1", "权限2"],
      "main_workflows": ["主要使用流程1", "主要使用流程2"]
    }
  ],
  "business_processes": [
    {
      "name": "业务流程名称",
      "description": "流程详细描述",
      "steps": [
        {
          "step_name": "步骤名称",
          "description": "步骤描述",
          "actor": "执行者",
          "inputs": ["输入1", "输入2"],
          "outputs": ["输出1", "输出2"],
          "business_rules": ["业务规则1", "业务规则2"]
        }
      ],
      "exception_handling": ["异常情况1", "异常情况2"],
      "performance_requirements": "性能要求描述"
    }
  ],
  "data_entities": [
    {
      "name": "实体名称",
      "description": "实体业务含义",
      "category": "核心实体|业务实体|配置实体|日志实体",
      "attributes": [
        {
          "name": "属性名",
          "type": "string|int|float|boolean|date|text|json",
          "required": true,
          "unique": false,
          "description": "属性业务含义",
          "constraints": ["约束条件1", "约束条件2"]
        }
      ],
      "relations": [
        {
          "target_entity": "目标实体名",
          "relation_type": "one-to-one|one-to-many|many-to-many",
          "description": "关系描述",
          "foreign_key": "外键字段名"
        }
      ],
      "indexes": ["需要建立索引的字段"],
      "business_rules": ["业务规则1", "业务规则2"]
    }
  ],
  "api_interfaces": [
    {
      "module": "功能模块",
      "endpoints": [
        {
          "method": "GET|POST|PUT|DELETE",
          "path": "/api/path",
          "description": "接口描述",
          "auth_required": true,
          "request_params": ["参数1", "参数2"],
          "response_format": "响应格式描述"
        }
      ]
    }
  ],
  "security_requirements": {
    "authentication": "认证方式描述",
    "authorization": "授权机制描述",
    "data_protection": ["数据保护措施1", "数据保护措施2"],
    "communication_security": "通信安全要求"
  },
  "performance_requirements": {
    "response_time": "响应时间要求",
    "concurrent_users": "并发用户数",
    "data_volume": "数据量要求",
    "availability": "可用性要求"
  },
  "development_phases": [
    {
      "phase_name": "阶段名称",
      "description": "阶段描述",
      "deliverables": ["交付物1", "交付物2"],
      "estimated_duration": "预估时间",
      "key_milestones": ["里程碑1", "里程碑2"]
    }
  ],
  "missing_info": [
    {
      "category": "business_rule|technical_spec|ui_ux|integration|other",
      "description": "缺失信息的具体描述",
      "impact": "对项目的影响",
      "priority": "高|中|低"
    }
  ],
  "completion_score": 0.7,
  "recommendations": [
    {
      "category": "技术选型|架构设计|开发流程|部署策略|其他",
      "recommendation": "具体建议内容",
      "reason": "建议理由"
    }
  ]
}

分析要求：
1. 将需求视为完整的软件项目进行全方位分析
2. 提供具体可执行的技术建议和架构方案
3. 识别所有核心业务流程和数据流
4. 提供详细的数据模型设计
5. 考虑安全性、性能、可扩展性等非功能性需求
6. 提供分阶段的开发计划
7. 完整度评分基于需求的详细程度和可实施性(0-1之间)
8. 重点关注项目的技术架构和实现路径`, requirement)
}

// buildQuestionsPrompt 构建问题生成的提示语
func (c *GeminiClient) buildQuestionsPrompt(analysis *RequirementAnalysis) string {
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

请确保问题具体、有针对性且容易理解。`, missingInfo)
}

// buildPUMLPrompt 构建PUML生成的提示语
func (c *GeminiClient) buildPUMLPrompt(analysis *RequirementAnalysis, diagramType PUMLType) string {
	var promptTemplate string
	
	switch diagramType {
	case PUMLTypeBusinessFlow:
		promptTemplate = `基于以下需求分析，生成业务流程PUML代码。

%s

请生成一个完整的业务流程活动图(Activity Diagram)，包括：
1. 主要业务流程
2. 决策节点
3. 并行处理
4. 异常处理流程

请严格按照以下JSON格式返回：
{
  "title": "业务流程图",
  "content": "@startuml\\n业务流程图\\n...\\n@enduml",
  "description": "图表说明"
}`

	case PUMLTypeArchitecture:
		promptTemplate = `基于以下需求分析，生成系统架构PUML代码。

%s

请生成一个完整的系统架构组件图(Component Diagram)，包括：
1. 系统主要组件
2. 组件间的依赖关系
3. 外部系统接口
4. 数据流向

请严格按照以下JSON格式返回：
{
  "title": "系统架构图",
  "content": "@startuml\\n系统架构图\\n...\\n@enduml",
  "description": "图表说明"
}`

	case PUMLTypeDataModel:
		promptTemplate = `基于以下需求分析，生成数据模型PUML代码。

%s

请生成一个完整的数据模型实体关系图(Entity Relationship Diagram)，包括：
1. 数据实体
2. 实体属性
3. 实体间关系
4. 主键和外键约束

请严格按照以下JSON格式返回：
{
  "title": "数据模型图",
  "content": "@startuml\\n数据模型图\\n...\\n@enduml",
  "description": "图表说明"
}`

	default:
		promptTemplate = `基于以下需求分析，生成PUML代码。

%s

请生成相应的PUML图表代码，严格按照以下JSON格式返回：
{
  "title": "图表标题",
  "content": "@startuml\\n图表内容\\n...\\n@enduml",
  "description": "图表说明"
}`
	}

	analysisContext := fmt.Sprintf(`
需求分析结果：
核心功能：%s
用户角色：%s
业务流程：%d个
数据实体：%d个
`, 
		strings.Join(analysis.CoreFunctions, ", "),
		strings.Join(analysis.Roles, ", "),
		len(analysis.BusinessProcesses),
		len(analysis.DataEntities))

	return fmt.Sprintf(promptTemplate, analysisContext)
}

// buildDocumentPrompt 构建文档生成的提示语
func (c *GeminiClient) buildDocumentPrompt(analysis *RequirementAnalysis) string {
	return fmt.Sprintf(`基于以下需求分析，生成完整的开发文档。

核心功能：%s
用户角色：%s
业务流程数：%d
数据实体数：%d

请按照以下JSON格式返回开发文档：
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
    "resources": "资源需求"
  },
  "tech_stack": {
    "backend": {
      "recommended": "推荐技术",
      "alternatives": ["备选方案"],
      "reason": "选择理由"
    },
    "frontend": {
      "recommended": "推荐技术",
      "alternatives": ["备选方案"],
      "reason": "选择理由"
    },
    "database": {
      "recommended": "推荐数据库",
      "alternatives": ["备选方案"],
      "reason": "选择理由"
    }
  },
  "database_design": {
    "tables": [
      {
        "name": "表名",
        "comment": "表说明",
        "columns": [
          {
            "name": "列名",
            "type": "数据类型",
            "length": 255,
            "nullable": false,
            "default": "",
            "comment": "列说明",
            "primary_key": true
          }
        ]
      }
    ]
  },
  "api_design": [
    {
      "path": "/api/endpoint",
      "method": "POST",
      "summary": "接口概述",
      "description": "接口详细描述",
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
}`, 
		strings.Join(analysis.CoreFunctions, ", "),
		strings.Join(analysis.Roles, ", "),
		len(analysis.BusinessProcesses),
		len(analysis.DataEntities))
}

// buildProjectChatPrompt 构建项目对话的提示语
func (c *GeminiClient) buildProjectChatPrompt(message, context string) string {
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

// parseAnalysisResponse 解析需求分析响应
func (c *GeminiClient) parseAnalysisResponse(content, originalText string) (*RequirementAnalysis, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var rawAnalysis struct {
		ProjectOverview struct {
			ProjectName  string   `json:"project_name"`
			ProjectType  string   `json:"project_type"`
			TargetUsers  []string `json:"target_users"`
			CoreValue    string   `json:"core_value"`
		} `json:"project_overview"`
		SystemArchitecture struct {
			ArchitecturePattern string   `json:"architecture_pattern"`
			FrontendTech        []string `json:"frontend_tech"`
			BackendTech         []string `json:"backend_tech"`
			DatabaseTech        []string `json:"database_tech"`
			DeploymentEnv       []string `json:"deployment_env"`
			ExternalServices    []string `json:"external_services"`
		} `json:"system_architecture"`
		CoreFunctions []struct {
			Name         string   `json:"name"`
			Description  string   `json:"description"`
			Priority     string   `json:"priority"`
			Complexity   string   `json:"complexity"`
			SubFunctions []string `json:"sub_functions"`
			Dependencies []string `json:"dependencies"`
		} `json:"core_functions"`
		UserRoles []struct {
			Name          string   `json:"name"`
			Description   string   `json:"description"`
			Permissions   []string `json:"permissions"`
			MainWorkflows []string `json:"main_workflows"`
		} `json:"user_roles"`
		BusinessProcesses []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Steps       []struct {
				StepName      string   `json:"step_name"`
				Description   string   `json:"description"`
				Actor         string   `json:"actor"`
				Inputs        []string `json:"inputs"`
				Outputs       []string `json:"outputs"`
				BusinessRules []string `json:"business_rules"`
			} `json:"steps"`
			ExceptionHandling        []string `json:"exception_handling"`
			PerformanceRequirements  string   `json:"performance_requirements"`
		} `json:"business_processes"`
		DataEntities []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Attributes  []struct {
				Name        string   `json:"name"`
				Type        string   `json:"type"`
				Required    bool     `json:"required"`
				Unique      bool     `json:"unique"`
				Description string   `json:"description"`
				Constraints []string `json:"constraints"`
			} `json:"attributes"`
			Relations []struct {
				TargetEntity string `json:"target_entity"`
				RelationType string `json:"relation_type"`
				Description  string `json:"description"`
				ForeignKey   string `json:"foreign_key"`
			} `json:"relations"`
			Indexes       []string `json:"indexes"`
			BusinessRules []string `json:"business_rules"`
		} `json:"data_entities"`
		ApiInterfaces []struct {
			Module    string `json:"module"`
			Endpoints []struct {
				Method         string   `json:"method"`
				Path           string   `json:"path"`
				Description    string   `json:"description"`
				AuthRequired   bool     `json:"auth_required"`
				RequestParams  []string `json:"request_params"`
				ResponseFormat string   `json:"response_format"`
			} `json:"endpoints"`
		} `json:"api_interfaces"`
		SecurityRequirements struct {
			Authentication      string   `json:"authentication"`
			Authorization       string   `json:"authorization"`
			DataProtection      []string `json:"data_protection"`
			CommunicationSecurity string `json:"communication_security"`
		} `json:"security_requirements"`
		PerformanceRequirements struct {
			ResponseTime     string `json:"response_time"`
			ConcurrentUsers  string `json:"concurrent_users"`
			DataVolume       string `json:"data_volume"`
			Availability     string `json:"availability"`
		} `json:"performance_requirements"`
		DevelopmentPhases []struct {
			PhaseName        string   `json:"phase_name"`
			Description      string   `json:"description"`
			Deliverables     []string `json:"deliverables"`
			EstimatedDuration string  `json:"estimated_duration"`
			KeyMilestones    []string `json:"key_milestones"`
		} `json:"development_phases"`
		MissingInfo []struct {
			Category    string `json:"category"`
			Description string `json:"description"`
			Impact      string `json:"impact"`
			Priority    string `json:"priority"`
		} `json:"missing_info"`
		CompletionScore float64 `json:"completion_score"`
		Recommendations []struct {
			Category       string `json:"category"`
			Recommendation string `json:"recommendation"`
			Reason         string `json:"reason"`
		} `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &rawAnalysis); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	analysis := &RequirementAnalysis{
		ID:              uuid.New().String(),
		OriginalText:    originalText,
		CompletionScore: rawAnalysis.CompletionScore,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 提取核心功能列表（向后兼容）
	for _, cf := range rawAnalysis.CoreFunctions {
		analysis.CoreFunctions = append(analysis.CoreFunctions, cf.Name)
	}

	// 提取用户角色列表（向后兼容）
	for _, ur := range rawAnalysis.UserRoles {
		analysis.Roles = append(analysis.Roles, ur.Name)
	}

	// 提取缺失信息列表（向后兼容）
	for _, mi := range rawAnalysis.MissingInfo {
		analysis.MissingInfo = append(analysis.MissingInfo, mi.Description)
	}

	// 转换业务流程
	for _, bp := range rawAnalysis.BusinessProcesses {
		process := BusinessProcess{
			Name:        bp.Name,
			Description: bp.Description,
		}
		
		// 提取步骤名称（向后兼容）
		for _, step := range bp.Steps {
			process.Steps = append(process.Steps, step.StepName)
		}
		
		// 提取参与者（向后兼容）
		actorMap := make(map[string]bool)
		for _, step := range bp.Steps {
			if step.Actor != "" && !actorMap[step.Actor] {
				process.Actors = append(process.Actors, step.Actor)
				actorMap[step.Actor] = true
			}
		}

		analysis.BusinessProcesses = append(analysis.BusinessProcesses, process)
	}

	// 转换数据实体
	for _, de := range rawAnalysis.DataEntities {
		entity := DataEntity{
			Name:        de.Name,
			Description: de.Description,
		}

		for _, attr := range de.Attributes {
			entity.Attributes = append(entity.Attributes, EntityAttribute{
				Name:        attr.Name,
				Type:        attr.Type,
				Required:    attr.Required,
				Description: attr.Description,
			})
		}

		for _, rel := range de.Relations {
			entity.Relations = append(entity.Relations, EntityRelation{
				TargetEntity: rel.TargetEntity,
				RelationType: rel.RelationType,
				Description:  rel.Description,
			})
		}

		analysis.DataEntities = append(analysis.DataEntities, entity)
	}

	return analysis, nil
}

// parseQuestionsResponse 解析问题生成响应
func (c *GeminiClient) parseQuestionsResponse(content string) ([]Question, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var response struct {
		Questions []struct {
			Category   string   `json:"category"`
			Content    string   `json:"content"`
			Options    []string `json:"options"`
			Priority   int      `json:"priority"`
			TargetInfo string   `json:"target_info"`
		} `json:"questions"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	var questions []Question
	for _, q := range response.Questions {
		questions = append(questions, Question{
			ID:          uuid.New().String(),
			Category:    q.Category,
			Content:     q.Content,
			Options:     q.Options,
			Priority:    q.Priority,
			TargetInfo:  q.TargetInfo,
		})
	}

	return questions, nil
}

// parsePUMLResponse 解析PUML生成响应
func (c *GeminiClient) parsePUMLResponse(content, projectID string, diagramType PUMLType) (*PUMLDiagram, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var response struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	return &PUMLDiagram{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		Type:        diagramType,
		Title:       response.Title,
		Content:     response.Content,
		Description: response.Description,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// parseDocumentResponse 解析文档生成响应
func (c *GeminiClient) parseDocumentResponse(content, projectID string) (*DevelopmentDocument, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var response DevelopmentDocument
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	response.ID = uuid.New().String()
	response.ProjectID = projectID
	response.Version = 1
	response.CreatedAt = time.Now()
	response.UpdatedAt = time.Now()

	return &response, nil
}

// parseProjectChatResponse 解析项目对话响应
func (c *GeminiClient) parseProjectChatResponse(content string) (*ProjectChatResponse, error) {
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

// buildStageDocumentPrompt 构建分阶段文档生成的提示语
func (c *GeminiClient) buildStageDocumentPrompt(analysis *RequirementAnalysis, documentType string) string {
	baseContext := fmt.Sprintf(`
基于以下需求分析生成专业的项目文档：

项目信息：
- 项目ID：%s
- 核心功能：%s
- 用户角色：%s
- 业务流程数：%d
- 数据实体数：%d

原始需求：
%s
`, 
		analysis.ProjectID,
		strings.Join(analysis.CoreFunctions, ", "),
		strings.Join(analysis.Roles, ", "),
		len(analysis.BusinessProcesses),
		len(analysis.DataEntities),
		analysis.OriginalText)

	var promptTemplate string
	
	switch documentType {
	case "requirements":
		promptTemplate = baseContext + `
请生成详细的【项目需求文档】，包含以下内容：

1. 项目概述
   - 项目背景和目标
   - 目标用户群体
   - 项目价值主张

2. 功能需求
   - 核心功能模块详细说明
   - 功能优先级排序
   - 用户故事和用例场景

3. 非功能需求
   - 性能要求
   - 安全性要求
   - 可用性要求
   - 兼容性要求

4. 系统约束
   - 技术约束
   - 业务约束
   - 时间约束

5. 风险评估
   - 技术风险
   - 业务风险
   - 缓解策略

请按照以下JSON格式返回：
{
  "title": "项目需求文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "requirements",
  "version": "1.0"
}`

	case "technical_spec":
		promptTemplate = baseContext + `
请生成详细的【技术规范文档】，包含以下内容：

1. 技术架构
   - 系统架构设计
   - 技术栈选择
   - 架构模式说明

2. 系统设计
   - 模块划分
   - 接口设计
   - 数据流设计

3. 开发规范
   - 编码标准
   - 命名规范
   - 文档规范

4. 部署架构
   - 环境配置
   - 部署流程
   - 监控方案

请按照以下JSON格式返回：
{
  "title": "技术规范文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "technical_spec",
  "version": "1.0"
}`

	case "api_design":
		promptTemplate = baseContext + `
请生成详细的【API接口设计文档】，包含以下内容：

1. API概述
   - 接口设计原则
   - 认证机制
   - 版本控制

2. 接口规范
   - RESTful API设计
   - 请求/响应格式
   - 错误码定义

3. 接口列表
   - 按模块分组的接口
   - 每个接口的详细说明
   - 请求参数和响应示例

4. 数据模型
   - 实体定义
   - 关系说明
   - 验证规则

请按照以下JSON格式返回：
{
  "title": "API接口设计文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "api_design",
  "version": "1.0"
}`

	case "database_design":
		promptTemplate = baseContext + `
请生成详细的【数据库设计文档】，包含以下内容：

1. 数据库概述
   - 数据库选型
   - 设计原则
   - 数据架构

2. 表结构设计
   - 表定义
   - 字段说明
   - 索引设计

3. 关系设计
   - 实体关系图
   - 外键约束
   - 关系说明

4. 数据字典
   - 表详细说明
   - 字段类型定义
   - 业务规则

请按照以下JSON格式返回：
{
  "title": "数据库设计文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "database_design",
  "version": "1.0"
}`

	case "development_process":
		promptTemplate = baseContext + `
请生成详细的【开发流程文档】，包含以下内容：

1. 开发流程
   - 开发阶段划分
   - 里程碑定义
   - 交付物清单

2. 团队协作
   - 角色职责
   - 协作流程
   - 沟通机制

3. 质量保证
   - 代码审查
   - 测试策略
   - 质量标准

4. 项目管理
   - 进度跟踪
   - 风险管理
   - 变更管理

请按照以下JSON格式返回：
{
  "title": "开发流程文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "development_process",
  "version": "1.0"
}`

	case "test_cases":
		promptTemplate = baseContext + `
请生成详细的【测试用例文档】，包含以下内容：

1. 测试策略
   - 测试范围
   - 测试类型
   - 测试环境

2. 功能测试用例
   - 正常流程测试
   - 异常流程测试
   - 边界值测试

3. 性能测试用例
   - 负载测试
   - 压力测试
   - 稳定性测试

4. 安全测试用例
   - 认证测试
   - 授权测试
   - 数据安全测试

请按照以下JSON格式返回：
{
  "title": "测试用例文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "test_cases",
  "version": "1.0"
}`

	case "deployment":
		promptTemplate = baseContext + `
请生成详细的【部署文档】，包含以下内容：

1. 部署架构
   - 环境规划
   - 服务器配置
   - 网络架构

2. 部署流程
   - 构建流程
   - 部署步骤
   - 回滚策略

3. 运维监控
   - 监控方案
   - 日志管理
   - 性能监控

4. 维护手册
   - 日常维护
   - 故障处理
   - 备份恢复

请按照以下JSON格式返回：
{
  "title": "部署文档",
  "content": "详细的markdown格式文档内容",
  "document_type": "deployment",
  "version": "1.0"
}`

	default:
		promptTemplate = baseContext + `
请生成对应类型的技术文档，按照以下JSON格式返回：
{
  "title": "文档标题",
  "content": "详细的markdown格式文档内容",
  "document_type": "` + documentType + `",
  "version": "1.0"
}`
	}
	
	return promptTemplate
}

// parseStageDocumentResponse 解析分阶段文档生成响应
func (c *GeminiClient) parseStageDocumentResponse(content, projectID, documentType string) (*DevelopmentDocument, error) {
	jsonContent := extractJSON(content)
	if jsonContent == "" {
		return nil, fmt.Errorf("响应中没有找到JSON格式的内容")
	}

	var response struct {
		Title        string `json:"title"`
		Content      string `json:"content"`
		DocumentType string `json:"document_type"`
		Version      string `json:"version"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	// 解析版本字符串为整数
	version := 1
	if response.Version != "" {
		if v, err := strconv.Atoi(strings.TrimPrefix(response.Version, "v")); err == nil {
			version = v
		}
	}

	// 创建简单的功能模块（基于文档类型）
	functionModule := FunctionModule{
		Name:         response.Title,
		Description:  "基于AI生成的" + documentType + "文档",
		SubModules:   []string{},
		Dependencies: []string{},
		Priority:     1,
		Complexity:   "medium",
		EstimatedHours: 8,
	}

	// 创建简单的开发计划
	developmentPlan := DevelopmentPlan{
		Phases: []DevelopmentPhase{
			{
				Name:        "文档阶段",
				Description: "AI生成的" + documentType + "文档",
				Tasks:       []string{"文档生成", "文档审核"},
				Duration:    "1天",
				Dependencies: []string{},
			},
		},
		Duration:  "1天",
		Resources: "AI自动生成",
	}

	// 创建空的技术栈推荐
	techStack := TechStackRecommendation{
		Backend:      TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
		Frontend:     TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
		Database:     TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
		Cache:        TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
		MessageQueue: TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
		Deployment:   TechChoice{Recommended: "未指定", Alternatives: []string{}, Reason: "需要进一步分析"},
	}

	document := &DevelopmentDocument{
		ID:               uuid.New().String(),
		ProjectID:        projectID,
		FunctionModules:  []FunctionModule{functionModule},
		DevelopmentPlan:  developmentPlan,
		TechStack:        techStack,
		DatabaseDesign:   DatabaseDesign{Tables: []TableDesign{}, Indexes: []IndexDesign{}, Relations: []RelationDesign{}},
		APIDesign:        []APIEndpoint{},
		Version:          version,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return document, nil
} 