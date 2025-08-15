package cmd

import (
	"context"
	"os"
	"os/signal"
	"pull2push/application"
	"pull2push/config"
	"pull2push/cron"
	"pull2push/event"
	"pull2push/logger"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "pull2push",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server that handles live streaming and proxy requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			logger.Error("Failed to load config", "error", err)
		}
		logger.Info("Config loaded successfully")

		// 1、创建服务管理器
		serviceManager, err := application.NewServiceManager(cfg)
		if err != nil {
			logger.Error("Failed to create service manager", "error", err)
			return err
		}
		logger.Info("Service manager created")

		// 2. HTTP 服务
		httpService := application.NewHTTPService(serviceManager.GetResource())
		if err := serviceManager.AddService(httpService); err != nil {
			logger.Error("Failed to add HTTP service", "error", err)
			return err
		}

		// 3. 任务管理器服务 - 没有特定依赖
		cronTaskManager := cron.NewCronTaskManager(serviceManager.GetResource())

		// 3.3 注册任务管理器服务
		if err := serviceManager.AddService(cronTaskManager); err != nil {
			logger.Error("Failed to add CronTaskManager service", "error", err)
			return err
		}

		// 启动所有服务
		if err := serviceManager.StartAll(ctx); err != nil {
			logger.Error("Failed to start services", "error", err)
			return err
		}
		logger.Info("All services started successfully")

		// 等待中断信号
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// 使用单独的退出通道
		exitCh := make(chan struct{})

		// 监听信号
		go func() {
			sig := <-sigCh
			logger.Info("Received termination signal", "signal", sig)

			// 发布系统关闭事件，通知所有模块
			serviceManager.GetEventBus().Publish(event.Event{
				Type:    event.SystemShutdown,
				Payload: "Shutdown triggered by signal",
			})

			// 取消主context，通知所有服务开始关闭
			cancel()

			// 优雅关闭所有服务
			logger.Info("=========Starting graceful shutdown...=========")

			// 创建退出超时计时器
			shutdownTimer := time.NewTimer(10 * time.Second)

			// 创建一个完成通道
			doneCh := make(chan struct{})

			// 在后台停止所有服务
			go func() {
				serviceManager.StopAll()
				logger.Info("All services stopped successfully")
				close(doneCh)
			}()

			// 等待服务停止或超时
			select {
			case <-doneCh:
				logger.Info("Graceful shutdown completed")
				if !shutdownTimer.Stop() {
					<-shutdownTimer.C
				}
			case <-shutdownTimer.C:
				logger.Error("Shutdown timed out, forcing exit")
			}

			// 通知主程序可以退出
			close(exitCh)
		}()

		// 等待退出信号或上下文取消
		select {
		case <-exitCh:
			logger.Info("Server shutdown completed")
		case <-ctx.Done():
			logger.Info("Context canceled, shutting down")
		}
		// 确保程序退出
		os.Exit(0)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file path")
	serveCmd.MarkFlagRequired("config")
}
