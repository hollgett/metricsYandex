package services

import (
	"context"
	"fmt"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

var (
	delta int64   = 0
	value float64 = 0
)

type metricHandler struct {
	repo repository.Repository
}

func New(repo repository.Repository) MetricHandler {
	return &metricHandler{repo: repo}
}

func (m *metricHandler) ValidateMetric(metric *models.Metrics) (int, error) {
	if len(metric.ID) == 0 {
		return 404, fmt.Errorf("name metric got nill")
	}
	if metric.MType != counter && metric.MType != gauge {
		return 400, fmt.Errorf("type metric got nill")
	}
	return 0, nil
}

func (m *metricHandler) CollectingMetric(metrics *models.Metrics) error {
	switch metrics.MType {
	case gauge:
		metrics.Delta = &delta
	case counter:
		metrics.Value = &value
	}
	return m.repo.Save(*metrics)
}

func (m *metricHandler) GetMetric(metric *models.Metrics) error {
	switch metric.MType {
	case gauge:
		metric.Delta = &delta
	case counter:
		metric.Value = &value
	}
	if err := m.repo.Get(metric); err != nil {
		return fmt.Errorf("get metric have error: %w", err)
	}
	return nil
}

func (m *metricHandler) GetMetricAll() (string, error) {
	listMetric, err := m.repo.GetAll()
	if err != nil {
		return "", fmt.Errorf("error repository get: %w", err)
	}
	bodyHead := `
		<table>
		<thead>
		<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>
		</thead>
	`
	var body string
	for _, v := range listMetric {
		switch v.MType {
		case gauge:
			body += fmt.Sprintf(`<tr><td>%v</td><td>%v</td></tr><br>`, v.ID, *v.Value)
		case counter:
			body += fmt.Sprintf(`<tr><td>%v</td><td>%v</td></tr><br>`, v.ID, *v.Delta)
		}

	}
	bodyBottom := `</table>`
	return fmt.Sprint(bodyHead, body, bodyBottom), nil
}

func (m *metricHandler) PingDB(ctx context.Context) error {
	return m.repo.Ping(ctx)
}

func (m *metricHandler) Batch(ctx context.Context, metrics []models.Metrics) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for i, _ := range metrics {
		switch metrics[i].MType {
		case gauge:
			metrics[i].Delta = &delta
		case counter:
			metrics[i].Value = &value
		}
	}
	if err := m.repo.Batch(ctx, metrics); err != nil {
		return fmt.Errorf("batch repo: %w", err)
	}
	return nil
}
