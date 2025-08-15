package broadcast

import (
	"errors"
	"fmt"
	"pull2push/core/broadcast"
	"sync"
)

// ====================== HLSBroker ======================

// HLSBroker 广播器
type HLSBroker struct {
	mutex sync.Mutex

	//  负责协调所有的 broadcast.Broadcaster 直接的工作
	broadcastMap map[string]broadcast.Broadcaster // key为直播房间号

}

func NewHLSBroker() *HLSBroker {
	return &HLSBroker{
		broadcastMap: make(map[string]broadcast.Broadcaster),
	}
}

func (hb *HLSBroker) AddBroadcaster(broadcastKey string, b broadcast.Broadcaster) {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	hb.broadcastMap[broadcastKey] = b
}

func (hb *HLSBroker) RemoveBroadcaster(broadcastKey string) {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	if _, ok := hb.broadcastMap[broadcastKey]; !ok {
		return
	}

	delete(hb.broadcastMap, broadcastKey)
}

// FindBroker 查询 Broker
func (fb *HLSBroker) FindBroadcaster(broadcastKey string) (broadcast.Broadcaster, error) {
	if val, ok := fb.broadcastMap[broadcastKey]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("未找到 %s 对应的 Broadcaster", broadcastKey))
}
