package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/config"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"go.uber.org/zap"
)

func setupRouters(h *api.APIMetric) *chi.Mux {
	rtr := chi.NewMux()
	rtr.Use(logger.RequestMiddleware)
	rtr.Use(logger.ResponseMiddleware)
	rtr.Use(api.ContentTypeMiddleware("text/plain", "", "application/json"))
	rtr.Get("/", h.GetMetricAll)
	rtr.Route("/value", func(r chi.Router) {
		r.Get("/{typeM}/{nameM}", h.GetMetric)
		r.Post("/", h.GetMetric)
	})
	rtr.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricPost)
		r.Post("/{typeM}/{nameM}/{valueM}", h.UpdateMetricPost)
	})

	return rtr
}

func NewServer(h *api.APIMetric, cfg *config.CommandAddr) *http.Server {
	logger.Log.Info("Server start",
		zap.String("address", cfg.Addr))
	r := setupRouters(h)
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}
}
