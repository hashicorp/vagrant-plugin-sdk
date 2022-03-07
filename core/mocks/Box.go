// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import (
	core "github.com/hashicorp/vagrant-plugin-sdk/core"
	mock "github.com/stretchr/testify/mock"
)

// Box is an autogenerated mock type for the Box type
type Box struct {
	mock.Mock
}

// AutomaticUpdateCheckAllowed provides a mock function with given fields:
func (_m *Box) AutomaticUpdateCheckAllowed() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BoxMetadata provides a mock function with given fields:
func (_m *Box) BoxMetadata() (map[string]interface{}, error) {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
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

// Compare provides a mock function with given fields: box
func (_m *Box) Compare(box core.Box) (int, error) {
	ret := _m.Called(box)

	var r0 int
	if rf, ok := ret.Get(0).(func(core.Box) int); ok {
		r0 = rf(box)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(core.Box) error); ok {
		r1 = rf(box)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields:
func (_m *Box) Destroy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Directory provides a mock function with given fields:
func (_m *Box) Directory() (string, error) {
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

// HasUpdate provides a mock function with given fields: version
func (_m *Box) HasUpdate(version string) (bool, error) {
	ret := _m.Called(version)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(version)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InUse provides a mock function with given fields: index
func (_m *Box) InUse(index core.TargetIndex) (bool, error) {
	ret := _m.Called(index)

	var r0 bool
	if rf, ok := ret.Get(0).(func(core.TargetIndex) bool); ok {
		r0 = rf(index)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(core.TargetIndex) error); ok {
		r1 = rf(index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Machines provides a mock function with given fields: index
func (_m *Box) Machines(index core.TargetIndex) ([]core.Machine, error) {
	ret := _m.Called(index)

	var r0 []core.Machine
	if rf, ok := ret.Get(0).(func(core.TargetIndex) []core.Machine); ok {
		r0 = rf(index)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Machine)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(core.TargetIndex) error); ok {
		r1 = rf(index)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Metadata provides a mock function with given fields:
func (_m *Box) Metadata() (core.BoxMetadata, error) {
	ret := _m.Called()

	var r0 core.BoxMetadata
	if rf, ok := ret.Get(0).(func() core.BoxMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.BoxMetadata)
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

// MetadataURL provides a mock function with given fields:
func (_m *Box) MetadataURL() (string, error) {
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

// Name provides a mock function with given fields:
func (_m *Box) Name() (string, error) {
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

// Provider provides a mock function with given fields:
func (_m *Box) Provider() (string, error) {
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

// Repackage provides a mock function with given fields: path
func (_m *Box) Repackage(path string) error {
	ret := _m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateInfo provides a mock function with given fields: version
func (_m *Box) UpdateInfo(version string) (bool, core.BoxMetadataMap, string, string, error) {
	ret := _m.Called(version)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(version)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 core.BoxMetadataMap
	if rf, ok := ret.Get(1).(func(string) core.BoxMetadataMap); ok {
		r1 = rf(version)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(core.BoxMetadataMap)
		}
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(string) string); ok {
		r2 = rf(version)
	} else {
		r2 = ret.Get(2).(string)
	}

	var r3 string
	if rf, ok := ret.Get(3).(func(string) string); ok {
		r3 = rf(version)
	} else {
		r3 = ret.Get(3).(string)
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(string) error); ok {
		r4 = rf(version)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

// Version provides a mock function with given fields:
func (_m *Box) Version() (string, error) {
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
