package ai

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockAIClient 模拟AI客户端
type MockAIClient struct {
	mock.Mock
	provider AIProvider
}

func (m *MockAIClient) AnalyzeRequirement(ctx context.Context, requirement string) (*RequirementAnalysis, error) {
	args := m.Called(ctx, requirement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RequirementAnalysis), args.Error(1)
}

func (m *MockAIClient) GenerateQuestions(ctx context.Context, analysis *RequirementAnalysis) ([]Question, error) {
	args := m.Called(ctx, analysis)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Question), args.Error(1)
}

func (m *MockAIClient) GeneratePUML(ctx context.Context, analysis *RequirementAnalysis, diagramType PUMLType) (*PUMLDiagram, error) {
	args := m.Called(ctx, analysis, diagramType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PUMLDiagram), args.Error(1)
}

func (m *MockAIClient) GenerateDocument(ctx context.Context, analysis *RequirementAnalysis) (*DevelopmentDocument, error) {
	args := m.Called(ctx, analysis)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DevelopmentDocument), args.Error(1)
}

func (m *MockAIClient) ProjectChat(ctx context.Context, message, context string) (*ProjectChatResponse, error) {
	args := m.Called(ctx, message, context)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ProjectChatResponse), args.Error(1)
}

func (m *MockAIClient) GetProvider() AIProvider {
	return m.provider
}

type MemoryCacheTestSuite struct {
	suite.Suite
	cache *MemoryCache
}

func (suite *MemoryCacheTestSuite) SetupTest() {
	suite.cache = NewMemoryCache()
}

func (suite *MemoryCacheTestSuite) TestMemoryCache_SetAndGet() {
	// Arrange
	key := "test_key"
	value := "test_value"
	ttl := time.Hour

	// Act
	suite.cache.Set(key, value, ttl)
	result, found := suite.cache.Get(key)

	// Assert
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, result)
}

func (suite *MemoryCacheTestSuite) TestMemoryCache_GetNonExistent() {
	// Act
	result, found := suite.cache.Get("nonexistent")

	// Assert
	assert.False(suite.T(), found)
	assert.Nil(suite.T(), result)
}

func (suite *MemoryCacheTestSuite) TestMemoryCache_ExpiredItem() {
	// Arrange
	key := "expired_key"
	value := "expired_value"
	ttl := time.Millisecond * 10

	// Act
	suite.cache.Set(key, value, ttl)
	time.Sleep(time.Millisecond * 20) // 等待过期
	result, found := suite.cache.Get(key)

	// Assert
	assert.False(suite.T(), found)
	assert.Nil(suite.T(), result)
}

func (suite *MemoryCacheTestSuite) TestMemoryCache_Delete() {
	// Arrange
	key := "delete_key"
	value := "delete_value"
	ttl := time.Hour

	suite.cache.Set(key, value, ttl)

	// Act
	suite.cache.Delete(key)
	result, found := suite.cache.Get(key)

	// Assert
	assert.False(suite.T(), found)
	assert.Nil(suite.T(), result)
}

func (suite *MemoryCacheTestSuite) TestMemoryCache_Clear() {
	// Arrange
	suite.cache.Set("key1", "value1", time.Hour)
	suite.cache.Set("key2", "value2", time.Hour)

	// Act
	suite.cache.Clear()

	// Assert
	_, found1 := suite.cache.Get("key1")
	_, found2 := suite.cache.Get("key2")
	assert.False(suite.T(), found1)
	assert.False(suite.T(), found2)
}

type AIManagerTestSuite struct {
	suite.Suite
	mockOpenAI  *MockAIClient
	mockGemini  *MockAIClient
	manager     *AIManager
}

func (suite *AIManagerTestSuite) SetupTest() {
	suite.mockOpenAI = &MockAIClient{provider: ProviderOpenAI}
	suite.mockGemini = &MockAIClient{provider: ProviderGemini}

	// 创建AI管理器并手动设置模拟客户端
	suite.manager = &AIManager{
		clients:         make(map[AIProvider]AIClient),
		defaultProvider: ProviderOpenAI,
		cache:           NewMemoryCache(),
	}

	suite.manager.clients[ProviderOpenAI] = suite.mockOpenAI
	suite.manager.clients[ProviderGemini] = suite.mockGemini
}

func (suite *AIManagerTestSuite) TearDownTest() {
	suite.mockOpenAI.AssertExpectations(suite.T())
	suite.mockGemini.AssertExpectations(suite.T())
}

func (suite *AIManagerTestSuite) TestNewAIManager_Success() {
	// Arrange
	config := AIManagerConfig{
		DefaultProvider: ProviderOpenAI,
		OpenAIConfig: &OpenAIConfig{
			APIKey: "test-key",
		},
		EnableCache: true,
	}

	// Act
	manager, err := NewAIManager(config)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), manager)
	assert.Equal(suite.T(), ProviderOpenAI, manager.defaultProvider)
	assert.NotNil(suite.T(), manager.cache)
}

func (suite *AIManagerTestSuite) TestNewAIManager_DefaultProviderNotConfigured() {
	// Arrange
	config := AIManagerConfig{
		DefaultProvider: ProviderClaude, // 未配置的提供商
		EnableCache:     true,
	}

	// Act
	manager, err := NewAIManager(config)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), manager)
	assert.Contains(suite.T(), err.Error(), "默认AI提供商")
}

func (suite *AIManagerTestSuite) TestGetClient_Success() {
	// Act
	client, err := suite.manager.GetClient(ProviderOpenAI)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.mockOpenAI, client)
}

func (suite *AIManagerTestSuite) TestGetClient_NotConfigured() {
	// Act
	client, err := suite.manager.GetClient(ProviderClaude)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), client)
	assert.Contains(suite.T(), err.Error(), "未配置")
}

func (suite *AIManagerTestSuite) TestGetDefaultClient() {
	// Act
	client := suite.manager.GetDefaultClient()

	// Assert
	assert.Equal(suite.T(), suite.mockOpenAI, client)
}

func (suite *AIManagerTestSuite) TestAnalyzeRequirement_Success() {
	// Arrange
	ctx := context.Background()
	requirement := "创建一个用户管理系统"
	
	expectedAnalysis := &RequirementAnalysis{
		ID:           "test-id",
		OriginalText: requirement,
		CoreFunctions: []string{"用户注册", "用户登录"},
		CompletionScore: 0.8,
	}

	suite.mockOpenAI.On("AnalyzeRequirement", ctx, requirement).Return(expectedAnalysis, nil)

	// Act
	result, err := suite.manager.AnalyzeRequirement(ctx, requirement)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAnalysis, result)
}

func (suite *AIManagerTestSuite) TestAnalyzeRequirement_WithSpecificProvider() {
	// Arrange
	ctx := context.Background()
	requirement := "创建一个用户管理系统"
	
	expectedAnalysis := &RequirementAnalysis{
		ID:           "test-id",
		OriginalText: requirement,
		CoreFunctions: []string{"用户注册", "用户登录"},
		CompletionScore: 0.8,
	}

	suite.mockGemini.On("AnalyzeRequirement", ctx, requirement).Return(expectedAnalysis, nil)

	// Act
	result, err := suite.manager.AnalyzeRequirement(ctx, requirement, ProviderGemini)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAnalysis, result)
}

func (suite *AIManagerTestSuite) TestAnalyzeRequirement_WithCache() {
	// Arrange
	ctx := context.Background()
	requirement := "创建一个用户管理系统"
	
	expectedAnalysis := &RequirementAnalysis{
		ID:           "test-id",
		OriginalText: requirement,
		CoreFunctions: []string{"用户注册", "用户登录"},
		CompletionScore: 0.8,
	}

	// 第一次调用AI客户端
	suite.mockOpenAI.On("AnalyzeRequirement", ctx, requirement).Return(expectedAnalysis, nil).Once()

	// Act - 第一次调用
	result1, err1 := suite.manager.AnalyzeRequirement(ctx, requirement)
	
	// Act - 第二次调用，应该从缓存返回
	result2, err2 := suite.manager.AnalyzeRequirement(ctx, requirement)

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.Equal(suite.T(), expectedAnalysis, result1)
	assert.Equal(suite.T(), expectedAnalysis, result2)
}

func (suite *AIManagerTestSuite) TestGenerateQuestions_Success() {
	// Arrange
	ctx := context.Background()
	analysis := &RequirementAnalysis{
		ID: "test-id",
		CoreFunctions: []string{"用户注册", "用户登录"},
	}
	
	expectedQuestions := []Question{
		{
			ID:       "q1",
			Category: "business_rule",
			Content:  "用户注册需要哪些字段？",
			Priority: 3,
		},
	}

	suite.mockOpenAI.On("GenerateQuestions", ctx, analysis).Return(expectedQuestions, nil)

	// Act
	result, err := suite.manager.GenerateQuestions(ctx, analysis)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedQuestions, result)
}

func (suite *AIManagerTestSuite) TestGeneratePUML_Success() {
	// Arrange
	ctx := context.Background()
	analysis := &RequirementAnalysis{
		ID: "test-id",
		CoreFunctions: []string{"用户注册", "用户登录"},
	}
	diagramType := PUMLTypeBusinessFlow
	
	expectedDiagram := &PUMLDiagram{
		ID:        "diagram-id",
		Type:      diagramType,
		Title:     "用户管理系统业务流程",
		Content:   "@startuml\n...\n@enduml",
	}

	suite.mockOpenAI.On("GeneratePUML", ctx, analysis, diagramType).Return(expectedDiagram, nil)

	// Act
	result, err := suite.manager.GeneratePUML(ctx, analysis, diagramType)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDiagram, result)
}

func (suite *AIManagerTestSuite) TestGenerateDocument_Success() {
	// Arrange
	ctx := context.Background()
	analysis := &RequirementAnalysis{
		ID: "test-id",
		CoreFunctions: []string{"用户注册", "用户登录"},
	}
	
	expectedDocument := &DevelopmentDocument{
		ID:        "doc-id",
		ProjectID: analysis.ID,
		FunctionModules: []FunctionModule{
			{
				Name:        "用户管理",
				Description: "负责用户注册、登录等功能",
				Priority:    1,
			},
		},
	}

	suite.mockOpenAI.On("GenerateDocument", ctx, analysis).Return(expectedDocument, nil)

	// Act
	result, err := suite.manager.GenerateDocument(ctx, analysis)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDocument, result)
}

func (suite *AIManagerTestSuite) TestProjectChat_Success() {
	// Arrange
	ctx := context.Background()
	message := "如何实现用户登录功能？"
	context := "用户管理系统项目"
	
	expectedResponse := &ProjectChatResponse{
		Message:              "可以使用JWT实现用户登录...",
		ShouldUpdateAnalysis: false,
		RelatedQuestions:     []string{"如何处理密码加密？"},
		Suggestions:          []string{"建议使用bcrypt加密密码"},
	}

	suite.mockOpenAI.On("ProjectChat", ctx, message, context).Return(expectedResponse, nil)

	// Act
	result, err := suite.manager.ProjectChat(ctx, message, context)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse, result)
}

func (suite *AIManagerTestSuite) TestListProviders() {
	// Act
	providers := suite.manager.ListProviders()

	// Assert
	assert.Contains(suite.T(), providers, ProviderOpenAI)
	assert.Contains(suite.T(), providers, ProviderGemini)
	assert.Len(suite.T(), providers, 2)
}

func (suite *AIManagerTestSuite) TestGetDefaultProvider() {
	// Act
	provider := suite.manager.GetDefaultProvider()

	// Assert
	assert.Equal(suite.T(), ProviderOpenAI, provider)
}

func (suite *AIManagerTestSuite) TestSetDefaultProvider_Success() {
	// Act
	err := suite.manager.SetDefaultProvider(ProviderGemini)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), ProviderGemini, suite.manager.defaultProvider)
}

func (suite *AIManagerTestSuite) TestSetDefaultProvider_ProviderNotConfigured() {
	// Act
	err := suite.manager.SetDefaultProvider(ProviderClaude)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "未配置")
	assert.Equal(suite.T(), ProviderOpenAI, suite.manager.defaultProvider) // 应该保持原值
}

func (suite *AIManagerTestSuite) TestClearCache() {
	// Arrange
	ctx := context.Background()
	requirement := "测试需求"
	
	expectedAnalysis := &RequirementAnalysis{ID: "test-id"}
	suite.mockOpenAI.On("AnalyzeRequirement", ctx, requirement).Return(expectedAnalysis, nil).Twice()

	// 先调用一次，填充缓存
	_, err := suite.manager.AnalyzeRequirement(ctx, requirement)
	assert.NoError(suite.T(), err)

	// Act
	suite.manager.ClearCache()

	// 再次调用，应该重新请求AI服务（因为缓存已清空）
	_, err = suite.manager.AnalyzeRequirement(ctx, requirement)

	// Assert
	assert.NoError(suite.T(), err)
	// mockOpenAI.On 中的 .Twice() 确保了被调用了两次
}

func (suite *AIManagerTestSuite) TestGenerateCacheKey() {
	// Arrange
	operation := "analyze"
	provider := ProviderOpenAI
	param1 := "param1"
	param2 := "param2"

	// Act
	key1 := suite.manager.generateCacheKey(operation, provider, param1, param2)
	key2 := suite.manager.generateCacheKey(operation, provider, param1, param2)
	key3 := suite.manager.generateCacheKey("different", provider, param1, param2)

	// Assert
	assert.Equal(suite.T(), key1, key2) // 相同参数应该生成相同的key
	assert.NotEqual(suite.T(), key1, key3) // 不同参数应该生成不同的key
	assert.NotEmpty(suite.T(), key1)
}

func (suite *AIManagerTestSuite) TestGetCacheStats() {
	// Act
	stats := suite.manager.GetCacheStats()

	// Assert
	assert.NotNil(suite.T(), stats)
	assert.Contains(suite.T(), stats, "enabled")
	assert.True(suite.T(), stats["enabled"].(bool))
}

func TestMemoryCacheTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryCacheTestSuite))
}

func TestAIManagerTestSuite(t *testing.T) {
	suite.Run(t, new(AIManagerTestSuite))
} 