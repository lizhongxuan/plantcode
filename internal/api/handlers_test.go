package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(req *model.CreateUserRequest) (*model.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) LoginUser(req *model.LoginRequest) (*service.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResponse), args.Error(1)
}

func (m *MockUserService) ValidateToken(token string) (*model.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUser(userID uuid.UUID) (*model.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(userID uuid.UUID, updates *service.UserUpdateRequest) (*model.User, error) {
	args := m.Called(userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateLastLogin(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// MockProjectService 模拟项目服务
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(userID uuid.UUID, req *model.CreateProjectRequest) (*model.Project, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectService) GetUserProjects(userID uuid.UUID, page, pageSize int) ([]*model.Project, utils.PaginationInfo, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, utils.PaginationInfo{}, args.Error(2)
	}
	return args.Get(0).([]*model.Project), args.Get(1).(utils.PaginationInfo), args.Error(2)
}

func (m *MockProjectService) GetProject(projectID uuid.UUID, userID uuid.UUID) (*model.Project, error) {
	args := m.Called(projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(projectID uuid.UUID, userID uuid.UUID, updates *service.ProjectUpdateRequest) (*model.Project, error) {
	args := m.Called(projectID, userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(projectID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(projectID, userID)
	return args.Error(0)
}

type HandlersTestSuite struct {
	suite.Suite
	mockUserService    *MockUserService
	mockProjectService *MockProjectService
	handlers           *Handlers
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.mockUserService = new(MockUserService)
	suite.mockProjectService = new(MockProjectService)
	suite.handlers = NewHandlers(suite.mockUserService, suite.mockProjectService)
}

func (suite *HandlersTestSuite) TearDownTest() {
	suite.mockUserService.AssertExpectations(suite.T())
	suite.mockProjectService.AssertExpectations(suite.T())
}

func (suite *HandlersTestSuite) TestHealth() {
	// Arrange
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.Health(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "服务运行正常", response.Message)
	assert.NotNil(suite.T(), response.Data)

	// 检查返回数据的结构
	data, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "healthy", data["status"])
	assert.Equal(suite.T(), "ai-dev-platform", data["service"])
	assert.Equal(suite.T(), "1.0.0", data["version"])
	assert.NotNil(suite.T(), data["timestamp"])
}

// 添加认证功能的单元测试

func (suite *HandlersTestSuite) TestRegisterUser_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Status:   model.UserStatusActive,
	}

	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	suite.mockUserService.On("RegisterUser", mock.MatchedBy(func(r *model.CreateUserRequest) bool {
		return r.Username == req.Username && r.Email == req.Email
	})).Return(user, nil)

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.RegisterUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "用户注册成功", response.Message)
}

func (suite *HandlersTestSuite) TestRegisterUser_InvalidRequest() {
	// Arrange
	request := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader([]byte("invalid json")))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.RegisterUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "无效的请求数据", response.Error)
}

func (suite *HandlersTestSuite) TestRegisterUser_DuplicateEmail() {
	// Arrange
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	suite.mockUserService.On("RegisterUser", mock.MatchedBy(func(r *model.CreateUserRequest) bool {
		return r.Email == req.Email
	})).Return(nil, fmt.Errorf("邮箱已被注册"))

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.RegisterUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "邮箱已被注册", response.Error)
}

func (suite *HandlersTestSuite) TestLoginUser_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Status:   model.UserStatusActive,
	}

	loginResponse := &service.LoginResponse{
		User:  user,
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
	}

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockUserService.On("LoginUser", mock.MatchedBy(func(r *model.LoginRequest) bool {
		return r.Email == req.Email && r.Password == req.Password
	})).Return(loginResponse, nil)

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.LoginUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "登录成功", response.Message)
	
	// 验证返回的数据包含token
	responseData := response.Data.(map[string]interface{})
	assert.Contains(suite.T(), responseData, "token")
	assert.Contains(suite.T(), responseData, "user")
}

func (suite *HandlersTestSuite) TestLoginUser_InvalidCredentials() {
	// Arrange
	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	suite.mockUserService.On("LoginUser", mock.MatchedBy(func(r *model.LoginRequest) bool {
		return r.Email == req.Email
	})).Return(nil, fmt.Errorf("用户不存在或密码错误"))

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.LoginUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "用户不存在或密码错误", response.Error)
}

func (suite *HandlersTestSuite) TestLoginUser_InvalidRequest() {
	// Arrange
	request := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader([]byte("invalid json")))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.LoginUser(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "无效的请求数据", response.Error)
}

func (suite *HandlersTestSuite) TestValidateToken_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Status:   model.UserStatusActive,
	}

	// 创建包含用户的上下文
	ctx := context.WithValue(context.Background(), UserContextKey, user)
	request := httptest.NewRequest("GET", "/api/auth/validate", nil)
	request = request.WithContext(ctx)
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.ValidateToken(recorder, request)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	
	var response utils.APIResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Token验证成功", response.Message)
	
	// 验证返回的用户数据
	responseData := response.Data.(map[string]interface{})
	assert.Equal(suite.T(), user.UserID.String(), responseData["user_id"])
	assert.Equal(suite.T(), user.Email, responseData["email"])
}

func (suite *HandlersTestSuite) TestGetCurrentUser() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Status:   model.UserStatusActive,
	}

	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("GET", "/api/users/me", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.GetCurrentUser(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "获取用户信息成功", response.Message)
}

func (suite *HandlersTestSuite) TestUpdateCurrentUser_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Status:   model.UserStatusActive,
	}

	updateReq := &service.UserUpdateRequest{
		FullName: utils.StringPtr("Updated Name"),
	}

	updatedUser := &model.User{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		FullName: "Updated Name",
		Status:   user.Status,
	}

	suite.mockUserService.On("UpdateUser", user.UserID, updateReq).Return(updatedUser, nil)

	body, _ := json.Marshal(updateReq)
	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("PUT", "/api/users/me", bytes.NewBuffer(body)).WithContext(ctx)
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	suite.handlers.UpdateCurrentUser(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "用户信息更新成功", response.Message)
}

func (suite *HandlersTestSuite) TestCreateProject_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	projectReq := &model.CreateProjectRequest{
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web",
	}

	expectedProject := &model.Project{
		ProjectID:   uuid.New(),
		UserID:      user.UserID,
		ProjectName: projectReq.ProjectName,
		Description: projectReq.Description,
		ProjectType: projectReq.ProjectType,
		Status:      "planning",
	}

	suite.mockProjectService.On("CreateProject", user.UserID, projectReq).Return(expectedProject, nil)

	body, _ := json.Marshal(projectReq)
	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(body)).WithContext(ctx)
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	suite.handlers.CreateProject(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "项目创建成功", response.Message)
}

func (suite *HandlersTestSuite) TestGetUserProjects_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	projects := []*model.Project{
		{
			ProjectID:   uuid.New(),
			UserID:      user.UserID,
			ProjectName: "项目1",
			Description: "项目1描述",
			ProjectType: "web",
			Status:      "planning",
		},
	}

	pagination := utils.PaginationInfo{
		Page:      1,
		PageSize:  10,
		Total:     1,
		TotalPage: 1,
	}

	suite.mockProjectService.On("GetUserProjects", user.UserID, 1, 10).Return(projects, pagination, nil)

	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("GET", "/api/projects", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.GetUserProjects(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response utils.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "获取项目列表成功", response.Message)
	assert.Equal(suite.T(), pagination, response.Pagination)
}

func (suite *HandlersTestSuite) TestGetProject_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	projectID := uuid.New()
	project := &model.Project{
		ProjectID:   projectID,
		UserID:      user.UserID,
		ProjectName: "测试项目",
		Description: "这是一个测试项目",
		ProjectType: "web",
		Status:      "planning",
	}

	suite.mockProjectService.On("GetProject", projectID, user.UserID).Return(project, nil)

	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("GET", fmt.Sprintf("/api/projects/%s", projectID), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.GetProject(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "获取项目详情成功", response.Message)
}

func (suite *HandlersTestSuite) TestGetProject_InvalidID() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("GET", "/api/projects/invalid-uuid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.GetProject(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "无效的项目ID格式", response.Error)
}

func (suite *HandlersTestSuite) TestGetProject_NotFound() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	projectID := uuid.New()

	suite.mockProjectService.On("GetProject", projectID, user.UserID).Return(nil, fmt.Errorf("项目不存在"))

	ctx := context.WithValue(context.Background(), UserContextKey, user)
	httpReq := httptest.NewRequest("GET", fmt.Sprintf("/api/projects/%s", projectID), nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	suite.handlers.GetProject(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), "项目不存在", response.Error)
}

func (suite *HandlersTestSuite) TestExtractIDFromPath() {
	testCases := []struct {
		name     string
		path     string
		prefix   string
		expected string
	}{
		{
			name:     "valid project ID",
			path:     "/api/projects/550e8400-e29b-41d4-a716-446655440000",
			prefix:   "/api/projects/",
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:     "path with additional segments",
			path:     "/api/projects/550e8400-e29b-41d4-a716-446655440000/documents",
			prefix:   "/api/projects/",
			expected: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:     "path without prefix",
			path:     "/other/path/550e8400-e29b-41d4-a716-446655440000",
			prefix:   "/api/projects/",
			expected: "",
		},
		{
			name:     "path equals prefix",
			path:     "/api/projects/",
			prefix:   "/api/projects/",
			expected: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := extractIDFromPath(tc.path, tc.prefix)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
} 