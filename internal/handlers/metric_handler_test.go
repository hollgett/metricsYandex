package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hollgett/metricsYandex.git/internal/mock"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"github.com/hollgett/metricsYandex.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_metricHandler_CollectingMetric(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	type want struct {
		expectedErr error
	}
	tests := []struct {
		name       string
		repository storage.Repositories
		requestURL models.Metrics
		want       want
	}{
		{
			name:       "positive test gauge",
			repository: simulateRepository(controller, nil),
			requestURL: models.Metrics{MType: "gauge",
				ID: "Alloc"},
			want: want{expectedErr: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricHandler := NewMetricHandler(tt.repository)
			var float64 = 5.343
			tt.requestURL.Value = &float64
			get := metricHandler.CollectingMetric(&tt.requestURL)
			assert.Equal(t, tt.want.expectedErr, get, "expected error not equal")
		})
	}
}

func simulateRepository(controller *gomock.Controller, err error) storage.Repositories {
	mockRepository := mock.NewMockRepositories(controller)

	// creating mock for method AddCounter
	mockRepository.EXPECT().
		AddCounter(gomock.Any(), gomock.Any()).
		Return(err).
		AnyTimes()

	// creating mock for method UpdateGauge
	mockRepository.EXPECT().
		UpdateGauge(gomock.Any(), gomock.Any()).
		Return(err).
		AnyTimes()

	return mockRepository
}
