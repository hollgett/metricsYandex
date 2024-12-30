package handlers

import (
	"errors"
	"fmt"

	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"github.com/hollgett/metricsYandex.git/internal/storage"
	"go.uber.org/zap"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type metricHandler struct {
	repo storage.Repositories
}

func NewMetricHandler(repo storage.Repositories) MetricHandler {
	return &metricHandler{repo: repo}
}

func ValidateTypeMetric(typeM string) error {
	if typeM != "counter" && typeM != "gauge" {
		return fmt.Errorf("type metric error, got: %s", typeM)
	}
	return nil
}

func ValidateNameMetric(nameM string) error {
	if len(nameM) == 0 {
		return fmt.Errorf("name metric error, got: %s", nameM)
	}
	return nil
}

func (m *metricHandler) CollectingMetric(metrics *models.Metrics) error {
	switch metrics.MType {
	case gauge:
		if err := m.repo.UpdateGauge(metrics.ID, *metrics.Value); err != nil {
			return fmt.Errorf("function UpdateGauge have error: %w", err)
		}
		return nil
	case counter:

		if err := m.repo.AddCounter(metrics.ID, *metrics.Delta); err != nil {
			return fmt.Errorf("function AddCounter have error: %w", err)
		}
		return nil
	default:
		return errors.New("default case, wrong type metric")
	}
}

func (m *metricHandler) GetMetric(metrics *models.Metrics) error {
	switch metrics.MType {
	case gauge:
		val, err := m.repo.GetMetricGauge(metrics.ID)
		if err != nil {
			return fmt.Errorf("get metric have error: %w", err)
		}
		metrics.Value = &val
		return nil
	case counter:
		val, err := m.repo.GetMetricCounter(metrics.ID)
		if err != nil {
			return fmt.Errorf("get metric have error: %w", err)
		}
		metrics.Delta = &val
		return nil
	default:
		return errors.New("case default, get metric error")
	}
}

func (m *metricHandler) GetMetricAll() (string, error) {
	listMetric, err := m.repo.GetMetricAll()
	if err != nil {
		return "", fmt.Errorf("error repository get: %w", err)
	}
	bodyHead := `<html>
    <head>
    <title></title>
    </head>
    <body>
	<table>
		<thead>
		<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>
		</thead>
	`
	var body string
	for i, v := range listMetric {
		logger.Log.Info(
			"GetMetricAll service got",
			zap.String("name", i),
			zap.String("value", v),
		)
		body += fmt.Sprintf(`<tr><td>%v</td><td>%v</td></tr><br>`, i, v)
	}
	bodyBottom := `</table>
			</body>
		</html>`
	return fmt.Sprint(bodyHead, body, bodyBottom), nil
}
