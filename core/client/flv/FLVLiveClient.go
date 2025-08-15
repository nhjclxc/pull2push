package flv

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"pull2push/core/broadcast"
	"runtime/debug"
)

// ====================== FLVLiveClient ======================

// FLVLiveClient 每一个前端页面有持有一个客户端对象
type FLVLiveClient struct {
	BrokerKey string        // 这个客户端的直播房间的唯一编号
	ClientId  string        // 这个客户端的id
	DataCh    chan []byte   // 这个客户端的一个只写通道
	CloseSig  chan struct{} // broker被关闭时，同时通知客户端关闭

	// http连接相关
	httpRequest         *http.Request
	responseWriter      io.Writer
	flusher             http.Flusher
	httpCloseSig        <-chan struct{} // 当这个请求被客户端主动被关闭时触发
	httpRequestCloseSig <-chan struct{} // 当这个请求被客户端主动被关闭时触发

	// 父级 broadcaster 相关的内容
	clientCloseSig      chan<- string                         // broker通过该信道监听客户端离线 【仅发送】
	broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE // broker被关闭时，同时通知客户端关闭 【仅接收】
}

func NewFLVLiveClient(c *gin.Context, brokerKey, clientId string, clientCloseSig chan<- string, broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE) (*FLVLiveClient, error) {
	// 创建一个带缓冲的双向通道，缓冲大小根据需求调节
	dataCh := make(chan []byte, 4096)

	// gin.ResponseWriter 是接口，不能用指针
	var writer io.Writer = c.Writer

	// 断言出 http.Flusher 接口，方便主动刷新数据
	flusher, ok := writer.(http.Flusher)
	if !ok {
		// 不支持刷新，可能无法做到流式推送
		log.Println("ResponseWriter does not support Flusher interface")
	}

	hc := FLVLiveClient{
		BrokerKey:           brokerKey,
		ClientId:            clientId,
		DataCh:              dataCh,
		CloseSig:            make(chan struct{}),
		httpRequest:         c.Request,
		responseWriter:      writer,
		flusher:             flusher,
		httpCloseSig:        c.Done(),
		httpRequestCloseSig: c.Request.Context().Done(),
		clientCloseSig:      clientCloseSig,
		broadcasterCloseSig: broadcasterCloseSig,
	}

	fmt.Println("客户端连接成功 ClientId = ", clientId)

	// 持续监控是否一些控制通道的消息
	go hc.Listen()

	return &hc, nil
}

// Listen 客户端监听器
func (flc *FLVLiveClient) Listen() {

	for {
		select {
		case data, ok := <-flc.DataCh:
			//fmt.Printf("接收到数据 len = %d \n", len(data))
			if ok {

				if flc.responseWriter == nil || flc.flusher == nil {
					log.Println("responseWriter 或 flusher 已经无效，退出")
					return
				}

				_, err := flc.responseWriter.Write(data)
				if err != nil {
					// 写出错，关闭连接
					return
				}
				flc.flusher.Flush()
			}
		case <-flc.httpCloseSig:
			flc.CloseSig <- struct{}{}
			// 收到关闭信号，退出循环
			fmt.Println("flc.httpCloseSig 收到客户端关闭信号，退出循环 ", flc.ClientId)

			// when client closes, remove it
			flc.clientCloseSig <- flc.ClientId

			close(flc.CloseSig)

			return
		case <-flc.httpRequestCloseSig:
			flc.CloseSig <- struct{}{}
			// 收到关闭信号，退出循环
			fmt.Println("<-flc.httpRequestCloseSig 收到客户端关闭信号，退出循环 ", flc.ClientId)

			// when client closes, remove it
			flc.clientCloseSig <- flc.ClientId

			close(flc.CloseSig)

			return
		}
	}
}

// GetDataChan 获取当前客户端的写通道
func (flc *FLVLiveClient) GetDataChan() chan []byte {
	return flc.DataCh
}
func (flc *FLVLiveClient) Broadcast(data []byte) {

	defer func() {
		if err := recover(); err != nil {
			// 这里写自定义日志处理  打印堆栈
			log.Printf("panic recovered: %v: %v\n stack trace:\n%s", err, debug.Stack())
		}
	}()

	if flc.responseWriter == nil || flc.flusher == nil {
		log.Println("responseWriter 或 flusher 已经无效，退出")
		return
	}

	_, err := flc.responseWriter.Write(data)
	if err != nil {
		// 写出错，关闭连接
		return
	}
	flc.flusher.Flush()

}
