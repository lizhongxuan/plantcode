package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
)

// TestSuite 测试套件基类
type TestSuite struct {
	suite.Suite
	DB          *gorm.DB
	Config      *config.Config
	TestUserID  uuid.UUID
	TestProject *model.Project
	Cleanup     []func()
}

// SetupSuite 测试套件初始化
func (s *TestSuite) SetupSuite() {
	// 设置测试配置
	s.Config = &config.Config{
		Env: "test",
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "test",
			Name:     "test_db",
		},
		JWT: config.JWTConfig{
			Secret:    "test-secret-key-for-testing",
			ExpiresIn: 86400, // 24小时，避免立即过期
		},
		AI: config.AIConfig{
			DefaultProvider: "openai",
			EnableCache:     true,
			CacheTTL:        300,
			OpenAIConfig: &config.OpenAIConfig{
				APIKey:       "test-key",
				DefaultModel: "gpt-3.5-turbo",
			},
		},
		CORS: config.CORSConfig{
			Origins:     []string{"http://localhost:3000"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"Content-Type", "Authorization"},
			Credentials: true,
		},
	}

	// 使用SQLite内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	s.Require().NoError(err)
	s.DB = db

	// 创建测试表
	err = s.DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Requirement{},
		&model.PUMLDiagram{},
		&model.Document{},
		&model.AsyncTask{},
		&model.UserAIConfig{},
	)
	s.Require().NoError(err)

	// 验证表是否创建成功
	var tableCount int64
	result := s.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tableCount)
	s.Require().NoError(result.Error)
	s.Require().True(tableCount > 0, "数据库表创建失败")

	// 创建测试用户
	s.createTestUser()
	s.createTestProject()
}

// TearDownSuite 测试套件清理
func (s *TestSuite) TearDownSuite() {
	// 执行清理函数
	for i := len(s.Cleanup) - 1; i >= 0; i-- {
		s.Cleanup[i]()
	}

	// 清理测试数据
	s.cleanupTestData()
}

// SetupTest 每个测试前的准备
func (s *TestSuite) SetupTest() {
	// 开始事务
	s.DB = s.DB.Begin()
}

// TearDownTest 每个测试后的清理
func (s *TestSuite) TearDownTest() {
	// 回滚事务
	s.DB.Rollback()
}

// createTestUser 创建测试用户
func (s *TestSuite) createTestUser() {
	s.TestUserID = uuid.New()
	testUser := &model.User{
		UserID:       s.TestUserID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		Status:       model.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.DB.Create(testUser).Error
	s.Require().NoError(err)
}

// createTestProject 创建测试项目
func (s *TestSuite) createTestProject() {
	s.TestProject = &model.Project{
		ProjectID:   uuid.New(),
		UserID:      s.TestUserID,
		ProjectName: "Test Project",
		Description: "Test project description",
		ProjectType: "web_application",
		Status:      model.ProjectStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.DB.Create(s.TestProject).Error
	s.Require().NoError(err)
}

// cleanupTestData 清理测试数据
func (s *TestSuite) cleanupTestData() {
	// SQLite内存数据库会自动清理，这里留空即可
}

// AddCleanup 添加清理函数
func (s *TestSuite) AddCleanup(fn func()) {
	s.Cleanup = append(s.Cleanup, fn)
}

// MockAIService 模拟AI服务
type MockAIService struct {
	service.AIService
	responses map[string]interface{}
}

// NewMockAIService 创建模拟AI服务
func NewMockAIService() *MockAIService {
	return &MockAIService{
		responses: make(map[string]interface{}),
	}
}

// SetResponse 设置模拟响应
func (m *MockAIService) SetResponse(method string, response interface{}) {
	m.responses[method] = response
}

// AssertExpected 断言辅助函数
func AssertExpected(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertNoError 断言无错误
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	assert.NoError(t, err, msgAndArgs...)
}

// AssertError 断言有错误
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	assert.Error(t, err, msgAndArgs...)
} 