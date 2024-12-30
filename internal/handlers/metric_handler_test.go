package handlers

import (
	"errors"
	"fmt"
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
		requestURL []string
		want       want
	}{
		{
			name:       "positive test gauge",
			repository: simulateRepository(controller, nil),
			requestURL: []string{"gauge", "Alloc", "4.34343e3"},
			want:       want{expectedErr: nil},
		},
		{
			name:       "positive test counter",
			repository: simulateRepository(controller, nil),
			requestURL: []string{"counter", "Alloc", "4"},
			want:       want{expectedErr: nil},
		},
		{
			name:       "negative test without type metric",
			repository: simulateRepository(controller, nil),
			requestURL: []string{"", "Alloc", "4"},
			want:       want{expectedErr: errors.New("wrong type metric")},
		},
		{
			name:       "negative test with string value counter",
			repository: simulateRepository(controller, nil),
			requestURL: []string{"counter", "Alloc", "abc"},
			want:       want{expectedErr: errors.New("wrong value")},
		},
		{
			name:       "negative test with string value gauge",
			repository: simulateRepository(controller, nil),
			requestURL: []string{"gauge", "Alloc", "abc"},
			want:       want{expectedErr: errors.New("wrong value")},
		},
		{
			name:       "negative test with error repository",
			repository: simulateRepository(controller, errors.New("name metric have nil")),
			requestURL: []string{"gauge", "Alloc", "4.3e34"},
			want:       want{expectedErr: fmt.Errorf("function UpdateGauge have error: %w", errors.New("name metric have nil"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricHandler := NewMetricHandler(tt.repository)
			get := metricHandler.CollectingMetric(&models.Metrics{})
			assert.Equal(t, tt.want.expectedErr, get, "expected error not equal")

			if get != nil {
				assert.Equal(t, tt.want.expectedErr, get, "expected error not equal: ", get.Error())
			}

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
