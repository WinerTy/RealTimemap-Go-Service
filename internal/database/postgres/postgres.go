package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewStorage(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	const fn = "storage.postgres.New"

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return pool, nil
}
