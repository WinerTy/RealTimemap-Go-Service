package postgres

import (
	"context"
	"log/slog"
	"realtimemap-service/internal/domain/mark"
	"realtimemap-service/internal/pkg/logger/sl"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgMarkRepository struct {
	db *pgxpool.Pool
}

func NewPgMarkRepository(pool *pgxpool.Pool) mark.Repository {
	return &PgMarkRepository{db: pool}
}

func (r *PgMarkRepository) GetByOwner(ctx context.Context, ownerID int) ([]*mark.Mark, error) {
	rows, err := r.db.Query(ctx, "SELECT id, owner_id, mark_name, ST_AsGeoJSON(geom) FROM marks WHERE owner_id = $1", ownerID)
	if err != nil {
		slog.Error("GetByOwner err:", sl.Err(err))
	}
	defer rows.Close()

	var marks []*mark.Mark

	for rows.Next() {
		var item mark.Mark

		err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Geom)
		if err != nil {
			slog.Error("GetByOwner err:", sl.Err(err))
		}
		marks = append(marks, &item)
	}

	if err := rows.Err(); err != nil {
		slog.Error("GetByOwner err:", sl.Err(err))
		return nil, err
	}

	return marks, nil
}
