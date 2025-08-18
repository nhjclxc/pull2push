package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	hlsBroadcast "pull2push/core/broadcast/hls"
	hlsBroker "pull2push/core/broker/hls"
	hlsClient "pull2push/core/client/hls"
	"strings"
)

// HLSService  Service 层
type HLSService struct {
	HLSBrokerPool *hlsBroker.HLSBroker

	BroadcasterKey string // 这个客户端的直播房间的唯一编号
	ClientId       string // 这个客户端的id
}

// ---------- HTTP 服务 ----------

// LiveHLS 处理 hls 的拉流转推
func (hs HLSService) LiveHLS(c *gin.Context, broadcasterKey, clientId string) error {
	hs.BroadcasterKey = broadcasterKey
	hs.ClientId = clientId

	findBroadcaster, err := hs.HLSBrokerPool.FindBroadcaster(broadcasterKey)
	if err != nil {
		return errors.New("直播不存在！！！" + err.Error())
	}

	findBroadcasterTemp, _ := findBroadcaster.(*hlsBroadcast.HLSBroadcaster)

	filepath := c.Param("filepath")
	//  "xxx/index.m3u8" 结尾的就是第一次请求，这时通过 HandleIndex 接口第一次返回本地缓存的数据片给前端使用
	if strings.HasSuffix(filepath, "/index.m3u8") {

		hlsLiveClient, err := hlsClient.NewHLSLiveClient(c, broadcasterKey, clientId, findBroadcasterTemp.ClientCloseSig, findBroadcasterTemp.BroadcasterCloseSig)
		if err != nil {
			return errors.New("客户端创建失败！！！" + err.Error())
		}
		findBroadcasterTemp.AddLiveClient(clientId, hlsLiveClient)

		// 第一次链接，返回最新的直播数据分片
		hlsLiveClient.HandleIndex(c.Writer, c.Request, findBroadcasterTemp)

		return nil
	}
	// xxx/2689.ts 表示此时前端不是第一次掉接口了，是来获取缓存数据分片的，那么通过 HandleSegment 来下载数据分片

	liveClient, err := findBroadcasterTemp.FindLiveClient(clientId)
	if err != nil {
		return errors.New("为查询到对应的客户端！！！" + err.Error())
	}
	hlsLiveClient, _ := liveClient.(*hlsClient.HLSLiveClient)
	// 返回本地缓存的数据分片
	hlsLiveClient.HandleSegment(c.Writer, c.Request, findBroadcasterTemp)
	// /live/hls/test-hls/c91b431e-ba21-47c9-8649-a05ce2490838/index.m3u8
	return nil

}
