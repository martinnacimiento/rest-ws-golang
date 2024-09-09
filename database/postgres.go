package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}
