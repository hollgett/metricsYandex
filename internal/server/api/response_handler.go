package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/hollgett/metricsYandex.git/internal/server/utils"
)

func RespondWithError(w http.ResponseWriter, code int, logMessage string, err error) {
	logger.LogErr(logMessage, err)
	http.Error(w, err.Error(), code)
}

func RespondWithSuccess(w http.ResponseWriter, contentT string, code int, response models.Metrics) {
	w.Header().Set("Content-Type", contentT)
	w.WriteHeader(code)
	switch contentT {
	case "application/json":
		if err := utils.EncoderJSON(w, response); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "EncoderJson", err)
		}
	default:
		switch response.MType {
		case "gauge":
			fmt.Fprint(w, strconv.FormatFloat(*response.Value, 'G', -1, 64))
		case "counter":
			fmt.Fprint(w, strconv.FormatInt(*response.Delta, 10))
		default:
			RespondWithError(w, http.StatusInternalServerError, "convertation value case default", nil)
		}
	}
	logger.LogAny("response", "data", response)
}



