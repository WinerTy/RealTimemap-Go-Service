package middleware

import (
	"log/slog"
	"net/http"
	"realtimemap-service/internal/pkg/logger/sl"
	"reflect"

	"github.com/gin-gonic/gin"
)

const (
	ValidatedParams = "validated_params"
	CacheKey        = "cache_key"
)

func NormalizeQueryParams(paramStruct interface{}) gin.HandlerFunc {
	if reflect.ValueOf(paramStruct).Kind() != reflect.Struct {
		slog.Error("NormalizeQueryParams paramStruct must be a struct")
		panic("paramStruct must be a struct")
	}
	return func(c *gin.Context) {
		params := reflect.New(reflect.TypeOf(paramStruct)).Interface()
		if err := c.ShouldBindQuery(params); err != nil {
			slog.Info("NormalizeQueryParams ShouldBindQuery err:", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		key := generateCacheKey(c.Request.URL.Path, params)

		c.Set(ValidatedParams, params)
		c.Set(CacheKey, key)
		c.Next()
	}
}
