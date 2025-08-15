package event

import (
	"sync"
)

type EventType string

// 事件类型定义
const (
	SystemSetUp    EventType = "SystemSetUp"
	SystemShutdown EventType = "SystemShutdown"
)

// Event 事件结构体定义
type Event struct {
	Type    EventType
	Payload any
}

type EventBus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan Event),
	}
}

func (b *EventBus) Subscribe(eventType EventType) chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 1)
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	return ch
}

func (b *EventBus) SubscribeMultiple(eventTypes []EventType) chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 10)
	for _, eventType := range eventTypes {
		b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	}
	return ch
}

func (b *EventBus) Unsubscribe(ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for eventType, channels := range b.subscribers {
		for i, subCh := range channels {
			if subCh == ch {
				b.subscribers[eventType] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
	}
	close(ch)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[event.Type]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
			}
		}
	}
}
