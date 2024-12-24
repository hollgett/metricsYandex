package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/api"
	"github.com/hollgett/metricsYandex.git/internal/config"
)

func NewServer(h *api.ApiMetric) *http.Server {
	cfg := config.InitConfig()
	rtr := chi.NewMux()
	rtr.Route("/", func(r chi.Router) {
		r.Get("/", h.GetMetricAll)
		r.Route("/value", func(r chi.Router) {
			r.Get("/{typeM}/{nameM}", h.GetMetric)
		})
		r.Route("/update", func(r chi.Router) {
			r.Use(h.ContentTypeMiddleware("text/plain", "text/plain; charset=utf-8", ""))
			r.Post("/{typeM}/{nameM}/{valueM}", h.UpdateMetricPost)
		})
	})
	return &http.Server{
		Addr:    cfg.Addr,
		Handler: rtr,
	}
}
