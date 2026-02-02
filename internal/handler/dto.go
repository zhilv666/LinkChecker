package handler

import "time"

type CheckReq struct {
	URL      string `json:"url"`
	Password string `json:"password"`
}

type CheckResp struct {
	Provider  string     `gorm:"type:text" json:"provider"`
	Title     string     `gorm:"type:text" json:"title"`
	Size      string     `gorm:"type:text" json:"size"`
	Author    string     `gorm:"type:text" json:"text"`
	ExpiredAt *time.Time `json:"expired_at"`
	URL       string     `gorm:"index;type:text" json:"url"`
	PWD       string     `gorm:"type:text" json:"pwd"`
}

type ListReq struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Keyword string `json:"Keyword"`
}
