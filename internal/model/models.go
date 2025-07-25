package model

import (
	"time"

	"github.com/google/uuid"
)

// 用户状态常量
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// 项目状态常量
const (
	ProjectStatusDraft     = "draft"
	ProjectStatusAnalyzing = "analyzing"
	ProjectStatusCompleted = "completed"
	ProjectStatusArchived  = "archived"
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
	AnalysisStatusPending            = "pending"
	AnalysisStatusInProgress         = "in_progress"
	AnalysisStatusQuestionsGenerated = "questions_generated"
	AnalysisStatusCompleted          = "completed"
	AnalysisStatusFailed             = "failed"
)

// 问题回答状态常量
const (
	AnswerStatusPending  = "pending"
	AnswerStatusAnswered = "answered"
	AnswerStatusSkipped  = "skipped"
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
	ConfigID     uuid.UUID `json:"config_id" gorm:"type:char(36);primaryKey;column:config_id" db:"config_id"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index;column:user_id" db:"user_id"`
	Provider     string    `json:"provider" gorm:"type:varchar(20);not null;column:provider" db:"provider"`
	OpenAIAPIKey string    `json:"openai_api_key,omitempty" gorm:"type:varchar(255);column:openai_api_key" db:"openai_api_key"`
	ClaudeAPIKey string    `json:"claude_api_key,omitempty" gorm:"type:varchar(255);column:claude_api_key" db:"claude_api_key"`
	GeminiAPIKey string    `json:"gemini_api_key,omitempty" gorm:"type:varchar(255);column:gemini_api_key" db:"gemini_api_key"`
	DefaultModel string    `json:"default_model" gorm:"type:varchar(50);not null;column:default_model" db:"default_model"`
	MaxTokens    int       `json:"max_tokens" gorm:"default:4096;column:max_tokens" db:"max_tokens"`
	IsActive     bool      `json:"is_active" gorm:"default:true;column:is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at" db:"updated_at"`
}

// TableName 指定表名
func (UserAIConfig) TableName() string {
	return "user_ai_configs"
}

// UpdateUserAIConfigRequest 更新用户AI配置请求
type UpdateUserAIConfigRequest struct {
	Provider     string `json:"provider" validate:"required"`
	OpenAIAPIKey string `json:"openai_api_key,omitempty"`
	ClaudeAPIKey string `json:"claude_api_key,omitempty"`
	GeminiAPIKey string `json:"gemini_api_key,omitempty"`
	DefaultModel string `json:"default_model" validate:"required"`
	MaxTokens    int    `json:"max_tokens" validate:"min=100,max=8192"`
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
	ProjectID    uuid.UUID      `json:"project_id"`
	Stage        int            `json:"stage"`
	GeneratedAt  time.Time      `json:"generated_at"`
	Documents    []*Document    `json:"documents"`
	PUMLDiagrams []*PUMLDiagram `json:"puml_diagrams"`
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

// 新增任务状态常量
const (
	// 异步任务状态
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"

	// 异步任务类型
	TaskTypeStageDocuments           = "stage_document_generation"
	TaskTypePUMLGeneration           = "puml_generation"
	TaskTypeDocumentGeneration       = "document_generation"
	TaskTypeRequirementAnalysis      = "requirement_analysis"
	TaskTypeCompleteProjectDocuments = "complete_project_documents" // 新增：一键生成完整项目文档

	// 阶段状态
	StageStatusNotStarted = "not_started"
	StageStatusInProgress = "in_progress"
	StageStatusCompleted  = "completed"
	StageStatusFailed     = "failed"
)

// ===== 请求类型 =====

// StartAsyncTaskRequest 启动异步任务请求
type StartAsyncTaskRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	TaskType  string    `json:"task_type" validate:"required"`
	Stage     *int      `json:"stage,omitempty"`    // 阶段文档生成时需要
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
	ProjectID uuid.UUID        `json:"project_id"`
	Stages    []*StageProgress `json:"stages"`
	Overall   struct {
		CompletionRate int    `json:"completion_rate"`
		Status         string `json:"status"`
	} `json:"overall"`
}

// ===== 文件夹结构管理模型 =====

// ProjectFolder 项目文件夹模型
type ProjectFolder struct {
	FolderID   uuid.UUID  `json:"folder_id" db:"folder_id"`
	ProjectID  uuid.UUID  `json:"project_id" db:"project_id"`
	FolderName string     `json:"folder_name" db:"folder_name"` // requirements, design, tasks
	FolderType string     `json:"folder_type" db:"folder_type"` // stage_folder
	ParentID   *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	SortOrder  int        `json:"sort_order" db:"sort_order"` // 排序顺序
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// ProjectDocument 项目文档模型（重构）
type ProjectDocument struct {
	DocumentID   uuid.UUID `json:"document_id" db:"document_id"`
	ProjectID    uuid.UUID `json:"project_id" db:"project_id"`
	FolderID     uuid.UUID `json:"folder_id" db:"folder_id"` // 所属文件夹
	DocumentName string    `json:"document_name" db:"document_name"`
	DocumentType string    `json:"document_type" db:"document_type"` // requirements_doc, design_doc, task_list, puml_diagram
	Content      string    `json:"content" db:"content"`             // markdown或其他格式内容
	Version      int       `json:"version" db:"version"`
	IsTemplate   bool      `json:"is_template" db:"is_template"` // 是否为模板
	CreatedBy    uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DocumentVersion 文档版本历史模型
type DocumentVersion struct {
	VersionID     uuid.UUID `json:"version_id" db:"version_id"`
	DocumentID    uuid.UUID `json:"document_id" db:"document_id"`
	VersionNumber int       `json:"version_number" db:"version_number"`
	Content       string    `json:"content" db:"content"`
	ChangeNote    string    `json:"change_note" db:"change_note"` // 变更说明
	ChangedBy     uuid.UUID `json:"changed_by" db:"changed_by"`   // 修改者
	ChangeType    string    `json:"change_type" db:"change_type"` // manual, ai_generated
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// AIConversation AI对话记录模型（重构聊天功能）
type AIConversation struct {
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	ProjectID      uuid.UUID `json:"project_id" db:"project_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Title          string    `json:"title" db:"title"`
	Context        string    `json:"context" db:"context"`     // 当前上下文信息（JSON）
	IsActive       bool      `json:"is_active" db:"is_active"` // 是否为活跃对话
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// AIMessage AI消息模型
type AIMessage struct {
	MessageID      uuid.UUID `json:"message_id" db:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"` // user, assistant, system
	Content        string    `json:"content" db:"content"`
	MessageType    string    `json:"message_type" db:"message_type"` // text, document_change, file_operation
	Metadata       string    `json:"metadata" db:"metadata"`         // 关联的文档或操作信息（JSON）
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// DocumentChange 文档变更记录模型
type DocumentChange struct {
	ChangeID      uuid.UUID  `json:"change_id" db:"change_id"`
	DocumentID    uuid.UUID  `json:"document_id" db:"document_id"`
	MessageID     *uuid.UUID `json:"message_id,omitempty" db:"message_id"` // 关联的AI消息
	ChangeType    string     `json:"change_type" db:"change_type"`         // create, update, delete
	OldContent    string     `json:"old_content" db:"old_content"`
	NewContent    string     `json:"new_content" db:"new_content"`
	ChangeSummary string     `json:"change_summary" db:"change_summary"` // 变更摘要
	IsAIGenerated bool       `json:"is_ai_generated" db:"is_ai_generated"`
	ChangedBy     uuid.UUID  `json:"changed_by" db:"changed_by"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// ===== 文件夹和文档相关常量 =====

const (
	// 文件夹类型
	FolderTypeStage = "stage_folder"

	// 标准文件夹名称
	FolderNameRequirements = "requirements"
	FolderNameDesign       = "design"
	FolderNameTasks        = "tasks"

	// 文档类型
	DocumentTypeRequirements = "requirements_doc"
	DocumentTypeDesign       = "design_doc"
	DocumentTypeTaskList     = "task_list"
	DocumentTypePUML         = "puml_diagram"
	DocumentTypeGeneral      = "general"

	// 变更类型
	ChangeTypeCreate = "create"
	ChangeTypeUpdate = "update"
	ChangeTypeDelete = "delete"

	// 消息角色
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"

	// 消息类型
	MessageTypeTextMsg        = "text"
	MessageTypeDocumentChange = "document_change"
	MessageTypeFileOperation  = "file_operation"
)

// ===== 文件夹和文档管理请求类型 =====

// CreateProjectFoldersRequest 创建项目文件夹请求
type CreateProjectFoldersRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
	ProjectID    uuid.UUID `json:"project_id" validate:"required"`
	FolderID     uuid.UUID `json:"folder_id" validate:"required"`
	DocumentName string    `json:"document_name" validate:"required,min=1,max=100"`
	DocumentType string    `json:"document_type" validate:"required"`
	Content      string    `json:"content"`
	IsTemplate   bool      `json:"is_template"`
}

// UpdateDocumentRequestNew 更新文档请求
type UpdateDocumentRequestNew struct {
	DocumentID uuid.UUID `json:"document_id" validate:"required"`
	Content    string    `json:"content" validate:"required"`
	ChangeNote string    `json:"change_note"`
}

// GetProjectStructureRequest 获取项目结构请求
type GetProjectStructureRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}

// ProjectStructureResponse 项目结构响应
type ProjectStructureResponse struct {
	ProjectID uuid.UUID                `json:"project_id"`
	Folders   []*ProjectFolderWithDocs `json:"folders"`
}

// ProjectFolderWithDocs 包含文档的文件夹
type ProjectFolderWithDocs struct {
	*ProjectFolder
	Documents []*ProjectDocument `json:"documents"`
}

// ===== AI对话相关请求类型 =====

// StartAIConversationRequest 开始AI对话请求
type StartAIConversationRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Title     string    `json:"title" validate:"required,min=1,max=100"`
}

// SendAIMessageRequest 发送AI消息请求
type SendAIMessageRequest struct {
	ConversationID uuid.UUID `json:"conversation_id" validate:"required"`
	Content        string    `json:"content" validate:"required"`
	MessageType    string    `json:"message_type"`
}

// AIConversationResponse AI对话响应
type AIConversationResponse struct {
	ConversationID uuid.UUID    `json:"conversation_id"`
	Messages       []*AIMessage `json:"messages"`
	Context        string       `json:"context"`
}

// ===== 文档变更相关请求类型 =====

// GetDocumentChangesRequest 获取文档变更请求
type GetDocumentChangesRequest struct {
	DocumentID uuid.UUID `json:"document_id" validate:"required"`
	Limit      *int      `json:"limit,omitempty"`
}

// DocumentChangesResponse 文档变更响应
type DocumentChangesResponse struct {
	DocumentID uuid.UUID          `json:"document_id"`
	Changes    []*DocumentChange  `json:"changes"`
	Versions   []*DocumentVersion `json:"versions"`
}

// RevertDocumentRequest 回滚文档请求
type RevertDocumentRequest struct {
	DocumentID    uuid.UUID `json:"document_id" validate:"required"`
	VersionNumber int       `json:"version_number" validate:"required"`
	ReasonNote    string    `json:"reason_note"`
}

// ===== Spec 工作流相关模型 =====

// ProjectSpec 项目规范模型
type ProjectSpec struct {
	SpecID       uuid.UUID `json:"spec_id" db:"spec_id"`
	ProjectID    uuid.UUID `json:"project_id" db:"project_id"`
	CurrentStage string    `json:"current_stage" db:"current_stage"` // requirements, design, tasks, implementation
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserStory 用户故事模型
type UserStory struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	RequirementsID     uuid.UUID `json:"requirements_id" db:"requirements_id"`
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	AcceptanceCriteria string    `json:"acceptance_criteria" db:"acceptance_criteria"` // JSON 数组
	Priority           string    `json:"priority" db:"priority"`
	StoryPoints        *int      `json:"story_points" db:"story_points"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// RequirementsDoc 需求文档模型
type RequirementsDoc struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	ProjectID                 uuid.UUID `json:"project_id" db:"project_id"`
	Content                   string    `json:"content" db:"content"`                                         // markdown 格式
	Assumptions               string    `json:"assumptions" db:"assumptions"`                                 // JSON 数组
	EdgeCases                 string    `json:"edge_cases" db:"edge_cases"`                                   // JSON 数组
	FunctionalRequirements    string    `json:"functional_requirements" db:"functional_requirements"`         // JSON 数组
	NonFunctionalRequirements string    `json:"non_functional_requirements" db:"non_functional_requirements"` // JSON 数组
	Version                   int       `json:"version" db:"version"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

// SpecPUMLDiagram Spec 相关的 PUML 图表模型（扩展原有的 PUMLDiagram）
type SpecPUMLDiagram struct {
	ID          uuid.UUID `json:"id" db:"id"`
	DesignID    uuid.UUID `json:"design_id" db:"design_id"`
	Title       string    `json:"title" db:"title"`
	Type        string    `json:"type" db:"type"` // sequence, class, activity, component, deployment, use_case
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TypeScriptInterface TypeScript 接口模型
type TypeScriptInterface struct {
	ID          uuid.UUID `json:"id" db:"id"`
	DesignID    uuid.UUID `json:"design_id" db:"design_id"`
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// APIEndpoint API 端点模型
type APIEndpoint struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DesignID     uuid.UUID `json:"design_id" db:"design_id"`
	Path         string    `json:"path" db:"path"`
	Method       string    `json:"method" db:"method"`
	Description  string    `json:"description" db:"description"`
	RequestBody  string    `json:"request_body" db:"request_body"`   // JSON
	ResponseBody string    `json:"response_body" db:"response_body"` // JSON
	Headers      string    `json:"headers" db:"headers"`             // JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DesignDoc 设计文档模型
type DesignDoc struct {
	ID                uuid.UUID `json:"id" db:"id"`
	ProjectID         uuid.UUID `json:"project_id" db:"project_id"`
	Content           string    `json:"content" db:"content"` // markdown 格式
	DatabaseSchema    string    `json:"database_schema" db:"database_schema"`
	ArchitectureNotes string    `json:"architecture_notes" db:"architecture_notes"` // JSON 数组
	Version           int       `json:"version" db:"version"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// DevelopmentTask 开发任务模型
type DevelopmentTask struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	TaskListID     uuid.UUID  `json:"task_list_id" db:"task_list_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Type           string     `json:"type" db:"type"` // feature, bug, refactor, test, docs
	Priority       string     `json:"priority" db:"priority"`
	Status         string     `json:"status" db:"status"` // todo, in_progress, review, done
	EstimatedHours int        `json:"estimated_hours" db:"estimated_hours"`
	ActualHours    *int       `json:"actual_hours" db:"actual_hours"`
	Assignee       string     `json:"assignee" db:"assignee"`
	Dependencies   string     `json:"dependencies" db:"dependencies"` // JSON 数组: task ids
	UserStoryID    *uuid.UUID `json:"user_story_id" db:"user_story_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// TestCase 测试用例模型
type TestCase struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	TaskListID     uuid.UUID  `json:"task_list_id" db:"task_list_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Type           string     `json:"type" db:"type"`   // unit, integration, e2e, api
	Steps          string     `json:"steps" db:"steps"` // JSON 数组
	ExpectedResult string     `json:"expected_result" db:"expected_result"`
	TaskID         *uuid.UUID `json:"task_id" db:"task_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// TaskListDoc 任务列表文档模型
type TaskListDoc struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	ProjectID           uuid.UUID `json:"project_id" db:"project_id"`
	Content             string    `json:"content" db:"content"` // markdown 格式
	EstimatedTotalHours int       `json:"estimated_total_hours" db:"estimated_total_hours"`
	Milestones          string    `json:"milestones" db:"milestones"` // JSON 数组
	Version             int       `json:"version" db:"version"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// Spec 阶段常量
const (
	SpecStageRequirements   = "requirements"
	SpecStageDesign         = "design"
	SpecStageTasks          = "tasks"
	SpecStageImplementation = "implementation"
)

// 任务类型常量
const (
	TaskTypeFeature  = "feature"
	TaskTypeBug      = "bug"
	TaskTypeRefactor = "refactor"
	TaskTypeTest     = "test"
	TaskTypeDocs     = "docs"
)

// 任务状态常量
const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusReview     = "review"
	TaskStatusDone       = "done"
)

// 测试用例类型常量
const (
	TestTypeUnit        = "unit"
	TestTypeIntegration = "integration"
	TestTypeE2E         = "e2e"
	TestTypeAPI         = "api"
)

// 优先级常量
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// ===== Spec 工作流请求类型 =====

// GenerateRequirementsRequest 生成需求请求
type GenerateRequirementsRequest struct {
	ProjectID      uuid.UUID `json:"project_id" validate:"required"`
	InitialPrompt  string    `json:"initial_prompt" validate:"required,min=10"`
	ProjectType    string    `json:"project_type" validate:"required"`
	TargetAudience string    `json:"target_audience,omitempty"`
	BusinessGoals  string    `json:"business_goals,omitempty"` // JSON 数组
}

// GenerateDesignRequest 生成设计请求
type GenerateDesignRequest struct {
	ProjectID         uuid.UUID `json:"project_id" validate:"required"`
	RequirementsID    uuid.UUID `json:"requirements_id" validate:"required"`
	FocusAreas        string    `json:"focus_areas,omitempty"`        // JSON 数组
	ArchitectureStyle string    `json:"architecture_style,omitempty"` // monolith, microservices, serverless
}

// GenerateTasksRequest 生成任务请求
type GenerateTasksRequest struct {
	ProjectID      uuid.UUID `json:"project_id" validate:"required"`
	RequirementsID uuid.UUID `json:"requirements_id" validate:"required"`
	DesignID       uuid.UUID `json:"design_id" validate:"required"`
	TeamSize       *int      `json:"team_size,omitempty"`
	SprintDuration *int      `json:"sprint_duration,omitempty"`
}

// UpdateSpecStageRequest 更新 Spec 阶段请求
type UpdateSpecStageRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Stage     string    `json:"stage" validate:"required"`
}

// SpecResponse Spec 响应基础结构
type SpecResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type CreatePUMLRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	Title     string `json:"title" validate:"required,min=1,max=100"`
	Content   string `json:"content" validate:"required"`
}

// RenderPUMLRequest 渲染PUML请求
type RenderPUMLRequest struct {
	Content string `json:"content" validate:"required"`
	Format  string `json:"format,omitempty"` // png, svg, txt
}

// GenerateImageRequest 生成图片请求
type GenerateImageRequest struct {
	Content string `json:"content" validate:"required"`
	Format  string `json:"format,omitempty"` // png, svg
}

// ValidatePUMLRequest 验证PUML请求
type ValidatePUMLRequest struct {
	Content string `json:"content" validate:"required"`
}

// PreviewPUMLRequest 预览PUML请求
type PreviewPUMLRequest struct {
	Content string `json:"content" validate:"required"`
}

// ExportPUMLRequest 导出PUML请求
type ExportPUMLRequest struct {
	PUMLIDs []string `json:"puml_ids" validate:"required"`
	Format  string   `json:"format" validate:"required"`
}
