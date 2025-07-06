package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"ai-dev-platform/internal/utils"

	"github.com/google/uuid"
)

// Handlers API处理器集合
type Handlers struct {
	userService    service.UserService
	projectService service.ProjectService
}

// NewHandlers 创建API处理器
func NewHandlers(userService service.UserService, projectService service.ProjectService) *Handlers {
	return &Handlers{
		userService:    userService,
		projectService: projectService,
	}
}

// Health 健康检查
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "ai-dev-platform",
		"version":   "1.0.0",
	}
	utils.WriteSuccessResponse(w, data, "服务运行正常")
}

// RegisterUser 用户注册
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	user, err := h.userService.RegisterUser(&req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteCreatedResponse(w, user, "用户注册成功")
}

// LoginUser 用户登录
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	response, err := h.userService.LoginUser(&req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, response, "登录成功")
}

// GetCurrentUser 获取当前用户信息
func (h *Handlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())
	utils.WriteSuccessResponse(w, user, "获取用户信息成功")
}

// UpdateCurrentUser 更新当前用户信息
func (h *Handlers) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	var req service.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	updatedUser, err := h.userService.UpdateUser(user.UserID, &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, updatedUser, "用户信息更新成功")
}

// CreateProject 创建项目
func (h *Handlers) CreateProject(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	var req model.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	project, err := h.projectService.CreateProject(user.UserID, &req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteCreatedResponse(w, project, "项目创建成功")
}

// GetUserProjects 获取用户项目列表
func (h *Handlers) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	// 获取分页参数
	page, pageSize := utils.GetPaginationParams(r)

	projects, pagination, err := h.projectService.GetUserProjects(user.UserID, page, pageSize)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WritePaginatedResponse(w, projects, pagination, "获取项目列表成功")
}

// GetProject 获取项目详情
func (h *Handlers) GetProject(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	// 从URL路径获取项目ID
	projectIDStr := extractIDFromPath(r.URL.Path, "/api/projects/")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID格式")
		return
	}

	project, err := h.projectService.GetProject(projectID, user.UserID)
	if err != nil {
		if err.Error() == "项目不存在" || err.Error() == "无权访问此项目" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteSuccessResponse(w, project, "获取项目详情成功")
}

// UpdateProject 更新项目
func (h *Handlers) UpdateProject(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	// 从URL路径获取项目ID
	projectIDStr := extractIDFromPath(r.URL.Path, "/api/projects/")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID格式")
		return
	}

	var req service.ProjectUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	project, err := h.projectService.UpdateProject(projectID, user.UserID, &req)
	if err != nil {
		if err.Error() == "项目不存在" || err.Error() == "无权修改此项目" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.WriteSuccessResponse(w, project, "项目更新成功")
}

// DeleteProject 删除项目
func (h *Handlers) DeleteProject(w http.ResponseWriter, r *http.Request) {
	user := MustGetUserFromContext(r.Context())

	// 从URL路径获取项目ID
	projectIDStr := extractIDFromPath(r.URL.Path, "/api/projects/")
	if projectIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "无效的项目ID格式")
		return
	}

	err = h.projectService.DeleteProject(projectID, user.UserID)
	if err != nil {
		if err.Error() == "项目不存在" || err.Error() == "无权删除此项目" {
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteSuccessResponse(w, nil, "项目删除成功")
}

// extractIDFromPath 从URL路径中提取ID
func extractIDFromPath(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	
	// 移除前缀
	id := path[len(prefix):]
	
	// 如果路径后面还有其他部分，只取第一部分
	if slashIndex := strings.IndexByte(id, '/'); slashIndex != -1 {
		id = id[:slashIndex]
	}
	
	return id
} 