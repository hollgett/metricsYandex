package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/server/api"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
)

func setupRouters(h *api.APIMetric, log logger.Logger) *chi.Mux {
	rtr := chi.NewMux()

	rtr.Use(log.RequestMiddleware,
		log.ResponseMiddleware,
		api.CompressMiddleware,
		api.ContentTypeMiddleware("text/plain", "", "application/json", "application/x-gzip"),
	)
	rtr.Get("/", h.GetMetricAll)
	rtr.Get("/ping", h.Ping)
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

func New(h *api.APIMetric, log logger.Logger, addr string) *http.Server {
	r := setupRouters(h, log)
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
