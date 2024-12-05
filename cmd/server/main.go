package main

import (
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/server"
	"github.com/hollgett/metricsYandex.git/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	handlers := handlers.NewMetricHandler(memStorage)
	api := api.NewApiMetric(handlers)
	server := server.NewServer(api)
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
