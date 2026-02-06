package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/app"
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
		cfg := configs.InitConfig()

		log.Init(log.Config{
			Level:     cfg.Log.Level,
			Filepath:  cfg.Log.Filepath,
			MaxSizeMB: cfg.Log.MaxSizeMB,
			MaxAgeDay: cfg.Log.MaxAgeDay,
			Backups:   cfg.Log.Backups,
			Compress:  cfg.Log.Compress,
		})
		logger := log.GetLogger()

		container, err := app.NewAppContainer(cfg, logger)
		if err != nil {
			logger.Fatal("App init failed", zap.Error(err))
		}
		defer container.Cleanup()

		log.Debug("Server configuration loaded", zap.Any("config", cfg))

		if cfg.Server.Debug {
			log.Debug("调试模式")
			gin.SetMode(gin.DebugMode)
		} else {
			log.Debug("生产模式")
			gin.SetMode(gin.ReleaseMode)
		}

		r := router.SetupRouter(cfg, logger, container)
		if cfg.Server.Port != 0 {
			port = fmt.Sprint(cfg.Server.Port)
		}

		srv := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}

		go func() {
			log.Info("Server is running", zap.String("addr", srv.Addr))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("Server start failed", zap.Error(err))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
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
