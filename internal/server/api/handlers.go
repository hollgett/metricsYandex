package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/services"
	"github.com/hollgett/metricsYandex.git/internal/server/utils"
)

const (
	jsonT string = "application/json"
	textT string = "text/plain"
)

type APIMetric struct {
	handler services.MetricHandler
}

func NewAPIMetric(h services.MetricHandler) *APIMetric {
	return &APIMetric{handler: h}
}

func (a *APIMetric) UpdateMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	logger.LogAny("UpdateMetricPlainText start", "value", metrics)
	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	ConvMetricVal(&metrics, chi.URLParam(r, "valueM"), w)

	if err := a.handler.CollectingMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "CollectingMetric plain text", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *APIMetric) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{}
	if err := utils.DecoderJSON(r.Body, &metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "decode json", err)
		return
	}
	logger.LogAny("UpdateMetricJSON start", "value", metrics)
	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	if err := a.handler.CollectingMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "CollectingMetric json", err)
		return
	}
	RespondWithSuccess(w, jsonT, http.StatusOK, metrics)
}

func (a *APIMetric) GetMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	logger.LogAny("getMetricPlainText start", "value", metrics)

	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}

	if err := a.handler.GetMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err)
		return
	}
	RespondWithSuccess(w, textT, http.StatusOK, metrics)
}

func (a *APIMetric) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics models.Metrics

	if err := utils.DecoderJSON(r.Body, &metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "getMetricJSON decode error", err)
		return
	}
	logger.LogAny("getMetricJSON start", "request param take", metrics)

	if code, err := a.handler.ValidateMetric(&metrics); err != nil {
		RespondWithError(w, code, "ValidateMetric plain text", err)
		return
	}
	if err := a.handler.GetMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err)
		return
	}
	RespondWithSuccess(w, jsonT, http.StatusOK, metrics)
}

func (a *APIMetric) GetMetricAll(w http.ResponseWriter, r *http.Request) {
	body, err := a.handler.GetMetricAll()
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "GetMetricAll", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
