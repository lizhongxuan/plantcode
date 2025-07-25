package controller

import (
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AIController AI功能控制器
type AIController struct {
	aiService *service.AIService
}

// NewAIController 创建AI控制器
func NewAIController(aiService *service.AIService) *AIController {
	return &AIController{
		aiService: aiService,
	}
}

// AnalyzeRequirement 分析业务需求
func (ac *AIController) AnalyzeRequirement(c *gin.Context) {
	log.InfofId(c, "AnalyzeRequirement: 开始处理需求分析请求")

	// 从上下文获取用户（通过JWT认证中间件设置）
	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "AnalyzeRequirement: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.AIAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "AnalyzeRequirement: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Requirement == "" {
		log.WarnfId(c, "AnalyzeRequirement: 需求描述不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "需求描述不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if len(req.Requirement) > 10000 {
		log.WarnfId(c, "AnalyzeRequirement: 需求描述过长")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "需求描述不能超过10000个字符",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "AnalyzeRequirement: 用户 %s 请求分析需求", user.UserID.String())

	// 调用AI服务进行需求分析
	result, err := ac.aiService.AnalyzeRequirementWithUser(c.Request.Context(), &req, user.UserID)
	if err != nil {
		log.ErrorfId(c, "AnalyzeRequirement: 需求分析失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "AnalyzeRequirement: 需求分析成功，分析ID: %s", result.RequirementID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "需求分析成功",
		"code":    http.StatusOK,
	})
}

// GetRequirementAnalysis 获取需求分析结果
func (ac *AIController) GetRequirementAnalysis(c *gin.Context) {
	log.InfofId(c, "GetRequirementAnalysis: 开始获取需求分析结果")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetRequirementAnalysis: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	analysisID := c.Param("id")
	if analysisID == "" {
		log.WarnfId(c, "GetRequirementAnalysis: 分析ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "分析ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 验证UUID格式
	analysisUUID, err := uuid.Parse(analysisID)
	if err != nil {
		log.WarnfId(c, "GetRequirementAnalysis: 无效的分析ID格式: %s", analysisID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的分析ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetRequirementAnalysis: 获取分析结果，分析ID: %s, 用户ID: %s", analysisID, user.UserID.String())

	// 调用服务获取分析结果
	result, err := ac.aiService.GetRequirementAnalysis(analysisUUID)
	if err != nil {
		log.ErrorfId(c, "GetRequirementAnalysis: 获取分析结果失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetRequirementAnalysis: 成功获取分析结果")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "获取分析结果成功",
		"code":    http.StatusOK,
	})
}

// GetRequirementAnalysesByProject 获取项目的需求分析列表
func (ac *AIController) GetRequirementAnalysesByProject(c *gin.Context) {
	log.InfofId(c, "GetRequirementAnalysesByProject: 开始获取项目需求分析列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetRequirementAnalysesByProject: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		log.WarnfId(c, "GetRequirementAnalysesByProject: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 验证UUID格式
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		log.WarnfId(c, "GetRequirementAnalysesByProject: 无效的项目ID格式: %s", projectID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetRequirementAnalysesByProject: 获取项目需求分析，项目ID: %s, 用户ID: %s", projectID, user.UserID.String())

	// 调用服务获取项目需求分析列表
	analyses, err := ac.aiService.GetRequirementAnalysesByProject(projectUUID)
	if err != nil {
		log.ErrorfId(c, "GetRequirementAnalysesByProject: 获取项目需求分析失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetRequirementAnalysesByProject: 成功获取项目需求分析列表，数量: %d", len(analyses))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analyses,
		"message": "获取项目需求分析成功",
		"code":    http.StatusOK,
	})
}

// GeneratePUML 生成PUML图表
func (ac *AIController) GeneratePUML(c *gin.Context) {
	log.InfofId(c, "GeneratePUML: 开始处理PUML生成请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GeneratePUML: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.GeneratePUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "GeneratePUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.AnalysisID == "" {
		log.WarnfId(c, "GeneratePUML: 分析ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "分析ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GeneratePUML: 用户 %s 请求生成PUML图表", user.UserID.String())

	// 调用AI服务生成PUML
	result, err := ac.aiService.GeneratePUMLWithUser(c.Request.Context(), &req, user.UserID)
	if err != nil {
		log.ErrorfId(c, "GeneratePUML: PUML生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GeneratePUML: PUML生成成功，图表ID: %s", result.DiagramID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML生成成功",
		"code":    http.StatusOK,
	})
}

// GetPUMLDiagramsByProjectID 获取项目的PUML图表列表
func (ac *AIController) GetPUMLDiagramsByProjectID(c *gin.Context) {
	log.InfofId(c, "GetPUMLDiagramsByProjectID: 开始获取项目PUML图表列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetPUMLDiagramsByProjectID: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		log.WarnfId(c, "GetPUMLDiagramsByProjectID: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetPUMLDiagramsByProjectID: 获取项目PUML图表，项目ID: %s, 用户ID: %s", projectID, user.UserID.String())

	// 调用服务获取项目PUML图表列表
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		log.WarnfId(c, "GetPUMLDiagramsByProjectID: 无效的项目ID格式: %s", projectID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	diagrams, err := ac.aiService.GetPUMLDiagramsByProjectID(projectUUID)
	if err != nil {
		log.ErrorfId(c, "GetPUMLDiagramsByProjectID: 获取项目PUML图表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetPUMLDiagramsByProjectID: 成功获取项目PUML图表列表，数量: %d", len(diagrams))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    diagrams,
		"message": "获取项目PUML图表成功",
		"code":    http.StatusOK,
	})
}

// UpdatePUML 更新PUML图表
func (ac *AIController) UpdatePUML(c *gin.Context) {
	log.InfofId(c, "UpdatePUML: 开始处理PUML更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdatePUML: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	diagramID := c.Param("id")
	if diagramID == "" {
		log.WarnfId(c, "UpdatePUML: 图表ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "图表ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	var req model.UpdatePUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "UpdatePUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdatePUML: 用户 %s 请求更新PUML图表: %s", user.UserID.String(), diagramID)

	// 调用服务更新PUML图表
	diagramUUID, err := uuid.Parse(diagramID)
	if err != nil {
		log.WarnfId(c, "UpdatePUML: 无效的图表ID格式: %s", diagramID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的图表ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	err = ac.aiService.UpdatePUMLDiagram(diagramUUID, &req)
	if err != nil {
		log.ErrorfId(c, "UpdatePUML: PUML更新失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "UpdatePUML: PUML更新成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PUML更新成功",
		"code":    http.StatusOK,
	})
}

// GenerateDocument 生成技术文档
func (ac *AIController) GenerateDocument(c *gin.Context) {
	log.InfofId(c, "GenerateDocument: 开始处理文档生成请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GenerateDocument: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.GenerateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "GenerateDocument: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.AnalysisID == "" {
		log.WarnfId(c, "GenerateDocument: 分析ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "分析ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GenerateDocument: 用户 %s 请求生成技术文档", user.UserID.String())

	// 调用AI服务生成文档
	result, err := ac.aiService.GenerateDocumentWithUser(c.Request.Context(), &req, user.UserID)
	if err != nil {
		log.ErrorfId(c, "GenerateDocument: 文档生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GenerateDocument: 文档生成成功，文档ID: %s", result.DocumentID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "文档生成成功",
		"code":    http.StatusOK,
	})
}

// GetDocumentsByProjectID 获取项目的技术文档列表
func (ac *AIController) GetDocumentsByProjectID(c *gin.Context) {
	log.InfofId(c, "GetDocumentsByProjectID: 开始获取项目技术文档列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetDocumentsByProjectID: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		log.WarnfId(c, "GetDocumentsByProjectID: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetDocumentsByProjectID: 获取项目技术文档，项目ID: %s, 用户ID: %s", projectID, user.UserID.String())

	// 调用服务获取项目技术文档列表
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		log.WarnfId(c, "GetDocumentsByProjectID: 无效的项目ID格式: %s", projectID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	documents, err := ac.aiService.GetDocumentsByProjectID(projectUUID)
	if err != nil {
		log.ErrorfId(c, "GetDocumentsByProjectID: 获取项目技术文档失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetDocumentsByProjectID: 成功获取项目技术文档列表，数量: %d", len(documents))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    documents,
		"message": "获取项目技术文档成功",
		"code":    http.StatusOK,
	})
}

// UpdateDocument 更新技术文档
func (ac *AIController) UpdateDocument(c *gin.Context) {
	log.InfofId(c, "UpdateDocument: 开始处理文档更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdateDocument: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		log.WarnfId(c, "UpdateDocument: 文档ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文档ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	var req model.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "UpdateDocument: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdateDocument: 用户 %s 请求更新技术文档: %s", user.UserID.String(), documentID)

	// 调用服务更新技术文档
	documentUUID, err := uuid.Parse(documentID)
	if err != nil {
		log.WarnfId(c, "UpdateDocument: 无效的文档ID格式: %s", documentID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的文档ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	err = ac.aiService.UpdateDocument(documentUUID, &req)
	if err != nil {
		log.ErrorfId(c, "UpdateDocument: 文档更新失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "UpdateDocument: 文档更新成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文档更新成功",
		"code":    http.StatusOK,
	})
}

// CreateChatSession 创建聊天会话 - 暂时不实现
func (ac *AIController) CreateChatSession(c *gin.Context) {
	log.InfofId(c, "CreateChatSession: 聊天功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "聊天功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// SendChatMessage 发送聊天消息 - 暂时不实现
func (ac *AIController) SendChatMessage(c *gin.Context) {
	log.InfofId(c, "SendChatMessage: 聊天功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "聊天功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GetChatMessages 获取聊天消息 - 暂时不实现
func (ac *AIController) GetChatMessages(c *gin.Context) {
	log.InfofId(c, "GetChatMessages: 聊天功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "聊天功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GetAIProviders 获取AI提供商列表
func (ac *AIController) GetAIProviders(c *gin.Context) {
	log.InfofId(c, "GetAIProviders: 开始获取AI提供商列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetAIProviders: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "GetAIProviders: 用户 %s 请求获取AI提供商列表", user.UserID.String())

	// 调用服务获取AI提供商列表
	providers := ac.aiService.GetAIProviders()

	log.InfofId(c, "GetAIProviders: 成功获取AI提供商列表")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    providers,
		"message": "获取AI提供商列表成功",
		"code":    http.StatusOK,
	})
}

// GetUserAIConfig 获取用户AI配置
func (ac *AIController) GetUserAIConfig(c *gin.Context) {
	log.InfofId(c, "GetUserAIConfig: 开始获取用户AI配置")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetUserAIConfig: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "GetUserAIConfig: 用户 %s 请求获取AI配置", user.UserID.String())

	// 调用服务获取用户AI配置
	config, err := ac.aiService.GetUserAIConfig(user.UserID)
	if err != nil {
		log.ErrorfId(c, "GetUserAIConfig: 获取用户AI配置失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetUserAIConfig: 成功获取用户AI配置")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
		"message": "获取用户AI配置成功",
		"code":    http.StatusOK,
	})
}

// UpdateUserAIConfig 更新用户AI配置
func (ac *AIController) UpdateUserAIConfig(c *gin.Context) {
	log.InfofId(c, "UpdateUserAIConfig: 开始处理用户AI配置更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdateUserAIConfig: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.UpdateUserAIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "UpdateUserAIConfig: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdateUserAIConfig: 用户 %s 请求更新AI配置", user.UserID.String())

	// 调用服务更新用户AI配置
	config, err := ac.aiService.UpdateUserAIConfig(user.UserID, &req)
	if err != nil {
		log.ErrorfId(c, "UpdateUserAIConfig: 更新用户AI配置失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "UpdateUserAIConfig: 用户AI配置更新成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
		"message": "用户AI配置更新成功",
		"code":    http.StatusOK,
	})
}

// TestAIConnection 测试AI连接
func (ac *AIController) TestAIConnection(c *gin.Context) {
	log.InfofId(c, "TestAIConnection: 开始处理AI连接测试请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "TestAIConnection: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.TestAIConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "TestAIConnection: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "TestAIConnection: 用户 %s 请求测试AI连接", user.UserID.String())

	// 调用服务测试AI连接
	result, err := ac.aiService.TestAIConnection(&req)
	if err != nil {
		log.ErrorfId(c, "TestAIConnection: AI连接测试失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "TestAIConnection: AI连接测试成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "AI连接测试成功",
		"code":    http.StatusOK,
	})
}

// GetAvailableModels 获取可用模型列表
func (ac *AIController) GetAvailableModels(c *gin.Context) {
	log.InfofId(c, "GetAvailableModels: 开始获取可用模型列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetAvailableModels: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		log.WarnfId(c, "GetAvailableModels: 提供商不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "提供商不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetAvailableModels: 用户 %s 请求获取 %s 的可用模型列表", user.UserID.String(), provider)

	// 调用服务获取可用模型列表
	models, err := ac.aiService.GetAvailableModels(provider)
	if err != nil {
		log.ErrorfId(c, "GetAvailableModels: 获取可用模型列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetAvailableModels: 成功获取可用模型列表")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    models,
		"message": "获取可用模型列表成功",
		"code":    http.StatusOK,
	})
}

// ProjectChat 项目上下文AI对话 - 暂时不实现
func (ac *AIController) ProjectChat(c *gin.Context) {
	log.InfofId(c, "ProjectChat: 项目对话功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "项目对话功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GenerateStageDocuments 生成阶段文档 - 暂时不实现
func (ac *AIController) GenerateStageDocuments(c *gin.Context) {
	log.InfofId(c, "GenerateStageDocuments: 阶段文档生成功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "阶段文档生成功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GenerateStageDocumentList 生成阶段文档列表 - 暂时不实现
func (ac *AIController) GenerateStageDocumentList(c *gin.Context) {
	log.InfofId(c, "GenerateStageDocumentList: 阶段文档列表生成功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "阶段文档列表生成功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}
