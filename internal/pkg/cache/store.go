package cache

import (
	"context"
	"time"
)

type Item struct {
	Value      []byte
	StatusCode int
	Headers    map[string][]string
	ExpiresAt  time.Time
}

type Store interface {
	Get(ctx context.Context, key string) (Item, bool)
	Set(ctx context.Context, key string, item Item, ttl time.Duration) error
}
