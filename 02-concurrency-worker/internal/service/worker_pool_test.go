package service

import (
	"testing"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/config"
	"github.com/yourname/02-concurrency-worker/internal/model"
)

func TestWorkerPoolCreation(t *testing.T) {
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			MinWorkers:        2,
			MaxWorkers:        10,
			EnableAutoScaling: false,
			ShutdownTimeout:   5 * time.Second,
		},
		Task: config.TaskConfig{
			QueueCapacity: 100,
		},
	}

	pool := NewWorkerPool(cfg)

	if pool.minWorkers != 2 {
		t.Errorf("Expected minWorkers to be 2, got %d", pool.minWorkers)
	}

	if pool.maxWorkers != 10 {
		t.Errorf("Expected maxWorkers to be 10, got %d", pool.maxWorkers)
	}

	if cap(pool.taskQueue) != 100 {
		t.Errorf("Expected taskQueue capacity to be 100, got %d", cap(pool.taskQueue))
	}
}

func TestTaskWrapperCreation(t *testing.T) {
	task := &model.Task{
		ID:        "test-task-1",
		Name:      "test_task",
		Payload:   "test payload",
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	// Create a simple task function
	taskFunc := func(task *model.Task) error {
		task.Result = "task completed"
		return nil
	}

	// Create task wrapper
	wrapper := &TaskWrapper{
		Task:     task,
		TaskFunc: taskFunc,
	}

	if wrapper.Task != task {
		t.Error("Expected Task to match")
	}

	if wrapper.TaskFunc == nil {
		t.Error("Expected TaskFunc to not be nil")
	}
}

func TestGenerateWorkerID(t *testing.T) {
	id1 := generateWorkerID()
	id2 := generateWorkerID()

	if id1 == id2 {
		t.Error("Expected generated worker IDs to be different")
	}

	if len(id1) == 0 {
		t.Error("Expected generated worker ID to not be empty")
	}
}