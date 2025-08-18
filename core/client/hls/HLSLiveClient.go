package flv

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pull2push/core/broadcast"
	hlsBroadcast "pull2push/core/broadcast/hls"
	"strings"
)

// ====================== HLSLiveClient ======================

// HLSLiveClient 每一个前端页面有持有一个客户端对象
type HLSLiveClient struct {
	BroadcasterKey string      // 这个客户端的直播房间的唯一编号
	ClientId       string      // 这个客户端的id
	DataCh         chan []byte // 这个客户端的一个只写通道

	// http连接相关
	httpCloseSig        <-chan struct{} // 当这个请求被客户端主动被关闭时触发
	httpRequestCloseSig <-chan struct{} // 当这个请求被客户端主动被关闭时触发

	// 父级 broadcaster 相关的内容
	clientCloseSig      chan<- string                         // broker通过该信道监听客户端离线 【仅发送】
	broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE // broker被关闭时，同时通知客户端关闭 【仅接收】
}

func NewHLSLiveClient(c *gin.Context, broadcasterKey, clientId string, clientCloseSig chan<- string, broadcasterCloseSig <-chan broadcast.BROADCAST_CLOSE_TYPE) (*HLSLiveClient, error) {

	hlc := HLSLiveClient{
		BroadcasterKey:      broadcasterKey,
		ClientId:            clientId,
		httpCloseSig:        c.Done(),
		httpRequestCloseSig: c.Request.Context().Done(),
		clientCloseSig:      clientCloseSig,
		broadcasterCloseSig: broadcasterCloseSig,
	}

	fmt.Println("HLS 客户端连接成功 ClientId = ", clientId)

	// 开启状态监听
	go hlc.Listen()

	return &hlc, nil
}

func (hlc *HLSLiveClient) Listen() {

	for {
		select {
		case <-hlc.httpCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("hlc.httpCloseSig 收到客户端关闭信号，退出循环 ", hlc.ClientId)
			//
			//// when client closes, remove it
			//hlc.clientCloseSig <- hlc.ClientId

			return
		case <-hlc.httpRequestCloseSig:
			//// 收到关闭信号，退出循环
			//fmt.Println("<-hlc.httpRequestCloseSig 收到客户端关闭信号，退出循环 ", hlc.ClientId)
			//
			//// when client closes, remove it
			//hlc.clientCloseSig <- hlc.ClientId

			return
		case <-hlc.broadcasterCloseSig:
			return
		}
	}

}

func (hlc *HLSLiveClient) Broadcast(data []byte) {

}

// GetDataChan 获取当前客户端的写通道
func (hlc *HLSLiveClient) GetDataChan() chan []byte {
	return hlc.DataCh
}

func (hlc *HLSLiveClient) HandleIndex(w http.ResponseWriter, r *http.Request, findBroadcasterTemp *hlsBroadcast.HLSBroadcaster) {
	// /live/hls/{broadcasterKey}/{clientID}/index.m3u8
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if parts[1] != "live" || parts[len(parts)-1] != "index.m3u8" {
		http.NotFound(w, r)
		return
	}
	if findBroadcasterTemp.StreamState0 == nil {
		http.NotFound(w, r)
		return
	}

	segs, seqStart, targetDur, discont := findBroadcasterTemp.StreamState0.Snapshot()
	pl, err := hlc.buildMediaPlaylist(segs, seqStart, targetDur, discont, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(pl))
}

func (hlc *HLSLiveClient) HandleSegment(w http.ResponseWriter, r *http.Request, findBroadcasterTemp *hlsBroadcast.HLSBroadcaster) {
	// /api/live/hls/{broadcasterKey}/{clientID}/{seg.ts|m4s}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	filename := parts[len(parts)-1]
	if parts[1] != "live" || !(strings.HasSuffix(filename, ".ts") || strings.HasSuffix(filename, ".m4s")) {
		http.NotFound(w, r)
		return
	}

	if findBroadcasterTemp.StreamState0 == nil {
		http.NotFound(w, r)
		return
	}

	findBroadcasterTemp.StreamState0.Mu.RLock()
	defer findBroadcasterTemp.StreamState0.Mu.RUnlock()

	var seg *hlsBroadcast.Segment
	findBroadcasterTemp.StreamState0.Segments.Do(func(v any) {
		if v == nil {
			return
		}
		ss := v.(*hlsBroadcast.Segment)
		if ss != nil && ss.LocalName == filename {
			seg = ss
		}
	})
	if seg == nil {
		http.NotFound(w, r)
		return
	}

	// 内容类型根据后缀猜测
	if strings.HasSuffix(filename, ".ts") {
		w.Header().Set("Content-Type", "video/mp2t")
	} else if strings.HasSuffix(filename, ".m4s") || strings.HasSuffix(filename, ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	_, _ = w.Write(seg.Data)
}

// buildMediaPlaylist HTTP 播放列表生成与分片访问
// 生成 HLS 标准的分片信息，顺序排列，标明时长和断点。
// 返回给播放器标准 HLS 播放列表。
func (hlc *HLSLiveClient) buildMediaPlaylist(segs []*hlsBroadcast.Segment, seqStart uint64, targetDur float64, discont bool, r *http.Request) (string, error) {

	// handleIndex 负责根据当前 StreamState 缓存的分片快照，生成标准 HLS 播放列表文本。
	// 播放列表里指向本地缓存的分片文件名（seq.ts 或 seq.m4s）。
	// handleSegment 负责根据请求的分片名返回对应的分片字节流。

	if len(segs) == 0 {
		// 空列表也要有基本头信息，避免播放器报错
		return "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:6\n#EXT-X-MEDIA-SEQUENCE:0\n", nil
	}
	var b strings.Builder
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
	b.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", int(targetDur+0.5)))
	b.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", seqStart))
	// 可选：I-Frame only、MAP 等根据上游情况补充
	if discont {
		b.WriteString("#EXT-X-DISCONTINUITY-SEQUENCE:1\n")
	}

	base := fmt.Sprintf("/api/live/hls/" + hlc.BroadcasterKey + "/" + hlc.ClientId + "/")
	for _, s := range segs {
		if s == nil {
			continue
		}
		if s.Discont {
			b.WriteString("#EXT-X-DISCONTINUITY\n")
		}
		b.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n", s.Dur))
		b.WriteString(base + s.LocalName + "\n")
	}
	return b.String(), nil
}
