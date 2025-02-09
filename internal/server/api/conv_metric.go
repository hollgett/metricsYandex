package api

import (
	"net/http"
	"strconv"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

func ConvMetricVal(metrics *models.Metrics, valueM string, w http.ResponseWriter) {
	switch metrics.MType {
	case "gauge":
		val, err := strconv.ParseFloat(valueM, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "gauge parse", err)
		}
		metrics.Value = &val
	case "counter":
		val, err := strconv.ParseInt(valueM, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "counter parse", err)
		}
		metrics.Delta = &val
	}
}
