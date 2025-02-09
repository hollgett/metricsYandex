package main

import (
	"context"
	"log"
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

func RequestMetric(mem *service.Metrics, client *api.Client, ctx context.Context, cancel context.CancelFunc) error {
	lastSendTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			//update metrics
			time.Sleep(time.Duration(config.AgentConfig.PollInterval) * time.Second)
			mem.UpdateMetrics()
			//send metrics
			if time.Since(lastSendTime) >= time.Duration(config.AgentConfig.ReportInterval)*time.Second {
				data := mem.GetMetric()
				if err := client.SendMetricsJSON(ctx, data, 2, 1*time.Second); err != nil {
					logger.Log.Info("send metric", zap.Error(err))
					cancel()
				}
				lastSendTime = time.Now()
			}
		}

	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handlerShutDown(cancel)
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("init logger: %v", err)
		cancel()
	}

	if err := config.InitConfig(); err != nil {
		log.Fatalf("init config: %v", err)
		cancel()
	}
	logger.Log.Info("agent config", zap.Any("value", config.AgentConfig))

	memS := service.NewMemStorage()
	client := api.NewClientResty("Content-Encoding", "gzip", false)
	go func() {
		if err := RequestMetric(memS, client, ctx, cancel); err != nil {
			logger.Log.Info("request metric", zap.Error(err))
			cancel()
		}
	}()
	<-ctx.Done()
}

func handlerShutDown(cancel context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	sig := <-exit
	logger.Log.Info("close app", zap.Any("signal", sig))
	cancel()
}
