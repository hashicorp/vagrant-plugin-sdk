package plugin

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// Not secret, just to avoid plugins being launched manually. The
	// cookie value is a random SHA256 via /dev/urandom. This cookie value
	// must NEVER be changed or plugins will stop working.
	MagicCookieKey:   "VAGRANT_PLUGIN",
	MagicCookieValue: "1cf2a7e8cbbd6cec9ec78b952860dc65a7a9eae433b815bea3257bff2257b3a7",
}
var MapperFns []*argmapper.Func

// Plugins returns the list of available plugins and initializes them with
// the given components. This will panic if an invalid component is given.
func Plugins(opts ...Option) map[int]plugin.PluginSet {
	var c pluginConfig
	for _, opt := range opts {
		opt(&c)
	}

	// If we have no logger, we use the default
	if c.Logger == nil {
		c.Logger = hclog.L()
	}

	bp := &BasePlugin{
		Mappers: c.Mappers,
		Logger:  c.Logger,
		Wrapped: false,
	}

	// Save this so we can update it before we finish
	info := &PluginInfoPlugin{BasePlugin: bp}

	// Build our plugin types
	result := map[int]plugin.PluginSet{
		1: {
			"command":      &CommandPlugin{BasePlugin: bp.Clone()},
			"communicator": &CommunicatorPlugin{BasePlugin: bp.Clone()},
			"config":       &ConfigPlugin{},
			"guest":        &GuestPlugin{BasePlugin: bp.Clone()},
			"host":         &HostPlugin{BasePlugin: bp.Clone()},
			"mapper":       &MapperPlugin{BasePlugin: bp.Clone()},
			"plugininfo":   info,
			"provider":     &ProviderPlugin{BasePlugin: bp.Clone()},
			"push":         &PushPlugin{BasePlugin: bp.Clone()},
			"provisioner":  &ProvisionerPlugin{BasePlugin: bp.Clone()},
			"syncedfolder": &SyncedFolderPlugin{BasePlugin: bp.Clone()},
		},
	}

	t := []component.Type{}

	// Set the various field values
	for _, c := range c.Components {
		for typ, ptr := range component.TypeMap {
			pTyp := reflect.TypeOf(ptr)
			cTyp := reflect.TypeOf(c)
			if cTyp.Implements(pTyp.Elem()) {
				t = append(t, typ)
			}
		}
		if err := setFieldValue(result, c); err != nil {
			panic(err)
		}
	}

	// Set plugin info before we finish
	info.Impl = &pluginInfo{types: t, name: c.Name, priority: c.Priority}

	return result
}

// pluginConfig is used to configure Plugins via Option calls.
type pluginConfig struct {
	Name       string
	Priority   int
	Components []interface{}
	Mappers    []*argmapper.Func
	Logger     hclog.Logger
}

// Option configures Plugins
type Option func(*pluginConfig)

// WithComponents sets the components to configure for the plugins.
// This will append to the components.
func WithComponents(cs ...interface{}) Option {
	return func(c *pluginConfig) { c.Components = append(c.Components, cs...) }
}

// WithMappers sets the mappers to configure for the plugins. This will
// append to the existing mappers.
func WithMappers(ms ...*argmapper.Func) Option {
	return func(c *pluginConfig) {
		c.Mappers = append(c.Mappers, ms...)
	}
}

// WithLogger sets the logger for the plugins.
func WithLogger(log hclog.Logger) Option {
	return func(c *pluginConfig) { c.Logger = log }
}

func WithName(n string) Option {
	return func(c *pluginConfig) { c.Name = n }
}

func WithPriority(n int) Option {
	return func(c *pluginConfig) { c.Priority = n }
}

// setFieldValue sets the given value c on any exported field of an available
// plugin that matches the type of c. An error is returned if c can't be
// assigned to ANY plugin type.
//
// preconditions:
//   - plugins in m are pointers to structs
func setFieldValue(m map[int]plugin.PluginSet, c interface{}) error {
	cv := reflect.ValueOf(c)
	ct := cv.Type()

	// Go through each pluginset
	once := false
	for _, set := range m {
		// Go through each plugin
		for _, p := range set {
			// Get the value, dereferencing the pointer. We expect
			// the value to be &SomeStruct{} so we must deref once.
			v := reflect.ValueOf(p).Elem()

			// Go through all the fields
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)

				// If the field is valid and our component can be assigned
				// to it then we set the value directly. We continue setting
				// values because some values we set are available in multiple
				// plugins (loggers for example).
				if f.IsValid() && ct.AssignableTo(f.Type()) {
					f.Set(cv)
					once = true
				}
			}
		}
	}

	if !once {
		return fmt.Errorf("no plugin available for setting field of type %T", c)
	}

	return nil
}

func directSetFieldValue(dst interface{}, c interface{}) error {
	cv := reflect.ValueOf(c)
	ct := cv.Type()

	// Get the value, dereferencing the pointer. We expect
	// the value to be &SomeStruct{} so we must deref once.
	v := reflect.ValueOf(dst).Elem()

	// Go through all the fields
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		// If the field is valid and our component can be assigned
		// to it then we set the value directly. We continue setting
		// values because some values we set are available in multiple
		// plugins (loggers for example).
		if f.IsValid() && ct.AssignableTo(f.Type()) {
			f.Set(cv)
		}
	}

	return nil
}
