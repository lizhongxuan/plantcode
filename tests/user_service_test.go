package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/repository"
	"ai-dev-platform/internal/service"
)

// UserServiceTestSuite 用户服务测试套件
type UserServiceTestSuite struct {
	TestSuite
	userService service.UserService
	repo        repository.Repository
}

// SetupSuite 初始化测试套件
func (s *UserServiceTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	
	// 创建Database实例
	db := &repository.Database{
		GORM: s.DB,
		Redis: nil, // 测试中不使用Redis
	}
	
	// 创建repository
	s.repo = repository.NewMySQLRepository(db)
	
	// 创建service
	s.userService = service.NewUserService(s.repo, s.Config)
}

// TestRegisterUser 测试用户注册
func (s *UserServiceTestSuite) TestRegisterUser() {
	// 准备测试数据
	createReq := &model.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
	}

	// 执行测试
	user, err := s.userService.RegisterUser(createReq)

	// 验证结果
	s.NoError(err)
	s.NotNil(user)
	s.Equal(createReq.Username, user.Username)
	s.Equal(createReq.Email, user.Email)
	s.Equal(createReq.FullName, user.FullName)
	s.Equal(model.UserStatusActive, user.Status)
	s.NotEmpty(user.UserID)
	s.Empty(user.PasswordHash) // 返回时应该清空密码哈希
}

// TestRegisterUserDuplicate 测试注册重复用户
func (s *UserServiceTestSuite) TestRegisterUserDuplicate() {
	// 先创建一个用户
	createReq := &model.CreateUserRequest{
		Username: "duplicateuser",
		Email:    "duplicate@example.com",
		Password: "password123",
		FullName: "Duplicate User",
	}
	
	_, err := s.userService.RegisterUser(createReq)
	s.NoError(err)
	
	// 尝试创建相同邮箱的用户
	duplicateReq := &model.CreateUserRequest{
		Username: "anotheruser",
		Email:    "duplicate@example.com", // 相同的邮箱
		Password: "password123",
		FullName: "Another User",
	}

	// 执行测试
	user, err := s.userService.RegisterUser(duplicateReq)

	// 验证结果
	s.Error(err)
	s.Nil(user)
	s.Contains(err.Error(), "邮箱已被注册")
}

// TestLoginUser 测试用户登录
func (s *UserServiceTestSuite) TestLoginUser() {
	// 先创建一个用户
	createReq := &model.CreateUserRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "loginpassword",
		FullName: "Login User",
	}
	
	createdUser, err := s.userService.RegisterUser(createReq)
	s.NoError(err)
	
	// 准备登录数据
	loginReq := &model.LoginRequest{
		Email:    "login@example.com",
		Password: "loginpassword",
	}

	// 执行测试
	loginResp, err := s.userService.LoginUser(loginReq)

	// 验证结果
	s.NoError(err)
	s.NotNil(loginResp)
	s.Equal(createdUser.UserID, loginResp.User.UserID)
	s.Equal(createdUser.Username, loginResp.User.Username)
	s.NotEmpty(loginResp.Token)
	s.Empty(loginResp.User.PasswordHash) // 返回时应该清空密码哈希
}

// TestLoginUserInvalidCredentials 测试无效凭证登录
func (s *UserServiceTestSuite) TestLoginUserInvalidCredentials() {
	// 准备测试数据
	loginReq := &model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	// 执行测试
	loginResp, err := s.userService.LoginUser(loginReq)

	// 验证结果
	s.Error(err)
	s.Nil(loginResp)
	s.Contains(err.Error(), "用户不存在或密码错误")
}

// TestGetUser 测试获取用户信息
func (s *UserServiceTestSuite) TestGetUser() {
	// 先创建一个用户
	createReq := &model.CreateUserRequest{
		Username: "getuser",
		Email:    "getuser@example.com",
		Password: "password123",
		FullName: "Get User",
	}
	
	createdUser, err := s.userService.RegisterUser(createReq)
	s.NoError(err)

	// 执行测试
	user, err := s.userService.GetUser(createdUser.UserID)

	// 验证结果
	s.NoError(err)
	s.NotNil(user)
	s.Equal(createdUser.UserID, user.UserID)
	s.Equal(createdUser.Username, user.Username)
	s.Equal(createdUser.Email, user.Email)
	s.Empty(user.PasswordHash) // 返回时应该清空密码哈希
}

// TestGetUserNotFound 测试获取不存在的用户
func (s *UserServiceTestSuite) TestGetUserNotFound() {
	// 准备测试数据
	nonExistentID := uuid.New()

	// 执行测试
	user, err := s.userService.GetUser(nonExistentID)

	// 验证结果
	s.Error(err)
	s.Nil(user)
	s.Contains(err.Error(), "获取用户失败")
}

// TestUpdateUser 测试更新用户信息
func (s *UserServiceTestSuite) TestUpdateUser() {
	// 先创建一个用户
	createReq := &model.CreateUserRequest{
		Username: "updateuser",
		Email:    "updateuser@example.com",
		Password: "password123",
		FullName: "Update User",
	}
	
	createdUser, err := s.userService.RegisterUser(createReq)
	s.NoError(err)
	
	// 准备更新数据
	newFullName := "Updated User Name"
	newEmail := "updated@example.com"
	updateReq := &service.UserUpdateRequest{
		FullName: &newFullName,
		Email:    &newEmail,
	}

	// 执行测试
	user, err := s.userService.UpdateUser(createdUser.UserID, updateReq)

	// 验证结果
	s.NoError(err)
	s.NotNil(user)
	s.Equal(newFullName, user.FullName)
	s.Equal(newEmail, user.Email)
	s.Equal(createdUser.UserID, user.UserID)
}

// TestValidateToken 测试验证Token
func (s *UserServiceTestSuite) TestValidateToken() {
	// 先创建一个用户并登录
	createReq := &model.CreateUserRequest{
		Username: "tokenuser",
		Email:    "tokenuser@example.com",
		Password: "password123",
		FullName: "Token User",
	}
	
	createdUser, err := s.userService.RegisterUser(createReq)
	s.NoError(err)
	
	// 登录获取token
	loginReq := &model.LoginRequest{
		Email:    "tokenuser@example.com",
		Password: "password123",
	}
	
	loginResp, err := s.userService.LoginUser(loginReq)
	s.NoError(err)
	s.NotEmpty(loginResp.Token)
	
	// 验证token
	user, err := s.userService.ValidateToken(loginResp.Token)
	
	// 验证结果
	s.NoError(err)
	s.NotNil(user)
	s.Equal(createdUser.UserID, user.UserID)
	s.Equal(createdUser.Username, user.Username)
	s.Empty(user.PasswordHash) // 返回时应该清空密码哈希
}

// TestValidateInvalidToken 测试验证无效Token
func (s *UserServiceTestSuite) TestValidateInvalidToken() {
	// 准备无效token
	invalidToken := "invalid.token.here"

	// 执行测试
	user, err := s.userService.ValidateToken(invalidToken)

	// 验证结果
	s.Error(err)
	s.Nil(user)
	s.Contains(err.Error(), "无效的令牌")
}

// TestSuite 运行器
func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
} 