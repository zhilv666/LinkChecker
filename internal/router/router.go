package router

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/handler"
	"github.com/zhilv666/linkchecker/internal/middleware"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/netdisk/baidu"
	"github.com/zhilv666/linkchecker/internal/netdisk/quark"
	"github.com/zhilv666/linkchecker/internal/repo"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/log"
	"github.com/zhilv666/linkchecker/pkg/request"
	"github.com/zhilv666/linkchecker/web"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func SetupRouter(cfg *configs.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	cache := cache.New(&cache.Config{})
	client := request.NewRestyClient()
	manager := netdisk.NewManager(cache, baidu.New(client), quark.New(client))
	linkRepo := repo.NewLinkRepo(db)
	linkService := service.NewLinkService(linkRepo, manager)
	linkHandler := handler.NewLinkHandler(linkService)

	limiter := middleware.NewIPRateLimiter(rate.Every(1*time.Second), 3)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowHeaders:     cfg.Cors.AllowHeaders,
		ExposeHeaders:    []string{"Context-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	dist, err := fs.Sub(web.Public, "dist")
	if err != nil {
		log.Fatal("failed to read dist dir", zap.Error(err))
	}
	distHTTP := http.FS(dist)
	subJs, err := fs.Sub(dist, "js")
	if err != nil {
		log.Fatal("failed to read js dir in dist", zap.Error(err))
	}
	router.StaticFS("/js", http.FS(subJs))

	subAssets, err := fs.Sub(dist, "assets")
	if err != nil {
		log.Fatal("failed to read assets dir in dist", zap.Error(err))
	}
	router.StaticFS("/assets", http.FS(subAssets))

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
		linkGroup.POST("/", linkHandler.CheckOne)
		linkGroup.GET("/list", linkHandler.ListWithPageSize)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(404, "页面不存在")
	})

	return router
}
