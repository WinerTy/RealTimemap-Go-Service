package mark

import (
	"net/http"
	"realtimemap-service/internal/domain/mark"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	service mark.Service
}

func InitMarkRoutes(g *gin.RouterGroup, service mark.Service) {
	r := &Routes{service: service}

	marksRoutes := g.Group("/")
	{
		marksRoutes.GET("/", r.GetNearest)
	}
}

func (r *Routes) GetNearest(c *gin.Context) {
	var userFilter mark.FilterRequest

	if err := c.ShouldBindQuery(&userFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userFilter.SetDefault()
	response, err := r.service.GetNearestMark(c.Request.Context(), userFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
