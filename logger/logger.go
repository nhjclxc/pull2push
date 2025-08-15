package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"pull2push/config"
	"sync"
	"time"
)

// go1.21以后内置的slog
var logger *slog.Logger

func init() {
	// 默认初始化为标准输出，使用文本格式
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// RollingFileWriter 支持每天一个文件，且超过大小时切分
type RollingFileWriter struct {
	dir       string
	baseName  string
	maxSize   int64 // bytes
	curDate   string
	curFile   *os.File
	curSize   int64
	fileIndex int
	mu        sync.Mutex
}

func NewRollingFileWriter(dir string, baseName string, maxSizeMB int64) (*RollingFileWriter, error) {
	w := &RollingFileWriter{
		dir:      dir,
		baseName: baseName,
		maxSize:  maxSizeMB * 1024 * 1024,
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	if err := w.rotate(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RollingFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	// 日期变了，重新建文件
	if today != w.curDate {
		w.fileIndex = 0
		w.curDate = today
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	// 超过大小，切分
	if w.curSize+int64(len(p)) > w.maxSize {
		w.fileIndex++
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.curFile.Write(p)
	w.curSize += int64(n)
	return
}

func (w *RollingFileWriter) rotate() error {
	if w.curFile != nil {
		_ = w.curFile.Close()
	}
	var fileName string
	if w.fileIndex == 0 {
		fileName = fmt.Sprintf("%s.log", w.curDate)
	} else {
		fileName = fmt.Sprintf("%s.%d.log", w.curDate, w.fileIndex)
	}
	path := filepath.Join(w.dir, fileName)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	w.curFile = file
	w.curSize = 0
	if info, _ := file.Stat(); info != nil {
		w.curSize = info.Size()
	}
	return nil
}

func init() {
	// 默认控制台输出
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// InitLogger 可以重新配置logger
func InitLogger(cfg *config.LogConfig) error {
	var writer io.Writer
	if cfg.FilePath == "stdout" {
		writer = os.Stdout
	} else {
		dir := filepath.Dir(cfg.FilePath)
		baseName := filepath.Base(cfg.FilePath)
		rollingWriter, err := NewRollingFileWriter(dir, baseName, 100) // 默认 100MB
		if err != nil {
			return fmt.Errorf("create rolling log writer failed: %w", err)
		}
		writer = io.MultiWriter(os.Stdout, rollingWriter)
	}

	// 设置日志级别
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(writer, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// 便捷日志函数
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// WithContext 添加上下文信息
func WithContext(ctx ...any) *slog.Logger {
	return logger.With(ctx...)
}
