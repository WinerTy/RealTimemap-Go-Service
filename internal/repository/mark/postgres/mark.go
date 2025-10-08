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
            m.id,
            m.mark_name,
            m.owner_id,
            m.additional_info,
            m.photo,
            ST_AsGeoJSON(m.geom) as geom,
            m.is_ended,
            (m.start_at + (m.duration * INTERVAL '1 hour')) as end_at,
            c.id as category_id,
            c.category_name,
            c.color,
            c.icon
        FROM
            marks m
        LEFT JOIN categories c ON m.category_id = c.id
        WHERE
            m.geohash = ANY($1)
            AND ST_DWithin(
                m.geom,
                ST_SetSRID(ST_MakePoint($2, $3), $4),
                $5
            )
            AND m.start_at <= $7
            AND (m.start_at + (m.duration * INTERVAL '1 hour')) > $6
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
		err := rows.Scan(&item.ID, &item.Name, &item.OwnerID, &item.AdditionalInfo,
			&item.Photo, &item.Geom, &item.IsEnded, &item.EndAt,
			&item.Category.ID, &item.Category.Name, &item.Category.Color, &item.Category.Icon)
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
