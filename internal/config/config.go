package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config 应用配置结构
type Config struct {
	// 服务器配置
	Port string
	Env  string

	// 数据库配置
	Database DatabaseConfig

	// Redis配置
	Redis RedisConfig

	// JWT配置
	JWT JWTConfig

	// AI服务配置
	AI AIConfig

	// PUML服务配置
	PUML PUMLConfig

	// CORS配置
	CORS CORSConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxConnections  int
	MaxIdleConn     int
	ConnMaxLifetime int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string
	ExpiresIn int // 秒
}

// AIConfig AI服务配置
type AIConfig struct {
	DefaultProvider string        `json:"default_provider" mapstructure:"default_provider"`
	EnableCache     bool          `json:"enable_cache" mapstructure:"enable_cache"`
	CacheTTL        time.Duration `json:"cache_ttl" mapstructure:"cache_ttl"`
	OpenAIConfig    *OpenAIConfig `json:"openai_config" mapstructure:"openai_config"`
	ClaudeConfig    *ClaudeConfig `json:"claude_config" mapstructure:"claude_config"`
	GeminiConfig    *GeminiConfig `json:"gemini_config" mapstructure:"gemini_config"`
}

// PUMLConfig puml服务相关配置
type PUMLConfig struct {
	ServerURL string `json:"server_url" mapstructure:"server_url"`
}

// OpenAIConfig OpenAI相关配置
type OpenAIConfig struct {
	APIKey       string `json:"api_key" mapstructure:"api_key"`
	DefaultModel string `json:"default_model" mapstructure:"default_model"`
}

// ClaudeConfig Claude相关配置
type ClaudeConfig struct {
	APIKey       string `json:"api_key" mapstructure:"api_key"`
	DefaultModel string `json:"default_model" mapstructure:"default_model"`
}

// GeminiConfig Gemini相关配置
type GeminiConfig struct {
	APIKey       string `json:"api_key" mapstructure:"api_key"`
	DefaultModel string `json:"default_model" mapstructure:"default_model"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
}

var cfg *Config

// Load 加载配置
func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("GO_ENV", "development"),

		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", "lzx234258"),
			Name:            getEnv("DB_NAME", "aicode"),
			MaxConnections:  getEnvInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConn:     getEnvInt("DB_MAX_IDLE", 10),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 3600),
		},

		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvInt("REDIS_DB", 0),
			MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 5),
		},

		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "ai-dev-platform-secret"),
			ExpiresIn: getEnvInt("JWT_EXPIRES_IN", 86400), // 24小时
		},

		AI: AIConfig{
			DefaultProvider: "openai",
			EnableCache:     true,
			CacheTTL:        60 * time.Minute,
			OpenAIConfig: &OpenAIConfig{
				APIKey:       os.Getenv("OPENAI_API_KEY"),
				DefaultModel: "gpt-4",
			},
			ClaudeConfig: &ClaudeConfig{
				APIKey:       os.Getenv("CLAUDE_API_KEY"),
				DefaultModel: "claude-2",
			},
			GeminiConfig: &GeminiConfig{
				APIKey:       os.Getenv("GEMINI_API_KEY"),
				DefaultModel: "gemini-pro",
			},
		},
		PUML: PUMLConfig{
			ServerURL: "http://localhost:8888",
		},

		CORS: CORSConfig{
			Origins:     []string{"http://localhost:3000", "http://localhost:8080"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"Content-Type", "Authorization", "X-Requested-With"},
			Credentials: true,
		},
	}

	// 验证必要配置
	validateConfig(cfg)

	return cfg
}

func validateConfig(cfg *Config) {
	if cfg.Env == "production" {
		if cfg.JWT.Secret == "ai-dev-platform-secret" {
			log.Println("警告: JWT密钥未在生产环境中更改")
		}
		if cfg.AI.OpenAIConfig.APIKey == "" && cfg.AI.ClaudeConfig.APIKey == "" && cfg.AI.GeminiConfig.APIKey == "" {
			log.Println("警告: 在生产环境中未设置任何AI服务密钥")
		}
	}

	if cfg.AI.OpenAIConfig.APIKey == "" && cfg.AI.ClaudeConfig.APIKey == "" && cfg.AI.GeminiConfig.APIKey == "" {
		log.Println("警告: 未设置AI服务密钥，部分功能可能无法使用")
	}
}

// getEnv 获取环境变量，如果不存在返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return c.Database.User + ":" + c.Database.Password +
		"@tcp(" + c.Database.Host + ":" + c.Database.Port + ")/" +
		c.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + c.Redis.Port
}
