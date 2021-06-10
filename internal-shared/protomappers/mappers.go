package protomappers

import (
	"context"
	"io"

	"github.com/DavidGamba/go-getoptions/option"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	plugincomponent "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/component"
	plugincore "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/core"
	pluginterminal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/terminal"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	Basis,
	BasisProto,
	CommandInfo,
	CommandInfoProto,
	DatadirBasis,
	DatadirBasisProto,
	DatadirProject,
	DatadirProjectProto,
	DatadirTarget,
	DatadirTargetProto,
	DatadirComponent,
	DatadirComponentProto,
	Flags,
	FlagsProto,
	JobInfo,
	JobInfoProto,
	Logger,
	LoggerProto,
	MapToProto,
	Metadata,
	MetadataProto,
	MetadataSet,
	MetadataSetProto,
	Project,
	ProjectProto,
	ProtoToMap,
	State,
	StateProto,
	StateBag,
	StateBagProto,
	Target,
	TargetProto,
	TargetMachine,
	TargetMachineProto,
	TerminalUI,
	TerminalUIProto,
	MachineState,
	MachineStateProto,
	Box,
	BoxProto,
	Host,
	HostProto,
}

func HostProto(
	input component.Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Host, error) {
	p := &plugincomponent.HostPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    input,
	}

	id, err := wrapClient(p, internal)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Args_Host{
		StreamId: id,
	}, nil
}

func Host(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Host, error) {
	p := &plugincomponent.HostPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}
	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Host), nil
}

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

func MachineState(input *vagrant_plugin_sdk.Args_Target_Machine_State) (*core.MachineState, error) {
	var result core.MachineState
	return &result, mapstructure.Decode(input, &result)
}

func MachineStateProto(input *core.MachineState) (*vagrant_plugin_sdk.Args_Target_Machine_State, error) {
	var result vagrant_plugin_sdk.Args_Target_Machine_State
	return &result, mapstructure.Decode(input, &result)
}

func Box(input *vagrant_plugin_sdk.Args_Target_Machine_Box) (*core.Box, error) {
	var result core.Box
	return &result, mapstructure.Decode(input, &result)
}

func BoxProto(input *core.Box) (*vagrant_plugin_sdk.Args_Target_Machine_Box, error) {
	var result vagrant_plugin_sdk.Args_Target_Machine_Box
	return &result, mapstructure.Decode(input, &result)
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

func DatadirTargetProto(input *datadir.Target) *vagrant_plugin_sdk.Args_DataDir_Target {
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

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}
	return client.(terminal.UI), nil
}

func TerminalUIProto(
	ui terminal.UI,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_TerminalUI, error) {
	// Create our plugin
	p := &pluginterminal.UIPlugin{
		Impl:    ui,
		Mappers: internal.Mappers,
		Logger:  log,
	}

	id, err := wrapClient(p, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_TerminalUI{
		StreamId: id}, nil
}

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

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.StateBag), nil
}

func StateBagProto(
	bag core.StateBag,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_StateBag, error) {
	// Create our plugin
	p := &plugincore.StateBagPlugin{
		Impl:   bag,
		Logger: log,
	}

	id, err := wrapClient(p, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_StateBag{
		StreamId: id}, nil
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

func StateProto(s core.State) *vagrant_plugin_sdk.Args_Target_State {
	var state vagrant_plugin_sdk.Args_Target_State_State
	switch s {
	case core.CREATED:
		state = vagrant_plugin_sdk.Args_Target_State_CREATED
	case core.DESTROYED:
		state = vagrant_plugin_sdk.Args_Target_State_DESTROYED
	case core.PENDING:
		state = vagrant_plugin_sdk.Args_Target_State_PENDING
	default:
		state = vagrant_plugin_sdk.Args_Target_State_UNKNOWN
	}
	return &vagrant_plugin_sdk.Args_Target_State{
		State: state,
	}
}

func State(s *vagrant_plugin_sdk.Args_Target_State) (state core.State) {
	switch s.State {
	case vagrant_plugin_sdk.Args_Target_State_CREATED:
		state = core.CREATED
	case vagrant_plugin_sdk.Args_Target_State_DESTROYED:
		state = core.DESTROYED
	case vagrant_plugin_sdk.Args_Target_State_PENDING:
		state = core.PENDING
	default:
		state = core.UNKNOWN
	}
	return
}

func MetadataProto(m map[string]string) *vagrant_plugin_sdk.Args_MetadataSet {
	return &vagrant_plugin_sdk.Args_MetadataSet{
		Metadata: m,
	}
}

func Metadata(m *vagrant_plugin_sdk.Args_MetadataSet) map[string]string {
	return m.Metadata
}

func BasisProto(
	b core.Basis,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Basis, error) {
	bp := &plugincore.BasisPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    b,
	}

	id, err := wrapClient(bp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Basis{
		StreamId: id,
	}, nil
}

func Basis(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Basis,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Basis, error) {
	b := &plugincore.BasisPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, b, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Basis), nil
}

func ProjectProto(
	p core.Project,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Project, error) {
	pp := &plugincore.ProjectPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    p,
	}

	id, err := wrapClient(pp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Project{
		StreamId: id,
	}, nil
}

func Project(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Project,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Project, error) {
	p := &plugincore.ProjectPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Project), nil
}

func TargetProto(
	t core.Target,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Target, error) {
	tp := &plugincore.TargetPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    t,
	}

	id, err := wrapClient(tp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Target{
		StreamId: id,
	}, nil
}

func Target(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Target,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Target, error) {
	t := &plugincore.TargetPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, t, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Target), nil
}

func TargetMachineProto(
	m core.Machine,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Target_Machine, error) {
	mp := &plugincore.TargetMachinePlugin{
		Mappers:    internal.Mappers,
		Logger:     log,
		Impl:       m,
		TargetImpl: m,
	}

	id, err := wrapClient(mp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Target_Machine{
		StreamId: id,
	}, nil
}

func TargetMachine(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Target_Machine,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Machine, error) {
	m := &plugincore.TargetMachinePlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, m, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Machine), nil
}

type connInfo interface {
	GetStreamId() uint32
}

// When a core plugin is received, the proto will match the
// ConnInfo interface which provides the information needed
// setup a new client. Depending on the origin of the proto
// the client will either establish a direct connection to
// the service, or will connect via the broker.
func wrapConnect(
	ctx context.Context,
	p plugin.GRPCPlugin,
	i connInfo,
	internal *pluginargs.Internal,
) (interface{}, error) {
	conn, err := internal.Broker.Dial(i.GetStreamId())
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

	return client, nil
}

// This takes a plugin (which generally uses a client as the plugin implementation)
// and creates a new server for remote connections via the internal broker.
func wrapClient(p plugin.GRPCPlugin, internal *pluginargs.Internal) (uint32, error) {
	id := internal.Broker.NextId()
	errChan := make(chan error, 1)

	go internal.Broker.AcceptAndServe(id, func(opts []grpc.ServerOption) *grpc.Server {
		server := plugin.DefaultGRPCServer(opts)
		if err := p.GRPCServer(internal.Broker, server); err != nil {
			errChan <- err
			return nil
		}
		internal.Cleanup.Do(func() { server.GracefulStop() })
		close(errChan)
		return server
	})

	err := <-errChan
	if err != nil {
		return 0, err
	}

	return id, nil
}

func init() {
	for _, fn := range All {
		mFn, err := argmapper.NewFunc(fn)
		if err != nil {
			panic(err)
		}
		plugincore.MapperFns = append(plugincore.MapperFns, mFn)
	}
}
