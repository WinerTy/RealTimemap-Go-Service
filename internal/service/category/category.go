package category

import (
	"context"
	"realtimemap-service/internal/handler/dto"
	"realtimemap-service/internal/repository/category"
)

type serviceCategory struct {
	categoryRepo category.Repository
}

func NewServiceCategory(categoryRepo category.Repository) ServiceCategory {
	return &serviceCategory{categoryRepo}
}

func (s *serviceCategory) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	serializedData := dto.ToListCategoryResponse(categories)

	return serializedData, nil

}
