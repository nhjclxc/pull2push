package hls

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"pull2push/core/broadcast"
	"pull2push/core/client"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
整体功能概述
功能目标：
	从多个上游 HLS 流地址拉流（m3u8 + ts），
	缓存最近若干个分片（ring buffer），
	并通过本地 HTTP 服务将这些分片和对应的 m3u8 播放列表转发出去，
	支持多路同时拉取和复用。

多路复用实现
	关键是用全局 map[string]*StreamState 管理多路拉流状态，路由里根据 URL 路径选择对应流。
	每一路独立拉取并缓存，HTTP 根据 URL 返回对应的 m3u8 和分片。
	这样前端可以访问 /live/cam/index.m3u8、/live/news/index.m3u8，不同流完全隔离。

主要流程：拉流 + 缓存 + HTTP 服务


*/

// HLSBroadcaster 每个 直播地址 用一个 Broker 管理，里面管理了多个当前直播链接的客户端
type HLSBroadcaster struct {
	// 直播数据相关
	BrokerKey    string       // 直播房间的唯一编号
	upstreamURL  string       // 直播房间的上游拉流地址
	Variant      string       // 可选：固定选择带宽 id/分辨率（留空自动选最优）
	StreamState0 *StreamState // m3u8数据分片处理器

	// 状态控制相关
	BrokerCloseSig chan struct{} // 控制当前这个直播是否被关闭
	once           sync.Once
	ctx            context.Context

	// 客户端相关
	clientMutex    sync.Mutex                   // 客户端的异步操作控制器
	clientMap      map[string]client.LiveClient // map[clientId]LiveClient 存储这个broker里面所有的客户端
	ClientCloseSig chan string                  // 客户端关闭信号，当客户端主动关闭通知时，该信道被触发，输出的字符串为关闭的客户端编号clientId

}

func NewHLSBroadcaster(ctx context.Context, brokerKey, upstreamURL, variant string, buffer int) *HLSBroadcaster {
	if buffer == 0 {
		buffer = 3
	}
	hmb := HLSBroadcaster{
		BrokerKey:      brokerKey,
		upstreamURL:    upstreamURL,
		Variant:        variant,
		StreamState0:   NewStreamState(buffer),
		clientMap:      make(map[string]client.LiveClient),
		ctx:            ctx,
		BrokerCloseSig: make(chan struct{}),
		ClientCloseSig: make(chan string),
	}

	// 开始持续拉流
	go hmb.PullLoop(broadcast.BroadcasterOptional{})

	// 开启必要的状态监听
	go hmb.ListenStatus()

	return &hmb
}

// ---------- HLS 拉流逻辑 ----------

// resolveURL 处理相对 URI -> 绝对 URL
func resolveURL(base, ref string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	r, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	return u.ResolveReference(r).String(), nil
}

// pickVariant 选择主清单里最合适的变体
func pickVariant(master *m3u8.MasterPlaylist, prefer string) (*m3u8.Variant, error) {
	if len(master.Variants) == 0 {
		return nil, errors.New("no variants in master playlist")
	}
	// 优先匹配名称/分辨率/带宽包含 prefer 的条目
	if prefer != "" {
		for _, v := range master.Variants {
			label := []string{v.Name, v.Codecs}
			if v.Resolution != "" {
				label = append(label, v.Resolution)
			}
			label = append(label, strconv.FormatInt(int64(v.Bandwidth), 10))
			joined := strings.ToLower(strings.Join(label, ","))
			if strings.Contains(joined, strings.ToLower(prefer)) {
				return v, nil
			}
		}
	}
	// 否则选带宽最高的
	var best *m3u8.Variant
	for _, v := range master.Variants {
		if best == nil || v.Bandwidth > best.Bandwidth {
			best = v
		}
	}
	return best, nil
}

// localSegName 构造每一个m3u8数据分片的文件名
func localSegName(absURI string, seq uint64) string {
	// 统一以 seq+原始后缀 命名，便于本地播放器顺序请求
	u, err := url.Parse(absURI)
	if err != nil {
		return fmt.Sprintf("%d.bin", seq)
	}
	ext := path.Ext(u.Path)
	if ext == "" {
		ext = ".bin"
	}
	return fmt.Sprintf("%d%s", seq, ext)
}

// fetchOnce 拉取并解析一个 m3u8 文本
func (hb *HLSBroadcaster) fetchOnce(ctx context.Context, client *http.Client, u string) (m m3u8.Playlist, body []byte, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("User-Agent", "hls-relay/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("bad status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	p, listType, err := m3u8.Decode(*bytes.NewBuffer(b), true)
	if err != nil {
		return nil, nil, err
	}
	_ = listType
	return p, b, nil
}

// PullWorker 持续从上游拉取分片并写入 stream state

// PullLoop 持续去直播原地址拉流/数据
func (hb *HLSBroadcaster) PullLoop(bo broadcast.BroadcasterOptional) {
	/*
		这是一个后台 goroutine，用来持续从某个 HLS 上游地址拉取数据。
		它先请求 Master Playlist，如果是多码率流，选择合适变体变成 Media Playlist。
		定时轮询 Media Playlist（默认 800ms），发现新分片后下载。
		下载到分片后，调用 stream.PushSegment() 把它放入对应的 StreamState 环形缓存。

		核心点：
			通过 seen 维护已下载分片，避免重复下载。
			计算本地序列号 Seq，保证分片顺序。
			下载的分片保持原样字节，不做解码重封装，性能好且稳定。
	*/

	log.Printf("[pull:%s] start from %s", hb.BrokerKey, hb.upstreamURL)
	client := &http.Client{Timeout: 10 * time.Second}

	stream := hb.StreamState0
	seen := map[string]bool{}
	var mediaURL string

	// 初次处理 master/ media
	p, _, err := hb.fetchOnce(hb.ctx, client, hb.upstreamURL)
	if err != nil {
		log.Printf("[pull:%s] fetch master/media failed: %v", hb.BrokerKey, err)
		return
	}
	if mp, ok := p.(*m3u8.MasterPlaylist); ok {
		v, err := pickVariant(mp, hb.Variant)
		if err != nil {
			log.Printf("[pull:%s] no variant: %v", hb.BrokerKey, err)
			return
		}
		mediaURL, err = resolveURL(hb.upstreamURL, v.URI)
		if err != nil {
			log.Printf("[pull:%s] resolve media url: %v", hb.BrokerKey, err)
			return
		}
		log.Printf("[pull:%s] choose variant bw=%d res=%s uri=%s", hb.BrokerKey, v.Bandwidth, v.Resolution, mediaURL)
	} else if _, ok := p.(*m3u8.MediaPlaylist); ok {
		mediaURL = hb.upstreamURL
	} else {
		log.Printf("[pull:%s] unknown playlist type", hb.BrokerKey)
		return
	}

	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()

	var lastSeq uint64
	for {
		select {
		case <-hb.ctx.Done():
			log.Printf("[pull:%s] stop", hb.BrokerKey)
			return
		case <-ticker.C:
			p, _, err := hb.fetchOnce(hb.ctx, client, mediaURL)
			if err != nil {
				log.Printf("[pull:%s] fetch media: %v", hb.BrokerKey, err)
				continue
			}
			mp, ok := p.(*m3u8.MediaPlaylist)
			if !ok {
				log.Printf("[pull:%s] not media playlist", hb.BrokerKey)
				continue
			}

			// 更新 target duration
			if mp.TargetDuration > 0 {
				stream.Mu.Lock()
				stream.TargetDur = float64(mp.TargetDuration)
				stream.Mu.Unlock()
			}

			// 遍历新片段
			for _, seg := range mp.Segments {
				if seg == nil {
					continue
				}
				absURI, err := resolveURL(mediaURL, seg.URI)
				if err != nil {
					continue
				}
				if seen[absURI] {
					continue
				}

				// 估算 seq：用节目序列号 + 相对偏移（若提供）
				var seq uint64
				if mp.SeqNo != 0 {
					seq = uint64(mp.SeqNo) + uint64(seg.SeqId)
				} else {
					// 回退：自增
					lastSeq++
					seq = lastSeq
				}

				data, err := hb.download(hb.ctx, client, absURI)
				if err != nil {
					log.Printf("[pull:%s] seg dl: %v", hb.BrokerKey, err)
					continue
				}

				localName := localSegName(absURI, seq)
				fmt.Println("分片创建完成：.filename = ", localName)
				stream.PushSegment(&Segment{
					Seq:       seq,
					URI:       absURI,
					LocalName: localName,
					Data:      data,
					Dur:       seg.Duration,
					Discont:   seg.Discontinuity,
					AddedAt:   time.Now(),
				})

				seen[absURI] = true
				lastSeq = seq
			}
		}
	}
}

func (hb *HLSBroadcaster) download(ctx context.Context, client *http.Client, u string) ([]byte, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("User-Agent", "hls-relay/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %d for %s", resp.StatusCode, u)
	}
	return io.ReadAll(resp.Body)
}

// AddLiveClient 添加客户端
func (hb *HLSBroadcaster) AddLiveClient(clientId string, client client.LiveClient) {
	hb.clientMutex.Lock()
	defer hb.clientMutex.Unlock()

	hb.clientMap[clientId] = client
}

// RemoveLiveClient 移除客户端
func (hb *HLSBroadcaster) RemoveLiveClient(clientId string) {
	hb.clientMutex.Lock()
	defer hb.clientMutex.Unlock()
	liveClient := hb.clientMap[clientId]
	if liveClient == nil {
		return
	}
	delete(hb.clientMap, clientId)

	//// 如果没有客户端并且想释放 broker，可关闭 stopCh 让 PullLoop 停止（本示例保留 broker，防止频繁断开上游）
	//if remaining == 0 {
	//	// optionally stop pulling after idle timeout. For simplicity we keep running.
	//}
}

// FindLiveClient 查询 LiveClient
func (hb *HLSBroadcaster) FindLiveClient(clientId string) (client.LiveClient, error) {
	if val, ok := hb.clientMap[clientId]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("未找到 %s 对应的 LiveClient", clientId))

}

// UpdateSourceURL 支持切换直播原地址
func (hb *HLSBroadcaster) UpdateSourceURL(newSourceURL string) {

}

// ListenStatus 监听当前直播的必要状态
func (hb *HLSBroadcaster) ListenStatus() {
	for {
		select {
		case clientId := <-hb.ClientCloseSig:
			// 监听客户端离开消息
			hb.RemoveLiveClient(clientId)
			fmt.Printf("HLSBroadcaster.ListenStatus.RemoveLiveClient.clientId %s successful.", clientId)
		case <-hb.BrokerCloseSig:
			// 直播被关闭
			close(hb.BrokerCloseSig)
		}

	}
}

// Broadcast2LiveClient 原地址拉取到数据之后广播给客户端
func (hb *HLSBroadcaster) Broadcast2LiveClient(data []byte) {

}
