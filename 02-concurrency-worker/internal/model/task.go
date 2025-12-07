package model

import (
	"time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
)

// Task 任务模型
type Task struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Payload     interface{} `json:"payload"`
	Status      TaskStatus  `json:"status"`
	Result      interface{} `json:"result"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	Progress    int         `json:"progress"` // 0-100
}

// TaskFunc 任务执行函数类型
type TaskFunc func(task *Task) error