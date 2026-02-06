package repo

import (
	"context"

	"github.com/zhilv666/linkchecker/internal/model"
	"gorm.io/gorm"
)

type ILinkRepo interface {
	Create(ctx context.Context, link *model.LinkRecord) error
	Save(ctx context.Context, link *model.LinkRecord) error
	GetList(ctx context.Context, page, size int, keyword string) ([]*model.LinkRecord, int64, error)
	GetByRawUrl(ctx context.Context, url string) (*model.LinkRecord, error)
	GetByUrl(ctx context.Context, url string) (*model.LinkRecord, error)
	UpdateStatus(ctx context.Context, id uint, status int) error
}

type linkRepo struct {
	db *gorm.DB
}

func NewLinkRepo(db *gorm.DB) *linkRepo {
	return &linkRepo{db: db}
}

func (r *linkRepo) Create(ctx context.Context, link *model.LinkRecord) error {
	return r.db.WithContext(ctx).Create(&link).Error
}

func (r *linkRepo) Save(ctx context.Context, link *model.LinkRecord) error {
	return r.db.WithContext(ctx).Save(link).Error
}

func (r *linkRepo) GetList(ctx context.Context, page, size int, keyword string) ([]*model.LinkRecord, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.LinkRecord{})

	if keyword != "" {
		linkStr := "%" + keyword + "%"
		query = query.Where("title LIKE ?", linkStr)
	}

	query = query.Where("status = 1")
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var links []*model.LinkRecord
	err = query.
		Order("updated_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&links).Error
	return links, count, err
}

func (r *linkRepo) GetByRawUrl(ctx context.Context, url string) (*model.LinkRecord, error) {
	var link model.LinkRecord
	err := r.db.WithContext(ctx).Model(&model.LinkRecord{}).Where("raw_url = ?", url).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, err
}

func (r *linkRepo) GetByUrl(ctx context.Context, url string) (*model.LinkRecord, error) {
	var link model.LinkRecord
	err := r.db.WithContext(ctx).Model(&model.LinkRecord{}).Where("url = ?", url).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, err
}

func (r *linkRepo) UpdateStatus(ctx context.Context, id uint, status int) error {
	query := r.db.WithContext(ctx).Model(&model.LinkRecord{})
	return query.Where("id = ?", id).Update("status", status).Error
}
