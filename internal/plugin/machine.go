package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// MachinePlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Machine component type.
type MachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    core.Machine      // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *MachinePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterMachineServiceServer(s, &machineServer{
		baseServer: &baseServer{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	})
	return nil
}

func (p *MachinePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &machineClient{
		baseClient: &baseClient{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// machineClient is an implementation of core.Machine over gRPC.
type machineClient struct {
	*baseClient

	client proto.MachineServiceClient
}

func (c *machineClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *machineClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *machineClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *machineClient) ActionFunc() interface{} {
	spec, err := c.client.ActionSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Action(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) Action(machine core.Machine) (bool, error) {
	f := c.ActionFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) CommunicateFunc() interface{} {
	spec, err := c.client.CommunicateSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Communicate(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) Communicate(machine core.Machine) (bool, error) {
	f := c.CommunicateFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) GuestFunc() interface{} {
	spec, err := c.client.GuestSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Guest(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) Guest(machine core.Machine) (bool, error) {
	f := c.GuestFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) SetIDFunc() interface{} {
	spec, err := c.client.SetIDSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.SetID(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) SetID(machine core.Machine) (bool, error) {
	f := c.SetIDFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) IndexUUIDFunc() interface{} {
	spec, err := c.client.IndexUUIDSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.IndexUUID(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) IndexUUID(machine core.Machine) (bool, error) {
	f := c.IndexUUIDFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) InspectFunc() interface{} {
	spec, err := c.client.InspectSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Inspect(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) Inspect(machine core.Machine) (bool, error) {
	f := c.InspectFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) ReloadFunc() interface{} {
	spec, err := c.client.ReloadSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Reload(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) Reload(machine core.Machine) (bool, error) {
	f := c.ReloadFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) ConnectionInfoFunc() interface{} {
	spec, err := c.client.ConnectionInfoSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.ConnectionInfo(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) ConnectionInfo(machine core.Machine) (bool, error) {
	f := c.ConnectionInfoFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) StateFunc() interface{} {
	spec, err := c.client.StateSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.State(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) State(machine core.Machine) (bool, error) {
	f := c.StateFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) UIDFunc() interface{} {
	spec, err := c.client.UIDSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.UID(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) UID(machine core.Machine) (bool, error) {
	f := c.UIDFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *machineClient) SyncedFoldersFunc() interface{} {
	spec, err := c.client.SyncedFoldersSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.SyncedFolders(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *machineClient) SyncedFolders(machine core.Machine) (bool, error) {
	f := c.SyncedFoldersFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// machineServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type machineServer struct {
	*baseServer

	Impl core.Machine
}

func (s *machineServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *machineServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *machineServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *machineServer) ActionSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ActionFunc())
}

func (s *machineServer) Action(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_ActionResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.ActionFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_ActionResp{Success: raw.(bool)}, nil
}

func (s *machineServer) CommunicateSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.CommunicateFunc())
}

func (s *machineServer) Communicate(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_CommunicateResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.CommunicateFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_CommunicateResp{Communicator: raw.(*proto.Communicator)}, nil
}

func (s *machineServer) GuestSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.GuestFunc())
}

func (s *machineServer) Guest(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_GuestResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.GuestFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_GuestResp{Guest: raw.(*proto.Guest)}, nil
}

func (s *machineServer) SetIDSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.SetIDFunc())
}

func (s *machineServer) SetID(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_SetIDResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.SetIDFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_SetIDResp{Success: raw.(bool)}, nil
}

func (s *machineServer) IndexUUIDSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.IndexUUIDFunc())
}

func (s *machineServer) IndexUUID(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_IndexUUIDResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.IndexUUIDFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_IndexUUIDResp{Uuid: raw.(string)}, nil
}

func (s *machineServer) InspectSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.InspectFunc())
}

func (s *machineServer) Inspect(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_InspectResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.InspectFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_InspectResp{Inspect: raw.(string)}, nil
}

func (s *machineServer) ReloadSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ReloadFunc())
}

func (s *machineServer) Reload(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_ReloadResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.ReloadFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_ReloadResp{Id: raw.(string)}, nil
}

func (s *machineServer) ConnectionInfoSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ConnectionInfoFunc())
}

func (s *machineServer) ConnectionInfo(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_ConnectionInfoResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.ConnectionInfoFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	// Maybe needs to be a Machine_ConnectionInfoResp_Winrm?
	return &proto.Machine_ConnectionInfoResp{Connection: raw.(*proto.Machine_ConnectionInfoResp_Ssh)}, nil
}

func (s *machineServer) StateSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.StateFunc())
}

func (s *machineServer) State(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_StateResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.StateFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_StateResp{State: raw.(*proto.MachineState)}, nil
}

func (s *machineServer) UIDSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.UIDFunc())
}

func (s *machineServer) UID(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_UIDResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.UIDFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_UIDResp{Uid: raw.(string)}, nil
}

func (s *machineServer) SyncedFoldersSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "machine"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.SyncedFoldersFunc())
}

func (s *machineServer) SyncedFolders(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Machine_SyncedFoldersResp, error) {
	_, err := s.callLocalDynamicFunc(s.Impl.SyncedFoldersFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Machine_SyncedFoldersResp{}, nil
}

var (
	_ plugin.Plugin              = (*MachinePlugin)(nil)
	_ plugin.GRPCPlugin          = (*MachinePlugin)(nil)
	_ proto.MachineServiceServer = (*machineServer)(nil)
	_ core.Machine               = (*machineClient)(nil)
)
