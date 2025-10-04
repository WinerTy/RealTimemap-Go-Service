package category

import (
	"context"
	"realtimemap-service/internal/pkg/pagination"
)

type Service interface {
	GetAll(ctx context.Context, page, pageSize int) (*pagination.Response[CategoryResponse], error)
}

type Repository interface {
	GetAll(ctx context.Context, pageSize, offset int) ([]*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id int) (bool, error)
}
