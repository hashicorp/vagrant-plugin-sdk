package core

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

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

func (b *boxMetadataClient) BoxName() (name string) {
	n, err := b.client.Name(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return n.Name
}

func (b *boxMetadataClient) Version(version string, opts *core.BoxProvider) (ver *core.BoxVersion, err error) {
	var boxMetadataOpts *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
	err = mapstructure.Decode(opts, &boxMetadataOpts)
	if err != nil {
		return nil, err
	}
	v, err := b.client.Version(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_VersionRequest{
			Version: version, Opts: boxMetadataOpts,
		},
	)
	if err != nil {
		return nil, err
	}
	var result *core.BoxVersion
	return result, mapstructure.Decode(v, &result)
}

func (b *boxMetadataClient) ListVersions(opts ...*core.BoxProvider) (versions []string, err error) {
	var boxMetadataOpts *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
	err = mapstructure.Decode(opts, &boxMetadataOpts)
	if err != nil {
		return nil, err
	}
	v, err := b.client.ListVersions(b.Ctx, boxMetadataOpts)
	if err != nil {
		return nil, err
	}
	return v.Versions, nil
}

func (b *boxMetadataClient) Provider(version string, name string) (provider *core.BoxProvider, err error) {
	p, err := b.client.Provider(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_ProviderRequest{Version: version, Name: name},
	)
	if err != nil {
		return nil, err
	}
	var result *core.BoxProvider
	return result, mapstructure.Decode(p, &result)
}

func (b *boxMetadataClient) ListProviders(version string) (providers []string, err error) {
	p, err := b.client.ListProviders(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_ListProvidersRequest{Version: version},
	)
	if err != nil {
		return nil, err
	}
	return p.Providers, nil
}

func (b *boxMetadataClient) Matches(version string, name string, provider *core.BoxProvider) (matches bool, err error) {
	var boxMetadataOpts *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
	err = mapstructure.Decode(provider, &boxMetadataOpts)
	if err != nil {
		return false, err
	}
	m, err := b.client.Matches(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_MatchesRequest{
			Version: version, Name: name, Provider: boxMetadataOpts,
		},
	)
	if err != nil {
		return false, err
	}
	return m.Matches, nil
}

func (b *boxMetadataClient) MatchesAny(version string, name string, provider ...*core.BoxProvider) (matches bool, err error) {
	opts := []*vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts{}
	var boxMetadataOpts *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
	for _, p := range provider {
		err = mapstructure.Decode(p, &boxMetadataOpts)
		if err != nil {
			return false, err
		}
		opts = append(opts, boxMetadataOpts)
	}
	m, err := b.client.MatchesAny(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_MatchesAnyRequest{
			Version: version, Name: name, Providers: opts,
		},
	)
	if err != nil {
		return false, err
	}
	return m.Matches, nil
}

type boxMetadataServer struct {
	*vplugin.BaseServer

	Impl core.BoxMetadata
	vagrant_plugin_sdk.UnimplementedBoxMetadataServiceServer
}

func (b *boxMetadataServer) BoxName(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.BoxMetadata_NameResponse, err error) {
	return &vagrant_plugin_sdk.BoxMetadata_NameResponse{Name: b.Impl.BoxName()}, nil
}

func (b *boxMetadataServer) Version(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_VersionRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_VersionResponse, err error) {
	var opts *core.BoxProvider
	err = mapstructure.Decode(in.Opts, &opts)
	if err != nil {
		return nil, err
	}
	v, err := b.Impl.Version(in.Version, opts)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.BoxMetadata_VersionResponse{
		Version: v.Version, Status: v.Status, Description: v.Description,
	}, nil
}

func (b *boxMetadataServer) ListVersions(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts,
) (r *vagrant_plugin_sdk.BoxMetadata_ListVersionsResponse, err error) {
	var opts *core.BoxProvider
	err = mapstructure.Decode(in, &opts)
	if err != nil {
		return nil, err
	}
	v, err := b.Impl.ListVersions(opts)
	return &vagrant_plugin_sdk.BoxMetadata_ListVersionsResponse{Versions: v}, nil
}

func (b *boxMetadataServer) Provider(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_ProviderRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_ProviderResponse, err error) {
	p, err := b.Impl.Provider(in.Version, in.Name)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.BoxMetadata_ProviderResponse{
		Name: p.Name, Url: p.Url, Checksum: p.Checksum, ChecksumType: p.ChecksumType,
	}, nil
}

func (b *boxMetadataServer) ListProviders(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_ListProvidersRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_ListProvidersResponse, err error) {
	p, err := b.Impl.ListProviders(in.Version)
	return &vagrant_plugin_sdk.BoxMetadata_ListProvidersResponse{Providers: p}, nil
}

func (b *boxMetadataServer) Matches(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_MatchesRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_MatchesResponse, err error) {
	var provider *core.BoxProvider
	err = mapstructure.Decode(in, &provider)
	if err != nil {
		return nil, err
	}
	m, err := b.Impl.Matches(in.Version, in.Name, provider)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.BoxMetadata_MatchesResponse{
		Matches: m,
	}, nil
}

func (b *boxMetadataServer) MatchesAny(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_MatchesAnyRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_MatchesResponse, err error) {
	providers := []*core.BoxProvider{}
	var provider *core.BoxProvider
	for _, p := range in.Providers {
		err = mapstructure.Decode(p, &provider)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	m, err := b.Impl.MatchesAny(in.Version, in.Name, providers...)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.BoxMetadata_MatchesResponse{
		Matches: m,
	}, nil
}

var (
	_ plugin.Plugin     = (*BoxMetadataPlugin)(nil)
	_ plugin.GRPCPlugin = (*BoxMetadataPlugin)(nil)
	_ core.BoxMetadata  = (*boxMetadataClient)(nil)
)
