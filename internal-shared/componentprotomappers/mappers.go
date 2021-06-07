package componentprotomappers

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	plugininternal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	HostProto,
	Host,
}

func HostProto(
	input component.Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (*vagrant_plugin_sdk.Args_Host, error) {
	p := &plugininternal.HostPlugin{
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

	// t := input.(plugininternal.PluginClient).ServerConfig()
	// if t != nil {
	// }
	// return &vagrant_plugin_sdk.Args_Host{
	// 	TargetAddress: "",
	// }, nil
}

func Host(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Host, error) {
	p := &plugininternal.HostPlugin{
		Mappers: internal.Mappers,
		Logger:  log,
	}
	client, err := wrapConnect(ctx, p, input, internal)
	if err != nil {
		return nil, err
	}

	return client.(core.Host), nil

	// timeout := 5 * time.Second
	// // Create a new cancellation context so we can cancel in the case of an error
	// ctx, cancel := context.WithTimeout(ctx, timeout)
	// defer cancel()

	// // Connect to the local server
	// conn, err := grpc.DialContext(ctx, input.TargetAddress,
	// 	grpc.WithBlock(),
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// internal.Cleanup.Do(func() { conn.Close() })

	// client, err := p.GRPCClient(ctx, internal.Broker, conn)
	// return client.(core.Host), nil
}

type connInfo interface {
	GetStreamId() uint32
}

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
