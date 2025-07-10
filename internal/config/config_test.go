package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	originalEnvs map[string]string
}

func (suite *ConfigTestSuite) SetupTest() {
	// 保存原始环境变量
	suite.originalEnvs = make(map[string]string)
	envKeys := []string{
		"PORT", "GO_ENV", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"DB_MAX_CONNECTIONS", "DB_MAX_IDLE", "DB_CONN_MAX_LIFETIME",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"JWT_SECRET", "JWT_EXPIRES_IN",
		"AI_PROVIDER", "OPENAI_API_KEY", "CLAUDE_API_KEY",
		"AI_TIMEOUT", "AI_MAX_RETRIES", "AI_DEFAULT_MODEL", "AI_MAX_TOKENS",
	}
	
	for _, key := range envKeys {
		suite.originalEnvs[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
}

func (suite *ConfigTestSuite) TearDownTest() {
	// 恢复原始环境变量
	for key, value := range suite.originalEnvs {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

func (suite *ConfigTestSuite) TestLoad_DefaultValues() {
	// Arrange & Act
	cfg := Load()

	// Assert
	assert.Equal(suite.T(), "8080", cfg.Port)
	assert.Equal(suite.T(), "development", cfg.Env)
	
	// Database defaults
	assert.Equal(suite.T(), "localhost", cfg.Database.Host)
	assert.Equal(suite.T(), "3306", cfg.Database.Port)
	assert.Equal(suite.T(), "root", cfg.Database.User)
	assert.Equal(suite.T(), "lzx234258", cfg.Database.Password)
	assert.Equal(suite.T(), "aicode", cfg.Database.Name)
	assert.Equal(suite.T(), 100, cfg.Database.MaxConnections)
	assert.Equal(suite.T(), 10, cfg.Database.MaxIdleConn)
	assert.Equal(suite.T(), 3600, cfg.Database.ConnMaxLifetime)
	
	// Redis defaults
	assert.Equal(suite.T(), "localhost", cfg.Redis.Host)
	assert.Equal(suite.T(), "6379", cfg.Redis.Port)
	assert.Equal(suite.T(), "", cfg.Redis.Password)
	assert.Equal(suite.T(), 0, cfg.Redis.DB)
	assert.Equal(suite.T(), 3, cfg.Redis.MaxRetries)
	assert.Equal(suite.T(), 10, cfg.Redis.PoolSize)
	assert.Equal(suite.T(), 5, cfg.Redis.MinIdleConns)
	
	// JWT defaults
	assert.Equal(suite.T(), "ai-dev-platform-secret", cfg.JWT.Secret)
	assert.Equal(suite.T(), 86400, cfg.JWT.ExpiresIn)
	
	// AI defaults
	assert.Equal(suite.T(), "openai", cfg.AI.Provider)
	assert.Equal(suite.T(), "", cfg.AI.OpenAIKey)
	assert.Equal(suite.T(), "", cfg.AI.ClaudeKey)
	assert.Equal(suite.T(), 30, cfg.AI.Timeout)
	assert.Equal(suite.T(), 3, cfg.AI.MaxRetries)
	assert.Equal(suite.T(), "gpt-3.5-turbo", cfg.AI.DefaultModel)
	assert.Equal(suite.T(), 2048, cfg.AI.MaxTokens)
	
	// CORS defaults
	assert.Equal(suite.T(), []string{"http://localhost:3000", "http://localhost:8080"}, cfg.CORS.Origins)
	assert.Equal(suite.T(), []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, cfg.CORS.Methods)
	assert.True(suite.T(), cfg.CORS.Credentials)
}

func (suite *ConfigTestSuite) TestLoad_CustomEnvironmentValues() {
	// Arrange
	os.Setenv("PORT", "9090")
	os.Setenv("GO_ENV", "production")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_MAX_CONNECTIONS", "200")
	os.Setenv("REDIS_HOST", "redis.example.com")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("JWT_SECRET", "custom-secret")
	os.Setenv("JWT_EXPIRES_IN", "3600")
	os.Setenv("AI_PROVIDER", "claude")
	os.Setenv("OPENAI_API_KEY", "sk-test-openai")
	os.Setenv("CLAUDE_API_KEY", "sk-test-claude")
	os.Setenv("AI_TIMEOUT", "60")
	os.Setenv("AI_MAX_RETRIES", "5")
	os.Setenv("AI_DEFAULT_MODEL", "claude-3-opus")
	os.Setenv("AI_MAX_TOKENS", "4096")

	// Act
	cfg := Load()

	// Assert
	assert.Equal(suite.T(), "9090", cfg.Port)
	assert.Equal(suite.T(), "production", cfg.Env)
	assert.Equal(suite.T(), "db.example.com", cfg.Database.Host)
	assert.Equal(suite.T(), "5432", cfg.Database.Port)
	assert.Equal(suite.T(), "testuser", cfg.Database.User)
	assert.Equal(suite.T(), "testpass", cfg.Database.Password)
	assert.Equal(suite.T(), "testdb", cfg.Database.Name)
	assert.Equal(suite.T(), 200, cfg.Database.MaxConnections)
	assert.Equal(suite.T(), "redis.example.com", cfg.Redis.Host)
	assert.Equal(suite.T(), "6380", cfg.Redis.Port)
	assert.Equal(suite.T(), 1, cfg.Redis.DB)
	assert.Equal(suite.T(), "custom-secret", cfg.JWT.Secret)
	assert.Equal(suite.T(), 3600, cfg.JWT.ExpiresIn)
	assert.Equal(suite.T(), "claude", cfg.AI.Provider)
	assert.Equal(suite.T(), "sk-test-openai", cfg.AI.OpenAIKey)
	assert.Equal(suite.T(), "sk-test-claude", cfg.AI.ClaudeKey)
	assert.Equal(suite.T(), 60, cfg.AI.Timeout)
	assert.Equal(suite.T(), 5, cfg.AI.MaxRetries)
	assert.Equal(suite.T(), "claude-3-opus", cfg.AI.DefaultModel)
	assert.Equal(suite.T(), 4096, cfg.AI.MaxTokens)
}

func (suite *ConfigTestSuite) TestLoad_InvalidIntegerEnvironmentValues() {
	// Arrange
	os.Setenv("DB_MAX_CONNECTIONS", "invalid")
	os.Setenv("REDIS_DB", "not-a-number")
	os.Setenv("JWT_EXPIRES_IN", "abc")

	// Act
	cfg := Load()

	// Assert - 应该使用默认值
	assert.Equal(suite.T(), 100, cfg.Database.MaxConnections)
	assert.Equal(suite.T(), 0, cfg.Redis.DB)
	assert.Equal(suite.T(), 86400, cfg.JWT.ExpiresIn)
}

func (suite *ConfigTestSuite) TestGetEnv() {
	// Test with existing environment variable
	os.Setenv("TEST_KEY", "test_value")
	result := getEnv("TEST_KEY", "default")
	assert.Equal(suite.T(), "test_value", result)

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTING_KEY", "default")
	assert.Equal(suite.T(), "default", result)

	// Clean up
	os.Unsetenv("TEST_KEY")
}

func (suite *ConfigTestSuite) TestGetEnvInt() {
	// Test with valid integer
	os.Setenv("TEST_INT", "123")
	result := getEnvInt("TEST_INT", 456)
	assert.Equal(suite.T(), 123, result)

	// Test with invalid integer
	os.Setenv("TEST_INT", "invalid")
	result = getEnvInt("TEST_INT", 456)
	assert.Equal(suite.T(), 456, result)

	// Test with non-existing environment variable
	result = getEnvInt("NON_EXISTING_INT", 789)
	assert.Equal(suite.T(), 789, result)

	// Clean up
	os.Unsetenv("TEST_INT")
}

func (suite *ConfigTestSuite) TestIsDevelopment() {
	// Test development environment
	cfg := &Config{Env: "development"}
	assert.True(suite.T(), cfg.IsDevelopment())

	// Test production environment
	cfg = &Config{Env: "production"}
	assert.False(suite.T(), cfg.IsDevelopment())

	// Test other environment
	cfg = &Config{Env: "staging"}
	assert.False(suite.T(), cfg.IsDevelopment())
}

func (suite *ConfigTestSuite) TestIsProduction() {
	// Test production environment
	cfg := &Config{Env: "production"}
	assert.True(suite.T(), cfg.IsProduction())

	// Test development environment
	cfg = &Config{Env: "development"}
	assert.False(suite.T(), cfg.IsProduction())

	// Test other environment
	cfg = &Config{Env: "staging"}
	assert.False(suite.T(), cfg.IsProduction())
}

func (suite *ConfigTestSuite) TestGetDSN() {
	// Arrange
	cfg := &Config{
		Database: DatabaseConfig{
			User:     "testuser",
			Password: "testpass",
			Host:     "localhost",
			Port:     "3306",
			Name:     "testdb",
		},
	}

	// Act
	dsn := cfg.GetDSN()

	// Assert
	expected := "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	assert.Equal(suite.T(), expected, dsn)
}

func (suite *ConfigTestSuite) TestGetRedisAddr() {
	// Arrange
	cfg := &Config{
		Redis: RedisConfig{
			Host: "redis.example.com",
			Port: "6379",
		},
	}

	// Act
	addr := cfg.GetRedisAddr()

	// Assert
	expected := "redis.example.com:6379"
	assert.Equal(suite.T(), expected, addr)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
} 