package utils

import (
	"encoding/json"

	"github.com/hollgett/metricsYandex.git/internal/agent/models"
)

func Marshal(requestdata models.Metrics) ([]byte, error) {
	return json.Marshal(requestdata)
}
