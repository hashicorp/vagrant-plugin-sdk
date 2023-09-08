// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// AuthFunc provides a mock function with given fields:
func (_m *Authenticator) AuthFunc() interface{} {
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

// ValidateAuthFunc provides a mock function with given fields:
func (_m *Authenticator) ValidateAuthFunc() interface{} {
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

type mockConstructorTestingTNewAuthenticator interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthenticator creates a new instance of Authenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthenticator(t mockConstructorTestingTNewAuthenticator) *Authenticator {
	mock := &Authenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
