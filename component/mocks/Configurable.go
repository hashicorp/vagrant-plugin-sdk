// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Configurable is an autogenerated mock type for the Configurable type
type Configurable struct {
	mock.Mock
}

// Config provides a mock function with given fields:
func (_m *Configurable) Config() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
