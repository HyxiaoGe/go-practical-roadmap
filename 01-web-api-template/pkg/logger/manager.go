package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	// GlobalLogger 全局日志记录器实例
	GlobalLogger Logger

	// once 确保只初始化一次
	once sync.Once
)

// InitLogger 初始化全局日志记录器
func InitLogger(level string, format string, outputPath string) error {
	var err error
	once.Do(func() {
		GlobalLogger, err = NewLogger(level, format, outputPath)
	})
	return err
}

// Debug 记录调试日志
func Debug(msg string, fields ...zap.Field) {
	if GlobalLogger != nil {
		GlobalLogger.Debug(msg, fields...)
	}
}

// Info 记录信息日志
func Info(msg string, fields ...zap.Field) {
	if GlobalLogger != nil {
		GlobalLogger.Info(msg, fields...)
	}
}

// Warn 记录警告日志
func Warn(msg string, fields ...zap.Field) {
	if GlobalLogger != nil {
		GlobalLogger.Warn(msg, fields...)
	}
}

// Error 记录错误日志
func Error(msg string, fields ...zap.Field) {
	if GlobalLogger != nil {
		GlobalLogger.Error(msg, fields...)
	}
}

// Fatal 记录致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	if GlobalLogger != nil {
		GlobalLogger.Fatal(msg, fields...)
	}
}