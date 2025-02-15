package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const postgresDriveName = "pgx"

// host=localhost user=postgres password=root dbname=USvideos sslmode=disable
func Connect(dbDSN string) (*sqlx.DB, error) {
	db, err := sqlx.Open(postgresDriveName, dbDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
