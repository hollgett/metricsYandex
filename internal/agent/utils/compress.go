package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func CompressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("write started: %w", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}
	return b.Bytes(), nil
}
