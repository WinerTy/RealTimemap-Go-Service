package category

import (
	"context"
	"realtimemap-service/internal/handler/dto"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]dto.CategoryResponse, error)
}
