package hls

import (
	"container/ring"
	"sync"
	"time"
)

// Segment 每一个m3u8数据分片的数据对象，代表 HLS 的一个 TS 或 fMP4 分片。记录了下载地址和本地暴露的名字，数据内容和时长。
type Segment struct {
	Seq       uint64    // 分片序列号（递增）
	URI       string    // 上游绝对地址（下载用）
	LocalName string    // 本地暴露的文件名（如 seq.ts 或 seq.m4s）
	Data      []byte    // 分片字节
	Dur       float64   // 分片时长，秒
	Discont   bool      // 是否断点分片
	AddedAt   time.Time // 拉取时间
}

// StreamState 每一路拉流任务维护一个 StreamState，存放它的分片缓存和元数据。
// 用 ring.Ring 实现固定容量的循环队列，保持缓存窗口。
type StreamState struct {
	Mu        sync.RWMutex
	Segments  *ring.Ring // 环形缓冲，存放最近 N 个分片，元素为 *Segment 或 nil
	Cap       int        // 缓冲分片数
	TargetDur float64    // HLS 目标分片时长
	SeqStart  uint64     // 本地播放列表起始序列号
	LastSeq   uint64     // 最新分片序列号（递增）
	LastMod   time.Time  // 最后更新时间
	Discont   bool       // 是否有断点续播
}

// NewStreamState 创建每一个直播的拉流缓冲区对象
func NewStreamState(cap int) *StreamState {
	return &StreamState{
		Segments:  ring.New(cap),
		Cap:       cap,
		TargetDur: 6,
		SeqStart:  0,
		LastSeq:   0,
	}
}

// PushSegment 在环形缓冲里追加一个分片
func (s *StreamState) PushSegment(seg *Segment) {
	/*
		维护固定容量缓存，最新的分片覆盖最旧的。
		更新本地播放列表序列号区间。
		保护并发安全（互斥锁）。
	*/

	s.Mu.Lock()
	defer s.Mu.Unlock()

	// 移动指针到下一格并覆盖
	s.Segments = s.Segments.Next()
	s.Segments.Value = seg

	if s.SeqStart == 0 {
		// 第一次写入
		s.SeqStart = seg.Seq
	}
	// 计算当前窗口的“起始序列号” = 最新序列 - (cap-1)
	if seg.Seq+1 >= uint64(s.Cap) {
		s.SeqStart = seg.Seq + 1 - uint64(s.Cap)
	}
	s.LastSeq = seg.Seq
	s.LastMod = time.Now()
	if seg.Discont {
		s.Discont = true
	}
}

// Snapshot 返回按序的窗口分片拷贝（只读）
func (s *StreamState) Snapshot() (segs []*Segment, seqStart uint64, targetDur float64, discont bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	segs = make([]*Segment, 0, s.Cap)
	// 从环形缓冲按时间顺序读出
	tmp := s.Segments
	tmp.Do(func(v any) {
		if v == nil {
			return
		}
		seg := v.(*Segment)
		// 只保留窗口内的、非空数据
		if seg != nil && len(seg.Data) > 0 {
			segs = append(segs, seg)
		}
	})
	return segs, s.SeqStart, s.TargetDur, s.Discont
}
