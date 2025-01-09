package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/jsonutil"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"go.uber.org/zap"
)

const (
	jsonT string = "application/json"
	textT string = "text/plain"
)

type ApiMetric struct {
	handler handlers.MetricHandler
}

func NewApiMetric(h handlers.MetricHandler) *ApiMetric {
	return &ApiMetric{handler: h}
}

func (a *ApiMetric) UpdateMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	logger.LogInfo("UpdateMetricPlainText start", zap.Any("value", metrics), zap.String("content type", r.Header.Get("Content-Type")))
	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "ValidateTypeMetric plain text", err.Error())
		return
	}
	if err := handlers.ValidateNameMetric(metrics.ID); err != nil {
		RespondWithError(w, http.StatusNotFound, "ValidateNameMetric plain text", err.Error())
		return
	}
	ConvMetricVal(&metrics, chi.URLParam(r, "valueM"), w)

	if err := a.handler.CollectingMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "CollectingMetric plain text", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *ApiMetric) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{}
	if err := jsonutil.DecoderJSON(r.Body, &metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "decode json", err.Error())
		return
	}
	logger.LogInfo("UpdateMetricPostJson start", zap.Any("value", metrics), zap.String("content type", r.Header.Get("Content-Type")))
	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "ValidateTypeMetric json", err.Error())
		return
	}
	if err := handlers.ValidateNameMetric(metrics.ID); err != nil {
		RespondWithError(w, http.StatusNotFound, "ValidateNameMetric json", err.Error())
		return
	}
	if err := a.handler.CollectingMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "CollectingMetric json", err.Error())
		return
	}
	RespondWithSuccess(w, jsonT, http.StatusOK, metrics)
}

func (a *ApiMetric) GetMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	logger.LogInfo("getMetricPlainText start", zap.Any("request param take", metrics))

	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "getMetricPlainText: ValidateTypeMetric error", err.Error())
		return
	}

	if err := a.handler.GetMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err.Error())
		return
	}
	RespondWithSuccess(w, textT, http.StatusOK, metrics)
}

func (a *ApiMetric) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics models.Metrics

	if err := jsonutil.DecoderJSON(r.Body, &metrics); err != nil {
		RespondWithError(w, http.StatusBadRequest, "getMetricJSON decode error", err.Error())
		return
	}
	logger.LogInfo("getMetricJSON start", zap.Any("request param take", metrics), zap.String("content type", r.Header.Get("Content-Type")))

	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "getMetricJSON: ValidateTypeMetric error", err.Error())
		return
	}

	if err := a.handler.GetMetric(&metrics); err != nil {
		RespondWithError(w, http.StatusNotFound, "getMetricJSON: GetMetric service error", err.Error())
		return
	}
	RespondWithSuccess(w, jsonT, http.StatusOK, metrics)
}

func (a *ApiMetric) GetMetricAll(w http.ResponseWriter, r *http.Request) {
	body, err := a.handler.GetMetricAll()
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "GetMetricAll", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
