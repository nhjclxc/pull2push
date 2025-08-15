package flv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"pull2push/core/broadcast"
	"pull2push/core/client"
	"sync"
	"time"
)

// ====================== FLVBroadcaster ======================

// FLVBroadcaster 每个 直播地址 用一个 Broker 管理，里面管理了多个当前直播链接的客户端
type FLVBroadcaster struct {
	BroadcasterKey string      // 直播房间的唯一编号
	UpstreamURL    string      // 直播房间的上游拉流地址
	DataCh         chan []byte // 上游拉流缓存的数据
	flvParser      *FLVParser  // 创建FLV解析器 (启用调试模式)
	HeaderMutex    sync.RWMutex
	HeaderBytes    []byte
	HeaderParsed   bool

	// 状态控制相关
	BroadcasterCloseSig chan broadcast.BROADCAST_CLOSE_TYPE // 控制当前这个直播是否被关闭
	stopSig             chan struct{}                       // 控制当前这个直播是否被关闭
	once                sync.Once

	// 客户端相关
	clientMutex    sync.Mutex                   // 客户端的异步操作控制器
	clientMap      map[string]client.LiveClient // map[clientId]LiveClient 存储这个broker里面所有的客户端
	ClientCloseSig chan string                  // 客户端关闭信号，当客户端主动关闭通知时，该信道被触发，输出的字符串为关闭的客户端编号clientId

}

func NewFLVBroadcaster(broadcasterKey, upstreamURL string) *FLVBroadcaster {
	b := FLVBroadcaster{
		BroadcasterKey: broadcasterKey,
		UpstreamURL:    upstreamURL,
		DataCh:         make(chan []byte, 4096), // 带缓冲，缓存 4096 个数据包
		clientMap:      make(map[string]client.LiveClient),
		stopSig:        make(chan struct{}),
	}

	// start pulling loop
	go b.PullLoop(broadcast.BroadcasterOptional{})
	fmt.Printf("\n FLVBroadcaster = %#v \n", &b)

	return &b
}

func (fb *FLVBroadcaster) AddLiveClient(clientId string, client client.LiveClient) {
	fb.clientMutex.Lock()
	defer fb.clientMutex.Unlock()

	fb.clientMap[clientId] = client
}

// RemoveLiveClient 移除客户端
func (fb *FLVBroadcaster) RemoveLiveClient(clientId string) {
	fb.clientMutex.Lock()
	defer fb.clientMutex.Unlock()
	liveClient := fb.clientMap[clientId]
	if liveClient == nil {
		return
	}
	delete(fb.clientMap, clientId)

	//// 如果没有客户端并且想释放 broker，可关闭 stopCh 让 PullLoop 停止（本示例保留 broker，防止频繁断开上游）
	//if remaining == 0 {
	//	// optionally stop pulling after idle timeout. For simplicity we keep running.
	//}
}

// FindLiveClient 查询 LiveClient
func (fb *FLVBroadcaster) FindLiveClient(clientId string) (client.LiveClient, error) {
	if val, ok := fb.clientMap[clientId]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("未找到 %s 对应的 LiveClient", clientId))
}

// UpdateSourceURL 支持切换直播原地址
func (fb *FLVBroadcaster) UpdateSourceURL(newSourceURL string) {}

// ListenStatus 监听当前直播的必要状态
func (fb *FLVBroadcaster) ListenStatus() {}

// PullLoop 持续去服务端拉流
func (fb *FLVBroadcaster) PullLoop(bo broadcast.BroadcasterOptional) {
	backoff := time.Second
	for {
		log.Println("dial upstream", fb.UpstreamURL)
		req, _ := http.NewRequest("GET", fb.UpstreamURL, nil)
		// add headers typical for FLV
		req.Header.Set("User-Agent", "Go-Relay-Flv/1.0")
		req.Header.Set("Accept", "*/*")
		client := &http.Client{
			Timeout: 0, // streaming
			Transport: &http.Transport{
				// keep-alive
				DialContext: (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			},
		}
		// 拉流
		resp, err := client.Do(req)
		// 失败重试
		if err != nil {
			time.Sleep(backoff)
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Println("upstream bad status:", resp.Status)
			resp.Body.Close()
			time.Sleep(backoff)
			continue
		}

		//fb.doFlvParse(resp.Body)

		// 成功连接，重置 backoff
		backoff = time.Second

		// 读取本次拉到的流数据，并且进行数据分发
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				// copy bytes to avoid race
				cp := make([]byte, n)
				copy(cp, buf[:n])
				fb.Broadcast2LiveClient(cp)
			}
			if err != nil {
				if err == io.EOF {
					log.Println("upstream EOF, reconnecting", "退出拉流过程")
					break
				} else if "unexpected EOF" == err.Error() {
					log.Println("upstream read error:", "上游直播关闭", "退出拉流过程")
					break
				} else {
					log.Println("upstream read error:", err)
				}
				resp.Body.Close()
				break
			}
		}

		// 如果 stop 信号被触发，可以退出（此实现未触发 stop）
		select {
		case <-fb.stopSig:
			return
		default:
		}

		// small backoff before reconnect
		time.Sleep(500 * time.Millisecond)
	}
}

func (fb *FLVBroadcaster) doFlvParse(body io.ReadCloser) {

	ctx := context.Background()
	_ = ctx
	// 创建FLV解析器 (启用调试模式)
	fb.flvParser = NewFLVParser(true)

	fmt.Println("开始解析FLV流...")
	fmt.Println("等待解析视频和音频标签...")

	// 使用带超时的方法解析FLV头部和初始标签
	err := fb.flvParser.ParseInitialTags(ctx, body)
	if err != nil {
		fmt.Printf("解析FLV失败: %v\n", err)

		panic(err)
	}
	fb.flvParser.PrintRequiredTags()

	// 获取必要的FLV头和标签字节
	headerBytes, err := fb.flvParser.GetRequiredTagsBytes()
	if err != nil {
		body.Close()
		fmt.Println("获取FLV头部字节失败: %v", err)
	}

	// 保存头部字节供后续使用
	fb.HeaderMutex.Lock()
	fb.HeaderBytes = headerBytes
	fb.HeaderParsed = true
	fb.HeaderMutex.Unlock()

	slog.Info("FLV头部解析完成，等待客户端连接开始转发", "bytes", len(headerBytes))

}
func (fb *FLVBroadcaster) Broadcast2LiveClient(data []byte) {
	//log.Println("FLVBroadcaster.Broadcast.BrokerKey:", fb.BrokerKey)
	if len(fb.clientMap) == 0 {
		return
	}
	fb.clientMutex.Lock()
	// 复制 client list to avoid holding lock during send
	clients := make([]client.LiveClient, 0, len(fb.clientMap))
	for _, c := range fb.clientMap {
		clients = append(clients, c)
	}
	fb.clientMutex.Unlock()

	// 广播数据 send non-blocking (drop if client's channel full)
	for _, c := range clients {
		c.Broadcast(data)
		////fmt.Printf("send non-blocking %#v \n\n", c.dataCh)
		//select {
		//case c.dataCh <- data:
		//	//fmt.Println("发送数据成功")
		//default:
		//	// 客户端慢，丢包并继续
		//}
	}
}
