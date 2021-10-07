package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// SyncedFolderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the SyncedFolder component type.
type SyncedFolderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.SyncedFolder // Impl is the concrete implementation
	*BasePlugin
}

func (p *SyncedFolderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := p.NewServer(broker)
	vagrant_plugin_sdk.RegisterSyncedFolderServiceServer(s, &syncedFolderServer{
		Impl:       p.Impl,
		BaseServer: bs,
		capabilityServer: &capabilityServer{
			BaseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "synced_folder",
		},
	})
	return nil
}

func (p *SyncedFolderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	bc := p.NewClient(ctx, broker)
	client := vagrant_plugin_sdk.NewSyncedFolderServiceClient(c)
	return &syncedFolderClient{
		BaseClient: bc,
		client:     client,
		capabilityClient: &capabilityClient{
			client:     client,
			BaseClient: bc,
		},
	}, nil
}

// syncedFolderClient is an implementation of component.SyncedFolder over gRPC.
type syncedFolderClient struct {
	*BaseClient
	*capabilityClient
	client vagrant_plugin_sdk.SyncedFolderServiceClient
}

func (c *syncedFolderClient) Config() (interface{}, error) {
	return configStructCall(c.Ctx, c.client)
}

func (c *syncedFolderClient) ConfigSet(v interface{}) error {
	return configureCall(c.Ctx, c.client, v)
}

func (c *syncedFolderClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.Ctx, c.client)
}

func (c *syncedFolderClient) SyncedFolderFunc() interface{} {
	//TODO
	return nil
}

// syncedFolderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type syncedFolderServer struct {
	*BaseServer
	*capabilityServer

	Impl component.SyncedFolder
}

func (s *syncedFolderServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *syncedFolderServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *syncedFolderServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *syncedFolderServer) UsableSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.UsableFunc())
}

func (s *syncedFolderServer) Usable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.SyncedFolder_UsableResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.UsableFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.SyncedFolder_UsableResp{
		Usable: raw.(bool)}, nil
}

func (s *syncedFolderServer) EnableSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.EnableFunc())
}

func (s *syncedFolderServer) Enable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*empty.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.EnableFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *syncedFolderServer) DisableSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.DisableFunc())
}

func (s *syncedFolderServer) Disable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*empty.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.DisableFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *syncedFolderServer) CleanupSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.CleanupFunc())
}

func (s *syncedFolderServer) Cleanup(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*empty.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.CleanupFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

var (
	_ plugin.Plugin                                = (*SyncedFolderPlugin)(nil)
	_ plugin.GRPCPlugin                            = (*SyncedFolderPlugin)(nil)
	_ vagrant_plugin_sdk.SyncedFolderServiceServer = (*syncedFolderServer)(nil)
	_ component.SyncedFolder                       = (*syncedFolderClient)(nil)
	_ core.SyncedFolder                            = (*syncedFolderClient)(nil)
)
