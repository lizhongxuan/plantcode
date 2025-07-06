package api

import (
	"net/http"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/service"
)

// Router API路由器
type Router struct {
	config      *config.Config
	handlers    *Handlers
	middlewares *Middlewares
}

// NewRouter 创建路由器
func NewRouter(cfg *config.Config, userService service.UserService, projectService service.ProjectService) *Router {
	handlers := NewHandlers(userService, projectService)
	middlewares := NewMiddlewares(cfg, userService)

	return &Router{
		config:      cfg,
		handlers:    handlers,
		middlewares: middlewares,
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

// GetHandler 获取处理器（用于外部调用）
func (router *Router) GetHandler() http.Handler {
	return router.SetupRoutes()
} 