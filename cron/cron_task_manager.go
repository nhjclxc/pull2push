package cron

import (
	"context"
	"pull2push/event"
	"pull2push/logger"
	"pull2push/resource"
	"sync"
	"time"
)

// CronTask 定义一个定时任务接口
type CronTask interface {

	// Name 返回任务名称
	Name() string

	// Execute 执行任务
	Execute(ctx context.Context) error

	// GetInterval 返回任务执行间隔
	GetInterval() time.Duration

	// Enable 是否启用
	Enable() bool
}

// CronTaskManager 定时任务管理器
type CronTaskManager struct {
	config        *resource.Resource
	eventBus      *event.EventBus
	stopCh        chan struct{}
	isRunning     bool
	tasks         []CronTask
	mu            sync.Mutex
	eventHandlers []chan event.Event
}

func NewCronTaskManager(res *resource.Resource) *CronTaskManager {
	return &CronTaskManager{
		config:        res,
		stopCh:        make(chan struct{}),
		tasks:         make([]CronTask, 0),
		eventHandlers: make([]chan event.Event, 0),
	}
}

// Name 返回服务名称
func (tm *CronTaskManager) Name() string {
	return "cron_task_manager"
}

// AddTask 向管理器添加一个定时任务
func (tm *CronTaskManager) AddTask(task CronTask) {

	if !task.Enable() {
		logger.Info("CronTask is disabled", "task", task.Name())
		return
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks = append(tm.tasks, task)
	logger.Info("Added task to manager", "task", task.Name(), "interval", task.GetInterval())
}

// 注册事件处理器，用于跟踪
func (tm *CronTaskManager) registerEventHandler(ch chan event.Event) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.eventHandlers = append(tm.eventHandlers, ch)
}

// SetEventBus 设置事件总线
func (tm *CronTaskManager) SetEventBus(bus *event.EventBus) {
	tm.eventBus = bus
}

// SetResources 设置资源
func (tm *CronTaskManager) SetResources(res *resource.Resource) {
	tm.config = res
}

// Start 启动任务管理器服务
func (tm *CronTaskManager) Start(ctx context.Context) error {
	if tm.isRunning {
		return nil
	}

	tm.isRunning = true

	// 订阅系统关闭事件
	tm.subscribeToEvents()

	// 启动所有定时任务
	tm.startAllTasks(ctx)

	logger.Info("TaskManager started")
	return nil
}

// 订阅事件
func (tm *CronTaskManager) subscribeToEvents() {
	// 订阅系统关闭事件
	shutdownCh := tm.eventBus.Subscribe(event.SystemShutdown)
	tm.registerEventHandler(shutdownCh)
	go func() {
		for {
			select {
			case <-shutdownCh:
				logger.Info("TaskManager received system shutdown event")
				// 停止服务
				tm.Stop()
			case <-tm.stopCh:
				return
			}
		}
	}()
}

// Stop 停止任务管理器服务
func (tm *CronTaskManager) Stop() error {
	if !tm.isRunning {
		return nil
	}

	// 发送停止信号
	close(tm.stopCh)
	tm.isRunning = false

	// 取消订阅所有事件
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, ch := range tm.eventHandlers {
		tm.eventBus.Unsubscribe(ch)
	}
	tm.eventHandlers = nil

	logger.Info("TaskManager stopped")
	return nil
}

// 启动所有定时任务
func (tm *CronTaskManager) startAllTasks(ctx context.Context) {
	tm.mu.Lock()
	tasks := make([]CronTask, len(tm.tasks))
	copy(tasks, tm.tasks)
	tm.mu.Unlock()

	for _, task := range tasks {
		// 为每个任务启动一个独立的goroutine
		go tm.runTask(ctx, task)
	}
}

// 运行单个定时任务
func (tm *CronTaskManager) runTask(ctx context.Context, task CronTask) {
	ticker := time.NewTicker(task.GetInterval())
	defer ticker.Stop()

	taskName := task.Name()
	logger.Info("Starting task", "task", taskName, "interval", task.GetInterval())

	tm.executeTask(ctx, task)

	for {
		select {
		case <-ticker.C:
			tm.executeTask(ctx, task)
		case <-tm.stopCh:
			logger.Info("CronTask stopped", "task", taskName)
			return
		case <-ctx.Done():
			logger.Info("Context canceled, stopping task", "task", taskName)
			return
		}
	}
}

// 执行单个任务并处理错误
func (tm *CronTaskManager) executeTask(ctx context.Context, task CronTask) {
	taskName := task.Name()

	// 创建一个带超时的上下文
	// 为每个任务设置一个默认的超时时间，防止任务执行时间过长
	taskCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logger.Info("Executing task", "task", taskName)

	// 执行任务
	err := task.Execute(taskCtx)
	if err != nil {
		logger.Error("CronTask execution failed", "task", taskName, "error", err)
	} else {
		logger.Info("CronTask executed successfully", "task", taskName)
	}
}
