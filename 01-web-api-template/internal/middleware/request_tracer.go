package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-practical-roadmap/01-web-api-template/pkg/logger"
	"go.uber.org/zap"
	"math/rand"
)

// RequestTracerMiddleware 请求追踪中间件
func RequestTracerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()

		// 将请求ID添加到上下文
		c.Set("request_id", requestID)

		// 记录请求开始
		start := time.Now()
		method := c.Request.Method
		url := c.Request.URL.Path

		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("url", url),
		)

		// 处理请求
		c.Next()

		// 记录请求结束
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("url", url),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
		)
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的请求ID生成器
	return strconv.Itoa(rand.Intn(1000000))
}