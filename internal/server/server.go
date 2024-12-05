package server

import (
	"net/http"
	"github.com/hollgett/metricsYandex.git/internal/api"
)

func NewServer (h *api.ApiMetric) *http.Server {
	rtr := http.NewServeMux()
	rtr.Handle(`/`, h.CheckURLMiddleware(http.HandlerFunc(h.UpdateMetricPost)))

	return &http.Server{
		Addr: `:8080`,
		Handler: rtr,
	}
}