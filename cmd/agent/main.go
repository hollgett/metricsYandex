package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hollgett/metricsYandex.git/internal/agent/api"
	"github.com/hollgett/metricsYandex.git/internal/agent/config"
	"github.com/hollgett/metricsYandex.git/internal/agent/logger"
	"github.com/hollgett/metricsYandex.git/internal/agent/service"
	"go.uber.org/zap"
)

func RequestMetric(mem *service.Metrics, client *api.Client, done chan bool) error {
	lastSendTime := time.Now()
	for {
		select {
		case <-done:
			return nil
		default:
			//update metrics
			time.Sleep(time.Duration(config.AgentConfig.PollInterval) * time.Second)
			mem.UpdateMetrics()
			//send metrics
			if time.Since(lastSendTime) >= time.Duration(config.AgentConfig.ReportInterval)*time.Second {
				data := mem.GetMetric()
				if err := client.SendMetricsJSON(data); err != nil {
					done <- true
					return err
				}
				lastSendTime = time.Now()
			}
		}

	}
}

func main() {
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}

	if err := config.InitConfig(); err != nil {
		logger.Log.Info("config init", zap.Error(err))
		panic(err)
	}
	logger.Log.Info("config start", zap.Any("value", config.AgentConfig))

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)
	go func() {
		sig := <-exit
		logger.Log.Info("close app", zap.Any("signal", sig))
		done <- true
	}()

	memS := service.NewMemStorage()
	client := api.NewClientResty("Content-Encoding", "gzip", false)
	go func() {
		if err := RequestMetric(memS, client, done); err != nil {
			logger.Log.Info("request metric", zap.Error(err))
			done <- true
		}
	}()
	<-done
	close(done)
	fmt.Scanln()
}
