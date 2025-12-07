package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yourname/02-concurrency-worker/internal/model"
	"github.com/yourname/02-concurrency-worker/internal/pkg/logger"
	"go.uber.org/zap"
)

// Worker 工作者结构体
type Worker struct {
	ID         string
	TaskQueue  chan *TaskWrapper
	WorkerPool *WorkerPool
	QuitChan   chan bool
	Context    context.Context
	Cancel     context.CancelFunc
}

// TaskWrapper 任务包装器
type TaskWrapper struct {
	Task     *model.Task
	TaskFunc model.TaskFunc
	Context  context.Context
	Cancel   context.CancelFunc
}

// NewWorker 创建新的工作者
func NewWorker(workerPool *WorkerPool, id string) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	worker := &Worker{
		ID:         id,
		TaskQueue:  make(chan *TaskWrapper),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
		Context:    ctx,
		Cancel:     cancel,
	}

	return worker
}

// Start 启动工作者
func (w *Worker) Start() {
	w.WorkerPool.wg.Add(1)

	go func() {
		defer w.WorkerPool.wg.Done()

		// 注册工作者的任务队列到工作者池
		w.WorkerPool.workerQueue <- w.TaskQueue

		for {
			select {
			case taskWrapper := <-w.TaskQueue:
				// 处理任务
				w.processTask(taskWrapper)

				// 任务处理完成后，将工作者的任务队列重新注册到工作者池
				select {
				case w.WorkerPool.workerQueue <- w.TaskQueue:
				case <-w.QuitChan:
					return
				}

			case <-w.QuitChan:
				// 收到退出信号
				logger.Info("Worker stopping", zap.String("worker_id", w.ID))
				return

			case <-w.Context.Done():
				// 上下文取消
				logger.Info("Worker context cancelled", zap.String("worker_id", w.ID))
				return
			}
		}
	}()
}

// Stop 停止工作者
func (w *Worker) Stop() {
	close(w.QuitChan)
	w.Cancel()
}

// processTask 处理任务
func (w *Worker) processTask(wrapper *TaskWrapper) {
	task := wrapper.Task

	// 更新任务状态
	task.Status = model.TaskRunning
	startTime := time.Now()
	task.StartedAt = &startTime

	// 通知任务状态更新
	w.WorkerPool.notifyTaskUpdate(task)

	// 执行任务
	err := wrapper.TaskFunc(task)
	if err != nil {
		task.Status = model.TaskFailed
		task.Error = err.Error()
	} else {
		task.Status = model.TaskCompleted
	}

	completionTime := time.Now()
	task.CompletedAt = &completionTime

	// 通知任务完成
	w.WorkerPool.notifyTaskUpdate(task)

	// 清理上下文
	wrapper.Cancel()
}

// generateWorkerID 生成工作者ID
func generateWorkerID() string {
	return fmt.Sprintf("worker-%d", time.Now().UnixNano())
}