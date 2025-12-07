package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/config"
	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"go.uber.org/zap"
)

// WorkerPool 工作者池
type WorkerPool struct {
	minWorkers     int
	maxWorkers     int
	currentWorkers int
	taskQueue      chan *TaskWrapper
	workerQueue    chan chan *TaskWrapper
	workers        map[string]*Worker
	quitChan       chan bool
	wg             sync.WaitGroup
	mutex          sync.RWMutex
	config         *config.WorkerConfig
	taskUpdateChan chan *model.Task
}

// NewWorkerPool 创建新的工作者池
func NewWorkerPool(cfg *config.Config) *WorkerPool {
	pool := &WorkerPool{
		minWorkers:     cfg.Worker.MinWorkers,
		maxWorkers:     cfg.Worker.MaxWorkers,
		currentWorkers: 0,
		taskQueue:      make(chan *TaskWrapper, cfg.Task.QueueCapacity),
		workerQueue:    make(chan chan *TaskWrapper, cfg.Worker.MaxWorkers),
		workers:        make(map[string]*Worker),
		quitChan:       make(chan bool),
		config:         &cfg.Worker,
		taskUpdateChan: make(chan *model.Task, 100),
	}

	return pool
}

// Start 启动工作者池
func (wp *WorkerPool) Start() {
	// 启动最小数量的工作者
	for i := 0; i < wp.minWorkers; i++ {
		worker := NewWorker(wp, generateWorkerID())
		wp.mutex.Lock()
		wp.workers[worker.ID] = worker
		wp.currentWorkers++
		wp.mutex.Unlock()
		worker.Start()
		logger.Info("Worker started", zap.String("worker_id", worker.ID))
	}

	// 启动任务分发器
	go wp.dispatch()

	// 如果启用了自动扩展，则启动自动扩展器
	if wp.config.EnableAutoScaling {
		go wp.autoscale()
	}

	logger.Info("Worker pool started",
		zap.Int("min_workers", wp.minWorkers),
		zap.Int("max_workers", wp.maxWorkers),
		zap.Int("initial_workers", wp.minWorkers))
}

// Stop 停止工作者池
func (wp *WorkerPool) Stop() error {
	logger.Info("Stopping worker pool")

	// 发送停止信号
	close(wp.quitChan)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), wp.config.ShutdownTimeout)
	defer cancel()

	// 等待所有工作者完成或超时
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有工作者正常结束
		logger.Info("All workers stopped gracefully")
		return nil
	case <-ctx.Done():
		// 超时
		logger.Warn("Worker pool shutdown timeout",
			zap.Duration("timeout", wp.config.ShutdownTimeout))
		return ctx.Err()
	}
}

// SubmitTask 提交任务
func (wp *WorkerPool) SubmitTask(task *model.Task, taskFunc model.TaskFunc) error {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), wp.config.ShutdownTimeout)

	// 包装任务
	wrapper := &TaskWrapper{
		Task:     task,
		TaskFunc: taskFunc,
		Context:  ctx,
		Cancel:   cancel,
	}

	// 尝试将任务放入队列
	select {
	case wp.taskQueue <- wrapper:
		logger.Debug("Task submitted to queue", zap.String("task_id", task.ID))
		return nil
	case <-ctx.Done():
		cancel()
		return fmt.Errorf("failed to submit task: %w", ctx.Err())
	}
}

// dispatch 任务分发器
func (wp *WorkerPool) dispatch() {
	for {
		select {
		case task := <-wp.taskQueue:
			// 尝试将任务分发给可用的工作者
			select {
			case workerTaskQueue := <-wp.workerQueue:
				workerTaskQueue <- task
				logger.Debug("Task dispatched to worker",
					zap.String("task_id", task.Task.ID),
					zap.String("worker_id", getWorkerIDFromTaskQueue(workerTaskQueue)))
			default:
				// 没有可用的工作者，可能需要扩展
				logger.Warn("No available workers, task queued", zap.String("task_id", task.Task.ID))
				// 将任务重新放回队列
				go func() {
					wp.taskQueue <- task
				}()
			}
		case <-wp.quitChan:
			logger.Info("Dispatcher stopping")
			return
		}
	}
}

// autoscale 自动扩展工作者数量
func (wp *WorkerPool) autoscale() {
	ticker := time.NewTicker(wp.config.ScaleCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wp.checkAndScale()
		case <-wp.quitChan:
			logger.Info("Autoscaler stopping")
			return
		}
	}
}

// checkAndScale 检查并调整工作者数量
func (wp *WorkerPool) checkAndScale() {
	wp.mutex.RLock()
	currentWorkers := wp.currentWorkers
	queueLength := len(wp.taskQueue)
	wp.mutex.RUnlock()

	// 计算队列使用率
	queueUsage := float64(queueLength) / float64(cap(wp.taskQueue))

	logger.Debug("Autoscale check",
		zap.Int("current_workers", currentWorkers),
		zap.Int("queue_length", queueLength),
		zap.Float64("queue_usage", queueUsage))

	// 检查是否需要扩展
	if queueUsage > wp.config.ScaleUpThreshold && currentWorkers < wp.maxWorkers {
		// 需要扩展
		newWorkers := min(wp.maxWorkers-currentWorkers, 5) // 每次最多扩展5个
		for i := 0; i < newWorkers; i++ {
			wp.addWorker()
		}
		logger.Info("Scaled up workers",
			zap.Int("added", newWorkers),
			zap.Int("total", currentWorkers+newWorkers))
	} else if queueUsage < wp.config.ScaleDownThreshold && currentWorkers > wp.minWorkers {
		// 需要缩减
		removeWorkers := min(currentWorkers-wp.minWorkers, 2) // 每次最多缩减2个
		for i := 0; i < removeWorkers; i++ {
			wp.removeWorker()
		}
		logger.Info("Scaled down workers",
			zap.Int("removed", removeWorkers),
			zap.Int("total", currentWorkers-removeWorkers))
	}
}

// addWorker 添加工作者
func (wp *WorkerPool) addWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.currentWorkers >= wp.maxWorkers {
		return
	}

	worker := NewWorker(wp, generateWorkerID())
	wp.workers[worker.ID] = worker
	wp.currentWorkers++
	worker.Start()

	logger.Info("Added worker", zap.String("worker_id", worker.ID))
}

// removeWorker 移除工作者
func (wp *WorkerPool) removeWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.currentWorkers <= wp.minWorkers {
		return
	}

	// 找到一个工作者并停止它
	for id, worker := range wp.workers {
		worker.Stop()
		delete(wp.workers, id)
		wp.currentWorkers--
		logger.Info("Removed worker", zap.String("worker_id", id))
		break
	}
}

// notifyTaskUpdate 通知任务更新
func (wp *WorkerPool) notifyTaskUpdate(task *model.Task) {
	select {
	case wp.taskUpdateChan <- task:
		// 通知成功
	default:
		// 通道满了，丢弃通知
		logger.Warn("Task update channel full, dropping notification",
			zap.String("task_id", task.ID))
	}
}

// GetTaskUpdateChan 获取任务更新通道
func (wp *WorkerPool) GetTaskUpdateChan() <-chan *model.Task {
	return wp.taskUpdateChan
}

// getWorkerIDFromTaskQueue 从任务队列获取工作者ID（辅助函数）
func getWorkerIDFromTaskQueue(taskQueue chan *TaskWrapper) string {
	// 这是一个简化的实现，实际应用中可能需要更好的方式来识别工作者
	return "unknown"
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}