package v1

import (
	"fmt"
	"net/http"
	"realtimemap-service/internal/service/category"

	"github.com/gin-gonic/gin"
)

type CategoryRoutes struct {
	categoryService category.CategoryService
}

func InitCategoryRoutes(g *gin.RouterGroup, categoryService category.CategoryService) {
	r := &CategoryRoutes{categoryService}

	categoryRoutes := g.Group("/category")
	{
		categoryRoutes.GET("/", r.GetAll)
	}
}

func (r *CategoryRoutes) GetAll(c *gin.Context) {
	categories, err := r.categoryService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  fmt.Errorf("could not get categories: %w", err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   categories,
	})
	return
}
