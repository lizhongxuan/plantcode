package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// APIResponse 标准API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// PaginatedResponse 带分页的响应
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
	Message    string         `json:"message,omitempty"`
	Code       int            `json:"code"`
}

// 密码哈希参数
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

// JWT Claims 结构
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

// 响应工具函数

// WriteSuccessResponse 写入成功响应
func WriteSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Code:    http.StatusOK,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// WriteCreatedResponse 写入创建成功响应
func WriteCreatedResponse(w http.ResponseWriter, data interface{}, message string) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Message: message,
		Code:    http.StatusCreated,
	}
	writeJSONResponse(w, http.StatusCreated, response)
}

// WriteErrorResponse 写入错误响应
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}
	writeJSONResponse(w, statusCode, response)
}

// WritePaginatedResponse 写入分页响应
func WritePaginatedResponse(w http.ResponseWriter, data interface{}, pagination PaginationInfo, message string) {
	response := PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Message:    message,
		Code:       http.StatusOK,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// writeJSONResponse 写入JSON响应
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 密码处理函数

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}

	// 使用Argon2id哈希密码
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// 组合盐和哈希值
	result := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%x$%x",
		argon2.Version, argonMemory, argonTime, argonThreads, salt, hash)

	return result, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false
	}

	// 解析参数
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	// 解析盐
	salt := make([]byte, saltLen)
	_, err = fmt.Sscanf(parts[4], "%x", &salt)
	if err != nil {
		return false
	}

	// 解析哈希值
	hash := make([]byte, argonKeyLen)
	_, err = fmt.Sscanf(parts[5], "%x", &hash)
	if err != nil {
		return false
	}

	// 重新计算哈希
	otherHash := argon2.IDKey([]byte(password), salt, time, memory, threads, argonKeyLen)

	// 使用恒定时间比较防止时序攻击
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

// JWT处理函数

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID uuid.UUID, username, email, secretKey string, expiresIn int) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-dev-platform",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("令牌无效")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("无效的令牌声明")
	}

	return claims, nil
}

// ExtractBearerToken 从Authorization头中提取Bearer令牌
func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// 通用工具函数

// GenerateUUID 生成UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// ParseUUID 解析UUID字符串
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// GetPaginationParams 获取分页参数
func GetPaginationParams(r *http.Request) (page, pageSize int) {
	page = 1
	pageSize = 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			if parsed > 100 {
				pageSize = 100
			} else {
				pageSize = parsed
			}
		}
	}

	return page, pageSize
}

// CalculatePagination 计算分页信息
func CalculatePagination(page, pageSize int, total int64) PaginationInfo {
	totalPage := 0
	if total > 0 {
		totalPage = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return PaginationInfo{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}

// ValidateEmail 邮箱验证
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	
	// 基本格式检查
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	localPart := parts[0]
	domain := parts[1]
	
	// 检查本地部分和域名部分是否为空
	if localPart == "" || domain == "" {
		return false
	}
	
	// 检查域名是否包含点号且不是以点号开头或结尾
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}
	
	// 检查域名各部分不为空
	domainParts := strings.Split(domain, ".")
	for _, part := range domainParts {
		if part == "" {
			return false
		}
	}
	
	return true
}

// SanitizeString 清理字符串（移除前后空格、制表符、换行符）
func SanitizeString(s string) string {
	// 先移除前后的空白字符
	s = strings.TrimSpace(s)
	// 替换内部的制表符和换行符为空格
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	// 合并多个连续空格为一个
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

// IsValidProjectType 验证项目类型
func IsValidProjectType(projectType string) bool {
	validTypes := []string{
		"web_application",
		"mobile_app",
		"desktop_app",
		"api_service",
		"microservice",
		"data_analysis",
		"ai_ml_project",
		"other",
	}

	for _, vt := range validTypes {
		if projectType == vt {
			return true
		}
	}
	return false
}

// TimePtr 返回时间指针
func TimePtr(t time.Time) *time.Time {
	return &t
}

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// IntPtr 返回整数指针
func IntPtr(i int) *int {
	return &i
}

// BoolPtr 返回布尔指针
func BoolPtr(b bool) *bool {
	return &b
}

// ContainsString 检查字符串切片是否包含指定字符串
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveString 从字符串切片中移除指定字符串
func RemoveString(slice []string, item string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// FormatFileSize 格式化文件大小
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// GetClientIP 获取客户端IP地址
func GetClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		ips := strings.Split(xf, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}

	// 使用RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
} 