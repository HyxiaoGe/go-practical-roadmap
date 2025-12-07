package service

import (
	"fmt"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/config"
	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"go.uber.org/zap"
)

// TaskManager 任务管理器
type TaskManager struct {
	pool     *WorkerPool
	registry *TaskRegistry
	config   *config.TaskConfig
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(cfg *config.Config) *TaskManager {
	pool := NewWorkerPool(cfg)
	registry := NewTaskRegistry()

	return &TaskManager{
		pool:     pool,
		registry: registry,
		config:   &cfg.Task,
	}
}

// Start 启动任务管理器
func (tm *TaskManager) Start() {
	tm.pool.Start()
	logger.Info("Task manager started")
}

// Stop 停止任务管理器
func (tm *TaskManager) Stop() error {
	logger.Info("Stopping task manager")
	return tm.pool.Stop()
}

// SubmitTask 提交任务
func (tm *TaskManager) SubmitTask(name string, payload interface{}) (*model.Task, error) {
	// 创建任务
	task := &model.Task{
		ID:        generateTaskID(),
		Name:      name,
		Payload:   payload,
		Status:    model.TaskPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}

	// 存储任务到注册表
	tm.registry.Store(task)

	// 获取任务执行函数
	taskFunc := tm.getTaskFunction(name)

	// 提交任务到工作者池
	err := tm.pool.SubmitTask(task, taskFunc)
	if err != nil {
		task.Status = model.TaskFailed
		task.Error = err.Error()
		tm.registry.Update(task)
		return nil, fmt.Errorf("failed to submit task: %w", err)
	}

	logger.Info("Task submitted", zap.String("task_id", task.ID), zap.String("task_name", name))
	return task, nil
}

// GetTask 获取任务
func (tm *TaskManager) GetTask(taskID string) (*model.Task, bool) {
	return tm.registry.Get(taskID)
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() []*model.Task {
	return tm.registry.GetAll()
}

// GetTasksByStatus 根据状态获取任务
func (tm *TaskManager) GetTasksByStatus(status model.TaskStatus) []*model.Task {
	return tm.registry.GetByStatus(status)
}

// CancelTask 取消任务
func (tm *TaskManager) CancelTask(taskID string) error {
	task, exists := tm.registry.Get(taskID)
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// 如果任务正在运行，尝试取消
	if task.Status == model.TaskRunning {
		// 注意：这里我们无法真正取消正在运行的goroutine，
		// 但可以在任务函数中检查上下文来实现协作式取消
		logger.Info("Task cancellation requested", zap.String("task_id", taskID))
	}

	// 更新任务状态
	task.Status = model.TaskCancelled
	now := time.Now()
	task.CompletedAt = &now
	tm.registry.Update(task)

	logger.Info("Task cancelled", zap.String("task_id", taskID))
	return nil
}

// GetTaskUpdateChan 获取任务更新通道
func (tm *TaskManager) GetTaskUpdateChan() <-chan *model.Task {
	return tm.pool.GetTaskUpdateChan()
}

// CleanupCompletedTasks 清理已完成的任务
func (tm *TaskManager) CleanupCompletedTasks() int {
	before := time.Now().Add(-tm.config.CleanupCompletedTasksAfter)
	return tm.registry.CleanupCompleted(before)
}

// GetTaskStats 获取任务统计信息
func (tm *TaskManager) GetTaskStats() map[model.TaskStatus]int {
	return tm.registry.GetStats()
}

// getTaskFunction 根据任务名称获取任务执行函数
func (tm *TaskManager) getTaskFunction(name string) model.TaskFunc {
	switch name {
	case "example_task":
		return tm.exampleTask
	case "long_running_task":
		return tm.longRunningTask
	default:
		return tm.defaultTask
	}
}

// exampleTask 示例任务函数
func (tm *TaskManager) exampleTask(task *model.Task) error {
	logger.Info("Executing example task", zap.String("task_id", task.ID))

	// 模拟一些工作
	time.Sleep(100 * time.Millisecond)

	// 设置结果
	task.Result = map[string]interface{}{
		"message": "Task completed successfully",
		"data":    task.Payload,
	}

	return nil
}

// longRunningTask 长时间运行的任务函数
func (tm *TaskManager) longRunningTask(task *model.Task) error {
	logger.Info("Executing long running task", zap.String("task_id", task.ID))

	// 模拟长时间运行的任务
	for i := 0; i < 10; i++ {
		// 检查是否被取消（在实际应用中应该检查上下文）
		time.Sleep(100 * time.Millisecond)
		task.Progress = (i + 1) * 10
	}

	// 设置结果
	task.Result = map[string]interface{}{
		"message": "Long running task completed",
		"steps":   10,
	}

	return nil
}

// defaultTask 默认任务函数
func (tm *TaskManager) defaultTask(task *model.Task) error {
	logger.Info("Executing default task", zap.String("task_id", task.ID))

	// 简单的任务处理
	result := fmt.Sprintf("Processed payload: %v", task.Payload)
	task.Result = result

	return nil
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}