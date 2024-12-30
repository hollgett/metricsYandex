package api

import (
	"net/http"
	"strconv"

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

func Conv2StrWithResp(metrics *models.Metrics, w http.ResponseWriter) {
	switch metrics.MType {
	case "gauge":
		RespondWithSuccess(w, "text/plain", http.StatusOK, strconv.FormatFloat(*metrics.Value, 'G', -1, 64))
	case "counter":
		RespondWithSuccess(w, "text/plain", http.StatusOK, strconv.FormatInt(*metrics.Delta, 10))
	default:
		RespondWithError(w, http.StatusInternalServerError, "convertation value case default", "")
	}
}
