package resource

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pull2push/config"
	"pull2push/repository"
)

// Resource 管理全局资源
type Resource struct {
	DB             *gorm.DB
	Redis          *redis.Client
	DeviceID       string
	Config         *config.Config
	LiveRepository *repository.LiveRepository
	ShutdownSignal chan struct{}
}

// NewResource 创建并初始化资源
func NewResource(config *config.Config) (*Resource, error) {
	r := &Resource{
		Config:         config,
		ShutdownSignal: make(chan struct{}),
	}

	//// 初始化数据库
	//if config.Common.DBConn != "" {
	//	db, err := gorm.Open(mysql.Open(config.Common.DBConn), &gorm.Config{})
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to connect to database: %w", err)
	//	}
	//	r.DB = db
	//}
	//
	//// 初始化 Redis
	//if config.Common.RedisAddr != "" {
	//	r.Redis = redis.NewClient(&redis.Options{
	//		Addr: config.Common.RedisAddr,
	//	})
	//	if err := r.Redis.Ping(context.Background()).Err(); err != nil {
	//		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	//	}
	//}

	// 初始化 LiveRepository
	r.LiveRepository = repository.NewLiveRepository()

	return r, nil
}

// Close 关闭所有资源
func (r *Resource) Close() error {
	if r.Redis != nil {
		if err := r.Redis.Close(); err != nil {
			return fmt.Errorf("failed to close redis: %w", err)
		}
	}

	if r.DB != nil {
		sqlDB, err := r.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	// 发送关闭信号
	close(r.ShutdownSignal)

	return nil
}
