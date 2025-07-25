package controller

import (
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AsyncController 异步任务控制器
type AsyncController struct {
	asyncTaskService *service.AsyncTaskService
	aiService        *service.AIService
}

// NewAsyncController 创建异步任务控制器
func NewAsyncController(asyncTaskService *service.AsyncTaskService, aiService *service.AIService) *AsyncController {
	return &AsyncController{
		asyncTaskService: asyncTaskService,
		aiService:        aiService,
	}
}

// StartStageDocumentGeneration 启动阶段文档生成任务 - 暂时不实现
func (ac *AsyncController) StartStageDocumentGeneration(c *gin.Context) {
	log.InfofId(c, "StartStageDocumentGeneration: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// StartCompleteProjectDocumentGeneration 启动完整项目文档生成任务 - 暂时不实现
func (ac *AsyncController) StartCompleteProjectDocumentGeneration(c *gin.Context) {
	log.InfofId(c, "StartCompleteProjectDocumentGeneration: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GetTaskStatus 获取任务状态 - 暂时不实现
func (ac *AsyncController) GetTaskStatus(c *gin.Context) {
	log.InfofId(c, "GetTaskStatus: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// PollTaskStatus 轮询任务状态 - 暂时不实现
func (ac *AsyncController) PollTaskStatus(c *gin.Context) {
	log.InfofId(c, "PollTaskStatus: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GetProjectProgress 获取项目进度 - 暂时不实现
func (ac *AsyncController) GetProjectProgress(c *gin.Context) {
	log.InfofId(c, "GetProjectProgress: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}

// GetStageDocuments 获取阶段文档 - 暂时不实现
func (ac *AsyncController) GetStageDocuments(c *gin.Context) {
	log.InfofId(c, "GetStageDocuments: 异步任务功能暂时不可用")
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error":   "异步任务功能暂时不可用",
		"code":    http.StatusNotImplemented,
	})
}