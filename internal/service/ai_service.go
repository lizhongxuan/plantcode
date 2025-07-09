package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"ai-dev-platform/internal/ai"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/repository"

	"github.com/google/uuid"
)

// AIService AI服务层
type AIService struct {
	aiManager *ai.AIManager
	aiRepo    *repository.AIRepository
	mysqlRepo *repository.MySQLRepository
}

// NewAIService 创建AI服务
func NewAIService(aiManager *ai.AIManager, aiRepo *repository.AIRepository, mysqlRepo *repository.MySQLRepository) *AIService {
	return &AIService{
		aiManager: aiManager,
		aiRepo:    aiRepo,
		mysqlRepo: mysqlRepo,
	}
}

// ===== 需求分析相关服务 =====

// AnalyzeRequirementWithUser 基于用户AI配置分析业务需求
func (s *AIService) AnalyzeRequirementWithUser(ctx context.Context, req *model.AIAnalysisRequest, userID uuid.UUID) (*model.Requirement, error) {
	// 验证项目是否存在
	_, err := s.mysqlRepo.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	// 获取用户AI配置
	userConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil {
		// 如果没有用户配置，使用默认配置
		log.Printf("获取用户AI配置失败，使用默认配置: %v", err)
		return s.AnalyzeRequirement(ctx, req)
	}

	// 确定使用的provider
	provider := ai.AIProvider(userConfig.Provider)
	var apiKey string

	switch userConfig.Provider {
	case "openai":
		apiKey = userConfig.OpenAIAPIKey
	case "claude":
		apiKey = userConfig.ClaudeAPIKey
	case "gemini":
		apiKey = userConfig.GeminiAPIKey
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", userConfig.Provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未配置%s的API密钥，请先在设置中配置", userConfig.Provider)
	}

	// 创建临时AI客户端配置
	clientConfig := ai.AIManagerConfig{
		DefaultProvider: provider,
		EnableCache:     true,
		CacheTTL:        time.Hour,
	}

	// 根据provider设置对应的客户端配置
	switch provider {
	case ai.ProviderOpenAI:
		clientConfig.OpenAIConfig = &ai.OpenAIConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	case ai.ProviderClaude:
		// Claude客户端暂时不可用
		return nil, fmt.Errorf("Claude客户端暂时不可用")
	case ai.ProviderGemini:
		clientConfig.GeminiConfig = &ai.GeminiConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	}

	// 创建用户特定的AI管理器
	tempAIManager, err := ai.NewAIManager(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AI管理器失败: %w", err)
	}

	// 调用AI分析
	analysis, err := tempAIManager.AnalyzeRequirement(ctx, req.Requirement, provider)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 转换为数据库模型并保存
	return s.saveAnalysisResult(req, analysis, tempAIManager, provider)
}

// AnalyzeRequirement 分析业务需求（原有方法，作为兼容性保留）
func (s *AIService) AnalyzeRequirement(ctx context.Context, req *model.AIAnalysisRequest) (*model.Requirement, error) {
	// 验证项目是否存在
	_, err := s.mysqlRepo.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	// 确定AI提供商
	provider := ai.ProviderOpenAI
	if req.Provider != "" {
		provider = ai.AIProvider(req.Provider)
	}

	// 调用AI分析
	analysis, err := s.aiManager.AnalyzeRequirement(ctx, req.Requirement, provider)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 转换为数据库模型并保存
	return s.saveAnalysisResult(req, analysis, s.aiManager, provider)
}

// saveAnalysisResult 保存分析结果到数据库的公共方法
func (s *AIService) saveAnalysisResult(req *model.AIAnalysisRequest, analysis *ai.RequirementAnalysis, aiManager *ai.AIManager, provider ai.AIProvider) (*model.Requirement, error) {
	// 转换为数据库模型
	dbAnalysis := &model.Requirement{
		RequirementID:     uuid.New(),
		ProjectID:         req.ProjectID,
		RawRequirement:    req.Requirement,
		CompletenessScore: analysis.CompletionScore,
		AnalysisStatus:    model.AnalysisStatusCompleted,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 序列化结构化需求
	structuredReq := map[string]interface{}{
		"core_functions":     analysis.CoreFunctions,
		"roles":              analysis.Roles,
		"business_processes": analysis.BusinessProcesses,
		"data_entities":      analysis.DataEntities,
	}
	structuredJSON, err := json.Marshal(structuredReq)
	if err != nil {
		return nil, fmt.Errorf("序列化结构化需求失败: %w", err)
	}
	dbAnalysis.StructuredRequirement = string(structuredJSON)

	// 序列化缺失信息
	missingInfoJSON, err := json.Marshal(analysis.MissingInfo)
	if err != nil {
		return nil, fmt.Errorf("序列化缺失信息失败: %w", err)
	}
	dbAnalysis.MissingInfoTypes = string(missingInfoJSON)

	// 保存到数据库
	err = s.aiRepo.CreateRequirementAnalysis(dbAnalysis)
	if err != nil {
		return nil, fmt.Errorf("保存需求分析失败: %w", err)
	}

	// 如果有缺失信息，生成补充问题
	if len(analysis.MissingInfo) > 0 {
		go func() {
			questions, err := aiManager.GenerateQuestions(context.Background(), analysis, provider)
			if err != nil {
				log.Printf("生成补充问题失败: %v", err)
				return
			}

			// 保存问题到数据库
			for _, question := range questions {
				dbQuestion := &model.Question{
					QuestionID:       uuid.New(),
					RequirementID:    dbAnalysis.RequirementID,
					QuestionText:     question.Content,
					QuestionCategory: question.Category,
					PriorityLevel:    question.Priority,
					AnswerStatus:     model.QuestionStatusPending,
					CreatedAt:        time.Now(),
				}

				if err := s.aiRepo.CreateQuestion(dbQuestion); err != nil {
					log.Printf("保存问题失败: %v", err)
				}
			}
		}()
	}

	return dbAnalysis, nil
}

// GetRequirementAnalysis 获取需求分析
func (s *AIService) GetRequirementAnalysis(analysisID uuid.UUID) (*model.Requirement, error) {
	return s.aiRepo.GetRequirementAnalysis(analysisID)
}

// GetRequirementAnalysesByProject 获取项目的需求分析列表
func (s *AIService) GetRequirementAnalysesByProject(projectID uuid.UUID) ([]*model.Requirement, error) {
	return s.aiRepo.GetRequirementAnalysesByProject(projectID)
}

// GetQuestions 获取需求分析的补充问题
func (s *AIService) GetQuestions(requirementID uuid.UUID) ([]*model.Question, error) {
	return s.aiRepo.GetQuestions(requirementID)
}

// AnswerQuestion 回答补充问题
func (s *AIService) AnswerQuestion(questionID uuid.UUID, answer string) error {
	return s.aiRepo.AnswerQuestion(questionID, answer)
}

// ===== PUML图表生成相关服务 =====

// GeneratePUML 生成PUML图表
func (s *AIService) GeneratePUML(ctx context.Context, req *model.GeneratePUMLRequest) (*model.PUMLDiagram, error) {
	// 获取需求分析
	analysisUUID, err := uuid.Parse(req.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("无效的分析ID: %w", err)
	}

	dbAnalysis, err := s.aiRepo.GetRequirementAnalysis(analysisUUID)
	if err != nil {
		return nil, fmt.Errorf("获取需求分析失败: %w", err)
	}

	// 解析结构化需求
	var structuredReq map[string]interface{}
	if err := json.Unmarshal([]byte(dbAnalysis.StructuredRequirement), &structuredReq); err != nil {
		return nil, fmt.Errorf("解析结构化需求失败: %w", err)
	}

	// 构建AI分析对象
	analysis := &ai.RequirementAnalysis{
		ID:           req.AnalysisID,
		ProjectID:    dbAnalysis.ProjectID.String(),
		OriginalText: dbAnalysis.RawRequirement,
	}

	// 提取核心功能
	if coreFuncs, ok := structuredReq["core_functions"].([]interface{}); ok {
		for _, fn := range coreFuncs {
			if funcStr, ok := fn.(string); ok {
				analysis.CoreFunctions = append(analysis.CoreFunctions, funcStr)
			}
		}
	}

	// 确定AI提供商
	provider := ai.ProviderOpenAI
	if req.Provider != "" {
		provider = ai.AIProvider(req.Provider)
	}

	// 调用AI生成PUML
	diagramType := ai.PUMLType(req.DiagramType)
	diagram, err := s.aiManager.GeneratePUML(ctx, analysis, diagramType, provider)
	if err != nil {
		return nil, fmt.Errorf("AI生成PUML失败: %w", err)
	}

	// 转换为数据库模型
	dbDiagram := &model.PUMLDiagram{
		DiagramID:   uuid.New(),
		ProjectID:   dbAnalysis.ProjectID,
		DiagramType: req.DiagramType,
		DiagramName: diagram.Title,
		PUMLContent: diagram.Content,
		Version:     1,
		IsValidated: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	err = s.aiRepo.CreatePUMLDiagram(dbDiagram)
	if err != nil {
		return nil, fmt.Errorf("保存PUML图表失败: %w", err)
	}

	return dbDiagram, nil
}

// GetPUMLDiagram 获取PUML图表
func (s *AIService) GetPUMLDiagram(diagramID uuid.UUID) (*model.PUMLDiagram, error) {
	return s.aiRepo.GetPUMLDiagram(diagramID)
}

// GetPUMLDiagramsByProject 获取项目的PUML图表列表
func (s *AIService) GetPUMLDiagramsByProject(projectID uuid.UUID) ([]*model.PUMLDiagram, error) {
	return s.aiRepo.GetPUMLDiagramsByProject(projectID)
}

// UpdatePUMLDiagram 更新PUML图表
func (s *AIService) UpdatePUMLDiagram(diagramID uuid.UUID, req *model.UpdatePUMLRequest) error {
	// 获取现有图表
	diagram, err := s.aiRepo.GetPUMLDiagram(diagramID)
	if err != nil {
		return fmt.Errorf("获取PUML图表失败: %w", err)
	}

	// 更新字段
	if req.Title != "" {
		diagram.DiagramName = req.Title
	}
	diagram.PUMLContent = req.Content
	if req.Description != "" {
		diagram.ValidationFeedback = req.Description
	}
	diagram.Version++

	// 保存更新
	return s.aiRepo.UpdatePUMLDiagram(diagram)
}

// CreatePUML 创建PUML图表
func (s *AIService) CreatePUML(diagram *model.PUMLDiagram) error {
	return s.aiRepo.CreatePUMLDiagram(diagram)
}

// DeletePUMLDiagram 删除PUML图表
func (s *AIService) DeletePUMLDiagram(diagramID uuid.UUID) error {
	return s.mysqlRepo.DeletePUMLDiagram(diagramID)
}

// ===== 文档生成相关服务 =====

// GenerateDocument 生成开发文档
func (s *AIService) GenerateDocument(ctx context.Context, req *model.GenerateDocumentRequest) (*model.Document, error) {
	// 获取需求分析
	analysisUUID, err := uuid.Parse(req.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("无效的分析ID: %w", err)
	}

	dbAnalysis, err := s.aiRepo.GetRequirementAnalysis(analysisUUID)
	if err != nil {
		return nil, fmt.Errorf("获取需求分析失败: %w", err)
	}

	// 解析结构化需求
	var structuredReq map[string]interface{}
	if err := json.Unmarshal([]byte(dbAnalysis.StructuredRequirement), &structuredReq); err != nil {
		return nil, fmt.Errorf("解析结构化需求失败: %w", err)
	}

	// 构建AI分析对象
	analysis := &ai.RequirementAnalysis{
		ID:           req.AnalysisID,
		ProjectID:    dbAnalysis.ProjectID.String(),
		OriginalText: dbAnalysis.RawRequirement,
	}

	// 提取核心功能
	if coreFuncs, ok := structuredReq["core_functions"].([]interface{}); ok {
		for _, fn := range coreFuncs {
			if funcStr, ok := fn.(string); ok {
				analysis.CoreFunctions = append(analysis.CoreFunctions, funcStr)
			}
		}
	}

	// 确定AI提供商
	provider := ai.ProviderOpenAI
	if req.Provider != "" {
		provider = ai.AIProvider(req.Provider)
	}

	// 调用AI生成文档
	document, err := s.aiManager.GenerateDocument(ctx, analysis, provider)
	if err != nil {
		return nil, fmt.Errorf("AI生成文档失败: %w", err)
	}

	// 序列化文档内容
	contentJSON, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("序列化文档内容失败: %w", err)
	}

	// 转换为数据库模型
	dbDocument := &model.Document{
		DocumentID:   uuid.New(),
		ProjectID:    dbAnalysis.ProjectID,
		DocumentType: "development_plan",
		DocumentName: "AI生成的开发文档",
		Content:      string(contentJSON),
		Format:       "json",
		Version:      1,
		GeneratedAt:  time.Now(),
		IsFinal:      false,
	}

	// 保存到数据库
	err = s.aiRepo.CreateDocument(dbDocument)
	if err != nil {
		return nil, fmt.Errorf("保存生成文档失败: %w", err)
	}

	return dbDocument, nil
}

// GetDocument 获取生成的文档
func (s *AIService) GetDocument(documentID uuid.UUID) (*model.Document, error) {
	return s.aiRepo.GetDocument(documentID)
}

// GetDocumentsByProject 获取项目的文档列表
func (s *AIService) GetDocumentsByProject(projectID uuid.UUID) ([]*model.Document, error) {
	return s.aiRepo.GetDocumentsByProject(projectID)
}

// UpdateDocument 更新文档
func (s *AIService) UpdateDocument(documentID uuid.UUID, req *model.UpdateDocumentRequest) error {
	// 获取现有文档
	document, err := s.aiRepo.GetDocument(documentID)
	if err != nil {
		return fmt.Errorf("获取文档失败: %w", err)
	}

	// 更新字段
	if req.Title != "" {
		document.DocumentName = req.Title
	}
	document.Content = req.Content
	document.Version++

	// 保存更新
	return s.aiRepo.UpdateDocument(document)
}

// ===== 对话相关服务 =====

// CreateChatSession 创建对话会话
func (s *AIService) CreateChatSession(req *model.ChatSessionCreateRequest, userID uuid.UUID) (*model.ChatSession, error) {
	var projectUUID uuid.UUID
	var err error

	// 如果提供了项目ID，验证项目是否存在
	if req.ProjectID != "" {
		projectUUID, err = uuid.Parse(req.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("无效的项目ID: %w", err)
		}

		_, err = s.mysqlRepo.GetProjectByID(projectUUID)
		if err != nil {
			return nil, fmt.Errorf("项目不存在: %w", err)
		}
	} else {
		// 对于通用聊天，使用空的UUID
		projectUUID = uuid.Nil
	}

	// 创建会话
	session := &model.ChatSession{
		SessionID:   uuid.New(),
		ProjectID:   projectUUID,
		UserID:      userID,
		SessionType: model.SessionTypeRequirementAnalysis,
		StartedAt:   time.Now(),
		Status:      model.ChatSessionStatusActive,
		Context:     "{}",
	}

	err = s.aiRepo.CreateChatSession(session)
	if err != nil {
		return nil, fmt.Errorf("创建对话会话失败: %w", err)
	}

	return session, nil
}

// SendChatMessage 发送聊天消息
func (s *AIService) SendChatMessage(ctx context.Context, req *model.SendChatMessageRequest) (*model.ChatMessage, error) {
	// 验证会话是否存在
	sessionUUID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("无效的会话ID: %w", err)
	}

	session, err := s.aiRepo.GetChatSession(sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("会话不存在: %w", err)
	}

	if session.Status != model.ChatSessionStatusActive {
		return nil, fmt.Errorf("会话已结束")
	}

	// 创建用户消息
	userMessage := &model.ChatMessage{
		MessageID:      uuid.New(),
		SessionID:      sessionUUID,
		SenderType:     req.Role,
		MessageContent: req.Content,
		MessageType:    model.MessageTypeText,
		Metadata:       "{}",
		Timestamp:      time.Now(),
		Processed:      false,
	}

	// 保存用户消息
	err = s.aiRepo.CreateChatMessage(userMessage)
	if err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 如果是用户消息，生成AI回复
	if req.Role == model.MessageRoleUser {
		go func() {
			// 这里可以调用AI服务生成回复
			// 当前为简化实现，仅返回确认消息
			aiMessage := &model.ChatMessage{
				MessageID:      uuid.New(),
				SessionID:      sessionUUID,
				SenderType:     model.MessageRoleAssistant,
				MessageContent: "我收到了您的消息：" + req.Content + "。我正在分析中，请稍后...",
				MessageType:    model.MessageTypeText,
				Metadata:       "{}",
				Timestamp:      time.Now(),
				Processed:      true,
			}

			if err := s.aiRepo.CreateChatMessage(aiMessage); err != nil {
				log.Printf("保存AI消息失败: %v", err)
			}
		}()
	}

	return userMessage, nil
}

// GetChatMessages 获取对话消息
func (s *AIService) GetChatMessages(sessionID uuid.UUID) ([]*model.ChatMessage, error) {
	return s.aiRepo.GetChatMessages(sessionID)
}

// GetChatSessionsByProject 获取项目的对话会话列表
func (s *AIService) GetChatSessionsByProject(projectID uuid.UUID) ([]*model.ChatSession, error) {
	return s.aiRepo.GetChatSessionsByProject(projectID)
}

// ===== 系统管理相关服务 =====

// GetAIProviders 获取可用的AI提供商列表
func (s *AIService) GetAIProviders() []string {
	providers := s.aiManager.ListProviders()
	result := make([]string, len(providers))
	for i, p := range providers {
		result[i] = string(p)
	}
	return result
}

// GetDefaultAIProvider 获取默认AI提供商
func (s *AIService) GetDefaultAIProvider() string {
	return string(s.aiManager.GetDefaultProvider())
}

// SetDefaultAIProvider 设置默认AI提供商
func (s *AIService) SetDefaultAIProvider(provider string) error {
	return s.aiManager.SetDefaultProvider(ai.AIProvider(provider))
}

// GetCacheStats 获取缓存统计信息
func (s *AIService) GetCacheStats() map[string]interface{} {
	return s.aiManager.GetCacheStats()
}

// ClearCache 清空AI缓存
func (s *AIService) ClearCache() {
	s.aiManager.ClearCache()
}

// ===== 用户AI配置管理相关服务 =====

// GetUserAIConfig 获取用户AI配置
func (s *AIService) GetUserAIConfig(userID uuid.UUID) (*model.UserAIConfig, error) {
	return s.aiRepo.GetUserAIConfig(userID)
}

// UpdateUserAIConfig 更新用户AI配置
func (s *AIService) UpdateUserAIConfig(userID uuid.UUID, req *model.UpdateUserAIConfigRequest) (*model.UserAIConfig, error) {
	// 检查是否已有配置
	existingConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, fmt.Errorf("检查现有配置失败: %w", err)
	}

	var config *model.UserAIConfig
	now := time.Now()

	if existingConfig != nil {
		// 更新现有配置
		config = existingConfig
		config.Provider = req.Provider
		config.DefaultModel = req.DefaultModel
		config.MaxTokens = req.MaxTokens
		config.UpdatedAt = now
		config.IsActive = true

		// 更新API密钥（只更新非空的密钥）
		if req.OpenAIAPIKey != "" {
			config.OpenAIAPIKey = req.OpenAIAPIKey
		}
		if req.ClaudeAPIKey != "" {
			config.ClaudeAPIKey = req.ClaudeAPIKey
		}
		if req.GeminiAPIKey != "" {
			config.GeminiAPIKey = req.GeminiAPIKey
		}

		err = s.aiRepo.UpdateUserAIConfig(config)
	} else {
		// 创建新配置
		config = &model.UserAIConfig{
			ConfigID:     uuid.New(),
			UserID:       userID,
			Provider:     req.Provider,
			OpenAIAPIKey: req.OpenAIAPIKey,
			ClaudeAPIKey: req.ClaudeAPIKey,
			GeminiAPIKey: req.GeminiAPIKey,
			DefaultModel: req.DefaultModel,
			MaxTokens:    req.MaxTokens,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		err = s.aiRepo.CreateUserAIConfig(config)
	}

	if err != nil {
		return nil, fmt.Errorf("保存AI配置失败: %w", err)
	}

	return config, nil
}

// TestAIConnection 测试AI连接
func (s *AIService) TestAIConnection(req *model.TestAIConnectionRequest) (*model.AIConnectionTestResult, error) {
	start := time.Now()

	result := &model.AIConnectionTestResult{
		Provider: req.Provider,
		Model:    req.Model,
	}

	switch req.Provider {
	case "openai":
		err := s.testOpenAIConnection(req.APIKey, req.Model)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "OpenAI连接测试成功"
		}
	case "claude":
		err := s.testClaudeConnection(req.APIKey, req.Model)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "Claude连接测试成功"
		}
	case "gemini":
		err := s.testGeminiConnection(req.APIKey, req.Model)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "Gemini连接测试成功"
		}
	default:
		result.Success = false
		result.Message = "不支持的AI提供商"
	}

	result.Latency = time.Since(start).Milliseconds()
	return result, nil
}

// testOpenAIConnection 测试OpenAI连接
func (s *AIService) testOpenAIConnection(apiKey, model string) error {
	// 这里应该调用OpenAI API进行简单的测试请求
	// 暂时模拟测试
	if apiKey == "" {
		return fmt.Errorf("OpenAI API密钥不能为空")
	}

	if !strings.HasPrefix(apiKey, "sk-") {
		return fmt.Errorf("OpenAI API密钥格式无效")
	}

	// 模拟网络请求延迟
	time.Sleep(500 * time.Millisecond)

	// 这里可以添加真实的OpenAI API测试
	// 例如发送一个简单的completion请求

	return nil
}

// testClaudeConnection 测试Claude连接
func (s *AIService) testClaudeConnection(apiKey, model string) error {
	// 这里应该调用Claude API进行简单的测试请求
	// 暂时模拟测试
	if apiKey == "" {
		return fmt.Errorf("Claude API密钥不能为空")
	}

	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return fmt.Errorf("Claude API密钥格式无效")
	}

	// 模拟网络请求延迟
	time.Sleep(800 * time.Millisecond)

	// 这里可以添加真实的Claude API测试

	return nil
}

// testGeminiConnection 测试Gemini连接
func (s *AIService) testGeminiConnection(apiKey, model string) error {
	// 这里应该调用Gemini API进行简单的测试请求
	// 暂时模拟测试
	if apiKey == "" {
		return fmt.Errorf("Gemini API密钥不能为空")
	}

	// Gemini API密钥格式通常以"AIza"开头
	if !strings.HasPrefix(apiKey, "AIza") {
		return fmt.Errorf("Gemini API密钥格式无效")
	}

	// 模拟网络请求延迟
	time.Sleep(600 * time.Millisecond)

	// 这里可以添加真实的Gemini API测试
	// 例如发送一个简单的generateContent请求

	return nil
}

// GetAvailableModels 获取可用的AI模型列表
func (s *AIService) GetAvailableModels(provider string) ([]string, error) {
	switch provider {
	case "openai":
		return []string{
			"gpt-4",
			"gpt-4-turbo",
			"gpt-4-turbo-preview",
			"gpt-3.5-turbo",
			"gpt-3.5-turbo-16k",
		}, nil
	case "claude":
		return []string{
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
			"claude-2.1",
			"claude-2.0",
		}, nil
	case "gemini":
		return []string{
			"gemini-2.5-pro",
			"gemini-1.5-pro",
			"gemini-1.5-flash",
			"gemini-pro",
			"gemini-pro-vision",
			"gemini-ultra",
		}, nil
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", provider)
	}
}

// ===== 项目上下文AI对话服务 =====

// ProjectChatResponse 项目AI对话响应
type ProjectChatResponse struct {
	Message          string             `json:"message"`
	UpdatedAnalysis  *model.Requirement `json:"updated_analysis,omitempty"`
	Suggestions      []string           `json:"suggestions,omitempty"`
	RelatedQuestions []string           `json:"related_questions,omitempty"`
}

// ProjectChat 项目上下文AI对话 - 使用用户AI配置
func (s *AIService) ProjectChat(ctx context.Context, projectID uuid.UUID, message, context string, userID uuid.UUID) (*ProjectChatResponse, error) {
	// 验证项目是否存在
	project, err := s.mysqlRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	// 获取用户AI配置
	userConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("获取AI配置失败，请先在设置中配置AI服务: %w", err)
	}

	// 确定使用的provider
	provider := ai.AIProvider(userConfig.Provider)
	var apiKey string

	switch userConfig.Provider {
	case "openai":
		apiKey = userConfig.OpenAIAPIKey
	case "claude":
		apiKey = userConfig.ClaudeAPIKey
	case "gemini":
		apiKey = userConfig.GeminiAPIKey
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", userConfig.Provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未配置%s的API密钥，请先在设置中配置", userConfig.Provider)
	}

	// 创建临时AI客户端配置
	clientConfig := ai.AIManagerConfig{
		DefaultProvider: provider,
		EnableCache:     true,
		CacheTTL:        time.Hour,
	}

	// 根据provider设置对应的客户端配置
	switch provider {
	case ai.ProviderOpenAI:
		clientConfig.OpenAIConfig = &ai.OpenAIConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	case ai.ProviderClaude:
		// Claude客户端暂时不可用
		return nil, fmt.Errorf("Claude客户端暂时不可用")
	case ai.ProviderGemini:
		clientConfig.GeminiConfig = &ai.GeminiConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	}

	// 创建用户特定的AI管理器
	tempAIManager, err := ai.NewAIManager(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AI管理器失败: %w", err)
	}

	// 获取项目的需求分析数据作为上下文
	analyses, err := s.aiRepo.GetRequirementAnalysesByProject(projectID)
	if err != nil {
		log.Printf("获取项目需求分析失败: %v", err)
		// 继续执行，允许没有分析数据的对话
	}

	// 构建对话上下文
	contextData := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        project.ProjectName,
			"description": project.Description,
			"type":        project.ProjectType,
			"status":      project.Status,
		},
		"conversation_context": context,
	}

	// 如果有需求分析数据，添加到上下文中
	if len(analyses) > 0 {
		latestAnalysis := analyses[0]

		// 解析结构化需求
		var structuredReq map[string]interface{}
		if err := json.Unmarshal([]byte(latestAnalysis.StructuredRequirement), &structuredReq); err == nil {
			contextData["requirement_analysis"] = structuredReq
		}
		contextData["completeness_score"] = latestAnalysis.CompletenessScore
	}

	// 序列化上下文数据
	contextJSON, err := json.Marshal(contextData)
	if err != nil {
		return nil, fmt.Errorf("序列化上下文数据失败: %w", err)
	}

	// 调用AI进行对话 - 使用用户配置的AI提供商
	response, err := tempAIManager.ProjectChat(ctx, message, string(contextJSON), provider)
	if err != nil {
		return nil, fmt.Errorf("AI对话失败: %w", err)
	}

	// 构建响应
	chatResponse := &ProjectChatResponse{
		Message: response.Message,
	}

	// 如果AI建议更新需求分析，处理更新
	if response.ShouldUpdateAnalysis && len(analyses) > 0 {
		// 这里可以根据AI的建议更新需求分析
		// 暂时先返回现有的分析数据
		chatResponse.UpdatedAnalysis = analyses[0]
	}

	// 添加AI建议的相关问题
	if len(response.RelatedQuestions) > 0 {
		chatResponse.RelatedQuestions = response.RelatedQuestions
	}

	// 添加AI提供的建议
	if len(response.Suggestions) > 0 {
		chatResponse.Suggestions = response.Suggestions
	}

	return chatResponse, nil
}

// ===== 分阶段文档生成相关服务 =====

// GenerateStageDocuments 分阶段生成项目文档
func (s *AIService) GenerateStageDocuments(ctx context.Context, req *model.GenerateStageDocumentsRequest, userID uuid.UUID) (*model.StageDocumentsResult, error) {
	// 验证项目是否存在
	project, err := s.mysqlRepo.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	// 获取项目的需求分析
	analyses, err := s.aiRepo.GetRequirementAnalysesByProject(req.ProjectID)
	if err != nil || len(analyses) == 0 {
		return nil, fmt.Errorf("请先完成项目需求分析")
	}

	latestAnalysis := analyses[0]

	// 获取用户AI配置
	userConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("获取AI配置失败，请先在设置中配置AI服务: %w", err)
	}

	// 确定使用的provider
	provider := ai.AIProvider(userConfig.Provider)
	var apiKey string

	switch userConfig.Provider {
	case "openai":
		apiKey = userConfig.OpenAIAPIKey
	case "claude":
		apiKey = userConfig.ClaudeAPIKey
	case "gemini":
		apiKey = userConfig.GeminiAPIKey
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", userConfig.Provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未配置%s的API密钥，请先在设置中配置", userConfig.Provider)
	}

	// 创建临时AI客户端配置
	clientConfig := ai.AIManagerConfig{
		DefaultProvider: provider,
		EnableCache:     true,
		CacheTTL:        time.Hour,
	}

	// 根据provider设置对应的客户端配置
	switch provider {
	case ai.ProviderOpenAI:
		clientConfig.OpenAIConfig = &ai.OpenAIConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	case ai.ProviderClaude:
		// Claude客户端暂时不可用
		return nil, fmt.Errorf("Claude客户端暂时不可用")
	case ai.ProviderGemini:
		clientConfig.GeminiConfig = &ai.GeminiConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	}

	// 创建用户特定的AI管理器
	tempAIManager, err := ai.NewAIManager(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AI管理器失败: %w", err)
	}

	// 构建AI分析对象
	var structuredReq map[string]interface{}
	if err := json.Unmarshal([]byte(latestAnalysis.StructuredRequirement), &structuredReq); err != nil {
		return nil, fmt.Errorf("解析结构化需求失败: %w", err)
	}

	analysis := &ai.RequirementAnalysis{
		ID:           latestAnalysis.RequirementID.String(),
		ProjectID:    req.ProjectID.String(),
		OriginalText: latestAnalysis.RawRequirement,
	}

	// 提取核心功能
	if coreFuncs, ok := structuredReq["core_functions"].([]interface{}); ok {
		for _, fn := range coreFuncs {
			if funcStr, ok := fn.(string); ok {
				analysis.CoreFunctions = append(analysis.CoreFunctions, funcStr)
			}
		}
	}

	result := &model.StageDocumentsResult{
		ProjectID:    req.ProjectID,
		Stage:        req.Stage,
		GeneratedAt:  time.Now(),
		Documents:    make([]*model.Document, 0),
		PUMLDiagrams: make([]*model.PUMLDiagram, 0),
	}

	// 根据不同阶段生成对应的文档
	switch req.Stage {
	case 1:
		// 第一阶段：项目需求文档 + 系统架构图 + 业务流程图 + 数据模型图 + 交互流程图
		log.Printf("开始生成第一阶段文档 (项目: %s)...", project.ProjectName)
		// TODO: 使用tempAIManager实现第一阶段文档生成
		_ = tempAIManager // 临时使用变量避免编译错误

	case 2:
		// 第二阶段：技术规范文档 + API设计 + 数据库设计
		log.Printf("开始生成第二阶段文档 (项目: %s)...", project.ProjectName)
		// TODO: 使用tempAIManager实现第二阶段文档生成
		_ = tempAIManager // 临时使用变量避免编译错误

	case 3:
		// 第三阶段：开发流程文档 + 测试用例文档 + 部署文档
		log.Printf("开始生成第三阶段文档 (项目: %s)...", project.ProjectName)
		// TODO: 使用tempAIManager实现第三阶段文档生成
		_ = tempAIManager // 临时使用变量避免编译错误

	default:
		return nil, fmt.Errorf("无效的阶段编号，支持的阶段：1、2、3")
	}

	return result, nil
}

// GenerateSpecificStageDocuments 生成指定阶段的特定文档（用于一键生成完整项目文档）
func (s *AIService) GenerateSpecificStageDocuments(ctx context.Context, req *model.GenerateStageDocumentsRequest, userID uuid.UUID, documentNames []string) (*model.StageDocumentsResult, error) {
	// 注意documentNames参数将在将来的实现中使用
	log.Printf("生成指定文档: %v", documentNames)
	
	// 调用现有的GenerateStageDocuments方法
	return s.GenerateStageDocuments(ctx, req, userID)
}

// GeneratePUMLWithUser 使用用户配置生成PUML图表
func (s *AIService) GeneratePUMLWithUser(ctx context.Context, req *model.GeneratePUMLRequest, userID uuid.UUID) (*model.PUMLDiagram, error) {
	// 获取用户AI配置
	userConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("获取AI配置失败，请先在设置中配置AI服务: %w", err)
	}

	// 确定使用的provider
	provider := ai.AIProvider(userConfig.Provider)
	var apiKey string

	switch userConfig.Provider {
	case "openai":
		apiKey = userConfig.OpenAIAPIKey
	case "claude":
		apiKey = userConfig.ClaudeAPIKey
	case "gemini":
		apiKey = userConfig.GeminiAPIKey
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", userConfig.Provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未配置%s的API密钥，请先在设置中配置", userConfig.Provider)
	}

	// 创建临时AI客户端配置
	clientConfig := ai.AIManagerConfig{
		DefaultProvider: provider,
		EnableCache:     true,
		CacheTTL:        time.Hour,
	}

	// 根据provider设置对应的客户端配置
	switch provider {
	case ai.ProviderOpenAI:
		clientConfig.OpenAIConfig = &ai.OpenAIConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	case ai.ProviderClaude:
		// Claude客户端暂时不可用
		return nil, fmt.Errorf("Claude客户端暂时不可用")
	case ai.ProviderGemini:
		clientConfig.GeminiConfig = &ai.GeminiConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	}

	// 创建用户特定的AI管理器
	tempAIManager, err := ai.NewAIManager(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AI管理器失败: %w", err)
	}

	// 获取需求分析数据
	analysisID, err := uuid.Parse(req.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("无效的分析ID: %w", err)
	}

	analysis, err := s.aiRepo.GetRequirementAnalysis(analysisID)
	if err != nil {
		return nil, fmt.Errorf("获取需求分析失败: %w", err)
	}

	// 构建AI分析对象
	var structuredReq map[string]interface{}
	if err := json.Unmarshal([]byte(analysis.StructuredRequirement), &structuredReq); err != nil {
		return nil, fmt.Errorf("解析结构化需求失败: %w", err)
	}

	aiAnalysis := &ai.RequirementAnalysis{
		ID:           analysis.RequirementID.String(),
		ProjectID:    analysis.ProjectID.String(),
		OriginalText: analysis.RawRequirement,
	}

	// 提取核心功能
	if coreFuncs, ok := structuredReq["core_functions"].([]interface{}); ok {
		for _, fn := range coreFuncs {
			if funcStr, ok := fn.(string); ok {
				aiAnalysis.CoreFunctions = append(aiAnalysis.CoreFunctions, funcStr)
			}
		}
	}

	// 使用用户配置的AI管理器生成PUML
	pumlDiagram, err := tempAIManager.GeneratePUML(ctx, aiAnalysis, ai.PUMLType(req.DiagramType), provider)
	if err != nil {
		return nil, fmt.Errorf("AI生成PUML失败: %w", err)
	}

	// 创建数据库记录
	diagram := &model.PUMLDiagram{
		DiagramID:   uuid.New(),
		ProjectID:   analysis.ProjectID,
		DiagramType: req.DiagramType,
		DiagramName: fmt.Sprintf("%s图表", req.DiagramType),
		PUMLContent: pumlDiagram.Content,
		Version:     1,
		Stage:       1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	if err := s.aiRepo.CreatePUMLDiagram(diagram); err != nil {
		return nil, fmt.Errorf("保存PUML图表失败: %w", err)
	}

	return diagram, nil
}

// GenerateDocumentWithUser 使用用户配置生成开发文档
func (s *AIService) GenerateDocumentWithUser(ctx context.Context, req *model.GenerateDocumentRequest, userID uuid.UUID) (*model.Document, error) {
	// 获取用户AI配置
	userConfig, err := s.aiRepo.GetUserAIConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("获取AI配置失败，请先在设置中配置AI服务: %w", err)
	}

	// 确定使用的provider
	provider := ai.AIProvider(userConfig.Provider)
	var apiKey string

	switch userConfig.Provider {
	case "openai":
		apiKey = userConfig.OpenAIAPIKey
	case "claude":
		apiKey = userConfig.ClaudeAPIKey
	case "gemini":
		apiKey = userConfig.GeminiAPIKey
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", userConfig.Provider)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("未配置%s的API密钥，请先在设置中配置", userConfig.Provider)
	}

	// 创建临时AI客户端配置
	clientConfig := ai.AIManagerConfig{
		DefaultProvider: provider,
		EnableCache:     true,
		CacheTTL:        time.Hour,
	}

	// 根据provider设置对应的客户端配置
	switch provider {
	case ai.ProviderOpenAI:
		clientConfig.OpenAIConfig = &ai.OpenAIConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	case ai.ProviderClaude:
		// Claude客户端暂时不可用
		return nil, fmt.Errorf("Claude客户端暂时不可用")
	case ai.ProviderGemini:
		clientConfig.GeminiConfig = &ai.GeminiConfig{
			APIKey: apiKey,
			Model:  userConfig.DefaultModel,
		}
	}

	// 创建用户特定的AI管理器
	tempAIManager, err := ai.NewAIManager(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AI管理器失败: %w", err)
	}

	// 获取需求分析数据
	analysisID, err := uuid.Parse(req.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("无效的分析ID: %w", err)
	}

	analysis, err := s.aiRepo.GetRequirementAnalysis(analysisID)
	if err != nil {
		return nil, fmt.Errorf("获取需求分析失败: %w", err)
	}

	// 构建AI分析对象
	var structuredReq map[string]interface{}
	if err := json.Unmarshal([]byte(analysis.StructuredRequirement), &structuredReq); err != nil {
		return nil, fmt.Errorf("解析结构化需求失败: %w", err)
	}

	aiAnalysis := &ai.RequirementAnalysis{
		ID:           analysis.RequirementID.String(),
		ProjectID:    analysis.ProjectID.String(),
		OriginalText: analysis.RawRequirement,
	}

	// 提取核心功能
	if coreFuncs, ok := structuredReq["core_functions"].([]interface{}); ok {
		for _, fn := range coreFuncs {
			if funcStr, ok := fn.(string); ok {
				aiAnalysis.CoreFunctions = append(aiAnalysis.CoreFunctions, funcStr)
			}
		}
	}

	// 使用用户配置的AI管理器生成文档
	aiDocument, err := tempAIManager.GenerateDocument(ctx, aiAnalysis, provider)
	if err != nil {
		return nil, fmt.Errorf("AI生成文档失败: %w", err)
	}

	// 将AI文档转换为JSON字符串
	documentJSON, err := json.Marshal(aiDocument)
	if err != nil {
		return nil, fmt.Errorf("序列化文档内容失败: %w", err)
	}

	// 创建数据库记录
	document := &model.Document{
		DocumentID:   uuid.New(),
		ProjectID:    analysis.ProjectID,
		DocumentType: "development",
		DocumentName: "开发文档",
		Content:      string(documentJSON),
		Format:       "json",
		Version:      1,
		Stage:        1,
		GeneratedAt:  time.Now(),
	}

	// 保存到数据库
	if err := s.aiRepo.CreateDocument(document); err != nil {
		return nil, fmt.Errorf("保存文档失败: %w", err)
	}

	return document, nil
}
