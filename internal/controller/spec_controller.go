package controller

import (
	"net/http"

	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SpecController Spec工作流控制器
type SpecController struct {
	specService *service.SpecService
}

// NewSpecController 创建Spec控制器
func NewSpecController(specService *service.SpecService) *SpecController {
	return &SpecController{
		specService: specService,
	}
}

// GetSpec 获取项目规格
func (sc *SpecController) GetSpec(c *gin.Context) {
	log.InfofId(c, "GetSpec: Spec工作流功能暂未完全实现")
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Spec工作流功能暂未完全实现，请使用项目详情页面的其他功能",
		"code":    http.StatusNotImplemented,
		"message": "建议使用AI助手功能进行需求分析和文档生成",
	})
}

// CreateRequirements 创建需求
func (sc *SpecController) CreateRequirements(c *gin.Context) {
	// 获取项目ID
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		log.ErrorfId(c, "Invalid project ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 获取用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		log.ErrorfId(c, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		log.ErrorfId(c, "Invalid user ID type in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid user context",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	// 解析请求体
	var req model.GenerateRequirementsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 设置项目ID
	req.ProjectID = projectID

	// 检查SpecService是否初始化
	if sc.specService == nil {
		log.ErrorfId(c, "SpecService not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Service not available",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// 调用服务生成需求文档
	reqDoc, err := sc.specService.GenerateRequirements(c.Request.Context(), userID, &req)
	if err != nil {
		log.ErrorfId(c, "Failed to generate requirements: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate requirements document",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "Successfully generated requirements document for project %s", projectID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reqDoc,
		"message": "Requirements document generated successfully",
	})
}

// CreateDesign 创建设计
func (sc *SpecController) CreateDesign(c *gin.Context) {
	log.InfofId(c, "CreateDesign: Spec工作流功能暂未完全实现")
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Spec工作流功能暂未完全实现，请使用AI功能进行设计分析",
		"code":    http.StatusNotImplemented,
		"message": "建议使用项目详情页面的AI助手进行设计文档生成",
	})
}

// CreateTasks 创建任务
func (sc *SpecController) CreateTasks(c *gin.Context) {
	log.InfofId(c, "CreateTasks: Spec工作流功能暂未完全实现")
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Spec工作流功能暂未完全实现，请使用AI功能进行任务规划",
		"code":    http.StatusNotImplemented,
		"message": "建议使用项目详情页面的AI助手进行任务分解",
	})
}

// UpdateSpec 更新规格
func (sc *SpecController) UpdateSpec(c *gin.Context) {
	log.InfofId(c, "UpdateSpec: Spec工作流功能暂未完全实现")
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Spec工作流功能暂未完全实现",
		"code":    http.StatusNotImplemented,
		"message": "功能开发中，敬请期待",
	})
}

// DeleteSpec 删除规格
func (sc *SpecController) DeleteSpec(c *gin.Context) {
	log.InfofId(c, "DeleteSpec: Spec工作流功能暂未完全实现")
	
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "Spec工作流功能暂未完全实现",
		"code":    http.StatusNotImplemented,
		"message": "功能开发中，敬请期待",
	})
}