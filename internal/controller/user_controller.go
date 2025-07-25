package controller

import (
	"ai-dev-platform/internal/log"
	"ai-dev-platform/internal/model"
	"ai-dev-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterUser 用户注册
func (uc *UserController) RegisterUser(c *gin.Context) {
	log.InfofId(c, "RegisterUser: 开始处理用户注册请求")

	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "RegisterUser: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求数据",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数校验
	if req.Email == "" {
		log.WarnfId(c, "RegisterUser: 邮箱不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "邮箱不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if req.Password == "" {
		log.WarnfId(c, "RegisterUser: 密码不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "密码不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if len(req.Password) < 6 {
		log.WarnfId(c, "RegisterUser: 密码长度不能少于6位")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "密码长度不能少于6位",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "RegisterUser: 用户注册请求，邮箱: %s", req.Email)

	user, err := uc.userService.RegisterUser(&req)
	if err != nil {
		log.ErrorfId(c, "RegisterUser: 用户注册失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "RegisterUser: 用户注册成功，用户ID: %s", user.UserID.String())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
		"message": "用户注册成功",
		"code":    http.StatusCreated,
	})
}

// LoginUser 用户登录
func (uc *UserController) LoginUser(c *gin.Context) {
	log.InfofId(c, "LoginUser: 开始处理用户登录请求")

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "LoginUser: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求数据",
			"code":    http.StatusBadRequest,
		})
		return
	}

	// 参数校验
	if req.Email == "" {
		log.WarnfId(c, "LoginUser: 邮箱不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "邮箱不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	if req.Password == "" {
		log.WarnfId(c, "LoginUser: 密码不能为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "密码不能为空",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "LoginUser: 用户登录请求，邮箱: %s", req.Email)

	response, err := uc.userService.LoginUser(&req)
	if err != nil {
		log.ErrorfId(c, "LoginUser: 用户登录失败: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "LoginUser: 用户登录成功，用户ID: %s", response.User.UserID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "登录成功",
		"code":    http.StatusOK,
	})
}

// ValidateToken 验证Token有效性
func (uc *UserController) ValidateToken(c *gin.Context) {
	log.InfofId(c, "ValidateToken: 开始处理Token验证请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "ValidateToken: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "ValidateToken: Token验证成功，用户ID: %s", user.UserID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "Token验证成功",
		"code":    http.StatusOK,
	})
}

// GetCurrentUser 获取当前用户信息
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	log.InfofId(c, "GetCurrentUser: 开始获取当前用户信息")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "GetCurrentUser: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	log.InfofId(c, "GetCurrentUser: 获取用户信息成功，用户ID: %s", user.UserID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"message": "获取用户信息成功",
		"code":    http.StatusOK,
	})
}

// UpdateCurrentUser 更新当前用户信息
func (uc *UserController) UpdateCurrentUser(c *gin.Context) {
	log.InfofId(c, "UpdateCurrentUser: 开始处理用户信息更新请求")

	user, ok := ginUserFromContext(c)
	if !ok {
		log.WarnfId(c, "UpdateCurrentUser: 认证信息无效")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "认证信息无效",
			"code":    http.StatusUnauthorized,
		})
		return
	}

	var req service.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.ErrorfId(c, "UpdateCurrentUser: 请求数据解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的请求数据",
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdateCurrentUser: 用户信息更新请求，用户ID: %s", user.UserID.String())

	updatedUser, err := uc.userService.UpdateUser(user.UserID, &req)
	if err != nil {
		log.ErrorfId(c, "UpdateCurrentUser: 用户信息更新失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	log.InfofId(c, "UpdateCurrentUser: 用户信息更新成功，用户ID: %s", updatedUser.UserID.String())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedUser,
		"message": "用户信息更新成功",
		"code":    http.StatusOK,
	})
}

// ginUserFromContext 从gin上下文获取用户信息
func ginUserFromContext(c *gin.Context) (*model.User, bool) {
	value, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := value.(*model.User)
	return user, ok
}