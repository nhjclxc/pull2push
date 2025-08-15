package base

import (
	"pull2push/event"
	"pull2push/logger"
	"pull2push/resource"
	"sync"
)

// BaseController 包含所有 controller 共用的依赖
type BaseController struct {
	Resource      *resource.Resource
	eventBus      *event.EventBus             // 新增事件总线
	responseMap   map[string]chan interface{} // 请求-响应映射
	mu            sync.Mutex                  // 保护responseMap的互斥锁
	eventHandlers []chan event.Event          // 所有订阅的事件通道
	shutdownCh    chan struct{}               // 通知事件处理器退出
}

func NewBaseController(res *resource.Resource) *BaseController {
	if res == nil {
		panic("res cannot be nil")
	}

	controller := &BaseController{
		Resource:      res,
		responseMap:   make(map[string]chan interface{}),
		eventHandlers: make([]chan event.Event, 0),
		shutdownCh:    make(chan struct{}),
	}
	return controller
}

// SetEventBus 设置事件总线
func (c *BaseController) SetEventBus(bus *event.EventBus) {
	c.eventBus = bus
	c.setupEventHandlers()

	// 订阅系统关闭事件
	shutdownCh := c.eventBus.Subscribe(event.SystemShutdown)
	c.eventHandlers = append(c.eventHandlers, shutdownCh)
	go func() {
		for {
			select {
			case <-shutdownCh:
				logger.Info("BaseController received system shutdown event")
				c.Cleanup()
				return
			case <-c.shutdownCh:
				return
			}
		}
	}()
}

// 设置事件处理
func (c *BaseController) setupEventHandlers() {

}

// Cleanup 清理所有资源
func (c *BaseController) Cleanup() {
	logger.Info("BaseController cleaning up resources")

	// 关闭关闭通道，通知所有处理器退出
	close(c.shutdownCh)

	// 取消所有事件订阅
	if c.eventBus != nil {
		for _, ch := range c.eventHandlers {
			c.eventBus.Unsubscribe(ch)
		}
	}

	// 清理响应映射
	c.mu.Lock()
	for requestID, respCh := range c.responseMap {
		close(respCh)
		delete(c.responseMap, requestID)
	}
	c.mu.Unlock()

	c.eventHandlers = nil

	logger.Info("BaseController cleanup completed")
}
