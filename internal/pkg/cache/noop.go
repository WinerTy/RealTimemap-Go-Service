package cache

import (
	"context"
	"time"
)

type NoOpCache struct{}

func NewNoOpCache() Store {
	return &NoOpCache{}
}
func (n *NoOpCache) Get(_ context.Context, _ string) (Item, bool) {
	return Item{}, false
}

// Set просто игнорирует данные и всегда возвращает успех.
func (n *NoOpCache) Set(_ context.Context, _ string, _ Item, _ time.Duration) error {
	return nil
}
