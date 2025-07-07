package api

import (
	"net/http"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/service"
)

// Router API路由器
type Router struct {
	config       *config.Config
	handlers     *Handlers
	aiHandlers   *AIHandlers
	pumlHandlers *PUMLHandlers
	middlewares  *Middlewares
}

// NewRouter 创建路由器
func NewRouter(cfg *config.Config, userService service.UserService, projectService service.ProjectService, aiService *service.AIService, pumlService *service.PUMLService) *Router {
	handlers := NewHandlers(userService, projectService)
	aiHandlers := NewAIHandlers(aiService)
	pumlHandlers := NewPUMLHandlers(pumlService)
	middlewares := NewMiddlewares(cfg, userService)

	return &Router{
		config:       cfg,
		handlers:     handlers,
		aiHandlers:   aiHandlers,
		pumlHandlers: pumlHandlers,
		middlewares:  middlewares,
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
		http.HandlerFunc(router.aiHandlers.GetRequirementAnalysesByProject),
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

	// 项目上下文AI对话
	mux.Handle("/api/ai/chat", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.aiHandlers.ProjectChat,
		})),
		authMiddlewares...,
	))

	// ===== PUML渲染和编辑功能 =====
	// PUML渲染
	mux.Handle("/api/puml/render", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.RenderPUML,
		})),
		authMiddlewares...,
	))

	// PUML语法验证
	mux.Handle("/api/puml/validate", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.ValidatePUML,
		})),
		authMiddlewares...,
	))

	// PUML预览
	mux.Handle("/api/puml/preview", Apply(
		http.HandlerFunc(router.methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: router.pumlHandlers.PreviewPUML,
		})),
		authMiddlewares...,
	))

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

// GetHandler 获取处理器（用于外部调用）
func (router *Router) GetHandler() http.Handler {
	return router.SetupRoutes()
} 