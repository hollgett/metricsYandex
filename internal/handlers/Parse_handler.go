package handlers

import (
	"fmt"
	"strconv"
)

// take gauge type metric, parse value to float64
func GaugeParse(value string) (float64, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("gauge case parse got value: %s, error: %w", value, err)
	}
	return val, nil
}

// take counter type metric, parse value to int64
func CounterParse(value string) (int64, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("counter case parse got value: %s, error: %w", value, err)
	}
	return val, nil
}
