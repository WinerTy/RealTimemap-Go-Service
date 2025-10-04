package category

import (
	"context"
	model "realtimemap-service/internal/app/models/category"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*model.Category, error)
	GetByID(ctx context.Context, id int) (*model.Category, error)
}

type PgCategoryRepository struct {
	db *pgxpool.Pool
}

// NewPgCategoryRepository Генератор
func NewPgCategoryRepository(pool *pgxpool.Pool) Repository {
	return &PgCategoryRepository{db: pool}
}

// GetAll функция слоя репозиториев делает запрос к бд на получение актуальных Категорий
// TODO сделать обработку JSON поля на фотки
func (r *PgCategoryRepository) GetAll(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, category_name, color FROM categories WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var category model.Category

		err := rows.Scan(&category.ID, &category.Name, &category.Color)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByID TODO написать метод по id
func (r *PgCategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	return nil, nil
}
