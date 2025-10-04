package category

import "context"

type Service interface {
	GetAll(ctx context.Context, page, pageSize int) (*PaginationCategoryResponse, error)
}

type Repository interface {
	GetAll(ctx context.Context, pageSize, offset int) ([]*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	Count(ctx context.Context) (int, error)
}
