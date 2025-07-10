package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	mockUserService *MockUserService
	middlewares     *Middlewares
	cfg             *config.Config
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.mockUserService = new(MockUserService)
	suite.cfg = &config.Config{
		CORS: config.CORSConfig{
			Origins:     []string{"http://localhost:3000", "https://example.com"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"Content-Type", "Authorization"},
			Credentials: true,
		},
		Env: "development",
	}
	suite.middlewares = NewMiddlewares(suite.cfg, suite.mockUserService)
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	suite.mockUserService.AssertExpectations(suite.T())
}

// 创建测试处理器
func testHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
}

func (suite *MiddlewareTestSuite) TestRequestID() {
	// Arrange
	handler := suite.middlewares.RequestID(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "test response", w.Body.String())
	
	// 检查是否设置了X-Request-ID头
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(suite.T(), requestID)
	
	// 验证是否为有效的UUID
	_, err := uuid.Parse(requestID)
	assert.NoError(suite.T(), err)
}

func (suite *MiddlewareTestSuite) TestLogging() {
	// Arrange
	// 先添加RequestID中间件，因为Logging依赖它
	handler := suite.middlewares.RequestID(suite.middlewares.Logging(testHandler()))
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "test response", w.Body.String())
	// 日志输出在这个测试中很难直接验证，但我们可以确保中间件没有破坏请求处理
}

func (suite *MiddlewareTestSuite) TestCORS_AllowedOrigin() {
	// Arrange
	handler := suite.middlewares.CORS(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(suite.T(), "86400", w.Header().Get("Access-Control-Max-Age"))
}

func (suite *MiddlewareTestSuite) TestCORS_DisallowedOrigin() {
	// Arrange
	// 修改配置为非开发环境
	suite.cfg.Env = "production"
	handler := suite.middlewares.CORS(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://malicious.com")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Empty(suite.T(), w.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *MiddlewareTestSuite) TestCORS_DevelopmentAllowsAll() {
	// Arrange
	suite.cfg.Env = "development"
	handler := suite.middlewares.CORS(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://any-origin.com")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "http://any-origin.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *MiddlewareTestSuite) TestCORS_PreflightRequest() {
	// Arrange
	handler := suite.middlewares.CORS(testHandler())
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
	assert.Equal(suite.T(), "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(suite.T(), w.Body.String()) // OPTIONS请求不应该有响应体
}

func (suite *MiddlewareTestSuite) TestAuth_Success() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	suite.mockUserService.On("ValidateToken", "valid-token").Return(user, nil)

	handler := suite.middlewares.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证用户是否正确设置到上下文中
		contextUser := MustGetUserFromContext(r.Context())
		assert.Equal(suite.T(), user.UserID, contextUser.UserID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "authenticated", w.Body.String())
}

func (suite *MiddlewareTestSuite) TestAuth_MissingToken() {
	// Arrange
	handler := suite.middlewares.Auth(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "未提供认证令牌")
}

func (suite *MiddlewareTestSuite) TestAuth_InvalidTokenFormat() {
	// Arrange
	handler := suite.middlewares.Auth(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "无效的认证令牌格式")
}

func (suite *MiddlewareTestSuite) TestAuth_InvalidToken() {
	// Arrange
	suite.mockUserService.On("ValidateToken", "invalid-token").Return(nil, fmt.Errorf("token invalid"))

	handler := suite.middlewares.Auth(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "无效的认证令牌")
}

func (suite *MiddlewareTestSuite) TestOptionalAuth_WithValidToken() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Status:   model.UserStatusActive,
	}

	suite.mockUserService.On("ValidateToken", "valid-token").Return(user, nil)

	handler := suite.middlewares.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextUser, exists := GetUserFromContext(r.Context())
		assert.True(suite.T(), exists)
		assert.Equal(suite.T(), user.UserID, contextUser.UserID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "authenticated", w.Body.String())
}

func (suite *MiddlewareTestSuite) TestOptionalAuth_WithoutToken() {
	// Arrange
	handler := suite.middlewares.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, exists := GetUserFromContext(r.Context())
		assert.False(suite.T(), exists)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no auth"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "no auth", w.Body.String())
}

func (suite *MiddlewareTestSuite) TestOptionalAuth_WithInvalidToken() {
	// Arrange
	suite.mockUserService.On("ValidateToken", "invalid-token").Return(nil, fmt.Errorf("token invalid"))

	handler := suite.middlewares.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, exists := GetUserFromContext(r.Context())
		assert.False(suite.T(), exists) // 无效token应该被忽略
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no auth"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "no auth", w.Body.String())
}

func (suite *MiddlewareTestSuite) TestRateLimiting_AllowedRequests() {
	// Arrange
	handler := suite.middlewares.RateLimiting(5)(testHandler()) // 每分钟5个请求
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	// Act & Assert - 发送5个请求，都应该成功
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}
}

func (suite *MiddlewareTestSuite) TestRateLimiting_ExceedLimit() {
	// Arrange
	handler := suite.middlewares.RateLimiting(2)(testHandler()) // 每分钟2个请求
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	// Act - 发送3个请求
	// 前两个应该成功
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusOK, w.Code)
	}

	// 第三个应该被限制
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "请求过于频繁")
}

func (suite *MiddlewareTestSuite) TestRecovery() {
	// Arrange
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	
	handler := suite.middlewares.Recovery(panicHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "服务器内部错误")
}

func (suite *MiddlewareTestSuite) TestContentType() {
	// Arrange
	handler := suite.middlewares.ContentType(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
}

func (suite *MiddlewareTestSuite) TestTimeout() {
	// Arrange
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // 超过超时时间
		w.WriteHeader(http.StatusOK)
	})
	
	handler := suite.middlewares.Timeout(100*time.Millisecond)(slowHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusRequestTimeout, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "请求超时")
}

func (suite *MiddlewareTestSuite) TestSecurity() {
	// Arrange
	handler := suite.middlewares.Security(testHandler())
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(suite.T(), "deny", w.Header().Get("X-Frame-Options"))
	assert.Equal(suite.T(), "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(suite.T(), "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}

func (suite *MiddlewareTestSuite) TestGetUserFromContext() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
	}
	
	ctx := context.WithValue(context.Background(), UserContextKey, user)

	// Act
	resultUser, exists := GetUserFromContext(ctx)

	// Assert
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), user, resultUser)
}

func (suite *MiddlewareTestSuite) TestGetUserFromContext_NotExists() {
	// Arrange
	ctx := context.Background()

	// Act
	resultUser, exists := GetUserFromContext(ctx)

	// Assert
	assert.False(suite.T(), exists)
	assert.Nil(suite.T(), resultUser)
}

func (suite *MiddlewareTestSuite) TestGetRequestIDFromContext() {
	// Arrange
	requestID := "test-request-id"
	ctx := context.WithValue(context.Background(), RequestIDKey, requestID)

	// Act
	result := GetRequestIDFromContext(ctx)

	// Assert
	assert.Equal(suite.T(), requestID, result)
}

func (suite *MiddlewareTestSuite) TestGetRequestIDFromContext_NotExists() {
	// Arrange
	ctx := context.Background()

	// Act
	result := GetRequestIDFromContext(ctx)

	// Assert
	assert.Empty(suite.T(), result)
}

func (suite *MiddlewareTestSuite) TestMustGetUserFromContext() {
	// Arrange
	user := &model.User{
		UserID:   uuid.New(),
		Username: "testuser",
	}
	
	ctx := context.WithValue(context.Background(), UserContextKey, user)

	// Act
	resultUser := MustGetUserFromContext(ctx)

	// Assert
	assert.Equal(suite.T(), user, resultUser)
}

func (suite *MiddlewareTestSuite) TestMustGetUserFromContext_Panic() {
	// Arrange
	ctx := context.Background()

	// Act & Assert
	assert.Panics(suite.T(), func() {
		MustGetUserFromContext(ctx)
	})
}

func (suite *MiddlewareTestSuite) TestChain() {
	// Arrange
	var order []string
	
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware1")
			next.ServeHTTP(w, r)
		})
	}
	
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware2")
			next.ServeHTTP(w, r)
		})
	}
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Act
	chainedHandler := Chain(middleware1, middleware2)(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	chainedHandler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), []string{"middleware1", "middleware2", "handler"}, order)
}

func (suite *MiddlewareTestSuite) TestApply() {
	// Arrange
	var order []string
	
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware1")
			next.ServeHTTP(w, r)
		})
	}
	
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware2")
			next.ServeHTTP(w, r)
		})
	}
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Act
	appliedHandler := Apply(handler, middleware1, middleware2)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	appliedHandler.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), []string{"middleware1", "middleware2", "handler"}, order)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
} 