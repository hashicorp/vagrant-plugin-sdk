package core

import (
	"context"
	"errors"
	"reflect"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

var errNotImplemented = errors.New("not implemented")

type TargetMachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl       core.Machine
	TargetImpl core.Target
	*vplugin.BasePlugin
}

// Implements plugin.GRPCPlugin
func (t *TargetMachinePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewTargetMachineServiceClient(c)
	bc := t.NewClient(ctx, broker, nil)
	return &targetMachineClient{
		client:     cl,
		BaseClient: bc,
		targetClient: &targetClient{
			client:     cl,
			BaseClient: bc,
		},
	}, nil
}

func (t *TargetMachinePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := t.NewServer(broker, nil)

	vagrant_plugin_sdk.RegisterTargetMachineServiceServer(s, &targetMachineServer{
		Impl:       t.Impl,
		BaseServer: bs,
		targetServer: &targetServer{
			Impl:       t.TargetImpl,
			BaseServer: bs,
		},
	})

	return nil
}

// Machine implements core.Machine interface
type targetMachineClient struct {
	*vplugin.BaseClient
	*targetClient

	client vagrant_plugin_sdk.TargetMachineServiceClient
}

type targetMachineServer struct {
	*vplugin.BaseServer
	*targetServer

	Impl core.Machine
	vagrant_plugin_sdk.UnsafeTargetMachineServiceServer
}

func (t *targetMachineClient) Guest() (g core.Guest, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get guest",
				"error", err,
			)
		}
	}()

	guestResp, err := t.client.Guest(t.Ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(guestResp, (*core.Guest)(nil),
		argmapper.Typed(t.Ctx))
	if err == nil {
		g = result.(core.Guest)
	}

	return
}

func (t *targetMachineClient) MachineState() (state *core.MachineState, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine state",
				"error", err,
			)
		}
	}()

	r, err := t.client.GetState(t.Ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(r, (**core.MachineState)(nil),
		argmapper.Typed(t.Ctx))
	if err == nil {
		state = result.(*core.MachineState)
	}

	return
}

func (t *targetMachineClient) SetMachineState(state *core.MachineState) (err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to set machine state",
				"error", err,
			)
		}
	}()

	stateArg, err := t.Map(
		state,
		(*vagrant_plugin_sdk.Args_Target_Machine_State)(nil),
		argmapper.Typed(t.Ctx),
	)
	_, err = t.client.SetState(
		t.Ctx,
		&vagrant_plugin_sdk.Target_Machine_SetStateRequest{
			State: stateArg.(*vagrant_plugin_sdk.Args_Target_Machine_State),
		},
	)
	return
}

func (t *targetMachineClient) Inspect() (printable string, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to inspect machine",
				"error", err,
			)
		}
	}()

	name, err := t.Name()
	provider, err := t.Provider()
	printable = "#<" + reflect.TypeOf(t).String() + ": " + name + " (" + reflect.TypeOf(provider).String() + ")>"
	return
}

func (t *targetMachineClient) ConnectionInfo() (info *core.ConnectionInfo, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine connection info",
				"error", err,
			)
		}
	}()

	connResp, err := t.client.ConnectionInfo(t.Ctx, &empty.Empty{})
	return info, mapstructure.Decode(connResp, &info)
}

func (t *targetMachineClient) UID() (id string, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine uid",
				"error", err,
			)
		}
	}()

	uidResp, err := t.client.UID(t.Ctx, &empty.Empty{})
	id = uidResp.UserId
	return
}

func (t *targetMachineClient) SyncedFolders() (folders []*core.MachineSyncedFolder, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine synced folders",
				"error", err,
			)
		}
	}()

	sfResp, err := t.client.SyncedFolders(t.Ctx, &empty.Empty{})
	folders = []*core.MachineSyncedFolder{}
	for _, folder := range sfResp.SyncedFolders {
		var fp, f interface{}
		fp, err = t.Map(folder.Plugin, (*core.SyncedFolder)(nil), argmapper.Typed(t.Ctx))
		if err != nil {
			return nil, err
		}
		f, err = t.Map(folder.Folder, (*core.Folder)(nil), argmapper.Typed(t.Ctx))
		if err != nil {
			return nil, err
		}
		folders = append(folders, &core.MachineSyncedFolder{
			Plugin: fp.(core.SyncedFolder),
			Folder: f.(*core.Folder),
		})
	}

	return
}

func (t *targetMachineClient) ID() (id string, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine id",
				"error", err,
			)
		}
	}()

	r, err := t.client.GetID(t.Ctx, &empty.Empty{})
	if err == nil {
		id = r.Id
	}

	return
}

func (t *targetMachineClient) SetID(id string) (err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to set machine id",
				"error", err,
			)
		}
	}()

	_, err = t.client.SetID(t.Ctx, &vagrant_plugin_sdk.Target_Machine_SetIDRequest{
		Id: id})
	return
}

func (t *targetMachineClient) Box() (b core.Box, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine box",
				"error", err,
			)
		}
	}()

	r, err := t.client.Box(t.Ctx, &empty.Empty{})
	if err != nil {
		return
	}

	result, err := t.Map(r, (*core.Box)(nil),
		argmapper.Typed(t.Ctx))
	if err == nil {
		b = result.(core.Box)
	}

	return
}

// Machine Server

func (t *targetMachineServer) ConnectionInfo(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.Target_Machine_ConnectionInfoResponse, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine connnection info",
				"error", err,
			)
		}
	}()

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

func (t *targetMachineServer) SyncedFolders(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.Target_Machine_SyncedFoldersResponse, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine synced folders",
				"error", err,
			)
		}
	}()

	syncedFolders, err := t.Impl.SyncedFolders()
	if err != nil {
		return nil, err
	}

	sf := []*vagrant_plugin_sdk.Target_Machine_SyncedFoldersResponse_MachineSyncedFolder{}
	for _, folder := range syncedFolders {
		var plg, f interface{}
		plg, err = t.Map(folder.Plugin, (**vagrant_plugin_sdk.Args_SyncedFolder)(nil), argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		f, err = t.Map(folder.Folder, (**vagrant_plugin_sdk.Vagrantfile_SyncedFolder)(nil), argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		sf = append(sf, &vagrant_plugin_sdk.Target_Machine_SyncedFoldersResponse_MachineSyncedFolder{
			Plugin: plg.(*vagrant_plugin_sdk.Args_SyncedFolder),
			Folder: f.(*vagrant_plugin_sdk.Vagrantfile_SyncedFolder),
		})
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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine uid",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get guest",
				"error", err,
			)
		}
	}()

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
		t.Logger.Error("failed to get machine id",
			"error", err,
		)

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
		t.Logger.Error("failed to get machine id",
			"error", err,
		)

		return nil, err
	}

	return &empty.Empty{}, nil
}

func (t *targetMachineServer) GetState(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Target_Machine_State, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine state",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to set machine state",
				"error", err,
			)
		}
	}()

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
) (r *vagrant_plugin_sdk.Args_Box, err error) {
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get machine box",
				"error", err,
			)
		}
	}()

	b, err := t.Impl.Box()
	if err != nil {
		return
	}

	result, err := t.Map(b, (**vagrant_plugin_sdk.Args_Box)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Box)
	}

	return
}

var (
	_ plugin.Plugin     = (*TargetMachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*TargetMachinePlugin)(nil)
	_ core.Machine      = (*targetMachineClient)(nil)
)
