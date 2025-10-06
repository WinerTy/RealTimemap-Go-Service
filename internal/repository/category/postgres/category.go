package postgres

import (
	"context"
	"realtimemap-service/internal/domain/category"
	repository "realtimemap-service/internal/domain/category"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCategoryRepository struct {
	db *pgxpool.Pool
}

// NewPgCategoryRepository Генератор
func NewPgCategoryRepository(pool *pgxpool.Pool) repository.Repository {
	return &PgCategoryRepository{db: pool}
}

// GetAll функция слоя репозиториев делает запрос к бд на получение актуальных Категорий
// TODO сделать обработку JSON поля на фотки
func (r *PgCategoryRepository) GetAll(ctx context.Context, pageSize, offset int) ([]*category.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, category_name, color FROM categories "+
		"WHERE is_active = true "+
		"LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*category.Category, 0)
	for rows.Next() {
		var c category.Category

		err := rows.Scan(&c.ID, &c.Name, &c.Color)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByID TODO написать метод по id
func (r *PgCategoryRepository) GetByID(ctx context.Context, id int) (*category.Category, error) {
	return nil, nil
}

func (r *PgCategoryRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE is_active = true").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PgCategoryRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exist bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
