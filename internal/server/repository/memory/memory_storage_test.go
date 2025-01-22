package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetGauge(t *testing.T) {
	tests := []struct {
		name  string
		ms    *MemStorage
		nameV string
		value float64
		err   string
	}{
		{"positive", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 54.3, ""},
		{"negative name", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "", 54.3, "name gauge have nil"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.SetGauge(tt.nameV, tt.value); err != nil {
				assert.EqualError(t, err, tt.err, "error not equal")
			} else {
				assert.Contains(t, tt.ms.gauge, tt.nameV, "value does't exists")
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	tests := []struct {
		name  string
		ms    *MemStorage
		nameV string
		want  float64
		err   string
	}{
		{"positive", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 54.3, ""},
		{"negative name", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "", 54.3, "name gauge have nil"},
		{"negative data", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 0, "data does't exist"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == "" {
				tt.ms.gauge[tt.nameV] = tt.want
			}
			got, err := tt.ms.GetGauge(tt.nameV)
			if err != nil {
				assert.EqualError(t, err, tt.err, "error not equal")
			} else {
				assert.Equal(t, got, tt.want, "expected got not equal")
			}
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name  string
		ms    *MemStorage
		nameV string
		value int64
		err   string
	}{
		{"positive", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 54, ""},
		{"negative name", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "", 54, "name counter have nil"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.AddCounter(tt.nameV, tt.value); err != nil {
				assert.EqualError(t, err, tt.err, "error not equal")
			} else {
				assert.Contains(t, tt.ms.counter, tt.nameV, "value does't exists")
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	tests := []struct {
		name  string
		ms    *MemStorage
		nameV string
		want  int64
		err   string
	}{
		{"positive", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 54, ""},
		{"negative name", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "", 54, "name counter have nil"},
		{"negative data", &MemStorage{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		}, "test", 0, "data does't exist"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == "" {
				tt.ms.counter[tt.nameV] = tt.want
			}
			got, err := tt.ms.GetCounter(tt.nameV)
			if err != nil {
				assert.EqualError(t, err, tt.err, "error not equal")
			} else {
				assert.Equal(t, got, tt.want, "expected got not equal")
			}
		})
	}
}
