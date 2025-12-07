package service

import (
	"sync"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"go.uber.org/zap"
)

// TaskRegistry 任务注册表
type TaskRegistry struct {
	tasks map[string]*model.Task
	mutex sync.RWMutex
}

// NewTaskRegistry 创建新的任务注册表
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		tasks: make(map[string]*model.Task),
	}
}

// Store 存储任务
func (tr *TaskRegistry) Store(task *model.Task) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tr.tasks[task.ID] = task
	logger.Debug("Task stored", zap.String("task_id", task.ID))
}

// Get 获取任务
func (tr *TaskRegistry) Get(taskID string) (*model.Task, bool) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	task, exists := tr.tasks[taskID]
	return task, exists
}

// GetAll 获取所有任务
func (tr *TaskRegistry) GetAll() []*model.Task {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	tasks := make([]*model.Task, 0, len(tr.tasks))
	for _, task := range tr.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetByStatus 根据状态获取任务
func (tr *TaskRegistry) GetByStatus(status model.TaskStatus) []*model.Task {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	tasks := make([]*model.Task, 0)
	for _, task := range tr.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// Update 更新任务
func (tr *TaskRegistry) Update(task *model.Task) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tr.tasks[task.ID] = task
	logger.Debug("Task updated", zap.String("task_id", task.ID))
}

// Delete 删除任务
func (tr *TaskRegistry) Delete(taskID string) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	delete(tr.tasks, taskID)
	logger.Debug("Task deleted", zap.String("task_id", taskID))
}

// CleanupCompleted 清理已完成的任务
func (tr *TaskRegistry) CleanupCompleted(before time.Time) int {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	count := 0
	for id, task := range tr.tasks {
		if (task.Status == model.TaskCompleted || task.Status == model.TaskFailed) &&
			task.CompletedAt != nil && task.CompletedAt.Before(before) {
			delete(tr.tasks, id)
			count++
		}
	}

	if count > 0 {
		logger.Info("Cleaned up completed tasks", zap.Int("count", count))
	}

	return count
}

// GetStats 获取任务统计信息
func (tr *TaskRegistry) GetStats() map[model.TaskStatus]int {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	stats := make(map[model.TaskStatus]int)
	for _, task := range tr.tasks {
		stats[task.Status]++
	}

	return stats
}