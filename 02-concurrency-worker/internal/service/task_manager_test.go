package service

import (
	"testing"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/config"
	"github.com/yourname/02-concurrency-worker/internal/model"
)

func TestTaskManagerCreation(t *testing.T) {
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			MinWorkers:        2,
			MaxWorkers:        10,
			EnableAutoScaling: false,
			ShutdownTimeout:   5 * time.Second,
		},
		Task: config.TaskConfig{
			QueueCapacity:              100,
			DefaultTaskTimeout:         30 * time.Second,
			CleanupCompletedTasksAfter: 1 * time.Hour,
		},
	}

	manager := NewTaskManager(cfg)

	if manager.pool == nil {
		t.Error("Expected pool to not be nil")
	}

	if manager.registry == nil {
		t.Error("Expected registry to not be nil")
	}

	if manager.config.QueueCapacity != 100 {
		t.Errorf("Expected QueueCapacity to be 100, got %d", manager.config.QueueCapacity)
	}
}

func TestGenerateTaskID(t *testing.T) {
	id1 := generateTaskID()
	id2 := generateTaskID()

	if id1 == id2 {
		t.Error("Expected generated task IDs to be different")
	}

	if len(id1) == 0 {
		t.Error("Expected generated task ID to not be empty")
	}
}

func TestTaskFunctions(t *testing.T) {
	manager := &TaskManager{}

	// Test example task function
	task := &model.Task{
		ID:      "test-task-1",
		Payload: "test payload",
	}

	err := manager.exampleTask(task)
	if err != nil {
		t.Errorf("Expected exampleTask to not return error, got %v", err)
	}

	if task.Result == nil {
		t.Error("Expected task result to be set")
	}

	// Test default task function
	task2 := &model.Task{
		ID:      "test-task-2",
		Payload: "test payload 2",
	}

	err = manager.defaultTask(task2)
	if err != nil {
		t.Errorf("Expected defaultTask to not return error, got %v", err)
	}

	if task2.Result == nil {
		t.Error("Expected task result to be set")
	}
}