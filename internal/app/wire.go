package app

import (
	"github.com/zhilv666/linkchecker/configs"
	"github.com/zhilv666/linkchecker/internal/handler"
	"github.com/zhilv666/linkchecker/internal/model"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/netdisk/baidu"
	"github.com/zhilv666/linkchecker/internal/netdisk/quark"
	"github.com/zhilv666/linkchecker/internal/repo"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/db"
	"github.com/zhilv666/linkchecker/pkg/request"
	"go.uber.org/zap"
)

type AppContainer struct {
	LinkHandler *handler.LinkHandler
	LinkService *service.LinkService
	Cleanup     func()
}

func NewAppContainer(cfg *configs.Config, logger *zap.Logger) (*AppContainer, error) {
	dbCfg := &db.Config{
		Type:          cfg.Database.Type,
		DSN:           cfg.Database.GetDSN(),
		MaxIdleConns:  cfg.Database.MaxIdleConns,
		MaxOpenConns:  cfg.Database.MaxOpenConns,
		MaxLifetime:   cfg.Database.MaxLifetime,
		TablePrefix:   cfg.Database.TablePrefix,
		SingularTable: cfg.Database.SingularTable,
		Debug:         cfg.Database.Debug,
	}

	gormDB, cleanup, err := db.New(dbCfg, logger)
	if err != nil {
		return nil, err
	}
	err = gormDB.AutoMigrate(
		new(model.LinkRecord),
	)
	if err != nil {
		return nil, err
	}

	memCache := cache.New(&cache.Config{})
	httpClient := request.NewClient(nil)

	ndManager := netdisk.NewManager(memCache,
		baidu.New(httpClient),
		quark.New(httpClient),
	)

	linkRepo := repo.NewLinkRepo(gormDB)
	linkSvc := service.NewLinkService(linkRepo, ndManager)
	linkH := handler.NewLinkHandler(linkSvc)

	return &AppContainer{
		LinkHandler: linkH,
		LinkService: linkSvc,
		Cleanup: func() {
			cleanup()
		},
	}, nil
}
