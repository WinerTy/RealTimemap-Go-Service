package category

import (
	"context"
	"realtimemap-service/internal/handler/dto"
)

type ServiceCategory interface {
	GetAll(ctx context.Context) ([]dto.CategoryResponse, error)
}
