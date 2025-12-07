package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志记录器接口
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

// zapLogger zap日志记录器实现
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger 创建新的日志记录器
func NewLogger(level string, format string, outputPath string) (Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// 创建encoder配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建输出路径
	var outputs []zapcore.WriteSyncer
	if outputPath == "stdout" {
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	} else {
		// 使用lumberjack进行日志轮转
		lumberJackLogger := &lumberjack.Logger{
			Filename:   outputPath,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
		}
		outputs = append(outputs, zapcore.AddSync(lumberJackLogger))
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(outputs...), zapLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return &zapLogger{logger: logger}, nil
}

// Debug 记录调试日志
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info 记录信息日志
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn 记录警告日志
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error 记录错误日志
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal 记录致命错误日志
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Sync 同步日志缓冲区
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}