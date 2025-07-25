package service

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ai-dev-platform/internal/config"
	"ai-dev-platform/internal/model"
	"github.com/google/uuid"
)

// PUMLService PlantUML渲染服务
type PUMLService struct {
	serverURL    string
	onlineRenderURL string
	httpClient   *http.Client
	enableCache  bool
	cache        map[string]*RenderResult // 简单内存缓存
}

// RenderResult 渲染结果
type RenderResult struct {
	ImageData   []byte    `json:"image_data"`
	Format      string    `json:"format"`
	URL         string    `json:"url,omitempty"`
	RenderedAt  time.Time `json:"rendered_at"`
	CacheKey    string    `json:"cache_key"`
}

// RenderOptions 渲染选项
type RenderOptions struct {
	Format     string `json:"format"`      // png, svg, txt
	Theme      string `json:"theme"`       // 主题 
	DPI        int    `json:"dpi"`         // DPI设置
	UseCache   bool   `json:"use_cache"`   // 是否使用缓存
	ServerMode bool   `json:"server_mode"` // 是否使用服务器模式
}

// ValidationResult PUML语法验证结果
type ValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// NewPUMLService 创建新的PUML服务
func NewPUMLService(cfg *config.PUMLConfig) *PUMLService {
	pumlServerURL := cfg.ServerURL
	if pumlServerURL == "" {
		// 使用官方在线服务器
		pumlServerURL = "http://www.plantuml.com/plantuml"
	}

	return &PUMLService{
		serverURL:    pumlServerURL,
		onlineRenderURL: fmt.Sprintf("%s/svg", pumlServerURL),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		enableCache: true,
		cache:      make(map[string]*RenderResult),
	}
}

// RenderPUMLOnline 使用POST请求在线渲染PUML，返回SVG字符串
func (s *PUMLService) RenderPUMLOnline(pumlCode string) (string, error) {
	// 创建POST请求
	req, err := http.NewRequest("POST", s.onlineRenderURL, strings.NewReader(pumlCode))
	if err != nil {
		return "", fmt.Errorf("创建PlantUML请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求PlantUML服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("PlantUML服务返回错误: %d - %s", resp.StatusCode, string(body))
	}

	// 读取响应数据
	svgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取渲染结果失败: %w", err)
	}

	return string(svgData), nil
}


// RenderPUML 渲染PUML代码为图像
func (s *PUMLService) RenderPUML(pumlCode string, options *RenderOptions) (*RenderResult, error) {
	if options == nil {
		options = &RenderOptions{
			Format:     "png",
			UseCache:   true,
			ServerMode: true,
		}
	}

	// 生成缓存键
	cacheKey := s.generateCacheKey(pumlCode, options)
	
	// 检查缓存
	if options.UseCache && s.enableCache {
		if cached, exists := s.cache[cacheKey]; exists {
			return cached, nil
		}
	}

	var result *RenderResult
	var err error

	if options.ServerMode {
		// 使用在线服务器渲染
		result, err = s.renderWithServer(pumlCode, options)
	} else {
		// 使用本地渲染（需要本地PlantUML环境）
		result, err = s.renderLocally(pumlCode, options)
	}

	if err != nil {
		return nil, err
	}

	result.CacheKey = cacheKey
	result.RenderedAt = time.Now()

	// 缓存结果
	if options.UseCache && s.enableCache {
		s.cache[cacheKey] = result
	}

	return result, nil
}

// ValidatePUML 验证PUML语法
func (s *PUMLService) ValidatePUML(pumlCode string) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 基本语法检查
	lines := strings.Split(pumlCode, "\n")
	
	hasStart := false
	hasEnd := false
	bracketCount := 0
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1
		
		// 检查开始标记
		if strings.HasPrefix(line, "@startuml") {
			if hasStart {
				result.Errors = append(result.Errors, fmt.Sprintf("行 %d: 重复的 @startuml 标记", lineNum))
				result.IsValid = false
			}
			hasStart = true
		}
		
		// 检查结束标记
		if strings.HasPrefix(line, "@enduml") {
			if hasEnd {
				result.Errors = append(result.Errors, fmt.Sprintf("行 %d: 重复的 @enduml 标记", lineNum))
				result.IsValid = false
			}
			hasEnd = true
		}
		
		// 检查括号匹配
		bracketCount += strings.Count(line, "{") - strings.Count(line, "}")
		
		// 检查常见错误
		if strings.Contains(line, "->") && !strings.Contains(line, ":") && 
		   !strings.Contains(line, "[") && !strings.Contains(line, "participant") {
			result.Warnings = append(result.Warnings, fmt.Sprintf("行 %d: 箭头可能缺少标签", lineNum))
		}
	}
	
	// 检查必需的标记
	if !hasStart {
		result.Errors = append(result.Errors, "缺少 @startuml 开始标记")
		result.IsValid = false
	}
	
	if !hasEnd {
		result.Errors = append(result.Errors, "缺少 @enduml 结束标记")
		result.IsValid = false
	}
	
	// 检查括号平衡
	if bracketCount != 0 {
		result.Errors = append(result.Errors, "括号不匹配")
		result.IsValid = false
	}
	
	return result
}

// renderWithServer 使用在线服务器渲染
func (s *PUMLService) renderWithServer(pumlCode string, options *RenderOptions) (*RenderResult, error) {
	// 将PUML代码编码为PlantUML服务器格式
	encoded, err := s.encodePUML(pumlCode)
	if err != nil {
		return nil, fmt.Errorf("编码PUML失败: %w", err)
	}
	
	// 构建请求URL
	format := options.Format
	if format == "" {
		format = "png"
	}
	
	url := fmt.Sprintf("%s/%s/%s", s.serverURL, format, encoded)
	
	// 发送HTTP请求
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求PlantUML服务器失败: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PlantUML服务器返回错误: %d", resp.StatusCode)
	}
	
	// 读取响应数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取渲染结果失败: %w", err)
	}
	
	return &RenderResult{
		ImageData: imageData,
		Format:    format,
		URL:       url,
	}, nil
}

// renderLocally 本地渲染（需要本地PlantUML环境）
func (s *PUMLService) renderLocally(pumlCode string, options *RenderOptions) (*RenderResult, error) {
	// 这里需要调用本地的PlantUML jar文件
	// 为了简化，目前返回错误，提示需要配置本地环境
	return nil, fmt.Errorf("本地渲染需要配置PlantUML环境，当前仅支持在线渲染")
}

// encodePUML 将PUML代码编码为PlantUML服务器格式
func (s *PUMLService) encodePUML(pumlCode string) (string, error) {
	// PlantUML服务器使用特殊的编码格式
	// 1. 使用zlib压缩
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(pumlCode))
	if err != nil {
		return "", err
	}
	w.Close()
	
	// 2. Base64编码
	compressed := b.Bytes()
	encoded := base64.StdEncoding.EncodeToString(compressed)
	
	// 3. 替换字符以符合URL格式
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")
	
	return encoded, nil
}

// generateCacheKey 生成缓存键
func (s *PUMLService) generateCacheKey(pumlCode string, options *RenderOptions) string {
	return fmt.Sprintf("%s_%s_%d", 
		base64.StdEncoding.EncodeToString([]byte(pumlCode)), 
		options.Format, 
		options.DPI)
}

// ClearCache 清空缓存
func (s *PUMLService) ClearCache() {
	s.cache = make(map[string]*RenderResult)
}

// GetCacheStats 获取缓存统计
func (s *PUMLService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_size": len(s.cache),
		"cache_enabled": s.enableCache,
	}
}

// ===== Controller需要的方法 =====

// CreatePUML 创建PUML图表
func (s *PUMLService) CreatePUML(userID uuid.UUID, req *model.CreatePUMLRequest) (*model.PUMLDiagram, error) {
	// 验证PUML语法
	validation := s.ValidatePUMLString(req.Content)
	if !validation.IsValid {
		return nil, fmt.Errorf("PUML语法错误: %v", validation.Errors)
	}

	// 创建PUML图表记录（这里应该调用repository层，暂时返回模拟数据）
	diagram := &model.PUMLDiagram{
		DiagramID:   uuid.New(),
		ProjectID:   uuid.MustParse(req.ProjectID),
		DiagramName: req.Title,
		PUMLContent: req.Content,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return diagram, nil
}

// GetProjectPUMLs 获取项目PUML图表列表
func (s *PUMLService) GetProjectPUMLs(userID uuid.UUID, projectID string) ([]*model.PUMLDiagram, error) {
	// 这里应该调用repository层获取数据，暂时返回空列表
	return []*model.PUMLDiagram{}, nil
}

// UpdatePUMLDiagram 更新PUML图表
func (s *PUMLService) UpdatePUMLDiagram(userID uuid.UUID, pumlID string, req *model.UpdatePUMLRequest) (*model.PUMLDiagram, error) {
	// 验证PUML语法
	validation := s.ValidatePUMLString(req.Content)
	if !validation.IsValid {
		return nil, fmt.Errorf("PUML语法错误: %v", validation.Errors)
	}

	// 更新PUML图表（这里应该调用repository层，暂时返回模拟数据）
	diagram := &model.PUMLDiagram{
		DiagramID:   uuid.MustParse(pumlID),
		DiagramName: req.Title,
		PUMLContent: req.Content,
		UpdatedAt:   time.Now(),
	}

	return diagram, nil
}

// DeletePUML 删除PUML图表
func (s *PUMLService) DeletePUML(userID uuid.UUID, pumlID string) error {
	// 这里应该调用repository层删除数据
	return nil
}

// RenderPUMLImage 渲染PUML图片
func (s *PUMLService) RenderPUMLImage(req *model.RenderPUMLRequest) (*RenderResult, error) {
	options := &RenderOptions{
		Format:     req.Format,
		UseCache:   true,
		ServerMode: true,
	}
	if options.Format == "" {
		options.Format = "png"
	}

	return s.RenderPUML(req.Content, options)
}

// RenderPUMLOnlineFromRequest 在线渲染PUML（适配controller接口）
func (s *PUMLService) RenderPUMLOnlineFromRequest(req *model.RenderPUMLRequest) (string, error) {
	return s.RenderPUMLOnline(req.Content)
}

// GenerateImage 生成图片
func (s *PUMLService) GenerateImage(req *model.GenerateImageRequest) (*RenderResult, error) {
	options := &RenderOptions{
		Format:     req.Format,
		UseCache:   true,
		ServerMode: true,
	}
	if options.Format == "" {
		options.Format = "png"
	}

	return s.RenderPUML(req.Content, options)
}

// ValidatePUMLString 验证PUML语法（重命名避免方法签名冲突）
func (s *PUMLService) ValidatePUMLString(pumlCode string) *ValidationResult {
	return s.ValidatePUML(pumlCode)
}

// ValidatePUMLFromRequest 验证PUML语法（适配controller接口）
func (s *PUMLService) ValidatePUMLFromRequest(req *model.ValidatePUMLRequest) (*ValidationResult, error) {
	result := s.ValidatePUML(req.Content)
	return result, nil
}

// PreviewPUML 预览PUML
func (s *PUMLService) PreviewPUML(req *model.PreviewPUMLRequest) (*RenderResult, error) {
	options := &RenderOptions{
		Format:     "svg",
		UseCache:   false, // 预览不使用缓存
		ServerMode: true,
	}

	return s.RenderPUML(req.Content, options)
}

// ExportPUML 导出PUML
func (s *PUMLService) ExportPUML(userID uuid.UUID, req *model.ExportPUMLRequest) (interface{}, error) {
	// 这里应该根据format类型导出不同格式的文件
	// 暂时返回基本响应
	return map[string]interface{}{
		"exported_count": len(req.PUMLIDs),
		"format":        req.Format,
		"user_id":       userID,
	}, nil
}

// GetPUMLStats 获取PUML统计信息
func (s *PUMLService) GetPUMLStats(userID uuid.UUID) (interface{}, error) {
	// 这里应该调用repository层获取统计数据
	return map[string]interface{}{
		"total_diagrams": 0,
		"user_id":       userID,
		"cache_stats":   s.GetCacheStats(),
	}, nil
}

// ClearPUMLCache 清空PUML缓存
func (s *PUMLService) ClearPUMLCache(userID uuid.UUID) error {
	s.ClearCache()
	return nil
} 