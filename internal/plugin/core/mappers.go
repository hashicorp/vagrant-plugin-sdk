package core

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	plugincore "github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

var MapperFns = []interface{}{
	MetadataProto,
	Metadata,
	StateProto,
	State,
	TargetDataDirProto,
	TargetDataDir,
}

func StateProto(s plugincore.State) *vagrant_plugin_sdk.Args_Target_State {
	var state vagrant_plugin_sdk.Args_Target_State_State
	switch s {
	case plugincore.CREATED:
		state = vagrant_plugin_sdk.Args_Target_State_CREATED
	case plugincore.DESTROYED:
		state = vagrant_plugin_sdk.Args_Target_State_DESTROYED
	case plugincore.PENDING:
		state = vagrant_plugin_sdk.Args_Target_State_PENDING
	default:
		state = vagrant_plugin_sdk.Args_Target_State_UNKNOWN
	}
	return &vagrant_plugin_sdk.Args_Target_State{
		State: state,
	}
}

func State(s *vagrant_plugin_sdk.Args_Target_State) (state plugincore.State) {
	switch s.State {
	case vagrant_plugin_sdk.Args_Target_State_CREATED:
		state = plugincore.CREATED
	case vagrant_plugin_sdk.Args_Target_State_DESTROYED:
		state = plugincore.DESTROYED
	case vagrant_plugin_sdk.Args_Target_State_PENDING:
		state = plugincore.PENDING
	default:
		state = plugincore.UNKNOWN
	}
	return
}

func TargetDataDirProto(d *datadir.Target) *vagrant_plugin_sdk.Args_DataDir_Target {
	return &vagrant_plugin_sdk.Args_DataDir_Target{
		RootDir:  d.RootDir().String(),
		CacheDir: d.CacheDir().String(),
		DataDir:  d.DataDir().String(),
		TempDir:  d.TempDir().String(),
	}
}

func TargetDataDir(d *vagrant_plugin_sdk.Args_DataDir_Target) *datadir.Target {
	return &datadir.Target{
		Dir: datadir.NewBasicDir(d.RootDir, d.CacheDir, d.DataDir, d.TempDir),
	}
}

func MetadataProto(m map[string]string) *vagrant_plugin_sdk.Args_MetadataSet {
	return &vagrant_plugin_sdk.Args_MetadataSet{
		Metadata: m,
	}
}

func Metadata(m *vagrant_plugin_sdk.Args_MetadataSet) map[string]string {
	return m.Metadata
}

func ProjectProto(
	p plugincore.Project,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Project, error) {
	pp := &ProjectPlugin{
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
) (plugincore.Project, error) {
	p := &ProjectPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(plugincore.Project), nil
}

func TargetProto(
	t plugincore.Target,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Target, error) {
	tp := &TargetPlugin{
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
) (plugincore.Target, error) {
	t := &TargetPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}

	client, err := wrapConnect(ctx, t, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(plugincore.Target), nil
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
