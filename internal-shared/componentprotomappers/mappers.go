package componentprotomappers

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"

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

	conn, err := internal.Broker.Dial(input.GetStreamId())
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

	return client.(core.Host), nil
}
