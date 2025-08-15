package repository

import (
	"fmt"
	"sync"
	"time"
)

// LiveInfo 直播信息
type LiveInfo struct {
	ID        string    // 直播ID
	Name      string    // 直播名称
	SourceURL string    // 源地址
	BeginAt   time.Time // 开始时间
	EndAt     time.Time // 结束时间
	Status    string    // 状态：created, active, ended, error
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间

	// 视频元数据
	Width     float64 // 视频宽度
	Height    float64 // 视频高度
	FrameRate float64 // 帧率
}

// LiveRepository 直播信息仓库
type LiveRepository struct {
	mu    sync.RWMutex
	lives map[string]*LiveInfo // key: LiveID
}

// NewLiveRepository 创建新的直播信息仓库
func NewLiveRepository() *LiveRepository {
	return &LiveRepository{
		lives: make(map[string]*LiveInfo),
	}
}

// Save 保存直播信息
func (r *LiveRepository) Save(info *LiveInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if info.ID == "" {
		return fmt.Errorf("live ID cannot be empty")
	}

	if _, exists := r.lives[info.ID]; exists {
		return fmt.Errorf("live with ID %s already exists", info.ID)
	}

	now := time.Now()
	if info.CreatedAt.IsZero() {
		info.CreatedAt = now
	}
	info.UpdatedAt = now

	// 设置视频元数据的默认值
	if info.Width == 0 {
		info.Width = 1280
	}
	if info.Height == 0 {
		info.Height = 720
	}
	if info.FrameRate == 0 {
		info.FrameRate = 30
	}

	r.lives[info.ID] = info
	return nil
}

// Get 获取直播信息
func (r *LiveRepository) Get(id string) (*LiveInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.lives[id]
	return info, exists
}

// List 列出所有直播信息
func (r *LiveRepository) List() []*LiveInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lives := make([]*LiveInfo, 0, len(r.lives))
	for _, info := range r.lives {
		lives = append(lives, info)
	}
	return lives
}

// Delete 删除直播信息
func (r *LiveRepository) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.lives, id)
}

// UpdateStatus 更新直播状态
func (r *LiveRepository) UpdateStatus(id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, exists := r.lives[id]
	if !exists {
		return fmt.Errorf("live with ID %s not found", id)
	}

	info.Status = status
	info.UpdatedAt = time.Now()
	return nil
}

// GetByStatus 根据状态获取直播信息
func (r *LiveRepository) GetByStatus(status string) []*LiveInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*LiveInfo
	for _, info := range r.lives {
		if info.Status == status {
			result = append(result, info)
		}
	}
	return result
}

// GetActive 获取当前活跃的直播
func (r *LiveRepository) GetActive() []*LiveInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	var result []*LiveInfo
	for _, info := range r.lives {
		if info.BeginAt.Before(now) && info.EndAt.After(now) && info.Status == "active" {
			result = append(result, info)
		}
	}
	return result
}

// CleanExpired 清理过期的直播
func (r *LiveRepository) CleanExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, info := range r.lives {
		if info.EndAt.Before(now) {
			delete(r.lives, id)
		}
	}
}
