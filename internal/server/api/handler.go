package api

import (
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/services"
)

type APIMetric struct {
	handler services.MetricHandler
	log     logger.Logger
}

func New(h services.MetricHandler, log logger.Logger) *APIMetric {
	return &APIMetric{handler: h, log: log}
}
