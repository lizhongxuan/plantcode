package controller

import (
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PUMLController PUML功能控制器
type PUMLController struct {
	pumlService *service.PUMLService
	aiService   *service.AIService
}

// NewPUMLController 创建PUML控制器
func NewPUMLController(pumlService *service.PUMLService, aiService *service.AIService) *PUMLController {
	return &PUMLController{
		pumlService: pumlService,
		aiService:   aiService,
	}
}

// CreatePUML 创建PUML图表
func (pc *PUMLController) CreatePUML(c *gin.Context) {
	log.InfofId(c, "CreatePUML: 开始处理PUML创建请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "CreatePUML: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.CreatePUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "CreatePUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.ProjectID == "" {
		log.WarnfId(c, "CreatePUML: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if req.Title == "" {
		log.WarnfId(c, "CreatePUML: 图表标题不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "图表标题不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if req.Content == "" {
		log.WarnfId(c, "CreatePUML: PUML内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "CreatePUML: 用户 %s 请求创建PUML图表", user.UserID.String())

	// 调用服务创建PUML图表
	result, err := pc.pumlService.CreatePUML(user.UserID, &req)
	if err != nil {
		log.ErrorfId(c, "CreatePUML: PUML创建失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "CreatePUML: PUML创建成功，图表ID: %s", result.DiagramID.String())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML创建成功",
		"code":    http.StatusCreated,
	})
}

// GetProjectPUMLs 获取项目PUML图表列表
func (pc *PUMLController) GetProjectPUMLs(c *gin.Context) {
	log.InfofId(c, "GetProjectPUMLs: 开始获取项目PUML图表列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetProjectPUMLs: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		log.WarnfId(c, "GetProjectPUMLs: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetProjectPUMLs: 获取项目PUML图表，项目ID: %s, 用户ID: %s", projectID, user.UserID.String())

	// 调用服务获取项目PUML图表列表
	pumls, err := pc.pumlService.GetProjectPUMLs(user.UserID, projectID)
	if err != nil {
		log.ErrorfId(c, "GetProjectPUMLs: 获取项目PUML图表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetProjectPUMLs: 成功获取项目PUML图表列表，数量: %d", len(pumls))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pumls,
		"message": "获取项目PUML图表成功",
		"code":    http.StatusOK,
	})
}

// UpdatePUMLDiagram 更新PUML图表
func (pc *PUMLController) UpdatePUMLDiagram(c *gin.Context) {
	log.InfofId(c, "UpdatePUMLDiagram: 开始处理PUML图表更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdatePUMLDiagram: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	pumlID := c.Param("pumlId")
	if pumlID == "" {
		log.WarnfId(c, "UpdatePUMLDiagram: PUML图表ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML图表ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	var req model.UpdatePUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "UpdatePUMLDiagram: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdatePUMLDiagram: 用户 %s 请求更新PUML图表: %s", user.UserID.String(), pumlID)

	// 调用服务更新PUML图表
	result, err := pc.pumlService.UpdatePUMLDiagram(user.UserID, pumlID, &req)
	if err != nil {
		log.ErrorfId(c, "UpdatePUMLDiagram: PUML图表更新失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "UpdatePUMLDiagram: PUML图表更新成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML图表更新成功",
		"code":    http.StatusOK,
	})
}

// DeletePUML 删除PUML图表
func (pc *PUMLController) DeletePUML(c *gin.Context) {
	log.InfofId(c, "DeletePUML: 开始处理PUML图表删除请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "DeletePUML: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	pumlID := c.Param("pumlId")
	if pumlID == "" {
		log.WarnfId(c, "DeletePUML: PUML图表ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML图表ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "DeletePUML: 用户 %s 请求删除PUML图表: %s", user.UserID.String(), pumlID)

	// 调用服务删除PUML图表
	err := pc.pumlService.DeletePUML(user.UserID, pumlID)
	if err != nil {
		log.ErrorfId(c, "DeletePUML: PUML图表删除失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "DeletePUML: PUML图表删除成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PUML图表删除成功",
		"code":    http.StatusOK,
	})
}

// RenderPUMLImage 渲染PUML图片
func (pc *PUMLController) RenderPUMLImage(c *gin.Context) {
	log.InfofId(c, "RenderPUMLImage: 开始处理PUML图片渲染请求")

	var req model.RenderPUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "RenderPUMLImage: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Content == "" {
		log.WarnfId(c, "RenderPUMLImage: PUML内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "RenderPUMLImage: 开始渲染PUML图片")

	// 调用服务渲染PUML图片
	result, err := pc.pumlService.RenderPUMLImage(&req)
	if err != nil {
		log.ErrorfId(c, "RenderPUMLImage: PUML图片渲染失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "RenderPUMLImage: PUML图片渲染成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML图片渲染成功",
		"code":    http.StatusOK,
	})
}

// RenderPUMLOnline 在线渲染PUML
func (pc *PUMLController) RenderPUMLOnline(c *gin.Context) {
	log.InfofId(c, "RenderPUMLOnline: 开始处理在线PUML渲染请求")

	var req model.RenderPUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "RenderPUMLOnline: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Content == "" {
		log.WarnfId(c, "RenderPUMLOnline: PUML内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "RenderPUMLOnline: 开始在线渲染PUML")

	// 调用服务在线渲染PUML
	result, err := pc.pumlService.RenderPUMLOnlineFromRequest(&req)
	if err != nil {
		log.ErrorfId(c, "RenderPUMLOnline: 在线PUML渲染失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "RenderPUMLOnline: 在线PUML渲染成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "在线PUML渲染成功",
		"code":    http.StatusOK,
	})
}

// GenerateImage 生成图片
func (pc *PUMLController) GenerateImage(c *gin.Context) {
	log.InfofId(c, "GenerateImage: 开始处理图片生成请求")

	var req model.GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "GenerateImage: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Content == "" {
		log.WarnfId(c, "GenerateImage: 内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GenerateImage: 开始生成图片")

	// 调用服务生成图片
	result, err := pc.pumlService.GenerateImage(&req)
	if err != nil {
		log.ErrorfId(c, "GenerateImage: 图片生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GenerateImage: 图片生成成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "图片生成成功",
		"code":    http.StatusOK,
	})
}

// ValidatePUML 验证PUML语法
func (pc *PUMLController) ValidatePUML(c *gin.Context) {
	log.InfofId(c, "ValidatePUML: 开始处理PUML语法验证请求")

	var req model.ValidatePUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "ValidatePUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Content == "" {
		log.WarnfId(c, "ValidatePUML: PUML内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "ValidatePUML: 开始验证PUML语法")

	// 调用服务验证PUML语法
	result, err := pc.pumlService.ValidatePUMLFromRequest(&req)
	if err != nil {
		log.ErrorfId(c, "ValidatePUML: PUML语法验证失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "ValidatePUML: PUML语法验证成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML语法验证成功",
		"code":    http.StatusOK,
	})
}

// PreviewPUML 预览PUML
func (pc *PUMLController) PreviewPUML(c *gin.Context) {
	log.InfofId(c, "PreviewPUML: 开始处理PUML预览请求")

	var req model.PreviewPUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "PreviewPUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数验证
	if req.Content == "" {
		log.WarnfId(c, "PreviewPUML: PUML内容不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PUML内容不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "PreviewPUML: 开始预览PUML")

	// 调用服务预览PUML
	result, err := pc.pumlService.PreviewPUML(&req)
	if err != nil {
		log.ErrorfId(c, "PreviewPUML: PUML预览失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "PreviewPUML: PUML预览成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML预览成功",
		"code":    http.StatusOK,
	})
}

// ExportPUML 导出PUML
func (pc *PUMLController) ExportPUML(c *gin.Context) {
	log.InfofId(c, "ExportPUML: 开始处理PUML导出请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "ExportPUML: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req model.ExportPUMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "ExportPUML: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "ExportPUML: 用户 %s 请求导出PUML", user.UserID.String())

	// 调用服务导出PUML
	result, err := pc.pumlService.ExportPUML(user.UserID, &req)
	if err != nil {
		log.ErrorfId(c, "ExportPUML: PUML导出失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "ExportPUML: PUML导出成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "PUML导出成功",
		"code":    http.StatusOK,
	})
}

// GetPUMLStats 获取PUML统计信息
func (pc *PUMLController) GetPUMLStats(c *gin.Context) {
	log.InfofId(c, "GetPUMLStats: 开始获取PUML统计信息")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetPUMLStats: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "GetPUMLStats: 用户 %s 请求获取PUML统计信息", user.UserID.String())

	// 调用服务获取PUML统计信息
	stats, err := pc.pumlService.GetPUMLStats(user.UserID)
	if err != nil {
		log.ErrorfId(c, "GetPUMLStats: 获取PUML统计信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetPUMLStats: 成功获取PUML统计信息")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "获取PUML统计信息成功",
		"code":    http.StatusOK,
	})
}

// ClearPUMLCache 清空PUML缓存
func (pc *PUMLController) ClearPUMLCache(c *gin.Context) {
	log.InfofId(c, "ClearPUMLCache: 开始清空PUML缓存")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "ClearPUMLCache: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "ClearPUMLCache: 用户 %s 请求清空PUML缓存", user.UserID.String())

	// 调用服务清空PUML缓存
	err := pc.pumlService.ClearPUMLCache(user.UserID)
	if err != nil {
		log.ErrorfId(c, "ClearPUMLCache: 清空PUML缓存失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "ClearPUMLCache: PUML缓存清空成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PUML缓存清空成功",
		"code":    http.StatusOK,
	})
}
