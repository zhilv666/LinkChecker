package netdisk

import (
	"encoding/json"
	"fmt"
	"time"
)

// ShareStatus 定义网盘的状态枚举
type ShareStatus int

const (
	StatusUnknown      ShareStatus = iota // 未知链接/解析失败
	StatusValid                           // 链接有效
	StatusDeleted                         // 文件已删除/链接不存在
	StatusExpired                         // 链接已过期
	StatusNeedPassword                    // 需要提取码（未提供或错误）
	StatusBanned                          // 内容违规被封禁
)

// String 用于打印状态的文本描述
func (s ShareStatus) String() string {
	switch s {
	case StatusValid:
		return "有效"
	case StatusDeleted:
		return "已删除"
	case StatusExpired:
		return "已过期"
	case StatusNeedPassword:
		return "需要密码"
	case StatusBanned:
		return "已封禁"
	default:
		return "未知"
	}
}

// ShareInfo 包含解析后的详情信息
type ShareInfo struct {
	Status        ShareStatus `json:"status"`
	Provider      string      `json:"provider"`
	Title         string      `json:"title,omitempty"`
	Size          string      `json:"size,omitempty"`
	Author        string      `json:"author,omitempty"`
	ExpiredAt     *time.Time  `json:"expired_at,omitempty"`
	RawUrl        string      `json:"raw_url"`
	Password      string      `json:"password"`
	NormalizedUrl string      `json:"normalized_url,omitempty"`
}

func (s ShareInfo) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("Error marshalling ShareInfo: %v", err)
	}
	return string(data)
}
