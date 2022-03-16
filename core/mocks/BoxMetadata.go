// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import (
	core "github.com/hashicorp/vagrant-plugin-sdk/core"
	mock "github.com/stretchr/testify/mock"
)

// BoxMetadata is an autogenerated mock type for the BoxMetadata type
type BoxMetadata struct {
	mock.Mock
}

// ListProviders provides a mock function with given fields: version
func (_m *BoxMetadata) ListProviders(version string) ([]string, error) {
	ret := _m.Called(version)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVersions provides a mock function with given fields: opts
func (_m *BoxMetadata) ListVersions(opts ...core.BoxMetadataOpts) ([]string, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []string
	if rf, ok := ret.Get(0).(func(...core.BoxMetadataOpts) []string); ok {
		r0 = rf(opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...core.BoxMetadataOpts) error); ok {
		r1 = rf(opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields:
func (_m *BoxMetadata) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Provider provides a mock function with given fields: version, name
func (_m *BoxMetadata) Provider(version string, name string) (core.BoxVersionProviderData, error) {
	ret := _m.Called(version, name)

	var r0 core.BoxVersionProviderData
	if rf, ok := ret.Get(0).(func(string, string) core.BoxVersionProviderData); ok {
		r0 = rf(version, name)
	} else {
		r0 = ret.Get(0).(core.BoxVersionProviderData)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(version, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Version provides a mock function with given fields: version, opts
func (_m *BoxMetadata) Version(version string, opts core.BoxMetadataOpts) (core.BoxVersionData, error) {
	ret := _m.Called(version, opts)

	var r0 core.BoxVersionData
	if rf, ok := ret.Get(0).(func(string, core.BoxMetadataOpts) core.BoxVersionData); ok {
		r0 = rf(version, opts)
	} else {
		r0 = ret.Get(0).(core.BoxVersionData)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, core.BoxMetadataOpts) error); ok {
		r1 = rf(version, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
