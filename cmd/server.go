package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/db"
	"github.com/zhilv666/linkchecker/internal/router"
	"github.com/zhilv666/linkchecker/pkg/log"
	"go.uber.org/zap"
)

var (
	port string // 接收命令行参数
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 Web 服务",
	Long:  `启动 LinkChecker 的后端 API 服务，默认监听 3000 端口`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 初始化配置
		cfg := configs.InitConfig()

		// 2. 初始化日志 (必须紧跟配置之后，且在打印任何日志之前)
		log.Init(log.Config{
			Level:     cfg.Log.Level,
			Filepath:  cfg.Log.Filepath,
			MaxSizeMB: cfg.Log.MaxSizeMB,
			MaxAgeDay: cfg.Log.MaxAgeDay,
			Backups:   cfg.Log.Backups,
			Compress:  cfg.Log.Compress,
		})

		// 3. 打印调试信息 (此时日志系统已就绪，可以正确写入文件)
		log.Debug("Server configuration loaded", zap.Any("config", cfg))

		// 4. 初始化数据库
		if err := db.Init(cfg.Database); err != nil {
			// 使用 Fatal 记录错误并退出，比 panic 更优雅，且会有 structured log
			log.Fatal("Failed to initialize database", zap.Error(err))
		}

		// 5. 初始化路由
		r := router.SetupRouter(cfg)

		// 6. 定义 HTTP Server
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}

		// 7. 启动服务 (在 Goroutine 中启动)
		go func() {
			log.Info("Server is running", zap.String("addr", srv.Addr))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("Server start failed", zap.Error(err))
			}
		}()

		// 8. 优雅停机 (Graceful Shutdown) 逻辑
		quit := make(chan os.Signal, 1)
		// 监听中断信号 (Ctrl+C) 和 终止信号 (kill)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞直到接收到信号
		<-quit
		log.Info("Shutting down server...")

		// 创建一个 5 秒超时的 Context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 尝试优雅关闭
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown", zap.Error(err))
		}

		log.Info("Server exited properly")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to listen on")
}
