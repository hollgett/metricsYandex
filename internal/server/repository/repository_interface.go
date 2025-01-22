package repository

import (
	"net/http"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

//go:generate mockgen -source=repository_interface.go -destination=../mock/repository.go -package=mock
type Repository interface {
	Save(data models.Metrics) error
	Get(metric *models.Metrics) error
	GetAll() ([]models.Metrics, error)
	Close() error
	UpdateSyncMiddleware(next http.Handler) http.Handler
	UpdateTicker(interval int, done chan bool)
}
