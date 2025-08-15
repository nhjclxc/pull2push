package flv

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pull2push/core/broadcast"
)

// ====================== CameraLiveClient ======================

// CameraLiveClient 每一个前端页面有持有一个客户端对象
type CameraLiveClient struct {
	BroadcasterKey string      // 这个客户端的直播房间的唯一编号
	ClientId       string      // 这个客户端的id
	dataCh         chan []byte // 这个客户端的一个只写通道

	// http连接相关
	httpCloseSig        <-chan struct{} // 当这个请求被客户端主动被关闭时触发
	httpRequestCloseSig <-chan struct{} // 当这个请求被客户端主动被关闭时触发

	// 父级 broadcaster 相关的内容
	clientCloseSig      chan<- string                         // broker通过该信道监听客户端离线 【仅发送】
	broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE // broker被关闭时，同时通知客户端关闭 【仅接收】
}

func NewCameraLiveClient(c *gin.Context, broadcasterKey, clientId string, clientCloseSig chan<- string, broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE) (*CameraLiveClient, error) {

	clc := CameraLiveClient{
		BroadcasterKey:      broadcasterKey,
		ClientId:            clientId,
		dataCh:              make(chan []byte, 1024),
		httpCloseSig:        c.Done(),
		httpRequestCloseSig: c.Request.Context().Done(),
		clientCloseSig:      clientCloseSig,
		broadcasterCloseSig: broadcasterCloseSig,
	}

	fmt.Println("HLS 客户端连接成功 ClientId = ", clientId)

	// 开启状态监听
	go clc.Listen()

	return &clc, nil
}

// GetDataChan 获取当前客户端的写通道
func (clc *CameraLiveClient) GetDataChan() chan []byte {
	return clc.dataCh
}

func (clc *CameraLiveClient) Listen() {

	for {
		select {
		case <-clc.httpCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("clc.httpCloseSig 收到客户端关闭信号，退出循环 ", clc.ClientId)
			//
			//// when client closes, remove it
			//clc.clientCloseSig <- clc.ClientId

			return
		case <-clc.httpRequestCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("<-clc.httpRequestCloseSig 收到客户端关闭信号，退出循环 ", clc.ClientId)
			//
			//// when client closes, remove it
			//clc.clientCloseSig <- clc.ClientId

			return
		case <-clc.broadcasterCloseSig:
			return
		}
	}

}

func (clc *CameraLiveClient) Broadcast(data []byte) {

}
