package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	flvBroadcast "pull2push/core/broadcast/flv"
	flvBroker "pull2push/core/broker/flv"
	flvClient "pull2push/core/client/flv"
)

// FLVService FLV拉流转推 Service 层
type FLVService struct {
	FLVBrokerPool *flvBroker.FLVBroker
}

// ---------- HTTP 服务 ----------

// LiveFlv 处理 flv 的拉流转推
func (fs *FLVService) LiveFlv(c *gin.Context, broadcasterKey, clientId string) {

	findBroadcaster, err := fs.FLVBrokerPool.FindBroadcaster(broadcasterKey)
	if err != nil {
		fmt.Printf("未找到对应的广播器 %s \n", broadcasterKey)
		return
	}

	findBroadcasterTemp, _ := findBroadcaster.(*flvBroadcast.FLVBroadcaster)

	// 阻塞客户端
	//<-c.Request.Context().Done()

	//// 或者使用以下逻辑
	c.Stream(func(w io.Writer) bool {

		liveFLVClient, err := flvClient.NewFLVLiveClient(c, broadcasterKey, clientId, findBroadcasterTemp.ClientCloseSig, findBroadcasterTemp.BroadcasterCloseSig)
		if err != nil {
			c.JSON(500, err)
			return false
		}

		findBroadcasterTemp.AddLiveClient(clientId, liveFLVClient)

		// 这里写数据推送逻辑，或者直接阻塞直到连接关闭
		<-c.Request.Context().Done()
		return false
	})

}
