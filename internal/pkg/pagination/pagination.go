package pagination

type PaginationResponse[T any] struct {
	Result     []T `json:"result"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func Offset(page, pageSize int) int {
	return (page - 1) * pageSize
}
