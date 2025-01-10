package api

import (
	"net/http"

	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/models"
)

func ConvMetricVal(metrics *models.Metrics, valueM string, w http.ResponseWriter) {
	switch metrics.MType {
	case "gauge":
		val, err := handlers.GaugeParse(valueM)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "gauge parse", err.Error())
		}
		metrics.Value = &val
	case "counter":
		val, err := handlers.CounterParse(valueM)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "counter parse", err.Error())
		}
		metrics.Delta = &val
	}
}
