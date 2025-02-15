package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
)

var (
	ErrMetric = errors.New("unknown metrics")
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func New() repository.Repository {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) Get(metric *models.Metrics) error {
	switch metric.MType {
	case "gauge":
		val, err := ms.getGauge(metric.ID)
		if err != nil {
			return fmt.Errorf("get gauge err: %w", err)
		}
		metric.Value = &val
	case "counter":
		val, err := ms.getCounter(metric.ID)
		if err != nil {
			return fmt.Errorf("get counter err: %w", err)
		}
		metric.Delta = &val
	default:
		return ErrMetric
	}
	return nil
}

func (ms MemStorage) GetAll() ([]models.Metrics, error) {
	var data []models.Metrics
	if len(ms.gauge) == 0 || len(ms.counter) == 0 {
		return nil, errors.New("data storage empty")
	}
	if len(ms.gauge) != 0 {
		for name, value := range ms.gauge {
			val := value
			data = append(data, models.Metrics{
				ID:    name,
				MType: "gauge",
				Value: &val,
			})
		}
	}
	if len(ms.counter) != 0 {
		for name, value := range ms.counter {
			val := value
			data = append(data, models.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &val,
			})
		}
	}
	return data, nil
}

func (ms *MemStorage) Save(data models.Metrics) error {
	switch data.MType {
	case "gauge":
		if err := ms.setGauge(data.ID, *data.Value); err != nil {
			return fmt.Errorf("set gauge err: %w", err)
		}
	case "counter":
		if err := ms.addCounter(data.ID, *data.Delta); err != nil {
			return fmt.Errorf("add counter err: %w", err)
		}
	default:
		return ErrMetric
	}
	return nil
}

func (ms *MemStorage) setGauge(name string, val float64) error {
	if len(name) == 0 {
		return errors.New("name gauge have nil")
	}
	ms.gauge[name] = val
	return nil
}

func (ms MemStorage) getGauge(name string) (float64, error) {
	if len(name) == 0 {
		return 0, errors.New("name gauge have nil")
	}
	v, ok := ms.gauge[name]
	if !ok {
		return 0, errors.New("data does't exist")
	}
	return v, nil
}

func (ms *MemStorage) addCounter(name string, val int64) error {
	if len(name) == 0 {
		return errors.New("name counter have nil")
	}
	ms.counter[name] += val
	return nil
}

func (ms MemStorage) getCounter(name string) (int64, error) {
	if len(name) == 0 {
		return 0, errors.New("name counter have nil")
	}
	v, ok := ms.counter[name]
	if !ok {
		return 0, errors.New("data does't exist")
	}
	return v, nil
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (ms *MemStorage) Batch(ctx context.Context, metrics []models.Metrics) error {
	if err := ctx.Err(); err != nil {
		return ctx.Err()
	}
	for _, v := range metrics {
		switch v.MType {
		case "gauge":
			ms.setGauge(v.ID, *v.Value)
		case "counter":
			ms.addCounter(v.ID, *v.Delta)
		default:
			return ErrMetric
		}
	}
	return nil
}

func (ms *MemStorage) Close() error {
	return nil
}
