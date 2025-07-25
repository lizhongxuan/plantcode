package service

import (
	"fmt"
	"testing"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockRepository 模拟Repository
type MockRepository struct {
	mock.Mock
}

// 实现Repository接口的用户相关方法
func (m *MockRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(userID uuid.UUID) (*model.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) UpdateUserLastLogin(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// 实现Repository接口的项目相关方法
func (m *MockRepository) CreateProject(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockRepository) GetProjectsByUserID(userID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetProjectByID(projectID uuid.UUID) (*model.Project, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockRepository) UpdateProject(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockRepository) DeleteProject(projectID uuid.UUID) error {
	args := m.Called(projectID)
	return args.Error(0)
}

// 实现Repository接口的其他方法（空实现）
func (m *MockRepository) CreateRequirementAnalysis(requirement *model.Requirement) error {
	return nil
}
func (m *MockRepository) GetRequirementByProjectID(projectID uuid.UUID) (*model.Requirement, error) {
	return nil, nil
}
func (m *MockRepository) UpdateRequirementAnalysis(requirement *model.Requirement) error {
	return nil
}
func (m *MockRepository) CreateChatSession(session *model.ChatSession) error { return nil }
func (m *MockRepository) GetChatSessionByProjectID(projectID uuid.UUID) (*model.ChatSession, error) {
	return nil, nil
}
func (m *MockRepository) CreateChatMessage(message *model.ChatMessage) error { return nil }
func (m *MockRepository) GetChatMessagesBySessionID(sessionID uuid.UUID, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	return nil, 0, nil
}
func (m *MockRepository) EndChatSession(sessionID uuid.UUID) error                        { return nil }
func (m *MockRepository) CreateQuestion(question *model.Question) error                  { return nil }
func (m *MockRepository) GetQuestionsByRequirementID(requirementID uuid.UUID) ([]*model.Question, error) {
	return nil, nil
}
func (m *MockRepository) AnswerQuestion(questionID uuid.UUID, answer string) error { return nil }
func (m *MockRepository) CreatePUMLDiagram(diagram *model.PUMLDiagram) error       { return nil }
func (m *MockRepository) GetPUMLDiagramsByProjectID(projectID uuid.UUID) ([]*model.PUMLDiagram, error) {
	return nil, nil
}
func (m *MockRepository) UpdatePUMLDiagram(diagram *model.PUMLDiagram) error { return nil }
func (m *MockRepository) DeletePUMLDiagram(diagramID uuid.UUID) error       { return nil }
func (m *MockRepository) CreateDocument(document *model.Document) error     { return nil }
func (m *MockRepository) GetDocumentsByProjectID(projectID uuid.UUID) ([]*model.Document, error) {
	return nil, nil
}
func (m *MockRepository) UpdateDocument(document *model.Document) error { return nil }
func (m *MockRepository) DeleteDocument(documentID uuid.UUID) error     { return nil }
func (m *MockRepository) CreateBusinessModule(module *model.BusinessModule) error {
	return nil
}
func (m *MockRepository) GetBusinessModulesByProjectID(projectID uuid.UUID) ([]*model.BusinessModule, error) {
	return nil, nil
}
func (m *MockRepository) UpdateBusinessModule(module *model.BusinessModule) error {
	return nil
}
func (m *MockRepository) DeleteBusinessModule(moduleID uuid.UUID) error { return nil }
func (m *MockRepository) CreateCommonModule(module *model.CommonModule) error {
	return nil
}
func (m *MockRepository) GetCommonModulesByCategory(category string, page, pageSize int) ([]*model.CommonModule, int64, error) {
	return nil, 0, nil
}
func (m *MockRepository) GetCommonModuleByID(moduleID uuid.UUID) (*model.CommonModule, error) {
	return nil, nil
}
func (m *MockRepository) UpdateCommonModule(module *model.CommonModule) error {
	return nil
}
func (m *MockRepository) DeleteCommonModule(moduleID uuid.UUID) error        { return nil }
func (m *MockRepository) CreateAsyncTask(task *model.AsyncTask) error        { return nil }
func (m *MockRepository) GetAsyncTask(taskID uuid.UUID) (*model.AsyncTask, error) {
	return nil, nil
}
func (m *MockRepository) UpdateAsyncTask(task *model.AsyncTask) error { return nil }
func (m *MockRepository) GetTasksByProject(projectID uuid.UUID, taskType string) ([]*model.AsyncTask, error) {
	return nil, nil
}
func (m *MockRepository) CreateStageProgress(progress *model.StageProgress) error {
	return nil
}
func (m *MockRepository) GetStageProgress(projectID uuid.UUID) ([]*model.StageProgress, error) {
	return nil, nil
}
func (m *MockRepository) UpdateStageProgress(progress *model.StageProgress) error {
	return nil
}
func (m *MockRepository) GetStageProgressByStage(projectID uuid.UUID, stage int) (*model.StageProgress, error) {
	return nil, nil
}

// UserAIConfig 相关方法
func (m *MockRepository) GetUserAIConfig(userID uuid.UUID) (*model.UserAIConfig, error) {
	return nil, nil
}
func (m *MockRepository) CreateUserAIConfig(config *model.UserAIConfig) error {
	return nil
}
func (m *MockRepository) UpdateUserAIConfig(config *model.UserAIConfig) error {
	return nil
}
func (m *MockRepository) DeleteUserAIConfig(userID uuid.UUID) error {
	return nil
}

// 扩展方法（用于兼容性）
func (m *MockRepository) GetRequirementAnalysis(analysisID uuid.UUID) (*model.Requirement, error) {
	return nil, nil
}
func (m *MockRepository) GetRequirementAnalysesByProject(projectID uuid.UUID) ([]*model.Requirement, error) {
	return nil, nil
}
func (m *MockRepository) GetChatSession(sessionID uuid.UUID) (*model.ChatSession, error) {
	return nil, nil
}
func (m *MockRepository) GetChatSessionsByProject(projectID uuid.UUID) ([]*model.ChatSession, error) {
	return nil, nil
}
func (m *MockRepository) GetChatMessages(sessionID uuid.UUID) ([]*model.ChatMessage, error) {
	return nil, nil
}
func (m *MockRepository) GetPUMLDiagram(diagramID uuid.UUID) (*model.PUMLDiagram, error) {
	return nil, nil
}
func (m *MockRepository) GetDocument(documentID uuid.UUID) (*model.Document, error) {
	return nil, nil
}
func (m *MockRepository) GetQuestions(requirementID uuid.UUID) ([]*model.Question, error) {
	return nil, nil
}

func (m *MockRepository) Health() error { return nil }

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo    *MockRepository
	userService UserService
	cfg         *config.Config
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 3600,
		},
	}
	suite.userService = NewUserService(suite.mockRepo, suite.cfg)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestRegisterUser_Success() {
	// Arrange
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	// Mock repository expectations
	suite.mockRepo.On("GetUserByEmail", req.Email).Return(nil, fmt.Errorf("用户不存在"))
	suite.mockRepo.On("CreateUser", mock.MatchedBy(func(user *model.User) bool {
		return user.Username == req.Username &&
			user.Email == req.Email &&
			user.FullName == req.FullName &&
			user.Status == model.UserStatusActive &&
			user.PasswordHash != "" &&
			user.PasswordHash != req.Password // 密码应该被哈希
	})).Return(nil)

	// Act
	user, err := suite.userService.RegisterUser(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), req.Username, user.Username)
	assert.Equal(suite.T(), req.Email, user.Email)
	assert.Equal(suite.T(), req.FullName, user.FullName)
	assert.Equal(suite.T(), model.UserStatusActive, user.Status)
	assert.Equal(suite.T(), "", user.PasswordHash) // 不应该返回密码哈希
	assert.NotEqual(suite.T(), uuid.Nil, user.UserID)
}

func (suite *UserServiceTestSuite) TestRegisterUser_EmailAlreadyExists() {
	// Arrange
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	existingUser := &model.User{
		UserID: uuid.New(),
		Email:  req.Email,
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil)

	// Act
	user, err := suite.userService.RegisterUser(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "邮箱已被注册")
}

func (suite *UserServiceTestSuite) TestRegisterUser_InvalidInput() {
	testCases := []struct {
		name string
		req  *model.CreateUserRequest
	}{
		{
			name: "empty username",
			req: &model.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
		},
		{
			name: "invalid email",
			req: &model.CreateUserRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				FullName: "Test User",
			},
		},
		{
			name: "short password",
			req: &model.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "123",
				FullName: "Test User",
			},
		},
		{
			name: "empty full name",
			req: &model.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Act
			user, err := suite.userService.RegisterUser(tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Contains(t, err.Error(), "输入验证失败")
		})
	}
}

func (suite *UserServiceTestSuite) TestRegisterUser_DatabaseError() {
	// Arrange
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(nil, fmt.Errorf("用户不存在"))
	suite.mockRepo.On("CreateUser", mock.Anything).Return(fmt.Errorf("数据库错误"))

	// Act
	user, err := suite.userService.RegisterUser(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "创建用户失败")
}

func (suite *UserServiceTestSuite) TestLoginUser_Success() {
	// Arrange
	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	existingUser := &model.User{
		UserID:       uuid.New(),
		Username:     "testuser",
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     "Test User",
		Status:       model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil)
	suite.mockRepo.On("UpdateUserLastLogin", existingUser.UserID).Return(nil)

	// Act
	response, err := suite.userService.LoginUser(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.NotNil(suite.T(), response.User)
	assert.NotEmpty(suite.T(), response.Token)
	assert.Equal(suite.T(), existingUser.UserID, response.User.UserID)
	assert.Equal(suite.T(), existingUser.Username, response.User.Username)
	assert.Equal(suite.T(), existingUser.Email, response.User.Email)
	assert.Equal(suite.T(), "", response.User.PasswordHash) // 不应该返回密码哈希
}

func (suite *UserServiceTestSuite) TestLoginUser_InvalidCredentials() {
	testCases := []struct {
		name string
		req  *model.LoginRequest
	}{
		{
			name: "empty email",
			req: &model.LoginRequest{
				Email:    "",
				Password: "password123",
			},
		},
		{
			name: "empty password",
			req: &model.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
		},
		{
			name: "invalid email format",
			req: &model.LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Act
			response, err := suite.userService.LoginUser(tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, response)
		})
	}
}

func (suite *UserServiceTestSuite) TestLoginUser_UserNotFound() {
	// Arrange
	req := &model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(nil, fmt.Errorf("用户不存在"))

	// Act
	response, err := suite.userService.LoginUser(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "用户不存在或密码错误")
}

func (suite *UserServiceTestSuite) TestLoginUser_WrongPassword() {
	// Arrange
	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hashedPassword, _ := utils.HashPassword("correctpassword")
	existingUser := &model.User{
		UserID:       uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil)

	// Act
	response, err := suite.userService.LoginUser(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "用户不存在或密码错误")
}

func (suite *UserServiceTestSuite) TestLoginUser_InactiveUser() {
	// Arrange
	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	existingUser := &model.User{
		UserID:       uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       model.UserStatusInactive, // 非活跃状态
	}

	suite.mockRepo.On("GetUserByEmail", req.Email).Return(existingUser, nil)

	// Act
	response, err := suite.userService.LoginUser(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "账户已被停用")
}

func (suite *UserServiceTestSuite) TestValidateToken_Success() {
	// Arrange
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"

	// 生成有效token
	token, err := utils.GenerateJWT(userID, username, email, suite.cfg.JWT.Secret, suite.cfg.JWT.ExpiresIn)
	assert.NoError(suite.T(), err)

	existingUser := &model.User{
		UserID:   userID,
		Username: username,
		Email:    email,
		Status:   model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	
	// Act
	user, err := suite.userService.ValidateToken(token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.Equal(suite.T(), username, user.Username)
	assert.Equal(suite.T(), email, user.Email)
	assert.Equal(suite.T(), "", user.PasswordHash)
}

func (suite *UserServiceTestSuite) TestValidateToken_InvalidToken() {
	// Arrange
	invalidToken := "invalid.token.format"

	// Act
	user, err := suite.userService.ValidateToken(invalidToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "无效的令牌")
}

func (suite *UserServiceTestSuite) TestValidateToken_UserNotFound() {
	// Arrange
	userID := uuid.New()
	token, err := utils.GenerateJWT(userID, "testuser", "test@example.com", suite.cfg.JWT.Secret, suite.cfg.JWT.ExpiresIn)
	assert.NoError(suite.T(), err)

	suite.mockRepo.On("GetUserByID", userID).Return(nil, fmt.Errorf("用户不存在"))

	// Act
	user, err := suite.userService.ValidateToken(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
}

func (suite *UserServiceTestSuite) TestValidateToken_InactiveUser() {
	// Arrange
	userID := uuid.New()
	token, err := utils.GenerateJWT(userID, "testuser", "test@example.com", suite.cfg.JWT.Secret, suite.cfg.JWT.ExpiresIn)
	assert.NoError(suite.T(), err)

	existingUser := &model.User{
		UserID: userID,
		Status: model.UserStatusInactive,
	}

	suite.mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	
	// Act
	user, err := suite.userService.ValidateToken(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "账户已被停用")
}

func (suite *UserServiceTestSuite) TestGetUser_Success() {
	// Arrange
	userID := uuid.New()
	existingUser := &model.User{
		UserID:       userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		Status:       model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	
	// Act
	user, err := suite.userService.GetUser(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.UserID)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), "", user.PasswordHash) // 不应该返回密码哈希
}

func (suite *UserServiceTestSuite) TestGetUser_NotFound() {
	// Arrange
	userID := uuid.New()
	suite.mockRepo.On("GetUserByID", userID).Return(nil, fmt.Errorf("用户不存在"))

	// Act
	user, err := suite.userService.GetUser(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "获取用户失败")
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	// Arrange
	userID := uuid.New()
	newUsername := "newusername"
	newEmail := "new@example.com"
	newFullName := "New Full Name"
	newPreferences := `{"theme": "dark"}`

	updates := &UserUpdateRequest{
		Username:    &newUsername,
		Email:       &newEmail,
		FullName:    &newFullName,
		Preferences: &newPreferences,
	}

	existingUser := &model.User{
		UserID:   userID,
		Username: "oldusername",
		Email:    "old@example.com",
		FullName: "Old Full Name",
		Status:   model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	// Mock GetUserByEmail 检查新邮箱是否被其他用户使用
	suite.mockRepo.On("GetUserByEmail", newEmail).Return(nil, fmt.Errorf("用户不存在"))
	suite.mockRepo.On("UpdateUser", mock.MatchedBy(func(user *model.User) bool {
		return user.UserID == userID &&
			user.Username == newUsername &&
			user.Email == newEmail &&
			user.FullName == newFullName &&
			user.Preferences == newPreferences
	})).Return(nil)

	// Act
	user, err := suite.userService.UpdateUser(userID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), newUsername, user.Username)
	assert.Equal(suite.T(), newEmail, user.Email)
	assert.Equal(suite.T(), newFullName, user.FullName)
	assert.Equal(suite.T(), newPreferences, user.Preferences)
	assert.Equal(suite.T(), "", user.PasswordHash)
}

func (suite *UserServiceTestSuite) TestUpdateUser_PartialUpdate() {
	// Arrange
	userID := uuid.New()
	newUsername := "newusername"

	updates := &UserUpdateRequest{
		Username: &newUsername,
		// 其他字段为nil，应该保持原值
	}

	existingUser := &model.User{
		UserID:      userID,
		Username:    "oldusername",
		Email:       "test@example.com",
		FullName:    "Test User",
		Preferences: `{"theme": "light"}`,
		Status:      model.UserStatusActive,
	}

	suite.mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	suite.mockRepo.On("UpdateUser", mock.MatchedBy(func(user *model.User) bool {
		return user.UserID == userID &&
			user.Username == newUsername &&
			user.Email == existingUser.Email && // 应该保持原值
			user.FullName == existingUser.FullName && // 应该保持原值
			user.Preferences == existingUser.Preferences // 应该保持原值
	})).Return(nil)

	// Act
	user, err := suite.userService.UpdateUser(userID, updates)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), newUsername, user.Username)
	assert.Equal(suite.T(), existingUser.Email, user.Email)
	assert.Equal(suite.T(), existingUser.FullName, user.FullName)
	assert.Equal(suite.T(), existingUser.Preferences, user.Preferences)
}

func (suite *UserServiceTestSuite) TestUpdateUser_UserNotFound() {
	// Arrange
	userID := uuid.New()
	updates := &UserUpdateRequest{
		Username: utils.StringPtr("newusername"),
	}

	suite.mockRepo.On("GetUserByID", userID).Return(nil, fmt.Errorf("用户不存在"))

	// Act
	user, err := suite.userService.UpdateUser(userID, updates)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "用户不存在")
}

func (suite *UserServiceTestSuite) TestUpdateLastLogin_Success() {
	// Arrange
	userID := uuid.New()
	suite.mockRepo.On("UpdateUserLastLogin", userID).Return(nil)

	// Act
	err := suite.userService.UpdateLastLogin(userID)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestUpdateLastLogin_Error() {
	// Arrange
	userID := uuid.New()
	suite.mockRepo.On("UpdateUserLastLogin", userID).Return(fmt.Errorf("数据库错误"))

	// Act
	err := suite.userService.UpdateLastLogin(userID)

	// Assert
	assert.Error(suite.T(), err)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
} 