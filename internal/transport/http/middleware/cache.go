package middleware

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

const HeaderCache = "X-Cache"

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CacheMiddleware(store cache.Store, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверка на GET Запрос
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// ПОлучение кэша по ключу на основе валидированных параметров
		cacheKey, exists := c.Get(CacheKey)
		if !exists {
			c.Next()
			return
		}

		// Преобразование
		key := cacheKey.(string)
		ctx := c.Request.Context()

		// Проверка Кэша, и отдаем кэшированный ответ
		if item, found := store.Get(ctx, key); found {
			c.Header(HeaderCache, "HIT")
			for h, v := range item.Headers {
				c.Writer.Header()[h] = v
			}
			c.Data(item.StatusCode, c.GetHeader("Content-Type"), item.Value)
			c.Abort()
			return
		}

		// Записываем ответ в кэш для последущих запросов
		blw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Продолжаем запрос
		c.Next()
		c.Header(HeaderCache, "MISS")

		if c.Writer.Status() == http.StatusOK {
			c.Header(HeaderCache, "HIT")
			item := cache.CacheItem{
				Value:      blw.body.Bytes(),
				StatusCode: c.Writer.Status(),
				Headers:    c.Writer.Header().Clone(),
			}

			err := store.Set(ctx, key, item, duration)
			if err != nil {
				slog.Error("Error caching item:", sl.Err(err))
			}
		}
	}
}

func generateCacheKey(path string, params interface{}) string {
	v := reflect.ValueOf(params).Elem()
	t := v.Type()

	queryValues := url.Values{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("form")
		if tag == "" {
			continue
		}
		queryValues.Add(tag, fmt.Sprintf("%v", field.Interface()))
	}

	return fmt.Sprintf("%s?%s", path, queryValues.Encode())
}
