package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/multistep"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// ProviderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provider component type.
type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Provider // Impl is the concrete implementation
	Mappers []*argmapper.Func  // Mappers
	Logger  hclog.Logger       // Logger
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServiceServer(s, &providerServer{
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

func (p *ProviderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &providerClient{
		baseClient: &baseClient{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// providerClient is an implementation of component.Provider over gRPC.
type providerClient struct {
	*baseClient

	client proto.ProviderServiceClient
}

func (c *providerClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *providerClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *providerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *providerClient) UsableFunc() interface{} {
	spec, err := c.client.UsableSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.Usable(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.IsUsable, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *providerClient) Usable() (bool, error) {
	f := c.UsableFunc()
	raw, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *providerClient) InitFunc() interface{} {
	spec, err := c.client.InitSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Init(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *providerClient) Init(machine core.Machine) (bool, error) {
	f := c.UsableFunc()
	_, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f,
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *providerClient) InstalledFunc() interface{} {
	spec, err := c.client.InstalledSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.Installed(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.IsInstalled, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *providerClient) Installed() (bool, error) {
	f := c.InstalledFunc()
	raw, err := c.callRemoteDynamicFunc(context.Background(), nil, (*bool)(nil), f)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *providerClient) ActionUpFunc() interface{} {
	spec, err := c.client.ActionUpSpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		resp, err := c.client.ActionUp(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Result, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *providerClient) ActionUp(machine core.Machine, state multistep.StateBag) (interface{}, error) {
	f := c.ActionUpFunc()
	raw, err := c.callRemoteDynamicFunc(context.Background(), nil, (interface{})(nil), f,
		argmapper.Typed(machine),
		argmapper.Typed(state),
	)
	if err != nil {
		return false, err
	}

	return raw, nil
}

// providerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type providerServer struct {
	*baseServer

	Impl component.Provider
}

func (s *providerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *providerServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *providerServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *providerServer) UsableSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.UsableFunc())
}

func (s *providerServer) Usable(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Provider_UsableResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.UsableFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Provider_UsableResp{IsUsable: raw.(bool)}, nil
}

func (s *providerServer) InstalledSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.InstalledFunc())
}

func (s *providerServer) Installed(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Provider_InstalledResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.InstalledFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Provider_InstalledResp{IsInstalled: raw.(bool)}, nil
}

func (s *providerServer) ActionUpSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if s.Impl == nil {
		return nil, status.Errorf(codes.Unimplemented, "plugin does not implement: provider")
	}

	return funcspec.Spec(s.Impl.ActionUpFunc(),
		argmapper.Logger(s.Logger),
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(s.internal()),
	)
	// if err := isImplemented(s, "provider"); err != nil {
	// 	return nil, err
	// }

	// return s.generateSpec(s.Impl.ActionUpFunc())
}

func (s *providerServer) ActionUp(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Provider_ActionResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.ActionUpFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Provider_ActionResp{Result: raw.(*anypb.Any)}, nil
}

func (s *providerServer) InitSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.InitFunc())
}

func (s *providerServer) Init(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*empty.Empty, error) {
	_, err := s.callLocalDynamicFunc(s.Impl.InitFunc(), args.Args, (*interface{})(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

var (
	_ plugin.Plugin               = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin           = (*ProviderPlugin)(nil)
	_ proto.ProviderServiceServer = (*providerServer)(nil)
	_ component.Provider          = (*providerClient)(nil)
)
