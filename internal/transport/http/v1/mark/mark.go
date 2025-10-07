package mark

import (
	"net/http"
	"realtimemap-service/internal/domain/mark"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmcloughlin/geohash"
)

type Routes struct {
	repo mark.Repository
}

type Filter struct {
	// Геоданные
	Longitude float64 `form:"longitude" binding:"required,longitude"`
	Latitude  float64 `form:"latitude" binding:"required,latitude"`
	Radius    float64 `form:"radius" binding:"omitempty,gt=0,max=10000"`
	SRID      int     `form:"srid" binding:"omitempty,gt=0"`

	// Временные
	ReferenceTime time.Time `form:"date" binding:"omitempty"  time_format:"2006-01-02"`
	DurationHours int       `form:"duration" binding:"omitempty,oneof=12 24 48"`
	ShowEnded     bool      `form:"show_ended"`
}

func (f *Filter) SetTimeRange() {
	if f.DurationHours == 0 {
		f.DurationHours = mark.DefaultDuration
	}
	if f.ReferenceTime.IsZero() {
		f.ReferenceTime = time.Now()
	}
}

func (f *Filter) SetDefault() {
	if f.Radius == 0 {
		f.Radius = mark.DefaultRadius
	}
	if f.SRID == 0 {
		f.SRID = mark.DefaultSRID
	}
	f.SetTimeRange()
}

func InitMarkRoutes(g *gin.RouterGroup, repo mark.Repository) {
	r := &Routes{repo: repo}

	marksRoutes := g.Group("/")
	{
		marksRoutes.GET("/", r.GetNearest)
	}
}

func (r *Routes) GetNearest(c *gin.Context) {
	var userFilter Filter

	if err := c.ShouldBindQuery(&userFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userFilter.SetDefault()
	domainFilter := mark.NewFilter(
		userFilter.Latitude,
		userFilter.Longitude,
		userFilter.Radius,
		userFilter.SRID,
		userFilter.DurationHours,
		userFilter.ShowEnded,
		userFilter.ReferenceTime,
	)
	geo := geohash.EncodeWithPrecision(userFilter.Latitude, userFilter.Longitude, 5)
	res, err := r.repo.GetNearestMarks(c.Request.Context(), domainFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message":       "GetNearest",
		"user_filter":   userFilter,
		"domain_filter": domainFilter,
		"response":      res,
		"debug":         geo,
		"neighbors":     geohash.Neighbors(geo),
	})
}
