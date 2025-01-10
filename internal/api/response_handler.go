package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hollgett/metricsYandex.git/internal/jsonutil"
	"github.com/hollgett/metricsYandex.git/internal/logger"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"go.uber.org/zap"
)

func RespondWithError(w http.ResponseWriter, code int, logMessage string, err string) {
	logger.LogInfo(logMessage, zap.String("error", err))
	http.Error(w, err, code)
}

func RespondWithSuccess(w http.ResponseWriter, contentT string, code int, response models.Metrics) {
	w.Header().Set("Content-Type", contentT)
	w.WriteHeader(code)
	switch contentT {
	case "application/json":
		if err := jsonutil.EncoderJSON(w, response); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "EncoderJson", err.Error())
		}
	default:
		switch response.MType {
		case "gauge":
			fmt.Fprint(w, strconv.FormatFloat(*response.Value, 'G', -1, 64))
		case "counter":
			fmt.Fprint(w, strconv.FormatInt(*response.Delta, 10))
		default:
			RespondWithError(w, http.StatusInternalServerError, "convertation value case default", "")
		}
	}
	logger.LogInfo("response", zap.String("type", contentT), zap.Any("data", response))
}
