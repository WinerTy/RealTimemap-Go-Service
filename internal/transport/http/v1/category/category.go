package category

import (
	"fmt"
	"net/http"
	"realtimemap-service/internal/domain/category"
	"realtimemap-service/internal/pkg/cache"
	"realtimemap-service/internal/transport/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryRoutes struct {
	categoryService category.Service
}

type Params struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func InitCategoryRoutes(g *gin.RouterGroup, categoryService category.Service, store cache.Store) {
	r := &CategoryRoutes{categoryService}

	cache1min := middleware.CacheMiddleware(store, time.Minute)
	categoryRoutes := g.Group("/category")

	categoryRoutes.Use(cache1min)
	{
		categoryRoutes.GET("/", r.GetAll)
	}
}

func (r *CategoryRoutes) GetAll(c *gin.Context) {
	var params Params

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.PageSize <= 0 {
		params.PageSize = 15
	}

	if params.PageSize > 150 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PageSize must be greater than 100"})
		return
	}

	response, err := r.categoryService.GetAll(c.Request.Context(), params.Page, params.PageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  fmt.Errorf("could not get categories: %w", err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   response,
	})
	return
}
