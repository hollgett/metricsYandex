package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type gauge float64
type counter int64

const (
	pollInterval   time.Duration = 2 * time.Second
	reportInterval time.Duration = 10 * time.Second
)

type metrics struct {
	PollCount counter
	metrics   map[string]gauge
}

func randomValue() gauge {
	return gauge(rand.Intn(1000))
}

func (m *metrics) updateMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.PollCount++
	m.metrics["RandomValue"] = randomValue()
	m.metrics["BuckHashSys"] = gauge(memStats.BuckHashSys)
	m.metrics["Alloc"] = gauge(memStats.Alloc)
	m.metrics["Frees"] = gauge(memStats.Frees)
	m.metrics["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
	m.metrics["GCSys"] = gauge(memStats.GCSys)
	m.metrics["HeapAlloc"] = gauge(memStats.HeapAlloc)
	m.metrics["HeapIdle"] = gauge(memStats.HeapIdle)
	m.metrics["HeapInuse"] = gauge(memStats.HeapInuse)
	m.metrics["HeapObjects"] = gauge(memStats.HeapObjects)
	m.metrics["HeapReleased"] = gauge(memStats.HeapReleased)
	m.metrics["HeapSys"] = gauge(memStats.HeapSys)
	m.metrics["LastGC"] = gauge(memStats.LastGC)
	m.metrics["Lookups"] = gauge(memStats.Lookups)
	m.metrics["MCacheInuse"] = gauge(memStats.MCacheInuse)
	m.metrics["MCacheSys"] = gauge(memStats.MCacheSys)
	m.metrics["MSpanInuse"] = gauge(memStats.MSpanInuse)
	m.metrics["MSpanSys"] = gauge(memStats.MSpanSys)
	m.metrics["Mallocs"] = gauge(memStats.Mallocs)
	m.metrics["NextGC"] = gauge(memStats.NextGC)
	m.metrics["NumForcedGC"] = gauge(memStats.NumForcedGC)
	m.metrics["NumGC"] = gauge(memStats.NumGC)
	m.metrics["OtherSys"] = gauge(memStats.OtherSys)
	m.metrics["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
	m.metrics["StackInuse"] = gauge(memStats.StackInuse)
	m.metrics["StackSys"] = gauge(memStats.StackSys)
	m.metrics["Sys"] = gauge(memStats.Sys)
	m.metrics["TotalAlloc"] = gauge(memStats.TotalAlloc)
}

func (m metrics) sendMetrics() error {
	fmt.Printf("send metrics PollCounter: %v\n", m.PollCount)
	http.Post(fmt.Sprintf("http://localhost:8080/update/counter/PollCount/%v", m.PollCount), "text/plain", nil)
	for i, v := range m.metrics {
		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%s/%v", i, v), "text/plain", nil)
		if err != nil {
			return errors.New(fmt.Sprint("Error with request to server", err.Error()))
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		fmt.Println(resp.Status)

	}
	return nil
}

func main() {
	currentMetrics := metrics{
		PollCount: 0,
		metrics:   make(map[string]gauge),
	}
	lastSendTime := time.Now()
	for {
		//update metrics
		time.Sleep(pollInterval)
		currentMetrics.updateMetrics()
		// fmt.Printf("time: %v, update metrics\n", time.Since(lastSendTime))
		// fmt.Printf("%+v\n", currentMetrics)
		//send metrics
		if time.Since(lastSendTime) >= reportInterval {
			err := currentMetrics.sendMetrics()
			if err != nil {
				fmt.Println(err)
				return
			}
			lastSendTime = time.Now()
		}
	}
}
