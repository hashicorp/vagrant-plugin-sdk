// Package component has the interfaces for all the components that
// can be implemented. A component is the broad term used to describe
// all providers, provisioners, etc.
//
// Many component interfaces have functions named `XFunc` where "X" is some
// operation and the return value is "interface{}". These functions should return
// a method handle to the function implementing that operation. This pattern is
// done so that we can support custom typed operations that take and return
// full rich types for an operation. We use a minimal dependency-injection
// framework (see internal/mapper) to call these functions.
package component

import (
	"fmt"
	"strings"
)

// Type is an enum of all the types of components supported.
// This isn't used directly in this package but is used by other packages
// to reference the component types.
type Type uint

const (
	InvalidType       Type = iota // Invalid
	CommandType                   // Command
	CommunicatorType              // Communicator
	GuestType                     // Guest
	HostType                      // Host
	ProviderType                  // Provider
	ProvisionerType               // Provisioner
	SyncedFolderType              // SyncedFolder
	AuthenticatorType             // Authenticator
	LogPlatformType               // LogPlatform
	LogViewerType                 // LogViewer
	MapperType                    // Mapper
	ConfigType                    // Config
	PluginInfoType                // PluginInfo
	PushType                      // Push
	maxType
)

// TypeMap is a mapping of Type to the nil pointer to the interface of that
// type. This can be used with libraries such as mapper.
var TypeMap = map[Type]interface{}{
	AuthenticatorType: (*Authenticator)(nil),
	CommandType:       (*Command)(nil),
	CommunicatorType:  (*Communicator)(nil),
	ConfigType:        (*Config)(nil),
	GuestType:         (*Guest)(nil),
	HostType:          (*Host)(nil),
	LogPlatformType:   (*LogPlatform)(nil),
	LogViewerType:     (*LogViewer)(nil),
	PluginInfoType:    (*PluginInfo)(nil),
	ProviderType:      (*Provider)(nil),
	ProvisionerType:   (*Provisioner)(nil),
	SyncedFolderType:  (*SyncedFolder)(nil),
	PushType:          (*Push)(nil),
}

func FindComponent(name string) (interface{}, error) {
	for k, v := range TypeMap {
		if k.String() == name ||
			strings.ToLower(k.String()) == name ||
			strings.ToLower(k.String()) == strings.ReplaceAll(name, "_", "") {
			return v, nil
		}
	}
	return nil, fmt.Errorf("failed to find component for name '%s'", name)
}

func FindType(name string) (Type, error) {
	for k, _ := range TypeMap {
		if k.String() == name ||
			strings.ToLower(k.String()) == name ||
			strings.ToLower(k.String()) == strings.ReplaceAll(name, "_", "") {
			return k, nil
		}
	}
	return maxType, fmt.Errorf("failed to find component for name '%s'", name)
}

type PluginInfo interface {
	ComponentTypes() []Type
	Name() string
}

type CommandFlags []*CommandFlag

func (c CommandFlags) Display() string {
	var pad int
	opts := make([]string, len(c))
	desc := make([]string, len(c))

	for i, f := range c {
		if f.ShortName != "" {
			opts[i] = "-" + f.ShortName + ","
		} else {
			opts[i] = "   "
		}
		if f.Type == FlagBool {
			opts[i] = opts[i] + " --[no-]" + f.LongName
		} else {
			opts[i] = opts[i] + " --" + f.LongName + " VALUE"
		}
		desc[i] = f.Description
		if len(opts[i]) > pad {
			pad = len(opts[i])
		}
	}

	pad += 11
	var d string

	for i := 0; i < len(opts); i++ {
		d = d + fmt.Sprintf("%4s%-*s%s\n", "", pad, opts[i], desc[i])
	}

	return d
}

type FlagType uint

const (
	FlagString FlagType = 1 << iota // String
	FlagBool                        // Bool
)

type CommandFlag struct {
	LongName     string
	ShortName    string
	Description  string
	DefaultValue string
	Type         FlagType
}

type CommandInfo struct {
	Name        string
	Help        string
	Synopsis    string
	Flags       []*CommandFlag
	Subcommands []*CommandInfo
}

type CommandParams struct {
	Flags     map[string]interface{}
	Arguments []string
}

type Direct struct {
	Arguments []interface{}
}

type Command interface {
	// Execute a command
	ExecuteFunc([]string) interface{}
	// Retruns command info
	CommandInfoFunc() interface{}
}

type Config interface {
}

type Communicator interface {
	// Checks if machine can be used with communicator
	MatchFunc() interface{}
	// Initialize communicator with machine
	InitFunc() interface{}
	// Check if communicator is ready
	ReadyFunc() interface{}
	// Wait for communicator to become ready for given seconds
	WaitForReadyFunc() interface{}
	// Download file from guest path to local path
	DownloadFunc() interface{}
	// Upload file from local path to guest path
	UploadFunc() interface{}
	// Run command
	ExecuteFunc() interface{}
	// Run privileged command
	PrivilegedExecuteFunc() interface{}
	// Run a test command on the guest
	TestFunc() interface{}
	// Reset the communicator. Close and re-establish connection where required.
	ResetFunc() interface{}
}

type CapabilityPlatform interface {
	// Test if capability is available
	HasCapabilityFunc() interface{}
	// Run a capability
	CapabilityFunc(capName string) interface{}
}

type Guest interface {
	// Detect if machine is supported guest
	GuestDetectFunc() interface{}
	// List of parent host names
	ParentFunc() interface{}
	// Test if capability is available
	HasCapabilityFunc() interface{}
	// Run a capability
	CapabilityFunc(capName string) interface{}
}

type Host interface {
	// Detect if machine is supported host
	HostDetectFunc() interface{}
	// List of parent host names
	ParentFunc() interface{}
	// Test if capability is available
	HasCapabilityFunc() interface{}
	// Run a capability
	CapabilityFunc(capName string) interface{}
}

type Provider interface {
	// Check if provider is usable
	UsableFunc() interface{}
	// Check if the provider is installed
	InstalledFunc() interface{}
	// Run an action by name
	ActionFunc(actionName string) interface{}
	// Called when the machine id is changed
	MachineIdChangedFunc() interface{}
	// Get SSH info
	SshInfoFunc() interface{}
	// Get target state
	StateFunc() interface{}

	// Test if capability is available
	HasCapabilityFunc() interface{}
	// Run a capability
	CapabilityFunc(capName string) interface{}
}

type Provisioner interface {
}

type Push interface {
	// Executes a named push strategy
	PushFunc() interface{}
}

type SyncedFolder interface {
	// Determines if an implementation is usable
	UsableFunc() interface{}
	// Called after the machine is booted and networks are setup
	// Adds folders without removing any existing ones
	EnableFunc() interface{}
	// Removes folders from a running machine
	DisableFunc() interface{}
	// Called after destroying a machine
	CleanupFunc() interface{}

	// Test if capability is available
	HasCapabilityFunc() interface{}
	// Run a capability
	CapabilityFunc(capName string) interface{}
}

type MetadataSet struct {
	Metadata map[string]string
}

// Authenticator is responsible for authenticating different types of plugins.
type Authenticator interface {
	// AuthFunc should return the method for getting credentials for a
	// plugin. This should return AuthResult.
	AuthFunc() interface{}

	// ValidateAuthFunc should return the method for validating authentication
	// credentials for the plugin
	ValidateAuthFunc() interface{}
}

// JobInfo is available to plugins to get information about the context
// in which a job is executing.
type JobInfo struct {
	// Id is the ID of the job that is executing this plugin operation.
	// If this is empty then it means that the execution is happening
	// outside of a job.
	Id string

	// Local is true if the operation is running locally on a machine
	// alongside the invocation. This can be used to determine if you can
	// do things such as open browser windows, read user files, etc.
	Local bool

	// Workspace is the workspace that this job is executing in. This should
	// be used by plugins to properly isolate resources from each other.
	// TODO(spox): this needs to be removed
	Workspace string
}

// AuthResult is the return value expected from Authenticator.AuthFunc.
type AuthResult struct {
	// Authenticated when true means that the plugin should now be authenticated
	// (given the other fields in this struct). If ValidateAuth is called,
	// it should succeed. If this is false, the auth method may have printed
	// help text or some other information, but it didn't authenticate. However,
	// this is not an error.
	Authenticated bool
}

type NamedCapability struct {
	Capability string
}
