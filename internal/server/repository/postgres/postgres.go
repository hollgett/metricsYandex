package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hollgett/metricsYandex.git/internal/server/database"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	gauge       = "gauge"
	counter     = "counter"
	tableScheme = `CREATE TABLE IF NOT EXISTS metrics (
"name" VARCHAR(70) PRIMARY KEY,
"type" VARCHAR(10) NOT NULL,
"delta" BIGINT NOT NULL DEFAULT 0,
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
	getMetricValue = `SELECT delta, value FROM metrics WHERE name = :name`
	getMetricAll   = `SELECT * FROM metrics`
)

var (
	ErrMetricTypeUnknown = errors.New("unknown metric")
	retryTimeSleep       = []int{1, 3, 5}
)

type Postgres struct {
	db         *sqlx.DB
	log        logger.Logger
	saveStmt   *sqlx.Stmt
	getStmt    *sqlx.NamedStmt
	getAllStmt *sqlx.Stmt
}

func New(ctx context.Context, dsn string, log logger.Logger) (repository.Repository, error) {
	db, err := database.Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}
	if _, err := db.ExecContext(ctx, tableScheme); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	saveStmt, err := db.PreparexContext(ctx, saveAndUpdQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare save: %w", err)
	}
	getStmt, err := db.PrepareNamedContext(ctx, getMetricValue)
	if err != nil {
		return nil, fmt.Errorf("prepare get value: %w", err)
	}
	getAll, err := db.PreparexContext(ctx, getMetricAll)
	if err != nil {
		return nil, fmt.Errorf("prepare get all: %w", err)
	}

	return &Postgres{
		db:         db,
		log:        log,
		saveStmt:   saveStmt,
		getStmt:    getStmt,
		getAllStmt: getAll,
	}, nil
}

func (pg *Postgres) Save(metric models.Metrics) error {
	for _, timeSleep := range retryTimeSleep {
		if _, err := pg.saveStmt.Exec(metric.ID, metric.MType, metric.Delta, metric.Value); err != nil {
			if isRetry(err) {
				pg.log.LogAny("retry request", "code", err)
				time.Sleep(time.Duration(timeSleep) * time.Second)
				continue
			}
			return err
		}
	}
	return nil
}

func (pg *Postgres) Get(metric *models.Metrics) error {
	val := new(float64)
	delta := new(int64)
	row := pg.getStmt.QueryRow(metric)
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

func (pg *Postgres) Batch(ctx context.Context, metrics []models.Metrics) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ctxTx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := pg.db.BeginTx(ctxTx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctxTx, saveAndUpdQuery)
	if err != nil {
		return err
	}
	for _, v := range metrics {
		if _, err := stmt.Exec(v.ID, v.MType, v.Delta, v.Value); err != nil {
			pg.log.LogAny("BATCH ERROR", "VALUE", v)
			pg.log.LogAny("FULL REQUEST BATCH", "value", metrics)
			return err
		}
	}

	return tx.Commit()
}

func isRetry(err error) bool {

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.ConnectionDoesNotExist, pgerrcode.ConnectionFailure:
			return true

		}
	}
	return false
}
