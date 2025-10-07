package mark

import "realtimemap-service/internal/domain/mark"

type Service struct {
	markRepo mark.Repository
}

func NewService(markRepo mark.Repository) mark.Service {
	return &Service{markRepo: markRepo}
}
