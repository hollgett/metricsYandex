package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hollgett/metricsYandex.git/internal/server/database"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
)

const (
	gauge       = "gauge"
	counter     = "counter"
	tableScheme = `CREATE TABLE IF NOT EXISTS metrics (
"name" VARCHAR(70) PRIMARY KEY,
"type" VARCHAR(10) NOT NULL,
"delta" INT NOT NULL DEFAULT 0,
"value" double precision NOT NULL DEFAULT 0,
CONSTRAINT unique_id_type UNIQUE (name, type)
);`
	saveAndUpdQuery = `INSERT INTO metrics (name, type, delta, value) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (name, type) 
DO UPDATE SET 
    delta = CASE 
                WHEN EXCLUDED.type = 'counter' THEN metrics.delta + EXCLUDED.delta 
                ELSE EXCLUDED.delta 
            END,
    value = CASE 
                WHEN EXCLUDED.type = 'gauge' THEN EXCLUDED.value 
                ELSE EXCLUDED.value 
            END;`
	getMetricValue = `SELECT delta, value FROM metrics WHERE name = $1`
	getMetricAll   = `SELECT * FROM metrics`
)

var (
	ErrMetricTypeUnknown = errors.New("unknown metric")
)

type Postgres struct {
	db         *sql.DB
	saveStmt   *sql.Stmt
	getStmt    *sql.Stmt
	getAllStmt *sql.Stmt
}

func New(ctx context.Context, dsn string) (repository.Repository, error) {
	db, err := database.Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}
	if _, err := db.ExecContext(ctx, tableScheme); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	saveStmt, err := db.PrepareContext(ctx, saveAndUpdQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare save: %w", err)
	}
	getStmt, err := db.PrepareContext(ctx, getMetricValue)
	if err != nil {
		return nil, fmt.Errorf("prepare get value: %w", err)
	}
	getAll, err := db.PrepareContext(ctx, getMetricAll)
	if err != nil {
		return nil, fmt.Errorf("prepare get all: %w", err)
	}

	return &Postgres{
		db:         db,
		saveStmt:   saveStmt,
		getStmt:    getStmt,
		getAllStmt: getAll,
	}, nil
}

func (pg *Postgres) Save(data models.Metrics) error {
	switch data.MType {
	case gauge:

		if _, err := pg.saveStmt.Exec(data.ID, data.MType, 0, data.Value); err != nil {
			return err
		}
		return nil
	case counter:
		if _, err := pg.saveStmt.Exec(data.ID, data.MType, data.Delta, 0); err != nil {
			return err
		}
		return nil
	default:
		return ErrMetricTypeUnknown
	}
}
func (pg *Postgres) Get(metric *models.Metrics) error {
	val := new(float64)
	delta := new(int64)
	row := pg.getStmt.QueryRow(metric.ID)
	if row.Err() != nil {
		return row.Err()
	}
	row.Scan(delta, val)
	metric.Delta = delta
	metric.Value = val
	return nil
}
func (pg *Postgres) GetAll() ([]models.Metrics, error) {
	rows, err := pg.getAllStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var metrics []models.Metrics
	for rows.Next() {
		var metric models.Metrics
		delta := new(int64)
		value := new(float64)
		err := rows.Scan(&metric.ID, &metric.MType, delta, value)
		if err != nil {
			return nil, err
		}
		metric.Delta = delta
		metric.Value = value
		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}
func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}
func (pg *Postgres) Close() error {
	var errs []error
	if pg.saveStmt != nil {
		if err := pg.saveStmt.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close saveStmt: %w", err))
		}
	}
	if pg.getStmt != nil {
		if err := pg.getStmt.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close getStmt: %w", err))
		}
	}
	if pg.getAllStmt != nil {
		if err := pg.getAllStmt.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close getAllStmt: %w", err))
		}
	}
	if err := pg.db.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close db: %w", err))
	}
	return errors.Join(errs...)
}
