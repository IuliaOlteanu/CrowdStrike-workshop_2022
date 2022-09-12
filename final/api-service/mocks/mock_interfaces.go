// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "exam/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockQueue) Add(batch []domain.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", batch)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockQueueMockRecorder) Add(batch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockQueue)(nil).Add), batch)
}
