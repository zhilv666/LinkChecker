package netdisk

import (
	"errors"
	"time"

	"github.com/zhilv666/linkchecker/pkg/cache"
)

// Checker 定义每个网盘需要实现的接口
type Provider interface {
	// Name 返回网盘名字
	Name() string

	// Match 用来判断该链接是否属于该网盘
	Match(url string) bool

	// Check 执行具体检测逻辑
	// password 是可选的提取码，如果没有传空字符串
	Check(rawUrl, password string) (*ShareInfo, error)
}

// Manager 管理所有网盘解析器
type Manager struct {
	providers []Provider
	cache     cache.Cache
}

// NewManager 创建并注册默认支持的网盘
func NewManager(cacheInstance cache.Cache, providers ...Provider) *Manager {
	return &Manager{
		providers: providers,
		cache:     cacheInstance,
	}
}

// Register 动态注册新的网盘解析器
func (m *Manager) Register(c Provider) {
	m.providers = append(m.providers, c)
}

func (m *Manager) Check(url, password string) (*ShareInfo, error) {
	if result, found, err := m.cache.Get(url); err == nil && found {
		info, ok := result.(*ShareInfo)
		if ok {
			return info, nil
		}
	}

	for _, checker := range m.providers {
		if checker.Match(url) {
			info, err := checker.Check(url, password)
			if err != nil {
				return nil, err
			}
			m.cache.Set(url, info, 10*time.Minute)
			return info, nil
		}
	}
	return nil, errors.New("不支持的网盘链接")
}
