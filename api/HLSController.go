package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pull2push/api/base"
	hlsBroker "pull2push/core/broker/hls"
	"pull2push/service"
)

// HLSController 处理所有流相关的请求
type HLSController struct {
	*base.BaseController
	hlsService *service.HLSService
}

// NewHLSController 创建一个新的 HLSController
func NewHLSController(base *base.BaseController, hlsBrokerPool *hlsBroker.HLSBroker) *HLSController {
	return &HLSController{
		BaseController: base,
		hlsService:     &service.HLSService{HLSBrokerPool: hlsBrokerPool},
	}
}

// LiveHLS 处理 hls 的拉流转推
func (hc HLSController) LiveHLS(c *gin.Context) {

	broadcasterKey := c.Param("broadcasterKey")
	clientId := c.Param("clientId")

	err := hc.hlsService.LiveHLS(c, broadcasterKey, clientId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"msg":  "请求成功",
	//})
}
