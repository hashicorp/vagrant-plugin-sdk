package core

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type DownloaderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.Downloader
	*vplugin.BasePlugin
}

func (p *DownloaderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterDownloaderServiceServer(s, &downloaderServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

func (p *DownloaderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &downloaderClient{
		client:     vagrant_plugin_sdk.NewDownloaderServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

type downloaderClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.DownloaderServiceClient
}

func (c *downloaderClient) Download() error {
	_, err := c.client.Download(c.Ctx, &emptypb.Empty{})
	return err
}

func (c *downloaderClient) Source() (string, error) {
	url, err := c.client.Source(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
func (c *downloaderClient) Destination() (string, error) {
	path, err := c.client.Destination(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

type downloaderServer struct {
	*vplugin.BaseServer

	Impl core.Downloader
	vagrant_plugin_sdk.UnimplementedDownloaderServiceServer
}

func (s *downloaderServer) Download(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *emptypb.Empty, err error) {
	return
}

func (s *downloaderServer) Source(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_URL, err error) {
	return
}

func (s *downloaderServer) Destination(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_Path, err error) {
	return
}

var (
	_ plugin.Plugin                              = (*DownloaderPlugin)(nil)
	_ plugin.GRPCPlugin                          = (*DownloaderPlugin)(nil)
	_ vagrant_plugin_sdk.DownloaderServiceServer = (*downloaderServer)(nil)
	_ core.Downloader                            = (*downloaderClient)(nil)
)
