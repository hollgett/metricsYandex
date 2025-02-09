package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const postgresDriveName = "pgx"

// host=localhost user=postgres password=root dbname=USvideos sslmode=disable
func Connect(dbDSN string) (*sql.DB, error) {
	db, err := sql.Open(postgresDriveName, dbDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
