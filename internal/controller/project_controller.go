package controller

import (
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectController 项目控制器
type ProjectController struct {
	projectService service.ProjectService
}

// NewProjectController 创建项目控制器
func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// CreateProject 创建项目
func (pc *ProjectController) CreateProject(c *gin.Context) {
	log.InfofId(c, "CreateProject: 开始处理项目创建请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "CreateProject: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "CreateProject: 用户认证成功，用户ID: %s", user.UserID.String())

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "CreateProject: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求数据",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数校验
	if req.ProjectName == "" {
		log.WarnfId(c, "CreateProject: 项目名称不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目名称不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if len(req.ProjectName) > 100 {
		log.WarnfId(c, "CreateProject: 项目名称长度不能超过100个字符")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目名称长度不能超过100个字符",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if req.ProjectType == "" {
		log.WarnfId(c, "CreateProject: 项目类型不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目类型不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if len(req.Description) > 500 {
		log.WarnfId(c, "CreateProject: 项目描述长度不能超过500个字符")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目描述长度不能超过500个字符",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "CreateProject: 项目创建请求数据: name=%s, type=%s", req.ProjectName, req.ProjectType)

	project, err := pc.projectService.CreateProject(user.UserID, &req)
	if err != nil {
		log.ErrorfId(c, "CreateProject: 项目创建失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "CreateProject: 项目创建成功，项目ID: %s", project.ProjectID.String())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    project,
		"message": "项目创建成功",
		"code":    http.StatusCreated,
	})
}

// GetUserProjects 获取用户项目列表
func (pc *ProjectController) GetUserProjects(c *gin.Context) {
	log.InfofId(c, "GetUserProjects: 开始获取用户项目列表")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetUserProjects: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "GetUserProjects: 用户认证成功，用户ID: %s", user.UserID.String())

	// 获取分页参数并校验
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			if parsed > 1000 { // 防止过大的页码
				page = 1000
			} else {
				page = parsed
			}
		} else {
			log.WarnfId(c, "GetUserProjects: 无效的页码参数: %s", p)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "无效的页码参数",
				"code":    http.StatusBadRequest,
			})
			return
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			if parsed > 100 {
				pageSize = 100
			} else {
				pageSize = parsed
			}
		} else {
			log.WarnfId(c, "GetUserProjects: 无效的页面大小参数: %s", ps)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "无效的页面大小参数",
				"code":    http.StatusBadRequest,
			})
			return
		}
	}

	log.InfofId(c, "GetUserProjects: 分页参数 page=%d, pageSize=%d", page, pageSize)

	projects, pagination, err := pc.projectService.GetUserProjects(user.UserID, page, pageSize)
	if err != nil {
		log.ErrorfId(c, "GetUserProjects: 获取项目列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	log.InfofId(c, "GetUserProjects: 成功获取项目列表，项目数量: %d", len(projects))

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       projects,
		"pagination": pagination,
		"message":    "获取项目列表成功",
		"code":       http.StatusOK,
	})
}

// GetProject 获取项目详情
func (pc *ProjectController) GetProject(c *gin.Context) {
	log.InfofId(c, "GetProject: 开始获取项目详情")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetProject: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectIDStr := c.Param("id")
	if projectIDStr == "" {
		log.WarnfId(c, "GetProject: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 解析UUID
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		log.WarnfId(c, "GetProject: 无效的项目ID格式: %s", projectIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "GetProject: 获取项目详情，项目ID: %s, 用户ID: %s", projectID.String(), user.UserID.String())

	// 调用服务层获取项目详情
	project, err := pc.projectService.GetProject(projectID, user.UserID)
	if err != nil {
		log.ErrorfId(c, "GetProject: 获取项目详情失败: %v", err)
		
		// 根据错误类型返回不同的HTTP状态码
		statusCode := http.StatusInternalServerError
		if err.Error() == "无权访问此项目" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "项目不存在" || err.Error() == "获取项目失败: 项目不存在" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    statusCode,
		})
		return
	}

	log.InfofId(c, "GetProject: 成功获取项目详情，项目名称: %s", project.ProjectName)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
		"message": "获取项目详情成功",
		"code":    http.StatusOK,
	})
}

// UpdateProject 更新项目
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	log.InfofId(c, "UpdateProject: 开始处理项目更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdateProject: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectIDStr := c.Param("id")
	if projectIDStr == "" {
		log.WarnfId(c, "UpdateProject: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 解析UUID
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		log.WarnfId(c, "UpdateProject: 无效的项目ID格式: %s", projectIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	var updates service.ProjectUpdateRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.ErrorfId(c, "UpdateProject: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求数据",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdateProject: 更新项目，项目ID: %s, 用户ID: %s", projectID.String(), user.UserID.String())

	// 调用服务层更新项目
	project, err := pc.projectService.UpdateProject(projectID, user.UserID, &updates)
	if err != nil {
		log.ErrorfId(c, "UpdateProject: 更新项目失败: %v", err)
		
		// 根据错误类型返回不同的HTTP状态码
		statusCode := http.StatusInternalServerError
		if err.Error() == "无权访问此项目" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "项目不存在" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    statusCode,
		})
		return
	}

	log.InfofId(c, "UpdateProject: 项目更新成功，项目名称: %s", project.ProjectName)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
		"message": "项目更新成功",
		"code":    http.StatusOK,
	})
}

// DeleteProject 删除项目
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	log.InfofId(c, "DeleteProject: 开始处理项目删除请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "DeleteProject: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	projectIDStr := c.Param("id")
	if projectIDStr == "" {
		log.WarnfId(c, "DeleteProject: 项目ID不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "项目ID不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 解析UUID
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		log.WarnfId(c, "DeleteProject: 无效的项目ID格式: %s", projectIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的项目ID格式",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "DeleteProject: 删除项目，项目ID: %s, 用户ID: %s", projectID.String(), user.UserID.String())

	// 调用服务层删除项目
	err = pc.projectService.DeleteProject(projectID, user.UserID)
	if err != nil {
		log.ErrorfId(c, "DeleteProject: 删除项目失败: %v", err)
		
		// 根据错误类型返回不同的HTTP状态码
		statusCode := http.StatusInternalServerError
		if err.Error() == "无权访问此项目" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "项目不存在" {
			statusCode = http.StatusNotFound
		}
		
		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    statusCode,
		})
		return
	}

	log.InfofId(c, "DeleteProject: 项目删除成功，项目ID: %s", projectID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "项目删除成功",
		"code":    http.StatusOK,
	})
}
