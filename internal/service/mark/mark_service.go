package mark

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/domain/mark"
)

type Service struct {
	markRepo mark.Repository
}

func NewService(markRepo mark.Repository) mark.Service {
	return &Service{markRepo: markRepo}
}

func (s *Service) GetNearestMark(ctx context.Context, req mark.FilterRequest) ([]mark.Response, error) {
	filters := mark.NewFilter(
		req.Latitude,
		req.Longitude,
		req.Radius,
		req.SRID,
		req.DurationHours,
		req.ShowEnded,
		req.ReferenceTime,
	)

	marks, err := s.markRepo.GetNearestMarks(ctx, filters)
	if err != nil {
		slog.Error("GetNearestMarks err", err)
		return nil, err
	}

	response := mark.ToListMarkResponse(marks)

	return response, nil
}
