package config

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/pkg/logger/sl"

	"github.com/redis/go-redis/v9"
)

func SetupRedis(ctx context.Context, cfg *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Url,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		slog.Error("Error connecting to redis ", sl.Err(err))
		return nil, err
	}

	return rdb, nil
}
