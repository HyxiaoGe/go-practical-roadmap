package api

import (
	"github.com/yourname/02-concurrency-worker/internal/pkg/websocket"
	"github.com/yourname/02-concurrency-worker/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(router *gin.Engine, taskManager *service.TaskManager, hub *websocket.Hub) {
	// 创建处理器
	handlers := NewHandlers(taskManager)
	wsHandler := NewWebSocketHandler(hub)

	// 健康检查端点
	router.GET("/health", handlers.HealthCheck)

	// API v1路由组
	v1 := router.Group("/api/v1")
	{
		// 任务相关路由
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", handlers.SubmitTask)           // 提交任务
			tasks.GET("/:id", handlers.GetTask)          // 获取任务详情
			tasks.GET("", handlers.GetAllTasks)          // 获取所有任务
			tasks.DELETE("/:id", handlers.CancelTask)    // 取消任务
			tasks.GET("/status/stats", handlers.GetTaskStats) // 获取任务统计
		}

		// WebSocket路由
		v1.GET("/ws/tasks", wsHandler.HandleWebSocket)
	}
}