package config

import (
	"testing"
	"time"
)

func TestConfigStructs(t *testing.T) {
	// Test ServerConfig
	serverCfg := ServerConfig{
		Port: 8080,
		Host: "localhost",
		Mode: "debug",
	}

	if serverCfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", serverCfg.Port)
	}

	if serverCfg.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", serverCfg.Host)
	}

	if serverCfg.Mode != "debug" {
		t.Errorf("Expected Mode to be 'debug', got '%s'", serverCfg.Mode)
	}

	// Test WorkerConfig
	workerCfg := WorkerConfig{
		MinWorkers:              5,
		MaxWorkers:              50,
		EnableAutoScaling:       true,
		ScaleUpThreshold:        0.8,
		ScaleDownThreshold:      0.3,
		ScaleCheckInterval:      30 * time.Second,
		ShutdownTimeout:         30 * time.Second,
	}

	if workerCfg.MinWorkers != 5 {
		t.Errorf("Expected MinWorkers to be 5, got %d", workerCfg.MinWorkers)
	}

	if !workerCfg.EnableAutoScaling {
		t.Error("Expected EnableAutoScaling to be true")
	}

	if workerCfg.ScaleUpThreshold != 0.8 {
		t.Errorf("Expected ScaleUpThreshold to be 0.8, got %f", workerCfg.ScaleUpThreshold)
	}

	// Test TaskConfig
	taskCfg := TaskConfig{
		QueueCapacity:                1000,
		MaxConcurrentTasks:           100,
		DefaultTaskTimeout:           5 * time.Minute,
		CleanupCompletedTasksAfter:   1 * time.Hour,
	}

	if taskCfg.QueueCapacity != 1000 {
		t.Errorf("Expected QueueCapacity to be 1000, got %d", taskCfg.QueueCapacity)
	}

	if taskCfg.DefaultTaskTimeout != 5*time.Minute {
		t.Errorf("Expected DefaultTaskTimeout to be 5 minutes, got %v", taskCfg.DefaultTaskTimeout)
	}

	// Test LoggerConfig
	loggerCfg := LoggerConfig{
		Level:    "debug",
		Format:   "console",
		Output:   "stdout",
		FilePath: "./logs/worker.log",
	}

	if loggerCfg.Level != "debug" {
		t.Errorf("Expected Level to be 'debug', got '%s'", loggerCfg.Level)
	}

	if loggerCfg.Format != "console" {
		t.Errorf("Expected Format to be 'console', got '%s'", loggerCfg.Format)
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
			Mode: "debug",
		},
		Worker: WorkerConfig{
			MinWorkers:        5,
			MaxWorkers:        50,
			EnableAutoScaling: true,
		},
		Task: TaskConfig{
			QueueCapacity: 1000,
		},
		Logger: LoggerConfig{
			Level: "debug",
		},
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected Server.Port to be 8080, got %d", cfg.Server.Port)
	}

	if cfg.Worker.MinWorkers != 5 {
		t.Errorf("Expected Worker.MinWorkers to be 5, got %d", cfg.Worker.MinWorkers)
	}

	if cfg.Task.QueueCapacity != 1000 {
		t.Errorf("Expected Task.QueueCapacity to be 1000, got %d", cfg.Task.QueueCapacity)
	}

	if cfg.Logger.Level != "debug" {
		t.Errorf("Expected Logger.Level to be 'debug', got '%s'", cfg.Logger.Level)
	}
}