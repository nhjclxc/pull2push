package broadcast

import (
	"errors"
	"fmt"
	"pull2push/core/broadcast"
	"sync"
)

// ====================== FLVBroker ======================

// FLVBroker 广播器
type FLVBroker struct {
	mutex sync.Mutex

	//  负责协调所有的拉流链接 broadcast.Broadcaster 直接的工作
	broadcastMap map[string]broadcast.Broadcaster // key为直播房间号

}

func NewFLVBroker() *FLVBroker {
	return &FLVBroker{
		broadcastMap: make(map[string]broadcast.Broadcaster),
	}
}

func (fb *FLVBroker) AddBroadcaster(broadcastKey string, b broadcast.Broadcaster) {
	fb.mutex.Lock()
	defer fb.mutex.Unlock()

	fb.broadcastMap[broadcastKey] = b
}

func (fb *FLVBroker) RemoveBroadcaster(broadcastKey string) {
	fb.mutex.Lock()
	defer fb.mutex.Unlock()

	if _, ok := fb.broadcastMap[broadcastKey]; !ok {
		return
	}

	delete(fb.broadcastMap, broadcastKey)
}

// FindBroker 查询 Broker
func (fb *FLVBroker) FindBroadcaster(broadcastKey string) (broadcast.Broadcaster, error) {
	if val, ok := fb.broadcastMap[broadcastKey]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("未找到 %s 对应的 Broadcaster", broadcastKey))
}
