// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Host is an autogenerated mock type for the Host type
type Host struct {
	mock.Mock
}

// CapabilityFunc provides a mock function with given fields: capName
func (_m *Host) CapabilityFunc(capName string) interface{} {
	ret := _m.Called(capName)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(capName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// DetectFunc provides a mock function with given fields:
func (_m *Host) DetectFunc() interface{} {
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

// HasCapabilityFunc provides a mock function with given fields:
func (_m *Host) HasCapabilityFunc() interface{} {
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
