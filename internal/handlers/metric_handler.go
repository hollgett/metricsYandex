package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hollgett/metricsYandex.git/internal/storage"
)

type metricHandler struct {
	repo storage.Repositories
}

func NewMetricHandler(repo storage.Repositories) MetricHandler {
	return &metricHandler{repo: repo}
}

const (
	gauge   string = "gauge"
	counter string = "counter"
)

func (m *metricHandler) CollectingMetric(requestParam []string) error {
	typeM, nameM, valueM := requestParam[0], requestParam[1], requestParam[2]
	switch typeM {
	case gauge:
		val, err := strconv.ParseFloat(valueM, 64)
		if err != nil {
			return errors.New("wrong value")
		}
		if err := m.repo.UpdateGauge(nameM, val); err != nil {
			return fmt.Errorf("function UpdateGauge have error: %w", err)
		}
	case counter:
		val, err := strconv.ParseInt(valueM, 10, 64)
		if err != nil {
			return errors.New("wrong value")
		}
		if err := m.repo.AddCounter(nameM, val); err != nil {
			return fmt.Errorf("function AddCounter have error: %w", err)
		}
	default:
		return errors.New("wrong type metric")
	}
	return nil
}

func (m *metricHandler) GetMetric(requestParam []string) (string, error) {
	typeM, nameM := requestParam[0], requestParam[1]
	fmt.Println(typeM, nameM)
	switch typeM {
	case gauge:
		val, err := m.repo.GetMetricGauge(nameM)
		if err != nil {
			return "", fmt.Errorf("get metric have error: %w", err)
		}
		return strconv.FormatFloat(val, 'G', -1, 64), nil
	case counter:
		val, err := m.repo.GetMetricCounter(nameM)
		if err != nil {
			return "", fmt.Errorf("get metric have error: %w", err)
		}
		return strconv.FormatInt(val, 10), nil
	}
	return "", errors.New("get metric error")
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
	bodyBottom := `</table>
			</body>
		</html>`
	var body string
	for i, v := range listMetric {
		body += fmt.Sprintf(`<tr><td>%v</td><td>%v</td></tr>`+"\r", i, v)
	}
	return fmt.Sprint(bodyHead, body, bodyBottom), nil
}
