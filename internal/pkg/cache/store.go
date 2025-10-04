package cache

import (
	"context"
	"time"
)

type CacheItem struct {
	Value      []byte
	StatusCode int
	Headers    map[string][]string
	ExpiresAt  time.Time
}

type Store interface {
	Get(ctx context.Context, key string) (CacheItem, bool)
	Set(ctx context.Context, key string, item CacheItem, ttl time.Duration) error
}
