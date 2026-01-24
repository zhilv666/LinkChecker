package cache

import "time"

type CacheType int
type CacheValue any

const (
	Memory CacheType = iota
	Redis
)

type Cache interface {
	// Get 获取缓存中的数据, 如果不存在返回 false 和错误
	Get(key string) (CacheValue, bool, error)

	// Set 设置缓存数据, 指定过期时间
	Set(key string, value CacheValue, ttl time.Duration) error

	// Delete 删除缓存数据
	Delete(key string) error
}

// 创建缓存配置
type Config struct {
	Type CacheType
	Addr string
}

// New 新建缓存, 默认使用内存作为缓存
func New(conf *Config) Cache {
	switch conf.Type {
	case Memory:
		return NewMemoryCache()
	case Redis:
		return NewRedisCache(conf.Addr)
	default:
		return NewMemoryCache()
	}
}
