package api

import (
	"net/http"
	"strings"

	"github.com/hollgett/metricsYandex.git/internal/handlers"
)

type ApiMetric struct {
	handler handlers.MetricHandler
}

func NewApiMetric(h handlers.MetricHandler) *ApiMetric {
	return &ApiMetric{handler: h}
}

func (a *ApiMetric) UpdateMetricPost(w http.ResponseWriter, r *http.Request) {
	err := a.handler.CollectingMetric(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *ApiMetric) CheckURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// fmt.Print(r.URL.Path,": ")
		arrURL := strings.Split(r.URL.Path[1:], "/")
		methodM, typeM, nameM := arrURL[0], arrURL[1], arrURL[2]
		// fmt.Printf("arr : %#v\n", arrURL)
		if methodM != "update" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if typeM != "counter" && typeM != "gauge" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(nameM) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.StripPrefix(`/update/`, next).ServeHTTP(w, r)
	})
}
