package websocket

import (
	"testing"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
)

func TestHubCreation(t *testing.T) {
	hub := NewHub()

	if hub.clients == nil {
		t.Error("Expected clients map to not be nil")
	}

	if hub.broadcast == nil {
		t.Error("Expected broadcast channel to not be nil")
	}

	if hub.register == nil {
		t.Error("Expected register channel to not be nil")
	}

	if hub.unregister == nil {
		t.Error("Expected unregister channel to not be nil")
	}
}

func TestTaskUpdateMessageCreation(t *testing.T) {
	message := &TaskUpdateMessage{
		TaskID:    "test-task-1",
		Status:    model.TaskRunning,
		Progress:  50,
		Result:    "test result",
		Error:     "",
		Timestamp: time.Now(),
	}

	if message.TaskID != "test-task-1" {
		t.Errorf("Expected TaskID to be 'test-task-1', got '%s'", message.TaskID)
	}

	if message.Status != model.TaskRunning {
		t.Errorf("Expected Status to be TaskRunning, got '%s'", message.Status)
	}

	if message.Progress != 50 {
		t.Errorf("Expected Progress to be 50, got %d", message.Progress)
	}

	if message.Result != "test result" {
		t.Errorf("Expected Result to be 'test result', got '%v'", message.Result)
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{5, 3, 3},
		{2, 8, 2},
		{4, 4, 4},
	}

	for _, test := range tests {
		result := min(test.a, test.b)
		if result != test.expected {
			t.Errorf("Expected min(%d, %d) to be %d, got %d", test.a, test.b, test.expected, result)
		}
	}
}