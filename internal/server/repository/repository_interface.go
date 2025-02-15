package repository

import (
	"context"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

//go:generate mockgen -source=repository_interface.go -destination=../mock/repository.go -package=mock
type Repository interface {
	Save(metric models.Metrics) error
	Get(metric *models.Metrics) error
	GetAll() ([]models.Metrics, error)
	Ping(ctx context.Context) error
	Batch(ctx context.Context, metrics []models.Metrics) error
	Close() error
}
