package broadcast

import (
	"github.com/gin-gonic/gin"
	"pull2push/core/client"
)

// ====================== Broadcaster ======================

// Broadcaster 广播器接口
type Broadcaster interface {

	// AddLiveClient 新增客户端
	AddLiveClient(clientId string, liveClient client.LiveClient)

	// RemoveLiveClient 移除客户端
	RemoveLiveClient(clientId string)

	// FindLiveClient 查询 LiveClient
	FindLiveClient(clientId string) (client.LiveClient, error)

	// ListenStatus 监听当前直播的必要状态
	ListenStatus()

	// PullLoop 持续去直播原地址拉流/数据
	PullLoop(BroadcasterOptional)

	// Broadcast2LiveClient 原地址拉取到数据之后广播给客户端
	Broadcast2LiveClient(data []byte)

	// UpdateSourceURL 支持切换直播原地址
	UpdateSourceURL(newSourceURL string)
}

// BroadcasterOptional broker配置选项
type BroadcasterOptional struct {
	GinContext *gin.Context
}

type BROADCAST_CLOSE_TYPE int

const (
	// BrokerStarted 直播开始
	BrokerStarted BROADCAST_CLOSE_TYPE = 1

	// BrokerEnd 直播结束
	BrokerEnd BROADCAST_CLOSE_TYPE = 2

	// BrokerClosed 直播被关闭
	BrokerClosed BROADCAST_CLOSE_TYPE = 3
)
