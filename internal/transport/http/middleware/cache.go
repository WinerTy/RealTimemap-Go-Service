package middleware

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/pkg/logger/sl"
	"time"

	"github.com/gin-gonic/gin"
)

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

		// Формирования ключа и контекста

		key := c.Request.URL.RequestURI()
		ctx := c.Request.Context()

		// Проверка Кэша, и отдаем кэшированный ответ
		if item, found := store.Get(ctx, key); found {
			c.Header("Cache-Control", "max-age="+fmt.Sprint(duration/time.Second))
			c.Header("X-Cache", "HIT")
			for h, v := range item.Headers {
				c.Writer.Header()[h] = v
			}
			c.Data(item.StatusCode, c.GetHeader("Content-Type"), item.Value)
			c.Abort()
			return
		}

		blw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Продолжаем запрос
		c.Next()

		slog.Debug("Запрос, Кэша нет")

		c.Header("X-Cache", "MISS")

		if c.Writer.Status() == http.StatusOK {
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
