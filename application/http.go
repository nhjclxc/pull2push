package application

import (
	"context"
	"errors"
	"net/http"
	"pull2push/api"
	"pull2push/api/base"
	"pull2push/config"
	cameraBroadcast "pull2push/core/broadcast/camera"
	flvBroadcast "pull2push/core/broadcast/flv"
	hlsBroadcast "pull2push/core/broadcast/hls"
	cameraBroker "pull2push/core/broker/camera"
	flvBroker "pull2push/core/broker/flv"
	hlsBroker "pull2push/core/broker/hls"
	"pull2push/event"
	"pull2push/logger"
	"pull2push/middleware"
	"pull2push/resource"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPService struct {
	server         *http.Server
	engine         *gin.Engine
	config         *config.Config
	resources      *resource.Resource
	eventBus       *event.EventBus
	baseController *base.BaseController
	shutdownCh     chan struct{}      // 关闭通道
	eventHandlers  []chan event.Event // 用于跟踪所有事件处理 goroutine
	mu             sync.Mutex         // 保护 eventHandlers

	cameraBrokerPool *cameraBroker.CameraBroker
	flvBrokerPool    *flvBroker.FLVBroker
	hlsBrokerPool    *hlsBroker.HLSBroker
}

// NewHTTPService 创建 HTTP 服务
func NewHTTPService(res *resource.Resource) *HTTPService {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// 添加 CORS 中间件
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "android-device-id", "Content-Type", "Accept", "Authorization", "X-Token", "Device-Id", "request-time", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	service := &HTTPService{
		engine:           engine,
		config:           res.Config,
		resources:        res,
		baseController:   base.NewBaseController(res),
		cameraBrokerPool: cameraBroker.NewCameraBroker(),
		flvBrokerPool:    flvBroker.NewFLVBroker(),
		hlsBrokerPool:    hlsBroker.NewHLSBroker(),
	}
	return service
}

// Name 返回服务名称
func (s *HTTPService) Name() string {
	return "http_service"
}

// Dependencies 返回依赖的其他服务
func (s *HTTPService) Dependencies() []string {
	// HTTP 服务依赖于 StreamHub 服务，通过事件通信
	return []string{"stream_hub"}
}

// 注册事件处理器，用于跟踪
func (s *HTTPService) registerEventHandler(ch chan event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventHandlers = append(s.eventHandlers, ch)
}

func (s *HTTPService) setupControllers() {
	// 创建基础控制器
	s.baseController = base.NewBaseController(s.resources)
	// 设置事件总线
	s.baseController.SetEventBus(s.eventBus)

	// 创建其他控制器
	// ...
}

func (s *HTTPService) setupRoutes() {

	if s.resources.Config.Common.Debug {
		s.engine.Use(gin.Logger())
	}
	s.engine.Use(middleware.GlobalPanicRecovery())

	s.engine.GET("/api/ping", func(c *gin.Context) { c.JSON(http.StatusOK, "pong") })

	/*
	   使用ffmpeg推流：ffmpeg -re -i demo.flv -c copy -f flv rtmp://192.168.203.182/live/livestream

	   拉流
	   ● RTMP (by VLC): rtmp://192.168.203.182/live/livestream
	   ● H5(HTTP-FLV): http://192.168.203.182:8080/live/livestream.flv
	   ● H5(HLS): http://192.168.203.182:8080/live/livestream.m3u8

	   	ffmpeg -re -i demo.flv \
	   	    -c:v libx264 -preset veryfast -tune zerolatency \
	   	    -g 25 -keyint_min 25 \
	   	    -c:a aac -ar 44100 -b:a 128k \
	   	    -f flv rtmp://192.168.203.182/live/livestream
	*/

	flvPull2pushRouter := s.engine.Group("/api/live/flv")
	{

		flvBroadcasterKey := "test-flv"
		flvUpstreamURL := "http://192.168.203.182:8080/live/livestream.flv"

		testFlvBroadcast := flvBroadcast.NewFLVBroadcaster(flvBroadcasterKey, flvUpstreamURL)
		s.flvBrokerPool.AddBroadcaster(flvBroadcasterKey, testFlvBroadcast)
		_ = flvPull2pushRouter

		flvController := api.NewFLVController(s.baseController, s.flvBrokerPool)

		// 使用ffmpeg推流：ffmpeg -re -i demo.flv -c copy -f flv rtmp://192.168.203.182/live/livestream

		// flv启动顺序：1、go服务器，2、前端页面，3、使用ffmpeg推流

		// http://localhost:8080/api/live/flv/test-flv/729119c9-0711-4ef8-b60e-6c2dca5b1a11
		// http://localhost:8080/api/live/flv/test-flv/123
		flvPull2pushRouter.GET("/:broadcasterKey/:clientId", flvController.LiveFlv)
	}

	hlsPull2pushRouter := s.engine.Group("/api/live/hls")
	{

		ctx, _ := context.WithCancel(context.Background())
		//defer cancel()

		hlsBroadcastKey := "test-hls"
		hlsUpstreamURL := "http://192.168.203.182:8080/live/livestream.m3u8"

		var hlsBroadcastTemp *hlsBroadcast.HLSBroadcaster = hlsBroadcast.NewHLSBroadcaster(ctx, hlsBroadcastKey, hlsUpstreamURL, "", 3)
		s.hlsBrokerPool.AddBroadcaster(hlsBroadcastKey, hlsBroadcastTemp)

		hlsController := api.NewHLSController(s.baseController, s.hlsBrokerPool)

		// hls启动顺序：1、使用ffmpeg推流，2、go服务器，3、前端页面
		// 使用ffmpeg推流：ffmpeg -re -i demo.flv -c copy -f flv rtmp://192.168.203.182/live/livestream

		// hls要提供两个接口，一个是 index.m3u8用于客户端第一次调用的时候获取最新数据分片消息的，有助于第二个接口来获取最新的分片数据
		// 一个是 类似 2689.ts 的接口，用于给客户端请求具体的流数据
		// http://localhost:8080/live/hls/:brokerKey/:clientId/index.m3u8
		// http://localhost:8080/live/hls/:brokerKey/:clientId/2689.ts
		hlsPull2pushRouter.GET("/:broadcasterKey/:clientId/*filepath", hlsController.LiveHLS)
	}

	cameraPull2pushRouter := s.engine.Group("/api/live/camera")
	{

		broadcasterKey := "test-camera"
		testCameraCameraBroadcast := cameraBroadcast.NewCameraBroadcaster(broadcasterKey, 150)
		s.cameraBrokerPool.AddBroadcaster(broadcasterKey, testCameraCameraBroadcast)

		cameraController := api.NewCameraController(s.baseController, s.cameraBrokerPool)

		// camera启动顺序：1、go服务器，2、前端页面，3、使用ffmpeg推流

		// ffmpeg -f avfoundation -framerate 30 -video_size 640x480 -i "0:0" -vcodec libx264 -preset veryfast -tune zerolatency -g 30 -acodec aac -ar 44100 -ac 2 -f flv "http://127.0.0.1:8080/api/live/camera/ingest/test-camera"
		// http://127.0.0.1:8080/api/live/camera/ingest/test-camera
		// 摄像头推流
		cameraPull2pushRouter.POST("/ingest/:broadcasterKey", cameraController.ExecutePush)

		// http://127.0.0.1:8080/api/live/camera/test-camera/123
		// 客户端拉流
		cameraPull2pushRouter.GET("/:broadcasterKey/:clientId", cameraController.ExecutePull)
	}

}

func (s *HTTPService) Start(ctx context.Context) error {
	s.setupRoutes()
	s.setupControllers()
	s.server = &http.Server{
		Addr:    ":" + s.config.HTTP.Port,
		Handler: s.engine,
	}
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// 订阅相关事件
	s.subscribeToEvents()
	logger.Info("Starting HTTP server on port", "port", s.config.HTTP.Port)
	return nil
}

// 订阅事件
func (s *HTTPService) subscribeToEvents() {
	// 订阅系统关闭事件
	shutdownCh := s.eventBus.Subscribe(event.SystemShutdown)
	s.registerEventHandler(shutdownCh)
	go func() {
		for {
			select {
			case <-shutdownCh:
				logger.Info("HTTP service received system shutdown event")
				// 立即关闭 HTTP 服务器
				if s.server != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					s.server.Shutdown(ctx)
				}
			case <-s.shutdownCh:
				return
			}
		}
	}()
}

func (s *HTTPService) Stop() error {
	// 关闭通道，通知所有事件处理 goroutine 退出
	close(s.shutdownCh)

	// 关闭 HTTP 服务器
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down HTTP server", "error", err)
			return err
		}
	}

	// 取消订阅所有事件
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.eventHandlers {
		s.eventBus.Unsubscribe(ch)
	}
	s.eventHandlers = nil
	logger.Info("HTTP service stopped successfully")
	return nil
}

func (s *HTTPService) SetEventBus(bus *event.EventBus) {
	s.eventBus = bus
}

func (s *HTTPService) SetResources(res *resource.Resource) {
	s.resources = res
}
