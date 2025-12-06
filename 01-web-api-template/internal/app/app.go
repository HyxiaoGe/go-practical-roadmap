package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-practical-roadmap/01-web-api-template/internal/api"
	"go-practical-roadmap/01-web-api-template/internal/config"
	"go-practical-roadmap/01-web-api-template/internal/model"
	"go-practical-roadmap/01-web-api-template/pkg/db"
	"go-practical-roadmap/01-web-api-template/pkg/logger"
	"go.uber.org/zap"
)

// App 应用结构体
type App struct {
	server *http.Server
}

// NewApp 创建新的应用实例
func NewApp() (*App, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.FilePath); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 连接数据库
	if err := db.Connect(cfg.Database.DSN, cfg.Database.Driver); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库
	if err := autoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &App{}, nil
}

// autoMigrate 自动迁移数据库
func autoMigrate() error {
	// 获取数据库实例
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 自动迁移模型
	err := database.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Database migration completed successfully")
	return nil
}

// Run 启动应用
func (a *App) Run() error {
	// 创建路由
	router := api.SetupRoutes()

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 启动服务器
	logger.Info("Starting server", zap.String("addr", a.server.Addr))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop 停止应用
func (a *App) Stop() error {
	logger.Info("Shutting down server...")

	// 创建优雅关闭的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// 关闭数据库连接
	if err := db.Close(); err != nil {
		return fmt.Errorf("database close failed: %w", err)
	}

	// 同步日志
	if logger.GlobalLogger != nil {
		logger.GlobalLogger.Sync()
	}

	logger.Info("Server exited properly")
	return nil
}

// WaitForInterrupt 等待中断信号
func (a *App) WaitForInterrupt() {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received")
}