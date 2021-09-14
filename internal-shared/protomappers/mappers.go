package protomappers

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/DavidGamba/go-getoptions/option"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cacher"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/pluginclient"
	plugincomponent "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
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
	Box,
	BoxProto,
	Host,
	HostProto,
	Guest,
	GuestProto,
	CommandInfo,
	CommandInfoProto,
	CommandInfoFromResponse,
	Communicator,
	CommunicatorProto,
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
	MachineState,
	MachineStateProto,
	MapToProto,
	Metadata,
	MetadataProto,
	MetadataSet,
	MetadataSetProto,
	NamedCapability,
	NamedCapabilityProto,
	Project,
	ProjectProto,
	ProtoToMap,
	Provider,
	ProviderProto,
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
	TargetIndex,
	TargetIndexProto,
}

func NamedCapability(
	input *vagrant_plugin_sdk.Args_NamedCapability,
) *component.NamedCapability {
	return &component.NamedCapability{
		Capability: input.Capability,
	}
}

func NamedCapabilityProto(
	input *component.NamedCapability,
) *vagrant_plugin_sdk.Args_NamedCapability {
	return &vagrant_plugin_sdk.Args_NamedCapability{
		Capability: input.Capability,
	}
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

	internal.Logger.Trace("wrapping host plugin",
		"host", input)

	id, ep, err := wrapClient(input, p, internal)
	if err != nil {
		internal.Logger.Warn("failed to wrap host plugin",
			"host", input,
			"error", err)

		return nil, err
	}
	return &vagrant_plugin_sdk.Args_Host{
		Network:  ep.Network(),
		Target:   ep.String(),
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
	internal.Logger.Trace("connecting to wrapped host plugin",
		"connection-info", input)

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		internal.Logger.Warn("failed to connect to wrapped host plugin",
			"connection-info", input,
			"error", err)

		return nil, err
	}

	return client.(core.Host), nil
}

func GuestProto(
	input component.Guest,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Guest, error) {
	p := &plugincomponent.GuestPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    input,
	}

	internal.Logger.Trace("wrapping guest plugin", "guest", input)
	id, ep, err := wrapClient(input, p, internal)
	if err != nil {
		internal.Logger.Warn("failed to wrap guest plugin", "guest", input, "error", err)
		return nil, err
	}
	return &vagrant_plugin_sdk.Args_Guest{
		Network:  ep.Network(),
		Target:   ep.String(),
		StreamId: id,
	}, nil
}

func Guest(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Guest,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Guest, error) {
	p := &plugincomponent.GuestPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}
	internal.Logger.Trace("connecting to wrapped guest plugin", "connection-info", input)
	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		internal.Logger.Warn("failed to connect to wrapped guest plugin", "connection-info", input, "error", err)
		return nil, err
	}

	return client.(core.Guest), nil
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
	dir := datadir.NewBasicDir(input.ConfigDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Basis{Dir: dir}
}

func DatadirProject(input *vagrant_plugin_sdk.Args_DataDir_Project) *datadir.Project {
	dir := datadir.NewBasicDir(input.ConfigDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Project{Dir: dir}
}

func DatadirTarget(input *vagrant_plugin_sdk.Args_DataDir_Target) *datadir.Target {
	dir := datadir.NewBasicDir(input.ConfigDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Target{Dir: dir}
}

func DatadirComponent(input *vagrant_plugin_sdk.Args_DataDir_Project) *datadir.Component {
	dir := datadir.NewBasicDir(input.ConfigDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Component{Dir: dir}
}

func DatadirBasisProto(input *datadir.Basis) *vagrant_plugin_sdk.Args_DataDir_Basis {
	return &vagrant_plugin_sdk.Args_DataDir_Basis{
		CacheDir:  input.CacheDir().String(),
		DataDir:   input.DataDir().String(),
		TempDir:   input.TempDir().String(),
		ConfigDir: input.ConfigDir().String(),
	}
}

func DatadirProjectProto(input *datadir.Project) *vagrant_plugin_sdk.Args_DataDir_Project {
	return &vagrant_plugin_sdk.Args_DataDir_Project{
		CacheDir:  input.CacheDir().String(),
		DataDir:   input.DataDir().String(),
		TempDir:   input.TempDir().String(),
		ConfigDir: input.ConfigDir().String(),
	}
}

func DatadirTargetProto(input *datadir.Target) *vagrant_plugin_sdk.Args_DataDir_Target {
	return &vagrant_plugin_sdk.Args_DataDir_Target{
		CacheDir:  input.CacheDir().String(),
		DataDir:   input.DataDir().String(),
		TempDir:   input.TempDir().String(),
		ConfigDir: input.ConfigDir().String(),
	}
}

func DatadirComponentProto(input *datadir.Component) *vagrant_plugin_sdk.Args_DataDir_Component {
	return &vagrant_plugin_sdk.Args_DataDir_Component{
		CacheDir:  input.CacheDir().String(),
		DataDir:   input.DataDir().String(),
		TempDir:   input.TempDir().String(),
		ConfigDir: input.ConfigDir().String(),
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

	internal.Logger.Trace("connecting to wrapped ui", "stream_id", input.StreamId)
	client, err := wrapConnect(ctx, p, input, internal)

	if err != nil {
		internal.Logger.Warn("failed to connect to wrapped ui", "steam_id", input.StreamId, "error", err)
		return nil, err
	}

	internal.Logger.Trace("connected to wrapped ui", "ui", client, "stream_id", input.StreamId)
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

	internal.Logger.Trace("wrapping ui", "ui", ui)
	id, ep, err := wrapClient(ui, p, internal)

	if err != nil {
		internal.Logger.Trace("failed to wrap ui", "ui", ui, "error", err)
		return nil, err
	}

	internal.Logger.Trace("wrapped ui", "ui", ui, "stream_id", id, "endpoint", ep)
	return &vagrant_plugin_sdk.Args_TerminalUI{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}, nil
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

	id, ep, err := wrapClient(bag, p, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_StateBag{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}, nil
}

func CommandInfoFromResponse(
	input *vagrant_plugin_sdk.Command_CommandInfoResp,
) *vagrant_plugin_sdk.Command_CommandInfo {
	return input.CommandInfo
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

	id, ep, err := wrapClient(b, bp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Basis{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String(),
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

func CommunicatorProto(
	c component.Communicator,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Communicator, error) {
	cp := &plugincomponent.CommunicatorPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    c,
	}

	id, ep, err := wrapClient(c, cp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Communicator{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String(),
	}, nil
}

func Communicator(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Communicator,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Communicator, error) {
	p := &plugincomponent.CommunicatorPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Communicator), nil
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

	id, ep, err := wrapClient(p, pp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Project{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String(),
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

func SyncedFolderProto(
	s component.SyncedFolder,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_SyncedFolder, error) {
	sp := &plugincomponent.SyncedFolderPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    s,
	}

	id, endpoint, err := wrapClient(s, sp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_SyncedFolder{
		StreamId: id,
		Network:  endpoint.Network(),
		Target:   endpoint.String(),
	}, nil
}

func SyncedFolder(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Provider,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.SyncedFolder, error) {
	s := &plugincomponent.SyncedFolderPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, s, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.SyncedFolder), nil
}

func ProviderProto(
	t component.Provider,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Provider, error) {
	tp := &plugincomponent.ProviderPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    t,
	}

	id, endpoint, err := wrapClient(t, tp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Provider{
		StreamId: id,
		Network:  endpoint.Network(),
		Target:   endpoint.String(),
	}, nil
}

func Provider(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Provider,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Provider, error) {
	t := &plugincomponent.ProviderPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, t, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Provider), nil
}

func TargetProto(
	t core.Target,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Target, error) {
	rid, err := t.ResourceId()
	if err != nil {
		return nil, err
	}
	if at := internal.Cache.Get(rid); at != nil {
		log.Warn("using cached target value", "value", at)
		return at.(*vagrant_plugin_sdk.Args_Target), nil
	}

	tp := &plugincore.TargetPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    t,
	}

	id, endpoint, err := wrapClient(t, tp, internal)
	if err != nil {
		return nil, err
	}

	proto := &vagrant_plugin_sdk.Args_Target{
		StreamId: id,
		Network:  endpoint.Network(),
		Target:   endpoint.String(),
	}

	log.Warn("registering target proto to cache",
		"rid", rid,
		"proto", proto,
	)
	internal.Cache.Register(rid, proto)
	return proto, nil
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

	id, ep, err := wrapClient(m, mp, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Target_Machine{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String(),
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

func TargetIndexProto(
	t core.TargetIndex,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_TargetIndex, error) {
	ti := &plugincore.TargetIndexPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
		Impl:    t,
	}

	id, ep, err := wrapClient(t, ti, internal)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_TargetIndex{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String(),
	}, nil
}

func TargetIndex(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_TargetIndex,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.TargetIndex, error) {
	ti := &plugincore.TargetIndexPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, ti, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.TargetIndex), nil
}

type connInfo interface {
	GetStreamId() uint32
	GetNetwork() string
	GetTarget() string
}

type hasTarget interface {
	SetTarget(net.Addr)
	Target() net.Addr
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
	internal.Logger.Trace("connecting to wrapped plugin",
		"plugin", hclog.Fmt("%T", p),
		"connection", i,
		"broker", hclog.Fmt("%p", internal.Broker))

	var err error
	var conn *grpc.ClientConn
	var addr net.Addr
	if target := i.GetTarget(); target != "" {
		switch i.GetNetwork() {
		case "tcp":
			addr, err = net.ResolveTCPAddr("tcp", target)
		case "unix":
			addr, err = net.ResolveUnixAddr("unix", target)
		default:
			return nil, fmt.Errorf(
				"Unknown target address type: %s", i.GetNetwork())
		}

		internal.Logger.Trace("connecting to wrapped plugin via direct target",
			"plugin", hclog.Fmt("%T", p),
			"target", target)

		// TODO(spox): grab dial options from pluginclient config
		conn, err = grpc.Dial("unused", grpc.WithDialer(
			func(_ string, _ time.Duration) (net.Conn, error) {
				return net.Dial(i.GetNetwork(), target)
			}), grpc.WithInsecure(),
		)
	} else {
		internal.Logger.Trace("connecting to wrapped plugin via broker",
			"plugin", hclog.Fmt("%T", p),
			"stream_id", i.GetStreamId(),
			"broker", hclog.Fmt("%p", internal.Broker))

		conn, err = internal.Broker.Dial(i.GetStreamId())
	}
	if err != nil {
		internal.Logger.Warn("failed to connect to wrapped plugin",
			"plugin", hclog.Fmt("%T", p),
			"connection", i,
			"broker", hclog.Fmt("%p", internal.Broker),
			"error", err)

		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	client, err := p.GRPCClient(ctx, internal.Broker, conn)
	if err != nil {
		internal.Logger.Warn("failed to create client for wrapped plugin",
			"plugin", hclog.Fmt("%T", p),
			"connection", i,
			"broker", hclog.Fmt("%p", internal.Broker),
			"error", err)

		return nil, err
	}

	if closer, ok := client.(io.Closer); ok {
		internal.Cleanup.Do(func() { closer.Close() })
	}

	if cache, ok := client.(cacher.HasCache); ok {
		cache.SetCache(internal.Cache)
	}

	internal.Logger.Trace("new client built for wrapped plugin",
		"plugin", hclog.Fmt("%T", p),
		"client", client,
		"connection", i,
		"broker", hclog.Fmt("%p", internal.Broker))

	if addr != nil {
		if ec, ok := client.(hasTarget); ok {
			internal.Logger.Trace("setting direct target on new client",
				"plugin", hclog.Fmt("%T", p),
				"target", addr)

			ec.SetTarget(addr)
		} else {
			internal.Logger.Trace("client does not support direct targets for wrapped plugins",
				"plugin", hclog.Fmt("%T", p),
				"client", hclog.Fmt("%T", client))
		}
	}

	return client, nil
}

// This takes a plugin (which generally uses a client as the plugin implementation)
// and creates a new server for remote connections via the internal broker.
func wrapClient(
	impl interface{},
	p plugin.GRPCPlugin,
	internal *pluginargs.Internal,
) (id uint32, target net.Addr, err error) {
	// If an existing target exists for the implementation, use
	// that value for where to connect
	if iep, ok := impl.(hasTarget); ok {
		if target = iep.Target(); target != nil {
			internal.Logger.Trace("using preset wrapped plugin target",
				"plugin", hclog.Fmt("%T", p),
				"target", target)

			return
		}
	} else {
		internal.Logger.Warn("implementation does not support direct targets for wrapped plugins",
			"plugin", hclog.Fmt("%T", p),
			"implementation", hclog.Fmt("%T", impl),
		)
	}

	// Fetch the next available steam ID from the broker
	id = internal.Broker.NextId()

	// Since we want to register the target endpoint directly for
	// access off the configured broker, we need to get the listener
	// and setup the server directly instead of letting the plugin
	// library handle it for us
	l, err := internal.Broker.Accept(id)
	if err != nil {
		internal.Logger.Warn("failed to establish connection stream",
			"error", err)

		return
	}
	target = l.Addr()

	// Grab the shared plugin configuration so the expected
	// server configuration can be applied
	config := pluginclient.ClientConfig(internal.Logger)
	sopts := []grpc.ServerOption{}
	if config.TLSConfig != nil {
		sopts = append(sopts, grpc.Creds(credentials.NewTLS(config.TLSConfig)))
	}

	internal.Logger.Trace("starting listener for wrapped plugin",
		"broker", hclog.Fmt("%p", internal.Broker),
		"plugin", hclog.Fmt("%T", p),
		"stream_id", id,
		"target", target)

	server := plugin.DefaultGRPCServer(sopts)
	if err = p.GRPCServer(internal.Broker, server); err != nil {
		return
	}

	// Register a shutdown of this wrapped plugin server in our
	// cleanup so we don't leave it hanging around when closed
	internal.Cleanup.Do(func() {
		internal.Logger.Trace("shutting down listener for wrapped plugin",
			"broker", hclog.Fmt("%p", internal.Broker),
			"plugin", hclog.Fmt("%T", p),
			"stream_id", id,
			"target", target)

		server.GracefulStop()
	})

	// Start serving
	go server.Serve(l)

	return
}

func init() {
	for _, fn := range All {
		mFn, err := argmapper.NewFunc(fn)
		if err != nil {
			panic(err)
		}
		plugincore.MapperFns = append(plugincore.MapperFns, mFn)
		plugincomponent.MapperFns = append(plugincomponent.MapperFns, mFn)
		plugincomponent.ProtomapperAllMap[reflect.TypeOf(fn)] = struct{}{}
	}
}
