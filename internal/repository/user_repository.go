package repository

import (
	"fmt"
	"time"

	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateUser 创建用户
func (r *MySQLRepository) CreateUser(user *model.User) error {
	if r.db.GORM == nil {
		return fmt.Errorf("数据库连接不可用，请检查MySQL服务状态")
	}

	// 设置创建和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := r.db.GORM.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// GetUserByEmail 通过邮箱获取用户
func (r *MySQLRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	
	err := r.db.GORM.Where("email = ? AND status != ?", email, "deleted").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// GetUserByID 通过ID获取用户
func (r *MySQLRepository) GetUserByID(userID uuid.UUID) (*model.User, error) {
	var user model.User
	
	err := r.db.GORM.Where("user_id = ? AND status != ?", userID, "deleted").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (r *MySQLRepository) UpdateUser(user *model.User) error {
	// 设置更新时间
	user.UpdatedAt = time.Now()

	result := r.db.GORM.Model(user).Where("user_id = ?", user.UserID).Updates(map[string]interface{}{
		"username":     user.Username,
		"email":        user.Email,
		"full_name":    user.FullName,
		"updated_at":   user.UpdatedAt,
		"preferences":  user.Preferences,
	})

	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或未更新")
	}

	return nil
}

// UpdateUserLastLogin 更新用户最后登录时间
func (r *MySQLRepository) UpdateUserLastLogin(userID uuid.UUID) error {
	now := time.Now()
	
	result := r.db.GORM.Model(&model.User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"last_login":  &now,
		"updated_at":  now,
	})

	if result.Error != nil {
		return fmt.Errorf("更新用户登录时间失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// CreateProject 创建项目
func (r *MySQLRepository) CreateProject(project *model.Project) error {
	// 设置创建和更新时间
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	if err := r.db.GORM.Create(project).Error; err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}

	return nil
}

// GetProjectsByUserID 获取用户的项目列表
func (r *MySQLRepository) GetProjectsByUserID(userID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.GORM.Model(&model.Project{}).Where("user_id = ? AND status != ?", userID, "deleted").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询项目总数失败: %w", err)
	}

	// 获取项目列表
	if err := r.db.GORM.Where("user_id = ? AND status != ?", userID, "deleted").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("查询项目列表失败: %w", err)
	}

	return projects, total, nil
}

// GetProjectByID 通过ID获取项目
func (r *MySQLRepository) GetProjectByID(projectID uuid.UUID) (*model.Project, error) {
	var project model.Project
	
	err := r.db.GORM.Where("project_id = ? AND status != ?", projectID, "deleted").First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("项目不存在")
		}
		return nil, fmt.Errorf("查询项目失败: %w", err)
	}

	return &project, nil
}

// UpdateProject 更新项目
func (r *MySQLRepository) UpdateProject(project *model.Project) error {
	// 设置更新时间
	project.UpdatedAt = time.Now()

	result := r.db.GORM.Model(project).Where("project_id = ?", project.ProjectID).Updates(map[string]interface{}{
		"project_name":          project.ProjectName,
		"description":           project.Description,
		"project_type":          project.ProjectType,
		"status":               project.Status,
		"updated_at":           project.UpdatedAt,
		"completion_percentage": project.CompletionPercentage,
		"settings":             project.Settings,
	})

	if result.Error != nil {
		return fmt.Errorf("更新项目失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("项目不存在或未更新")
	}

	return nil
}

// DeleteProject 删除项目（软删除）
func (r *MySQLRepository) DeleteProject(projectID uuid.UUID) error {
	now := time.Now()
	
	result := r.db.GORM.Model(&model.Project{}).Where("project_id = ?", projectID).Updates(map[string]interface{}{
		"status":     "deleted",
		"updated_at": now,
	})

	if result.Error != nil {
		return fmt.Errorf("删除项目失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("项目不存在")
	}

	return nil
}

// Health 健康检查
func (r *MySQLRepository) Health() error {
	return r.db.Health()
}