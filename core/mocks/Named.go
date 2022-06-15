// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Named is an autogenerated mock type for the Named type
type Named struct {
	mock.Mock
}

// PluginName provides a mock function with given fields:
func (_m *Named) PluginName() (string, error) {
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

// SetPluginName provides a mock function with given fields: _a0
func (_m *Named) SetPluginName(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewNamedT interface {
	mock.TestingT
	Cleanup(func())
}

// NewNamed creates a new instance of Named. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNamed(t NewNamedT) *Named {
	mock := &Named{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
