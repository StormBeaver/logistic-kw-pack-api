// Code generated by MockGen. DO NOT EDIT.
// Source: route255/logistic-kw-pack-api/internal/app/sender (interfaces: EventSender)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	model "route255/logistic-kw-pack-api/internal/model"

	gomock "github.com/golang/mock/gomock"
)

// MockEventSender is a mock of EventSender interface.
type MockEventSender struct {
	ctrl     *gomock.Controller
	recorder *MockEventSenderMockRecorder
}

// MockEventSenderMockRecorder is the mock recorder for MockEventSender.
type MockEventSenderMockRecorder struct {
	mock *MockEventSender
}

// NewMockEventSender creates a new mock instance.
func NewMockEventSender(ctrl *gomock.Controller) *MockEventSender {
	mock := &MockEventSender{ctrl: ctrl}
	mock.recorder = &MockEventSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventSender) EXPECT() *MockEventSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEventSender) Send(arg0 *model.PackEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEventSenderMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEventSender)(nil).Send), arg0)
}
