package tests

import (
	"fmt"
	"os"
	"time"

	"ai-dev-platform/internal/config"
)

// TestConfig 测试配置
type TestConfig struct {
	Database DatabaseTestConfig `yaml:"database"`
	Redis    RedisTestConfig    `yaml:"redis"`
	JWT      JWTTestConfig      `yaml:"jwt"`
	AI       AITestConfig       `yaml:"ai"`
}

// DatabaseTestConfig 数据库测试配置
type DatabaseTestConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// RedisTestConfig Redis测试配置
type RedisTestConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTTestConfig JWT测试配置
type JWTTestConfig struct {
	Secret    string        `yaml:"secret"`
	ExpiresIn time.Duration `yaml:"expires_in"`
}

// AITestConfig AI测试配置
type AITestConfig struct {
	Provider      string `yaml:"provider"`
	MockResponses bool   `yaml:"mock_responses"`
	Timeout       int    `yaml:"timeout"`
}

// GetTestConfig 获取测试配置
func GetTestConfig() *config.Config {
	return &config.Config{
		Port: "8080",
		Env:  "test",
		Database: config.DatabaseConfig{
			Host:            getTestEnvOrDefault("TEST_DB_HOST", "localhost"),
			Port:            "3306",
			User:            getTestEnvOrDefault("TEST_DB_USER", "test"),
			Password:        getTestEnvOrDefault("TEST_DB_PASSWORD", "test"),
			Name:            getTestEnvOrDefault("TEST_DB_NAME", "ai_dev_platform_test"),
			MaxConnections:  5,
			MaxIdleConn:     2,
			ConnMaxLifetime: 300,
		},
		Redis: config.RedisConfig{
			Host:         getTestEnvOrDefault("TEST_REDIS_HOST", "localhost"),
			Port:         "6379",
			Password:     getTestEnvOrDefault("TEST_REDIS_PASSWORD", ""),
			DB:           1, // 使用不同的DB避免与开发环境冲突
			MaxRetries:   3,
			PoolSize:     5,
			MinIdleConns: 2,
		},
		JWT: config.JWTConfig{
			Secret:    getTestEnvOrDefault("TEST_JWT_SECRET", "test-secret-key-for-testing-only"),
			ExpiresIn: 3600,
		},
		AI: config.AIConfig{
			DefaultProvider: getTestEnvOrDefault("TEST_AI_PROVIDER", "openai"),
			EnableCache:     true,
			CacheTTL:        300,
			OpenAIConfig: &config.OpenAIConfig{
				APIKey:       getTestEnvOrDefault("TEST_OPENAI_KEY", "test-key"),
				DefaultModel: "gpt-3.5-turbo",
			},
		},
		CORS: config.CORSConfig{
			Origins:     []string{"*"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"*"},
			Credentials: true,
		},
	}
}

// getTestEnvOrDefault 获取测试环境变量或返回默认值
func getTestEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestConstants 测试常量
var TestConstants = struct {
	// 用户测试数据
	TestUser struct {
		Username string
		Email    string
		Password string
		FullName string
	}
	
	// 项目测试数据
	TestProject struct {
		Name        string
		Description string
		Type        string
	}
	
	// JWT测试数据
	TestJWT struct {
		ValidToken   string
		InvalidToken string
		ExpiredToken string
	}
	
	// 时间相关
	TestTimeout time.Duration
	
	// 分页测试数据
	TestPagination struct {
		DefaultPage     int
		DefaultPageSize int
		MaxPageSize     int
	}
}{
	TestUser: struct {
		Username string
		Email    string
		Password string
		FullName string
	}{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword123",
		FullName: "Test User",
	},
	
	TestProject: struct {
		Name        string
		Description string
		Type        string
	}{
		Name:        "测试项目",
		Description: "这是一个用于测试的项目",
		Type:        "web",
	},
	
	TestJWT: struct {
		ValidToken   string
		InvalidToken string
		ExpiredToken string
	}{
		ValidToken:   "valid.jwt.token",
		InvalidToken: "invalid.jwt.token",
		ExpiredToken: "expired.jwt.token",
	},
	
	TestTimeout: time.Second * 5,
	
	TestPagination: struct {
		DefaultPage     int
		DefaultPageSize int
		MaxPageSize     int
	}{
		DefaultPage:     1,
		DefaultPageSize: 10,
		MaxPageSize:     100,
	},
}

// DatabaseURL 生成测试数据库连接URL
func DatabaseURL(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
}

// RedisURL 生成测试Redis连接URL
func RedisURL(cfg *config.Config) string {
	if cfg.Redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d",
			cfg.Redis.Password,
			cfg.Redis.Host,
			cfg.Redis.Port,
			cfg.Redis.DB,
		)
	}
	return fmt.Sprintf("redis://%s:%s/%d",
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.DB,
	)
}

// IsTestEnvironment 检查是否为测试环境
func IsTestEnvironment() bool {
	return os.Getenv("GO_ENV") == "test" || 
		   os.Getenv("ENVIRONMENT") == "test" ||
		   os.Getenv("APP_ENV") == "test"
}

// SetTestEnvironment 设置测试环境变量
func SetTestEnvironment() {
	os.Setenv("GO_ENV", "test")
	os.Setenv("ENVIRONMENT", "test")
}

// CleanTestEnvironment 清理测试环境变量
func CleanTestEnvironment() {
	testEnvKeys := []string{
		"GO_ENV", "ENVIRONMENT", "APP_ENV",
		"TEST_DB_HOST", "TEST_DB_USER", "TEST_DB_PASSWORD", "TEST_DB_NAME",
		"TEST_REDIS_HOST", "TEST_REDIS_PASSWORD",
		"TEST_JWT_SECRET", "TEST_AI_PROVIDER", "TEST_OPENAI_KEY", "TEST_CLAUDE_KEY",
	}
	
	for _, key := range testEnvKeys {
		os.Unsetenv(key)
	}
} 