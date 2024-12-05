package handlers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hollgett/metricsYandex.git/internal/mock"
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
		requestURL string
		want       want
	}{
		{
			name:       "positive test gauge",
			repository: simulateRepository(controller, nil),
			requestURL: "gauge/Alloc/4.34343e3",
			want:       want{expectedErr: nil},
		},
		{
			name:       "positive test counter",
			repository: simulateRepository(controller, nil),
			requestURL: "counter/Alloc/4",
			want:       want{expectedErr: nil},
		},
		{
			name:       "negative test without name",
			repository: simulateRepository(controller, nil),
			requestURL: "counter//4",
			want:       want{expectedErr: errors.New("wrong request")},
		},
		{
			name:       "negative test without type metric",
			repository: simulateRepository(controller, nil),
			requestURL: "/Alloc/4",
			want:       want{expectedErr: errors.New("wrong type metric")},
		},
		{
			name:       "negative test without value metric",
			repository: simulateRepository(controller, nil),
			requestURL: "counter/Alloc/",
			want:       want{expectedErr: errors.New("wrong value")},
		},
		{
			name:       "negative test without request data",
			repository: simulateRepository(controller, nil),
			requestURL: "",
			want:       want{expectedErr: errors.New("wrong request")},
		},
		{
			name:       "negative test with string value counter",
			repository: simulateRepository(controller, nil),
			requestURL: "counter/Alloc/abc",
			want:       want{expectedErr: errors.New("wrong value")},
		},
		{
			name:       "negative test with string value gauge",
			repository: simulateRepository(controller, nil),
			requestURL: "gauge/Alloc/abc",
			want:       want{expectedErr: errors.New("wrong value")},
		},
		{
			name:       "negative test with error repository",
			repository: simulateRepository(controller, errors.New("name metric have nil")),
			requestURL: "gauge/Alloc/4.3e34",
			want:       want{expectedErr: errors.Join(errors.New("function UpdateGauge have error: "), errors.New("name metric have nil"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricHandler := NewMetricHandler(tt.repository)
			got := metricHandler.CollectingMetric(tt.requestURL)
			assert.Equal(t, tt.want.expectedErr, got, "expected error not equal")

			if got != nil {
				assert.Equal(t, tt.want.expectedErr, got, "expected error not equal: ", got.Error())
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
