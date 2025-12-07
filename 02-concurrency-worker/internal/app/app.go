package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/api"
	"github.com/yourname/02-concurrency-worker/internal/config"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"github.com/yourname/02-concurrency-worker/internal/pkg/websocket"
	"github.com/yourname/02-concurrency-worker/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Application 应用程序结构体
type Application struct {
	config      *config.Config
	taskManager *service.TaskManager
	hub         *websocket.Hub
	server      *http.Server
}

// NewApplication 创建新的应用程序实例
func NewApplication() (*Application, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.Output); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 创建任务管理器
	taskManager := service.NewTaskManager(cfg)

	// 创建WebSocket中心
	hub := websocket.NewHub()

	return &Application{
		config:      cfg,
		taskManager: taskManager,
		hub:         hub,
	}, nil
}

// Run 运行应用程序
func (a *Application) Run() error {
	// 启动任务管理器
	a.taskManager.Start()

	// 启动WebSocket中心
	go a.hub.Run()

	// 启动任务状态监听器
	go a.listenTaskUpdates()

	// 设置Gin模式
	if a.config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin路由器
	router := gin.Default()

	// 设置路由
	api.SetupRoutes(router, a.taskManager, a.hub)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port),
		Handler: router,
	}

	logger.Info("Starting server",
		zap.String("address", a.server.Addr),
		zap.String("mode", a.config.Server.Mode))

	// 启动HTTP服务器
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", zap.Error(err))
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// WaitForInterrupt 等待中断信号
func (a *Application) WaitForInterrupt() {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-quit
	logger.Info("Shutdown signal received")

	// 执行优雅关闭
	a.Shutdown()
}

// Shutdown 优雅关闭应用程序
func (a *Application) Shutdown() {
	logger.Info("Shutting down application...")

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}
	}

	// 停止任务管理器
	if err := a.taskManager.Stop(); err != nil {
		logger.Error("Failed to stop task manager", zap.Error(err))
	}

	logger.Info("Application shutdown completed")
}

// listenTaskUpdates 监听任务更新
func (a *Application) listenTaskUpdates() {
	taskUpdateChan := a.taskManager.GetTaskUpdateChan()

	for {
		select {
		case task, ok := <-taskUpdateChan:
			if !ok {
				return
			}

			// 广播任务更新到WebSocket客户端
			a.hub.BroadcastTaskUpdate(task)

		case <-time.After(1 * time.Second):
			// 定期清理已完成的任务
			count := a.taskManager.CleanupCompletedTasks()
			if count > 0 {
				logger.Info("Cleaned up completed tasks", zap.Int("count", count))
			}
		}
	}
}