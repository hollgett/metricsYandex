package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hollgett/metricsYandex.git/internal/server/api"
	"github.com/hollgett/metricsYandex.git/internal/server/config"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
	"github.com/hollgett/metricsYandex.git/internal/server/server"
	"github.com/hollgett/metricsYandex.git/internal/server/services"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)
	go func() {
		sig := <-signalChan
		logger.LogAny("close app", "signal", sig)
		done <- true
	}()
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}

	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	logger.LogAny("server start", "config", config.Config)

	storage, err := repository.NewCompositeStorage(config.Config.PathFileStorage, config.Config.Restore)
	if err != nil {
		panic(err)
	}

	if config.Config.StorageInterval > 0 {
		go storage.UpdateTicker(config.Config.StorageInterval, done)
	}

	handlers := services.NewMetricHandler(storage)
	api := api.NewAPIMetric(handlers)
	server := server.NewServer(api, storage)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.LogErr("server", err)
			done <- true
		}
	}()

	<-done
	close(done)
	if err := storage.Close(); err != nil {
		logger.LogErr("storage close", err)
	}
}
