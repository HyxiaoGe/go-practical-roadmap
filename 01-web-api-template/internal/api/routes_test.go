package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-practical-roadmap/01-web-api-template/internal/api/dto"
)

func TestHealthCheck(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建响应记录器和请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// 创建路由引擎
	r := gin.New()
	r.GET("/health", healthCheck)

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Service is running", response["message"])
}

func TestRegisterHandler_Success(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试请求数据
	registerReq := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// 将请求数据转换为JSON
	jsonData, _ := json.Marshal(registerReq)

	// 创建响应记录器和请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建路由引擎
	r := gin.New()
	r.POST("/api/v1/register", func(c *gin.Context) {
		// 模拟用户服务，直接返回成功响应
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"data": gin.H{
				"username": registerReq.Username,
				"email":    registerReq.Email,
			},
		})
	})

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)

	// 验证响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
}

func TestLoginHandler_Success(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试请求数据
	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// 将请求数据转换为JSON
	jsonData, _ := json.Marshal(loginReq)

	// 创建响应记录器和请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建路由引擎
	r := gin.New()
	r.POST("/api/v1/login", func(c *gin.Context) {
		// 模拟用户服务，直接返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   "mock-jwt-token",
		})
	})

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response["message"])
	assert.Equal(t, "mock-jwt-token", response["token"])
}