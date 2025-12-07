package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
)

func TestTaskRegistry_StoreAndGet(t *testing.T) {
	registry := NewTaskRegistry()

	task := &model.Task{
		ID:        "test-task-1",
		Name:      "test_task",
		Payload:   "test payload",
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	// Store task
	registry.Store(task)

	// Retrieve task
	retrievedTask, exists := registry.Get("test-task-1")
	if !exists {
		t.Error("Expected task to exist")
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", task.ID, retrievedTask.ID)
	}

	if retrievedTask.Name != task.Name {
		t.Errorf("Expected Name to be '%s', got '%s'", task.Name, retrievedTask.Name)
	}
}

func TestTaskRegistry_GetAll(t *testing.T) {
	registry := NewTaskRegistry()

	// Store multiple tasks
	for i := 0; i < 3; i++ {
		task := &model.Task{
			ID:        fmt.Sprintf("test-task-%d", i),
			Name:      fmt.Sprintf("test_task_%d", i),
			Payload:   fmt.Sprintf("test payload %d", i),
			Status:    model.TaskPending,
			CreatedAt: time.Now(),
			Progress:  0,
		}
		registry.Store(task)
	}

	// Get all tasks
	tasks := registry.GetAll()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestTaskRegistry_GetByStatus(t *testing.T) {
	registry := NewTaskRegistry()

	// Store tasks with different statuses
	pendingTask := &model.Task{
		ID:        "pending-task",
		Name:      "pending_task",
		Payload:   "test payload",
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	runningTask := &model.Task{
		ID:        "running-task",
		Name:      "running_task",
		Payload:   "test payload",
		Status:    model.TaskRunning,
		CreatedAt: time.Now(),
		Progress:  50,
	}

	registry.Store(pendingTask)
	registry.Store(runningTask)

	// Get tasks by status
	pendingTasks := registry.GetByStatus(model.TaskPending)
	if len(pendingTasks) != 1 {
		t.Errorf("Expected 1 pending task, got %d", len(pendingTasks))
	}

	runningTasks := registry.GetByStatus(model.TaskRunning)
	if len(runningTasks) != 1 {
		t.Errorf("Expected 1 running task, got %d", len(runningTasks))
	}
}

func TestTaskRegistry_Update(t *testing.T) {
	registry := NewTaskRegistry()

	task := &model.Task{
		ID:        "test-task-1",
		Name:      "test_task",
		Payload:   "test payload",
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	// Store task
	registry.Store(task)

	// Update task
	task.Status = model.TaskRunning
	task.Progress = 50
	registry.Update(task)

	// Retrieve updated task
	updatedTask, exists := registry.Get("test-task-1")
	if !exists {
		t.Error("Expected task to exist")
	}

	if updatedTask.Status != model.TaskRunning {
		t.Errorf("Expected Status to be TaskRunning, got '%s'", updatedTask.Status)
	}

	if updatedTask.Progress != 50 {
		t.Errorf("Expected Progress to be 50, got %d", updatedTask.Progress)
	}
}

func TestTaskRegistry_Delete(t *testing.T) {
	registry := NewTaskRegistry()

	task := &model.Task{
		ID:        "test-task-1",
		Name:      "test_task",
		Payload:   "test payload",
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	// Store task
	registry.Store(task)

	// Delete task
	registry.Delete("test-task-1")

	// Try to retrieve deleted task
	_, exists := registry.Get("test-task-1")
	if exists {
		t.Error("Expected task to not exist after deletion")
	}
}

func TestTaskRegistry_GetStats(t *testing.T) {
	registry := NewTaskRegistry()

	// Store tasks with different statuses
	tasks := []*model.Task{
		{ID: "task-1", Status: model.TaskPending},
		{ID: "task-2", Status: model.TaskRunning},
		{ID: "task-3", Status: model.TaskCompleted},
		{ID: "task-4", Status: model.TaskFailed},
		{ID: "task-5", Status: model.TaskCancelled},
		{ID: "task-6", Status: model.TaskPending},
		{ID: "task-7", Status: model.TaskRunning},
	}

	for _, task := range tasks {
		registry.Store(task)
	}

	// Get stats
	stats := registry.GetStats()

	if stats[model.TaskPending] != 2 {
		t.Errorf("Expected 2 pending tasks, got %d", stats[model.TaskPending])
	}

	if stats[model.TaskRunning] != 2 {
		t.Errorf("Expected 2 running tasks, got %d", stats[model.TaskRunning])
	}

	if stats[model.TaskCompleted] != 1 {
		t.Errorf("Expected 1 completed task, got %d", stats[model.TaskCompleted])
	}

	if stats[model.TaskFailed] != 1 {
		t.Errorf("Expected 1 failed task, got %d", stats[model.TaskFailed])
	}

	if stats[model.TaskCancelled] != 1 {
		t.Errorf("Expected 1 cancelled task, got %d", stats[model.TaskCancelled])
	}
}