package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"realtimemap-service/internal/pkg/logger/sl"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var NoRedisCli error

func NewRedisCache(client *redis.Client) Store {
	return &RedisCache{client: client}
}

func (rc *RedisCache) Get(ctx context.Context, key string) (CacheItem, bool) {
	val, err := rc.client.Get(ctx, key).Result()

	// Проверка найден ли значение по ключу
	if errors.Is(err, redis.Nil) { // Ключ не найден
		return CacheItem{}, false
	} else if err != nil { // Другая ошибка
		return CacheItem{}, false
	}

	// Парсинг значения
	var item CacheItem
	err = json.Unmarshal([]byte(val), &item)
	if err != nil {
		slog.Error("Error unmarshalling cache item:", sl.Err(err))
		return CacheItem{}, false
	}
	return item, true
}

// Set Функция для добавления в Redis
func (rc *RedisCache) Set(ctx context.Context, key string, item CacheItem, ttl time.Duration) error {
	data, err := json.Marshal(item)
	if err != nil {
		slog.Error("Error marshalling cache item:", sl.Err(err))
		return err
	}
	return rc.client.Set(ctx, key, data, ttl).Err()
}
