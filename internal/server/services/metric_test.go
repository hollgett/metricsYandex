package services

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hollgett/metricsYandex.git/internal/server/mock"
	"github.com/hollgett/metricsYandex.git/internal/server/models"
	"github.com/stretchr/testify/assert"
)

func Test_metricHandler_ValidateMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		m          MetricHandler
		metric     models.Metrics
		wantStatus int
		wantErr    error
	}{
		{"positive", repositoryMock(ctrl), models.Metrics{ID: "test", MType: "gauge"}, 0, nil},
		{"without", repositoryMock(ctrl), models.Metrics{ID: "", MType: ""}, 404, fmt.Errorf("name metric got nill")},
		{"without name", repositoryMock(ctrl), models.Metrics{ID: "", MType: "gauge"}, 404, fmt.Errorf("name metric got nill")},
		{"without type", repositoryMock(ctrl), models.Metrics{ID: "test", MType: "gauges"}, 400, fmt.Errorf("type metric got nill")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := tt.m.ValidateMetric(&tt.metric)
			if err != nil {
				assert.Equal(t, tt.wantErr, err, "error not equal")
			}
			assert.Equal(t, tt.wantStatus, status, "return status not equal")
		})
	}
}

func repositoryMock(ctrl *gomock.Controller) MetricHandler {
	controller := mock.NewMockRepository(ctrl)

	return NewMetricHandler(controller)
}
