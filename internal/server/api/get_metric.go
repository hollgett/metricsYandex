package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

func (a *APIMetric) GetMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	a.log.LogAny("getMetricPlainText start", "value", metrics)

	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		a.RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}

	if err := a.handler.GetMetric(&metrics); err != nil {
		a.RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err)
		return
	}
	a.RespondWithSuccessText(w, http.StatusOK, metrics)
}

func (a *APIMetric) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		a.RespondWithError(w, http.StatusBadRequest, "getMetricJSON decode error", err)
		return
	}
	a.log.LogAny("getMetricJSON start", "request param take", metrics)

	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		a.RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	if err := a.handler.GetMetric(&metrics); err != nil {
		a.RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err)
		return
	}
	a.RespondWithSuccessJson(w, http.StatusOK, metrics)
}

func (a *APIMetric) GetMetricAll(w http.ResponseWriter, r *http.Request) {
	body, err := a.handler.GetMetricAll()
	if err != nil {
		a.RespondWithError(w, http.StatusBadRequest, "GetMetricAll", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
