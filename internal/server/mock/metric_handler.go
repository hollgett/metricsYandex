// Code generated by MockGen. DO NOT EDIT.
// Source: metric_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/hollgett/metricsYandex.git/internal/server/models"
)

// MockMetricHandler is a mock of MetricHandler interface.
type MockMetricHandler struct {
	ctrl     *gomock.Controller
	recorder *MockMetricHandlerMockRecorder
}

// MockMetricHandlerMockRecorder is the mock recorder for MockMetricHandler.
type MockMetricHandlerMockRecorder struct {
	mock *MockMetricHandler
}

// NewMockMetricHandler creates a new mock instance.
func NewMockMetricHandler(ctrl *gomock.Controller) *MockMetricHandler {
	mock := &MockMetricHandler{ctrl: ctrl}
	mock.recorder = &MockMetricHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricHandler) EXPECT() *MockMetricHandlerMockRecorder {
	return m.recorder
}

// Batch mocks base method.
func (m *MockMetricHandler) Batch(ctx context.Context, metrics []models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Batch", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// Batch indicates an expected call of Batch.
func (mr *MockMetricHandlerMockRecorder) Batch(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Batch", reflect.TypeOf((*MockMetricHandler)(nil).Batch), ctx, metrics)
}

// CollectingMetric mocks base method.
func (m *MockMetricHandler) CollectingMetric(metric *models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollectingMetric", metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// CollectingMetric indicates an expected call of CollectingMetric.
func (mr *MockMetricHandlerMockRecorder) CollectingMetric(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollectingMetric", reflect.TypeOf((*MockMetricHandler)(nil).CollectingMetric), metric)
}

// GetMetric mocks base method.
func (m *MockMetricHandler) GetMetric(metric *models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricHandlerMockRecorder) GetMetric(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricHandler)(nil).GetMetric), metric)
}

// GetMetricAll mocks base method.
func (m *MockMetricHandler) GetMetricAll() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricAll")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricAll indicates an expected call of GetMetricAll.
func (mr *MockMetricHandlerMockRecorder) GetMetricAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricAll", reflect.TypeOf((*MockMetricHandler)(nil).GetMetricAll))
}

// PingDB mocks base method.
func (m *MockMetricHandler) PingDB(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingDB", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingDB indicates an expected call of PingDB.
func (mr *MockMetricHandlerMockRecorder) PingDB(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingDB", reflect.TypeOf((*MockMetricHandler)(nil).PingDB), ctx)
}

// ValidateMetric mocks base method.
func (m *MockMetricHandler) ValidateMetric(metric *models.Metrics) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateMetric", metric)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateMetric indicates an expected call of ValidateMetric.
func (mr *MockMetricHandlerMockRecorder) ValidateMetric(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateMetric", reflect.TypeOf((*MockMetricHandler)(nil).ValidateMetric), metric)
}
