package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"

	"go.uber.org/zap"
)

type metrics struct {
	PollCount int64
	metrics   map[string]float64
}

func randomValue() float64 {
	return rand.Float64()
}

func (m *metrics) updateMetrics() {
	runtime.ReadMemStats(&memStats)
	m.PollCount++
	m.metrics["RandomValue"] = randomValue()
	m.metrics["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.metrics["Alloc"] = float64(memStats.Alloc)
	m.metrics["Frees"] = float64(memStats.Frees)
	m.metrics["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.metrics["GCSys"] = float64(memStats.GCSys)
	m.metrics["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.metrics["HeapIdle"] = float64(memStats.HeapIdle)
	m.metrics["HeapInuse"] = float64(memStats.HeapInuse)
	m.metrics["HeapObjects"] = float64(memStats.HeapObjects)
	m.metrics["HeapReleased"] = float64(memStats.HeapReleased)
	m.metrics["HeapSys"] = float64(memStats.HeapSys)
	m.metrics["LastGC"] = float64(memStats.LastGC)
	m.metrics["Lookups"] = float64(memStats.Lookups)
	m.metrics["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.metrics["MCacheSys"] = float64(memStats.MCacheSys)
	m.metrics["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.metrics["MSpanSys"] = float64(memStats.MSpanSys)
	m.metrics["Mallocs"] = float64(memStats.Mallocs)
	m.metrics["NextGC"] = float64(memStats.NextGC)
	m.metrics["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.metrics["NumGC"] = float64(memStats.NumGC)
	m.metrics["OtherSys"] = float64(memStats.OtherSys)
	m.metrics["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.metrics["StackInuse"] = float64(memStats.StackInuse)
	m.metrics["StackSys"] = float64(memStats.StackSys)
	m.metrics["Sys"] = float64(memStats.Sys)
	m.metrics["TotalAlloc"] = float64(memStats.TotalAlloc)
}

func (m metrics) sendMetricsJSON() error {
	metricsP := Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &m.PollCount,
	}
	resp, err := clientPost(metricsP)
	if err != nil {
		if resp.Body() != nil {
			return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\t\tresponse data \"%v\"\r\n\t\tresponse body: \"%v\"", err, metricsP, resp.String())
		}
		return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
	}
	log.Info("PollCount", zap.Any("value", metricsP.Delta), zap.String("status", resp.Status()))
	for i, v := range m.metrics {
		metricsG := Metrics{
			ID:    i,
			MType: "gauge",
			Value: &v,
		}
		resp, err := clientPost(metricsG)
		if err != nil {
			if resp.Body() != nil {
				return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\t\tresponse data \"%v\"\r\n\t\tresponse body: \"%v\"", err, metricsP, resp.String())
			}
			return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
		}
		log.Info("request", zap.Any("value", metricsG), zap.String("status", resp.Status()))
	}
	return nil
}
