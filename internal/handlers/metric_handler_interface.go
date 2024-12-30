package handlers

import "github.com/hollgett/metricsYandex.git/internal/models"

//go:generate mockgen -source=metric_handler_interface.go -destination=../mock/metric_handler.go -package=mock
type MetricHandler interface {
	CollectingMetric(metrics *models.Metrics) error
	GetMetric(metrics *models.Metrics) error
	GetMetricAll() (string, error)
}
