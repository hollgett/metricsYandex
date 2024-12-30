package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	cfg            *AgentArgs
	pollInterval   time.Duration
	reportInterval time.Duration
)

var clientOnce *resty.Client

func newClientResty() {
	clientOnce = resty.New().
		SetBaseURL("http://"+cfg.Addr).
		SetHeader("Content-Type", "application/json")

}

type metrics struct {
	PollCount int64
	metrics   map[string]float64
}

func randomValue() float64 {
	return rand.Float64()
}

func (m *metrics) updateMetrics() {
	var memStats runtime.MemStats
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
	fmt.Printf("send metrics PollCounter: %v\n", m.PollCount)
	resp, err := clientOnce.R().
		SetBody(Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &m.PollCount,
		}).
		Post("/update/")
	if err != nil {
		fmt.Println("error catch: ", err.Error())
		if resp.Body() != nil {
			return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\tresponse body: %v", err, resp.String())
		}
		return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
	}

	fmt.Println(resp.StatusCode())
	for i, v := range m.metrics {
		resp, err := clientOnce.R().
			SetBody(Metrics{
				ID:    i,
				MType: "gauge",
				Value: &v,
			}).
			Post("/update/")

		if err != nil {
			if resp.Body() != nil {
				return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\tresponse body: %v", err, resp.String())
			}
			return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
		}
		fmt.Println(resp.StatusCode())
	}
	return nil
}

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	cfg = InitConfig()
	log.Info("agent start", zap.Any("config", cfg))

	newClientResty()
	pollInterval = time.Duration(cfg.PollInterval) * time.Second
	reportInterval = time.Duration(cfg.ReportInterval) * time.Second
	currentMetrics := metrics{
		PollCount: 0,
		metrics:   make(map[string]float64),
	}
	lastSendTime := time.Now()
	for {
		//update metrics
		time.Sleep(pollInterval)
		currentMetrics.updateMetrics()

		//send metrics
		if time.Since(lastSendTime) >= reportInterval {
			err := currentMetrics.sendMetricsJSON()
			if err != nil {
				fmt.Println("error catch: ", err)
				return
			}
			lastSendTime = time.Now()
		}
	}
}

// func (m metrics) sendMetricsPlainText() error {
// 	clientOnce = newClientResty()
// 	fmt.Printf("send metrics PollCounter: %v\n", m.PollCount)
// 	resp, err := clientOnce.R().Post(fmt.Sprintf("/update/counter/PollCount/%v", m.PollCount))
// 	if err != nil {
// 		if resp.Body() != nil {
// 			return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\tresponse body: %v", err, resp.String())
// 		}
// 		return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
// 	}
// 	fmt.Println(resp.StatusCode())
// 	for i, v := range m.metrics {
// 		resp, err := clientOnce.R().Post(fmt.Sprintf("/update/gauge/%s/%v", i, v))
// 		if err != nil {
// 			if resp.Body() != nil {
// 				return fmt.Errorf("error with request to server. \r\n\terror:%w\r\n\tresponse body: %v", err, resp.String())
// 			}
// 			return fmt.Errorf("error with request to server. \r\n\terror:%w", err)
// 		}
// 		fmt.Println(resp.StatusCode())
// 	}
// 	return nil
// }
