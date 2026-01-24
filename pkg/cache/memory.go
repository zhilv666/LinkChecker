package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	value  CacheValue
	expire time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]cacheItem),
	}
}

func (c *MemoryCache) Get(key string) (CacheValue, bool, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(item.expire) {
		return new(CacheValue), false, nil
	}

	return item.value, true, nil
}

func (c *MemoryCache) Set(key string, value CacheValue, ttl time.Duration) error {
	c.mu.Lock()
	c.items[key] = cacheItem{
		value:  value,
		expire: time.Now().Add(ttl),
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}
