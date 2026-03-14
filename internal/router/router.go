package router

import (
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/app"
	"github.com/zhilv666/linkchecker/internal/middleware"
	"github.com/zhilv666/linkchecker/pkg/log"
	"github.com/zhilv666/linkchecker/web"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func SetupRouter(cfg *configs.Config, logger *zap.Logger, app *app.AppContainer) *gin.Engine {
	router := gin.Default()

	limiter := middleware.NewIPRateLimiter(rate.Every(1*time.Second), 3)
	tokener := middleware.TokenMiddleware(cfg.Server.Token)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowHeaders:     cfg.Cors.AllowHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	dist, err := fs.Sub(web.Public, "dist")
	if err != nil {
		log.Fatal("failed to read dist dir", zap.Error(err))
	}
	distHTTP := http.FS(dist)
	router.StaticFS("/js", http.FS(mustSub(dist, "js")))
	router.StaticFS("/assets", http.FS(mustSub(dist, "assets")))
	router.StaticFileFS("/favicon.ico", "favicon.ico", distHTTP)

	router.GET("/", func(ctx *gin.Context) {
		http.ServeFileFS(ctx.Writer, ctx.Request, dist, "index.html")
	})

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	apiV1 := router.Group("/api/v1")
	linkGroup := apiV1.Group("/link")
	linkGroup.Use(middleware.RateMiddleware(limiter))
	{
		linkGroup.POST("/", app.LinkHandler.CheckOne)
		linkGroup.POST("/list", app.LinkHandler.ListWithPageSize)
	}
	apiV1.POST("/report", tokener, app.LinkHandler.Report)

	router.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.JSON(404, gin.H{"code": 404, "msg": "API not found"})
			return
		}
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		http.ServeFileFS(ctx.Writer, ctx.Request, dist, "index.html")
	})

	return router
}

// 辅助函数：简化 fs.Sub 的错误处理
func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic("failed to find static dir: " + dir)
	}
	return sub
}
