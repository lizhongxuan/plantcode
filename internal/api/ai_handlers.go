package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
)

// AIHandlers AI相关的HTTP处理器
type AIHandlers struct {
	aiService *service.AIService
}

// NewAIHandlers 创建AI处理器
func NewAIHandlers(aiService *service.AIService) *AIHandlers {
	return &AIHandlers{
		aiService: aiService,
	}
}

// ===== 需求分析相关接口 =====

// AnalyzeRequirement 分析业务需求
// POST /api/ai/analyze
func (h *AIHandlers) AnalyzeRequirement(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户（通过JWT认证中间件设置）
	user := MustGetUserFromContext(r.Context())

	var req model.AIAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 简单验证
	if req.Requirement == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "需求描述不能为空")
		return
	}

	// 调用服务层进行分析，传递用户ID
	analysis, err := h.aiService.AnalyzeRequirementWithUser(r.Context(), &req, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("需求分析失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, analysis, "需求分析完成")
}

// GetRequirementAnalysis 获取需求分析详情
// GET /api/ai/analysis/{analysisId}
func (h *AIHandlers) GetRequirementAnalysis(w http.ResponseWriter, r *http.Request) {
	analysisIDStr := r.PathValue("analysisId")
	if analysisIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少分析ID")
		return
	}

	analysisID, err := uuid.Parse(analysisIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的分析ID")
		return
	}

	analysis, err := h.aiService.GetRequirementAnalysis(analysisID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "需求分析不存在")
		return
	}

	utils.WriteSuccessResponse(w, analysis, "获取需求分析成功")
}

// GetRequirementAnalysesByProject 获取项目的需求分析列表
// GET /api/ai/analysis/project/{projectId}
func (h *AIHandlers) GetRequirementAnalysesByProject(w http.ResponseWriter, r *http.Request) {
	// 从URL路径获取项目ID
	projectIDStr := extractIDFromPath(r.URL.Path, "/api/ai/analysis/project/")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	analyses, err := h.aiService.GetRequirementAnalysesByProject(projectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "获取需求分析列表失败")
		return
	}

	utils.WriteSuccessResponse(w, analyses, "获取需求分析列表成功")
}

// ===== PUML图表相关接口 =====

// GeneratePUML 生成PUML图表
// POST /api/ai/puml/generate
func (h *AIHandlers) GeneratePUML(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户信息
	user := MustGetUserFromContext(r.Context())

	var req model.GeneratePUMLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 简单验证
	if req.AnalysisID == "" || req.DiagramType == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "分析ID和图表类型不能为空")
		return
	}

	// 使用用户配置生成PUML图表
	diagram, err := h.aiService.GeneratePUMLWithUser(r.Context(), &req, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("生成PUML图表失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, diagram, "PUML图表生成成功")
}

// GetPUMLDiagramsByProject 获取项目的PUML图表列表
// GET /api/ai/puml/project/{projectId}
func (h *AIHandlers) GetPUMLDiagramsByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectId")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	diagrams, err := h.aiService.GetPUMLDiagramsByProject(projectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "获取PUML图表列表失败")
		return
	}

	utils.WriteSuccessResponse(w, diagrams, "获取PUML图表列表成功")
}

// ===== 文档生成相关接口 =====

// GenerateDocument 生成开发文档
// POST /api/ai/document/generate
func (h *AIHandlers) GenerateDocument(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户信息
	user := MustGetUserFromContext(r.Context())

	var req model.GenerateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 简单验证
	if req.AnalysisID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "分析ID不能为空")
		return
	}

	// 使用用户配置生成文档
	document, err := h.aiService.GenerateDocumentWithUser(r.Context(), &req, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("生成开发文档失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, document, "开发文档生成成功")
}

// GetDocumentsByProject 获取项目的文档列表
// GET /api/ai/document/project/{projectId}
func (h *AIHandlers) GetDocumentsByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectId")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	documents, err := h.aiService.GetDocumentsByProject(projectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "获取文档列表失败")
		return
	}

	utils.WriteSuccessResponse(w, documents, "获取文档列表成功")
}

// ===== 对话相关接口 =====

// CreateChatSession 创建对话会话
// POST /api/ai/chat/session
func (h *AIHandlers) CreateChatSession(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户（通过JWT认证中间件设置）
	user := MustGetUserFromContext(r.Context())

	var req model.ChatSessionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证标题不能为空，project_id可以为空（用于通用聊天）
	if req.Title == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "标题不能为空")
		return
	}

	session, err := h.aiService.CreateChatSession(&req, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("创建对话会话失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, session, "对话会话创建成功")
}

// SendChatMessage 发送聊天消息
// POST /api/ai/chat/message
func (h *AIHandlers) SendChatMessage(w http.ResponseWriter, r *http.Request) {
	var req model.SendChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 简单验证
	if req.SessionID == "" || req.Content == "" || req.Role == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "会话ID、内容和角色不能为空")
		return
	}

	message, err := h.aiService.SendChatMessage(r.Context(), &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, message, "消息发送成功")
}

// GetChatMessages 获取对话消息列表
// GET /api/ai/chat/session/{sessionId}/messages
func (h *AIHandlers) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("sessionId")
	if sessionIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少会话ID")
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的会话ID")
		return
	}

	messages, err := h.aiService.GetChatMessages(sessionID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "获取对话消息失败")
		return
	}

	utils.WriteSuccessResponse(w, messages, "获取对话消息成功")
}

// ===== 系统管理相关接口 =====

// GetAIProviders 获取可用的AI提供商列表
// GET /api/ai/providers
func (h *AIHandlers) GetAIProviders(w http.ResponseWriter, r *http.Request) {
	providers := h.aiService.GetAIProviders()
	defaultProvider := h.aiService.GetDefaultAIProvider()

	response := map[string]interface{}{
		"providers":        providers,
		"default_provider": defaultProvider,
	}

	utils.WriteSuccessResponse(w, response, "获取AI提供商列表成功")
}

// GetCacheStats 获取缓存统计信息
// GET /api/ai/cache/stats
func (h *AIHandlers) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.aiService.GetCacheStats()
	utils.WriteSuccessResponse(w, stats, "获取缓存统计成功")
}

// ClearCache 清空AI缓存
// DELETE /api/ai/cache
func (h *AIHandlers) ClearCache(w http.ResponseWriter, r *http.Request) {
	h.aiService.ClearCache()
	utils.WriteSuccessResponse(w, map[string]string{"result": "success"}, "缓存清空成功")
}

// UpdatePUML 更新PUML图表
// PUT /api/ai/puml/{diagramId}
func (h *AIHandlers) UpdatePUML(w http.ResponseWriter, r *http.Request) {
	// 从URL路径提取diagram ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少图表ID")
		return
	}
	
	diagramIDStr := parts[4] // /api/ai/puml/{diagramId}
	if diagramIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少图表ID")
		return
	}

	diagramID, err := uuid.Parse(diagramIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的图表ID")
		return
	}

	var req model.UpdatePUMLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.Content == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "PUML内容不能为空")
		return
	}

	err = h.aiService.UpdatePUMLDiagram(diagramID, &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("更新PUML图表失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, map[string]string{"result": "success"}, "PUML图表更新成功")
}

// UpdateDocument 更新技术文档
// PUT /api/ai/document/{documentId}
func (h *AIHandlers) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	// 从URL路径提取document ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少文档ID")
		return
	}
	
	documentIDStr := parts[4] // /api/ai/document/{documentId}
	if documentIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少文档ID")
		return
	}

	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的文档ID")
		return
	}

	var req model.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.Content == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "文档内容不能为空")
		return
	}

	err = h.aiService.UpdateDocument(documentID, &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("更新技术文档失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, map[string]string{"result": "success"}, "技术文档更新成功")
}

// ===== 用户AI配置相关接口 =====

// GetUserAIConfig 获取用户AI配置
// GET /api/ai/config
func (h *AIHandlers) GetUserAIConfig(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户
	user := MustGetUserFromContext(r.Context())

	// 获取用户AI配置
	config, err := h.aiService.GetUserAIConfig(user.UserID)
	if err != nil {
		// 如果没有配置，返回默认配置
		config = &model.UserAIConfig{
			Provider:     "openai",
			DefaultModel: "gpt-4",
			MaxTokens:    2048,
			IsActive:     false,
		}
	}

	// 构建返回结果，包含API密钥的配置状态
	response := map[string]interface{}{
		"provider":      config.Provider,
		"default_model": config.DefaultModel,
		"max_tokens":    config.MaxTokens,
		"is_active":     config.IsActive,
		"created_at":    config.CreatedAt,
		"updated_at":    config.UpdatedAt,
		// API密钥状态（不返回实际密钥，只返回是否已配置）
		"api_keys_configured": map[string]bool{
			"openai": config.OpenAIAPIKey != "",
			"claude": config.ClaudeAPIKey != "",
			"gemini": config.GeminiAPIKey != "",
		},
		// 脱敏显示的API密钥（显示前4个字符和后4个字符）
		"api_keys_display": map[string]string{
			"openai": h.maskAPIKey(config.OpenAIAPIKey),
			"claude": h.maskAPIKey(config.ClaudeAPIKey),
			"gemini": h.maskAPIKey(config.GeminiAPIKey),
		},
	}

	utils.WriteSuccessResponse(w, response, "获取AI配置成功")
}

// maskAPIKey 脱敏显示API密钥
func (h *AIHandlers) maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return ""
	}
	
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	
	// 显示前4个字符 + 星号 + 后4个字符
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
}

// UpdateUserAIConfig 更新用户AI配置
// PUT /api/ai/config
func (h *AIHandlers) UpdateUserAIConfig(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())
	
	var req model.UpdateUserAIConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证提供商
	if req.Provider != "openai" && req.Provider != "claude" && req.Provider != "gemini" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的AI提供商")
		return
	}

	// 验证至少有一个API密钥
	if req.OpenAIAPIKey == "" && req.ClaudeAPIKey == "" && req.GeminiAPIKey == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "至少需要配置一个AI提供商的API密钥")
		return
	}

	config, err := h.aiService.UpdateUserAIConfig(user.UserID, &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("更新AI配置失败: %v", err))
		return
	}

	// 隐藏API密钥信息
	config.OpenAIAPIKey = ""
	config.ClaudeAPIKey = ""
	config.GeminiAPIKey = ""
	
	utils.WriteSuccessResponse(w, config, "AI配置更新成功")
}

// TestAIConnection 测试AI连接
// POST /api/ai/test-connection
func (h *AIHandlers) TestAIConnection(w http.ResponseWriter, r *http.Request) {
	var req model.TestAIConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.Provider == "" || req.APIKey == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "提供商和API密钥不能为空")
		return
	}

	result, err := h.aiService.TestAIConnection(&req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("测试连接失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, result, "连接测试完成")
}

// GetAvailableModels 获取可用的AI模型列表
// GET /api/ai/models/{provider}
func (h *AIHandlers) GetAvailableModels(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少提供商参数")
		return
	}

	models, err := h.aiService.GetAvailableModels(provider)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取模型列表失败: %v", err))
		return
	}

	response := map[string]interface{}{
		"provider": provider,
		"models":   models,
	}

	utils.WriteSuccessResponse(w, response, "获取模型列表成功")
}

// ===== 项目上下文AI对话接口 =====

// ProjectChat 项目上下文AI对话
// POST /api/ai/chat
func (h *AIHandlers) ProjectChat(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户信息
	user := MustGetUserFromContext(r.Context())

	var req struct {
		ProjectID string `json:"project_id" validate:"required"`
		Message   string `json:"message" validate:"required"`
		Context   string `json:"context"` // requirement_analysis, puml_editing, document_review
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.ProjectID == "" || req.Message == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "项目ID和消息内容不能为空")
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 设置默认上下文
	if req.Context == "" {
		req.Context = "requirement_analysis"
	}

	// 调用服务方法，传递用户ID
	response, err := h.aiService.ProjectChat(r.Context(), projectID, req.Message, req.Context, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("AI对话失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, response, "AI对话成功")
}

// GenerateStageDocuments 分阶段生成文档
// POST /api/ai/generate-stage-documents
func (h *AIHandlers) GenerateStageDocuments(w http.ResponseWriter, r *http.Request) {
	var req model.GenerateStageDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求参数
	if req.ProjectID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "项目ID不能为空")
		return
	}

	if req.Stage < 1 || req.Stage > 3 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "阶段必须是1、2或3")
		return
	}

	// 从上下文获取用户（通过JWT认证中间件设置）
	user := MustGetUserFromContext(r.Context())

	// 调用服务生成文档
	result, err := h.aiService.GenerateStageDocuments(r.Context(), &req, user.UserID)
	if err != nil {
		log.Printf("生成阶段文档失败: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "生成文档失败: "+err.Error())
		return
	}

	utils.WriteSuccessResponse(w, result, "分阶段文档生成成功")
}

// GenerateStageDocumentList 根据需求分析生成阶段文档列表
// POST /api/ai/generate-document-list
func (h *AIHandlers) GenerateStageDocumentList(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	var req struct {
		ProjectID string `json:"project_id" validate:"required"`
		Stage     int    `json:"stage" validate:"required,min=1,max=3"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.ProjectID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "项目ID不能为空")
		return
	}

	if req.Stage < 1 || req.Stage > 3 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "阶段必须是1、2或3")
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 创建阶段文档生成请求
	stageReq := &model.GenerateStageDocumentsRequest{
		ProjectID: projectID,
		Stage:     req.Stage,
	}

	// 调用AI服务生成阶段文档
	result, err := h.aiService.GenerateStageDocuments(r.Context(), stageReq, user.UserID)
	if err != nil {
		log.Printf("生成阶段文档列表失败: %v", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "生成文档列表失败: "+err.Error())
		return
	}

	utils.WriteSuccessResponse(w, result, "阶段文档列表生成成功")
} 