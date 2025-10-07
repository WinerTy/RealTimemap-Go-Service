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
		return nil, err
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

func (r *PgMarkRepository) GetNearestMarks(ctx context.Context, filter mark.Filter) ([]*mark.Mark, error) {
	query := `
		SELECT
			id,
			owner_id,
			mark_name,
			ST_AsGeoJSON(geom),
			is_ended
		FROM
			marks
		WHERE
		    (geohash LIKE ANY($1)) AND
			ST_DWithin(
				ST_Transform(geom, 3857),
				ST_Transform(ST_SetSRID(ST_MakePoint($2, $3), $4), 3857),
				$5
			)
			AND start_at <= $7
			AND (start_at + (duration * INTERVAL '1 hour')) > $6
		`

	rows, err := r.db.Query(ctx, query,
		filter.Geohash,           // $1
		filter.Longitude,         // $2
		filter.Latitude,          // $3
		filter.SRID,              // $4
		filter.Radius,            // $5
		filter.SearchWindowStart, // $6
		filter.SearchWindowEnd,   // $7
	)
	if err != nil {
		slog.Error("GetNearest err:", sl.Err(err))
		return nil, err
	}
	defer rows.Close()
	var marks []*mark.Mark
	for rows.Next() {
		var item mark.Mark
		err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Geom, &item.IsEnded)
		if err != nil {
			slog.Error("GetNearest err:", sl.Err(err))
			return nil, err
		}
		marks = append(marks, &item)
	}
	if err := rows.Err(); err != nil {
		slog.Error("GetNearest err:", sl.Err(err))
		return nil, err
	}
	return marks, nil
}
