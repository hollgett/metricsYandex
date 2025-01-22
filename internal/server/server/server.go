package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/server/api"
	"github.com/hollgett/metricsYandex.git/internal/server/config"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
)

func setupRouters(h *api.APIMetric, repo repository.Repository) *chi.Mux {
	rtr := chi.NewMux()

	rtr.Use(logger.RequestMiddleware,
		logger.ResponseMiddleware,
		api.CompressMiddleware,
		api.ContentTypeMiddleware("text/plain", "", "application/json", "application/x-gzip"),
	)
	if config.Config.StorageInterval == 0 {
		rtr.Use(repo.UpdateSyncMiddleware)
	}

	rtr.Get("/", h.GetMetricAll)
	rtr.Route("/value", func(r chi.Router) {
		r.Get("/{typeM}/{nameM}", h.GetMetricPlainText)
		r.Post("/", h.GetMetricJSON)
	})
	rtr.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricJSON)
		r.Post("/{typeM}/{nameM}/{valueM}", h.UpdateMetricPlainText)
	})

	return rtr
}

func NewServer(h *api.APIMetric, repo repository.Repository) *http.Server {
	r := setupRouters(h, repo)
	return &http.Server{
		Addr:    config.Config.Addr,
		Handler: r,
	}
}
