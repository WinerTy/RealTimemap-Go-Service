package category

import (
	"context"
	"realtimemap-service/internal/entity"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*entity.Category, error)
	GetByID(ctx context.Context, id int) (*entity.Category, error)
}
