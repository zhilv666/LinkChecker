package service

import (
	"context"
	"time"

	"github.com/zhilv666/linkchecker/internal/dto"
	"github.com/zhilv666/linkchecker/internal/model"
	"github.com/zhilv666/linkchecker/internal/netdisk"
	"github.com/zhilv666/linkchecker/internal/repo"
	"gorm.io/gorm"
)

type LinkService struct {
	repo    repo.ILinkRepo
	manager *netdisk.Manager
}

func NewLinkService(repo repo.ILinkRepo, manager *netdisk.Manager) *LinkService {
	return &LinkService{repo: repo, manager: manager}
}

func (s *LinkService) CheckAndSave(ctx context.Context, rawUrl, password string) (*model.LinkRecord, error) {
	info, err := s.manager.Check(rawUrl, password)
	if err != nil {
		return nil, err
	}

	req := &dto.ReportReq{
		Provider:  info.Provider,
		Title:     info.Title,
		Size:      info.Size,
		Author:    info.Author,
		Status:    info.Status,
		ExpiredAt: info.ExpiredAt,
		RawURL:    rawUrl,             // 当前输入的原始链接
		URL:       info.NormalizedUrl, // 标准化链接 (唯一键)
		PWD:       info.Password,
	}

	if err := s.SaveResult(ctx, req); err != nil {
		return nil, err
	}

	return s.repo.GetByUrl(ctx, info.NormalizedUrl)
}

func (s *LinkService) GetLinkList(ctx context.Context, page, size int, keyword string) ([]*model.LinkRecord, int64, error) {
	return s.repo.GetList(ctx, page, size, keyword)
}

func (s *LinkService) SaveResult(ctx context.Context, req *dto.ReportReq) error {
	existing, err := s.repo.GetByUrl(ctx, req.URL)

	if err == gorm.ErrRecordNotFound || existing == nil {
		link := &model.LinkRecord{
			Provider:  req.Provider,
			Title:     req.Title,
			Size:      req.Size,
			Author:    req.Author,
			ExpiredAt: req.ExpiredAt,
			RawURL:    []string{req.RawURL},
			URL:       req.URL,
			PWD:       req.PWD,
			Status:    req.Status,
		}
		return s.repo.Create(ctx, link)
	}

	if err != nil {
		return err
	}

	needUpdate := false

	if !contains(existing.RawURL, req.RawURL) {
		existing.RawURL = append(existing.RawURL, req.RawURL)
		needUpdate = true
	}

	if time.Since(existing.UpdatedAt) > 7*24*time.Hour {
		existing.Title = req.Title
		existing.Size = req.Size
		existing.Status = req.Status
		existing.ExpiredAt = req.ExpiredAt
		existing.Author = req.Author
		needUpdate = true
	} else {
		if !needUpdate {
			return nil
		}
	}

	if needUpdate {
		return s.repo.Save(ctx, existing)
	}
	return nil
}

// 辅助函数：判断切片是否包含字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
