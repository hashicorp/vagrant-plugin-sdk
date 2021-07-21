package mocks

import (
	"reflect"

	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

// AuthenticatorType: (*Authenticator)(nil),
// CommandType:       (*Command)(nil),
// CommunicatorType:  (*Communicator)(nil),
// ConfigType:        (*Config)(nil),
// GuestType:         (*Guest)(nil),
// HostType:          (*Host)(nil),
// LogPlatformType:   (*LogPlatform)(nil),
// LogViewerType:     (*LogViewer)(nil),
// ProviderType:      (*Provider)(nil),
// ProvisionerType:   (*Provisioner)(nil),
// SyncedFolderType:  (*SyncedFolder)(nil),

// ForType returns an implementation of the given type that supports mocking.
func ForType(t component.Type) interface{} {
	// Note that the tests in mocks_test.go verify that we support all types
	switch t {
	case component.AuthenticatorType:
		return &Authenticator{}

	case component.CommandType:
		return &Command{}

	case component.CommunicatorType:
		return &Communicator{}

	case component.ConfigType:
		return &Config{}

	case component.GuestType:
		return &Guest{}

	case component.HostType:
		return &Host{}

	case component.LogPlatformType:
		return &LogPlatform{}

	case component.LogViewerType:
		return &LogViewer{}

	case component.ProviderType:
		return &Provider{}

	case component.ProvisionerType:
		return &Provisioner{}

	case component.SyncedFolderType:
		return &SyncedFolder{}

	default:
		return nil
	}
}

// Mock returns the Mock field for the given interface. The interface value
// should be one of the mocks in this package. This will panic if an incorrect
// value is given, error checking is not done.
func Mock(v interface{}) *mock.Mock {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Interface {
		value = reflect.Indirect(value)
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	field := value.FieldByName("Mock")
	return field.Addr().Interface().(*mock.Mock)
}
