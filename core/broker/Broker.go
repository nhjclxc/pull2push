package broker

import "pull2push/core/broadcast"

// ====================== Broker ======================
// Broker ===>> Broadcaster ===>> LiveClient
// 每一个直播类型一个 Broker 对象，一个 Broker 对象管理了多个 Broadcaster，一个 Broadcaster 就是一个直播间，
//一个 Broadcaster 又管理了多个 LiveClient，一个 LiveClient 就是一个客户端页面

// Broker 直播类型管理器
type Broker interface {

	// AddBroadcaster 添加 Broker
	AddBroadcaster(brokerKey string, bc broadcast.Broadcaster)

	// RemoveBroadcaster 移除 Broker
	RemoveBroadcaster(brokerKey string)

	// FindBroadcaster 查询 Broker
	FindBroadcaster(brokerKey string) (broadcast.Broadcaster, error)
}
