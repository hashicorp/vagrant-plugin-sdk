package protomappers

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	plugincore "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/core"
	pluginterminal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/terminal"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/multistep"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
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
}

// TODO(spox): make sure these new mappers actually work
// func Machine(input *pb.Args_Machine) (*core.Machine, error) {
// 	var result core.Machine
// 	return &result, mapstructure.Decode(input, &result)
// }

// func MachineProto(input *core.Machine) (*pb.Args_Machine, error) {
// 	var result pb.Args_Machine
// 	return &result, mapstructure.Decode(intput, &result)
// }

// TODO(spox): end of mappers to validate

// JobInfo maps Args.JobInfo to component.JobInfo.
func JobInfo(input *pb.Args_JobInfo) (*component.JobInfo, error) {
	var result component.JobInfo
	return &result, mapstructure.Decode(input, &result)
}

// JobInfoProto
func JobInfoProto(input *component.JobInfo) (*pb.Args_JobInfo, error) {
	var result pb.Args_JobInfo
	return &result, mapstructure.Decode(input, &result)
}

func DatadirBasis(input *pb.Args_DataDir_Basis) *datadir.Basis {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Basis{Dir: dir}
}

func DatadirProject(input *pb.Args_DataDir_Project) *datadir.Project {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Project{Dir: dir}
}

func DatadirMachine(input *pb.Args_DataDir_Project) *datadir.Machine {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Machine{Dir: dir}
}

func DatadirComponent(input *pb.Args_DataDir_Project) *datadir.Component {
	dir := datadir.NewBasicDir(input.RootDir, input.CacheDir, input.DataDir, input.TempDir)
	return &datadir.Component{Dir: dir}
}

func DatadirBasisProto(input *datadir.Basis) *pb.Args_DataDir_Basis {
	return &pb.Args_DataDir_Basis{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirProjectProto(input *datadir.Project) *pb.Args_DataDir_Project {
	return &pb.Args_DataDir_Project{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirMachineProto(input *datadir.Project) *pb.Args_DataDir_Machine {
	return &pb.Args_DataDir_Machine{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

func DatadirComponentProto(input *datadir.Project) *pb.Args_DataDir_Component {
	return &pb.Args_DataDir_Component{
		CacheDir: input.CacheDir().String(),
		DataDir:  input.DataDir().String(),
		TempDir:  input.TempDir().String(),
		RootDir:  input.RootDir().String(),
	}
}

// Logger maps *pb.Args_Logger to an hclog.Logger
func Logger(input *pb.Args_Logger) hclog.Logger {
	// We use the default logger as the base. Within a plugin we always set
	// it so we can confidently use this. This lets plugins potentially mess
	// with this but that's a risk we have to take.
	return hclog.L().ResetNamed(input.Name)
}

func LoggerProto(log hclog.Logger) *pb.Args_Logger {
	return &pb.Args_Logger{
		Name: log.Name(),
	}
}

// TerminalUI maps *pb.Args_TerminalUI to an hclog.TerminalUI
func TerminalUI(
	ctx context.Context,
	input *pb.Args_TerminalUI,
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
) *pb.Args_TerminalUI {
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

	return &pb.Args_TerminalUI{StreamId: id}
}

// Machine maps *pb.Args_Machine to a core.Machine
func Machine(
	ctx context.Context,
	input *pb.Args_Machine,
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
	machineClient := mc.(*plugincore.MachineClient)
	rawMachine, err := machineClient.GetMachine(input.MachineId)
	if err != nil {
		return nil, err
	}

	machine := plugincore.NewMachine(machineClient, rawMachine)
	return machine, nil
}

func LabelSet(input *pb.Args_LabelSet) *component.LabelSet {
	return &component.LabelSet{
		Labels: input.Labels,
	}
}

func LabelSetProto(labels *component.LabelSet) *pb.Args_LabelSet {
	return &pb.Args_LabelSet{Labels: labels.Labels}
}

// StateBag maps StateBag proto to multistep.StateBag.
func StateBag(input *pb.StateBag) (*multistep.BasicStateBag, error) {
	var result multistep.BasicStateBag
	return &result, mapstructure.Decode(input, &result)
}

// StateBag
func StateBagProto(input *multistep.BasicStateBag) (*pb.StateBag, error) {
	var result pb.StateBag
	return &result, mapstructure.Decode(input, &result)
}
