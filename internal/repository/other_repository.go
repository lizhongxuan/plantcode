package repository

import (
	"fmt"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateRequirementAnalysis 创建需求分析
func (r *MySQLRepository) CreateRequirementAnalysis(requirement *model.Requirement) error {
	now := time.Now()
	requirement.CreatedAt = now
	requirement.UpdatedAt = now

	if err := r.db.GORM.Create(requirement).Error; err != nil {
		return fmt.Errorf("创建需求分析失败: %w", err)
	}

	return nil
}

// GetRequirementByProjectID 根据项目ID获取需求分析
func (r *MySQLRepository) GetRequirementByProjectID(projectID uuid.UUID) (*model.Requirement, error) {
	var requirement model.Requirement

	err := r.db.GORM.Where("project_id = ?", projectID).
		Order("created_at DESC").
		First(&requirement).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("需求分析不存在")
		}
		return nil, fmt.Errorf("查询需求分析失败: %w", err)
	}

	return &requirement, nil
}

// UpdateRequirementAnalysis 更新需求分析
func (r *MySQLRepository) UpdateRequirementAnalysis(requirement *model.Requirement) error {
	requirement.UpdatedAt = time.Now()

	result := r.db.GORM.Model(requirement).Where("requirement_id = ?", requirement.RequirementID).Updates(map[string]interface{}{
		"structured_requirement": requirement.StructuredRequirement,
		"completeness_score":     requirement.CompletenessScore,
		"analysis_status":        requirement.AnalysisStatus,
		"missing_info_types":     requirement.MissingInfoTypes,
		"updated_at":             requirement.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("更新需求分析失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("需求分析不存在或未更新")
	}

	return nil
}

// CreateChatSession 创建对话会话
func (r *MySQLRepository) CreateChatSession(session *model.ChatSession) error {
	now := time.Now()
	session.StartedAt = now

	if err := r.db.GORM.Create(session).Error; err != nil {
		return fmt.Errorf("创建对话会话失败: %w", err)
	}

	return nil
}

// GetChatSessionByProjectID 根据项目ID获取对话会话
func (r *MySQLRepository) GetChatSessionByProjectID(projectID uuid.UUID) (*model.ChatSession, error) {
	var session model.ChatSession

	err := r.db.GORM.Where("project_id = ?", projectID).
		Order("started_at DESC").
		First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("对话会话不存在")
		}
		return nil, fmt.Errorf("查询对话会话失败: %w", err)
	}

	return &session, nil
}

// CreateChatMessage 创建对话消息
func (r *MySQLRepository) CreateChatMessage(message *model.ChatMessage) error {
	message.Timestamp = time.Now()

	if err := r.db.GORM.Create(message).Error; err != nil {
		return fmt.Errorf("创建对话消息失败: %w", err)
	}

	return nil
}

// GetChatMessagesBySessionID 根据会话ID获取对话消息列表
func (r *MySQLRepository) GetChatMessagesBySessionID(sessionID uuid.UUID, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	var messages []*model.ChatMessage
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.GORM.Model(&model.ChatMessage{}).Where("session_id = ?", sessionID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询消息总数失败: %w", err)
	}

	// 获取消息列表
	if err := r.db.GORM.Where("session_id = ?", sessionID).
		Order("timestamp ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("查询消息列表失败: %w", err)
	}

	return messages, total, nil
}

// EndChatSession 结束对话会话
func (r *MySQLRepository) EndChatSession(sessionID uuid.UUID) error {
	now := time.Now()

	result := r.db.GORM.Model(&model.ChatSession{}).Where("session_id = ?", sessionID).Updates(map[string]interface{}{
		"ended_at": &now,
		"status":   "completed",
	})

	if result.Error != nil {
		return fmt.Errorf("结束对话会话失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("对话会话不存在")
	}

	return nil
}

// CreateQuestion 创建问题
func (r *MySQLRepository) CreateQuestion(question *model.Question) error {
	question.CreatedAt = time.Now()

	if err := r.db.GORM.Create(question).Error; err != nil {
		return fmt.Errorf("创建问题失败: %w", err)
	}

	return nil
}

// GetQuestionsByRequirementID 根据需求ID获取问题列表
func (r *MySQLRepository) GetQuestionsByRequirementID(requirementID uuid.UUID) ([]*model.Question, error) {
	var questions []*model.Question

	if err := r.db.GORM.Where("requirement_id = ?", requirementID).
		Order("priority_level DESC, created_at ASC").
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("查询问题列表失败: %w", err)
	}

	return questions, nil
}

// AnswerQuestion 回答问题
func (r *MySQLRepository) AnswerQuestion(questionID uuid.UUID, answer string) error {
	now := time.Now()

	result := r.db.GORM.Model(&model.Question{}).Where("question_id = ?", questionID).Updates(map[string]interface{}{
		"answer_text":   answer,
		"answer_status": model.AnswerStatusAnswered,
		"answered_at":   &now,
	})

	if result.Error != nil {
		return fmt.Errorf("回答问题失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("问题不存在")
	}

	return nil
}

// CreatePUMLDiagram 创建PUML图表
func (r *MySQLRepository) CreatePUMLDiagram(diagram *model.PUMLDiagram) error {
	now := time.Now()
	diagram.CreatedAt = now
	diagram.UpdatedAt = now

	if err := r.db.GORM.Create(diagram).Error; err != nil {
		return fmt.Errorf("创建PUML图表失败: %w", err)
	}

	return nil
}

// GetPUMLDiagramsByProjectID 根据项目ID获取PUML图表列表
func (r *MySQLRepository) GetPUMLDiagramsByProjectID(projectID uuid.UUID) ([]*model.PUMLDiagram, error) {
	var diagrams []*model.PUMLDiagram

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&diagrams).Error; err != nil {
		return nil, fmt.Errorf("查询PUML图表列表失败: %w", err)
	}

	return diagrams, nil
}

// UpdatePUMLDiagram 更新PUML图表
func (r *MySQLRepository) UpdatePUMLDiagram(diagram *model.PUMLDiagram) error {
	diagram.UpdatedAt = time.Now()

	result := r.db.GORM.Model(diagram).Where("diagram_id = ?", diagram.DiagramID).Updates(map[string]interface{}{
		"diagram_name":         diagram.DiagramName,
		"puml_content":         diagram.PUMLContent,
		"rendered_url":         diagram.RenderedURL,
		"version":              diagram.Version,
		"is_validated":         diagram.IsValidated,
		"validation_feedback":  diagram.ValidationFeedback,
		"updated_at":           diagram.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("更新PUML图表失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("PUML图表不存在或未更新")
	}

	return nil
}

// DeletePUMLDiagram 删除PUML图表
func (r *MySQLRepository) DeletePUMLDiagram(diagramID uuid.UUID) error {
	result := r.db.GORM.Delete(&model.PUMLDiagram{}, "diagram_id = ?", diagramID)

	if result.Error != nil {
		return fmt.Errorf("删除PUML图表失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("PUML图表不存在")
	}

	return nil
}

// CreateDocument 创建文档
func (r *MySQLRepository) CreateDocument(document *model.Document) error {
	document.GeneratedAt = time.Now()

	if err := r.db.GORM.Create(document).Error; err != nil {
		return fmt.Errorf("创建文档失败: %w", err)
	}

	return nil
}

// GetDocumentsByProjectID 根据项目ID获取文档列表
func (r *MySQLRepository) GetDocumentsByProjectID(projectID uuid.UUID) ([]*model.Document, error) {
	var documents []*model.Document

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("generated_at DESC").
		Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("查询文档列表失败: %w", err)
	}

	return documents, nil
}

// UpdateDocument 更新文档
func (r *MySQLRepository) UpdateDocument(document *model.Document) error {
	result := r.db.GORM.Model(document).Where("document_id = ?", document.DocumentID).Updates(map[string]interface{}{
		"document_name": document.DocumentName,
		"content":       document.Content,
		"format":        document.Format,
		"file_path":     document.FilePath,
		"version":       document.Version,
		"is_final":      document.IsFinal,
	})

	if result.Error != nil {
		return fmt.Errorf("更新文档失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("文档不存在或未更新")
	}

	return nil
}

// DeleteDocument 删除文档
func (r *MySQLRepository) DeleteDocument(documentID uuid.UUID) error {
	result := r.db.GORM.Delete(&model.Document{}, "document_id = ?", documentID)

	if result.Error != nil {
		return fmt.Errorf("删除文档失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("文档不存在")
	}

	return nil
}

// CreateBusinessModule 创建业务模块
func (r *MySQLRepository) CreateBusinessModule(module *model.BusinessModule) error {
	module.CreatedAt = time.Now()

	if err := r.db.GORM.Create(module).Error; err != nil {
		return fmt.Errorf("创建业务模块失败: %w", err)
	}

	return nil
}

// GetBusinessModulesByProjectID 根据项目ID获取业务模块列表
func (r *MySQLRepository) GetBusinessModulesByProjectID(projectID uuid.UUID) ([]*model.BusinessModule, error) {
	var modules []*model.BusinessModule

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&modules).Error; err != nil {
		return nil, fmt.Errorf("查询业务模块列表失败: %w", err)
	}

	return modules, nil
}

// UpdateBusinessModule 更新业务模块
func (r *MySQLRepository) UpdateBusinessModule(module *model.BusinessModule) error {
	result := r.db.GORM.Model(module).Where("module_id = ?", module.ModuleID).Updates(map[string]interface{}{
		"module_name":      module.ModuleName,
		"description":      module.Description,
		"module_type":      module.ModuleType,
		"complexity_level": module.ComplexityLevel,
		"business_logic":   module.BusinessLogic,
		"interfaces":       module.Interfaces,
		"dependencies":     module.Dependencies,
		"is_reusable":      module.IsReusable,
	})

	if result.Error != nil {
		return fmt.Errorf("更新业务模块失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("业务模块不存在或未更新")
	}

	return nil
}

// DeleteBusinessModule 删除业务模块
func (r *MySQLRepository) DeleteBusinessModule(moduleID uuid.UUID) error {
	result := r.db.GORM.Delete(&model.BusinessModule{}, "module_id = ?", moduleID)

	if result.Error != nil {
		return fmt.Errorf("删除业务模块失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("业务模块不存在")
	}

	return nil
}

// CreateCommonModule 创建通用模块
func (r *MySQLRepository) CreateCommonModule(module *model.CommonModule) error {
	now := time.Now()
	module.CreatedAt = now
	module.UpdatedAt = now

	if err := r.db.GORM.Create(module).Error; err != nil {
		return fmt.Errorf("创建通用模块失败: %w", err)
	}

	return nil
}

// GetCommonModulesByCategory 根据分类获取通用模块列表
func (r *MySQLRepository) GetCommonModulesByCategory(category string, page, pageSize int) ([]*model.CommonModule, int64, error) {
	var modules []*model.CommonModule
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	query := r.db.GORM.Model(&model.CommonModule{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询通用模块总数失败: %w", err)
	}

	// 获取模块列表
	query = r.db.GORM.Where("category = ? OR ? = ''", category, category).
		Order("downloads_count DESC, rating DESC, created_at DESC").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&modules).Error; err != nil {
		return nil, 0, fmt.Errorf("查询通用模块列表失败: %w", err)
	}

	return modules, total, nil
}

// GetCommonModuleByID 根据ID获取通用模块
func (r *MySQLRepository) GetCommonModuleByID(moduleID uuid.UUID) (*model.CommonModule, error) {
	var module model.CommonModule

	err := r.db.GORM.Where("common_module_id = ?", moduleID).First(&module).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("通用模块不存在")
		}
		return nil, fmt.Errorf("查询通用模块失败: %w", err)
	}

	return &module, nil
}

// UpdateCommonModule 更新通用模块
func (r *MySQLRepository) UpdateCommonModule(module *model.CommonModule) error {
	module.UpdatedAt = time.Now()

	result := r.db.GORM.Model(module).Where("common_module_id = ?", module.CommonModuleID).Updates(map[string]interface{}{
		"module_name":      module.ModuleName,
		"category":         module.Category,
		"description":      module.Description,
		"functionality":    module.Functionality,
		"interface_spec":   module.InterfaceSpec,
		"code_template":    module.CodeTemplate,
		"usage_examples":   module.UsageExamples,
		"version":          module.Version,
		"downloads_count":  module.DownloadsCount,
		"rating":           module.Rating,
		"tags":             module.Tags,
		"updated_at":       module.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("更新通用模块失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("通用模块不存在或未更新")
	}

	return nil
}

// DeleteCommonModule 删除通用模块
func (r *MySQLRepository) DeleteCommonModule(moduleID uuid.UUID) error {
	result := r.db.GORM.Delete(&model.CommonModule{}, "common_module_id = ?", moduleID)

	if result.Error != nil {
		return fmt.Errorf("删除通用模块失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("通用模块不存在")
	}

	return nil
}

// CreateAsyncTask 创建异步任务
func (r *MySQLRepository) CreateAsyncTask(task *model.AsyncTask) error {
	task.CreatedAt = time.Now()

	if err := r.db.GORM.Create(task).Error; err != nil {
		return fmt.Errorf("创建异步任务失败: %w", err)
	}

	return nil
}

// GetAsyncTask 根据ID获取异步任务
func (r *MySQLRepository) GetAsyncTask(taskID uuid.UUID) (*model.AsyncTask, error) {
	var task model.AsyncTask

	err := r.db.GORM.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("异步任务不存在")
		}
		return nil, fmt.Errorf("查询异步任务失败: %w", err)
	}

	return &task, nil
}

// UpdateAsyncTask 更新异步任务
func (r *MySQLRepository) UpdateAsyncTask(task *model.AsyncTask) error {
	result := r.db.GORM.Model(task).Where("task_id = ?", task.TaskID).Updates(map[string]interface{}{
		"status":        task.Status,
		"progress":      task.Progress,
		"result_data":   task.ResultData,
		"error_message": task.ErrorMessage,
		"started_at":    task.StartedAt,
		"completed_at":  task.CompletedAt,
		"metadata":      task.Metadata,
	})

	if result.Error != nil {
		return fmt.Errorf("更新异步任务失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("异步任务不存在或未更新")
	}

	return nil
}

// GetTasksByProject 根据项目ID获取异步任务列表
func (r *MySQLRepository) GetTasksByProject(projectID uuid.UUID, taskType string) ([]*model.AsyncTask, error) {
	var tasks []*model.AsyncTask

	query := r.db.GORM.Where("project_id = ?", projectID)
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	if err := query.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("查询异步任务列表失败: %w", err)
	}

	return tasks, nil
}

// CreateStageProgress 创建阶段进度
func (r *MySQLRepository) CreateStageProgress(progress *model.StageProgress) error {
	now := time.Now()
	progress.CreatedAt = now
	progress.UpdatedAt = now

	if err := r.db.GORM.Create(progress).Error; err != nil {
		return fmt.Errorf("创建阶段进度失败: %w", err)
	}

	return nil
}

// GetStageProgress 根据项目ID获取阶段进度列表
func (r *MySQLRepository) GetStageProgress(projectID uuid.UUID) ([]*model.StageProgress, error) {
	var progresses []*model.StageProgress

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("stage ASC").
		Find(&progresses).Error; err != nil {
		return nil, fmt.Errorf("查询阶段进度列表失败: %w", err)
	}

	return progresses, nil
}

// UpdateStageProgress 更新阶段进度
func (r *MySQLRepository) UpdateStageProgress(progress *model.StageProgress) error {
	progress.UpdatedAt = time.Now()

	result := r.db.GORM.Model(progress).Where("progress_id = ?", progress.ProgressID).Updates(map[string]interface{}{
		"status":          progress.Status,
		"completion_rate": progress.CompletionRate,
		"started_at":      progress.StartedAt,
		"completed_at":    progress.CompletedAt,
		"document_count":  progress.DocumentCount,
		"puml_count":      progress.PUMLCount,
		"last_task_id":    progress.LastTaskID,
		"updated_at":      progress.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("更新阶段进度失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("阶段进度不存在或未更新")
	}

	return nil
}

// GetStageProgressByStage 根据项目ID和阶段获取阶段进度
func (r *MySQLRepository) GetStageProgressByStage(projectID uuid.UUID, stage int) (*model.StageProgress, error) {
	var progress model.StageProgress

	err := r.db.GORM.Where("project_id = ? AND stage = ?", projectID, stage).First(&progress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("阶段进度不存在")
		}
		return nil, fmt.Errorf("查询阶段进度失败: %w", err)
	}

	return &progress, nil
}

// ===== 用户AI配置相关方法 =====

// GetUserAIConfig 获取用户AI配置
func (r *MySQLRepository) GetUserAIConfig(userID uuid.UUID) (*model.UserAIConfig, error) {
	var config model.UserAIConfig

	err := r.db.GORM.Where("user_id = ? AND is_active = ?", userID, true).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户AI配置不存在")
		}
		return nil, fmt.Errorf("查询用户AI配置失败: %w", err)
	}

	return &config, nil
}

// CreateUserAIConfig 创建用户AI配置
func (r *MySQLRepository) CreateUserAIConfig(config *model.UserAIConfig) error {
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	if err := r.db.GORM.Create(config).Error; err != nil {
		return fmt.Errorf("创建用户AI配置失败: %w", err)
	}

	return nil
}

// UpdateUserAIConfig 更新用户AI配置
func (r *MySQLRepository) UpdateUserAIConfig(config *model.UserAIConfig) error {
	config.UpdatedAt = time.Now()

	result := r.db.GORM.Model(config).Where("config_id = ?", config.ConfigID).Updates(map[string]interface{}{
		"provider":       config.Provider,
		"openai_api_key": config.OpenAIAPIKey,
		"claude_api_key": config.ClaudeAPIKey,
		"gemini_api_key": config.GeminiAPIKey,
		"default_model":  config.DefaultModel,
		"max_tokens":     config.MaxTokens,
		"updated_at":     config.UpdatedAt,
	})

	if result.Error != nil {
		return fmt.Errorf("更新用户AI配置失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户AI配置不存在或未更新")
	}

	return nil
}

// DeleteUserAIConfig 删除用户AI配置（软删除）
func (r *MySQLRepository) DeleteUserAIConfig(userID uuid.UUID) error {
	now := time.Now()

	result := r.db.GORM.Model(&model.UserAIConfig{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": now,
	})

	if result.Error != nil {
		return fmt.Errorf("删除用户AI配置失败: %w", result.Error)
	}

	return nil
}

// ===== 扩展方法（用于兼容性） =====

// GetRequirementAnalysis 根据ID获取需求分析
func (r *MySQLRepository) GetRequirementAnalysis(analysisID uuid.UUID) (*model.Requirement, error) {
	var requirement model.Requirement

	err := r.db.GORM.Where("requirement_id = ?", analysisID).First(&requirement).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("需求分析不存在")
		}
		return nil, fmt.Errorf("查询需求分析失败: %w", err)
	}

	return &requirement, nil
}

// GetRequirementAnalysesByProject 根据项目ID获取需求分析列表
func (r *MySQLRepository) GetRequirementAnalysesByProject(projectID uuid.UUID) ([]*model.Requirement, error) {
	var requirements []*model.Requirement

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&requirements).Error; err != nil {
		return nil, fmt.Errorf("查询需求分析列表失败: %w", err)
	}

	return requirements, nil
}

// GetChatSession 根据ID获取对话会话
func (r *MySQLRepository) GetChatSession(sessionID uuid.UUID) (*model.ChatSession, error) {
	var session model.ChatSession

	err := r.db.GORM.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("对话会话不存在")
		}
		return nil, fmt.Errorf("查询对话会话失败: %w", err)
	}

	return &session, nil
}

// GetChatSessionsByProject 根据项目ID获取对话会话列表
func (r *MySQLRepository) GetChatSessionsByProject(projectID uuid.UUID) ([]*model.ChatSession, error) {
	var sessions []*model.ChatSession

	if err := r.db.GORM.Where("project_id = ?", projectID).
		Order("started_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("查询对话会话列表失败: %w", err)
	}

	return sessions, nil
}

// GetChatMessages 根据会话ID获取对话消息列表
func (r *MySQLRepository) GetChatMessages(sessionID uuid.UUID) ([]*model.ChatMessage, error) {
	var messages []*model.ChatMessage

	if err := r.db.GORM.Where("session_id = ?", sessionID).
		Order("timestamp ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询对话消息列表失败: %w", err)
	}

	return messages, nil
}

// GetPUMLDiagram 根据ID获取PUML图表
func (r *MySQLRepository) GetPUMLDiagram(diagramID uuid.UUID) (*model.PUMLDiagram, error) {
	var diagram model.PUMLDiagram

	err := r.db.GORM.Where("diagram_id = ?", diagramID).First(&diagram).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("PUML图表不存在")
		}
		return nil, fmt.Errorf("查询PUML图表失败: %w", err)
	}

	return &diagram, nil
}

// GetDocument 根据ID获取文档
func (r *MySQLRepository) GetDocument(documentID uuid.UUID) (*model.Document, error) {
	var document model.Document

	err := r.db.GORM.Where("document_id = ?", documentID).First(&document).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文档不存在")
		}
		return nil, fmt.Errorf("查询文档失败: %w", err)
	}

	return &document, nil
}

// GetQuestions 根据需求ID获取问题列表（别名方法）
func (r *MySQLRepository) GetQuestions(requirementID uuid.UUID) ([]*model.Question, error) {
	return r.GetQuestionsByRequirementID(requirementID)
}