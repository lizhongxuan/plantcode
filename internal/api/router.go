package api

import (
	"net/http"
	"strings"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/service"
)

// Router 路由管理器
type Router struct {
	config        *config.Config
	handlers      *Handlers
	aiHandlers    *AIHandlers
	pumlHandlers  *PUMLHandlers
	asyncHandlers *AsyncHandlers
	middlewares   *Middlewares
	httpHandler   http.Handler // 添加缓存的handler
}

// NewRouter 创建路由器
func NewRouter(cfg *config.Config, userService service.UserService, projectService service.ProjectService, aiService *service.AIService, pumlService *service.PUMLService, asyncTaskService *service.AsyncTaskService) *Router {
	handlers := NewHandlers(userService, projectService)
	aiHandlers := NewAIHandlers(aiService)
	pumlHandlers := NewPUMLHandlers(pumlService, aiService)
	asyncHandlers := NewAsyncHandlers(asyncTaskService, aiService)
	middlewares := NewMiddlewares(cfg, userService)

	return &Router{
		config:        cfg,
		handlers:      handlers,
		aiHandlers:    aiHandlers,
		pumlHandlers:  pumlHandlers,
		asyncHandlers: asyncHandlers,
		middlewares:   middlewares,
	}
}

// SetupRoutes 设置路由
func (router *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// 公共中间件
	commonMiddlewares := []Middleware{
		router.middlewares.RequestID,
		router.middlewares.Logging,
		router.middlewares.Recovery,
		router.middlewares.Security,
		router.middlewares.CORS,
		router.middlewares.ContentType,
		router.middlewares.Timeout(30 * time.Second),
	}

	// 认证中间件
	authMiddlewares := append(commonMiddlewares, router.middlewares.Auth)

	// 限流中间件
	rateLimitMiddleware := router.middlewares.RateLimiting(60) // 每分钟60次请求

	// ===== 健康检查 =====
	mux.Handle("/health", Apply(
		http.HandlerFunc(router.handlers.Health),
		append(commonMiddlewares, rateLimitMiddleware)...,
	))

	// ===== 用户认证相关 =====
	// 用户注册
	mux.Handle("/api/auth/register", Apply(
		http.HandlerFunc(router.handlers.RegisterUser),
		append(commonMiddlewares, rateLimitMiddleware)...,
	))

	// 用户登录
	mux.Handle("/api/auth/login", Apply(
		http.HandlerFunc(router.handlers.LoginUser),
		append(commonMiddlewares, rateLimitMiddleware)...,
	))

	// 验证token
	mux.Handle("/api/auth/validate", Apply(
		http.HandlerFunc(router.handlers.ValidateToken),
		authMiddlewares...,
	))

	// ===== 用户管理 =====
	// 获取当前用户信息
	mux.Handle("/api/user/profile", Apply(
		http.HandlerFunc(router.handlers.GetCurrentUser),
		authMiddlewares...,
	))

	// 更新当前用户信息
	mux.Handle("/api/user/profile/update", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPut: router.handlers.UpdateCurrentUser,
		})),
		authMiddlewares...,
	))

	// ===== 项目管理 =====
	// 创建项目
	mux.Handle("/api/projects", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.handlers.CreateProject,
		})),
		authMiddlewares...,
	))

	// 获取用户项目列表
	mux.Handle("/api/projects/list", Apply(
		http.HandlerFunc(router.handlers.GetUserProjects),
		authMiddlewares...,
	))

	// 项目详情路由 - 处理带ID的路径
	mux.Handle("/api/projects/", Apply(
		http.HandlerFunc(router.projectHandler),
		authMiddlewares...,
	))

	// ===== AI功能相关 =====
	// 需求分析
	mux.Handle("/api/ai/analyze", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.AnalyzeRequirement,
		})),
		authMiddlewares...,
	))

	// 获取项目需求分析
	mux.Handle("/api/ai/analysis/project/", Apply(
		http.HandlerFunc(router.aiAnalysisProjectHandler),
		authMiddlewares...,
	))

	// 获取需求分析问题
	mux.Handle("/api/ai/analysis/", Apply(
		http.HandlerFunc(router.aiAnalysisHandler),
		authMiddlewares...,
	))

	// 回答问题
	mux.Handle("/api/ai/questions/", Apply(
		http.HandlerFunc(router.aiQuestionHandler),
		authMiddlewares...,
	))

	// PUML生成
	mux.Handle("/api/ai/puml/generate", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.GeneratePUML,
		})),
		authMiddlewares...,
	))

	// 获取项目PUML图表
	mux.Handle("/api/ai/puml/project/", Apply(
		http.HandlerFunc(router.aiHandlers.GetPUMLDiagramsByProject),
		authMiddlewares...,
	))

	// 更新PUML图表
	mux.Handle("/api/ai/puml/", Apply(
		http.HandlerFunc(router.aiPUMLHandler),
		authMiddlewares...,
	))

	// 文档生成
	mux.Handle("/api/ai/document/generate", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.GenerateDocument,
		})),
		authMiddlewares...,
	))

	// 获取项目技术文档
	mux.Handle("/api/ai/document/project/", Apply(
		http.HandlerFunc(router.aiHandlers.GetDocumentsByProject),
		authMiddlewares...,
	))

	// 更新技术文档
	mux.Handle("/api/ai/document/", Apply(
		http.HandlerFunc(router.aiDocumentHandler),
		authMiddlewares...,
	))

	// 对话会话管理
	mux.Handle("/api/ai/chat/session", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.CreateChatSession,
		})),
		authMiddlewares...,
	))

	// 发送消息
	mux.Handle("/api/ai/chat/message", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.SendChatMessage,
		})),
		authMiddlewares...,
	))

	// 获取对话消息
	mux.Handle("/api/ai/chat/session/", Apply(
		http.HandlerFunc(router.aiChatHandler),
		authMiddlewares...,
	))

	// AI管理接口
	mux.Handle("/api/ai/providers", Apply(
		http.HandlerFunc(router.aiHandlers.GetAIProviders),
		authMiddlewares...,
	))

	// 用户AI配置管理
	mux.Handle("/api/ai/config", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodGet: router.aiHandlers.GetUserAIConfig,
			http.MethodPut: router.aiHandlers.UpdateUserAIConfig,
		})),
		authMiddlewares...,
	))

	// AI连接测试
	mux.Handle("/api/ai/test-connection", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.TestAIConnection,
		})),
		authMiddlewares...,
	))

	// 获取可用模型列表
	mux.Handle("/api/ai/models/", Apply(
		http.HandlerFunc(router.aiHandlers.GetAvailableModels),
		authMiddlewares...,
	))

	// 项目上下文AI对话 (需要认证以获取用户AI配置)
	mux.Handle("/api/ai/chat", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.ProjectChat,
		})),
		authMiddlewares...,
	))

	// 分阶段文档生成
	mux.Handle("/api/ai/generate-stage-documents", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.GenerateStageDocuments,
		})),
		authMiddlewares...,
	))

	// 根据需求分析生成阶段文档列表
	mux.Handle("/api/ai/generate-document-list", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.GenerateStageDocumentList,
		})),
		authMiddlewares...,
	))

	// ===== 异步任务管理 =====
	// 启动阶段文档生成任务
	mux.Handle("/api/async/stage-documents", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.asyncHandlers.StartStageDocumentGeneration,
		})),
		authMiddlewares...,
	))

	// 启动完整项目文档生成任务
	mux.Handle("/api/async/complete-project-documents", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.asyncHandlers.StartCompleteProjectDocumentGeneration,
		})),
		authMiddlewares...,
	))

	// 获取任务状态
	mux.Handle("/api/async/tasks/", Apply(
		http.HandlerFunc(router.asyncTaskHandler),
		authMiddlewares...,
	))

	// 获取项目进度状态 - 修复路由冲突
	mux.Handle("/api/async/projects/", Apply(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 解析路径判断是进度查询还是文档查询
			path := r.URL.Path
			if strings.Contains(path, "/progress") {
				router.projectProgressHandler(w, r)
			} else if strings.Contains(path, "/stages/") && strings.Contains(path, "/documents") {
				router.stageDocumentsHandler(w, r)
			} else {
				router.asyncProjectHandler(w, r)
			}
		}),
		authMiddlewares...,
	))

	// ===== PUML渲染和编辑功能 (移除认证要求用于测试) =====
	// PUML渲染
	mux.Handle("/api/puml/render", 
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.RenderPUMLImage,
		})),
	)

	// 在线PUML渲染 (新增)
	mux.Handle("/api/puml/render-online",
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.RenderPUMLOnlineHandler,
		})),
	)

	// PUML图片生成 (兼容前端API)
	mux.Handle("/api/puml/generate-image", 
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.GenerateImage,
		})),
	)

	// PUML语法验证
	mux.Handle("/api/puml/validate", 
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.ValidatePUML,
		})),
	)

	// PUML预览
	mux.Handle("/api/puml/preview", 
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.PreviewPUML,
		})),
	)

	// PUML导出
	mux.Handle("/api/puml/export", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.ExportPUML,
		})),
		authMiddlewares...,
	))

	// PUML服务统计
	mux.Handle("/api/puml/stats", Apply(
		http.HandlerFunc(router.pumlHandlers.GetPUMLStats),
		authMiddlewares...,
	))

	// 清空PUML缓存
	mux.Handle("/api/puml/cache/clear", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.ClearPUMLCache,
		})),
		authMiddlewares...,
	))

	// ===== 项目PUML管理路由 =====
	// 获取项目PUML图表列表
	mux.Handle("/api/puml/project/", Apply(
		http.HandlerFunc(router.pumlHandlers.GetProjectPUMLs),
		authMiddlewares...,
	))

	// 创建PUML图表
	mux.Handle("/api/puml/create", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.CreatePUML,
		})),
		authMiddlewares...,
	))

	// PUML图表管理（更新、删除）
	mux.Handle("/api/puml/", Apply(
		http.HandlerFunc(router.pumlManagementHandler),
		authMiddlewares...,
	))

	// ===== 静态文件服务 (开发环境) =====
	if router.config.IsDevelopment() {
		// 前端静态文件
		fs := http.FileServer(http.Dir("./web/public/"))
		mux.Handle("/", fs)
	}

	return mux
}

// methodHandler 方法处理器，根据HTTP方法分发请求
func (router *Router) methodHandler(methods map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, exists := methods[r.Method]; exists {
			handler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// projectHandler 项目处理器，处理项目相关的路由
func (router *Router) projectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		router.handlers.GetProject(w, r)
	case http.MethodPut:
		router.handlers.UpdateProject(w, r)
	case http.MethodDelete:
		router.handlers.DeleteProject(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// aiAnalysisHandler AI需求分析处理器
func (router *Router) aiAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		router.aiHandlers.GetRequirementAnalysis(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// aiQuestionHandler AI问题处理器
func (router *Router) aiQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// 暂时返回404，等待实现
	http.Error(w, "Not Found", http.StatusNotFound)
}

// aiChatHandler AI聊天处理器
func (router *Router) aiChatHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 检查路径是否以 /messages 结尾
		if r.URL.Path[len(r.URL.Path)-9:] == "/messages" {
			router.aiHandlers.GetChatMessages(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// aiPUMLHandler AI PUML处理器
func (router *Router) aiPUMLHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		router.aiHandlers.UpdatePUML(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// aiDocumentHandler AI文档处理器
func (router *Router) aiDocumentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		router.aiHandlers.UpdateDocument(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// aiAnalysisProjectHandler AI需求分析项目处理器
func (router *Router) aiAnalysisProjectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		router.aiHandlers.GetRequirementAnalysesByProject(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// asyncTaskHandler 异步任务处理器
func (router *Router) asyncTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 解析路径，判断是状态查询还是轮询
	if strings.HasSuffix(r.URL.Path, "/status") {
		switch r.Method {
		case http.MethodGet:
			router.asyncHandlers.GetTaskStatus(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasSuffix(r.URL.Path, "/poll") {
		switch r.Method {
		case http.MethodGet:
			router.asyncHandlers.PollTaskStatus(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// asyncProjectHandler 异步项目处理器
func (router *Router) asyncProjectHandler(w http.ResponseWriter, r *http.Request) {
	// 判断是获取进度还是获取阶段文档
	if strings.Contains(r.URL.Path, "/progress") {
		switch r.Method {
		case http.MethodGet:
			router.asyncHandlers.GetStageProgress(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.Contains(r.URL.Path, "/stages/") && strings.HasSuffix(r.URL.Path, "/documents") {
		switch r.Method {
		case http.MethodGet:
			router.asyncHandlers.GetStageDocuments(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// GetHandler 获取处理器（用于外部调用）
func (router *Router) GetHandler() http.Handler {
	// 如果handler已经初始化，直接返回缓存的handler
	if router.httpHandler != nil {
		return router.httpHandler
	}
	
	// 初始化并缓存handler
	router.httpHandler = router.SetupRoutes()
	return router.httpHandler
}

// pumlManagementHandler PUML图表管理处理器
func (router *Router) pumlManagementHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		router.pumlHandlers.UpdatePUMLDiagram(w, r)
	case http.MethodDelete:
		router.pumlHandlers.DeletePUML(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// projectProgressHandler 项目进度处理器
func (router *Router) projectProgressHandler(w http.ResponseWriter, r *http.Request) {
	// 解析路径参数: /api/async/projects/{projectId}/progress
	path := r.URL.Path
	if r.Method == http.MethodGet && strings.Contains(path, "/progress") {
		router.asyncHandlers.GetProjectProgress(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// stageDocumentsHandler 阶段文档处理器
func (router *Router) stageDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// 解析路径参数: /api/async/projects/{projectId}/stages/{stage}/documents
	path := r.URL.Path
	if r.Method == http.MethodGet && strings.Contains(path, "/stages/") && strings.Contains(path, "/documents") {
		router.asyncHandlers.GetStageDocuments(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
} 