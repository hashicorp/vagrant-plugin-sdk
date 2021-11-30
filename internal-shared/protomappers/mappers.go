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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cacher"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/pluginclient"
	plugincomponent "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	plugincore "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/core"
	pluginterminal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/terminal"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

var WellKnownTypes = []interface{}{
	Boolean,
	BooleanProto,
	Bytes,
	BytesProto,
	Double,
	DoubleProto,
	Float,
	FloatProto,
	Int32,
	Int32Proto,
	Int64,
	Int64Proto,
	String,
	StringProto,
	Timestamp,
	TimestampProto,
	UInt32,
	UInt32Proto,
	UInt64,
	UInt64Proto,
	ValueToBool,
	ValueToList,
	ValueToNull,
	ValueToNumber,
	ValueToString,
	ValueToStruct,
	ValueToString,
	// ValueProto,
}

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	Array,
	ArrayProto,
	Basis,
	BasisProto,
	Box,
	BoxProto,
	BoxCollection,
	BoxCollectionProto,
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
	Direct,
	DirectProto,
	Flags,
	FlagsProto,
	Hash,
	HashProto,
	JobInfo,
	JobInfoProto,
	Logger,
	LoggerProto,
	MachineProject,
	MachineState,
	MachineStateProto,
	MapToProto,
	Metadata,
	MetadataProto,
	MetadataSet,
	MetadataSetProto,
	NamedCapability,
	NamedCapabilityProto,
	Path,
	PathProto,
	Project,
	ProjectProto,
	ProtoToMap,
	Provider,
	ProviderProto,
	Seeds,
	SeedsProto,
	State,
	StateProto,
	StateBag,
	StateBagProto,
	SyncedFolder,
	SyncedFolderProto,
	VagrantfileSyncedFolderToFolder,
	FolderToVagrantfileSyncedFolder,
	Target,
	TargetProto,
	TargetIndex,
	TargetIndexProto,
	TargetMachine,
	TargetMachineProto,
	TargetProject,
	TerminalUI,
	TerminalUIProto,
}

// Known type mappers

func Boolean(
	input *wrapperspb.BoolValue,
) bool {
	return input.Value
}

func BooleanProto(
	input bool,
) *wrapperspb.BoolValue {
	return &wrapperspb.BoolValue{
		Value: input,
	}
}

func Bytes(
	input *wrapperspb.BytesValue,
) []byte {
	return input.Value
}

func BytesProto(
	input []byte,
) *wrapperspb.BytesValue {
	return &wrapperspb.BytesValue{
		Value: input,
	}
}

func Double(
	input *wrapperspb.DoubleValue,
) float64 {
	return input.Value
}

func DoubleProto(
	input float64,
) *wrapperspb.DoubleValue {
	return &wrapperspb.DoubleValue{
		Value: input,
	}
}

func Float(
	input *wrapperspb.FloatValue,
) float32 {
	return input.Value
}

func FloatProto(
	input float32,
) *wrapperspb.FloatValue {
	return &wrapperspb.FloatValue{
		Value: input,
	}
}

func Int32(
	input *wrapperspb.Int32Value,
) int32 {
	return input.Value
}

func Int32Proto(
	input int32,
) *wrapperspb.Int32Value {
	return &wrapperspb.Int32Value{
		Value: input,
	}
}

func Int64(
	input *wrapperspb.Int64Value,
) int64 {
	return input.Value
}

func Int64Proto(
	input int64,
) *wrapperspb.Int64Value {
	return &wrapperspb.Int64Value{
		Value: input,
	}
}

func String(
	input *wrapperspb.StringValue,
) string {
	return input.Value
}

func StringProto(
	input string,
) *wrapperspb.StringValue {
	return &wrapperspb.StringValue{
		Value: input,
	}
}

func Timestamp(
	input *timestamppb.Timestamp,
) time.Time {
	return input.AsTime()
}

func TimestampProto(
	input time.Time,
) *timestamppb.Timestamp {
	return timestamppb.New(input)
}

func UInt32(
	input *wrapperspb.UInt32Value,
) uint32 {
	return input.Value
}

func UInt32Proto(
	input uint32,
) *wrapperspb.UInt32Value {
	return &wrapperspb.UInt32Value{
		Value: input,
	}
}

func UInt64(
	input *wrapperspb.UInt64Value,
) uint64 {
	return input.Value
}

func UInt64Proto(
	input uint64,
) *wrapperspb.UInt64Value {
	return &wrapperspb.UInt64Value{
		Value: input,
	}
}

func ValueToBool(
	input *structpb.Value,
) (bool, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_BoolValue)(nil)) {
		return false, fmt.Errorf("value is not bool kind")
	}

	return input.GetBoolValue(), nil
}

func ValueToList(
	input *structpb.Value,
) ([]*structpb.Value, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_ListValue)(nil)) {
		return nil, fmt.Errorf("value is not list kind")
	}

	return input.GetListValue().Values, nil
}

func ValueToNull(
	input *structpb.Value,
) (interface{}, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_NullValue)(nil)) {
		return nil, fmt.Errorf("value is not null kind")
	}

	return nil, nil
}

func ValueToNumber(
	input *structpb.Value,
) (float64, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_NumberValue)(nil)) {
		return 0, fmt.Errorf("value is not number kind")
	}

	return input.GetNumberValue(), nil
}

func ValueToString(
	input *structpb.Value,
) (string, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_StringValue)(nil)) {
		return "", fmt.Errorf("value is not string kind")
	}
	return input.GetStringValue(), nil
}

func ValueToStruct(
	input *structpb.Value,
) (map[string]*structpb.Value, error) {
	if reflect.TypeOf(input.Kind) != reflect.TypeOf((*structpb.Value_StructValue)(nil)) {
		return nil, fmt.Errorf("value is not struct kind")
	}

	return input.GetStructValue().GetFields(), nil
}

// Custom mappers

func Seeds(
	input *vagrant_plugin_sdk.Args_Seeds,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*core.Seeds, error) {
	result := core.NewSeeds()
	t := make([]interface{}, len(input.Typed))
	for i, v := range input.Typed {
		t[i] = v
	}
	result.Typed = t
	for k, v := range input.Named {
		result.Named[k] = v
	}

	return result, nil
}

func SeedsProto(
	input *core.Seeds,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*vagrant_plugin_sdk.Args_Seeds, error) {
	result := &vagrant_plugin_sdk.Args_Seeds{
		Named: make(map[string]*anypb.Any),
	}
	t, err := DirectProto(
		&component.Direct{
			Arguments: input.Typed,
		},
		log, internal, ctx,
	)

	if err != nil {
		return nil, err
	}

	result.Typed = t.Arguments

	nv := make([]interface{}, len(input.Named))

	i := 0
	for _, v := range input.Named {
		nv[i] = v
	}

	t, err = DirectProto(
		&component.Direct{
			Arguments: nv,
		},
		log, internal, ctx,
	)

	if err != nil {
		return nil, err
	}

	i = 0
	for k, _ := range input.Named {
		if t.Arguments[i] != nil {
			result.Named[k] = t.Arguments[i]
		}
	}

	return result, nil
}

func Array(
	input *vagrant_plugin_sdk.Args_Array,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) ([]interface{}, error) {
	result, err := Direct(
		&vagrant_plugin_sdk.Args_Direct{
			Arguments: input.List,
		},
		log,
		internal,
		ctx,
	)

	if err != nil {
		return nil, err
	}

	return result.Arguments, nil
}

func ArrayProto(
	input []interface{},
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*vagrant_plugin_sdk.Args_Array, error) {
	result, err := DirectProto(
		&component.Direct{
			Arguments: input,
		},
		log,
		internal,
		ctx,
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Array{
		List: result.Arguments,
	}, nil
}

func Hash(
	input *vagrant_plugin_sdk.Args_Hash,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (result map[string]interface{}, err error) {
	result = make(map[string]interface{}, len(input.Fields))

	for k, v := range input.Fields {
		r, err := Direct(
			&vagrant_plugin_sdk.Args_Direct{Arguments: []*anypb.Any{v}},
			log,
			internal,
			ctx,
		)

		if err != nil {
			return nil, err
		}

		result[k] = r.Arguments[0]
	}

	return
}

func HashProto(
	input map[string]interface{},
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*vagrant_plugin_sdk.Args_Hash, error) {
	content := make(map[string]*anypb.Any, len(input))

	for k, v := range input {
		r, err := DirectProto(
			&component.Direct{
				Arguments: []interface{}{v},
			},
			log,
			internal,
			ctx,
		)

		if err != nil {
			return nil, err
		}

		content[k] = r.Arguments[0]
	}

	return &vagrant_plugin_sdk.Args_Hash{
		Fields: content,
	}, nil
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

func Direct(
	input *vagrant_plugin_sdk.Args_Direct,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*component.Direct, error) {
	args := make([]interface{}, len(input.Arguments))

	for i := 0; i < len(args); i++ {
		v := input.Arguments[i]
		// List item are Any values so start with decoding
		_, val, err := dynamic.DecodeAny(v)
		if err != nil {
			return nil, err
		}

		// First attempt to map the decoded value using
		// the well known type protos
		nv, err := dynamic.MapFromWellKnownProto(val.(proto.Message))

		// If we didn't generate an error, set the value and move on
		if err == nil {
			args[i] = nv
			continue
		}

		// Next attempt a blind map to convert the value into something
		// we may have a converter for
		nv, err = dynamic.BlindMap(val, internal.Mappers,
			argmapper.Typed(internal, ctx, log))

		// Again, if there's no error, set the value and move on
		if err == nil {
			args[i] = nv
			continue
		}

		// Log the mapping failure to the debug output
		log.Warn("failed to map decoded direct argument",
			"value", val,
			"error", err,
		)

		// Set the decoded value into the result set since it's
		// the best we can do
		args[i] = val
	}

	return &component.Direct{Arguments: args}, nil
}

func DirectProto(
	input *component.Direct,
	log hclog.Logger,
	internal *pluginargs.Internal,
	ctx context.Context,
) (*vagrant_plugin_sdk.Args_Direct, error) {
	list := make([]*anypb.Any, len(input.Arguments))
	for i := 0; i < len(list); i++ {
		arg := input.Arguments[i]
		var v interface{}
		v, err := dynamic.MapToWellKnownProto(arg)
		if err != nil {
			v, err = dynamic.UnknownMap(arg, (*proto.Message)(nil),
				internal.Mappers,
				argmapper.Typed(internal),
				argmapper.Typed(ctx),
				argmapper.Typed(log),
			)
			if err != nil {
				return nil, err
			}
		}

		if v == nil {
			list[i] = nil
		} else if list[i], err = dynamic.EncodeAny(v.(proto.Message)); err != nil {
			return nil, err
		}
	}

	return &vagrant_plugin_sdk.Args_Direct{
		Arguments: list,
	}, nil
}

func HostProto(
	input component.Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Host, error) {
	cid := fmt.Sprintf("%p", input)
	if ch := internal.Cache.Get(cid); ch != nil {
		return ch.(*vagrant_plugin_sdk.Args_Host), nil
	}
	p := &plugincomponent.HostPlugin{
		BasePlugin: basePlugin(input, internal),
		Impl:       input,
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
	proto := &vagrant_plugin_sdk.Args_Host{
		Network:  ep.Network(),
		Target:   ep.String(),
		StreamId: id,
	}
	internal.Cache.Register(cid, proto)
	return proto, nil
}

func Host(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Host, error) {
	p := &plugincomponent.HostPlugin{
		BasePlugin: basePlugin(nil, internal),
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
		BasePlugin: basePlugin(input, internal),
		Impl:       input,
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
		BasePlugin: basePlugin(nil, internal),
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

func Box(ctx context.Context,
	input *vagrant_plugin_sdk.Args_Box,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Box, error) {
	// Create our plugin
	p := &plugincore.BoxPlugin{
		BasePlugin: basePlugin(nil, internal),
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Box), nil
}

func BoxProto(
	box core.Box,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Box, error) {
	n, err := box.Name()
	if err != nil {
		return nil, err
	}
	v, err := box.Version()
	if err != nil {
		return nil, err
	}
	pr, err := box.Provider()
	if err != nil {
		return nil, err
	}
	cid := n + "-" + v + "-" + pr
	if ch := internal.Cache.Get(cid); ch != nil {
		return ch.(*vagrant_plugin_sdk.Args_Box), nil
	}

	log.Warn("failed to locate cached box", "cid", cid)

	// Create our plugin
	p := &plugincore.BoxPlugin{
		BasePlugin: basePlugin(box, internal),
		Impl:       box,
	}

	log.Warn("wrapping box to generate proto", "cid", cid)

	id, ep, err := wrapClient(box, p, internal)
	if err != nil {
		return nil, err
	}

	proto := &vagrant_plugin_sdk.Args_Box{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}

	log.Warn("registered box into cache", "cid", cid, "proto", proto, "cache", hclog.Fmt("%p", internal.Cache))
	internal.Cache.Register(cid, proto)

	return proto, nil
}

func BoxCollection(ctx context.Context,
	input *vagrant_plugin_sdk.Args_BoxCollection,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.BoxCollection, error) {
	// Create our plugin
	p := &plugincore.BoxCollectionPlugin{
		BasePlugin: basePlugin(nil, internal),
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.BoxCollection), nil
}

func BoxCollectionProto(
	boxCollection core.BoxCollection,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_BoxCollection, error) {
	cid := fmt.Sprintf("%p", boxCollection)
	if ch := internal.Cache.Get(cid); ch != nil {
		return ch.(*vagrant_plugin_sdk.Args_BoxCollection), nil
	}

	log.Warn("failed to locate cached box collection", "cid", cid)

	// Create our plugin
	p := &plugincore.BoxCollectionPlugin{
		BasePlugin: basePlugin(boxCollection, internal),
		Impl:       boxCollection,
	}

	log.Warn("wrapping box to generate proto", "cid", cid)

	id, ep, err := wrapClient(boxCollection, p, internal)
	if err != nil {
		return nil, err
	}

	proto := &vagrant_plugin_sdk.Args_BoxCollection{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}

	log.Warn("registered box collection into cache", "cid", cid, "proto", proto, "cache", hclog.Fmt("%p", internal.Cache))
	internal.Cache.Register(cid, proto)

	return proto, nil
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
		Logger:  log.ResetNamed("vagrant.wrapped"),
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
		BasePlugin: basePlugin(nil, internal),
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
	cid := fmt.Sprintf("%p", bag)
	if ch := internal.Cache.Get(cid); ch != nil {
		return ch.(*vagrant_plugin_sdk.Args_StateBag), nil
	}

	// Create our plugin
	p := &plugincore.StateBagPlugin{
		BasePlugin: basePlugin(bag, internal),
		Impl:       bag,
	}

	id, ep, err := wrapClient(bag, p, internal)
	if err != nil {
		return nil, err
	}

	proto := &vagrant_plugin_sdk.Args_StateBag{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}

	internal.Cache.Register(cid, proto)

	return proto, nil
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
		BasePlugin: basePlugin(b, internal),
		Impl:       b,
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
		BasePlugin: basePlugin(nil, internal),
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
		BasePlugin: basePlugin(c, internal),
		Impl:       c,
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
		BasePlugin: basePlugin(nil, internal),
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
	cid, err := p.ResourceId()
	if err != nil {
		log.Warn("failed to get resource ID from project", "error", err)
		cid = fmt.Sprintf("%p", p)
	}

	if ch := internal.Cache.Get(cid); ch != nil {
		return ch.(*vagrant_plugin_sdk.Args_Project), nil
	}

	pp := &plugincore.ProjectPlugin{
		BasePlugin: basePlugin(p, internal),
		Impl:       p,
	}

	id, ep, err := wrapClient(p, pp, internal)
	if err != nil {
		return nil, err
	}

	proto := &vagrant_plugin_sdk.Args_Project{
		StreamId: id,
		Network:  ep.Network(),
		Target:   ep.String()}

	internal.Cache.Register(cid, proto)

	return proto, nil
}

func Project(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Project,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Project, error) {
	cid := input.Target
	if cid != "" {
		if ch := internal.Cache.Get(cid); ch != nil {
			return ch.(core.Project), nil
		}
	}

	p := &plugincore.ProjectPlugin{
		BasePlugin: basePlugin(nil, internal),
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	internal.Cache.Register(cid, client)

	return client.(core.Project), nil
}

func VagrantfileSyncedFolderToFolder(
	f *vagrant_plugin_sdk.Vagrantfile_SyncedFolder,
) (*core.Folder, error) {
	var result *core.Folder
	return result, mapstructure.Decode(f, &result)
}

func FolderToVagrantfileSyncedFolder(
	f *core.Folder,
) (*vagrant_plugin_sdk.Vagrantfile_SyncedFolder, error) {
	var result *vagrant_plugin_sdk.Vagrantfile_SyncedFolder
	return result, mapstructure.Decode(f, &result)
}

func SyncedFolderProto(
	s component.SyncedFolder,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_SyncedFolder, error) {
	sp := &plugincomponent.SyncedFolderPlugin{
		BasePlugin: basePlugin(s, internal),
		Impl:       s,
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
		BasePlugin: basePlugin(nil, internal),
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
		BasePlugin: basePlugin(t, internal),
		Impl:       t,
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
		BasePlugin: basePlugin(nil, internal),
	}

	client, err := wrapConnect(ctx, t, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Provider), nil
}

func PathProto(
	p path.Path,
) *vagrant_plugin_sdk.Args_Path {
	return &vagrant_plugin_sdk.Args_Path{
		Path: p.String(),
	}
}

func Path(
	input *vagrant_plugin_sdk.Args_Path,
) path.Path {
	return path.NewPath(input.Path)
}

func MachineProject(
	m core.Machine,
) (core.Project, error) {
	return m.Project()
}

func TargetProject(
	t core.Target,
) (core.Project, error) {
	return t.Project()
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
		BasePlugin: basePlugin(t, internal),
		Impl:       t,
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
		BasePlugin: basePlugin(nil, internal),
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
		BasePlugin: basePlugin(m, internal),
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
		BasePlugin: basePlugin(nil, internal),
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
		BasePlugin: basePlugin(t, internal),
		Impl:       t,
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
		BasePlugin: basePlugin(nil, internal),
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

// creates a BasePlugin configuration using an optional
// source client and the internal args
func basePlugin(
	src interface{},
	internal *pluginargs.Internal,
) *plugincomponent.BasePlugin {
	if w, ok := src.(wrappable); ok {
		return w.Wrap()
	}
	return &plugincomponent.BasePlugin{
		Mappers: internal.Mappers,
		Logger:  internal.Logger,
		Cache:   internal.Cache,
		Wrapped: true,
	}
}

func init() {
	for _, fn := range All {
		mFn, err := argmapper.NewFunc(fn)
		if err != nil {
			panic(err)
		}
		plugincomponent.MapperFns = append(plugincomponent.MapperFns, mFn)
		plugincomponent.ProtomapperAllMap[reflect.TypeOf(fn)] = struct{}{}
	}
	for _, fn := range WellKnownTypes {
		mFn, err := argmapper.NewFunc(fn)
		if err != nil {
			panic(err)
		}
		dynamic.WellKnownTypeFns = append(dynamic.WellKnownTypeFns, mFn)
	}
}

type pluginMetadata interface {
	SetRequestMetadata(k, v string)
}

type wrappable interface {
	Wrap() *plugincomponent.BasePlugin
}
