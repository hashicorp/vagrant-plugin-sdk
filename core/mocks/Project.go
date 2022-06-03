// Code generated by mockery 2.12.3. DO NOT EDIT.

package mocks

import (
	core "github.com/hashicorp/vagrant-plugin-sdk/core"
	datadir "github.com/hashicorp/vagrant-plugin-sdk/datadir"

	mock "github.com/stretchr/testify/mock"

	path "github.com/hashicorp/vagrant-plugin-sdk/helper/path"

	terminal "github.com/hashicorp/vagrant-plugin-sdk/terminal"

	vagrant_plugin_sdk "github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// Project is an autogenerated mock type for the Project type
type Project struct {
	mock.Mock
}

// ActiveTargets provides a mock function with given fields:
func (_m *Project) ActiveTargets() ([]core.Target, error) {
	ret := _m.Called()

	var r0 []core.Target
	if rf, ok := ret.Get(0).(func() []core.Target); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Target)
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

// Boxes provides a mock function with given fields:
func (_m *Project) Boxes() (core.BoxCollection, error) {
	ret := _m.Called()

	var r0 core.BoxCollection
	if rf, ok := ret.Get(0).(func() core.BoxCollection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.BoxCollection)
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

// CWD provides a mock function with given fields:
func (_m *Project) CWD() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// Close provides a mock function with given fields:
func (_m *Project) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Config provides a mock function with given fields:
func (_m *Project) Config() (*vagrant_plugin_sdk.Vagrantfile_Vagrantfile, error) {
	ret := _m.Called()

	var r0 *vagrant_plugin_sdk.Vagrantfile_Vagrantfile
	if rf, ok := ret.Get(0).(func() *vagrant_plugin_sdk.Vagrantfile_Vagrantfile); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vagrant_plugin_sdk.Vagrantfile_Vagrantfile)
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

// DataDir provides a mock function with given fields:
func (_m *Project) DataDir() (*datadir.Project, error) {
	ret := _m.Called()

	var r0 *datadir.Project
	if rf, ok := ret.Get(0).(func() *datadir.Project); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*datadir.Project)
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

// DefaultPrivateKey provides a mock function with given fields:
func (_m *Project) DefaultPrivateKey() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// DefaultProvider provides a mock function with given fields: opts
func (_m *Project) DefaultProvider(opts *core.DefaultProviderOptions) (string, error) {
	ret := _m.Called(opts)

	var r0 string
	if rf, ok := ret.Get(0).(func(*core.DefaultProviderOptions) string); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*core.DefaultProviderOptions) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Home provides a mock function with given fields:
func (_m *Project) Home() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// Host provides a mock function with given fields:
func (_m *Project) Host() (core.Host, error) {
	ret := _m.Called()

	var r0 core.Host
	if rf, ok := ret.Get(0).(func() core.Host); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Host)
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

// LocalData provides a mock function with given fields:
func (_m *Project) LocalData() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// PrimaryTargetName provides a mock function with given fields:
func (_m *Project) PrimaryTargetName() (string, error) {
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

// ResourceId provides a mock function with given fields:
func (_m *Project) ResourceId() (string, error) {
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

// RootPath provides a mock function with given fields:
func (_m *Project) RootPath() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// Target provides a mock function with given fields: name, provider
func (_m *Project) Target(name string, provider string) (core.Target, error) {
	ret := _m.Called(name, provider)

	var r0 core.Target
	if rf, ok := ret.Get(0).(func(string, string) core.Target); ok {
		r0 = rf(name, provider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Target)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, provider)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TargetIds provides a mock function with given fields:
func (_m *Project) TargetIds() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// TargetIndex provides a mock function with given fields:
func (_m *Project) TargetIndex() (core.TargetIndex, error) {
	ret := _m.Called()

	var r0 core.TargetIndex
	if rf, ok := ret.Get(0).(func() core.TargetIndex); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TargetIndex)
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

// TargetNames provides a mock function with given fields:
func (_m *Project) TargetNames() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// Tmp provides a mock function with given fields:
func (_m *Project) Tmp() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

// UI provides a mock function with given fields:
func (_m *Project) UI() (terminal.UI, error) {
	ret := _m.Called()

	var r0 terminal.UI
	if rf, ok := ret.Get(0).(func() terminal.UI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(terminal.UI)
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

// VagrantfileName provides a mock function with given fields:
func (_m *Project) VagrantfileName() (string, error) {
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

// VagrantfilePath provides a mock function with given fields:
func (_m *Project) VagrantfilePath() (path.Path, error) {
	ret := _m.Called()

	var r0 path.Path
	if rf, ok := ret.Get(0).(func() path.Path); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(path.Path)
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

type NewProjectT interface {
	mock.TestingT
	Cleanup(func())
}

// NewProject creates a new instance of Project. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProject(t NewProjectT) *Project {
	mock := &Project{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
