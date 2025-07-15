package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
)

// PUMLHandlers PUML相关的API处理器
type PUMLHandlers struct {
	pumlService *service.PUMLService
	aiService   *service.AIService
}

// NewPUMLHandlers 创建新的PUML处理器
func NewPUMLHandlers(pumlService *service.PUMLService, aiService *service.AIService) *PUMLHandlers {
	return &PUMLHandlers{
		pumlService: pumlService,
		aiService:   aiService,
	}
}

// RenderPUML 渲染PUML图像
// POST /api/puml/render
func (h *PUMLHandlers) RenderPUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLCode string `json:"puml_code"`
		Format   string `json:"format"` // png, svg, txt
		Theme    string `json:"theme"`  // 主题
		DPI      int    `json:"dpi"`    // DPI设置
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

	result := map[string]interface{}{
		"success": true,
		"message": "PUML缓存已清空",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ===== 项目PUML管理相关处理器 =====

// GetProjectPUMLs 获取项目PUML图表列表
// GET /api/puml/project/{projectId}
func (h *PUMLHandlers) GetProjectPUMLs(w http.ResponseWriter, r *http.Request) {
	// 从URL路径提取项目ID
	projectIDStr := r.PathValue("projectId")
	if projectIDStr == "" {
		http.Error(w, "缺少项目ID", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "无效的项目ID", http.StatusBadRequest)
		return
	}

	// 获取查询参数
	stageStr := r.URL.Query().Get("stage")
	var stage *int
	if stageStr != "" {
		if s, err := strconv.Atoi(stageStr); err == nil {
			stage = &s
		}
	}

	// 调用AI服务获取PUML图表列表
	diagrams, err := h.aiService.GetPUMLDiagramsByProject(projectID)
	if err != nil {
		// 添加详细的错误日志
		fmt.Printf("GetPUMLDiagramsByProject error: %v\n", err)
		http.Error(w, "获取PUML图表列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果指定了阶段，进行过滤
	if stage != nil {
		filteredDiagrams := make([]*model.PUMLDiagram, 0)
		for _, diagram := range diagrams {
			if diagram.Stage == *stage {
				filteredDiagrams = append(filteredDiagrams, diagram)
			}
		}
		diagrams = filteredDiagrams
	}

	result := map[string]interface{}{
		"success": true,
		"data":    diagrams,
		"message": "获取PUML图表列表成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CreatePUML 创建PUML图表
// POST /api/puml/create
func (h *PUMLHandlers) CreatePUML(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID   string `json:"project_id"`
		Stage       int    `json:"stage"`
		DiagramType string `json:"diagram_type"`
		DiagramName string `json:"diagram_name"`
		PUMLContent string `json:"puml_content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.ProjectID == "" || req.DiagramName == "" || req.PUMLContent == "" {
		http.Error(w, "项目ID、图表名称和PUML内容不能为空", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		http.Error(w, "无效的项目ID", http.StatusBadRequest)
		return
	}

	// 创建PUML图表
	now := time.Now()
	diagram := &model.PUMLDiagram{
		DiagramID:   uuid.New(),
		ProjectID:   projectID,
		DiagramType: req.DiagramType,
		DiagramName: req.DiagramName,
		PUMLContent: req.PUMLContent,
		Stage:       req.Stage,
		Version:     1,
		IsValidated: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 调用AI服务创建PUML图表
	err = h.aiService.CreatePUML(diagram)
	if err != nil {
		http.Error(w, "创建PUML图表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"success": true,
		"data":    diagram,
		"message": "PUML图表创建成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// UpdatePUMLDiagram 更新PUML图表
// PUT /api/puml/{pumlId}
func (h *PUMLHandlers) UpdatePUMLDiagram(w http.ResponseWriter, r *http.Request) {
	// 从URL路径提取PUML ID
	pumlIDStr := r.PathValue("pumlId")
	if pumlIDStr == "" {
		http.Error(w, "缺少PUML图表ID", http.StatusBadRequest)
		return
	}

	pumlID, err := uuid.Parse(pumlIDStr)
	if err != nil {
		http.Error(w, "无效的PUML图表ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DiagramName *string `json:"diagram_name"`
		PUMLContent *string `json:"puml_content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 构建更新请求
	updateReq := &model.UpdatePUMLRequest{}
	if req.DiagramName != nil {
		updateReq.Title = *req.DiagramName
	}
	if req.PUMLContent != nil {
		updateReq.Content = *req.PUMLContent
	}

	// 调用AI服务更新PUML图表
	err = h.aiService.UpdatePUMLDiagram(pumlID, updateReq)
	if err != nil {
		http.Error(w, "更新PUML图表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"success": true,
		"message": "PUML图表更新成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeletePUML 删除PUML图表
// DELETE /api/puml/{pumlId}
func (h *PUMLHandlers) DeletePUML(w http.ResponseWriter, r *http.Request) {
	// 从URL路径提取PUML ID
	pumlIDStr := r.PathValue("pumlId")
	if pumlIDStr == "" {
		http.Error(w, "缺少PUML图表ID", http.StatusBadRequest)
		return
	}

	pumlID, err := uuid.Parse(pumlIDStr)
	if err != nil {
		http.Error(w, "无效的PUML图表ID", http.StatusBadRequest)
		return
	}

	// 调用AI服务删除PUML图表
	err = h.aiService.DeletePUMLDiagram(pumlID)
	if err != nil {
		http.Error(w, "删除PUML图表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"success": true,
		"message": "PUML图表删除成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// encodeBase64 将字节数组编码为base64字符串
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// RenderPUMLImage 渲染PUML为图片
// POST /api/puml/render
func (h *PUMLHandlers) RenderPUMLImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PUMLContent string `json:"puml_content"`
		Format      string `json:"format"` // png, svg
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	if req.PUMLContent == "" {
		http.Error(w, "PUML内容不能为空", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "png"
	}

	// 调用PUML服务渲染图片
	renderResult, err := h.pumlService.RenderPUML(req.PUMLContent, &service.RenderOptions{
		Format:     req.Format,
		ServerMode: true, // 确保使用在线服务器模式
		UseCache:   true,
	})
	if err != nil {
		http.Error(w, "渲染PUML图片失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"url":    renderResult.URL,
			"format": req.Format,
		},
		"message": "PUML图片渲染成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GenerateImage 生成PUML图片 (别名方法，兼容前端API)
// POST /api/puml/generate-image
func (h *PUMLHandlers) GenerateImage(w http.ResponseWriter, r *http.Request) {
	h.RenderPUMLImage(w, r)
}

// RenderPUMLOnlineHandler 在线渲染PUML图表（通过POST原始代码）
func (h *PUMLHandlers) RenderPUMLOnlineHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "PUML代码不能为空")
		return
	}

	// 调用服务
	svgData, err := h.pumlService.RenderPUMLOnline(req.Code)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "渲染PUML失败: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"message":   "渲染成功",
		"imageData": svgData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidatePUMLHandler 验证PUML语法
func (h *PUMLHandlers) ValidatePUMLHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "PUML代码不能为空")
		return
	}

	// 调用服务
	result := h.pumlService.ValidatePUML(req.Code)

	response := map[string]interface{}{
		"success": true,
		"message": "验证成功",
		"data":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
