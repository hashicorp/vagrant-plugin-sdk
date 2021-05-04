package componentprotomappers

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	plugininternal "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// All is the list of all mappers as raw function pointers.
var All = []interface{}{
	Host,
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

	client, err := p.GRPCClient(ctx, internal.Broker, conn)
	return client.(core.Host), nil
}
