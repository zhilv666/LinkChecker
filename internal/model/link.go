package model

import (
	"time"

	"github.com/zhilv666/linkchecker/internal/netdisk"
	"gorm.io/gorm"
)

type LinkRecord struct {
	gorm.Model

	Provider  string              `gorm:"type:text" json:"provider"`
	Title     string              `gorm:"type:text" json:"title"`
	Size      string              `gorm:"type:text" json:"size"`
	Author    string              `gorm:"type:text" json:"author"`
	Status    netdisk.ShareStatus `gorm:"type:text" json:"status"`
	ExpiredAt *time.Time          `json:"expired_at"`
	URL       string              `gorm:"uniqueIndex;type:text" json:"url"`
	PWD       string              `gorm:"type:text" json:"pwd"`
	RawURL    []string            `gorm:"type:text;serializer:json" json:"raw_url"`
}

// 增加这个方法，Go 的模板可以直接调用它
func (i LinkRecord) CreatedAtStr() string {
	// 2006-01-02 15:04:05 是 Go 固定的格式化布局，不能改数字
	return i.CreatedAt.Format("2006-01-02 15:04:05")
}
