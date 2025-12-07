package model

import (
	"testing"
	"time"
)

func TestTaskCreation(t *testing.T) {
	task := &Task{
		ID:        "test-task-1",
		Name:      "test_task",
		Payload:   "test payload",
		Status:    TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	if task.ID != "test-task-1" {
		t.Errorf("Expected ID to be 'test-task-1', got '%s'", task.ID)
	}

	if task.Name != "test_task" {
		t.Errorf("Expected Name to be 'test_task', got '%s'", task.Name)
	}

	if task.Status != TaskPending {
		t.Errorf("Expected Status to be TaskPending, got '%s'", task.Status)
	}

	if task.Progress != 0 {
		t.Errorf("Expected Progress to be 0, got %d", task.Progress)
	}
}

func TestTaskStatusEnum(t *testing.T) {
	statuses := []TaskStatus{
		TaskPending,
		TaskRunning,
		TaskCompleted,
		TaskFailed,
		TaskCancelled,
	}

	expected := []string{"pending", "running", "completed", "failed", "cancelled"}

	for i, status := range statuses {
		if string(status) != expected[i] {
			t.Errorf("Expected status '%s', got '%s'", expected[i], string(status))
		}
	}
}