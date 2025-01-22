package services

import (
	"fmt"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type metricHandler struct {
	repo repository.Repository
}

func NewMetricHandler(repo repository.Repository) MetricHandler {
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
	return m.repo.Save(*metrics)
}

func (m *metricHandler) GetMetric(metrics *models.Metrics) error {
	if err := m.repo.Get(metrics); err != nil {
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
