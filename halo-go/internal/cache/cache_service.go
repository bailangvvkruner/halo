package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	prefix string
}

func NewCacheService(addr, password string, db int, prefix string) (*CacheService, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return &CacheService{
		client: rdb,
		prefix: prefix,
	}, nil
}

func (c *CacheService) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return val, nil
	}
	return result, nil
}

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return c.client.Set(ctx, c.prefix+key, data, expiration).Err()
}

func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.prefix + k
	}
	return c.client.Del(ctx, prefixedKeys...).Err()
}

func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, c.prefix+key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *CacheService) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.prefix+key).Result()
}

func (c *CacheService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, c.prefix+key, expiration).Err()
}

func (c *CacheService) Close() error {
	return c.client.Close()
}

type MemoryCache struct {
	store map[string]cacheEntry
	mu    sync.RWMutex
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]cacheEntry),
	}
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.store[key]
	if !ok {
		return nil, false
	}
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		delete(m.store, key)
		return nil, false
	}
	return entry.value, true
}

func (m *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	m.store[key] = cacheEntry{value: value, expiration: exp}
}

func (m *MemoryCache) Delete(keys ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.store, k)
	}
}

func (m *MemoryCache) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.store[key]
	if !ok {
		return false
	}
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		delete(m.store, key)
		return false
	}
	return true
}
