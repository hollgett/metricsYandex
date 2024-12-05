package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestApiMetric_CheckURLMiddleware(t *testing.T) {
	type want struct {
		wantStatus int
	}
	tests := []struct {
		name    string
		method  string
		request string

		want want
	}{
		{
			name:    "positive test #1",
			method:  http.MethodPost,
			request: "http://localhost:8080/update/counter/test/34",
			want: want{
				wantStatus: http.StatusOK,
			},
		},
		{
			name:    "positive test #2",
			method:  http.MethodPost,
			request: "http://localhost:8080/update/gauge/test/34",
			want: want{
				wantStatus: http.StatusOK,
			},
		},
		{
			name:    "negative test method request",
			method:  http.MethodGet,
			request: "http://localhost:8080/update/counter/test/34",
			want: want{
				wantStatus: http.StatusMethodNotAllowed,
			},
		},
		{
			name:    "negative test metric method",
			method:  http.MethodPost,
			request: "http://localhost:8080/gupdate/counter/test/34",
			want: want{
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name:    "negative test metric type",
			method:  http.MethodPost,
			request: "http://localhost:8080/update/cofunter/test/34",
			want: want{
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name:    "negative test metric name",
			method:  http.MethodPost,
			request: "http://localhost:8080/update/counter//34",
			want: want{
				wantStatus: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			apiMetirc := &ApiMetric{}
			test := apiMetirc.CheckURLMiddleware(middleware)

			r := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			test.ServeHTTP(w, r)
			res := w.Result()

			assert.Equal(t, tt.want.wantStatus, res.StatusCode, "response status code not equal with expected")

		})
	}
}

func TestApiMetric_UpdateMetricPost(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	type want struct {
		wantStatus int
	}
	tests := []struct {
		name          string
		metricHandler handlers.MetricHandler
		want          want
	}{
		{
			name:          "positive test",
			metricHandler: simulateMetricHandler(controller, nil),
			want: want{
				wantStatus: http.StatusOK,
			},
		},
		{
			name:          "negative test",
			metricHandler: simulateMetricHandler(controller, errors.New("error")),
			want: want{
				wantStatus: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
			w := httptest.NewRecorder()
			api := &ApiMetric{
				handler: tt.metricHandler,
			}
			api.UpdateMetricPost(w, r)

			res := w.Result()

			assert.Equal(t, tt.want.wantStatus, res.StatusCode, "response status code not equal with expected")
		})
	}
}

func simulateMetricHandler(ctrl *gomock.Controller, err error) handlers.MetricHandler {
	mockHandler := mock.NewMockMetricHandler(ctrl)

	mockHandler.EXPECT().CollectingMetric(gomock.Any()).Return(err).AnyTimes()

	return mockHandler
}
