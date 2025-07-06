package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserContextKey 用户上下文键
	UserContextKey ContextKey = "user"
	// RequestIDKey 请求ID键
	RequestIDKey ContextKey = "request_id"
)

// Middleware 中间件函数类型
type Middleware func(http.Handler) http.Handler

// Middlewares 中间件集合
type Middlewares struct {
	config      *config.Config
	userService service.UserService
}

// NewMiddlewares 创建中间件集合
func NewMiddlewares(cfg *config.Config, userService service.UserService) *Middlewares {
	return &Middlewares{
		config:      cfg,
		userService: userService,
	}
}

// RequestID 请求ID中间件
func (m *Middlewares) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := utils.GenerateUUID().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging 日志记录中间件
func (m *Middlewares) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 包装ResponseWriter以捕获状态码
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// 获取请求ID
		requestID := ""
		if id := r.Context().Value(RequestIDKey); id != nil {
			requestID = id.(string)
		}
		
		// 获取客户端IP
		clientIP := utils.GetClientIP(r)
		
		// 执行请求
		next.ServeHTTP(wrapped, r)
		
		// 记录请求日志
		duration := time.Since(start)
		log.Printf("[%s] %s %s %s %d %v - %s",
			requestID,
			clientIP,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			r.UserAgent(),
		)
	})
}

// CORS 跨域中间件
func (m *Middlewares) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// 检查是否允许的源
		isAllowedOrigin := false
		for _, allowedOrigin := range m.config.CORS.Origins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				isAllowedOrigin = true
				break
			}
		}
		
		if isAllowedOrigin || m.config.IsDevelopment() {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.CORS.Methods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.config.CORS.Headers, ", "))
		
		if m.config.CORS.Credentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		
		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Auth JWT认证中间件
func (m *Middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "未提供认证令牌")
			return
		}
		
		// 提取Bearer令牌
		token := utils.ExtractBearerToken(authHeader)
		if token == "" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "无效的认证令牌格式")
			return
		}
		
		// 验证令牌
		user, err := m.userService.ValidateToken(token)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "无效的认证令牌")
			return
		}
		
		// 将用户信息添加到上下文
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth 可选认证中间件（令牌存在时验证，不存在时继续）
func (m *Middlewares) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := utils.ExtractBearerToken(authHeader)
			if token != "" {
				user, err := m.userService.ValidateToken(token)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					r = r.WithContext(ctx)
				}
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// RateLimiting 简单的内存限流中间件
func (m *Middlewares) RateLimiting(requestsPerMinute int) Middleware {
	// 简单的内存存储，生产环境应该使用Redis
	clients := make(map[string][]time.Time)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := utils.GetClientIP(r)
			now := time.Now()
			
			// 清理过期记录
			if times, exists := clients[clientIP]; exists {
				var validTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						validTimes = append(validTimes, t)
					}
				}
				clients[clientIP] = validTimes
			}
			
			// 检查请求数量
			if len(clients[clientIP]) >= requestsPerMinute {
				utils.WriteErrorResponse(w, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
				return
			}
			
			// 记录当前请求
			clients[clientIP] = append(clients[clientIP], now)
			
			next.ServeHTTP(w, r)
		})
	}
}

// Recovery 恢复中间件，捕获panic
func (m *Middlewares) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求ID
				requestID := ""
				if id := r.Context().Value(RequestIDKey); id != nil {
					requestID = id.(string)
				}
				
				log.Printf("[PANIC] [%s] %s %s: %v", requestID, r.Method, r.URL.Path, err)
				
				if !m.config.IsProduction() {
					utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("服务器内部错误: %v", err))
				} else {
					utils.WriteErrorResponse(w, http.StatusInternalServerError, "服务器内部错误")
				}
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// ContentType 内容类型中间件
func (m *Middlewares) ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 对于POST、PUT、PATCH请求，检查Content-Type
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				utils.WriteErrorResponse(w, http.StatusUnsupportedMediaType, "Content-Type必须是application/json")
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// Timeout 超时中间件
func (m *Middlewares) Timeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "请求超时")
	}
}

// Security 安全头中间件
func (m *Middlewares) Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 安全头
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		if m.config.IsProduction() {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		next.ServeHTTP(w, r)
	})
}

// responseWriter 包装ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	return user, ok
}

// GetRequestIDFromContext 从上下文获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// MustGetUserFromContext 从上下文获取用户信息（必须存在）
func MustGetUserFromContext(ctx context.Context) *model.User {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		panic("用户信息不存在于上下文中")
	}
	return user
}

// Chain 链式组合中间件
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply 应用中间件到处理器
func Apply(handler http.Handler, middlewares ...Middleware) http.Handler {
	return Chain(middlewares...)(handler)
} 