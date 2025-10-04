package dto

import "realtimemap-service/internal/entity"

type CategoryResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"category_name"`
	Color string `json:"color"`
}

func ToCategoryResponse(category *entity.Category) CategoryResponse {
	return CategoryResponse{
		ID:    category.ID,
		Name:  category.Name,
		Color: category.Color,
	}
}

func ToListCategoryResponse(categories []*entity.Category) []CategoryResponse {
	listCategoryResponse := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		listCategoryResponse[i] = ToCategoryResponse(category)
	}
	return listCategoryResponse
}
