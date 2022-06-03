// Code generated by mockery 2.12.3. DO NOT EDIT.

package mocks

import (
	core "github.com/hashicorp/vagrant-plugin-sdk/core"
	mock "github.com/stretchr/testify/mock"
)

// Host is an autogenerated mock type for the Host type
type Host struct {
	mock.Mock
}

// Capability provides a mock function with given fields: name, args
func (_m *Host) Capability(name string, args ...interface{}) (interface{}, error) {
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) interface{}); ok {
		r0 = rf(name, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(name, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Host) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Detect provides a mock function with given fields: state
func (_m *Host) Detect(state core.StateBag) (bool, error) {
	ret := _m.Called(state)

	var r0 bool
	if rf, ok := ret.Get(0).(func(core.StateBag) bool); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(core.StateBag) error); ok {
		r1 = rf(state)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasCapability provides a mock function with given fields: name
func (_m *Host) HasCapability(name string) (bool, error) {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Parent provides a mock function with given fields:
func (_m *Host) Parent() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PluginName provides a mock function with given fields:
func (_m *Host) PluginName() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Seed provides a mock function with given fields: _a0
func (_m *Host) Seed(_a0 *core.Seeds) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.Seeds) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Seeds provides a mock function with given fields:
func (_m *Host) Seeds() (*core.Seeds, error) {
	ret := _m.Called()

	var r0 *core.Seeds
	if rf, ok := ret.Get(0).(func() *core.Seeds); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Seeds)
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

// SetPluginName provides a mock function with given fields: _a0
func (_m *Host) SetPluginName(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewHostT interface {
	mock.TestingT
	Cleanup(func())
}

// NewHost creates a new instance of Host. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHost(t NewHostT) *Host {
	mock := &Host{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
