package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pull2push/api/base"
	cameraBroker "pull2push/core/broker/camera"
	"pull2push/service"
)

// CameraController 处理所有流相关的请求
type CameraController struct {
	*base.BaseController
	cameraService *service.CameraService
}

// NewCameraController 创建一个新的 CameraController
func NewCameraController(base *base.BaseController, cameraBrokerPool *cameraBroker.CameraBroker) *CameraController {
	return &CameraController{
		BaseController: base,
		cameraService:  &service.CameraService{CameraBrokerPool: cameraBrokerPool},
	}
}

// ExecutePush 处理摄像头推上来的流数据
func (cc *CameraController) ExecutePush(c *gin.Context) {
	//if !strings.HasPrefix(c.GetHeader("Content-Type"), "video/x-flv") {
	//	c.String(http.StatusBadRequest, "Content-Type must be video/x-flv")
	//	return
	//}
	broadcasterKey := c.Param("broadcasterKey")

	cc.cameraService.ExecutePush(c, broadcasterKey)

}

// ExecutePull 处理每一个链接上来的客户端的推流
func (cc *CameraController) ExecutePull(c *gin.Context) {

	c.Header("Content-Type", "video/x-flv")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	broadcasterKey := c.Param("broadcasterKey")
	clientId := c.Param("clientId")

	cc.cameraService.ExecutePull(c, broadcasterKey, clientId)
}
