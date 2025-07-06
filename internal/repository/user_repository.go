package repository

import (
	"database/sql"
	"fmt"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
)

// CreateUser 创建用户
func (r *MySQLRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (user_id, username, email, password_hash, full_name, created_at, updated_at, status, preferences)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.MySQL.Exec(query,
		user.UserID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		now,
		now,
		user.Status,
		user.Preferences,
	)
	
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	
	return nil
}

// GetUserByEmail 通过邮箱获取用户
func (r *MySQLRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, full_name, created_at, updated_at, last_login, status, preferences
		FROM users
		WHERE email = ? AND status != 'deleted'
	`
	
	row := r.db.MySQL.QueryRow(query, email)
	
	var user model.User
	var lastLogin sql.NullTime
	var preferences sql.NullString
	
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
		&user.Status,
		&preferences,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	
	user.LastLogin = convertNullTime(lastLogin)
	if preferences.Valid {
		user.Preferences = preferences.String
	}
	
	return &user, nil
}

// GetUserByID 通过ID获取用户
func (r *MySQLRepository) GetUserByID(userID uuid.UUID) (*model.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, full_name, created_at, updated_at, last_login, status, preferences
		FROM users
		WHERE user_id = ? AND status != 'deleted'
	`
	
	row := r.db.MySQL.QueryRow(query, userID)
	
	var user model.User
	var lastLogin sql.NullTime
	var preferences sql.NullString
	
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
		&user.Status,
		&preferences,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	
	user.LastLogin = convertNullTime(lastLogin)
	if preferences.Valid {
		user.Preferences = preferences.String
	}
	
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *MySQLRepository) UpdateUser(user *model.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, full_name = ?, updated_at = ?, preferences = ?
		WHERE user_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query,
		user.Username,
		user.Email,
		user.FullName,
		now,
		user.Preferences,
		user.UserID,
	)
	
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在或未更新")
	}
	
	return nil
}

// UpdateUserLastLogin 更新用户最后登录时间
func (r *MySQLRepository) UpdateUserLastLogin(userID uuid.UUID) error {
	query := `
		UPDATE users 
		SET last_login = ?, updated_at = ?
		WHERE user_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query, now, now, userID)
	
	if err != nil {
		return fmt.Errorf("更新用户登录时间失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}
	
	return nil
}

// CreateProject 创建项目
func (r *MySQLRepository) CreateProject(project *model.Project) error {
	query := `
		INSERT INTO projects (project_id, user_id, project_name, description, project_type, status, created_at, updated_at, completion_percentage, settings)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := r.db.MySQL.Exec(query,
		project.ProjectID,
		project.UserID,
		project.ProjectName,
		project.Description,
		project.ProjectType,
		project.Status,
		now,
		now,
		project.CompletionPercentage,
		project.Settings,
	)
	
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}
	
	return nil
}

// GetProjectsByUserID 获取用户的项目列表
func (r *MySQLRepository) GetProjectsByUserID(userID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	
	// 获取总数
	countQuery := `SELECT COUNT(*) FROM projects WHERE user_id = ? AND status != 'deleted'`
	var total int64
	err := r.db.MySQL.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询项目总数失败: %w", err)
	}
	
	// 获取项目列表
	query := `
		SELECT project_id, user_id, project_name, description, project_type, status, created_at, updated_at, completion_percentage, settings
		FROM projects
		WHERE user_id = ? AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.MySQL.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询项目列表失败: %w", err)
	}
	defer rows.Close()
	
	var projects []*model.Project
	for rows.Next() {
		var project model.Project
		var settings sql.NullString
		
		err := rows.Scan(
			&project.ProjectID,
			&project.UserID,
			&project.ProjectName,
			&project.Description,
			&project.ProjectType,
			&project.Status,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.CompletionPercentage,
			&settings,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("扫描项目数据失败: %w", err)
		}
		
		if settings.Valid {
			project.Settings = settings.String
		}
		
		projects = append(projects, &project)
	}
	
	return projects, total, nil
}

// GetProjectByID 通过ID获取项目
func (r *MySQLRepository) GetProjectByID(projectID uuid.UUID) (*model.Project, error) {
	query := `
		SELECT project_id, user_id, project_name, description, project_type, status, created_at, updated_at, completion_percentage, settings
		FROM projects
		WHERE project_id = ? AND status != 'deleted'
	`
	
	row := r.db.MySQL.QueryRow(query, projectID)
	
	var project model.Project
	var settings sql.NullString
	
	err := row.Scan(
		&project.ProjectID,
		&project.UserID,
		&project.ProjectName,
		&project.Description,
		&project.ProjectType,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.CompletionPercentage,
		&settings,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("项目不存在")
		}
		return nil, fmt.Errorf("查询项目失败: %w", err)
	}
	
	if settings.Valid {
		project.Settings = settings.String
	}
	
	return &project, nil
}

// UpdateProject 更新项目
func (r *MySQLRepository) UpdateProject(project *model.Project) error {
	query := `
		UPDATE projects 
		SET project_name = ?, description = ?, project_type = ?, status = ?, updated_at = ?, completion_percentage = ?, settings = ?
		WHERE project_id = ?
	`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query,
		project.ProjectName,
		project.Description,
		project.ProjectType,
		project.Status,
		now,
		project.CompletionPercentage,
		project.Settings,
		project.ProjectID,
	)
	
	if err != nil {
		return fmt.Errorf("更新项目失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("项目不存在或未更新")
	}
	
	return nil
}

// DeleteProject 删除项目（软删除）
func (r *MySQLRepository) DeleteProject(projectID uuid.UUID) error {
	query := `UPDATE projects SET status = 'deleted', updated_at = ? WHERE project_id = ?`
	
	now := time.Now()
	result, err := r.db.MySQL.Exec(query, now, projectID)
	
	if err != nil {
		return fmt.Errorf("删除项目失败: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("项目不存在")
	}
	
	return nil
}

// Health 健康检查
func (r *MySQLRepository) Health() error {
	return r.db.Health()
} 