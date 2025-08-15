package flv

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ====================== HLSLiveClient ======================

// HLSLiveClient 每一个前端页面有持有一个客户端对象
type HLSLiveClient struct {
	BrokerKey string      // 这个客户端的直播房间的唯一编号
	ClientId  string      // 这个客户端的id
	DataCh    chan []byte // 这个客户端的一个只写通道

	// http连接相关
	httpCloseSig        <-chan struct{} // 当这个请求被客户端主动被关闭时触发
	httpRequestCloseSig <-chan struct{} // 当这个请求被客户端主动被关闭时触发

	// 父级 Broker相关的内容
	clientCloseSig chan<- string   // broker通过该信道监听客户端离线 【仅发送】
	brokerCloseSig <-chan struct{} // broker被关闭时，同时通知客户端关闭 【仅接收】
}

func NewHLSLiveClient(c *gin.Context, brokerKey, clientId string, clientCloseSig chan<- string, brokerCloseSig <-chan struct{}) (*HLSLiveClient, error) {

	hlc := HLSLiveClient{
		BrokerKey:           brokerKey,
		ClientId:            clientId,
		httpCloseSig:        c.Done(),
		httpRequestCloseSig: c.Request.Context().Done(),
		clientCloseSig:      clientCloseSig,
		brokerCloseSig:      brokerCloseSig,
	}

	fmt.Println("HLS 客户端连接成功 ClientId = ", clientId)

	// 开启状态监听
	go hlc.Listen()

	return &hlc, nil
}

func (hlc *HLSLiveClient) Listen() {

	for {
		select {
		case <-hlc.httpCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("hlc.httpCloseSig 收到客户端关闭信号，退出循环 ", hlc.ClientId)
			//
			//// when client closes, remove it
			//hlc.clientCloseSig <- hlc.ClientId

			return
		case <-hlc.httpRequestCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("<-hlc.httpRequestCloseSig 收到客户端关闭信号，退出循环 ", hlc.ClientId)
			//
			//// when client closes, remove it
			//hlc.clientCloseSig <- hlc.ClientId

			return

		case <-hlc.brokerCloseSig:
			return
		}
	}

}

func (hlc *HLSLiveClient) Broadcast(data []byte) {

}

// GetDataChan 获取当前客户端的写通道
func (hlc *HLSLiveClient) GetDataChan() chan []byte {
	return hlc.DataCh
}
