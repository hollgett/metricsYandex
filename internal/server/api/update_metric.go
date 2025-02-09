package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

func (a *APIMetric) UpdateMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	a.log.LogAny("UpdateMetricPlainText", "value", metrics)
	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		a.RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	value := chi.URLParam(r, "valueM")
	switch metrics.MType {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			a.RespondWithError(w, http.StatusBadRequest, "gauge parse", err)
		}
		metrics.Value = &val
	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			a.RespondWithError(w, http.StatusBadRequest, "counter parse", err)
		}
		metrics.Delta = &val
	}

	if err := a.handler.CollectingMetric(&metrics); err != nil {
		a.RespondWithError(w, http.StatusBadRequest, "CollectingMetric plain text", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *APIMetric) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{}
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		a.RespondWithError(w, http.StatusBadRequest, "decode json", err)
		return
	}
	a.log.LogAny("UpdateMetricJSON start", "value", metrics)
	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		a.RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	if err := a.handler.CollectingMetric(&metrics); err != nil {
		a.RespondWithError(w, http.StatusBadRequest, "CollectingMetric json", err)
		return
	}
	a.RespondWithSuccessJson(w, http.StatusOK, metrics)
}
