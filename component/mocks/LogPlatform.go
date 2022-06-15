// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// LogPlatform is an autogenerated mock type for the LogPlatform type
type LogPlatform struct {
	mock.Mock
}

// LogsFunc provides a mock function with given fields:
func (_m *LogPlatform) LogsFunc() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

type NewLogPlatformT interface {
	mock.TestingT
	Cleanup(func())
}

// NewLogPlatform creates a new instance of LogPlatform. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogPlatform(t NewLogPlatformT) *LogPlatform {
	mock := &LogPlatform{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
