package config

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"

	"github.com/redis/go-redis/v9"
)

// InitCache Функция создает стратегию кэшировнаия для приложения на основе кофнига
func InitCache(ctx context.Context, cfg *Config) (cache.Store, *redis.Client) {
	switch cfg.CacheStrategy {
	case "redis":
		redisCli, err := SetupRedis(ctx, cfg)
		if err != nil {
			slog.Warn("Error setting up redis cache. Cache will be disabled", sl.Err(err))
			return cache.NewNoOpCache(), redisCli
		}
		return cache.NewRedisCache(redisCli), redisCli

	case "none":
		return cache.NewNoOpCache(), nil
	default:
		slog.Warn("Unknown cache strategy. Cache will be disabled", cfg.CacheStrategy)
		return cache.NewNoOpCache(), nil
	}
}
