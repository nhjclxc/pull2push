package client

// ====================== LiveClient ======================

// LiveClient 观看直播的客户端
type LiveClient interface {

	// Broadcast 服务端给客户端推流
	// data 当前接收到直播上游的新数据
	Broadcast(data []byte)

	// Listen 客户端监听器
	Listen()

	// GetDataChan 获取当前客户端的写通道
	GetDataChan() chan []byte
}
