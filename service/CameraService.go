package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pull2push/core/broadcast"
	cameraBroadcast "pull2push/core/broadcast/camera"
	cameraBroker "pull2push/core/broker/camera"
	cameraClient "pull2push/core/client/camera"
)

// CameraService 摄像头直播推流 Service 层
type CameraService struct {
	CameraBrokerPool *cameraBroker.CameraBroker
}

// ExecutePush ==================== HTTP ====================
// ExecutePush 处理摄像头推上来的流数据
func (cs *CameraService) ExecutePush(c *gin.Context, broadcasterKey string) {

	findBroadcaster, err := cs.CameraBrokerPool.FindBroadcaster(broadcasterKey)
	if err != nil {
		fmt.Printf("未找到对应的广播器 %s \n", broadcasterKey)
		return
	}

	// 开始不断接收推流
	findBroadcaster.PullLoop(broadcast.BroadcasterOptional{GinContext: c})

	cs.CameraBrokerPool.RemoveBroadcaster(broadcasterKey)
}

// ExecutePull 处理每一个链接上来的客户端的推流
func (cs *CameraService) ExecutePull(c *gin.Context, broadcasterKey, clientId string) {

	findBroadcaster, err := cs.CameraBrokerPool.FindBroadcaster(broadcasterKey)
	if err != nil {
		fmt.Printf("未找到对应的广播器 %s \n", broadcasterKey)
		return
	}

	findBroadcasterTemp, _ := findBroadcaster.(*cameraBroadcast.CameraBroadcaster)
	client, err := cameraClient.NewCameraLiveClient(c, broadcasterKey, clientId, findBroadcasterTemp.ClientCloseSig, findBroadcasterTemp.BroadcasterCloseSig)
	if err != nil {
		fmt.Println("NewCameraLiveClient 创建失败：", err)
		return
	}
	fmt.Println("NewCameraLiveClient 创建成功：clientId = ", clientId)
	findBroadcasterTemp.AddLiveClient(clientId, client)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	for {
		select {
		case pkt, ok := <-client.GetDataChan():
			if !ok {
				return
			}
			_, err := c.Writer.Write(pkt)
			if err != nil {
				return
			}
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
