package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-practical-roadmap/01-web-api-template/internal/api/dto"
	"go-practical-roadmap/01-web-api-template/internal/middleware"
	"go-practical-roadmap/01-web-api-template/internal/service"
	"go-practical-roadmap/01-web-api-template/pkg/logger"
	"go.uber.org/zap"
)

// SetupRoutes 设置Gin路由
func SetupRoutes(userService service.UserService) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 添加日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 添加自定义中间件
	r.Use(middleware.RequestTracerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 公开路由
	r.GET("/health", healthCheck)
	r.POST("/api/v1/register", func(c *gin.Context) {
		registerHandler(c, userService)
	})
	r.POST("/api/v1/login", func(c *gin.Context) {
		loginHandler(c, userService)
	})

	// 受保护的路由组
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware)
	{
		authorized.GET("/api/v1/profile", profileHandler)
	}

	return r
}

// healthCheck 健康检查端点
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
	logger.Info("Health check endpoint called")
}

// registerHandler 注册端点
func registerHandler(c *gin.Context, userService service.UserService) {
	var req dto.RegisterRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用用户服务注册
	user, err := userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    user,
	})
	logger.Info("User registration endpoint called",
		zap.String("username", req.Username),
		zap.String("email", req.Email))
}

// loginHandler 登录端点
func loginHandler(c *gin.Context, userService service.UserService) {
	var req dto.LoginRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用用户服务登录
	token, err := userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
	logger.Info("User login endpoint called",
		zap.String("username", req.Username))
}

// profileHandler 用户信息端点
func profileHandler(c *gin.Context) {
	// 模拟用户数据
	user := dto.UserProfileResponse{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile retrieved successfully",
		"data":    user,
	})
	logger.Info("User profile endpoint called",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username))
}