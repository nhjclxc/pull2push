package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// Config 总配置结构
type Config struct {
	Common CommonConfig `yaml:"common"`
	DB     DBConfig     `yaml:"db"`
	HTTP   HTTPConfig   `yaml:"http"`
	Live   LiveConfig   `yaml:"live"`
}

// CommonConfig 包含共享的配置项
type CommonConfig struct {
	Log   LogConfig `yaml:"log"`
	Debug bool      `yaml:"debug"`
}

type LogConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大尺寸（MB）
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件数量
}

// HTTPConfig HTTP服务特定配置
type HTTPConfig struct {
	Port      string `yaml:"port"`
	ProxyHost string `yaml:"proxy_host"` // 代理目标主机地址
}

// DBConfig HTTP服务特定配置
type DBConfig struct {
	//MysqlConf DBMysqlConfig `yaml:"mysql"`
	//RedisConf DBRedisConfig `yaml:"redis"`
}

// LiveConfig 直播配置
type LiveConfig struct {
	hlsPort    int `yaml:"hlsPort"`  // 80/443
	flvPort    int `yaml:"flvPort"`  // 8080/80/443
	rtmpPort   int `yaml:"rtmpPort"` // 1935
	cameraPort int `yaml:"cameraPort"`
}
