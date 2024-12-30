package api

import (
	"net/http"

	"github.com/hollgett/metricsYandex.git/internal/logger"
	"go.uber.org/zap"
)

func RespondWithError(w http.ResponseWriter, code int, logMessage string, err string) {
	logger.LogInfo(logMessage, zap.String("error catch", err))
	http.Error(w, err, code)

}

func RespondWithSuccess(w http.ResponseWriter, contentT string, code int, response string) {
	if len(contentT) != 0 {
		w.Header().Set("Content-Type", contentT)
	}
	w.WriteHeader(code)
	w.Write([]byte(response))
}
