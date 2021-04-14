package protomappers

import (
	"context"
	"io"
	"time"

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
	"github.com/hashicorp/vagrant-plugin-sdk/multistep"
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
	DatadirMachine,
	DatadirComponent,
	DatadirBasisProto,
	DatadirProjectProto,
	DatadirMachineProto,
	DatadirComponentProto,
	Logger,
	LoggerProto,
	TerminalUI,
	TerminalUIProto,
	LabelSet,
	LabelSetProto,
	StateBag,
	StateBagProto,
	Machine,
	Flags,
	FlagsProto,
	MapToProto,
	ProtoToMap,
	Project,
	// CommandInfo,
	CommandInfoProto,
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

func DatadirMachine(input *vagrant_plugin_sdk.Args_DataDir_Project) *datadir.Machine {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Machine{Dir: dir}
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

func DatadirMachineProto(input *datadir.Project) *vagrant_plugin_sdk.Args_DataDir_Machine {
	return &vagrant_plugin_sdk.Args_DataDir_Machine{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirComponentProto(input *datadir.Project) *vagrant_plugin_sdk.Args_DataDir_Component {
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
func Machine(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Machine,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*plugincore.Machine, error) {
	p := &plugincore.MachinePlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	timeout := 5 * time.Second
	// Create a new cancellation context so we can cancel in the case of an error
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to the local server
	conn, err := grpc.DialContext(ctx, input.ServerAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	mc, err := p.GRPCClient(ctx, internal.Broker, conn)
	machineClient, ok := mc.(*plugincore.MachineClient)
	if !ok {
		panic("failed to create machine client")
	}

	return plugincore.NewMachine(machineClient, input.ResourceId), nil
}

// Machine maps component.Machine to a *vagrant_plugin_sdk.Args_Machine
func MachineProto(
	machine *plugincore.Machine,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Machine, error) {
	return &vagrant_plugin_sdk.Args_Machine{
		ResourceId: machine.ResourceID,
		ServerAddr: machine.ServerAddr,
	}, nil
}

// Project maps *vagrant_plugin_sdk.Args_Project to a core.Project
func Project(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Project,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*plugincore.Project, error) {
	p := &plugincore.ProjectPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	timeout := 5 * time.Second
	// Create a new cancellation context so we can cancel in the case of an error
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Connect to the local server
	conn, err := grpc.DialContext(ctx, input.ServerAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	internal.Cleanup.Do(func() { conn.Close() })

	mc, err := p.GRPCClient(ctx, internal.Broker, conn)
	projectClient, ok := mc.(*plugincore.ProjectClient)
	if !ok {
		panic("failed to create machine client")
	}

	// TODO: decode input project into a project
	result := plugincore.NewProject(projectClient)
	mapstructure.Decode(input, &result)

	return result, nil
}

func LabelSet(input *vagrant_plugin_sdk.Args_LabelSet) *component.LabelSet {
	return &component.LabelSet{
		Labels: input.Labels,
	}
}

func LabelSetProto(labels *component.LabelSet) *vagrant_plugin_sdk.Args_LabelSet {
	return &vagrant_plugin_sdk.Args_LabelSet{Labels: labels.Labels}
}

// StateBag maps StateBag proto to multistep.StateBag.
func StateBag(input *vagrant_plugin_sdk.Args_StateBag) (*multistep.BasicStateBag, error) {
	var result multistep.BasicStateBag
	return &result, mapstructure.Decode(input, &result)
}

// StateBag
func StateBagProto(input *multistep.BasicStateBag) (*vagrant_plugin_sdk.Args_StateBag, error) {
	var result vagrant_plugin_sdk.Args_StateBag
	return &result, mapstructure.Decode(input, &result)
}

func CommandInfo(input *vagrant_plugin_sdk.Command_CommandInfo) (*core.CommandInfo, error) {
	flags, err := Flags(input.Flags)

	subcommands := []*core.CommandInfo{}
	for _, cmd := range input.Subcommands {
		subcommand, err := CommandInfo(cmd)
		if err != nil {
			return nil, err
		}
		subcommands = append(subcommands, subcommand)
	}

	result := &core.CommandInfo{
		Flags:       flags,
		Name:        input.Name,
		Help:        input.Help,
		Synopsis:    input.Synopsis,
		Subcommands: subcommands,
	}
	return result, err
}

func CommandInfoProto(input *core.CommandInfo) (*vagrant_plugin_sdk.Command_CommandInfo, error) {
	var result vagrant_plugin_sdk.Command_CommandInfo
	err := mapstructure.Decode(input, &result)
	if err != nil {
		return nil, err
	}
	result.Flags, err = FlagsProto(input.Flags)
	// TODO: subcommands
	// result.Subcommands = CommandInfoProto()
	return &result, err
}

func serverListener() (net.Listener, error) {
	if runtime.GOOS == "windows" {
		return serverListener_tcp()
	}

	return serverListener_unix()
}

func serverListener_tcp() (net.Listener, error) {
	envMinPort := os.Getenv("PLUGIN_MIN_PORT")
	envMaxPort := os.Getenv("PLUGIN_MAX_PORT")

	var minPort, maxPort int64
	var err error

	switch {
	case len(envMinPort) == 0:
		minPort = 0
	default:
		minPort, err = strconv.ParseInt(envMinPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get value from PLUGIN_MIN_PORT: %v", err)
		}
	}

	switch {
	case len(envMaxPort) == 0:
		maxPort = 0
	default:
		maxPort, err = strconv.ParseInt(envMaxPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get value from PLUGIN_MAX_PORT: %v", err)
		}
	}

	if minPort > maxPort {
		return nil, fmt.Errorf("PLUGIN_MIN_PORT value of %d is greater than PLUGIN_MAX_PORT value of %d", minPort, maxPort)
	}

	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			return listener, nil
		}
	}

	return nil, errors.New("Couldn't bind plugin TCP listener")
}

func serverListener_unix() (net.Listener, error) {
	tf, err := ioutil.TempFile("", "plugin")
	if err != nil {
		return nil, err
	}
	path := tf.Name()

	// Close the file and remove it because it has to not exist for
	// the domain socket.
	if err := tf.Close(); err != nil {
		return nil, err
	}
	if err := os.Remove(path); err != nil {
		return nil, err
	}

	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	// Wrap the listener in rmListener so that the Unix domain socket file
	// is removed on close.
	return &rmListener{
		Listener: l,
		Path:     path,
	}, nil
}

func (l *rmListener) Close() error {
	// Close the listener itself
	if err := l.Listener.Close(); err != nil {
		return err
	}

	// Remove the file
	return os.Remove(l.Path)
}

// rmListener is an implementation of net.Listener that forwards most
// calls to the listener but also removes a file as part of the close. We
// use this to cleanup the unix domain socket on close.
type rmListener struct {
	net.Listener
	Path string
}

func runServer(listener net.Listener, server *grpc.Server) {
	doneCh := make(chan struct{})
	// Here we use a run group to close this goroutine if the server is shutdown
	// or the broker is shutdown.
	var g run.Group
	{
		// Serve on the listener, if shutting down call GracefulStop.
		g.Add(func() error {
			return server.Serve(listener)
		}, func(err error) {
			server.GracefulStop()
		})
	}
	{
		// block on the closeCh or the doneCh. If we are shutting down close the
		// closeCh.
		closeCh := make(chan struct{})
		g.Add(func() error {
			select {
			case <-doneCh:
			case <-closeCh:
			}
			return nil
		}, func(err error) {
			close(closeCh)
		})
	}

	// Block until we are done
	g.Run()
}
