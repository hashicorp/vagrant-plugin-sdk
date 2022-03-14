// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	component "github.com/hashicorp/vagrant-plugin-sdk/component"

	mock "github.com/stretchr/testify/mock"
)

// Config is an autogenerated mock type for the Config type
type Config struct {
	mock.Mock
}

// Finalize provides a mock function with given fields: _a0
func (_m *Config) Finalize(_a0 *component.ConfigData) (*component.ConfigData, error) {
	ret := _m.Called(_a0)

	var r0 *component.ConfigData
	var r1 error
	if rf, ok := ret.Get(0).(func(*component.ConfigData) (*component.ConfigData, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*component.ConfigData) *component.ConfigData); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*component.ConfigData)
		}
	}

	if rf, ok := ret.Get(1).(func(*component.ConfigData) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Merge provides a mock function with given fields: base, toMerge
func (_m *Config) Merge(base *component.ConfigData, toMerge *component.ConfigData) (*component.ConfigData, error) {
	ret := _m.Called(base, toMerge)

	var r0 *component.ConfigData
	var r1 error
	if rf, ok := ret.Get(0).(func(*component.ConfigData, *component.ConfigData) (*component.ConfigData, error)); ok {
		return rf(base, toMerge)
	}
	if rf, ok := ret.Get(0).(func(*component.ConfigData, *component.ConfigData) *component.ConfigData); ok {
		r0 = rf(base, toMerge)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*component.ConfigData)
		}
	}

	if rf, ok := ret.Get(1).(func(*component.ConfigData, *component.ConfigData) error); ok {
		r1 = rf(base, toMerge)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields:
func (_m *Config) Register() (*component.ConfigRegistration, error) {
	ret := _m.Called()

	var r0 *component.ConfigRegistration
	var r1 error
	if rf, ok := ret.Get(0).(func() (*component.ConfigRegistration, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *component.ConfigRegistration); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*component.ConfigRegistration)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Struct provides a mock function with given fields:
func (_m *Config) Struct() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func() (interface{}, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfig creates a new instance of Config. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfig(t mockConstructorTestingTNewConfig) *Config {
	mock := &Config{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
