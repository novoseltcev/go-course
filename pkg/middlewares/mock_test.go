// Code generated by MockGen. DO NOT EDIT.
// Source: decrypt.go
//
// Generated by this command:
//
//	mockgen -source=decrypt.go -destination=./mock_test.go -package=middlewares_test
//

// Package middlewares_test is a generated GoMock package.
package middlewares_test

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockdecryptor is a mock of decryptor interface.
type Mockdecryptor struct {
	ctrl     *gomock.Controller
	recorder *MockdecryptorMockRecorder
	isgomock struct{}
}

// MockdecryptorMockRecorder is the mock recorder for Mockdecryptor.
type MockdecryptorMockRecorder struct {
	mock *Mockdecryptor
}

// NewMockdecryptor creates a new mock instance.
func NewMockdecryptor(ctrl *gomock.Controller) *Mockdecryptor {
	mock := &Mockdecryptor{ctrl: ctrl}
	mock.recorder = &MockdecryptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdecryptor) EXPECT() *MockdecryptorMockRecorder {
	return m.recorder
}

// Decrypt mocks base method.
func (m *Mockdecryptor) Decrypt(b []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", b)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt.
func (mr *MockdecryptorMockRecorder) Decrypt(b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*Mockdecryptor)(nil).Decrypt), b)
}
