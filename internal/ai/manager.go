package ai

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AIManager AI服务管理器
type AIManager struct {
	clients     map[AIProvider]AIClient
	defaultProvider AIProvider
	cache       AICache
	mutex       sync.RWMutex
}

// AIManagerConfig AI管理器配置
type AIManagerConfig struct {
	DefaultProvider AIProvider
	OpenAIConfig    *OpenAIConfig
	ClaudeConfig    *ClaudeConfig
	EnableCache     bool
	CacheTTL        time.Duration
}

// ClaudeConfig Claude配置（预留）
type ClaudeConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// AICache AI响应缓存接口
type AICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data   map[string]*cacheItem
	mutex  sync.RWMutex
}

type cacheItem struct {
	value     interface{}
	expireAt  time.Time
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		data: make(map[string]*cacheItem),
	}
	
	// 启动清理协程
	go cache.cleanup()
	
	return cache
}

// Get 获取缓存值
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return nil, false
	}
	
	if time.Now().After(item.expireAt) {
		delete(c.data, key)
		return nil, false
	}
	
	return item.value, true
}

// Set 设置缓存值
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.data[key] = &cacheItem{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

// Delete 删除缓存值
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.data, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.data = make(map[string]*cacheItem)
}

// cleanup 清理过期缓存
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expireAt) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// NewAIManager 创建AI管理器
func NewAIManager(config AIManagerConfig) (*AIManager, error) {
	manager := &AIManager{
		clients:         make(map[AIProvider]AIClient),
		defaultProvider: config.DefaultProvider,
	}
	
	// 初始化缓存
	if config.EnableCache {
		manager.cache = NewMemoryCache()
	}
	
	// 初始化OpenAI客户端
	if config.OpenAIConfig != nil {
		openAIClient := NewOpenAIClient(*config.OpenAIConfig)
		manager.clients[ProviderOpenAI] = openAIClient
	}
	
	// 初始化Claude客户端（预留）
	if config.ClaudeConfig != nil {
		// TODO: 实现Claude客户端
		// claudeClient := NewClaudeClient(*config.ClaudeConfig)
		// manager.clients[ProviderClaude] = claudeClient
	}
	
	// 验证默认提供商是否可用
	if _, exists := manager.clients[config.DefaultProvider]; !exists {
		return nil, fmt.Errorf("默认AI提供商 %s 未配置", config.DefaultProvider)
	}
	
	return manager, nil
}

// GetClient 获取指定提供商的客户端
func (m *AIManager) GetClient(provider AIProvider) (AIClient, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	client, exists := m.clients[provider]
	if !exists {
		return nil, fmt.Errorf("AI提供商 %s 未配置", provider)
	}
	
	return client, nil
}

// GetDefaultClient 获取默认客户端
func (m *AIManager) GetDefaultClient() AIClient {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	return m.clients[m.defaultProvider]
}

// AnalyzeRequirement 分析业务需求（带缓存）
func (m *AIManager) AnalyzeRequirement(ctx context.Context, requirement string, provider ...AIProvider) (*RequirementAnalysis, error) {
	// 确定使用的提供商
	targetProvider := m.defaultProvider
	if len(provider) > 0 {
		targetProvider = provider[0]
	}
	
	// 检查缓存
	if m.cache != nil {
		cacheKey := m.generateCacheKey("analyze", targetProvider, requirement)
		if cached, exists := m.cache.Get(cacheKey); exists {
			if analysis, ok := cached.(*RequirementAnalysis); ok {
				return analysis, nil
			}
		}
	}
	
	// 获取客户端
	client, err := m.GetClient(targetProvider)
	if err != nil {
		return nil, err
	}
	
	// 调用AI分析
	analysis, err := client.AnalyzeRequirement(ctx, requirement)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if m.cache != nil {
		cacheKey := m.generateCacheKey("analyze", targetProvider, requirement)
		m.cache.Set(cacheKey, analysis, 30*time.Minute)
	}
	
	return analysis, nil
}

// GenerateQuestions 生成补充问题（带缓存）
func (m *AIManager) GenerateQuestions(ctx context.Context, analysis *RequirementAnalysis, provider ...AIProvider) ([]Question, error) {
	// 确定使用的提供商
	targetProvider := m.defaultProvider
	if len(provider) > 0 {
		targetProvider = provider[0]
	}
	
	// 检查缓存
	if m.cache != nil {
		cacheKey := m.generateCacheKey("questions", targetProvider, analysis.ID)
		if cached, exists := m.cache.Get(cacheKey); exists {
			if questions, ok := cached.([]Question); ok {
				return questions, nil
			}
		}
	}
	
	// 获取客户端
	client, err := m.GetClient(targetProvider)
	if err != nil {
		return nil, err
	}
	
	// 调用AI生成问题
	questions, err := client.GenerateQuestions(ctx, analysis)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if m.cache != nil {
		cacheKey := m.generateCacheKey("questions", targetProvider, analysis.ID)
		m.cache.Set(cacheKey, questions, 15*time.Minute)
	}
	
	return questions, nil
}

// GeneratePUML 生成PUML图表（带缓存）
func (m *AIManager) GeneratePUML(ctx context.Context, analysis *RequirementAnalysis, diagramType PUMLType, provider ...AIProvider) (*PUMLDiagram, error) {
	// 确定使用的提供商
	targetProvider := m.defaultProvider
	if len(provider) > 0 {
		targetProvider = provider[0]
	}
	
	// 检查缓存
	if m.cache != nil {
		cacheKey := m.generateCacheKey("puml", targetProvider, analysis.ID, string(diagramType))
		if cached, exists := m.cache.Get(cacheKey); exists {
			if diagram, ok := cached.(*PUMLDiagram); ok {
				return diagram, nil
			}
		}
	}
	
	// 获取客户端
	client, err := m.GetClient(targetProvider)
	if err != nil {
		return nil, err
	}
	
	// 调用AI生成PUML
	diagram, err := client.GeneratePUML(ctx, analysis, diagramType)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if m.cache != nil {
		cacheKey := m.generateCacheKey("puml", targetProvider, analysis.ID, string(diagramType))
		m.cache.Set(cacheKey, diagram, 60*time.Minute)
	}
	
	return diagram, nil
}

// GenerateDocument 生成开发文档（带缓存）
func (m *AIManager) GenerateDocument(ctx context.Context, analysis *RequirementAnalysis, provider ...AIProvider) (*DevelopmentDocument, error) {
	// 确定使用的提供商
	targetProvider := m.defaultProvider
	if len(provider) > 0 {
		targetProvider = provider[0]
	}
	
	// 检查缓存
	if m.cache != nil {
		cacheKey := m.generateCacheKey("document", targetProvider, analysis.ID)
		if cached, exists := m.cache.Get(cacheKey); exists {
			if document, ok := cached.(*DevelopmentDocument); ok {
				return document, nil
			}
		}
	}
	
	// 获取客户端
	client, err := m.GetClient(targetProvider)
	if err != nil {
		return nil, err
	}
	
	// 调用AI生成文档
	document, err := client.GenerateDocument(ctx, analysis)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果
	if m.cache != nil {
		cacheKey := m.generateCacheKey("document", targetProvider, analysis.ID)
		m.cache.Set(cacheKey, document, 60*time.Minute)
	}
	
	return document, nil
}

// ListProviders 列出所有可用的AI提供商
func (m *AIManager) ListProviders() []AIProvider {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	providers := make([]AIProvider, 0, len(m.clients))
	for provider := range m.clients {
		providers = append(providers, provider)
	}
	
	return providers
}

// GetDefaultProvider 获取默认提供商
func (m *AIManager) GetDefaultProvider() AIProvider {
	return m.defaultProvider
}

// SetDefaultProvider 设置默认提供商
func (m *AIManager) SetDefaultProvider(provider AIProvider) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.clients[provider]; !exists {
		return fmt.Errorf("AI提供商 %s 未配置", provider)
	}
	
	m.defaultProvider = provider
	return nil
}

// ClearCache 清空缓存
func (m *AIManager) ClearCache() {
	if m.cache != nil {
		m.cache.Clear()
	}
}

// generateCacheKey 生成缓存键
func (m *AIManager) generateCacheKey(operation string, provider AIProvider, params ...string) string {
	data := []string{operation, string(provider)}
	data = append(data, params...)
	
	// 将参数序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// 如果序列化失败，使用简单的字符串拼接
		result := operation + ":" + string(provider)
		for _, param := range params {
			result += ":" + param
		}
		return result
	}
	
	// 使用MD5生成哈希
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("%x", hash)
}

// GetCacheStats 获取缓存统计信息
func (m *AIManager) GetCacheStats() map[string]interface{} {
	if m.cache == nil {
		return map[string]interface{}{
			"enabled": false,
		}
	}
	
	memCache, ok := m.cache.(*MemoryCache)
	if !ok {
		return map[string]interface{}{
			"enabled": true,
			"type":    "unknown",
		}
	}
	
	memCache.mutex.RLock()
	defer memCache.mutex.RUnlock()
	
	return map[string]interface{}{
		"enabled": true,
		"type":    "memory",
		"size":    len(memCache.data),
	}
}

// ProjectChat 项目上下文AI对话（带缓存）
func (m *AIManager) ProjectChat(ctx context.Context, message, context string, provider ...AIProvider) (*ProjectChatResponse, error) {
	// 确定使用的提供商
	targetProvider := m.defaultProvider
	if len(provider) > 0 {
		targetProvider = provider[0]
	}
	
	// 获取客户端
	client, err := m.GetClient(targetProvider)
	if err != nil {
		return nil, fmt.Errorf("获取AI客户端失败: %w", err)
	}
	
	// 调用客户端进行对话
	response, err := client.ProjectChat(ctx, message, context)
	if err != nil {
		return nil, fmt.Errorf("AI对话失败: %w", err)
	}
	
	return response, nil
} 