package core

import (
	"context"
	"errors"
	"reflect"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

var errNotImplemented = errors.New("not implemented")

type TargetMachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Mappers    []*argmapper.Func
	Logger     hclog.Logger
	Impl       core.Machine
	TargetImpl core.Target
}

// Implements plugin.GRPCPlugin
func (t *TargetMachinePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewTargetMachineServiceClient(c)
	b := &base{
		Mappers: t.Mappers,
		Logger:  t.Logger,
		Broker:  broker,
		Cleanup: &pluginargs.Cleanup{},
	}
	return &targetMachineClient{
		client: cl,
		ctx:    ctx,
		base:   b,
		targetClient: &targetClient{
			client: cl,
			ctx:    ctx,
			base:   b,
		},
	}, nil
}

func (t *TargetMachinePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	b := &base{
		Mappers: t.Mappers,
		Logger:  t.Logger,
		Broker:  broker,
		Cleanup: &pluginargs.Cleanup{},
	}

	vagrant_plugin_sdk.RegisterTargetMachineServiceServer(s, &targetMachineServer{
		Impl: t.Impl,
		base: b,
		targetServer: &targetServer{
			Impl: t.TargetImpl,
			base: b,
		},
	})

	return nil
}

// Machine implements core.Machine interface
type targetMachineClient struct {
	*base
	*targetClient

	ctx    context.Context
	client vagrant_plugin_sdk.TargetMachineServiceClient
}

type targetMachineServer struct {
	*base
	*targetServer

	Impl core.Machine
	vagrant_plugin_sdk.UnsafeTargetMachineServiceServer
}

func (t *targetMachineClient) Guest() (g core.Guest, err error) {
	guestResp, err := t.client.Guest(t.ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(guestResp, (*core.Guest)(nil),
		argmapper.Typed(t.ctx))
	if err == nil {
		g = result.(core.Guest)
	}

	return
}

func (t *targetMachineClient) MachineState() (state *core.MachineState, err error) {
	r, err := t.client.GetState(t.ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(r, (**core.MachineState)(nil),
		argmapper.Typed(t.ctx))
	if err == nil {
		state = result.(*core.MachineState)
	}

	return
}

func (t *targetMachineClient) SetMachineState(state *core.MachineState) (err error) {
	stateArg, err := t.Map(
		state,
		(*vagrant_plugin_sdk.Args_Target_Machine_State)(nil),
		argmapper.Typed(t.ctx),
	)
	_, err = t.client.SetState(
		t.ctx,
		&vagrant_plugin_sdk.Target_Machine_SetStateRequest{
			State: stateArg.(*vagrant_plugin_sdk.Args_Target_Machine_State),
		},
	)
	return
}

func (t *targetMachineClient) Inspect() (printable string, err error) {
	name, err := t.Name()
	provider, err := t.Provider()
	printable = "#<" + reflect.TypeOf(t).String() + ": " + name + " (" + reflect.TypeOf(provider).String() + ")>"
	return
}

func (t *targetMachineClient) Reload() (err error) {
	_, err = t.client.Reload(t.ctx, &empty.Empty{})
	return
}

func (t *targetMachineClient) ConnectionInfo() (info *core.ConnectionInfo, err error) {
	connResp, err := t.client.ConnectionInfo(t.ctx, &empty.Empty{})
	return info, mapstructure.Decode(connResp, &info)
}

func (t *targetMachineClient) UID() (id string, err error) {
	uidResp, err := t.client.UID(t.ctx, &empty.Empty{})
	id = uidResp.UserId
	return
}

func (t *targetMachineClient) SyncedFolders() (folders []core.SyncedFolder, err error) {
	sfResp, err := t.client.SyncedFolders(t.ctx, &empty.Empty{})
	folders = []core.SyncedFolder{}
	for _, folder := range sfResp.SyncedFolders {
		f, err := t.Map(folder, (*core.SyncedFolder)(nil), argmapper.Typed(t.ctx))
		if err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}

	return
}

func (t *targetMachineClient) ID() (id string, err error) {
	r, err := t.client.GetID(t.ctx, &empty.Empty{})
	if err == nil {
		id = r.Id
	}

	return
}

func (t *targetMachineClient) SetID(id string) (err error) {
	_, err = t.client.SetID(t.ctx, &vagrant_plugin_sdk.Target_Machine_SetIDRequest{
		Id: id})
	return
}

func (t *targetMachineClient) Box() (b *core.Box, err error) {
	r, err := t.client.Box(t.ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(r, (*core.Box)(nil),
		argmapper.Typed(t.ctx))
	if err == nil {
		b = result.(*core.Box)
	}

	return
}

// Machine Server

func (t *targetMachineServer) ConnectionInfo(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.Target_Machine_ConnectionInfoResponse, err error) {
	connInfo, err := t.Impl.ConnectionInfo()
	if err != nil {
		return nil, err
	}

	result, err := t.Map(connInfo, (**vagrant_plugin_sdk.Target_Machine_ConnectionInfoResponse)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		resp = result.(*vagrant_plugin_sdk.Target_Machine_ConnectionInfoResponse)
	}

	return
}

func (t *targetMachineServer) Reload(
	ctx context.Context,
	_ *empty.Empty,
) (e *empty.Empty, err error) {
	return &empty.Empty{}, t.Impl.Reload()
}

func (t *targetMachineServer) SyncedFolders(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.Target_Machine_SyncedFoldersResponse, err error) {
	syncedFolders, err := t.Impl.SyncedFolders()
	if err != nil {
		return nil, err
	}

	sf := []*vagrant_plugin_sdk.Args_SyncedFolder{}
	for _, folder := range syncedFolders {
		f, err := t.Map(folder, (**vagrant_plugin_sdk.Args_SyncedFolder)(nil), argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		sf = append(sf, f.(*vagrant_plugin_sdk.Args_SyncedFolder))
	}
	resp = &vagrant_plugin_sdk.Target_Machine_SyncedFoldersResponse{
		SyncedFolders: sf,
	}

	return
}

func (t *targetMachineServer) UID(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.Target_Machine_UIDResponse, err error) {
	uid, err := t.Impl.UID()
	if err != nil {
		return nil, err
	}

	resp = &vagrant_plugin_sdk.Target_Machine_UIDResponse{
		UserId: uid,
	}

	return
}

func (t *targetMachineServer) Guest(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Guest, err error) {
	guest, err := t.Impl.Guest()
	if err != nil {
		return nil, err
	}

	result, err := t.Map(guest, (**vagrant_plugin_sdk.Args_Guest)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Guest)
	}

	return
}

func (t *targetMachineServer) GetID(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_Machine_GetIDResponse, error) {
	id, err := t.Impl.ID()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_Machine_GetIDResponse{
		Id: id}, nil
}

func (t *targetMachineServer) SetID(
	ctx context.Context,
	in *vagrant_plugin_sdk.Target_Machine_SetIDRequest,
) (*empty.Empty, error) {
	if err := t.Impl.SetID(in.Id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (t *targetMachineServer) GetState(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Target_Machine_State, err error) {
	s, err := t.Impl.MachineState()
	if err != nil {
		return
	}

	result, err := t.Map(s, (**vagrant_plugin_sdk.Args_Target_Machine_State)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Target_Machine_State)
	}

	return
}

func (t *targetMachineServer) SetState(
	ctx context.Context,
	in *vagrant_plugin_sdk.Target_Machine_SetStateRequest,
) (e *empty.Empty, err error) {
	e = &empty.Empty{}
	s, err := t.Map(in.State, (**core.MachineState)(nil))
	if err != nil {
		return
	}
	err = t.Impl.SetMachineState(s.(*core.MachineState))

	return
}

func (t *targetMachineServer) Box(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Target_Machine_Box, err error) {
	b, err := t.Impl.Box()
	if err != nil {
		return
	}

	result, err := t.Map(b, (**vagrant_plugin_sdk.Args_Target_Machine_Box)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Target_Machine_Box)
	}

	return
}

var (
	_ plugin.Plugin     = (*TargetMachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*TargetMachinePlugin)(nil)
	_ core.Machine      = (*targetMachineClient)(nil)
)
