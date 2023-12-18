// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ConfigFn is an autogenerated mock type for the ConfigFn type
type ConfigFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: configStruct
func (_m *ConfigFn) Execute(configStruct interface{}) (interface{}, error) {
	ret := _m.Called(configStruct)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(configStruct)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(configStruct)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}