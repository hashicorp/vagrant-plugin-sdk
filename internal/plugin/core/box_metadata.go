package core

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type BoxMetadataPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.BoxMetadata
	*vplugin.BasePlugin
}

func (p *BoxMetadataPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &boxMetadataClient{
		client:     vagrant_plugin_sdk.NewBoxMetadataServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *BoxMetadataPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBoxMetadataServiceServer(s, &boxMetadataServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type boxMetadataClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.BoxMetadataServiceClient
}

func (b *boxMetadataClient) Name() (name string) {
	return
}

func (b *boxMetadataClient) Version(version string, opts *core.BoxMetadataOpts) (ver *core.BoxVersionData, err error) {
	return
}

func (b *boxMetadataClient) ListVersions(opts ...*core.BoxMetadataOpts) (versions []string, err error) {
	return
}

func (b *boxMetadataClient) Provider(version string, name string) (provider *core.BoxProviderData, err error) {
	return
}

func (b *boxMetadataClient) ListProviders(version string) (providers []string, err error) {
	return
}

type boxMetadataServer struct {
	*vplugin.BaseServer

	Impl core.BoxMetadata
	vagrant_plugin_sdk.UnimplementedBoxMetadataServiceServer
}

func (b *boxMetadataServer) Version(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_VersionRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_VersionResponse, err error) {
	return
}

var (
	_ plugin.Plugin     = (*BoxMetadataPlugin)(nil)
	_ plugin.GRPCPlugin = (*BoxMetadataPlugin)(nil)
	_ core.BoxMetadata  = (*boxMetadataClient)(nil)
)
