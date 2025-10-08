package category

import (
	"realtimemap-service/internal/pkg/entity"
)

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"category_name"`
	Color string `json:"color"`
	Icon  entity.Image
}
