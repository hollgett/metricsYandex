package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/hollgett/metricsYandex.git/internal/handlers"
	"github.com/hollgett/metricsYandex.git/internal/mock"
	"github.com/hollgett/metricsYandex.git/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRouter(h handlers.MetricHandler) chi.Router {
	api := NewAPIMetric(h)
	rtr := chi.NewMux()
	rtr.Route("/", func(r chi.Router) {
		r.Get("/", api.GetMetricAll)
		r.Route("/value", func(r chi.Router) {
			r.Get("/{typeM}/{nameM}", api.GetMetric)
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/{typeM}/{nameM}/{valueM}", api.UpdateMetricPost)
		})
	})
	return rtr
}

func testRequest(t *testing.T, method string, ts *httptest.Server, reqURL string) *resty.Response {
	client := resty.New().
		SetBaseURL(ts.URL).
		SetHeader("Content-Type", "text/plain")
	if method == http.MethodGet {
		resp, err := client.R().Get(reqURL)
		require.NoError(t, err)
		return resp
	}
	resp, err := client.R().Post(reqURL)
	require.NoError(t, err)
	return resp
}

func TestApiMetric_UpdateMetricPost(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name          string
		metricHandler handlers.MetricHandler
		method        string
		request       string
		wantStatus    int
	}{
		{
			name:          "gauge request",
			method:        http.MethodPost,
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/update/gauge/a/5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "counter request",
			method:        http.MethodPost,
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/update/counter/a/5",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "test without type metric",
			method:        http.MethodPost,
			request:       "/update//Alloc/5",
			metricHandler: simulateMetricHandler(controller, nil),
			wantStatus:    http.StatusBadRequest,
		},
		{
			name:          "test without name metric",
			method:        http.MethodPost,
			request:       "/update/counter//5",
			metricHandler: simulateMetricHandler(controller, nil),
			wantStatus:    http.StatusNotFound,
		},
		{
			name:          "test with service metric error",
			method:        http.MethodPost,
			request:       "/update/counter/Alloc/5",
			metricHandler: simulateMetricHandler(controller, errors.New("error")),
			wantStatus:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(testRouter(tt.metricHandler))
			defer testServer.Close()

			resp := testRequest(t, tt.method, testServer, tt.request)

			assert.Equal(t, tt.wantStatus, resp.StatusCode(), "response status code not equal with expected: "+string(resp.Body()))
		})
	}
}

func TestApiMetric_GetMetric(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tests := []struct {
		name          string
		metricHandler handlers.MetricHandler
		request       string
		expectedBody  string
		expectedCode  int
	}{
		{
			name:          "get metric counter",
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/value/counter/F",
			expectedBody:  "result",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "get metric gauge",
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/value/gauge/F",
			expectedBody:  "result",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "error metricHandler",
			metricHandler: simulateMetricHandler(controller, errors.New("error")),
			request:       "/value/gauge/F",
			expectedBody:  "result",
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "path error",
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/values/gauge/F",
			expectedBody:  "result",
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "error type metric",
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/value/gauges/F",
			expectedBody:  "result",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "error name metric",
			metricHandler: simulateMetricHandler(controller, nil),
			request:       "/value/gauges/",
			expectedBody:  "result",
			expectedCode:  http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(testRouter(tt.metricHandler))
			defer testServer.Close()

			resp := testRequest(t, http.MethodGet, testServer, tt.request)

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "response status code not equal with expected: "+resp.String())
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, "text/plain", resp.Header().Get("Content-Type"), "response header not equal with expected: "+resp.String())
				assert.Equal(t, tt.expectedBody, resp.String(), "response body not equal with expected: "+resp.String())
			}

		})
	}
}

func simulateMetricHandler(ctrl *gomock.Controller, err error) handlers.MetricHandler {
	mockHandler := mock.NewMockMetricHandler(ctrl)

	mockHandler.EXPECT().CollectingMetric(gomock.AssignableToTypeOf(models.Metrics{})).Return(err).AnyTimes()
	mockHandler.EXPECT().GetMetric(gomock.AssignableToTypeOf(models.Metrics{})).Return("result", err).AnyTimes()

	return mockHandler
}
