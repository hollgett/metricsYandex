package service

import (
	"math/rand/v2"
	"runtime"

	"github.com/hollgett/metricsYandex.git/internal/agent/models"
)

type Metrics struct {
	runtimeMem runtime.MemStats
	PollCount  int64
	metrics    map[string]float64
}

func NewMemStorage() *Metrics {
	return &Metrics{
		runtimeMem: runtime.MemStats{},
		PollCount:  0,
		metrics:    make(map[string]float64),
	}
}

func randomValue() float64 {
	return rand.Float64()
}

func (m *Metrics) UpdateMetrics() {
	runtime.ReadMemStats(&m.runtimeMem)
	m.PollCount++
	m.metrics["RandomValue"] = randomValue()
	m.metrics["BuckHashSys"] = float64(m.runtimeMem.BuckHashSys)
	m.metrics["Alloc"] = float64(m.runtimeMem.Alloc)
	m.metrics["Frees"] = float64(m.runtimeMem.Frees)
	m.metrics["GCCPUFraction"] = float64(m.runtimeMem.GCCPUFraction)
	m.metrics["GCSys"] = float64(m.runtimeMem.GCSys)
	m.metrics["HeapAlloc"] = float64(m.runtimeMem.HeapAlloc)
	m.metrics["HeapIdle"] = float64(m.runtimeMem.HeapIdle)
	m.metrics["HeapInuse"] = float64(m.runtimeMem.HeapInuse)
	m.metrics["HeapObjects"] = float64(m.runtimeMem.HeapObjects)
	m.metrics["HeapReleased"] = float64(m.runtimeMem.HeapReleased)
	m.metrics["HeapSys"] = float64(m.runtimeMem.HeapSys)
	m.metrics["LastGC"] = float64(m.runtimeMem.LastGC)
	m.metrics["Lookups"] = float64(m.runtimeMem.Lookups)
	m.metrics["MCacheInuse"] = float64(m.runtimeMem.MCacheInuse)
	m.metrics["MCacheSys"] = float64(m.runtimeMem.MCacheSys)
	m.metrics["MSpanInuse"] = float64(m.runtimeMem.MSpanInuse)
	m.metrics["MSpanSys"] = float64(m.runtimeMem.MSpanSys)
	m.metrics["Mallocs"] = float64(m.runtimeMem.Mallocs)
	m.metrics["NextGC"] = float64(m.runtimeMem.NextGC)
	m.metrics["NumForcedGC"] = float64(m.runtimeMem.NumForcedGC)
	m.metrics["NumGC"] = float64(m.runtimeMem.NumGC)
	m.metrics["OtherSys"] = float64(m.runtimeMem.OtherSys)
	m.metrics["PauseTotalNs"] = float64(m.runtimeMem.PauseTotalNs)
	m.metrics["StackInuse"] = float64(m.runtimeMem.StackInuse)
	m.metrics["StackSys"] = float64(m.runtimeMem.StackSys)
	m.metrics["Sys"] = float64(m.runtimeMem.Sys)
	m.metrics["TotalAlloc"] = float64(m.runtimeMem.TotalAlloc)
}

func (m Metrics) GetMetric() []models.Metrics {
	var sliceMetric []models.Metrics
	sliceMetric = append(sliceMetric, models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &m.PollCount,
	})
	for name, value := range m.metrics {
		val := value
		sliceMetric = append(sliceMetric, models.Metrics{

			ID:    name,
			MType: "gauge",
			Value: &val,
		})
	}
	return sliceMetric
}
