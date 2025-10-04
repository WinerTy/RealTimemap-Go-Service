package category

import (
	"context"
	"math"
	"realtimemap-service/internal/domain/category"

	categoryrepo "realtimemap-service/internal/domain/category"
	"realtimemap-service/internal/pkg/pagination"
	"sync"
)

type serviceCategory struct {
	categoryRepo categoryrepo.Repository
}

func NewServiceCategory(categoryRepo categoryrepo.Repository) category.Service {
	return &serviceCategory{categoryRepo}
}

func (s *serviceCategory) GetAll(ctx context.Context, page, pageSize int) (*category.PaginationCategoryResponse, error) {
	offset := pagination.Offset(page, pageSize)
	var categories []*category.Category
	var categoryErr, countErr error
	var total int
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		categories, categoryErr = s.categoryRepo.GetAll(ctx, pageSize, offset)
	}()

	go func() {
		defer wg.Done()
		total, countErr = s.categoryRepo.Count(ctx)
	}()
	wg.Wait()

	if categoryErr != nil {
		return nil, categoryErr
	}

	if countErr != nil {
		return nil, countErr
	}
	serializedData := category.ToListCategoryResponse(categories)
	response := &category.PaginationCategoryResponse{
		Result:     serializedData,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}
	return response, nil

}
