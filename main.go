package sdk

import (
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mattn/go-colorable"

	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/protomappers"
	sdkplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/plugincomponent"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/stdio"
)

//go:generate sh -c "protoc -I`go list -m -f \"{{.Dir}}\" github.com/mitchellh/protostructure` -I./3rdparty/proto/api-common-protos -Iproto --go-grpc_opt=module=github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk --go-grpc_out=require_unimplemented_servers=false:proto/vagrant_plugin_sdk/ --go_opt=module=github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk --go_out=proto/vagrant_plugin_sdk/ proto/vagrant_plugin_sdk/*.proto"
//go:generate sh -c "grpc_tools_ruby_protoc -I`go list -m -f \"{{.Dir}}\" github.com/mitchellh/protostructure` -I./3rdparty/proto/api-common-protos -Iproto --grpc_out=ruby-proto/ --ruby_out=ruby-proto/ ./proto/vagrant_plugin_sdk/*.proto protostructure.proto"

// Main is the primary entrypoint for plugins serving components. This
// function never returns; it blocks until the program is exited. This should
// be called immediately in main() in your plugin binaries, no prior setup
// should be done.
func Main(opts ...Option) {

	var c config

	// Default our mappers
	c.Mappers = append(c.Mappers, protomappers.All...)

	// Build config
	for _, opt := range opts {
		opt(&c)
	}

	// We have to rewrite the fatih/color package output/error writers
	// to be our plugin stdout/stderr. We use the color package a lot in
	// our UI and this causes the UI to work.
	color.Output = colorable.NewColorable(stdio.Stdout())
	color.Error = colorable.NewColorable(stdio.Stderr())

	// Create our logger. We also set this as the default logger in case
	// any other libraries are using hclog and our plugin doesn't properly
	// chain it along.
	log := c.Log
	if log == nil {
		log = hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Level:  hclog.Trace,
			Output: os.Stderr,
			Color:  hclog.AutoColor,

			// Critical that this is JSON-formatted. Since we're a plugin this
			// will enable the host to parse our logs and output them in a
			// structured way.
			JSONFormat: true,
		})
		hclog.SetDefault(log)
	}

	// Build up our mappers
	var mappers []*argmapper.Func
	for _, raw := range c.Mappers {
		// If the mapper is already a argmapper.Func, then we let that through as-is
		m, ok := raw.(*argmapper.Func)
		if !ok {
			var err error
			m, err = argmapper.NewFunc(raw, argmapper.Logger(plugincomponent.ArgmapperLogger))
			if err != nil {
				panic(err)
			}
		}

		mappers = append(mappers, m)
	}

	// Serve
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdkplugin.Handshake,
		VersionedPlugins: sdkplugin.Plugins(
			sdkplugin.WithComponents(c.Components...),
			sdkplugin.WithMappers(mappers...),
			sdkplugin.WithLogger(log),
		),
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     log,
		Test:       c.InProcess,
	})
}

// config is the configuration for Main. This can only be modified using
// Option implementations.
type config struct {
	// Components is the list of components to serve from the plugin.
	Components []interface{}

	// Mappers is the list of mapper functions.
	Mappers []interface{}

	InProcess *plugin.ServeTestConfig

	Log hclog.Logger
}

// Option modifies config. Zero or more can be passed to Main.
type Option func(*config)

func InProcess(tc *plugin.ServeTestConfig) Option {
	return func(c *config) { c.InProcess = tc }
}

func WithLogger(l hclog.Logger) Option {
	return func(c *config) { c.Log = l }
}

// WithComponents specifies a list of components to serve from the plugin
// binary. This will append to the list of components to serve. You can
// currently only serve at most one of each type of plugin.
func WithComponents(cs ...interface{}) Option {
	return func(c *config) { c.Components = append(c.Components, cs...) }
}

// WithMappers specifies a list of mappers to apply to the plugin.
//
// Mappers are functions that take zero or more arguments and return
// one result (optionally with an error). These can be used to convert argument
// types as needed for your plugin functions. For example, you can convert a
// proto type to a richer Go struct.
//
// Mappers must take zero or more arguments and return exactly one or two
// values where the second return type must be an error. Example:
//
//   func() *Value
//   func() (*Value, error)
//   -- the above with any arguments
//
// This will append the mappers to the list of available mappers. A set of
// default mappers is always included to convert from SDK proto types to
// richer Go structs.
func WithMappers(ms ...interface{}) Option {
	return func(c *config) { c.Mappers = append(c.Mappers, ms...) }
}
