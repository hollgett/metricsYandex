package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hollgett/metricsYandex.git/internal/server/models"
)

const (
	counter string = "counter"
	gauge   string = "gauge"
	jsonT   string = "application/json"
	textT   string = "text/plain"
)

func (a *APIMetric) RespondWithError(w http.ResponseWriter, code int, logMessage string, err error) {
	a.log.LogErr(logMessage, err)
	http.Error(w, err.Error(), code)
}

func (a *APIMetric) RespondWithSuccessJson(w http.ResponseWriter, code int, response models.Metrics) {
	w.Header().Set("Content-Type", jsonT)

	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(response); err != nil {
		a.RespondWithError(w, http.StatusInternalServerError, "EncoderJson", err)
		return
	}
	w.WriteHeader(code)
	buffer.WriteTo(w)
	a.log.LogAny("response", "data", response)
}

func (a *APIMetric) RespondWithSuccessText(w http.ResponseWriter, code int, response models.Metrics) {
	w.Header().Set("Content-Type", textT)
	var buffer bytes.Buffer
	switch response.MType {
	case gauge:
		fmt.Fprint(&buffer, strconv.FormatFloat(*response.Value, 'G', -1, 64))
	case counter:
		fmt.Fprint(&buffer, strconv.FormatInt(*response.Delta, 10))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "unknown metric type")
		return
	}
	w.WriteHeader(code)
	buffer.WriteTo(w)
	a.log.LogAny("response", "data", response)
}
