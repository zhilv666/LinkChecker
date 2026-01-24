package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(redisAddr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (c *RedisCache) Get(key string) (CacheValue, bool, error) {
	val, error := c.client.Get(key).Result()
	if error == redis.Nil {
		return new(CacheValue), false, nil
	}

	if error != nil {
		return new(CacheValue), false, error
	}

	return val, true, nil
}

func (c *RedisCache) Set(key string, value CacheValue, ttl time.Duration) error {
	return c.client.Set(key, value, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(key).Err()
}
