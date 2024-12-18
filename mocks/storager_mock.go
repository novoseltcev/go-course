// Code generated by MockGen. DO NOT EDIT.
// Source: storager.go
//
// Generated by this command:
//
//	mockgen -source=storager.go -destination=../../mocks/storager_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	schemas "github.com/novoseltcev/go-course/internal/schemas"
	gomock "go.uber.org/mock/gomock"
)

// MockMetricStorager is a mock of MetricStorager interface.
type MockMetricStorager struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStoragerMockRecorder
	isgomock struct{}
}

// MockMetricStoragerMockRecorder is the mock recorder for MockMetricStorager.
type MockMetricStoragerMockRecorder struct {
	mock *MockMetricStorager
}

// NewMockMetricStorager creates a new mock instance.
func NewMockMetricStorager(ctrl *gomock.Controller) *MockMetricStorager {
	mock := &MockMetricStorager{ctrl: ctrl}
	mock.recorder = &MockMetricStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricStorager) EXPECT() *MockMetricStoragerMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockMetricStorager) GetAll(ctx context.Context) ([]schemas.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]schemas.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMetricStoragerMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMetricStorager)(nil).GetAll), ctx)
}

// GetOne mocks base method.
func (m *MockMetricStorager) GetOne(ctx context.Context, id, mType string) (*schemas.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", ctx, id, mType)
	ret0, _ := ret[0].(*schemas.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockMetricStoragerMockRecorder) GetOne(ctx, id, mType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockMetricStorager)(nil).GetOne), ctx, id, mType)
}

// Ping mocks base method.
func (m *MockMetricStorager) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockMetricStoragerMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMetricStorager)(nil).Ping), ctx)
}

// Save mocks base method.
func (m *MockMetricStorager) Save(ctx context.Context, metric *schemas.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockMetricStoragerMockRecorder) Save(ctx, metric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockMetricStorager)(nil).Save), ctx, metric)
}

// SaveBatch mocks base method.
func (m *MockMetricStorager) SaveBatch(ctx context.Context, metrics []schemas.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockMetricStoragerMockRecorder) SaveBatch(ctx, metrics any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockMetricStorager)(nil).SaveBatch), ctx, metrics)
}
