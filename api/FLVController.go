package api

import (
	"github.com/gin-gonic/gin"
	"pull2push/api/base"
	flvBroker "pull2push/core/broker/flv"
	"pull2push/service"
)

// FLVController 处理所有流相关的请求
type FLVController struct {
	*base.BaseController
	flvService *service.FLVService
}

// NewFLVController 创建一个新的 FLVController
func NewFLVController(base *base.BaseController, flvBrokerPool *flvBroker.FLVBroker) *FLVController {
	return &FLVController{
		BaseController: base,
		flvService:     &service.FLVService{FLVBrokerPool: flvBrokerPool},
	}
}

// ExecutePush 处理摄像头推上来的流数据
func (fc *FLVController) LiveFlv(c *gin.Context) {
	c.Header("Content-Type", "video/x-flv")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 确保响应缓冲区被刷新
	c.Writer.Flush()

	broadcasterKey := c.Param("broadcasterKey")
	clientId := c.Param("clientId")

	fc.flvService.LiveFlv(c, broadcasterKey, clientId)

}
