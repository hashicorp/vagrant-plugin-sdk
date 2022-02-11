package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type prioritizedPluginServer struct {
	*BaseServer
}

func (s *prioritizedPluginServer) PluginPriority(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.PluginInfo_Priority, error) {
	return &vagrant_plugin_sdk.PluginInfo_Priority{
		Priority: 0, // TODO
	}, nil
}

func (s *prioritizedPluginServer) SetPluginPriority(
	ctx context.Context,
	p *vagrant_plugin_sdk.PluginInfo_Priority,
) (*empty.Empty, error) {
	// TODO
	return &empty.Empty{}, nil
}
