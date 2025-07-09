package model

import (
	"time"

	"github.com/google/uuid"
)

// User 用户模型
type User struct {
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Username    string     `json:"username" db:"username"`
	Email       string     `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // 不返回给前端
	FullName    string     `json:"full_name" db:"full_name"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin   *time.Time `json:"last_login" db:"last_login"`
	Status      string     `json:"status" db:"status"`
	Preferences string     `json:"preferences" db:"preferences"` // JSON字符串
}

// Project 项目模型
type Project struct {
	ProjectID             uuid.UUID `json:"project_id" db:"project_id"`
	UserID                uuid.UUID `json:"user_id" db:"user_id"`
	ProjectName           string    `json:"project_name" db:"project_name"`
	Description           string    `json:"description" db:"description"`
	ProjectType           string    `json:"project_type" db:"project_type"`
	Status                string    `json:"status" db:"status"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	CompletionPercentage  int       `json:"completion_percentage" db:"completion_percentage"`
	Settings              string    `json:"settings" db:"settings"` // JSON字符串
}

// Requirement 需求分析模型
type Requirement struct {
	RequirementID         uuid.UUID `json:"requirement_id" db:"requirement_id"`
	ProjectID             uuid.UUID `json:"project_id" db:"project_id"`
	RawRequirement        string    `json:"raw_requirement" db:"raw_requirement"`
	StructuredRequirement string    `json:"structured_requirement" db:"structured_requirement"` // JSON
	CompletenessScore     float64   `json:"completeness_score" db:"completeness_score"`
	AnalysisStatus        string    `json:"analysis_status" db:"analysis_status"`
	MissingInfoTypes      string    `json:"missing_info_types" db:"missing_info_types"` // JSON
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// ChatSession 对话会话模型
type ChatSession struct {
	SessionID   uuid.UUID  `json:"session_id" db:"session_id"`
	ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	SessionType string     `json:"session_type" db:"session_type"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	EndedAt     *time.Time `json:"ended_at" db:"ended_at"`
	Status      string     `json:"status" db:"status"`
	Context     string     `json:"context" db:"context"` // JSON字符串
}

// ChatMessage 对话消息模型
type ChatMessage struct {
	MessageID      uuid.UUID `json:"message_id" db:"message_id"`
	SessionID      uuid.UUID `json:"session_id" db:"session_id"`
	SenderType     string    `json:"sender_type" db:"sender_type"`     // user, ai, system
	MessageContent string    `json:"message_content" db:"message_content"`
	MessageType    string    `json:"message_type" db:"message_type"`   // text, question, answer
	Metadata       string    `json:"metadata" db:"metadata"`           // JSON字符串
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	Processed      bool      `json:"processed" db:"processed"`
}

// Question 补充问题模型
type Question struct {
	QuestionID       uuid.UUID  `json:"question_id" db:"question_id"`
	RequirementID    uuid.UUID  `json:"requirement_id" db:"requirement_id"`
	QuestionText     string     `json:"question_text" db:"question_text"`
	QuestionCategory string     `json:"question_category" db:"question_category"`
	PriorityLevel    int        `json:"priority_level" db:"priority_level"`
	AnswerText       string     `json:"answer_text" db:"answer_text"`
	AnswerStatus     string     `json:"answer_status" db:"answer_status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	AnsweredAt       *time.Time `json:"answered_at" db:"answered_at"`
}

// PUMLDiagram PUML图表模型
type PUMLDiagram struct {
	DiagramID          uuid.UUID `json:"diagram_id" db:"diagram_id"`
	ProjectID          uuid.UUID `json:"project_id" db:"project_id"`
	DiagramType        string    `json:"diagram_type" db:"diagram_type"` // business_flow, architecture, data_model
	DiagramName        string    `json:"diagram_name" db:"diagram_name"`
	PUMLContent        string    `json:"puml_content" db:"puml_content"`
	RenderedURL        string    `json:"rendered_url" db:"rendered_url"`
	Version            int       `json:"version" db:"version"`
	Stage              int       `json:"stage" db:"stage"`         // 新增：所属阶段 1,2,3
	TaskID             *uuid.UUID `json:"task_id,omitempty" db:"task_id"` // 新增：关联的异步任务ID
	IsValidated        bool      `json:"is_validated" db:"is_validated"`
	ValidationFeedback string    `json:"validation_feedback" db:"validation_feedback"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// BusinessModule 业务模块模型
type BusinessModule struct {
	ModuleID        uuid.UUID `json:"module_id" db:"module_id"`
	ProjectID       uuid.UUID `json:"project_id" db:"project_id"`
	ModuleName      string    `json:"module_name" db:"module_name"`
	Description     string    `json:"description" db:"description"`
	ModuleType      string    `json:"module_type" db:"module_type"`
	ComplexityLevel string    `json:"complexity_level" db:"complexity_level"`
	BusinessLogic   string    `json:"business_logic" db:"business_logic"` // JSON
	Interfaces      string    `json:"interfaces" db:"interfaces"`         // JSON
	Dependencies    string    `json:"dependencies" db:"dependencies"`     // JSON
	IsReusable      bool      `json:"is_reusable" db:"is_reusable"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// CommonModule 通用模块库模型
type CommonModule struct {
	CommonModuleID  uuid.UUID `json:"common_module_id" db:"common_module_id"`
	ModuleName      string    `json:"module_name" db:"module_name"`
	Category        string    `json:"category" db:"category"`
	Description     string    `json:"description" db:"description"`
	Functionality   string    `json:"functionality" db:"functionality"`     // JSON
	InterfaceSpec   string    `json:"interface_spec" db:"interface_spec"`   // JSON
	CodeTemplate    string    `json:"code_template" db:"code_template"`
	UsageExamples   string    `json:"usage_examples" db:"usage_examples"`   // JSON
	Version         string    `json:"version" db:"version"`
	DownloadsCount  int       `json:"downloads_count" db:"downloads_count"`
	Rating          float64   `json:"rating" db:"rating"`
	Tags            string    `json:"tags" db:"tags"`                       // JSON
	CreatedBy       uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Document 生成文档模型
type Document struct {
	DocumentID   uuid.UUID `json:"document_id" db:"document_id"`
	ProjectID    uuid.UUID `json:"project_id" db:"project_id"`
	DocumentType string    `json:"document_type" db:"document_type"`
	DocumentName string    `json:"document_name" db:"document_name"`
	Content      string    `json:"content" db:"content"`
	Format       string    `json:"format" db:"format"`
	FilePath     string    `json:"file_path" db:"file_path"`
	Version      int       `json:"version" db:"version"`
	Stage        int       `json:"stage" db:"stage"`         // 新增：所属阶段 1,2,3
	TaskID       *uuid.UUID `json:"task_id,omitempty" db:"task_id"` // 新增：关联的异步任务ID
	GeneratedAt  time.Time `json:"generated_at" db:"generated_at"`
	IsFinal      bool      `json:"is_final" db:"is_final"`
}

// 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusSuspended = "suspended"
)

// 项目状态常量
const (
	ProjectStatusDraft      = "draft"
	ProjectStatusAnalyzing  = "analyzing"
	ProjectStatusCompleted  = "completed"
	ProjectStatusArchived   = "archived"
)

// 对话会话类型常量
const (
	SessionTypeRequirementAnalysis = "requirement_analysis"
	SessionTypeDocumentReview      = "document_review"
	SessionTypePUMLEditing         = "puml_editing"
)

// 对话消息发送者类型常量
const (
	SenderTypeUser   = "user"
	SenderTypeAI     = "ai"
	SenderTypeSystem = "system"
)

// 消息类型常量
const (
	MessageTypeText     = "text"
	MessageTypeQuestion = "question"
	MessageTypeAnswer   = "answer"
	MessageTypeCommand  = "command"
)

// PUML图表类型常量
const (
	DiagramTypeBusinessFlow = "business_flow"
	DiagramTypeArchitecture = "architecture"
	DiagramTypeDataModel    = "data_model"
	DiagramTypeSequence     = "sequence"
)

// 分析状态常量
const (
	AnalysisStatusPending           = "pending"
	AnalysisStatusInProgress        = "in_progress"
	AnalysisStatusQuestionsGenerated = "questions_generated"
	AnalysisStatusCompleted         = "completed"
	AnalysisStatusFailed            = "failed"
)

// 问题回答状态常量
const (
	AnswerStatusPending   = "pending"
	AnswerStatusAnswered  = "answered"
	AnswerStatusSkipped   = "skipped"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	ProjectName string `json:"project_name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=1000"`
	ProjectType string `json:"project_type" validate:"required"`
}

// AnalyzeRequirementRequest 分析需求请求
type AnalyzeRequirementRequest struct {
	RawRequirement string `json:"raw_requirement" validate:"required,min=10"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	MessageContent string `json:"message_content" validate:"required"`
	MessageType    string `json:"message_type" validate:"required"`
}

// AnswerQuestionRequest 回答问题请求
type AnswerQuestionRequest struct {
	QuestionID uuid.UUID `json:"question_id" validate:"required"`
	Answer     string    `json:"answer" validate:"required"`
}

// ===== 第二阶段新增模型和请求类型 =====

// AIAnalysisRequest AI分析请求
type AIAnalysisRequest struct {
	ProjectID   uuid.UUID `json:"project_id" validate:"required"`
	Requirement string    `json:"requirement" validate:"required,min=10"`
	Provider    string    `json:"provider,omitempty"`
}

// GeneratePUMLRequest 生成PUML请求
type GeneratePUMLRequest struct {
	AnalysisID  string `json:"analysis_id" validate:"required"`
	DiagramType string `json:"diagram_type" validate:"required"`
	Provider    string `json:"provider,omitempty"`
}

// GenerateDocumentRequest 生成文档请求
type GenerateDocumentRequest struct {
	AnalysisID string `json:"analysis_id" validate:"required"`
	Provider   string `json:"provider,omitempty"`
}

// ChatSessionCreateRequest 创建对话会话请求
type ChatSessionCreateRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	Title     string `json:"title" validate:"required,min=1,max=100"`
}

// SendChatMessageRequest 发送聊天消息请求
type SendChatMessageRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Role      string `json:"role" validate:"required"`
}

// UpdatePUMLRequest 更新PUML请求
type UpdatePUMLRequest struct {
	Title       string `json:"title,omitempty"`
	Content     string `json:"content" validate:"required"`
	Description string `json:"description,omitempty"`
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content" validate:"required"`
}

// 新增状态常量
const (
	// 对话会话状态
	ChatSessionStatusActive    = "active"
	ChatSessionStatusCompleted = "completed"
	ChatSessionStatusArchived  = "archived"

	// 消息角色
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
	MessageRoleSystem    = "system"

	// 问题状态
	QuestionStatusPending  = "pending"
	QuestionStatusAnswered = "answered"
	QuestionStatusSkipped  = "skipped"

	// 问题分类
	QuestionCategoryBusinessRule      = "business_rule"
	QuestionCategoryExceptionHandling = "exception_handling"
	QuestionCategoryDataStructure     = "data_structure"
	QuestionCategoryExternalInterface = "external_interface"
	QuestionCategoryPerformance       = "performance_requirement"
	QuestionCategorySecurity          = "security_requirement"

	// AI提供商
	AIProviderOpenAI = "openai"
	AIProviderClaude = "claude"
	AIProviderGemini = "gemini"
)

// ===== 用户AI配置相关模型 =====

// UserAIConfig 用户AI配置
type UserAIConfig struct {
	ConfigID         uuid.UUID `json:"config_id" db:"config_id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	Provider         string    `json:"provider" db:"provider"`
	OpenAIAPIKey     string    `json:"openai_api_key,omitempty" db:"openai_api_key"`
	ClaudeAPIKey     string    `json:"claude_api_key,omitempty" db:"claude_api_key"`
	GeminiAPIKey     string    `json:"gemini_api_key,omitempty" db:"gemini_api_key"`
	DefaultModel     string    `json:"default_model" db:"default_model"`
	MaxTokens        int       `json:"max_tokens" db:"max_tokens"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateUserAIConfigRequest 更新用户AI配置请求
type UpdateUserAIConfigRequest struct {
	Provider      string `json:"provider" validate:"required"`
	OpenAIAPIKey  string `json:"openai_api_key,omitempty"`
	ClaudeAPIKey  string `json:"claude_api_key,omitempty"`
	GeminiAPIKey  string `json:"gemini_api_key,omitempty"`
	DefaultModel  string `json:"default_model" validate:"required"`
	MaxTokens     int    `json:"max_tokens" validate:"min=100,max=8192"`
}

// TestAIConnectionRequest 测试AI连接请求
type TestAIConnectionRequest struct {
	Provider string `json:"provider" validate:"required"`
	APIKey   string `json:"api_key" validate:"required"`
	Model    string `json:"model,omitempty"`
}

// AIConnectionTestResult AI连接测试结果
type AIConnectionTestResult struct {
	Success    bool   `json:"success"`
	Provider   string `json:"provider"`
	Model      string `json:"model,omitempty"`
	Message    string `json:"message"`
	Latency    int64  `json:"latency"` // 毫秒
	TokenUsage int    `json:"token_usage,omitempty"`
}

// ===== 分阶段文档生成相关模型 =====

// GenerateStageDocumentsRequest 分阶段文档生成请求
type GenerateStageDocumentsRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Stage     int       `json:"stage" validate:"required,min=1,max=3"`
}

// StageDocumentsResult 分阶段文档生成结果
type StageDocumentsResult struct {
	ProjectID     uuid.UUID      `json:"project_id"`
	Stage         int            `json:"stage"`
	GeneratedAt   time.Time      `json:"generated_at"`
	Documents     []*Document    `json:"documents"`
	PUMLDiagrams  []*PUMLDiagram `json:"puml_diagrams"`
}

// CompleteProjectDocumentsResult 完整项目文档生成结果
type CompleteProjectDocumentsResult struct {
	ProjectID         uuid.UUID             `json:"project_id"`
	GeneratedAt       time.Time             `json:"generated_at"`
	Stage1            *StageDocumentsResult `json:"stage1"`
	Stage2            *StageDocumentsResult `json:"stage2"`
	Stage3            *StageDocumentsResult `json:"stage3"`
	TotalDocuments    int                   `json:"total_documents"`
	TotalPUMLDiagrams int                   `json:"total_puml_diagrams"`
}

// ===== 新增异步任务管理模型 =====

// AsyncTask 异步任务表
type AsyncTask struct {
	TaskID      uuid.UUID  `json:"task_id" db:"task_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
	TaskType    string     `json:"task_type" db:"task_type"` // stage_document_generation, puml_generation, document_generation
	TaskName    string     `json:"task_name" db:"task_name"`
	Status      string     `json:"status" db:"status"`       // pending, running, completed, failed
	Progress    int        `json:"progress" db:"progress"`   // 0-100
	ResultData  string     `json:"result_data,omitempty" db:"result_data"` // JSON格式的结果数据
	ErrorMessage string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Metadata    string     `json:"metadata,omitempty" db:"metadata"` // JSON格式的任务元数据
}

// StageProgress 阶段进度跟踪
type StageProgress struct {
	ProgressID     uuid.UUID  `json:"progress_id" db:"progress_id"`
	ProjectID      uuid.UUID  `json:"project_id" db:"project_id"`
	Stage          int        `json:"stage" db:"stage"`           // 1, 2, 3
	Status         string     `json:"status" db:"status"`         // not_started, in_progress, completed, failed
	CompletionRate int        `json:"completion_rate" db:"completion_rate"` // 0-100
	StartedAt      *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	DocumentCount  int        `json:"document_count" db:"document_count"`
	PUMLCount      int        `json:"puml_count" db:"puml_count"`
	LastTaskID     *uuid.UUID `json:"last_task_id,omitempty" db:"last_task_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// 新增任务状态常量
const (
	// 异步任务状态
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"  
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"

	// 异步任务类型
	TaskTypeStageDocuments = "stage_document_generation"
	TaskTypePUMLGeneration = "puml_generation"
	TaskTypeDocumentGeneration = "document_generation"
	TaskTypeRequirementAnalysis = "requirement_analysis"
	TaskTypeCompleteProjectDocuments = "complete_project_documents" // 新增：一键生成完整项目文档

	// 阶段状态
	StageStatusNotStarted  = "not_started"
	StageStatusInProgress  = "in_progress"
	StageStatusCompleted   = "completed"
	StageStatusFailed      = "failed"
)

// ===== 请求类型 =====

// StartAsyncTaskRequest 启动异步任务请求
type StartAsyncTaskRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	TaskType  string    `json:"task_type" validate:"required"`
	Stage     *int      `json:"stage,omitempty"` // 阶段文档生成时需要
	Metadata  string    `json:"metadata,omitempty"` // JSON格式的额外参数
}

// AsyncTaskResponse 异步任务响应
type AsyncTaskResponse struct {
	TaskID   uuid.UUID `json:"task_id"`
	Status   string    `json:"status"`
	Progress int       `json:"progress"`
	Message  string    `json:"message,omitempty"`
}

// GetStageProgressRequest 获取阶段进度请求  
type GetStageProgressRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}

// StageProgressResponse 阶段进度响应
type StageProgressResponse struct {
	ProjectID uuid.UUID       `json:"project_id"`
	Stages    []*StageProgress `json:"stages"`
	Overall   struct {
		CompletionRate int    `json:"completion_rate"`
		Status         string `json:"status"`
	} `json:"overall"`
}