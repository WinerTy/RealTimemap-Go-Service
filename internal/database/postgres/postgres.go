package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func NewStorage(dbUrl string) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}
	return &Storage{Db: db}, nil
}
