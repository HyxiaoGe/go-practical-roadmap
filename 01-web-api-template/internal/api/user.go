package api

import (
	"net/http"

	"go-practical-roadmap/01-web-api-template/internal/service"
	"go-practical-roadmap/01-web-api-template/pkg/logger"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// Register 注册用户
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现用户注册逻辑
	// 解析请求参数
	// 调用UserService.Register
	// 返回响应

	logger.Info("User registration endpoint called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registration endpoint"))
}

// Login 用户登录
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现用户登录逻辑
	// 解析请求参数
	// 调用UserService.Login
	// 返回JWT令牌

	logger.Info("User login endpoint called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User login endpoint"))
}

// GetProfile 获取用户信息
func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取用户信息逻辑
	// 从JWT中获取用户信息
	// 返回用户详情

	logger.Info("Get user profile endpoint called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get user profile endpoint"))
}