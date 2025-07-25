package api

import (
	"ai-dev-platform/internal/api/middleware"
	"ai-dev-platform/internal/controller"
	"net/http"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// NewGinRouter 创建 Gin 路由器
func NewGinRouter(cfg *config.Config, userService service.UserService, projectService service.ProjectService, aiService *service.AIService, pumlService *service.PUMLService, asyncTaskService *service.AsyncTaskService, specService *service.SpecService) *gin.Engine {
	// 创建 Gin 引擎
	r := gin.New()

	// 添加自定义中间件
	r.Use(middleware.RequestIdRouter())    // 请求ID和链路追踪
	r.Use(middleware.RecoveryMiddleware()) // 恢复中间件
	r.Use(middleware.LoggingMiddleware())  // 日志中间件
	r.Use(middleware.SecurityMiddleware()) // 安全头中间件
	r.Use(middleware.CORSMiddleware(cfg))  // CORS中间件

	// 创建控制器
	userController := controller.NewUserController(userService)
	projectController := controller.NewProjectController(projectService)
	aiController := controller.NewAIController(aiService)
	pumlController := controller.NewPUMLController(pumlService, aiService)
	asyncController := controller.NewAsyncController(asyncTaskService, aiService)
	specController := controller.NewSpecController(specService)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API 路由组
	api := r.Group("/api")

	// 认证相关路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", userController.RegisterUser)
		auth.POST("/login", userController.LoginUser)
		auth.GET("/validate", middleware.AuthMiddleware(userService), userController.ValidateToken)
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(userService))
	{
		// 用户管理
		user := protected.Group("/user")
		{
			user.GET("/profile", userController.GetCurrentUser)
			user.PUT("/profile/update", userController.UpdateCurrentUser)
		}

		// 项目管理
		projects := protected.Group("/v1/projects")
		{
			projects.POST("", projectController.CreateProject)
			projects.GET("", projectController.GetUserProjects)
			projects.GET("/:id", projectController.GetProject)
			projects.PUT("/:id", projectController.UpdateProject)
			projects.DELETE("/:id", projectController.DeleteProject)
		}

		// AI 功能
		ai := protected.Group("/ai")
		{
			ai.POST("/analyze", aiController.AnalyzeRequirement)
			ai.GET("/analysis/:id", aiController.GetRequirementAnalysis)
			ai.GET("/analysis/project/:projectId", aiController.GetRequirementAnalysesByProject)
			ai.POST("/puml/generate", aiController.GeneratePUML)
			ai.GET("/puml/project/:projectId", aiController.GetPUMLDiagramsByProjectID)
			ai.PUT("/puml/:id", aiController.UpdatePUML)
			ai.POST("/document/generate", aiController.GenerateDocument)
			ai.GET("/document/project/:projectId", aiController.GetDocumentsByProjectID)
			ai.PUT("/document/:id", aiController.UpdateDocument)
			ai.POST("/chat/session", aiController.CreateChatSession)
			ai.POST("/chat/message", aiController.SendChatMessage)
			ai.GET("/chat/session/:sessionId/messages", aiController.GetChatMessages)
			ai.GET("/providers", aiController.GetAIProviders)
			ai.GET("/config", aiController.GetUserAIConfig)
			ai.PUT("/config", aiController.UpdateUserAIConfig)
			ai.POST("/test-connection", aiController.TestAIConnection)
			ai.GET("/models/:provider", aiController.GetAvailableModels)
			ai.POST("/chat", aiController.ProjectChat)
			ai.POST("/generate-stage-documents", aiController.GenerateStageDocuments)
			ai.POST("/generate-document-list", aiController.GenerateStageDocumentList)
		}

		// 异步任务
		async := protected.Group("/async")
		{
			async.POST("/stage-documents", asyncController.StartStageDocumentGeneration)
			async.POST("/complete-project-documents", asyncController.StartCompleteProjectDocumentGeneration)
			async.GET("/tasks/:taskId/status", asyncController.GetTaskStatus)
			async.GET("/tasks/:taskId/poll", asyncController.PollTaskStatus)
			async.GET("/projects/:projectId/progress", asyncController.GetProjectProgress)
			async.GET("/projects/:projectId/stages/:stage/documents", asyncController.GetStageDocuments)
		}

		// PUML 功能（需要认证）
		puml := protected.Group("/puml")
		{
			puml.POST("/create", pumlController.CreatePUML)
			puml.GET("/project/:projectId", pumlController.GetProjectPUMLs)
			puml.PUT("/:pumlId", pumlController.UpdatePUMLDiagram)
			puml.DELETE("/:pumlId", pumlController.DeletePUML)
			puml.POST("/export", pumlController.ExportPUML)
			puml.GET("/stats", pumlController.GetPUMLStats)
			puml.POST("/cache/clear", pumlController.ClearPUMLCache)
		}

		// Spec 工作流
		spec := protected.Group("/projects/:projectId/spec")
		{
			spec.GET("", specController.GetSpec)
			spec.POST("/requirements", specController.CreateRequirements)
			spec.POST("/design", specController.CreateDesign)
			spec.POST("/tasks", specController.CreateTasks)
			spec.PUT("", specController.UpdateSpec)
			spec.DELETE("", specController.DeleteSpec)
		}
	}

	// 公共 PUML 功能（不需要认证）
	pumlPublic := api.Group("/puml")
	{
		pumlPublic.POST("/render", pumlController.RenderPUMLImage)
		pumlPublic.POST("/render-online", pumlController.RenderPUMLOnline)
		pumlPublic.POST("/generate-image", pumlController.GenerateImage)
		pumlPublic.POST("/validate", pumlController.ValidatePUML)
		pumlPublic.POST("/preview", pumlController.PreviewPUML)
	}

	// 静态文件（开发环境）
	if cfg.IsDevelopment() {
		r.Static("/static", "./web/public")
		r.StaticFile("/", "./web/public/index.html")
	}

	return r
}
