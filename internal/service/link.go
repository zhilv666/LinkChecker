package service

import (
	"context"

	"github.com/zhilv666/linkchecker/internal/model"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/repo"
)

type ILinkService interface {
	CheckAndSave(ctx context.Context, url, password string) (*model.LinkRecord, error)
	GetLinkList(ctx context.Context, page, size int, keyword string) ([]*model.LinkRecord, int64, error)
}

type linkService struct {
	repo    repo.ILinkRepo
	manager *netdisk.Manager
}

func NewLinkService(repo repo.ILinkRepo, manager *netdisk.Manager) *linkService {
	return &linkService{repo: repo, manager: manager}
}

func (s *linkService) CheckAndSave(ctx context.Context, url, password string) (*model.LinkRecord, error) {
	exist, _ := s.repo.GetByRawUrl(ctx, url)
	if exist != nil {
		return exist, nil
	}

	info, err := s.manager.Check(url, password)
	if err != nil {
		return nil, err
	}

	newLink := &model.LinkRecord{
		Provider:  info.Provider,
		Title:     info.Title,
		Size:      info.Size,
		Author:    info.Author,
		Status:    info.Status,
		ExpiredAt: info.ExpiredAt,
		URL:       info.NormalizedUrl,
		RawURL:    info.RawUrl,
		PWD:       info.Password,
	}

	if err := s.repo.Create(ctx, newLink); err != nil {
		return nil, err
	}
	return newLink, err
}

func (s *linkService) GetLinkList(ctx context.Context, page, size int, keyword string) ([]*model.LinkRecord, int64, error) {
	return s.repo.GetList(ctx, page, size, keyword)
}
