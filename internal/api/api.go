package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"go.uber.org/zap"
)

type APIMetric struct {
	handler handlers.MetricHandler
}

func NewAPIMetric(h handlers.MetricHandler) *APIMetric {
	return &APIMetric{handler: h}
}

func (a *APIMetric) UpdateMetricPost(w http.ResponseWriter, r *http.Request) {
	contT := r.Header.Get("Content-Type")
	switch contT {
	case "application/json":
		a.UpdateMetricPostJson(w, r)
	default:
		a.UpdateMetricPlainText(w, r)
	}
}

func (a *APIMetric) UpdateMetricPlainText(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{
		ID:    chi.URLParam(r, "nameM"),
		MType: chi.URLParam(r, "typeM"),
	}
	logger.LogInfo("UpdateMetricPlainText start", zap.Any("value", metrics), zap.String("content type", r.Header.Get("Content-Type")))
	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "UpdateMetricPlainText: ValidateTypeMetric error", err.Error())
		return
	}
	if err := handlers.ValidateNameMetric(metrics.ID); err != nil {
		RespondWithError(w, http.StatusNotFound, "ValidateNameMetric error", err.Error())
		return
	}

	ConvMetricVal(&metrics, chi.URLParam(r, "valueM"), w)

	if err := a.handler.CollectingMetric(&metrics); err != nil {
		logger.LogInfo("UpdateMetricPlainText: CollectingMetric error", zap.Any("arguments", metrics), zap.String("error catch", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("UpdateMetricPlainText complete")
	w.WriteHeader(http.StatusOK)
}

func (a *APIMetric) UpdateMetricPostJson(w http.ResponseWriter, r *http.Request) {
	metrics := models.Metrics{}
	json.NewDecoder(r.Body).Decode(&metrics)
	logger.LogInfo("UpdateMetricPostJson start", zap.Any("value", metrics), zap.String("content type", r.Header.Get("Content-Type")))
	if err := handlers.ValidateTypeMetric(metrics.MType); err != nil {
		RespondWithError(w, http.StatusBadRequest, "UpdateMetricPostJson: ValidateTypeMetric error", err.Error())
		return
	}
	if err := handlers.ValidateNameMetric(metrics.ID); err != nil {
		RespondWithError(w, http.StatusNotFound, "UpdateMetricPostJson: ValidateNameMetric error", err.Error())
		return
	}
	if err := a.handler.CollectingMetric(&metrics); err != nil {
		logger.LogInfo("UpdateMetricPostJson: CollectingMetric error", zap.Any("arguments", metrics), zap.String("error catch", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
	logger.LogInfo("UpdateMetricPostJson complete", zap.Any("response", metrics))
}

func (a *APIMetric) GetMetric(w http.ResponseWriter, r *http.Request) {
	mediaT := r.Header.Get("Content-Type")
	switch mediaT {
	case "application/json":
		a.getMetricJSON(w, r)
	default:
		a.getMetricPlainText(w, r)
	}
}

func (a *APIMetric) getMetricPlainText(w http.ResponseWriter, r *http.Request) {
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
	logger.LogInfo("getMetricPlainText: GetMetric service complete", zap.Any("value", metrics))

	Conv2StrWithResp(&metrics, w)
}

func (a *APIMetric) getMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)

	logger.LogInfo("requestJSON complete", zap.Any("request param take", metrics))
}

func (a *APIMetric) GetMetricAll(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("GetMetricAll start")

	body, err := a.handler.GetMetricAll()
	if err != nil {
		logger.Log.Info(
			"GetMetricAll error",
			zap.String("error catch", err.Error()),
		)
		http.Error(w, "error get metric: "+err.Error(), http.StatusBadRequest)
		return
	}
	RespondWithSuccess(w, "", http.StatusOK, body)
}
