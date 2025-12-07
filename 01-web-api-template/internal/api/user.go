package api

import (
	"encoding/json"
	"net/http"

	"go-practical-roadmap/01-web-api-template/internal/api/dto"
	"go-practical-roadmap/01-web-api-template/internal/service"
	"go-practical-roadmap/01-web-api-template/pkg/logger"
	"go.uber.org/zap"
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
	var req dto.RegisterRequest

	// 解析请求参数
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 调用UserService.Register
	user, err := c.userService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"data":    user,
	})

	logger.Info("User registration endpoint called",
		zap.String("username", req.Username),
		zap.String("email", req.Email))
}

// Login 用户登录
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// 解析请求参数
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 调用UserService.Login
	token, err := c.userService.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 返回JWT令牌
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})

	logger.Info("User login endpoint called", zap.String("username", req.Username))
}

// GetProfile 获取用户信息
func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取用户信息逻辑
	// 从JWT中获取用户信息
	// 返回用户详情

	logger.Info("Get user profile endpoint called")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Get user profile endpoint",
	})
}