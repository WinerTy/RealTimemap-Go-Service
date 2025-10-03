package category

import (
	"database/sql"
	model "realtimemap-service/internal/app/models/category"
)

type Repository interface {
	GetAll() ([]*model.Category, error)
	//GetByID(id int)
}

type PgCategoryRepository struct {
	db *sql.DB
}

func NewPgCategoryRepository(db *sql.DB) Repository {
	return &PgCategoryRepository{db: db}
}

func (r *PgCategoryRepository) GetAll() ([]*model.Category, error) {
	rows, err := r.db.Query("SELECT id, category_name, color FROM categories WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]*model.Category, 0)
	for rows.Next() {
		var row model.Category
		err := rows.Scan(&row.ID, &row.Name, &row.Color)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &row)
	}
	return categories, nil
}
