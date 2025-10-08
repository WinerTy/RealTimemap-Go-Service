package mark

import (
	"log/slog"
	"realtimemap-service/internal/domain/category"
	"time"
)

type Response struct {
	ID             int               `json:"id"`
	Name           string            `json:"mark_name"`
	OwnerID        int               `json:"owner_id"`
	AdditionalInfo *string           `json:"additional_info"`
	Geom           Geometry          `json:"geom"`
	IsEnded        bool              `json:"is_ended"`
	DurationHours  int               `json:"duration"`
	EndAt          time.Time         `json:"end_at"`
	Photo          []string          `json:"photo"`
	Category       category.Response `json:"category"`
}

func ToMarkResponse(item *Mark) Response {
	photoUrls := make([]string, len(item.Photo))
	for i, photo := range item.Photo {
		if photoUrl := photo.BuildUrl(); photoUrl != "" {
			photoUrls[i] = photoUrl
		}
	}
	return Response{
		ID:             item.ID,
		Name:           item.Name,
		OwnerID:        item.OwnerID,
		AdditionalInfo: item.AdditionalInfo,
		Geom:           item.Geom,
		IsEnded:        item.IsEnded,
		DurationHours:  item.DurationHours,
		EndAt:          item.EndAt,
		Photo:          photoUrls,
		Category:       category.ToCategoryResponse(item.Category),
	}
}

func ToListMarkResponse(items []*Mark) []Response {
	slog.Info("Длинна ответа: ", len(items))
	responseList := make([]Response, len(items))
	for i, item := range items {
		responseList[i] = ToMarkResponse(item)
	}
	return responseList
}
