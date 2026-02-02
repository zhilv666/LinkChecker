package router

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/configs"
	dbMod "github.com/zhilv666/linkchecker/internal/db"
	"github.com/zhilv666/linkchecker/internal/handler"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/log"
	"github.com/zhilv666/linkchecker/web"
	"go.uber.org/zap"
)

func SetupRouter(cfg *configs.Config) *gin.Engine {
	router := gin.Default()
	cache := cache.New(&cache.Config{})
	db := dbMod.GetDB()
	linkRepo := dbMod.NewLinkDB(db)
	linkService := service.NewLinkService(linkRepo, cache)
	linkHandler := handler.NewLinkHandler(linkService)

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

	router.GET("/", func(ctx *gin.Context) {
		http.ServeFileFS(ctx.Writer, ctx.Request, dist, "index.html")
	})

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	apiV1 := router.Group("/api/v1")
	linkGroup := apiV1.Group("/link")
	{
		linkGroup.GET("/", linkHandler.CheckOne)
		linkGroup.GET("/list", linkHandler.ListWithPageSize)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.FileFromFS("index.html", distHTTP)
	})

	return router
}
