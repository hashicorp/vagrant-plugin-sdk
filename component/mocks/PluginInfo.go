// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	component "github.com/hashicorp/vagrant-plugin-sdk/component"
	mock "github.com/stretchr/testify/mock"
)

// PluginInfo is an autogenerated mock type for the PluginInfo type
type PluginInfo struct {
	mock.Mock
}

// ComponentOptions provides a mock function with given fields:
func (_m *PluginInfo) ComponentOptions() map[component.Type]interface{} {
	ret := _m.Called()

	var r0 map[component.Type]interface{}
	if rf, ok := ret.Get(0).(func() map[component.Type]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[component.Type]interface{})
		}
	}

	return r0
}

// ComponentTypes provides a mock function with given fields:
func (_m *PluginInfo) ComponentTypes() []component.Type {
	ret := _m.Called()

	var r0 []component.Type
	if rf, ok := ret.Get(0).(func() []component.Type); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]component.Type)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *PluginInfo) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewPluginInfo interface {
	mock.TestingT
	Cleanup(func())
}

// NewPluginInfo creates a new instance of PluginInfo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPluginInfo(t mockConstructorTestingTNewPluginInfo) *PluginInfo {
	mock := &PluginInfo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
