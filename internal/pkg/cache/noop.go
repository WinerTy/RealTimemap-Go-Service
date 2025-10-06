package cache

import (
	"context"
	"time"
)

type NoOpCache struct{}

func NewNoOpCache() Store {
	return &NoOpCache{}
}
func (n *NoOpCache) Get(_ context.Context, _ string) (CacheItem, bool) {
	return CacheItem{}, false
}

// Set просто игнорирует данные и всегда возвращает успех.
func (n *NoOpCache) Set(_ context.Context, _ string, _ CacheItem, _ time.Duration) error {
	return nil
}
