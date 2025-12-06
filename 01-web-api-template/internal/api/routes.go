package api

import (
	"net/http"

	"github.com/yourname/01-web-api-template/internal/middleware"
	"github.com/yourname/01-web-api-template/pkg/logger"
)

// SetupRoutes 设置路由
func SetupRoutes() http.Handler {
	// 创建一个简单的多路复用器
	mux := http.NewServeMux()

	// 公开路由
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/api/v1/register", registerHandler)
	mux.HandleFunc("/api/v1/login", loginHandler)

	// 受保护的路由
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/v1/profile", profileHandler)

	// 应用JWT认证中间件
	protectedHandler := middleware.JWTAuthMiddleware(protectedMux)

	// 注册受保护的路由
	mux.Handle("/api/v1/profile", protectedHandler)

	return mux
}

// healthCheck 健康检查端点
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "Service is running"}`))
	logger.Info("Health check endpoint called")
}

// registerHandler 注册端点
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User registered successfully"}`))
	logger.Info("User registration endpoint called")
}

// loginHandler 登录端点
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login successful", "token": "mock-jwt-token"}`))
	logger.Info("User login endpoint called")
}

// profileHandler 用户信息端点
func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User profile retrieved successfully", "user": {"id": 1, "username": "testuser"}}`))
	logger.Info("User profile endpoint called")
}