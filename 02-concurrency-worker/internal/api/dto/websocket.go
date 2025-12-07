package dto

import (
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
)

// TaskUpdateMessage 任务更新消息
type TaskUpdateMessage struct {
	TaskID    string           `json:"task_id"`
	Status    model.TaskStatus `json:"status"`
	Progress  int              `json:"progress,omitempty"`
	Result    interface{}      `json:"result,omitempty"`
	Error     string           `json:"error,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
}

// WebSocketResponse WebSocket响应
type WebSocketResponse struct {
	Type    string      `json:"type"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}