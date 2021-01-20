// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MessageHandler is an autogenerated mock type for the MessageHandler type
type MessageHandler struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *MessageHandler) Done() {
	_m.Called()
}

// Handle provides a mock function with given fields: rawMessage
func (_m *MessageHandler) Handle(rawMessage []byte) error {
	ret := _m.Called(rawMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(rawMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}