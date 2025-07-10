package service

import (
	"fmt"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/repository"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
)

// UserService 用户服务接口
type UserService interface {
	// 注册用户
	RegisterUser(req *model.CreateUserRequest) (*model.User, error)
	
	// 用户登录
	LoginUser(req *model.LoginRequest) (*LoginResponse, error)
	
	// 验证JWT令牌
	ValidateToken(token string) (*model.User, error)
	
	// 获取用户信息
	GetUser(userID uuid.UUID) (*model.User, error)
	
	// 更新用户信息
	UpdateUser(userID uuid.UUID, updates *UserUpdateRequest) (*model.User, error)
	
	// 更新用户最后登录时间
	UpdateLastLogin(userID uuid.UUID) error
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Username    *string `json:"username,omitempty"`
	Email       *string `json:"email,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
	Preferences *string `json:"preferences,omitempty"`
}

// userService 用户服务实现
type userService struct {
	repo   repository.Repository
	config *config.Config
}

// NewUserService 创建用户服务
func NewUserService(repo repository.Repository, cfg *config.Config) UserService {
	return &userService{
		repo:   repo,
		config: cfg,
	}
}

// RegisterUser 注册用户
func (s *userService) RegisterUser(req *model.CreateUserRequest) (*model.User, error) {
	// 验证输入
	if err := s.validateUserInput(req); err != nil {
		return nil, fmt.Errorf("输入验证失败: %w", err)
	}

	// 检查邮箱是否已存在
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	// 创建用户模型
	user := &model.User{
		UserID:       utils.GenerateUUID(),
		Username:     utils.SanitizeString(req.Username),
		Email:        utils.SanitizeString(req.Email),
		PasswordHash: hashedPassword,
		FullName:     utils.SanitizeString(req.FullName),
		Status:       model.UserStatusActive,
		Preferences:  `{}`, // 默认空JSON
	}

	// 保存到数据库
	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 不返回密码哈希
	user.PasswordHash = ""
	return user, nil
}

// LoginUser 用户登录
func (s *userService) LoginUser(req *model.LoginRequest) (*LoginResponse, error) {
	// 验证输入
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("邮箱和密码不能为空")
	}

	if !utils.ValidateEmail(req.Email) {
		return nil, fmt.Errorf("邮箱格式无效")
	}

	// 查找用户
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, fmt.Errorf("账户已被停用，请联系管理员")
	}

	// 验证密码
	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(
		user.UserID,
		user.Username,
		user.Email,
		s.config.JWT.Secret,
		s.config.JWT.ExpiresIn,
	)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 更新最后登录时间
	if err := s.repo.UpdateUserLastLogin(user.UserID); err != nil {
		// 记录错误但不影响登录
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}

	// 不返回密码哈希
	user.PasswordHash = ""

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// ValidateToken 验证JWT令牌
func (s *userService) ValidateToken(token string) (*model.User, error) {
	// 验证JWT令牌
	claims, err := utils.ValidateJWT(token, s.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("无效的令牌: %w", err)
	}

	// 获取用户信息
	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, fmt.Errorf("账户已被停用")
	}

	// 不返回密码哈希
	user.PasswordHash = ""
	return user, nil
}

// GetUser 获取用户信息
func (s *userService) GetUser(userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 不返回密码哈希
	user.PasswordHash = ""
	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(userID uuid.UUID, updates *UserUpdateRequest) (*model.User, error) {
	// 获取当前用户信息
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 应用更新
	if updates.Username != nil {
		username := utils.SanitizeString(*updates.Username)
		if len(username) < 3 || len(username) > 50 {
			return nil, fmt.Errorf("用户名长度必须在3-50个字符之间")
		}
		user.Username = username
	}

	if updates.Email != nil {
		email := utils.SanitizeString(*updates.Email)
		if !utils.ValidateEmail(email) {
			return nil, fmt.Errorf("邮箱格式无效")
		}
		
		// 检查邮箱是否已被其他用户使用
		if email != user.Email {
			existingUser, _ := s.repo.GetUserByEmail(email)
			if existingUser != nil && existingUser.UserID != userID {
				return nil, fmt.Errorf("邮箱已被其他用户使用")
			}
		}
		user.Email = email
	}

	if updates.FullName != nil {
		fullName := utils.SanitizeString(*updates.FullName)
		if len(fullName) < 2 || len(fullName) > 100 {
			return nil, fmt.Errorf("姓名长度必须在2-100个字符之间")
		}
		user.FullName = fullName
	}

	if updates.Preferences != nil {
		user.Preferences = *updates.Preferences
	}

	// 更新到数据库
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 不返回密码哈希
	user.PasswordHash = ""
	return user, nil
}

// UpdateLastLogin 更新用户最后登录时间
func (s *userService) UpdateLastLogin(userID uuid.UUID) error {
	return s.repo.UpdateUserLastLogin(userID)
}

// validateUserInput 验证用户输入
func (s *userService) validateUserInput(req *model.CreateUserRequest) error {
	// 验证用户名
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("用户名长度必须在3-50个字符之间")
	}

	// 验证邮箱
	if !utils.ValidateEmail(req.Email) {
		return fmt.Errorf("邮箱格式无效")
	}

	// 验证密码
	if len(req.Password) < 8 {
		return fmt.Errorf("密码长度不能少于8个字符")
	}

	// 验证姓名
	if len(req.FullName) < 2 || len(req.FullName) > 100 {
		return fmt.Errorf("姓名长度必须在2-100个字符之间")
	}

	return nil
}

// ProjectService 项目服务接口
type ProjectService interface {
	// 创建项目
	CreateProject(userID uuid.UUID, req *model.CreateProjectRequest) (*model.Project, error)
	
	// 获取用户项目列表
	GetUserProjects(userID uuid.UUID, page, pageSize int) ([]*model.Project, utils.PaginationInfo, error)
	
	// 获取项目详情
	GetProject(projectID uuid.UUID, userID uuid.UUID) (*model.Project, error)
	
	// 更新项目
	UpdateProject(projectID uuid.UUID, userID uuid.UUID, updates *ProjectUpdateRequest) (*model.Project, error)
	
	// 删除项目
	DeleteProject(projectID uuid.UUID, userID uuid.UUID) error
}

// ProjectUpdateRequest 项目更新请求
type ProjectUpdateRequest struct {
	ProjectName           *string `json:"project_name,omitempty"`
	Description           *string `json:"description,omitempty"`
	ProjectType           *string `json:"project_type,omitempty"`
	Status                *string `json:"status,omitempty"`
	CompletionPercentage  *int    `json:"completion_percentage,omitempty"`
	Settings              *string `json:"settings,omitempty"`
}

// projectService 项目服务实现
type projectService struct {
	repo repository.Repository
}

// NewProjectService 创建项目服务
func NewProjectService(repo repository.Repository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

// CreateProject 创建项目
func (s *projectService) CreateProject(userID uuid.UUID, req *model.CreateProjectRequest) (*model.Project, error) {
	// 验证输入
	if err := s.validateProjectInput(req); err != nil {
		return nil, fmt.Errorf("输入验证失败: %w", err)
	}

	// 创建项目模型
	project := &model.Project{
		ProjectID:            utils.GenerateUUID(),
		UserID:               userID,
		ProjectName:          utils.SanitizeString(req.ProjectName),
		Description:          utils.SanitizeString(req.Description),
		ProjectType:          req.ProjectType,
		Status:               model.ProjectStatusDraft,
		CompletionPercentage: 0,
		Settings:             `{}`, // 默认空JSON
	}

	// 保存到数据库
	if err := s.repo.CreateProject(project); err != nil {
		return nil, fmt.Errorf("创建项目失败: %w", err)
	}

	return project, nil
}

// GetUserProjects 获取用户项目列表
func (s *projectService) GetUserProjects(userID uuid.UUID, page, pageSize int) ([]*model.Project, utils.PaginationInfo, error) {
	projects, total, err := s.repo.GetProjectsByUserID(userID, page, pageSize)
	if err != nil {
		return nil, utils.PaginationInfo{}, fmt.Errorf("获取项目列表失败: %w", err)
	}

	pagination := utils.CalculatePagination(page, pageSize, total)
	return projects, pagination, nil
}

// GetProject 获取项目详情
func (s *projectService) GetProject(projectID uuid.UUID, userID uuid.UUID) (*model.Project, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}

	// 验证项目所有权
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问此项目")
	}

	return project, nil
}

// UpdateProject 更新项目
func (s *projectService) UpdateProject(projectID uuid.UUID, userID uuid.UUID, updates *ProjectUpdateRequest) (*model.Project, error) {
	// 获取当前项目信息
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	// 验证项目所有权
	if project.UserID != userID {
		return nil, fmt.Errorf("无权修改此项目")
	}

	// 应用更新
	if updates.ProjectName != nil {
		name := utils.SanitizeString(*updates.ProjectName)
		if len(name) < 1 || len(name) > 100 {
			return nil, fmt.Errorf("项目名称长度必须在1-100个字符之间")
		}
		project.ProjectName = name
	}

	if updates.Description != nil {
		description := utils.SanitizeString(*updates.Description)
		if len(description) > 1000 {
			return nil, fmt.Errorf("项目描述不能超过1000个字符")
		}
		project.Description = description
	}

	if updates.ProjectType != nil {
		if !utils.IsValidProjectType(*updates.ProjectType) {
			return nil, fmt.Errorf("无效的项目类型")
		}
		project.ProjectType = *updates.ProjectType
	}

	if updates.Status != nil {
		project.Status = *updates.Status
	}

	if updates.CompletionPercentage != nil {
		percentage := *updates.CompletionPercentage
		if percentage < 0 || percentage > 100 {
			return nil, fmt.Errorf("完成百分比必须在0-100之间")
		}
		project.CompletionPercentage = percentage
	}

	if updates.Settings != nil {
		project.Settings = *updates.Settings
	}

	// 更新到数据库
	if err := s.repo.UpdateProject(project); err != nil {
		return nil, fmt.Errorf("更新项目失败: %w", err)
	}

	return project, nil
}

// DeleteProject 删除项目
func (s *projectService) DeleteProject(projectID uuid.UUID, userID uuid.UUID) error {
	// 获取项目信息验证所有权
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return fmt.Errorf("项目不存在: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("无权删除此项目")
	}

	// 执行软删除
	if err := s.repo.DeleteProject(projectID); err != nil {
		return fmt.Errorf("删除项目失败: %w", err)
	}

	return nil
}

// validateProjectInput 验证项目输入
func (s *projectService) validateProjectInput(req *model.CreateProjectRequest) error {
	// 验证项目名称
	if len(req.ProjectName) < 1 || len(req.ProjectName) > 100 {
		return fmt.Errorf("项目名称长度必须在1-100个字符之间")
	}

	// 验证描述
	if len(req.Description) > 1000 {
		return fmt.Errorf("项目描述不能超过1000个字符")
	}

	// 验证项目类型
	if !utils.IsValidProjectType(req.ProjectType) {
		return fmt.Errorf("无效的项目类型")
	}

	return nil
} 