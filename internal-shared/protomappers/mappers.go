package protomappers

import (
	"context"
	"io"

	//	"strconv"
	//	"time"

	"github.com/DavidGamba/go-getoptions/option"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	plugincore "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/core"
	pluginterminal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/terminal"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"google.golang.org/protobuf/types/known/structpb"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	JobInfo,
	JobInfoProto,
	DatadirBasis,
	DatadirProject,
	DatadirTarget,
	DatadirComponent,
	DatadirBasisProto,
	DatadirProjectProto,
	DatadirTargetProto,
	DatadirComponentProto,
	Logger,
	LoggerProto,
	TerminalUI,
	TerminalUIProto,
	MetadataSet,
	MetadataSetProto,
	StateBag,
	StateBagProto,
	//	Machine,
	Flags,
	FlagsProto,
	MapToProto,
	ProtoToMap,
	CommandInfo,
	CommandInfoProto,
	StateBag,
	StateBagProto,
	//	Project,
}

// TODO(spox): make sure these new mappers actually work
// func Machine(input *vagrant_plugin_sdk.Args_Machine) (*core.Machine, error) {
// 	var result core.Machine
// 	return &result, mapstructure.Decode(input, &result)
// }

// func MachineProto(input *core.Machine) (*vagrant_plugin_sdk.Args_Machine, error) {
// 	var result vagrant_plugin_sdk.Args_Machine
// 	return &result, mapstructure.Decode(intput, &result)
// }

// TODO(spox): end of mappers to validate

// Flags maps
func Flags(input []*vagrant_plugin_sdk.Command_Flag) (opt []*option.Option, err error) {
	opt = []*option.Option{}
	// TODO: add short description as alias
	// https://godoc.org/github.com/DavidGamba/go-getoptions#GetOpt.Alias
	for _, f := range input {
		var newOpt *option.Option
		switch f.Type {
		case vagrant_plugin_sdk.Command_Flag_STRING:
			newOpt = option.New(f.LongName, option.StringType)
		case vagrant_plugin_sdk.Command_Flag_BOOL:
			newOpt = option.New(f.LongName, option.BoolType)
		}
		newOpt.Description = f.Description
		newOpt.DefaultStr = f.DefaultValue
		opt = append(opt, newOpt)
	}
	return opt, err
}

// Flags maps
func FlagsProto(input []*option.Option) (output []*vagrant_plugin_sdk.Command_Flag, err error) {
	output = []*vagrant_plugin_sdk.Command_Flag{}

	for _, f := range input {
		var flagType vagrant_plugin_sdk.Command_Flag_Type
		switch f.OptType {
		case option.StringType:
			flagType = vagrant_plugin_sdk.Command_Flag_STRING
		case option.BoolType:
			flagType = vagrant_plugin_sdk.Command_Flag_BOOL
		}

		// TODO: get aliases
		j := &vagrant_plugin_sdk.Command_Flag{
			LongName:     f.Name,
			ShortName:    f.Name,
			Description:  f.Description,
			DefaultValue: f.DefaultStr,
			Type:         flagType,
		}
		output = append(output, j)
	}
	return output, nil
}

func MapToProto(input map[string]interface{}) (*structpb.Struct, error) {
	return structpb.NewStruct(input)
}

func ProtoToMap(input *structpb.Struct) (map[string]interface{}, error) {
	return input.AsMap(), nil
}

// JobInfo maps Args.JobInfo to component.JobInfo.
func JobInfo(input *vagrant_plugin_sdk.Args_JobInfo) (*component.JobInfo, error) {
	var result component.JobInfo
	return &result, mapstructure.Decode(input, &result)
}

// JobInfoProto
func JobInfoProto(input *component.JobInfo) (*vagrant_plugin_sdk.Args_JobInfo, error) {
	var result vagrant_plugin_sdk.Args_JobInfo
	return &result, mapstructure.Decode(input, &result)
}

func DatadirBasis(input *vagrant_plugin_sdk.Args_DataDir_Basis) *datadir.Basis {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Basis{Dir: dir}
}

func DatadirProject(input *vagrant_plugin_sdk.Args_DataDir_Project) *datadir.Project {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Project{Dir: dir}
}

func DatadirTarget(input *vagrant_plugin_sdk.Args_DataDir_Target) *datadir.Target {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Target{Dir: dir}
}

func DatadirComponent(input *vagrant_plugin_sdk.Args_DataDir_Project) *datadir.Component {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Component{Dir: dir}
}

func DatadirBasisProto(input *datadir.Basis) *vagrant_plugin_sdk.Args_DataDir_Basis {
	return &vagrant_plugin_sdk.Args_DataDir_Basis{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirProjectProto(input *datadir.Project) *vagrant_plugin_sdk.Args_DataDir_Project {
	return &vagrant_plugin_sdk.Args_DataDir_Project{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirTargetProto(input *datadir.Project) *vagrant_plugin_sdk.Args_DataDir_Target {
	return &vagrant_plugin_sdk.Args_DataDir_Target{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirComponentProto(input *datadir.Component) *vagrant_plugin_sdk.Args_DataDir_Component {
	return &vagrant_plugin_sdk.Args_DataDir_Component{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

// Logger maps *vagrant_plugin_sdk.Args_Logger to an hclog.Logger
func Logger(input *vagrant_plugin_sdk.Args_Logger) hclog.Logger {
	// We use the default logger as the base. Within a plugin we always set
	// it so we can confidently use this. This lets plugins potentially mess
	// with this but that's a risk we have to take.
	return hclog.L().ResetNamed(input.Name)
}

func LoggerProto(log hclog.Logger) *vagrant_plugin_sdk.Args_Logger {
	return &vagrant_plugin_sdk.Args_Logger{
		Name: log.Name(),
	}
}

// TerminalUI maps *vagrant_plugin_sdk.Args_TerminalUI to an hclog.TerminalUI
func TerminalUI(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_TerminalUI,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (terminal.UI, error) {
	// Create our plugin
	p := &pluginterminal.UIPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	conn, err := internal.Broker.Dial(input.StreamId)
	if err != nil {
		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	client, err := p.GRPCClient(ctx, internal.Broker, conn)
	if err != nil {
		return nil, err
	}

	// Our UI should implement close since we have to stop streams and
	// such but we gate it here in case we ever change the implementation.
	if closer, ok := client.(io.Closer); ok {
		internal.Cleanup.Do(func() { closer.Close() })
	}

	return client.(terminal.UI), nil
}

func TerminalUIProto(
	ui terminal.UI,
	log hclog.Logger,
	internal *pluginargs.Internal,
) *vagrant_plugin_sdk.Args_TerminalUI {
	// Create our plugin
	p := &pluginterminal.UIPlugin{
		Impl:    ui,
		Mappers: internal.Mappers,
		Logger:  log,
	}

	id := internal.Broker.NextId()

	// Serve it
	go internal.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
		server := plugin.DefaultGRPCServer(opts)
		if err := p.GRPCServer(internal.Broker, server); err != nil {
			panic(err)
		}
		return server
	})

	return &vagrant_plugin_sdk.Args_TerminalUI{StreamId: id}
}

// Machine maps *vagrant_plugin_sdk.Args_Machine to a core.Machine
// func Machine(
// 	ctx context.Context,
// 	input *vagrant_plugin_sdk.Args_Machine,
// 	log hclog.Logger,
// 	internal *pluginargs.Internal,
// ) (*plugincore.MachineClient, error) {
// 	p := &plugincore.MachinePlugin{
// 		Mappers: internal.Mappers,
// 		Logger:  log,
// 	}

// 	id, err := strconv.Atoi(input.ServerAddr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	conn, err := internal.Broker.Dial(uint32(id))
// 	if err != nil {
// 		return nil, err
// 	}
// 	internal.Cleanup.Do(func() { conn.Close() })

// 	client, err := p.GRPCClient(ctx, internal.Broker, conn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if closer, ok := client.(io.Closer); ok {
// 		internal.Cleanup.Do(func() { closer.Close() })
// 	}

// 	return client.(*plugincore.MachineClient).ForResource(input.ResourceId), nil
// }

// Machine maps component.Machine to a *vagrant_plugin_sdk.Args_Machine
// func MachineProto(
// 	machine *plugincore.MachineClient,
// 	log hclog.Logger,
// 	internal *pluginargs.Internal,
// ) (*vagrant_plugin_sdk.Args_Machine, error) {
// 	p := &plugincore.MachinePlugin{
// 		Impl:    machine,
// 		Mappers: internal.Mappers,
// 		Logger:  log}

// 	id := internal.Broker.NextId()

// 	go internal.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
// 		server := plugin.DefaultGRPCServer(opts)
// 		if err := p.GRPCServer(internal.Broker, server); err != nil {
// 			panic(err)
// 		}
// 		return server
// 	})

// 	return &vagrant_plugin_sdk.Args_Machine{
// 		ResourceId: machine.ResourceID,
// 		ServerAddr: strconv.Itoa(int(id)),
// 	}, nil
// }

// Project maps *vagrant_plugin_sdk.Args_Project to a core.Project
// func Project(
// 	ctx context.Context,
// 	input *vagrant_plugin_sdk.Args_Project,
// 	log hclog.Logger,
// 	internal *pluginargs.Internal,
// ) (*plugincore.Project, error) {
// 	p := &plugincore.ProjectPlugin{
// 		Mappers: internal.Mappers,
// 		Logger:  log,
// 	}

// 	timeout := 5 * time.Second
// 	// Create a new cancellation context so we can cancel in the case of an error
// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	// Connect to the local server
// 	conn, err := grpc.DialContext(ctx, input.ServerAddr,
// 		grpc.WithBlock(),
// 		grpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	internal.Cleanup.Do(func() { conn.Close() })

// 	mc, err := p.GRPCClient(ctx, internal.Broker, conn)
// 	projectClient, ok := mc.(*plugincore.ProjectClient)
// 	if !ok {
// 		panic("failed to create machine client")
// 	}

// 	// TODO: decode input project into a project
// 	result := plugincore.NewProject(projectClient)
// 	mapstructure.Decode(input, &result)

// 	return result, nil
// }

func MetadataSet(input *vagrant_plugin_sdk.Args_MetadataSet) *component.MetadataSet {
	return &component.MetadataSet{
		Metadata: input.Metadata,
	}
}

func MetadataSetProto(meta *component.MetadataSet) *vagrant_plugin_sdk.Args_MetadataSet {
	return &vagrant_plugin_sdk.Args_MetadataSet{Metadata: meta.Metadata}
}

// StateBag maps StateBag proto to core.StateBag.
func StateBag(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_StateBag,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.StateBag, error) {
	// Create our plugin
	p := &plugincore.StateBagPlugin{
		Logger: log,
	}

	conn, err := internal.Broker.Dial(input.StreamId)
	if err != nil {
		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	client, err := p.GRPCClient(ctx, internal.Broker, conn)
	if err != nil {
		return nil, err
	}

	if closer, ok := client.(io.Closer); ok {
		internal.Cleanup.Do(func() { closer.Close() })
	}

	return client.(core.StateBag), nil
}

func StateBagProto(
	bag core.StateBag,
	log hclog.Logger,
	internal *pluginargs.Internal,
) *vagrant_plugin_sdk.Args_StateBag {
	// Create our plugin
	p := &plugincore.StateBagPlugin{
		Impl:   bag,
		Logger: log,
	}

	id := internal.Broker.NextId()

	go internal.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
		server := plugin.DefaultGRPCServer(opts)
		if err := p.GRPCServer(internal.Broker, server); err != nil {
			panic(err)
		}
		return server
	})

	return &vagrant_plugin_sdk.Args_StateBag{StreamId: id}
}

func CommandInfo(input *vagrant_plugin_sdk.Command_CommandInfo) (*component.CommandInfo, error) {
	flags, err := Flags(input.Flags)

	subcommands := []*component.CommandInfo{}
	for _, cmd := range input.Subcommands {
		subcommand, err := CommandInfo(cmd)
		if err != nil {
			return nil, err
		}
		subcommands = append(subcommands, subcommand)
	}

	result := &component.CommandInfo{
		Flags:       flags,
		Name:        input.Name,
		Help:        input.Help,
		Synopsis:    input.Synopsis,
		Subcommands: subcommands,
	}
	return result, err
}

func CommandInfoProto(input *component.CommandInfo) (*vagrant_plugin_sdk.Command_CommandInfo, error) {
	var result vagrant_plugin_sdk.Command_CommandInfo
	err := mapstructure.Decode(input, &result)
	if err != nil {
		return nil, err
	}
	result.Flags, err = FlagsProto(input.Flags)
	subcmds := []*vagrant_plugin_sdk.Command_CommandInfo{}
	for _, cmd := range input.Subcommands {
		toAdd, err := CommandInfoProto(cmd)
		if err != nil {
			return nil, err
		}
		subcmds = append(subcmds, toAdd)
	}
	result.Subcommands = subcmds
	return &result, err
}
