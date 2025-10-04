package category

type PaginationCategoryResponse struct {
	Result     []CategoryResponse `json:"result"`
	Total      int                `json:"total"`
	TotalPages int                `json:"total_pages"`
}
type CategoryResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"category_name"`
	Color string `json:"color"`
}

func ToCategoryResponse(category *Category) CategoryResponse {
	return CategoryResponse{
		ID:    category.ID,
		Name:  category.Name,
		Color: category.Color,
	}
}

func ToListCategoryResponse(categories []*Category) []CategoryResponse {
	listCategoryResponse := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		listCategoryResponse[i] = ToCategoryResponse(category)
	}
	return listCategoryResponse
}
