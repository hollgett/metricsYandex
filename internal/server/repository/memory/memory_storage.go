package memory

import (
	"errors"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

type Memory interface {
	SetGauge(name string, val float64) error
	GetGauge(name string) (float64, error)
	AddCounter(name string, val int64) error
	GetCounter(name string) (int64, error)
	GetAll() []models.Metrics
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemoryStorage() Memory {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) SetGauge(name string, val float64) error {
	if len(name) == 0 {
		return errors.New("name gauge have nil")
	}
	ms.gauge[name] = val
	return nil
}

func (ms MemStorage) GetGauge(name string) (float64, error) {
	if len(name) == 0 {
		return 0, errors.New("name gauge have nil")
	}
	v, ok := ms.gauge[name]
	if !ok {
		return 0, errors.New("data does't exist")
	}
	return v, nil
}

func (ms *MemStorage) AddCounter(name string, val int64) error {
	if len(name) == 0 {
		return errors.New("name counter have nil")
	}
	ms.counter[name] += val
	return nil
}

func (ms MemStorage) GetCounter(name string) (int64, error) {
	if len(name) == 0 {
		return 0, errors.New("name counter have nil")
	}
	v, ok := ms.counter[name]
	if !ok {
		return 0, errors.New("data does't exist")
	}
	return v, nil
}

func (ms MemStorage) GetAll() []models.Metrics {
	var data []models.Metrics
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
	return data
}
