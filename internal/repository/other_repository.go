package repository

import (
	"database/sql"
	"fmt"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
)

// CreateRequirementAnalysis 创建需求分析
func (r *MySQLRepository) CreateRequirementAnalysis(requirement *model.Requirement) error {
	query := `
		INSERT INTO requirement_analyses (requirement_id, project_id, raw_requirement, structured_requirement, completeness_score, analysis_status, missing_info_types, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.MySQL.Exec(query,
		requirement.RequirementID,
		requirement.ProjectID,
		requirement.RawRequirement,
		requirement.StructuredRequirement,
		requirement.CompletenessScore,
		requirement.AnalysisStatus,
		requirement.MissingInfoTypes,
		now,
		now,
	)
	
	if err != nil {
		return fmt.Errorf("创建需求分析失败: %w", err)
	}
	
	return nil
}

// GetRequirementByProjectID 根据项目ID获取需求分析
func (r *MySQLRepository) GetRequirementByProjectID(projectID uuid.UUID) (*model.Requirement, error) {
	query := `
		SELECT requirement_id, project_id, raw_requirement, structured_requirement, completeness_score, analysis_status, missing_info_types, created_at, updated_at
		FROM requirement_analyses
		WHERE project_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	row := r.db.MySQL.QueryRow(query, projectID)
	
	var requirement model.Requirement
	var structuredRequirement, missingInfoTypes sql.NullString
	
	err := row.Scan(
		&requirement.RequirementID,
		&requirement.ProjectID,
		&requirement.RawRequirement,
		&structuredRequirement,
		&requirement.CompletenessScore,
		&requirement.AnalysisStatus,
		&missingInfoTypes,
		&requirement.CreatedAt,
		&requirement.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("需求分析不存在")
		}
		return nil, fmt.Errorf("查询需求分析失败: %w", err)
	}
	
	if structuredRequirement.Valid {
		requirement.StructuredRequirement = structuredRequirement.String
	}
	if missingInfoTypes.Valid {
		requirement.MissingInfoTypes = missingInfoTypes.String
	}
	
	return &requirement, nil
}

// UpdateRequirementAnalysis 更新需求分析
func (r *MySQLRepository) UpdateRequirementAnalysis(requirement *model.Requirement) error {
	query := `
		UPDATE requirement_analyses 
		SET structured_requirement = ?, completeness_score = ?, analysis_status = ?, missing_info_types = ?, updated_at = ?
		WHERE requirement_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query,
		requirement.StructuredRequirement,
		requirement.CompletenessScore,
		requirement.AnalysisStatus,
		requirement.MissingInfoTypes,
		now,
		requirement.RequirementID,
	)
	
	if err != nil {
		return fmt.Errorf("更新需求分析失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("需求分析不存在或未更新")
	}
	
	return nil
}

// CreateChatSession 创建对话会话
func (r *MySQLRepository) CreateChatSession(session *model.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (session_id, project_id, user_id, session_type, started_at, status, context)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.MySQL.Exec(query,
		session.SessionID,
		session.ProjectID,
		session.UserID,
		session.SessionType,
		session.StartedAt,
		session.Status,
		session.Context,
	)
	
	if err != nil {
		return fmt.Errorf("创建对话会话失败: %w", err)
	}
	
	return nil
}

// GetChatSessionByProjectID 根据项目ID获取对话会话
func (r *MySQLRepository) GetChatSessionByProjectID(projectID uuid.UUID) (*model.ChatSession, error) {
	query := `
		SELECT session_id, project_id, user_id, session_type, started_at, ended_at, status, context
		FROM chat_sessions
		WHERE project_id = ? AND status = 'active'
		ORDER BY started_at DESC
		LIMIT 1
	`
	
	row := r.db.MySQL.QueryRow(query, projectID)
	
	var session model.ChatSession
	var endedAt sql.NullTime
	var context sql.NullString
	
	err := row.Scan(
		&session.SessionID,
		&session.ProjectID,
		&session.UserID,
		&session.SessionType,
		&session.StartedAt,
		&endedAt,
		&session.Status,
		&context,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("对话会话不存在")
		}
		return nil, fmt.Errorf("查询对话会话失败: %w", err)
	}
	
	session.EndedAt = convertNullTime(endedAt)
	if context.Valid {
		session.Context = context.String
	}
	
	return &session, nil
}

// CreateChatMessage 创建对话消息
func (r *MySQLRepository) CreateChatMessage(message *model.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (message_id, session_id, sender_type, message_content, message_type, metadata, timestamp, processed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.MySQL.Exec(query,
		message.MessageID,
		message.SessionID,
		message.SenderType,
		message.MessageContent,
		message.MessageType,
		message.Metadata,
		message.Timestamp,
		message.Processed,
	)
	
	if err != nil {
		return fmt.Errorf("创建对话消息失败: %w", err)
	}
	
	return nil
}

// GetChatMessagesBySessionID 根据会话ID获取对话消息
func (r *MySQLRepository) GetChatMessagesBySessionID(sessionID uuid.UUID, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	offset := (page - 1) * pageSize
	
	// 获取总数
	countQuery := `SELECT COUNT(*) FROM chat_messages WHERE session_id = ?`
	var total int64
	err := r.db.MySQL.QueryRow(countQuery, sessionID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询消息总数失败: %w", err)
	}
	
	// 获取消息列表
	query := `
		SELECT message_id, session_id, sender_type, message_content, message_type, metadata, timestamp, processed
		FROM chat_messages
		WHERE session_id = ?
		ORDER BY timestamp ASC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.MySQL.Query(query, sessionID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询消息列表失败: %w", err)
	}
	defer rows.Close()
	
	var messages []*model.ChatMessage
	for rows.Next() {
		var message model.ChatMessage
		var metadata sql.NullString
		
		err := rows.Scan(
			&message.MessageID,
			&message.SessionID,
			&message.SenderType,
			&message.MessageContent,
			&message.MessageType,
			&metadata,
			&message.Timestamp,
			&message.Processed,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("扫描消息数据失败: %w", err)
		}
		
		if metadata.Valid {
			message.Metadata = metadata.String
		}
		
		messages = append(messages, &message)
	}
	
	return messages, total, nil
}

// EndChatSession 结束对话会话
func (r *MySQLRepository) EndChatSession(sessionID uuid.UUID) error {
	query := `UPDATE chat_sessions SET ended_at = ?, status = 'ended' WHERE session_id = ?`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query, now, sessionID)
	
	if err != nil {
		return fmt.Errorf("结束对话会话失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("对话会话不存在")
	}
	
	return nil
}

// CreateQuestion 创建问题
func (r *MySQLRepository) CreateQuestion(question *model.Question) error {
	query := `
		INSERT INTO questions (question_id, requirement_id, question_text, question_category, priority_level, answer_text, answer_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.MySQL.Exec(query,
		question.QuestionID,
		question.RequirementID,
		question.QuestionText,
		question.QuestionCategory,
		question.PriorityLevel,
		question.AnswerText,
		question.AnswerStatus,
		now,
	)
	
	if err != nil {
		return fmt.Errorf("创建问题失败: %w", err)
	}
	
	return nil
}

// GetQuestionsByRequirementID 根据需求ID获取问题列表
func (r *MySQLRepository) GetQuestionsByRequirementID(requirementID uuid.UUID) ([]*model.Question, error) {
	query := `
		SELECT question_id, requirement_id, question_text, question_category, priority_level, answer_text, answer_status, created_at, answered_at
		FROM questions
		WHERE requirement_id = ?
		ORDER BY priority_level DESC, created_at ASC
	`
	
	rows, err := r.db.MySQL.Query(query, requirementID)
	if err != nil {
		return nil, fmt.Errorf("查询问题列表失败: %w", err)
	}
	defer rows.Close()
	
	var questions []*model.Question
	for rows.Next() {
		var question model.Question
		var answeredAt sql.NullTime
		
		err := rows.Scan(
			&question.QuestionID,
			&question.RequirementID,
			&question.QuestionText,
			&question.QuestionCategory,
			&question.PriorityLevel,
			&question.AnswerText,
			&question.AnswerStatus,
			&question.CreatedAt,
			&answeredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描问题数据失败: %w", err)
		}
		
		question.AnsweredAt = convertNullTime(answeredAt)
		questions = append(questions, &question)
	}
	
	return questions, nil
}

// AnswerQuestion 回答问题
func (r *MySQLRepository) AnswerQuestion(questionID uuid.UUID, answer string) error {
	query := `
		UPDATE questions 
		SET answer_text = ?, answer_status = 'answered', answered_at = ?
		WHERE question_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query, answer, now, questionID)
	
	if err != nil {
		return fmt.Errorf("回答问题失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("问题不存在")
	}
	
	return nil
}

// 以下方法提供简单的实现，可以在需要时进一步完善

// CreatePUMLDiagram 创建PUML图表
func (r *MySQLRepository) CreatePUMLDiagram(diagram *model.PUMLDiagram) error {
	return fmt.Errorf("功能待实现")
}

// GetPUMLDiagramsByProjectID 根据项目ID获取PUML图表
func (r *MySQLRepository) GetPUMLDiagramsByProjectID(projectID uuid.UUID) ([]*model.PUMLDiagram, error) {
	return nil, fmt.Errorf("功能待实现")
}

// UpdatePUMLDiagram 更新PUML图表
func (r *MySQLRepository) UpdatePUMLDiagram(diagram *model.PUMLDiagram) error {
	return fmt.Errorf("功能待实现")
}

// DeletePUMLDiagram 删除PUML图表
func (r *MySQLRepository) DeletePUMLDiagram(diagramID uuid.UUID) error {
	return fmt.Errorf("功能待实现")
}

// CreateDocument 创建文档
func (r *MySQLRepository) CreateDocument(document *model.Document) error {
	return fmt.Errorf("功能待实现")
}

// GetDocumentsByProjectID 根据项目ID获取文档
func (r *MySQLRepository) GetDocumentsByProjectID(projectID uuid.UUID) ([]*model.Document, error) {
	return nil, fmt.Errorf("功能待实现")
}

// UpdateDocument 更新文档
func (r *MySQLRepository) UpdateDocument(document *model.Document) error {
	return fmt.Errorf("功能待实现")
}

// DeleteDocument 删除文档
func (r *MySQLRepository) DeleteDocument(documentID uuid.UUID) error {
	return fmt.Errorf("功能待实现")
}

// CreateBusinessModule 创建业务模块
func (r *MySQLRepository) CreateBusinessModule(module *model.BusinessModule) error {
	return fmt.Errorf("功能待实现")
}

// GetBusinessModulesByProjectID 根据项目ID获取业务模块
func (r *MySQLRepository) GetBusinessModulesByProjectID(projectID uuid.UUID) ([]*model.BusinessModule, error) {
	return nil, fmt.Errorf("功能待实现")
}

// UpdateBusinessModule 更新业务模块
func (r *MySQLRepository) UpdateBusinessModule(module *model.BusinessModule) error {
	return fmt.Errorf("功能待实现")
}

// DeleteBusinessModule 删除业务模块
func (r *MySQLRepository) DeleteBusinessModule(moduleID uuid.UUID) error {
	return fmt.Errorf("功能待实现")
}

// CreateCommonModule 创建通用模块
func (r *MySQLRepository) CreateCommonModule(module *model.CommonModule) error {
	return fmt.Errorf("功能待实现")
}

// GetCommonModulesByCategory 根据分类获取通用模块
func (r *MySQLRepository) GetCommonModulesByCategory(category string, page, pageSize int) ([]*model.CommonModule, int64, error) {
	return nil, 0, fmt.Errorf("功能待实现")
}

// GetCommonModuleByID 根据ID获取通用模块
func (r *MySQLRepository) GetCommonModuleByID(moduleID uuid.UUID) (*model.CommonModule, error) {
	return nil, fmt.Errorf("功能待实现")
}

// UpdateCommonModule 更新通用模块
func (r *MySQLRepository) UpdateCommonModule(module *model.CommonModule) error {
	return fmt.Errorf("功能待实现")
}

// DeleteCommonModule 删除通用模块
func (r *MySQLRepository) DeleteCommonModule(moduleID uuid.UUID) error {
	return fmt.Errorf("功能待实现")
}

// ===== 异步任务相关操作 =====

// CreateAsyncTask 创建异步任务
func (r *MySQLRepository) CreateAsyncTask(task *model.AsyncTask) error {
	query := `
		INSERT INTO async_tasks (task_id, user_id, project_id, task_type, task_name, status, progress, result_data, error_message, created_at, started_at, completed_at, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.MySQL.Exec(query,
		task.TaskID,
		task.UserID,
		task.ProjectID,
		task.TaskType,
		task.TaskName,
		task.Status,
		task.Progress,
		task.ResultData,
		task.ErrorMessage,
		task.CreatedAt,
		convertTimePtr(task.StartedAt),
		convertTimePtr(task.CompletedAt),
		task.Metadata,
	)
	
	if err != nil {
		return fmt.Errorf("创建异步任务失败: %w", err)
	}
	
	return nil
}

// GetAsyncTask 获取异步任务
func (r *MySQLRepository) GetAsyncTask(taskID uuid.UUID) (*model.AsyncTask, error) {
	query := `
		SELECT task_id, user_id, project_id, task_type, task_name, status, progress, result_data, error_message, created_at, started_at, completed_at, metadata
		FROM async_tasks
		WHERE task_id = ?
	`
	
	var task model.AsyncTask
	var startedAt, completedAt sql.NullTime
	
	err := r.db.MySQL.QueryRow(query, taskID).Scan(
		&task.TaskID,
		&task.UserID,
		&task.ProjectID,
		&task.TaskType,
		&task.TaskName,
		&task.Status,
		&task.Progress,
		&task.ResultData,
		&task.ErrorMessage,
		&task.CreatedAt,
		&startedAt,
		&completedAt,
		&task.Metadata,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("任务不存在")
		}
		return nil, fmt.Errorf("获取异步任务失败: %w", err)
	}
	
	task.StartedAt = convertNullTime(startedAt)
	task.CompletedAt = convertNullTime(completedAt)
	
	return &task, nil
}

// UpdateAsyncTask 更新异步任务
func (r *MySQLRepository) UpdateAsyncTask(task *model.AsyncTask) error {
	query := `
		UPDATE async_tasks 
		SET status = ?, progress = ?, result_data = ?, error_message = ?, started_at = ?, completed_at = ?
		WHERE task_id = ?
	`
	
	result, err := r.db.MySQL.Exec(query,
		task.Status,
		task.Progress,
		task.ResultData,
		task.ErrorMessage,
		convertTimePtr(task.StartedAt),
		convertTimePtr(task.CompletedAt),
		task.TaskID,
	)
	
	if err != nil {
		return fmt.Errorf("更新异步任务失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("任务不存在")
	}
	
	return nil
}

// GetTasksByProject 获取项目的任务列表
func (r *MySQLRepository) GetTasksByProject(projectID uuid.UUID, taskType string) ([]*model.AsyncTask, error) {
	var query string
	var args []interface{}
	
	if taskType != "" {
		query = `
			SELECT task_id, user_id, project_id, task_type, task_name, status, progress, result_data, error_message, created_at, started_at, completed_at, metadata
			FROM async_tasks
			WHERE project_id = ? AND task_type = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{projectID, taskType}
	} else {
		query = `
			SELECT task_id, user_id, project_id, task_type, task_name, status, progress, result_data, error_message, created_at, started_at, completed_at, metadata
			FROM async_tasks
			WHERE project_id = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{projectID}
	}
	
	rows, err := r.db.MySQL.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}
	defer rows.Close()
	
	var tasks []*model.AsyncTask
	for rows.Next() {
		var task model.AsyncTask
		var startedAt, completedAt sql.NullTime
		
		err := rows.Scan(
			&task.TaskID,
			&task.UserID,
			&task.ProjectID,
			&task.TaskType,
			&task.TaskName,
			&task.Status,
			&task.Progress,
			&task.ResultData,
			&task.ErrorMessage,
			&task.CreatedAt,
			&startedAt,
			&completedAt,
			&task.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描任务数据失败: %w", err)
		}
		
		task.StartedAt = convertNullTime(startedAt)
		task.CompletedAt = convertNullTime(completedAt)
		tasks = append(tasks, &task)
	}
	
	return tasks, nil
}

// ===== 阶段进度相关操作 =====

// CreateStageProgress 创建阶段进度
func (r *MySQLRepository) CreateStageProgress(progress *model.StageProgress) error {
	query := `
		INSERT INTO stage_progress (progress_id, project_id, stage, status, completion_rate, started_at, completed_at, document_count, puml_count, last_task_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.MySQL.Exec(query,
		progress.ProgressID,
		progress.ProjectID,
		progress.Stage,
		progress.Status,
		progress.CompletionRate,
		convertTimePtr(progress.StartedAt),
		convertTimePtr(progress.CompletedAt),
		progress.DocumentCount,
		progress.PUMLCount,
		progress.LastTaskID,
		now,
		now,
	)
	
	if err != nil {
		return fmt.Errorf("创建阶段进度失败: %w", err)
	}
	
	return nil
}

// GetStageProgress 获取项目的所有阶段进度
func (r *MySQLRepository) GetStageProgress(projectID uuid.UUID) ([]*model.StageProgress, error) {
	query := `
		SELECT progress_id, project_id, stage, status, completion_rate, started_at, completed_at, document_count, puml_count, last_task_id, created_at, updated_at
		FROM stage_progress
		WHERE project_id = ?
		ORDER BY stage ASC
	`
	
	rows, err := r.db.MySQL.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("查询阶段进度失败: %w", err)
	}
	defer rows.Close()
	
	var progresses []*model.StageProgress
	for rows.Next() {
		var progress model.StageProgress
		var startedAt, completedAt sql.NullTime
		var lastTaskID sql.NullString
		
		err := rows.Scan(
			&progress.ProgressID,
			&progress.ProjectID,
			&progress.Stage,
			&progress.Status,
			&progress.CompletionRate,
			&startedAt,
			&completedAt,
			&progress.DocumentCount,
			&progress.PUMLCount,
			&lastTaskID,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描阶段进度数据失败: %w", err)
		}
		
		progress.StartedAt = convertNullTime(startedAt)
		progress.CompletedAt = convertNullTime(completedAt)
		if lastTaskID.Valid {
			taskID, _ := uuid.Parse(lastTaskID.String)
			progress.LastTaskID = &taskID
		}
		
		progresses = append(progresses, &progress)
	}
	
	return progresses, nil
}

// UpdateStageProgress 更新阶段进度
func (r *MySQLRepository) UpdateStageProgress(progress *model.StageProgress) error {
	query := `
		UPDATE stage_progress 
		SET status = ?, completion_rate = ?, started_at = ?, completed_at = ?, document_count = ?, puml_count = ?, last_task_id = ?, updated_at = ?
		WHERE progress_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query,
		progress.Status,
		progress.CompletionRate,
		convertTimePtr(progress.StartedAt),
		convertTimePtr(progress.CompletedAt),
		progress.DocumentCount,
		progress.PUMLCount,
		progress.LastTaskID,
		now,
		progress.ProgressID,
	)
	
	if err != nil {
		return fmt.Errorf("更新阶段进度失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("阶段进度不存在")
	}
	
	return nil
}

// GetStageProgressByStage 获取指定阶段的进度
func (r *MySQLRepository) GetStageProgressByStage(projectID uuid.UUID, stage int) (*model.StageProgress, error) {
	query := `
		SELECT progress_id, project_id, stage, status, completion_rate, started_at, completed_at, document_count, puml_count, last_task_id, created_at, updated_at
		FROM stage_progress
		WHERE project_id = ? AND stage = ?
	`
	
	var progress model.StageProgress
	var startedAt, completedAt sql.NullTime
	var lastTaskID sql.NullString
	
	err := r.db.MySQL.QueryRow(query, projectID, stage).Scan(
		&progress.ProgressID,
		&progress.ProjectID,
		&progress.Stage,
		&progress.Status,
		&progress.CompletionRate,
		&startedAt,
		&completedAt,
		&progress.DocumentCount,
		&progress.PUMLCount,
		&lastTaskID,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("阶段进度不存在")
		}
		return nil, fmt.Errorf("获取阶段进度失败: %w", err)
	}
	
	progress.StartedAt = convertNullTime(startedAt)
	progress.CompletedAt = convertNullTime(completedAt)
	if lastTaskID.Valid {
		taskID, _ := uuid.Parse(lastTaskID.String)
		progress.LastTaskID = &taskID
	}
	
	return &progress, nil
} 