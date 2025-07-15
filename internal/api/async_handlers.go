package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
)

// AsyncHandlers 异步任务相关的HTTP处理器
type AsyncHandlers struct {
	asyncTaskService *service.AsyncTaskService
	aiService        *service.AIService
}

// NewAsyncHandlers 创建异步任务处理器
func NewAsyncHandlers(asyncTaskService *service.AsyncTaskService, aiService *service.AIService) *AsyncHandlers {
	return &AsyncHandlers{
		asyncTaskService: asyncTaskService,
		aiService:        aiService,
	}
}

// StartStageDocumentGeneration 启动阶段文档生成任务
// POST /api/async/stage-documents
func (h *AsyncHandlers) StartStageDocumentGeneration(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户信息
	user := MustGetUserFromContext(r.Context())

	var req struct {
		ProjectID string `json:"project_id" validate:"required"`
		Stage     int    `json:"stage" validate:"required,min=1,max=3"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.ProjectID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "项目ID不能为空")
		return
	}

	if req.Stage < 1 || req.Stage > 3 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "阶段必须是1、2或3")
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 启动异步任务
	response, err := h.asyncTaskService.StartStageDocumentGeneration(projectID, user.UserID, req.Stage)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("启动文档生成任务失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, response, "文档生成任务已启动")
}

// GetTaskStatus 获取任务状态
// GET /api/async/tasks/{taskId}/status
func (h *AsyncHandlers) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.PathValue("taskId")
	if taskIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少任务ID")
		return
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的任务ID")
		return
	}

	// 获取任务状态
	response, err := h.asyncTaskService.GetTaskStatus(taskID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取任务状态失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, response, "获取任务状态成功")
}

// GetStageProgress 获取项目阶段进度
// GET /api/async/projects/{projectId}/progress
func (h *AsyncHandlers) GetStageProgress(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectId")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取阶段进度
	response, err := h.asyncTaskService.GetStageProgress(projectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取阶段进度失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, response, "获取阶段进度成功")
}

// GetStageDocuments 获取阶段文档列表
// GET /api/async/projects/{projectId}/stages/{stage}/documents
func (h *AsyncHandlers) GetStageDocuments(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetStageDocuments called with URL: %s\n", r.URL.Path)
	
	projectIDStr := r.PathValue("projectId")
	stageStr := r.PathValue("stage")
	
	fmt.Printf("Extracted projectId: %s, stage: %s\n", projectIDStr, stageStr)
	
	if projectIDStr == "" || stageStr == "" {
		fmt.Printf("Missing parameters - projectId: %s, stage: %s\n", projectIDStr, stageStr)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID或阶段参数")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		fmt.Printf("Invalid projectId format: %s, error: %v\n", projectIDStr, err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	stage, err := strconv.Atoi(stageStr)
	if err != nil || stage < 1 || stage > 3 {
		fmt.Printf("Invalid stage value: %s, error: %v\n", stageStr, err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的阶段参数")
		return
	}

	fmt.Printf("Processing request for projectId: %s, stage: %d\n", projectID, stage)

	// 获取项目的文档和PUML图表
	documents, err := h.aiService.GetDocumentsByProject(projectID)
	if err != nil {
		fmt.Printf("Error getting documents: %v\n", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取文档列表失败: %v", err))
		return
	}

	diagrams, err := h.aiService.GetPUMLDiagramsByProject(projectID)
	if err != nil {
		fmt.Printf("Error getting diagrams: %v\n", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取PUML图表列表失败: %v", err))
		return
	}

	// 过滤指定阶段的文档和图表
	stageDocuments := make([]*model.Document, 0)
	stageDiagrams := make([]*model.PUMLDiagram, 0)

	for _, doc := range documents {
		if doc.Stage == stage {
			stageDocuments = append(stageDocuments, doc)
		}
	}

	for _, diagram := range diagrams {
		if diagram.Stage == stage {
			stageDiagrams = append(stageDiagrams, diagram)
		}
	}

	response := map[string]interface{}{
		"project_id": projectID,
		"stage":      stage,
		"documents":  stageDocuments,
		"diagrams":   stageDiagrams,
	}

	fmt.Printf("Returning %d documents and %d diagrams for stage %d\n", len(stageDocuments), len(stageDiagrams), stage)
	utils.WriteSuccessResponse(w, response, "获取阶段文档成功")
}

// PollTaskStatus 轮询任务状态（支持长轮询）
// GET /api/async/tasks/{taskId}/poll?timeout=30
func (h *AsyncHandlers) PollTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.PathValue("taskId")
	if taskIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少任务ID")
		return
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的任务ID")
		return
	}

	// 解析超时参数（可选）
	timeoutStr := r.URL.Query().Get("timeout")
	timeout := 10 // 默认10秒
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 && t <= 60 {
			timeout = t
		}
	}

	// 简单的轮询实现，实际生产环境可以使用WebSocket或SSE
	response, err := h.asyncTaskService.GetTaskStatus(taskID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("获取任务状态失败: %v", err))
		return
	}

	// 设置缓存控制头
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Poll-Timeout", fmt.Sprintf("%d", timeout))

	utils.WriteSuccessResponse(w, response, "轮询任务状态成功")
}

// GetProjectProgress 获取项目阶段进度
// GET /api/async/projects/{projectId}/progress
func (h *AsyncHandlers) GetProjectProgress(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectId")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "缺少项目ID参数")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 这里可以调用service层获取项目进度信息
	// 暂时返回模拟数据
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"project_id": projectID,
			"stages": map[string]interface{}{
				"1": map[string]interface{}{
					"stage":      1,
					"name":       "需求梳理",
					"completed":  true,
					"documents":  3,
					"puml_diagrams": 2,
				},
				"2": map[string]interface{}{
					"stage":      2,
					"name":       "技术设计",
					"completed":  false,
					"documents":  1,
					"puml_diagrams": 1,
				},
				"3": map[string]interface{}{
					"stage":      3,
					"name":       "实施计划",
					"completed":  false,
					"documents":  0,
					"puml_diagrams": 0,
				},
			},
		},
		"message": "获取项目进度成功",
	}

	utils.WriteSuccessResponse(w, response, "获取项目进度成功")
}

// StartCompleteProjectDocumentGeneration 启动完整项目文档生成任务
// POST /api/async/complete-project-documents
func (h *AsyncHandlers) StartCompleteProjectDocumentGeneration(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户信息
	user := MustGetUserFromContext(r.Context())

	var req struct {
		ProjectID string `json:"project_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 验证请求
	if req.ProjectID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "项目ID不能为空")
		return
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 启动异步任务
	response, err := h.asyncTaskService.StartCompleteProjectDocumentGeneration(projectID, user.UserID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("启动完整项目文档生成任务失败: %v", err))
		return
	}

	utils.WriteSuccessResponse(w, response, "完整项目文档生成任务已启动")
} 