// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	n, err := b.client.BoxName(b.Ctx, &emptypb.Empty{})
	if err != nil {
		b.Logger.Error("failed to get box name",
			"error", err,
		)
		return
	}
	return n.Name
}

func (b *boxMetadataClient) LoadMetadata(url string) (err error) {
	_, err = b.client.LoadMetadata(b.Ctx, &vagrant_plugin_sdk.BoxMetadata_LoadMetadataRequest{Url: url})
	if err != nil {
		b.Logger.Error("failed to load metadata",
			"error", err,
		)
	}

	return
}

func (b *boxMetadataClient) Version(version string, opts ...*core.BoxProvider) (ver *core.BoxVersion, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box version",
				"error", err,
			)
		}
	}()

	boxMetadataOpts := []*vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts{}
	for _, o := range opts {
		var bmo *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
		err = mapstructure.Decode(o, &bmo)
		if err != nil {
			return nil, err
		}
		boxMetadataOpts = append(boxMetadataOpts, bmo)
	}

	v, err := b.client.Version(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_VersionQuery{
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
	defer func() {
		if err != nil {
			b.Logger.Error("failed to list box versions",
				"error", err,
			)
		}
	}()

	boxMetadataOpts := []*vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts{}
	for _, o := range opts {
		var bmo *vagrant_plugin_sdk.BoxMetadata_BoxMetadataOpts
		err = mapstructure.Decode(o, &bmo)
		if err != nil {
			return nil, err
		}
		boxMetadataOpts = append(boxMetadataOpts, bmo)
	}

	v, err := b.client.ListVersions(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_ListVersionsQuery{
			Opts: boxMetadataOpts,
		},
	)
	if err != nil {
		return nil, err
	}
	return v.Versions, nil
}

func (b *boxMetadataClient) Provider(version string, name string) (provider *core.BoxProvider, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box provider",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box provider list",
				"error", err,
			)
		}
	}()

	p, err := b.client.ListProviders(
		b.Ctx,
		&vagrant_plugin_sdk.BoxMetadata_ListProvidersRequest{Version: version},
	)
	if err != nil {
		return nil, err
	}
	return p.Providers, nil
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

func (b *boxMetadataServer) LoadMetadata(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_LoadMetadataRequest,
) (r *emptypb.Empty, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to load box metadata",
				"error", err,
			)
		}
	}()

	err = b.Impl.LoadMetadata(in.Url)
	r = &emptypb.Empty{}
	return
}

func (b *boxMetadataServer) Version(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_VersionQuery,
) (r *vagrant_plugin_sdk.BoxMetadata_VersionResponse, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box version",
				"error", err,
			)
		}
	}()

	opts := []*core.BoxProvider{}
	for _, o := range in.Opts {
		var decodedOpts *core.BoxProvider
		err = mapstructure.Decode(o, &decodedOpts)
		if err != nil {
			return nil, err
		}
		opts = append(opts, decodedOpts)
	}

	v, err := b.Impl.Version(in.Version, opts...)
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.BoxMetadata_VersionResponse{
		Version: v.Version, Status: v.Status, Description: v.Description,
	}, nil
}

func (b *boxMetadataServer) ListVersions(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_ListVersionsQuery,
) (r *vagrant_plugin_sdk.BoxMetadata_ListVersionsResponse, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box version list",
				"error", err,
			)
		}
	}()

	opts := []*core.BoxProvider{}
	for _, o := range in.Opts {
		var decodedOpts *core.BoxProvider
		err = mapstructure.Decode(o, &decodedOpts)
		if err != nil {
			return nil, err
		}
		opts = append(opts, decodedOpts)
	}
	v, err := b.Impl.ListVersions(opts...)
	return &vagrant_plugin_sdk.BoxMetadata_ListVersionsResponse{Versions: v}, nil
}

func (b *boxMetadataServer) Provider(
	ctx context.Context, in *vagrant_plugin_sdk.BoxMetadata_ProviderRequest,
) (r *vagrant_plugin_sdk.BoxMetadata_ProviderResponse, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box provider",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box provider list",
				"error", err,
			)
		}
	}()

	p, err := b.Impl.ListProviders(in.Version)
	return &vagrant_plugin_sdk.BoxMetadata_ListProvidersResponse{Providers: p}, nil
}

var (
	_ plugin.Plugin     = (*BoxMetadataPlugin)(nil)
	_ plugin.GRPCPlugin = (*BoxMetadataPlugin)(nil)
	_ core.BoxMetadata  = (*boxMetadataClient)(nil)
)
