package protomappers

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	pluginterminal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin/terminal"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	DatadirProject,
	DatadirApp,
	DatadirComponent,
	DatadirProjectProto,
	DatadirAppProto,
	DatadirComponentProto,
	Logger,
	LoggerProto,
	TerminalUI,
	TerminalUIProto,
}

// DatadirProject maps *pb.Args_DataDir_Project to *datadir.Project
func DatadirProject(input *pb.Args_DataDir_Project) *datadir.Project {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.Project{Dir: dir}
}

func DatadirProjectProto(input *datadir.Project) *pb.Args_DataDir_Project {
	return &pb.Args_DataDir_Project{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
	}
}

// DatadirApp maps *pb.Args_DataDir_App to *datadir.App
func DatadirApp(input *pb.Args_DataDir_App) *datadir.App {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.App{Dir: dir}
}

func DatadirAppProto(input *datadir.App) *pb.Args_DataDir_App {
	return &pb.Args_DataDir_App{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
	}
}

// DatadirComponent maps *pb.Args_DataDir_Component to *datadir.Component
func DatadirComponent(input *pb.Args_DataDir_Component) *datadir.Component {
	dir := datadir.NewBasicDir(input.CacheDir, input.DataDir)
	return &datadir.Component{Dir: dir}
}

func DatadirComponentProto(input *datadir.Component) *pb.Args_DataDir_Component {
	return &pb.Args_DataDir_Component{
		CacheDir: input.CacheDir(),
		DataDir:  input.DataDir(),
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
