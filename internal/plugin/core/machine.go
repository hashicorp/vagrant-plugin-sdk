package core

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
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
		targetServer: targetServer{
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
	targetServer

	Impl core.Machine
	vagrant_plugin_sdk.UnimplementedTargetMachineServiceServer
}

func (t *targetMachineClient) Communicate() (comm core.Communicator, err error) {

	// TODO
	return nil, errNotImplemented
}

func (t *targetMachineClient) Guest() (g core.Guest, err error) {
	// TODO
	return nil, errNotImplemented
}

func (t *targetMachineClient) MachineState() (state *core.MachineState, err error) {
	// TODO
	return nil, errNotImplemented
}

func (t *targetMachineClient) SetMachineState(state *core.MachineState) (err error) {
	return errNotImplemented
}

func (t *targetMachineClient) IndexUUID() (id string, err error) {
	// TODO
	return "", errNotImplemented
}

func (t *targetMachineClient) SetUUID(uuid string) (err error) {
	return errNotImplemented
}

func (t *targetMachineClient) Inspect() (printable string, err error) {
	// TODO
	return "", errNotImplemented
}

func (t *targetMachineClient) Reload() (err error) {
	// TODO
	return errNotImplemented
}

func (t *targetMachineClient) ConnectionInfo() (info *core.ConnectionInfo, err error) {
	// TODO
	return nil, errNotImplemented
}

func (t *targetMachineClient) UID() (user_id int, err error) {
	// TODO
	return 10, errNotImplemented
}

func (t *targetMachineClient) SyncedFolders() (folders []core.SyncedFolder, err error) {
	// TODO
	return nil, errNotImplemented
}

func (t *targetMachineClient) SetName(name string) (err error) {
	_, err = t.client.SetName(t.ctx, &vagrant_plugin_sdk.Target_Machine_SetNameRequest{
		Name: name})
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

func (t *targetMachineClient) Box() (b core.Box, err error) {
	r, err := t.client.Box(t.ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(r, (*core.Box)(nil),
		argmapper.Typed(t.ctx))
	if err == nil {
		b = result.(core.Box)
	}

	return
}

func (t *targetMachineClient) Provider() (p core.Provider, err error) {
	return nil, errNotImplemented
}

func (t *targetMachineClient) VagrantfileName() (name string, err error) {
	r, err := t.client.VagrantfileName(t.ctx, &empty.Empty{})
	if err == nil {
		name = r.Name
	}

	return
}

func (t *targetMachineClient) VagrantfilePath() (p path.Path, err error) {
	r, err := t.client.VagrantfilePath(t.ctx, &empty.Empty{})
	if err == nil {
		p = path.NewPath(r.Path)
	}

	return
}

func (t *targetMachineClient) UpdatedAt() (utime *time.Time, err error) {
	r, err := t.client.UpdatedAt(t.ctx, &empty.Empty{})
	if err == nil {
		ut := r.UpdatedAt.AsTime()
		utime = &ut
	}

	return
}

func (t *targetMachineServer) SetName(
	ctx context.Context,
	in *vagrant_plugin_sdk.Target_Machine_SetNameRequest,
) (*empty.Empty, error) {
	if err := t.Impl.SetName(in.Name); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
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

func (t *targetMachineServer) SetUUID(
	ctx context.Context,
	in *vagrant_plugin_sdk.Target_Machine_SetUUIDRequest,
) (*empty.Empty, error) {
	err := t.Impl.SetUUID(in.Uuid)
	return &empty.Empty{}, err
}

func (t *targetMachineServer) GetUUID(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_Machine_GetUUIDResponse, error) {
	uuid, err := t.Impl.IndexUUID()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_Machine_GetUUIDResponse{
		Uuid: uuid}, nil
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

func (t *targetMachineServer) Provider(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Provider, err error) {
	p, err := t.Impl.Provider()
	if err != nil {
		return
	}

	result, err := t.Map(p, (**vagrant_plugin_sdk.Args_Provider)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Provider)
	}

	return
}

func (t *targetMachineServer) VagrantfileName(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_Machine_VagrantfileNameResponse, error) {
	n, err := t.Impl.VagrantfileName()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_Machine_VagrantfileNameResponse{
		Name: n}, nil
}

func (t *targetMachineServer) VagrantfilePath(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_Machine_VagrantfilePathResponse, error) {
	n, err := t.Impl.VagrantfilePath()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_Machine_VagrantfilePathResponse{
		Path: n.String()}, nil
}

func (t *targetMachineServer) UpdatedAt(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_Machine_UpdatedAtResponse, error) {
	u, err := t.Impl.UpdatedAt()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_Machine_UpdatedAtResponse{
		UpdatedAt: timestamppb.New(*u)}, nil
}

func (s *targetMachineServer) DataDir(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_DataDir_Target, err error) {
	d, err := s.Impl.DataDir()
	if err != nil {
		return
	}
	result, err := s.Map(d, (**vagrant_plugin_sdk.Args_DataDir_Target)(nil))
	if err != nil {
		return
	}
	r = result.(*vagrant_plugin_sdk.Args_DataDir_Target)

	return
}

var (
	_ plugin.Plugin     = (*TargetMachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*TargetMachinePlugin)(nil)
	_ core.Machine      = (*targetMachineClient)(nil)
)
