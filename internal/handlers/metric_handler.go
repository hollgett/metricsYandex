package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/hollgett/metricsYandex.git/internal/storage"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

//go:generate mockgen -source=metric_handler.go -destination=../mock/metric_handler.go -package=mock
type MetricHandler interface {
	CollectingMetric(originalURL string) error
}

type metricHandler struct {
	repo storage.Repositories
}

func NewMetricHandler(repo storage.Repositories) MetricHandler {
	return &metricHandler{repo: repo}
}

func (m *metricHandler) CollectingMetric(originalURL string) error {
	// fmt.Print("CollectingMetric: ",originalURL," ")
	arrURL := strings.Split(originalURL, "/")
	// fmt.Printf("arr : %#v\n", arrURL)
	if len(arrURL) != 3 || len(arrURL[1]) == 0 {
		return errors.New("wrong request")
	}
	typeM, nameM, valueM := arrURL[0], arrURL[1], arrURL[2]
	switch typeM {
	case gauge:
		val, err := strconv.ParseFloat(valueM, 64)
		if err != nil {
			return errors.New("wrong value")
		}
		if err := m.repo.UpdateGauge(nameM, val); err != nil {
			return errors.Join(errors.New("function UpdateGauge have error: "), err)
		}
	case counter:
		val, err := strconv.ParseInt(valueM, 10, 64)
		if err != nil {
			return errors.New("wrong value")
		}
		if err := m.repo.AddCounter(nameM, val); err != nil {
			return errors.Join(errors.New("function AddCounter have error: "), err)
		}
	default:
		return errors.New("wrong type metric")
	}
	return nil
}
