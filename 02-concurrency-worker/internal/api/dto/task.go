package dto

import (
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
)

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	Name    string      `json:"name" binding:"required"`
	Payload interface{} `json:"payload"`
	Timeout int         `json:"timeout,omitempty"` // 超时时间（秒）
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Payload     interface{}       `json:"payload"`
	Status      model.TaskStatus  `json:"status"`
	Result      interface{}       `json:"result"`
	Error       string            `json:"error,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Progress    int               `json:"progress"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
	Total int             `json:"total"`
}

// TaskStatsResponse 任务统计响应
type TaskStatsResponse struct {
	Stats map[model.TaskStatus]int `json:"stats"`
}

// ConvertTaskToResponse 将任务模型转换为响应DTO
func ConvertTaskToResponse(task *model.Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Payload:     task.Payload,
		Status:      task.Status,
		Result:      task.Result,
		Error:       task.Error,
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Progress:    task.Progress,
	}
}

// ConvertTasksToResponse 将任务模型列表转换为响应DTO列表
func ConvertTasksToResponse(tasks []*model.Task) []*TaskResponse {
	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = ConvertTaskToResponse(task)
	}
	return responses
}