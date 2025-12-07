package api

import (
	"net/http"

	"github.com/yourname/02-concurrency-worker/internal/api/dto"
	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
)

// Handlers API处理器
type Handlers struct {
	taskManager *service.TaskManager
}

// NewHandlers 创建新的API处理器
func NewHandlers(taskManager *service.TaskManager) *Handlers {
	return &Handlers{
		taskManager: taskManager,
	}
}

// SubmitTask 提交任务
func (h *Handlers) SubmitTask(c *gin.Context) {
	var req dto.SubmitTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskManager.SubmitTask(req.Name, req.Payload)
	if err != nil {
		logger.Error("Failed to submit task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ConvertTaskToResponse(task)
	c.JSON(http.StatusCreated, response)
}

// GetTask 获取任务
func (h *Handlers) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	task, exists := h.taskManager.GetTask(taskID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	response := dto.ConvertTaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// GetAllTasks 获取所有任务
func (h *Handlers) GetAllTasks(c *gin.Context) {
	tasks := h.taskManager.GetAllTasks()

	response := &dto.TaskListResponse{
		Tasks: dto.ConvertTasksToResponse(tasks),
		Total: len(tasks),
	}

	c.JSON(http.StatusOK, response)
}

// GetTasksByStatus 根据状态获取任务
func (h *Handlers) GetTasksByStatus(c *gin.Context) {
	status := c.Query("status")

	tasks := h.taskManager.GetAllTasks()
	if status != "" {
		// 过滤任务
		filteredTasks := make([]*model.Task, 0)
		for _, task := range tasks {
			if string(task.Status) == status {
				filteredTasks = append(filteredTasks, task)
			}
		}
		tasks = filteredTasks
	}

	response := &dto.TaskListResponse{
		Tasks: dto.ConvertTasksToResponse(tasks),
		Total: len(tasks),
	}

	c.JSON(http.StatusOK, response)
}

// CancelTask 取消任务
func (h *Handlers) CancelTask(c *gin.Context) {
	taskID := c.Param("id")

	err := h.taskManager.CancelTask(taskID)
	if err != nil {
		logger.Error("Failed to cancel task", zap.String("task_id", taskID), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task cancelled successfully"})
}

// GetTaskStats 获取任务统计
func (h *Handlers) GetTaskStats(c *gin.Context) {
	stats := h.taskManager.GetTaskStats()

	response := &dto.TaskStatsResponse{
		Stats: stats,
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck 健康检查
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "Worker pool is running",
	})
}