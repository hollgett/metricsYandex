package main

import (
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/config"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/server"
	"github.com/hollgett/metricsYandex.git/internal/storage"
)

func main() {
	logger.InitLogger()

	cfg := config.InitConfig()

	memStorage := storage.NewMemStorage()
	handlers := handlers.NewMetricHandler(memStorage)
	api := api.NewAPIMetric(handlers)
	server := server.NewServer(api, cfg)
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
