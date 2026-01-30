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
		DoUpdates: clause.AssignmentColumns([]string{"title", "size", "author", "expired_at", "pwd"}),
	}).Create(&lr).Error
}

func (l *LinkDB) List(page, size int) (lrs []model.LinkRecord, count int64, err error) {
	lrDB := l.db.Model(&model.LinkRecord{})
	lrDB.Count(&count)
	err = lrDB.Order("updated_at DESC").Limit((page) * size).Where("status = 1").Find(&lrs).Error
	return lrs, count, err
}
