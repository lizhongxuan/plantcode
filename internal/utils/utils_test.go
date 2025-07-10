package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

// 密码哈希测试
func (suite *UtilsTestSuite) TestHashPassword() {
	// Arrange
	password := "testPassword123"

	// Act
	hash, err := HashPassword(password)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hash)
	assert.Contains(suite.T(), hash, "$argon2id$")
	assert.True(suite.T(), len(hash) > 50) // 哈希应该足够长
}

func (suite *UtilsTestSuite) TestVerifyPassword() {
	// Arrange
	password := "testPassword123"
	wrongPassword := "wrongPassword"
	
	hash, err := HashPassword(password)
	assert.NoError(suite.T(), err)

	// Act & Assert
	// 正确密码验证
	assert.True(suite.T(), VerifyPassword(password, hash))
	
	// 错误密码验证
	assert.False(suite.T(), VerifyPassword(wrongPassword, hash))
	
	// 无效哈希格式
	assert.False(suite.T(), VerifyPassword(password, "invalid-hash"))
	
	// 空密码
	assert.False(suite.T(), VerifyPassword("", hash))
	
	// 空哈希
	assert.False(suite.T(), VerifyPassword(password, ""))
}

func (suite *UtilsTestSuite) TestVerifyPassword_InvalidHashFormat() {
	// Arrange
	password := "test"
	invalidHashes := []string{
		"invalid",
		"$argon2id$v=19$",
		"$argon2id$v=19$m=65536,t=1,p=4$",
		"$argon2id$v=19$m=65536,t=1,p=4$salt$",
		"$argon2id$v=19$m=invalid,t=1,p=4$salt$hash",
	}

	// Act & Assert
	for _, hash := range invalidHashes {
		assert.False(suite.T(), VerifyPassword(password, hash), "应该拒绝无效格式: %s", hash)
	}
}

// JWT测试
func (suite *UtilsTestSuite) TestGenerateJWT() {
	// Arrange
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	secret := "test-secret"
	expiresIn := 3600

	// Act
	token, err := GenerateJWT(userID, username, email, secret, expiresIn)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	
	// 验证token结构
	parts := bytes.Split([]byte(token), []byte("."))
	assert.Equal(suite.T(), 3, len(parts)) // JWT应该有3部分
}

func (suite *UtilsTestSuite) TestValidateJWT() {
	// Arrange
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	secret := "test-secret"
	expiresIn := 3600

	token, err := GenerateJWT(userID, username, email, secret, expiresIn)
	assert.NoError(suite.T(), err)

	// Act
	claims, err := ValidateJWT(token, secret)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), username, claims.Username)
	assert.Equal(suite.T(), email, claims.Email)
	assert.Equal(suite.T(), "ai-dev-platform", claims.Issuer)
}

func (suite *UtilsTestSuite) TestValidateJWT_InvalidToken() {
	// Arrange
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	
	userID := uuid.New()
	token, err := GenerateJWT(userID, "test", "test@example.com", secret, 3600)
	assert.NoError(suite.T(), err)

	// Act & Assert
	// 错误的密钥
	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(suite.T(), err)
	
	// 无效的token格式
	_, err = ValidateJWT("invalid.token.format", secret)
	assert.Error(suite.T(), err)
	
	// 空token
	_, err = ValidateJWT("", secret)
	assert.Error(suite.T(), err)
}

func (suite *UtilsTestSuite) TestValidateJWT_ExpiredToken() {
	// Arrange - 创建一个已过期的token
	claims := &JWTClaims{
		UserID:   uuid.New(),
		Username: "test",
		Email:    "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 过期
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "ai-dev-platform",
		},
	}
	
	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(suite.T(), err)

	// Act
	_, err = ValidateJWT(tokenString, secret)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "token is expired")
}

// Bearer Token测试
func (suite *UtilsTestSuite) TestExtractBearerToken() {
	// Valid Bearer token
	token := ExtractBearerToken("Bearer abc123def456")
	assert.Equal(suite.T(), "abc123def456", token)
	
	// No Bearer prefix
	token = ExtractBearerToken("abc123def456")
	assert.Equal(suite.T(), "", token)
	
	// Only Bearer
	token = ExtractBearerToken("Bearer")
	assert.Equal(suite.T(), "", token)
	
	// Empty string
	token = ExtractBearerToken("")
	assert.Equal(suite.T(), "", token)
	
	// Different case
	token = ExtractBearerToken("bearer abc123def456")
	assert.Equal(suite.T(), "", token)
}

// UUID测试
func (suite *UtilsTestSuite) TestGenerateUUID() {
	// Act
	id1 := GenerateUUID()
	id2 := GenerateUUID()

	// Assert
	assert.NotEqual(suite.T(), id1, id2)
	assert.NotEqual(suite.T(), uuid.Nil, id1)
	assert.NotEqual(suite.T(), uuid.Nil, id2)
}

func (suite *UtilsTestSuite) TestParseUUID() {
	// Valid UUID
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id, err := ParseUUID(validUUID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), validUUID, id.String())
	
	// Invalid UUID
	_, err = ParseUUID("invalid-uuid")
	assert.Error(suite.T(), err)
	
	// Empty string
	_, err = ParseUUID("")
	assert.Error(suite.T(), err)
}

// 分页测试
func (suite *UtilsTestSuite) TestGetPaginationParams() {
	// Valid parameters
	req := httptest.NewRequest("GET", "/?page=2&page_size=20", nil)
	page, pageSize := GetPaginationParams(req)
	assert.Equal(suite.T(), 2, page)
	assert.Equal(suite.T(), 20, pageSize)
	
	// Default values
	req = httptest.NewRequest("GET", "/", nil)
	page, pageSize = GetPaginationParams(req)
	assert.Equal(suite.T(), 1, page)
	assert.Equal(suite.T(), 10, pageSize)
	
	// Invalid page
	req = httptest.NewRequest("GET", "/?page=invalid&page_size=abc", nil)
	page, pageSize = GetPaginationParams(req)
	assert.Equal(suite.T(), 1, page)
	assert.Equal(suite.T(), 10, pageSize)
	
	// Zero page
	req = httptest.NewRequest("GET", "/?page=0&page_size=0", nil)
	page, pageSize = GetPaginationParams(req)
	assert.Equal(suite.T(), 1, page)
	assert.Equal(suite.T(), 10, pageSize)
	
	// Large page size
	req = httptest.NewRequest("GET", "/?page=1&page_size=200", nil)
	page, pageSize = GetPaginationParams(req)
	assert.Equal(suite.T(), 1, page)
	assert.Equal(suite.T(), 100, pageSize) // 应该被限制到最大值
}

func (suite *UtilsTestSuite) TestCalculatePagination() {
	// Test case 1: 正常情况
	pagination := CalculatePagination(1, 10, 25)
	assert.Equal(suite.T(), 1, pagination.Page)
	assert.Equal(suite.T(), 10, pagination.PageSize)
	assert.Equal(suite.T(), int64(25), pagination.Total)
	assert.Equal(suite.T(), 3, pagination.TotalPage)
	
	// Test case 2: 恰好整除
	pagination = CalculatePagination(2, 10, 20)
	assert.Equal(suite.T(), 2, pagination.Page)
	assert.Equal(suite.T(), 10, pagination.PageSize)
	assert.Equal(suite.T(), int64(20), pagination.Total)
	assert.Equal(suite.T(), 2, pagination.TotalPage)
	
	// Test case 3: 无数据
	pagination = CalculatePagination(1, 10, 0)
	assert.Equal(suite.T(), 1, pagination.Page)
	assert.Equal(suite.T(), 10, pagination.PageSize)
	assert.Equal(suite.T(), int64(0), pagination.Total)
	assert.Equal(suite.T(), 0, pagination.TotalPage)
	
	// Test case 4: 单页数据小于页面大小
	pagination = CalculatePagination(1, 10, 5)
	assert.Equal(suite.T(), 1, pagination.Page)
	assert.Equal(suite.T(), 10, pagination.PageSize)
	assert.Equal(suite.T(), int64(5), pagination.Total)
	assert.Equal(suite.T(), 1, pagination.TotalPage)
}

// 验证函数测试
func (suite *UtilsTestSuite) TestValidateEmail() {
	// Valid emails
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+label@example.org",
		"123@example.com",
	}
	
	for _, email := range validEmails {
		assert.True(suite.T(), ValidateEmail(email), "应该验证通过: %s", email)
	}
	
	// Invalid emails
	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"test@",
		"test.example.com",
		"test@example",
		"test@.com",
	}
	
	for _, email := range invalidEmails {
		assert.False(suite.T(), ValidateEmail(email), "应该验证失败: %s", email)
	}
}

func (suite *UtilsTestSuite) TestSanitizeString() {
	// Normal string
	result := SanitizeString("  Hello World  ")
	assert.Equal(suite.T(), "Hello World", result)
	
	// String with tabs and newlines
	result = SanitizeString("  Hello\t\nWorld  ")
	assert.Equal(suite.T(), "Hello World", result)
	
	// Empty string
	result = SanitizeString("")
	assert.Equal(suite.T(), "", result)
	
	// Only whitespace
	result = SanitizeString("   \t\n   ")
	assert.Equal(suite.T(), "", result)
}

func (suite *UtilsTestSuite) TestIsValidProjectType() {
	// Valid types
	validTypes := []string{"web_application", "mobile_app", "desktop_app", "api_service", "microservice", "data_analysis", "ai_ml_project", "other"}
	for _, projectType := range validTypes {
		assert.True(suite.T(), IsValidProjectType(projectType), "应该是有效类型: %s", projectType)
	}
	
	// Invalid types
	invalidTypes := []string{"", "invalid", "web", "mobile", "WEB", "Mobile", "unknown"}
	for _, projectType := range invalidTypes {
		assert.False(suite.T(), IsValidProjectType(projectType), "应该是无效类型: %s", projectType)
	}
}

// 响应处理测试
func (suite *UtilsTestSuite) TestWriteSuccessResponse() {
	// Arrange
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	message := "Success"

	// Act
	WriteSuccessResponse(w, data, message)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
	
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), message, response.Message)
	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.NotNil(suite.T(), response.Data)
}

func (suite *UtilsTestSuite) TestWriteCreatedResponse() {
	// Arrange
	w := httptest.NewRecorder()
	data := map[string]string{"id": "123"}
	message := "Created successfully"

	// Act
	WriteCreatedResponse(w, data, message)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
	
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), message, response.Message)
	assert.Equal(suite.T(), http.StatusCreated, response.Code)
	assert.NotNil(suite.T(), response.Data)
}

func (suite *UtilsTestSuite) TestWriteErrorResponse() {
	// Arrange
	w := httptest.NewRecorder()
	message := "Something went wrong"

	// Act
	WriteErrorResponse(w, http.StatusBadRequest, message)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
	
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
	assert.Equal(suite.T(), message, response.Error)
	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)
	assert.Nil(suite.T(), response.Data)
}

func (suite *UtilsTestSuite) TestWritePaginatedResponse() {
	// Arrange
	w := httptest.NewRecorder()
	data := []string{"item1", "item2"}
	pagination := PaginationInfo{
		Page:      1,
		PageSize:  10,
		Total:     2,
		TotalPage: 1,
	}
	message := "Success"

	// Act
	WritePaginatedResponse(w, data, pagination, message)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
	
	var response PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), message, response.Message)
	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.NotNil(suite.T(), response.Data)
	assert.Equal(suite.T(), pagination, response.Pagination)
}

// 辅助函数测试
func (suite *UtilsTestSuite) TestPointerFunctions() {
	// TimePtr
	now := time.Now()
	timePtr := TimePtr(now)
	assert.NotNil(suite.T(), timePtr)
	assert.Equal(suite.T(), now, *timePtr)
	
	// StringPtr
	str := "test"
	strPtr := StringPtr(str)
	assert.NotNil(suite.T(), strPtr)
	assert.Equal(suite.T(), str, *strPtr)
	
	// IntPtr
	num := 42
	intPtr := IntPtr(num)
	assert.NotNil(suite.T(), intPtr)
	assert.Equal(suite.T(), num, *intPtr)
	
	// BoolPtr
	val := true
	boolPtr := BoolPtr(val)
	assert.NotNil(suite.T(), boolPtr)
	assert.Equal(suite.T(), val, *boolPtr)
}

func (suite *UtilsTestSuite) TestContainsString() {
	slice := []string{"apple", "banana", "orange"}
	
	// Item exists
	assert.True(suite.T(), ContainsString(slice, "apple"))
	assert.True(suite.T(), ContainsString(slice, "banana"))
	assert.True(suite.T(), ContainsString(slice, "orange"))
	
	// Item doesn't exist
	assert.False(suite.T(), ContainsString(slice, "grape"))
	assert.False(suite.T(), ContainsString(slice, ""))
	assert.False(suite.T(), ContainsString(slice, "Apple")) // case sensitive
	
	// Empty slice
	assert.False(suite.T(), ContainsString([]string{}, "apple"))
}

func (suite *UtilsTestSuite) TestRemoveString() {
	slice := []string{"apple", "banana", "orange", "banana"}
	
	// Remove existing item
	result := RemoveString(slice, "banana")
	expected := []string{"apple", "orange"}
	assert.Equal(suite.T(), expected, result)
	
	// Remove non-existing item
	result = RemoveString(slice, "grape")
	assert.Equal(suite.T(), slice, result)
	
	// Remove from empty slice
	result = RemoveString([]string{}, "apple")
	assert.Equal(suite.T(), []string{}, result)
}

func (suite *UtilsTestSuite) TestTruncateString() {
	// String shorter than max length
	result := TruncateString("hello", 10)
	assert.Equal(suite.T(), "hello", result)
	
	// String equal to max length
	result = TruncateString("hello", 5)
	assert.Equal(suite.T(), "hello", result)
	
	// String longer than max length
	result = TruncateString("hello world", 5)
	assert.Equal(suite.T(), "he...", result)
	
	// Empty string
	result = TruncateString("", 5)
	assert.Equal(suite.T(), "", result)
	
	// Zero max length
	result = TruncateString("hello", 0)
	assert.Equal(suite.T(), "", result)
}

func (suite *UtilsTestSuite) TestGetClientIP() {
	// Test X-Forwarded-For header
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")
	ip := GetClientIP(req)
	assert.Equal(suite.T(), "192.168.1.1", ip)
	
	// Test X-Real-IP header
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.2")
	ip = GetClientIP(req)
	assert.Equal(suite.T(), "192.168.1.2", ip)
	
	// Test RemoteAddr
	req = httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.3:8080"
	ip = GetClientIP(req)
	assert.Equal(suite.T(), "192.168.1.3", ip)
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
} 