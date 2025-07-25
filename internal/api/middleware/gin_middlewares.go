package middleware

import (
	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/service"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS中间件
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许的源
		isAllowedOrigin := false
		for _, allowedOrigin := range cfg.CORS.Origins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				isAllowedOrigin = true
				break
			}
		}

		if isAllowedOrigin || cfg.IsDevelopment() {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.Methods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.Headers, ", "))

		if cfg.CORS.Credentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Max-Age", "86400") // 24小时

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			log.WarnfId(c, "AuthMiddleware: 未提供认证令牌")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "未提供认证令牌",
				"code":    401,
			})
			c.Abort()
			return
		}

		// 提取Bearer令牌
		token := extractBearerToken(authHeader)
		if token == "" {
			log.WarnfId(c, "AuthMiddleware: 无效的认证令牌格式")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "无效的认证令牌格式",
				"code":    401,
			})
			c.Abort()
			return
		}

		// 验证令牌
		user, err := userService.ValidateToken(token)
		if err != nil {
			log.WarnfId(c, "AuthMiddleware: 无效的认证令牌: %v", err)
			c.JSON(401, gin.H{
				"success": false,
				"error":   "无效的认证令牌",
				"code":    401,
			})
			c.Abort()
			return
		}

		// 将用户信息添加到上下文
		c.Set("user", user)
		log.DebugfId(c, "AuthMiddleware: 用户认证成功，用户ID: %s", user.UserID.String())
		c.Next()
	}
}

// SecurityMiddleware 安全头中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "deny")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			traceID := "unknown"
			requestID := "unknown"

			// 使用param.Request获取headers
			if param.Request != nil {
				if tid := param.Request.Header.Get("X-Trace-ID"); tid != "" {
					traceID = tid
				}
				if rid := param.Request.Header.Get("X-Request-ID"); rid != "" {
					requestID = rid
				}
			}

			return fmt.Sprintf("[trace_id:%s] [request_id:%s] %s %s %s %d %v - %s\n",
				traceID,
				requestID,
				param.ClientIP,
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
			)
		},
	})
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.ErrorfId(c, "PANIC: %v", recovered)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "服务器内部错误",
			"code":    500,
		})
	})
}

// extractBearerToken 从Authorization头中提取Bearer令牌
func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
