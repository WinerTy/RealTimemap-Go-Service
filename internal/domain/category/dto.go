package category

type Response struct {
	ID    int    `json:"id"`
	Name  string `json:"category_name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func ToCategoryResponse(category *Category) Response {
	return Response{
		ID:    category.ID,
		Name:  category.Name,
		Color: category.Color,
		Icon:  category.Icon.BuildUrl(),
	}
}

func ToListCategoryResponse(categories []*Category) []Response {
	listCategoryResponse := make([]Response, len(categories))
	for i, category := range categories {
		listCategoryResponse[i] = ToCategoryResponse(category)
	}
	return listCategoryResponse
}
