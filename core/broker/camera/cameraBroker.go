package camera

import (
	"errors"
	"fmt"
	"pull2push/core/broadcast"
	"sync"
)

// ====================== CameraBroker ======================

// CameraBroker 广播器
type CameraBroker struct {
	mutex sync.Mutex

	//  负责协调所有的拉流链接 FLVStreamBroker 直接的工作
	broadcastMap map[string]broadcast.Broadcaster // key为直播房间号

	mu sync.Mutex
}

func NewCameraBroker() *CameraBroker {
	return &CameraBroker{
		broadcastMap: make(map[string]broadcast.Broadcaster),
	}
}

func (cb *CameraBroker) AddBroadcaster(broadcastKey string, b broadcast.Broadcaster) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.broadcastMap[broadcastKey] = b
}

func (cb *CameraBroker) RemoveBroadcaster(broadcastKey string) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if _, ok := cb.broadcastMap[broadcastKey]; !ok {
		return
	}

	delete(cb.broadcastMap, broadcastKey)
}

// FindBroker 查询 Broker
func (cb *CameraBroker) FindBroadcaster(broadcastKey string) (broadcast.Broadcaster, error) {
	if val, ok := cb.broadcastMap[broadcastKey]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("未找到 %s 对应的 Broadcaster", broadcastKey))
}
