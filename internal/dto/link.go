package dto

import (
	"time"

	"github.com/zhilv666/linkchecker/internal/netdisk"
)

type CheckReq struct {
	URL      string `json:"url" form:"url"`
	Password string `json:"password" form:"password"`
}

type CheckResp struct {
	Provider  string              `json:"provider"`
	Title     string              `json:"title"`
	Size      string              `json:"size"`
	Author    string              `json:"text"`
	ExpiredAt *time.Time          `json:"expired_at"`
	URL       string              `json:"url"`
	Status    netdisk.ShareStatus `json:"status"`
	PWD       string              `json:"pwd"`
}

type ListReq struct {
	Page    int    `json:"page" form:"page"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
}

type ReportReq struct {
	Provider  string              `json:"provider"`
	Title     string              `json:"title"`
	Size      string              `json:"size"`
	Author    string              `json:"author"`
	ExpiredAt *time.Time          `json:"expired_at"`
	RawURL    string              `json:"raw_url"`
	URL       string              `json:"url"`
	Status    netdisk.ShareStatus `json:"status"`
	PWD       string              `json:"pwd"`
}
