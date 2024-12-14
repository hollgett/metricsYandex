package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
)

type ApiMetric struct {
	handler handlers.MetricHandler
}

func NewApiMetric(h handlers.MetricHandler) *ApiMetric {
	return &ApiMetric{handler: h}
}

func validRequest(typeM, nameM string) (int, bool) {
	if typeM != "counter" && typeM != "gauge" {
		return http.StatusBadRequest, false
	}
	if len(nameM) == 0 {
		return http.StatusNotFound, false
	}
	return 0, true
}

func (a *ApiMetric) ContentTypeMiddleware(expectedType ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contT := r.Header.Get("Content-Type")
			for _, v := range expectedType {
				if contT == v {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "request content type unsupported", http.StatusUnsupportedMediaType)
		})
	}
}
func (a *ApiMetric) UpdateMetricPost(w http.ResponseWriter, r *http.Request) {
	requestParam := []string{chi.URLParam(r, "typeM"), chi.URLParam(r, "nameM"), chi.URLParam(r, "valueM")}
	if statusCode, ok := validRequest(requestParam[0], requestParam[1]); !ok {
		w.WriteHeader(statusCode)
		return
	}
	err := a.handler.CollectingMetric(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *ApiMetric) GetMetric(w http.ResponseWriter, r *http.Request) {
	requestParam := []string{chi.URLParam(r, "typeM"), chi.URLParam(r, "nameM")}
	if statusCode, ok := validRequest(requestParam[0], requestParam[1]); !ok {
		w.WriteHeader(statusCode)
		return
	}
	result, err := a.handler.GetMetric(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func (a *ApiMetric) GetMetricAll(w http.ResponseWriter, r *http.Request) {
	body, err := a.handler.GetMetricAll()
	if err != nil {
		http.Error(w, "error get metric: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
