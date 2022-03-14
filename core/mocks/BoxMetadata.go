// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	core "github.com/hashicorp/vagrant-plugin-sdk/core"
	mock "github.com/stretchr/testify/mock"
)

// BoxMetadata is an autogenerated mock type for the BoxMetadata type
type BoxMetadata struct {
	mock.Mock
}

// BoxName provides a mock function with given fields:
func (_m *BoxMetadata) BoxName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ListProviders provides a mock function with given fields: version
func (_m *BoxMetadata) ListProviders(version string) ([]string, error) {
	ret := _m.Called(version)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(version)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVersions provides a mock function with given fields: opts
func (_m *BoxMetadata) ListVersions(opts ...*core.BoxProvider) ([]string, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(...*core.BoxProvider) ([]string, error)); ok {
		return rf(opts...)
	}
	if rf, ok := ret.Get(0).(func(...*core.BoxProvider) []string); ok {
		r0 = rf(opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(...*core.BoxProvider) error); ok {
		r1 = rf(opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadMetadata provides a mock function with given fields: url
func (_m *BoxMetadata) LoadMetadata(url string) error {
	ret := _m.Called(url)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Provider provides a mock function with given fields: version, name
func (_m *BoxMetadata) Provider(version string, name string) (*core.BoxProvider, error) {
	ret := _m.Called(version, name)

	var r0 *core.BoxProvider
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*core.BoxProvider, error)); ok {
		return rf(version, name)
	}
	if rf, ok := ret.Get(0).(func(string, string) *core.BoxProvider); ok {
		r0 = rf(version, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.BoxProvider)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(version, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Version provides a mock function with given fields: version, opts
func (_m *BoxMetadata) Version(version string, opts ...*core.BoxProvider) (*core.BoxVersion, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, version)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *core.BoxVersion
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...*core.BoxProvider) (*core.BoxVersion, error)); ok {
		return rf(version, opts...)
	}
	if rf, ok := ret.Get(0).(func(string, ...*core.BoxProvider) *core.BoxVersion); ok {
		r0 = rf(version, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.BoxVersion)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...*core.BoxProvider) error); ok {
		r1 = rf(version, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBoxMetadata interface {
	mock.TestingT
	Cleanup(func())
}

// NewBoxMetadata creates a new instance of BoxMetadata. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBoxMetadata(t mockConstructorTestingTNewBoxMetadata) *BoxMetadata {
	mock := &BoxMetadata{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
