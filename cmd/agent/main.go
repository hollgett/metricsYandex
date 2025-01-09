package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	clientOnce *resty.Client
	log        *zap.Logger
	memStats   runtime.MemStats
)

func newClientResty() {
	clientOnce = resty.New().
		SetBaseURL(Cfg.Addr).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetDebug(true)
}

func clientPost(metric Metrics) (*resty.Response, error) {
	data, err := Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}
	return clientOnce.R().
		SetBody(data).
		Post("/update/")
}

func RequestMetric() {
	currentMetrics := metrics{
		PollCount: 0,
		metrics:   make(map[string]float64),
	}

	{
		lastSendTime := time.Now()
		for {
			//update metrics
			time.Sleep(time.Duration(Cfg.PollInterval) * time.Second)
			currentMetrics.updateMetrics()

			//send metrics
			if time.Since(lastSendTime) >= time.Duration(Cfg.ReportInterval)*time.Second {
				if err := currentMetrics.sendMetricsJSON(); err != nil {
					fmt.Println("error catch: ", err)
					return
				}
				lastSendTime = time.Now()
			}
		}
	}
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log = logger
	InitConfig()
	log.Info("agent start", zap.Any("config", Cfg))
	newClientResty()
	go RequestMetric()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	// Ожидание сигнала завершения
	sig := <-exit
	log.Info("signal exit", zap.Any("value", sig))
	fmt.Scan()
}
