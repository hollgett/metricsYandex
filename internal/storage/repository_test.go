package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type want struct {
		expectedGauge float64
		expectedError error
	}
	tests := []struct {
		name       string
		nameMetric string
		val        []float64
		want       want
	}{
		{
			name:       "Test positive UpdateGauge #1",
			nameMetric: "Alloc",
			val:        []float64{3},
			want:       want{expectedGauge: 3, expectedError: nil},
		},
		{
			name:       "Test positive with more value UpdateGauge #2",
			nameMetric: "Alloc",
			val:        []float64{3, 4, 6, 3},
			want:       want{expectedGauge: 3, expectedError: nil},
		},
		{
			name:       "Test positive without value UpdateGauge #3",
			nameMetric: "Alloc",
			val:        nil,
			want:       want{expectedGauge: 0, expectedError: nil},
		},
		{
			name:       "Test positive with more after point #4",
			nameMetric: "Alloc",
			val:        []float64{3.34343e34, 4.3443553535, 6.34343434, 3.434343e3},
			want:       want{expectedGauge: 3.434343e3, expectedError: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemStorage()
			var err error
			for _, v := range tt.val {
				if err == nil {
					err = repo.UpdateGauge(tt.nameMetric, v)
				}
			}
			storage := repo.(*MemStorage)
			if tt.want.expectedError != nil {
				assert.Equal(t, tt.want.expectedError, err, "expected error not equal")
				return
			}
			assert.Equal(t, tt.want.expectedGauge, storage.gauge[tt.nameMetric], "expected value not equal")
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type want struct {
		expectedGauge int64
		expectedError error
	}
	tests := []struct {
		name       string
		nameMetric string
		val        []int64
		want       want
	}{
		{
			name:       "Test positive AddCounter #1",
			nameMetric: "Alloc",
			val:        []int64{3},
			want:       want{expectedGauge: 3, expectedError: nil},
		},
		{
			name:       "Test positive with more value AddCounter #2",
			nameMetric: "Alloc",
			val:        []int64{3, 4, 6, 3},
			want:       want{expectedGauge: 16, expectedError: nil},
		},
		{
			name:       "Test positive without value AddCounter #3",
			nameMetric: "AAA",
			val:        nil,
			want:       want{expectedGauge: 0, expectedError: nil},
		},
		{
			name: "Test negative without name metric AddCounter #1",
			val:  []int64{3},
			want: want{expectedGauge: 0, expectedError: errors.New("name metric have nil")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemStorage()
			var err error
			for _, v := range tt.val {
				if err == nil {
					err = repo.AddCounter(tt.nameMetric, v)
				}
			}
			storage := repo.(*MemStorage)
			if tt.want.expectedError != nil {
				assert.Equal(t, tt.want.expectedError, err, "expected error not equal")
				return
			}
			assert.Equal(t, tt.want.expectedGauge, storage.counter[tt.nameMetric], "expected value not equal")
		})
	}
}
