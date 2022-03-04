package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type BoxPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.Box
	*vplugin.BasePlugin
}

func (p *BoxPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &boxClient{
		client:     vagrant_plugin_sdk.NewBoxServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *BoxPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBoxServiceServer(s, &boxServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type boxClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.BoxServiceClient
}

func (b *boxClient) AutomaticUpdateCheckAllowed() (allowed bool, err error) {
	r, err := b.client.AutomaticUpdateCheckAllowed(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return r.Allowed, nil
}

func (b *boxClient) Destroy() (err error) {
	_, err = b.client.Destroy(b.Ctx, &emptypb.Empty{})
	return
}

func (b *boxClient) Directory() (path string, err error) {
	dir, err := b.client.Directory(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return dir.Path, nil
}

func (b *boxClient) HasUpdate(version string) (updateAvailable bool, err error) {
	result, err := b.client.HasUpdate(
		b.Ctx,
		&vagrant_plugin_sdk.Box_HasUpdateRequest{Version: version},
	)
	if err != nil {
		return
	}
	return result.HasUpdate, nil
}

func (b *boxClient) InUse(index core.TargetIndex) (inUse bool, err error) {
	targetIndex, err := b.Map(index, (*vagrant_plugin_sdk.Args_TargetIndex)(nil), argmapper.Typed(b.Ctx))
	if err != nil {
		return
	}
	result, err := b.client.InUse(
		b.Ctx,
		targetIndex.(*vagrant_plugin_sdk.Args_TargetIndex),
	)
	if err != nil {
		return
	}
	return result.InUse, nil
}

func (b *boxClient) Machines(index core.TargetIndex) (machines []core.Machine, err error) {
	targetIndex, err := b.Map(index, (*vagrant_plugin_sdk.Args_TargetIndex)(nil), argmapper.Typed(b.Ctx))
	if err != nil {
		return
	}
	result, err := b.client.Machines(
		b.Ctx,
		targetIndex.(*vagrant_plugin_sdk.Args_TargetIndex),
	)
	if err != nil {
		return
	}
	machines = []core.Machine{}
	for _, m := range result.Machines {
		coreMachine, err := b.Map(m, (*core.Machine)(nil), argmapper.Typed(b.Ctx))
		if err != nil {
			return nil, err
		}
		machines = append(machines, coreMachine.(core.Machine))
	}
	return machines, nil
}

func (b *boxClient) Metadata() (metadata core.BoxMetadataMap, err error) {
	meta, err := b.client.Metadata(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	boxMetadataMap, err := b.Map(meta, (*core.BoxMetadataMap)(nil), argmapper.Typed(b.Ctx))
	if err != nil {
		return nil, err
	}
	return boxMetadataMap.(core.BoxMetadataMap), nil
}

func (b *boxClient) MetadataURL() (url string, err error) {
	result, err := b.client.MetadataURL(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return result.MetadataUrl, nil
}

func (b *boxClient) Name() (name string, err error) {
	result, err := b.client.Name(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return result.Name, nil
}

func (b *boxClient) Provider() (name string, err error) {
	result, err := b.client.Provider(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return result.Provider, nil
}

func (b *boxClient) Repackage(path string) (err error) {
	_, err = b.client.Repackage(
		b.Ctx,
		&vagrant_plugin_sdk.Args_Path{Path: path},
	)
	return
}

func (b *boxClient) Version() (version string, err error) {
	result, err := b.client.Version(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	return result.Version, nil
}

func (b *boxClient) Compare(box core.Box) (int, error) {
	return 0, nil
}

type boxServer struct {
	*vplugin.BaseServer

	Impl core.Box
	vagrant_plugin_sdk.UnimplementedBoxServiceServer
}

func (b *boxServer) AutomaticUpdateCheckAllowed(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_AutomaticUpdateCheckAllowedResponse, err error) {
	d, err := b.Impl.AutomaticUpdateCheckAllowed()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_AutomaticUpdateCheckAllowedResponse{
		Allowed: d,
	}, nil
}

func (b *boxServer) Destroy(
	ctx context.Context, in *emptypb.Empty,
) (*emptypb.Empty, error) {
	err := b.Impl.Destroy()
	return &emptypb.Empty{}, err
}

func (b *boxServer) HasUpdate(
	ctx context.Context, in *vagrant_plugin_sdk.Box_HasUpdateRequest,
) (r *vagrant_plugin_sdk.Box_HasUpdateResponse, err error) {
	result, err := b.Impl.HasUpdate(in.Version)
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_HasUpdateResponse{
		HasUpdate: result,
	}, nil
}

func (b *boxServer) InUse(
	ctx context.Context, in *vagrant_plugin_sdk.Args_TargetIndex,
) (r *vagrant_plugin_sdk.Box_InUseResponse, err error) {
	targetIndex, err := b.Map(in, (*core.TargetIndex)(nil), argmapper.Typed(ctx))
	if err != nil {
		return
	}

	result, err := b.Impl.InUse(targetIndex.(core.TargetIndex))
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_InUseResponse{
		InUse: result,
	}, nil
}

func (b *boxServer) Machines(
	ctx context.Context, in *vagrant_plugin_sdk.Args_TargetIndex,
) (r *vagrant_plugin_sdk.Box_MachinesResponse, err error) {
	targetIndex, err := b.Map(in, (*core.TargetIndex)(nil), argmapper.Typed(ctx))
	if err != nil {
		return
	}

	result, err := b.Impl.Machines(targetIndex.(core.TargetIndex))
	if err != nil {
		return
	}

	machines := []*vagrant_plugin_sdk.Args_Target_Machine{}
	for _, m := range result {
		machineArg, err := b.Map(m, (**vagrant_plugin_sdk.Args_Target_Machine)(nil), argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		machines = append(machines, machineArg.(*vagrant_plugin_sdk.Args_Target_Machine))
	}

	return &vagrant_plugin_sdk.Box_MachinesResponse{
		Machines: machines,
	}, nil
}

func (b *boxServer) Repackage(
	ctx context.Context, in *vagrant_plugin_sdk.Args_Path,
) (*emptypb.Empty, error) {
	err := b.Impl.Repackage(in.Path)
	return &emptypb.Empty{}, err
}

func (b *boxServer) Directory(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_Path, err error) {
	d, err := b.Impl.Directory()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Args_Path{
		Path: d,
	}, nil
}

func (b *boxServer) Metadata(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_MetadataResponse, err error) {
	meta, err := b.Impl.Metadata()
	if err != nil {
		return
	}

	metadataHash, err := b.Map(meta, (**vagrant_plugin_sdk.Args_Hash)(nil), argmapper.Typed(ctx))
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Box_MetadataResponse{
		Metadata: metadataHash.(*vagrant_plugin_sdk.Args_Hash),
	}, nil
}

func (b *boxServer) MetadataURL(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_MetadataUrlResponse, err error) {
	d, err := b.Impl.MetadataURL()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_MetadataUrlResponse{
		MetadataUrl: d,
	}, nil
}

func (b *boxServer) Name(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_NameResponse, err error) {
	d, err := b.Impl.Name()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_NameResponse{
		Name: d,
	}, nil
}

func (b *boxServer) Provider(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_ProviderResponse, err error) {
	d, err := b.Impl.Provider()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_ProviderResponse{
		Provider: d,
	}, nil
}

func (b *boxServer) Version(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.Box_VersionResponse, err error) {
	d, err := b.Impl.Version()
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_VersionResponse{
		Version: d,
	}, nil
}

func (b *boxServer) Compare(
	ctx context.Context, in *vagrant_plugin_sdk.Args_Box,
) (r *vagrant_plugin_sdk.Box_EqualityResponse, err error) {
	box, err := b.Map(in, (*core.Box)(nil), argmapper.Typed(ctx))
	if err != nil {
		return
	}

	result, err := b.Impl.Compare(box.(core.Box))
	if err != nil {
		return
	}

	return &vagrant_plugin_sdk.Box_EqualityResponse{
		Result: int32(result),
	}, nil
}

var (
	_ plugin.Plugin                       = (*BoxPlugin)(nil)
	_ plugin.GRPCPlugin                   = (*BoxPlugin)(nil)
	_ core.Box                            = (*boxClient)(nil)
	_ vagrant_plugin_sdk.BoxServiceServer = (*boxServer)(nil)
)
