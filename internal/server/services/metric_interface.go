package services

import (
	"context"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

//go:generate mockgen -source=metric_interface.go -destination=../mock/metric_handler.go -package=mock
type MetricHandler interface {
	CollectingMetric(metric *models.Metrics) error
	GetMetric(metric *models.Metrics) error
	GetMetricAll() (string, error)
	ValidateMetric(metric *models.Metrics) (int, error)
	Batch(ctx context.Context, metrics []models.Metrics) error
	PingDB(ctx context.Context) error
}
