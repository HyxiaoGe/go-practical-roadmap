package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Task     TaskConfig     `mapstructure:"task"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

// WorkerConfig Worker配置
type WorkerConfig struct {
	MinWorkers              int           `mapstructure:"min_workers"`
	MaxWorkers              int           `mapstructure:"max_workers"`
	EnableAutoScaling       bool          `mapstructure:"enable_auto_scaling"`
	ScaleUpThreshold        float64       `mapstructure:"scale_up_threshold"`
	ScaleDownThreshold      float64       `mapstructure:"scale_down_threshold"`
	ScaleCheckInterval      time.Duration `mapstructure:"scale_check_interval"`
	ShutdownTimeout         time.Duration `mapstructure:"shutdown_timeout"`
}

// TaskConfig 任务配置
type TaskConfig struct {
	QueueCapacity                int           `mapstructure:"queue_capacity"`
	MaxConcurrentTasks           int           `mapstructure:"max_concurrent_tasks"`
	DefaultTaskTimeout           time.Duration `mapstructure:"default_task_timeout"`
	CleanupCompletedTasksAfter   time.Duration `mapstructure:"cleanup_completed_tasks_after"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("worker.min_workers", 5)
	viper.SetDefault("worker.max_workers", 50)
	viper.SetDefault("worker.enable_auto_scaling", true)
	viper.SetDefault("worker.scale_up_threshold", 0.8)
	viper.SetDefault("worker.scale_down_threshold", 0.3)
	viper.SetDefault("worker.scale_check_interval", "30s")
	viper.SetDefault("worker.shutdown_timeout", "30s")

	viper.SetDefault("task.queue_capacity", 1000)
	viper.SetDefault("task.max_concurrent_tasks", 100)
	viper.SetDefault("task.default_task_timeout", "5m")
	viper.SetDefault("task.cleanup_completed_tasks_after", "1h")

	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("logger.output", "stdout")

	// 设置环境变量前缀
	viper.SetEnvPrefix("APP")

	// 自动绑定环境变量
	viper.AutomaticEnv()

	// 设置环境变量替换规则
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}