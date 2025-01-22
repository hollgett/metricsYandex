package utils

import (
	"encoding/json"
	"io"
)

func DecoderJSON(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

func EncoderJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}
