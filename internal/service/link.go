package service

import (
	"github.com/zhilv666/linkchecker/internal/db"
	"github.com/zhilv666/linkchecker/internal/model"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/netdisk/baidu"
	"github.com/zhilv666/linkchecker/internal/netdisk/quark"
	"github.com/zhilv666/linkchecker/pkg/cache"
	"github.com/zhilv666/linkchecker/pkg/log"
	"github.com/zhilv666/linkchecker/pkg/request"
	"go.uber.org/zap"
)

type LinkService struct {
	linkDB *db.LinkDB
	cache  cache.Cache
}

func NewSubManager(cache cache.Cache) *netdisk.Manager {
	client := request.NewRestyClient()
	manager := netdisk.NewManager(cache,
		baidu.New(client),
		quark.New(client),
	)
	return manager
}

func NewLinkService(l *db.LinkDB, c cache.Cache) *LinkService {
	return &LinkService{
		linkDB: l,
		cache:  c,
	}
}

func (ls *LinkService) CheckOne(url, password string) (*model.LinkRecord, error) {
	manager := NewSubManager(ls.cache)
	info, err := manager.Check(url, password)

	if err != nil {
		log.Error("检测失败", zap.Error(err))
		return nil, err
	}

	lr := &model.LinkRecord{
		Provider:  info.Provider,
		Title:     info.Title,
		Size:      info.Size,
		Author:    info.Author,
		Status:    info.Status,
		ExpiredAt: info.ExpiredAt,
		URL:       info.NormalizedUrl,
		PWD:       info.Password,
	}
	if info.NormalizedUrl == "" {
		lr.URL = info.RawUrl
	}

	err = ls.linkDB.Create(lr)

	if err != nil {
		log.Error("添加数据失败", zap.Error(err))
		return nil, err
	}

	return lr, err
}

func (ls *LinkService) ListWithPageSize(page, size int, keyword string) ([]model.LinkRecord, int64, error) {
	return ls.linkDB.List(page, size, keyword)
}
