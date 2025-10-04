package config

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"

	"github.com/redis/go-redis/v9"
)

func SetupRedis(ctx context.Context, cfg *Config) (*redis.Client, error) {
	if !cfg.Redis.Use {
		// Флаг для того чтоб не создавать клиент рдиса
		slog.Info("redis not used")
		return nil, cache.NoRedisCli
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		slog.Error("Error connecting to redis ", sl.Err(err))
		return nil, err
	}

	return rdb, nil
}
