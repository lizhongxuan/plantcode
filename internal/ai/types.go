package ai

import (
	"context"
	"time"
)

// AIProvider 定义支持的AI服务提供商
type AIProvider string

const (
	ProviderOpenAI AIProvider = "openai"
	ProviderClaude AIProvider = "claude"
	ProviderGemini AIProvider = "gemini"
)

// AIClient 定义AI客户端的统一接口
type AIClient interface {
	// AnalyzeRequirement 分析业务需求
	AnalyzeRequirement(ctx context.Context, requirement string) (*RequirementAnalysis, error)
	
	// GenerateQuestions 基于分析结果生成补充问题
	GenerateQuestions(ctx context.Context, analysis *RequirementAnalysis) ([]Question, error)
	
	// GeneratePUML 生成PUML图表代码
	GeneratePUML(ctx context.Context, analysis *RequirementAnalysis, diagramType PUMLType) (*PUMLDiagram, error)
	
	// GenerateDocument 生成开发文档
	GenerateDocument(ctx context.Context, analysis *RequirementAnalysis) (*DevelopmentDocument, error)
	
	// ProjectChat 项目上下文AI对话
	ProjectChat(ctx context.Context, message, context string) (*ProjectChatResponse, error)
	
	// GetProvider 返回AI服务提供商类型
	GetProvider() AIProvider
}

// RequirementAnalysis 需求分析结果
type RequirementAnalysis struct {
	ID                string            `json:"id"`
	ProjectID         string            `json:"project_id"`
	OriginalText      string            `json:"original_text"`
	CoreFunctions     []string          `json:"core_functions"`     // 核心功能
	Roles             []string          `json:"roles"`              // 参与角色
	BusinessProcesses []BusinessProcess `json:"business_processes"` // 业务流程
	DataEntities      []DataEntity      `json:"data_entities"`      // 数据实体
	MissingInfo       []string          `json:"missing_info"`       // 缺失信息
	CompletionScore   float64           `json:"completion_score"`   // 完整度评分 (0-1)
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// BusinessProcess 业务流程
type BusinessProcess struct {
	Name        string   `json:"name"`        // 流程名称
	Description string   `json:"description"` // 流程描述
	Steps       []string `json:"steps"`       // 流程步骤
	Actors      []string `json:"actors"`      // 参与者
}

// DataEntity 数据实体
type DataEntity struct {
	Name        string              `json:"name"`        // 实体名称
	Description string              `json:"description"` // 实体描述
	Attributes  []EntityAttribute   `json:"attributes"`  // 属性列表
	Relations   []EntityRelation    `json:"relations"`   // 关系列表
}

// EntityAttribute 实体属性
type EntityAttribute struct {
	Name        string `json:"name"`        // 属性名
	Type        string `json:"type"`        // 数据类型
	Required    bool   `json:"required"`    // 是否必需
	Description string `json:"description"` // 属性描述
}

// EntityRelation 实体关系
type EntityRelation struct {
	TargetEntity string `json:"target_entity"` // 目标实体
	RelationType string `json:"relation_type"` // 关系类型 (one-to-one, one-to-many, many-to-many)
	Description  string `json:"description"`   // 关系描述
}

// Question 补充问题
type Question struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`    // 问题分类 (business_rule, exception_handling, etc.)
	Content     string   `json:"content"`     // 问题内容
	Options     []string `json:"options"`     // 可选答案（如果有）
	Priority    int      `json:"priority"`    // 优先级 (1-5)
	TargetInfo  string   `json:"target_info"` // 目标获取的信息类型
}

// PUMLType PUML图表类型
type PUMLType string

const (
	PUMLTypeBusinessFlow PUMLType = "business_flow" // 业务流程图
	PUMLTypeArchitecture PUMLType = "architecture"  // 系统架构图
	PUMLTypeSequence     PUMLType = "sequence"      // 序列图
	PUMLTypeClass        PUMLType = "class"         // 类图
	PUMLTypeDataModel    PUMLType = "data_model"    // 数据模型图 (ER图)
)

// PUMLDiagram PUML图表
type PUMLDiagram struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Type        PUMLType  `json:"type"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`     // PUML代码
	Description string    `json:"description"` // 图表说明
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DevelopmentDocument 开发文档
type DevelopmentDocument struct {
	ID               string                `json:"id"`
	ProjectID        string                `json:"project_id"`
	FunctionModules  []FunctionModule      `json:"function_modules"`  // 功能模块
	DevelopmentPlan  DevelopmentPlan       `json:"development_plan"`  // 开发计划
	TechStack        TechStackRecommendation `json:"tech_stack"`        // 技术选型
	DatabaseDesign   DatabaseDesign        `json:"database_design"`   // 数据库设计
	APIDesign        []APIEndpoint         `json:"api_design"`        // API设计
	Version          int                   `json:"version"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

// FunctionModule 功能模块
type FunctionModule struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	SubModules   []string `json:"sub_modules"`   // 子模块
	Dependencies []string `json:"dependencies"` // 依赖的其他模块
	Priority     int      `json:"priority"`     // 开发优先级
	Complexity   string   `json:"complexity"`   // 复杂度 (low, medium, high)
	EstimatedHours int    `json:"estimated_hours"` // 预估工时
}

// DevelopmentPlan 开发计划
type DevelopmentPlan struct {
	Phases    []DevelopmentPhase `json:"phases"`
	Duration  string             `json:"duration"`  // 总开发周期
	Resources string             `json:"resources"` // 资源需求
}

// DevelopmentPhase 开发阶段
type DevelopmentPhase struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tasks       []string `json:"tasks"`
	Duration    string   `json:"duration"`
	Dependencies []string `json:"dependencies"` // 依赖的前置阶段
}

// TechStackRecommendation 技术选型建议
type TechStackRecommendation struct {
	Backend      TechChoice `json:"backend"`
	Frontend     TechChoice `json:"frontend"`
	Database     TechChoice `json:"database"`
	Cache        TechChoice `json:"cache"`
	MessageQueue TechChoice `json:"message_queue"`
	Deployment   TechChoice `json:"deployment"`
}

// TechChoice 技术选择
type TechChoice struct {
	Recommended string   `json:"recommended"` // 推荐选择
	Alternatives []string `json:"alternatives"` // 备选方案
	Reason      string   `json:"reason"`      // 选择理由
}

// DatabaseDesign 数据库设计
type DatabaseDesign struct {
	Tables []TableDesign `json:"tables"`
	Indexes []IndexDesign `json:"indexes"`
	Relations []RelationDesign `json:"relations"`
}

// TableDesign 表设计
type TableDesign struct {
	Name        string        `json:"name"`
	Comment     string        `json:"comment"`
	Columns     []ColumnDesign `json:"columns"`
}

// ColumnDesign 列设计
type ColumnDesign struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Length     int    `json:"length,omitempty"`
	Nullable   bool   `json:"nullable"`
	Default    string `json:"default,omitempty"`
	Comment    string `json:"comment"`
	PrimaryKey bool   `json:"primary_key"`
	AutoIncrement bool `json:"auto_increment"`
}

// IndexDesign 索引设计
type IndexDesign struct {
	Name    string   `json:"name"`
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// RelationDesign 关系设计
type RelationDesign struct {
	Name           string `json:"name"`
	FromTable      string `json:"from_table"`
	FromColumn     string `json:"from_column"`
	ToTable        string `json:"to_table"`
	ToColumn       string `json:"to_column"`
	OnDelete       string `json:"on_delete"` // CASCADE, SET NULL, RESTRICT
	OnUpdate       string `json:"on_update"` // CASCADE, SET NULL, RESTRICT
}

// APIEndpoint API端点设计
type APIEndpoint struct {
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Parameters  []APIParameter    `json:"parameters"`
	RequestBody *APIRequestBody   `json:"request_body,omitempty"`
	Responses   map[string]APIResponse `json:"responses"`
	Tags        []string          `json:"tags"`
}

// APIParameter API参数
type APIParameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // query, path, header
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// APIRequestBody API请求体
type APIRequestBody struct {
	Description string     `json:"description"`
	Schema      APISchema  `json:"schema"`
	Required    bool       `json:"required"`
}

// APIResponse API响应
type APIResponse struct {
	Description string    `json:"description"`
	Schema      APISchema `json:"schema"`
}

// APISchema API数据架构
type APISchema struct {
	Type       string               `json:"type"`
	Properties map[string]APIProperty `json:"properties,omitempty"`
	Example    interface{}          `json:"example,omitempty"`
}

// APIProperty API属性
type APIProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Example     interface{} `json:"example,omitempty"`
}

// AIRequest AI请求
type AIRequest struct {
	Provider AIProvider  `json:"provider"`
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	MaxTokens int        `json:"max_tokens"`
	Temperature float64  `json:"temperature"`
}

// AIMessage AI消息
type AIMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// AIResponse AI响应
type AIResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Usage   AIUsage `json:"usage"`
	Model   string `json:"model"`
}

// AIUsage AI使用量统计
type AIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ProjectChatResponse 项目对话响应
type ProjectChatResponse struct {
	Message              string   `json:"message"`
	ShouldUpdateAnalysis bool     `json:"should_update_analysis"`
	RelatedQuestions     []string `json:"related_questions"`
	Suggestions          []string `json:"suggestions"`
	AnalysisUpdates      map[string]interface{} `json:"analysis_updates,omitempty"`
} 