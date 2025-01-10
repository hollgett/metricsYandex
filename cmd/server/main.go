package main

import (
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/config"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/server"
	"github.com/hollgett/metricsYandex.git/internal/storage"
	"go.uber.org/zap"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}

	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	logger.LogInfo("server start", zap.Any("cfg", config.Cfg))
	memStorage := storage.NewMemStorage()
	handlers := handlers.NewMetricHandler(memStorage)
	api := api.NewAPIMetric(handlers)
	server := server.NewServer(api)
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
