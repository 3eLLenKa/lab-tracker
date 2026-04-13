package postgres

import "database/sql"

type Postgres struct {
	db *sql.DB
}

func New() (*Postgres, error) {
	return nil, nil
}
