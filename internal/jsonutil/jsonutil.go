package jsonutil

import (
	"encoding/json"
	"io"

	"github.com/hollgett/metricsYandex.git/internal/models"
)

func DecoderJSON(r io.Reader, v *models.Metrics) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

func EncoderJSON(w io.Writer, v models.Metrics) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}
