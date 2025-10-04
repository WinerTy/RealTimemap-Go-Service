package postgres

import (
	"context"
	"realtimemap-service/internal/entity"
	"realtimemap-service/internal/repository/category"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCategoryRepository struct {
	db *pgxpool.Pool
}

// NewPgCategoryRepository Генератор
func NewPgCategoryRepository(pool *pgxpool.Pool) category.Repository {
	return &PgCategoryRepository{db: pool}
}

// GetAll функция слоя репозиториев делает запрос к бд на получение актуальных Категорий
// TODO сделать обработку JSON поля на фотки
func (r *PgCategoryRepository) GetAll(ctx context.Context) ([]*entity.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, category_name, color FROM categories WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*entity.Category, 0)
	for rows.Next() {
		var category entity.Category

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
func (r *PgCategoryRepository) GetByID(ctx context.Context, id int) (*entity.Category, error) {
	return nil, nil
}
