// Code generated by mockery 2.12.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// ActionFunc provides a mock function with given fields: actionName
func (_m *Provider) ActionFunc(actionName string) interface{} {
	ret := _m.Called(actionName)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(actionName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// CapabilityFunc provides a mock function with given fields: capName
func (_m *Provider) CapabilityFunc(capName string) interface{} {
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

// HasCapabilityFunc provides a mock function with given fields:
func (_m *Provider) HasCapabilityFunc() interface{} {
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

// InstalledFunc provides a mock function with given fields:
func (_m *Provider) InstalledFunc() interface{} {
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

// MachineIdChangedFunc provides a mock function with given fields:
func (_m *Provider) MachineIdChangedFunc() interface{} {
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

// SshInfoFunc provides a mock function with given fields:
func (_m *Provider) SshInfoFunc() interface{} {
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

// StateFunc provides a mock function with given fields:
func (_m *Provider) StateFunc() interface{} {
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

// UsableFunc provides a mock function with given fields:
func (_m *Provider) UsableFunc() interface{} {
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

type NewProviderT interface {
	mock.TestingT
	Cleanup(func())
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProvider(t NewProviderT) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
