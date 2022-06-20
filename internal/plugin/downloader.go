package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DownloaderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Downloader
	*BasePlugin
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
	cl := vagrant_plugin_sdk.NewDownloaderServiceClient(c)
	return &downloaderClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

type downloaderClient struct {
	*BaseClient

	client vagrant_plugin_sdk.DownloaderServiceClient
}

// func (c *downloaderClient) Config() (interface{}, error) {
// 	return configStructCall(c.Ctx, c.client)
// }

// func (c *downloaderClient) ConfigSet(v interface{}) error {
// 	return configureCall(c.Ctx, c.client, v)
// }

// func (c *downloaderClient) Documentation() (*docs.Documentation, error) {
// 	return documentationCall(c.Ctx, c.client)
// }

func (c *downloaderClient) DownloadFunc() interface{} {
	spec, err := c.client.DownloadSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Download(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return err
		}
		return nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *downloaderClient) Download() error {
	f := c.DownloadFunc()
	_, err := c.CallDynamicFunc(f, false, argmapper.Typed(c.Ctx))
	return err
}

type downloaderServer struct {
	*BaseServer

	Impl component.Downloader
	vagrant_plugin_sdk.UnsafeDownloaderServiceServer
}

// func (s *downloaderServer) ConfigStruct(
// 	ctx context.Context,
// 	empty *emptypb.Empty,
// ) (*vagrant_plugin_sdk.Config_StructResp, error) {
// 	return configStruct(s.Impl)
// }

// func (s *downloaderServer) Configure(
// 	ctx context.Context,
// 	req *vagrant_plugin_sdk.Config_ConfigureRequest,
// ) (*emptypb.Empty, error) {
// 	return configure(s.Impl, req)
// }

// func (s *downloaderServer) Documentation(
// 	ctx context.Context,
// 	empty *emptypb.Empty,
// ) (*vagrant_plugin_sdk.Config_Documentation, error) {
// 	return documentation(s.Impl)
// }

func (s *downloaderServer) DownloadSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "download"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.DownloadFunc())
}

func (s *downloaderServer) Download(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.DownloadFunc(), false, args.Args,
		argmapper.Typed(ctx))

	return &emptypb.Empty{}, err
}

var (
	_ plugin.Plugin                              = (*DownloaderPlugin)(nil)
	_ plugin.GRPCPlugin                          = (*DownloaderPlugin)(nil)
	_ vagrant_plugin_sdk.DownloaderServiceServer = (*downloaderServer)(nil)
	_ component.Downloader                       = (*downloaderClient)(nil)
	_ core.Downloader                            = (*downloaderClient)(nil)
	_ core.Seeder                                = (*downloaderClient)(nil)
)
