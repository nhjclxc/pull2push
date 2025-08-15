package application

import (
	"context"
	"fmt"
	"pull2push/config"
	"pull2push/event"
	"pull2push/logger"
	"pull2push/resource"
	"sync"
)

type Service interface {
	Start(ctx context.Context) error
	Stop() error

	SetEventBus(bus *event.EventBus)
	SetResources(res *resource.Resource)

	Name() string
}

// ServiceManager 管理所有服务的生命周期
type ServiceManager struct {
	services   map[string]Service
	eventBus   *event.EventBus
	resource   *resource.Resource
	mu         sync.Mutex
	startOrder []string // 记录服务启动顺序
}

// NewServiceManager 创建服务管理器
func NewServiceManager(config *config.Config) (*ServiceManager, error) {
	// 初始化日志
	if err := logger.InitLogger(&config.Common.Log); err != nil {
		return nil, fmt.Errorf("init logger failed: %w", err)
	}

	// 初始化资源
	r, err := resource.NewResource(config)
	if err != nil {
		return nil, err
	}

	return &ServiceManager{
		services:   make(map[string]Service),
		eventBus:   event.NewEventBus(),
		resource:   r,
		startOrder: make([]string, 0),
	}, nil
}

// GetResource 获取资源管理器
func (m *ServiceManager) GetResource() *resource.Resource {
	return m.resource
}

// GetEventBus 获取事件总线
func (m *ServiceManager) GetEventBus() *event.EventBus {
	return m.eventBus
}

// AddService 添加服务
func (m *ServiceManager) AddService(s Service) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := s.Name()
	if _, exists := m.services[name]; exists {
		return fmt.Errorf("service with name %s already exists", name)
	}

	// 设置事件总线和资源
	s.SetEventBus(m.eventBus)
	s.SetResources(m.resource)

	m.services[name] = s
	logger.Info("Service added", "name", name)
	return nil
}

// GetService 获取特定服务实例
func (m *ServiceManager) GetService(name string) (Service, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	service, ok := m.services[name]
	return service, ok
}

// StartAll 启动所有服务
func (m *ServiceManager) StartAll(ctx context.Context) error {
	// 按顺序启动服务
	for _, service := range m.services {
		logger.Info("Starting service", "name", service.Name())
		if err := service.Start(ctx); err != nil {
			return fmt.Errorf("failed to start service %s: %w", service.Name(), err)
		}
		logger.Info("Service started successfully", "name", service.Name())
	}

	return nil
}

// StopAll 停止所有服务
func (m *ServiceManager) StopAll() {
	// 发布系统关闭事件
	m.eventBus.Publish(event.Event{
		Type: event.SystemShutdown,
	})

	logger.Info("System shutdown event published")
	for _, service := range m.services {
		logger.Info("Stopping service", "name", service.Name())
		if err := service.Stop(); err != nil {
			logger.Error("Error stopping service", "name", service.Name(), "error", err)
		} else {
			logger.Info("Service stopped successfully", "name", service.Name())
		}
	}
	// 关闭资源
	if err := m.resource.Close(); err != nil {
		logger.Error("Error closing resources", "error", err)
	}

	logger.Info("All services stopped and resources cleaned up")
}
