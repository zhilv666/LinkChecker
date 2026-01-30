package router

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/zhilv666/linkchecker/configs"
	dbMod "github.com/zhilv666/linkchecker/internal/db"
	"github.com/zhilv666/linkchecker/internal/handler"
	localMiddleware "github.com/zhilv666/linkchecker/internal/middleware"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/web"
)

func SetupRouter(cfg *configs.Config) *chi.Mux {
	r := chi.NewRouter()
	cache := cache.New(&cache.Config{})
	db := dbMod.GetDB()
	linkRepo := dbMod.NewLinkDB(db)
	linkService := service.NewLinkService(linkRepo, cache)
	linkHandler := handler.NewLinkHandler(linkService)

	// webHandler := handler.NewWebHandler(linkService)

	// 1. 注册基础中间件
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Cors.AllowOrigins,
		AllowedMethods:   cfg.Cors.AllowMethods,
		AllowedHeaders:   cfg.Cors.AllowHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// workDir, _ := os.Getwd()
	// filesDir := http.Dir(filepath.Join(workDir, "web/static"))
	distFS, err := fs.Sub(web.Public, "dist")
	if err != nil {
		panic(err)
	}

	// 2. 注册自定义中间件
	r.Use(localMiddleware.LoggerMiddleware())

	// FileServer(r, "/static", filesDir)
	r.Handle("/*", SpaHandler(distFS))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.Write([]byte("pong"))
	})

	// r.Get("/", webHandler.Index)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/link", func(r chi.Router) {
				r.Get("/", linkHandler.CheckOne)
				r.Get("/list", linkHandler.ListWithPageSize)
			})
		})
	})

	return r
}

func SpaHandler(staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")

		// 1. 检查文件是否存在
		f, err := staticFS.Open(path)
		if err != nil {
			// ❌ 文件不存在 -> 返回 index.html

			// 【修正点】使用 fs.ReadFile 读取文件内容为 []byte
			indexBytes, err := fs.ReadFile(staticFS, "index.html")
			if err != nil {
				http.Error(w, "index.html missing", 500)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			// 现在传入的是 []byte，这就对了
			w.Write(indexBytes)
			return
		}

		// ✅ 文件存在
		f.Close() // 只是为了检查是否存在，查完就关掉
		http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
	}
}

// // FileServer 是一个辅助函数，用于在 chi 中方便地提供静态文件服务
// func FileServer(r chi.Router, path string, root http.FileSystem) {
// 	if strings.ContainsAny(path, "{}*") {
// 		panic("FileServer does not permit any URL parameters.")
// 	}

// 	if path != "/" && path[len(path)-1] != '/' {
// 		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
// 		path += "/"
// 	}
// 	path += "*"

// 	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
// 		rctx := chi.RouteContext(r.Context())
// 		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
// 		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
// 		fs.ServeHTTP(w, r)
// 	})
// }
