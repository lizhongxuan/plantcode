package repository

import (
	"database/sql"
	"fmt"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
)

// AIRepository AI分析相关的数据访问层
type AIRepository struct {
	db *sql.DB
}

// NewAIRepository 创建AI Repository
func NewAIRepository(db *sql.DB) *AIRepository {
	return &AIRepository{
		db: db,
	}
}

// ===== 需求分析相关操作 =====

// CreateRequirementAnalysis 创建需求分析记录
func (r *AIRepository) CreateRequirementAnalysis(analysis *model.Requirement) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		INSERT INTO requirement_analyses (
			requirement_id, project_id, raw_requirement, 
			structured_requirement, completeness_score, analysis_status,
			missing_info_types, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		analysis.RequirementID,
		analysis.ProjectID,
		analysis.RawRequirement,
		analysis.StructuredRequirement,
		analysis.CompletenessScore,
		analysis.AnalysisStatus,
		analysis.MissingInfoTypes,
		analysis.CreatedAt,
		analysis.UpdatedAt,
	)

	return err
}

// GetRequirementAnalysis 获取需求分析
func (r *AIRepository) GetRequirementAnalysis(analysisID uuid.UUID) (*model.Requirement, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT requirement_id, project_id, raw_requirement, 
			   structured_requirement, completeness_score, analysis_status,
			   missing_info_types, created_at, updated_at
		FROM requirement_analyses 
		WHERE requirement_id = ?
	`

	var analysis model.Requirement
	err := r.db.QueryRow(query, analysisID).Scan(
		&analysis.RequirementID,
		&analysis.ProjectID,
		&analysis.RawRequirement,
		&analysis.StructuredRequirement,
		&analysis.CompletenessScore,
		&analysis.AnalysisStatus,
		&analysis.MissingInfoTypes,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &analysis, nil
}

// GetRequirementAnalysesByProject 获取项目的所有需求分析
func (r *AIRepository) GetRequirementAnalysesByProject(projectID uuid.UUID) ([]*model.Requirement, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT requirement_id, project_id, raw_requirement, 
			   structured_requirement, completeness_score, analysis_status,
			   missing_info_types, created_at, updated_at
		FROM requirement_analyses 
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*model.Requirement
	for rows.Next() {
		var analysis model.Requirement
		err := rows.Scan(
			&analysis.RequirementID,
			&analysis.ProjectID,
			&analysis.RawRequirement,
			&analysis.StructuredRequirement,
			&analysis.CompletenessScore,
			&analysis.AnalysisStatus,
			&analysis.MissingInfoTypes,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, &analysis)
	}

	return analyses, nil
}

// UpdateRequirementAnalysis 更新需求分析
func (r *AIRepository) UpdateRequirementAnalysis(analysis *model.Requirement) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		UPDATE requirement_analyses 
		SET structured_requirement = ?, completeness_score = ?, 
			analysis_status = ?, missing_info_types = ?, updated_at = ?
		WHERE requirement_id = ?
	`

	_, err := r.db.Exec(query,
		analysis.StructuredRequirement,
		analysis.CompletenessScore,
		analysis.AnalysisStatus,
		analysis.MissingInfoTypes,
		time.Now(),
		analysis.RequirementID,
	)

	return err
}

// ===== 对话会话相关操作 =====

// CreateChatSession 创建对话会话
func (r *AIRepository) CreateChatSession(session *model.ChatSession) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		INSERT INTO chat_sessions (
			session_id, project_id, user_id, session_type,
			started_at, status, context
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		session.SessionID,
		session.ProjectID,
		session.UserID,
		session.SessionType,
		session.StartedAt,
		session.Status,
		session.Context,
	)

	return err
}

// GetChatSession 获取对话会话
func (r *AIRepository) GetChatSession(sessionID uuid.UUID) (*model.ChatSession, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT session_id, project_id, user_id, session_type,
			   started_at, ended_at, status, context
		FROM chat_sessions 
		WHERE session_id = ?
	`

	var session model.ChatSession
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.SessionID,
		&session.ProjectID,
		&session.UserID,
		&session.SessionType,
		&session.StartedAt,
		&session.EndedAt,
		&session.Status,
		&session.Context,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetChatSessionsByProject 获取项目的对话会话列表
func (r *AIRepository) GetChatSessionsByProject(projectID uuid.UUID) ([]*model.ChatSession, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT session_id, project_id, user_id, session_type,
			   started_at, ended_at, status, context
		FROM chat_sessions 
		WHERE project_id = ?
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.ChatSession
	for rows.Next() {
		var session model.ChatSession
		err := rows.Scan(
			&session.SessionID,
			&session.ProjectID,
			&session.UserID,
			&session.SessionType,
			&session.StartedAt,
			&session.EndedAt,
			&session.Status,
			&session.Context,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// ===== 对话消息相关操作 =====

// CreateChatMessage 创建对话消息
func (r *AIRepository) CreateChatMessage(message *model.ChatMessage) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		INSERT INTO chat_messages (
			message_id, session_id, sender_type, message_content,
			message_type, metadata, timestamp, processed
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		message.MessageID,
		message.SessionID,
		message.SenderType,
		message.MessageContent,
		message.MessageType,
		message.Metadata,
		message.Timestamp,
		message.Processed,
	)

	return err
}

// GetChatMessages 获取会话的所有消息
func (r *AIRepository) GetChatMessages(sessionID uuid.UUID) ([]*model.ChatMessage, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT message_id, session_id, sender_type, message_content,
			   message_type, metadata, timestamp, processed
		FROM chat_messages 
		WHERE session_id = ?
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		var message model.ChatMessage
		err := rows.Scan(
			&message.MessageID,
			&message.SessionID,
			&message.SenderType,
			&message.MessageContent,
			&message.MessageType,
			&message.Metadata,
			&message.Timestamp,
			&message.Processed,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// ===== 问题相关操作 =====

// CreateQuestion 创建补充问题
func (r *AIRepository) CreateQuestion(question *model.Question) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		INSERT INTO questions (
			question_id, requirement_id, question_text, question_category,
			priority_level, answer_status, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		question.QuestionID,
		question.RequirementID,
		question.QuestionText,
		question.QuestionCategory,
		question.PriorityLevel,
		question.AnswerStatus,
		question.CreatedAt,
	)

	return err
}

// GetQuestions 获取需求分析的所有问题
func (r *AIRepository) GetQuestions(requirementID uuid.UUID) ([]*model.Question, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT question_id, requirement_id, question_text, question_category,
			   priority_level, answer_text, answer_status, created_at, answered_at
		FROM questions 
		WHERE requirement_id = ?
		ORDER BY priority_level DESC, created_at ASC
	`

	rows, err := r.db.Query(query, requirementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*model.Question
	for rows.Next() {
		var question model.Question
		err := rows.Scan(
			&question.QuestionID,
			&question.RequirementID,
			&question.QuestionText,
			&question.QuestionCategory,
			&question.PriorityLevel,
			&question.AnswerText,
			&question.AnswerStatus,
			&question.CreatedAt,
			&question.AnsweredAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}

	return questions, nil
}

// AnswerQuestion 回答问题
func (r *AIRepository) AnswerQuestion(questionID uuid.UUID, answer string) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		UPDATE questions 
		SET answer_text = ?, answer_status = ?, answered_at = ?
		WHERE question_id = ?
	`

	now := time.Now()
	_, err := r.db.Exec(query, answer, model.QuestionStatusAnswered, now, questionID)
	return err
}

// ===== PUML图表相关操作 =====

// CreatePUMLDiagram 创建PUML图表
func (r *AIRepository) CreatePUMLDiagram(diagram *model.PUMLDiagram) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	// 使用更简单的插入语句，只包含基本字段
	query := `
		INSERT INTO puml_diagrams (
			diagram_id, project_id, diagram_type, diagram_name,
			puml_content, version, is_validated, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		diagram.DiagramID,
		diagram.ProjectID,
		diagram.DiagramType,
		diagram.DiagramName,
		diagram.PUMLContent,
		diagram.Version,
		diagram.IsValidated,
		diagram.CreatedAt,
		diagram.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("Insert error: %v\n", err)
		return err
	}

	return nil
}

// GetPUMLDiagram 获取PUML图表
func (r *AIRepository) GetPUMLDiagram(diagramID uuid.UUID) (*model.PUMLDiagram, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT diagram_id, project_id, diagram_type, diagram_name,
			   puml_content, rendered_url, version, stage, task_id, is_validated,
			   validation_feedback, created_at, updated_at
		FROM puml_diagrams 
		WHERE diagram_id = ?
	`

	var diagram model.PUMLDiagram
	err := r.db.QueryRow(query, diagramID).Scan(
		&diagram.DiagramID,
		&diagram.ProjectID,
		&diagram.DiagramType,
		&diagram.DiagramName,
		&diagram.PUMLContent,
		&diagram.RenderedURL,
		&diagram.Version,
		&diagram.Stage,
		&diagram.TaskID,
		&diagram.IsValidated,
		&diagram.ValidationFeedback,
		&diagram.CreatedAt,
		&diagram.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &diagram, nil
}

// GetPUMLDiagramsByProject 获取项目的所有PUML图表
func (r *AIRepository) GetPUMLDiagramsByProject(projectID uuid.UUID) ([]*model.PUMLDiagram, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	// 首先检查表结构
	checkQuery := "DESCRIBE puml_diagrams"
	rows, err := r.db.Query(checkQuery)
	if err != nil {
		fmt.Printf("Failed to describe table: %v\n", err)
	} else {
		fmt.Println("puml_diagrams table structure:")
		for rows.Next() {
			var field, fieldType, null, key, defaultVal, extra sql.NullString
			rows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
			fmt.Printf("  %s %s\n", field.String, fieldType.String)
		}
		rows.Close()
	}

	// 尝试一个更简单的查询
	simpleQuery := `
		SELECT diagram_id, project_id, diagram_type, diagram_name,
			   puml_content, version, is_validated, created_at, updated_at
		FROM puml_diagrams 
		WHERE project_id = ?
		ORDER BY created_at DESC
	`

	fmt.Printf("Trying simple query first...\n")
	rows, err = r.db.Query(simpleQuery, projectID)
	if err != nil {
		fmt.Printf("Simple query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var diagrams []*model.PUMLDiagram
	for rows.Next() {
		var diagram model.PUMLDiagram
		err := rows.Scan(
			&diagram.DiagramID,
			&diagram.ProjectID,
			&diagram.DiagramType,
			&diagram.DiagramName,
			&diagram.PUMLContent,
			&diagram.Version,
			&diagram.IsValidated,
			&diagram.CreatedAt,
			&diagram.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("Simple scan error: %v\n", err)
			return nil, err
		}
		// 设置默认值
		diagram.Stage = 1
		diagram.TaskID = nil
		diagram.RenderedURL = ""
		diagram.ValidationFeedback = ""
		
		diagrams = append(diagrams, &diagram)
	}

	fmt.Printf("Found %d diagrams\n", len(diagrams))
	return diagrams, nil
}

// UpdatePUMLDiagram 更新PUML图表
func (r *AIRepository) UpdatePUMLDiagram(diagram *model.PUMLDiagram) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		UPDATE puml_diagrams 
		SET diagram_name = ?, puml_content = ?, version = ?, 
			is_validated = ?, validation_feedback = ?, updated_at = ?
		WHERE diagram_id = ?
	`

	_, err := r.db.Exec(query,
		diagram.DiagramName,
		diagram.PUMLContent,
		diagram.Version,
		diagram.IsValidated,
		diagram.ValidationFeedback,
		time.Now(),
		diagram.DiagramID,
	)

	return err
}

// ===== 生成文档相关操作 =====

// CreateDocument 创建生成的文档
func (r *AIRepository) CreateDocument(document *model.Document) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		INSERT INTO generated_documents (
			document_id, project_id, document_type, document_name,
			content, format, version, generated_at, is_final
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		document.DocumentID,
		document.ProjectID,
		document.DocumentType,
		document.DocumentName,
		document.Content,
		document.Format,
		document.Version,
		document.GeneratedAt,
		document.IsFinal,
	)

	return err
}

// GetDocument 获取生成的文档
func (r *AIRepository) GetDocument(documentID uuid.UUID) (*model.Document, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT document_id, project_id, document_type, document_name,
			   content, format, file_path, version, generated_at, is_final
		FROM generated_documents 
		WHERE document_id = ?
	`

	var document model.Document
	err := r.db.QueryRow(query, documentID).Scan(
		&document.DocumentID,
		&document.ProjectID,
		&document.DocumentType,
		&document.DocumentName,
		&document.Content,
		&document.Format,
		&document.FilePath,
		&document.Version,
		&document.GeneratedAt,
		&document.IsFinal,
	)

	if err != nil {
		return nil, err
	}

	return &document, nil
}

// GetDocumentsByProject 获取项目的所有生成文档
func (r *AIRepository) GetDocumentsByProject(projectID uuid.UUID) ([]*model.Document, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT document_id, project_id, document_type, document_name,
			   content, format, file_path, version, generated_at, is_final
		FROM generated_documents 
		WHERE project_id = ?
		ORDER BY generated_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*model.Document
	for rows.Next() {
		var document model.Document
		err := rows.Scan(
			&document.DocumentID,
			&document.ProjectID,
			&document.DocumentType,
			&document.DocumentName,
			&document.Content,
			&document.Format,
			&document.FilePath,
			&document.Version,
			&document.GeneratedAt,
			&document.IsFinal,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	return documents, nil
}

// UpdateDocument 更新生成的文档
func (r *AIRepository) UpdateDocument(document *model.Document) error {
	if r.db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	query := `
		UPDATE generated_documents 
		SET document_name = ?, content = ?, version = ?, is_final = ?
		WHERE document_id = ?
	`

	_, err := r.db.Exec(query,
		document.DocumentName,
		document.Content,
		document.Version,
		document.IsFinal,
		document.DocumentID,
	)

	return err
}

// ===== 用户AI配置相关操作 =====

// GetUserAIConfig 获取用户AI配置
func (r *AIRepository) GetUserAIConfig(userID uuid.UUID) (*model.UserAIConfig, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	query := `
		SELECT config_id, user_id, provider, openai_api_key, claude_api_key, gemini_api_key,
		       default_model, max_tokens, is_active, created_at, updated_at
		FROM user_ai_configs 
		WHERE user_id = ? AND is_active = true
	`
	
	var config model.UserAIConfig
	err := r.db.QueryRow(query, userID).Scan(
		&config.ConfigID,
		&config.UserID,
		&config.Provider,
		&config.OpenAIAPIKey,
		&config.ClaudeAPIKey,
		&config.GeminiAPIKey,
		&config.DefaultModel,
		&config.MaxTokens,
		&config.IsActive,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user ai config not found")
		}
		return nil, fmt.Errorf("failed to get user ai config: %w", err)
	}
	
	return &config, nil
}

// CreateUserAIConfig 创建用户AI配置
func (r *AIRepository) CreateUserAIConfig(config *model.UserAIConfig) error {
	query := `
		INSERT INTO user_ai_configs (
			config_id, user_id, provider, openai_api_key, claude_api_key, gemini_api_key,
			default_model, max_tokens, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		config.ConfigID,
		config.UserID,
		config.Provider,
		config.OpenAIAPIKey,
		config.ClaudeAPIKey,
		config.GeminiAPIKey,
		config.DefaultModel,
		config.MaxTokens,
		config.IsActive,
		config.CreatedAt,
		config.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user ai config: %w", err)
	}
	
	return nil
}

// UpdateUserAIConfig 更新用户AI配置
func (r *AIRepository) UpdateUserAIConfig(config *model.UserAIConfig) error {
	query := `
		UPDATE user_ai_configs 
		SET provider = ?, openai_api_key = ?, claude_api_key = ?, gemini_api_key = ?,
		    default_model = ?, max_tokens = ?, updated_at = ?
		WHERE config_id = ?
	`
	
	result, err := r.db.Exec(query,
		config.Provider,
		config.OpenAIAPIKey,
		config.ClaudeAPIKey,
		config.GeminiAPIKey,
		config.DefaultModel,
		config.MaxTokens,
		config.UpdatedAt,
		config.ConfigID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user ai config: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user ai config not found")
	}
	
	return nil
}

// DeleteUserAIConfig 删除用户AI配置（软删除）
func (r *AIRepository) DeleteUserAIConfig(userID uuid.UUID) error {
	query := `
		UPDATE user_ai_configs 
		SET is_active = false, updated_at = ?
		WHERE user_id = ?
	`
	
	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user ai config: %w", err)
	}
	
	return nil
} 