package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/config"
	"github.com/hollgett/metricsYandex.git/internal/logger"
)

func setupRouters(h *api.ApiMetric) *chi.Mux {
	rtr := chi.NewMux()
	rtr.Use(logger.RequestMiddleware)
	rtr.Use(logger.ResponseMiddleware)
	rtr.Use(api.ContentTypeMiddleware("text/plain", "", "application/json"))
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

func NewServer(h *api.ApiMetric) *http.Server {
	r := setupRouters(h)
	return &http.Server{
		Addr:    config.Cfg.Addr,
		Handler: r,
	}
}
