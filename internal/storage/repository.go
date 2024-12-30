package storage

import (
	"errors"
	"fmt"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() Repositories {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}

func (m *MemStorage) UpdateGauge(nameMetric string, val float64) error {
	m.gauge[nameMetric] = val
	return nil
}

func (m *MemStorage) AddCounter(nameMetric string, val int64) error {
	m.counter[nameMetric] += val
	return nil
}

func (m *MemStorage) GetMetricGauge(nameM string) (float64, error) {
	if v, ok := m.gauge[nameM]; ok {
		return v, nil
	}
	return 0, errors.New("data doesn't exist")
}

func (m *MemStorage) GetMetricCounter(nameM string) (int64, error) {
	if v, ok := m.counter[nameM]; ok {
		return v, nil
	}
	return 0, errors.New("data doesn't exist")
}

func (m *MemStorage) GetMetricAll() (map[string]string, error) {
	list := make(map[string]string, len(m.counter)+len(m.gauge))
	for i, v := range m.gauge {
		list[i] = fmt.Sprintf("%v", v)
	}
	for i, v := range m.counter {
		list[i] = fmt.Sprintf("%v", v)
	}
	if len(list) == 0 {
		return nil, errors.New("data doesn't exist")
	}
	return list, nil
}
