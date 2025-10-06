package pagination

import "math"

type Response[T any] struct {
	Result     []T `json:"result"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPaginationResponse[T any](result []T, total, pageSize int) *Response[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}
	return &Response[T]{
		Result:     result,
		Total:      total,
		TotalPages: totalPages,
	}
}

// Offset Вычисления offset на основе параметров
func Offset(page, pageSize int) int {
	return (page - 1) * pageSize
}
