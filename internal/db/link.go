package db

import (
	"github.com/zhilv666/linkchecker/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LinkDB struct {
	db *gorm.DB
}

func NewLinkDB(db *gorm.DB) *LinkDB {
	return &LinkDB{
		db: db,
	}
}

func (l *LinkDB) Create(lr *model.LinkRecord) error {
	return l.db.Model(&model.LinkRecord{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "size", "author", "expired_at", "pwd", "updated_at"}),
	}).Create(&lr).Error
}

func (l *LinkDB) List(page, size int, keyword string) (lrs []model.LinkRecord, count int64, err error) {
	query := l.db.Model(&model.LinkRecord{})
	if keyword != "" {
		likeStr := "%" + keyword + "%"
		query = query.Where("title LIKE ?", likeStr)
	}

	query = query.Where("status = 1")

	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("updated_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&lrs).Error
	return lrs, count, err
}
