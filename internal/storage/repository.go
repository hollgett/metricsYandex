package storage

import(
	"errors"
)
//go:generate mockgen -source=repository.go -destination=../mock/repository.go -package=mock
type Repositories interface {
	UpdateGauge(nameMetric string, val float64) error
	AddCounter(nameMetric string, val int64) error
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() Repositories {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}

func (m *MemStorage) UpdateGauge(nameMetric string, val float64) error {
	if len(nameMetric) == 0 {
		return errors.New("name metric have nil")
	}
	m.gauge[nameMetric] = val
	return nil
}

func (m *MemStorage) AddCounter(nameMetric string, val int64) error {
	if len(nameMetric) == 0 {
		return errors.New("name metric have nil")
	}
	m.counter[nameMetric] += val
	return nil
}

