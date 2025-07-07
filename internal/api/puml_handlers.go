package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"ai-dev-platform/internal/service"

	"github.com/google/uuid"
)

// PUMLHandlers PUML相关的API处理器
type PUMLHandlers struct {
	pumlService *service.PUMLService
}

// NewPUMLHandlers 创建新的PUML处理器
func NewPUMLHandlers(pumlService *service.PUMLService) *PUMLHandlers {
	return &PUMLHandlers{
		pumlService: pumlService,
	}
}

// RenderPUML 渲染PUML图像
// POST /api/puml/render
func (h *PUMLHandlers) RenderPUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLCode string `json:"puml_code"`
		Format   string `json:"format"`   // png, svg, txt
		Theme    string `json:"theme"`    // 主题
		DPI      int    `json:"dpi"`      // DPI设置
		UseCache bool   `json:"use_cache"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PUMLCode == "" {
		http.Error(w, "PUML代码不能为空", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.Format == "" {
		req.Format = "png"
	}

	options := &service.RenderOptions{
		Format:     req.Format,
		Theme:      req.Theme,
		DPI:        req.DPI,
		UseCache:   req.UseCache,
		ServerMode: true,
	}

	result, err := h.pumlService.RenderPUML(req.PUMLCode, options)
	if err != nil {
		http.Error(w, "渲染失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 根据格式设置响应头
	switch req.Format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case "txt":
		w.Header().Set("Content-Type", "text/plain")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(result.ImageData)))
	w.Write(result.ImageData)
}

// ValidatePUML 验证PUML语法
// POST /api/puml/validate
func (h *PUMLHandlers) ValidatePUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLCode string `json:"puml_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PUMLCode == "" {
		http.Error(w, "PUML代码不能为空", http.StatusBadRequest)
		return
	}

	result := h.pumlService.ValidatePUML(req.PUMLCode)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// PreviewPUML 预览PUML图表（返回图片URL或base64）
// POST /api/puml/preview
func (h *PUMLHandlers) PreviewPUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLCode   string `json:"puml_code"`
		Format     string `json:"format"`      // png, svg
		ReturnType string `json:"return_type"` // url, base64
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PUMLCode == "" {
		http.Error(w, "PUML代码不能为空", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.Format == "" {
		req.Format = "png"
	}
	if req.ReturnType == "" {
		req.ReturnType = "url"
	}

	options := &service.RenderOptions{
		Format:     req.Format,
		UseCache:   true,
		ServerMode: true,
	}

	result, err := h.pumlService.RenderPUML(req.PUMLCode, options)
	if err != nil {
		http.Error(w, "预览失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"format":      result.Format,
		"rendered_at": result.RenderedAt,
		"cache_key":   result.CacheKey,
	}

	if req.ReturnType == "base64" {
		// 返回base64编码的图片
		response["data"] = fmt.Sprintf("data:image/%s;base64,%s", 
			result.Format, 
			encodeBase64(result.ImageData))
	} else {
		// 返回图片URL
		response["url"] = result.URL
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExportPUML 导出PUML图表
// POST /api/puml/export
func (h *PUMLHandlers) ExportPUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLCode string `json:"puml_code"`
		Format   string `json:"format"`   // png, svg, pdf
		Filename string `json:"filename"` // 导出文件名
		DPI      int    `json:"dpi"`      // DPI设置
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "解析请求失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.PUMLCode == "" {
		http.Error(w, "PUML代码不能为空", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.Format == "" {
		req.Format = "png"
	}
	if req.Filename == "" {
		req.Filename = fmt.Sprintf("diagram_%s.%s", uuid.New().String()[:8], req.Format)
	}
	if req.DPI == 0 {
		req.DPI = 300 // 高质量导出
	}

	options := &service.RenderOptions{
		Format:     req.Format,
		DPI:        req.DPI,
		UseCache:   false, // 导出时不使用缓存，确保最新结果
		ServerMode: true,
	}

	result, err := h.pumlService.RenderPUML(req.PUMLCode, options)
	if err != nil {
		http.Error(w, "导出失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置下载响应头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", req.Filename))
	
	switch req.Format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(result.ImageData)))
	w.Write(result.ImageData)
}

// GetPUMLStats 获取PUML服务统计信息
// GET /api/puml/stats
func (h *PUMLHandlers) GetPUMLStats(w http.ResponseWriter, r *http.Request) {
	stats := h.pumlService.GetCacheStats()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// ClearPUMLCache 清空PUML缓存
// POST /api/puml/cache/clear
func (h *PUMLHandlers) ClearPUMLCache(w http.ResponseWriter, r *http.Request) {
	h.pumlService.ClearCache()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "缓存已清空",
	})
}

// encodeBase64 将字节数组编码为base64字符串
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
} 