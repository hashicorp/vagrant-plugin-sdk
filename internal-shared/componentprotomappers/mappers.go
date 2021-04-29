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
	Host,
	HostProto,
}

func Host(
	ctx context.Context,
	input *vagrant_plugin_sdk.Args_Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) (core.Host, error) {
	// Create our plugin
	p := &plugininternal.HostPlugin{
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

	return client.(core.Host), nil
}

func HostProto(
	host component.Host,
	log hclog.Logger,
	internal *pluginargs.Internal,
) *vagrant_plugin_sdk.Args_Host {
	// Create our plugin
	p := &plugininternal.HostPlugin{
		Impl:    host,
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

	return &vagrant_plugin_sdk.Args_Host{StreamId: id}
}
