package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// HostPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Host component type.
type HostPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Host    // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *HostPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterHostServiceServer(s, &hostServer{
		Impl: p.Impl,
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

func (p *HostPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &hostClient{
		client: proto.NewHostServiceClient(c),
		baseClient: &baseClient{
			ctx: context.Background(),
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// hostClient is an implementation of component.Host over gRPC.
type hostClient struct {
	*baseClient

	client proto.HostServiceClient
}

func (c *hostClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *hostClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *hostClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *hostClient) DetectFunc() interface{} {
	spec, err := c.client.DetectSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.Detect(ctx, &proto.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}

		return resp.Detected, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *hostClient) Detect() (bool, error) {
	f := c.DetectFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*bool)(nil), f)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

// hostServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type hostServer struct {
	*baseServer

	Impl component.Host
}

func (s *hostServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *hostServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *hostServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *hostServer) DetectSpec(
	ctx context.Context,
	args *empty.Empty,
) (*proto.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.DetectFunc())
}

func (s *hostServer) Detect(
	ctx context.Context,
	args *proto.FuncSpec_Args,
) (*proto.Host_DetectResp, error) {
	s.Logger.Debug("running the detect function on the server to call real implementation")
	raw, err := s.callLocalDynamicFunc(s.Impl.DetectFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &proto.Host_DetectResp{Detected: raw.(bool)}, nil
}

var (
	_ plugin.Plugin           = (*HostPlugin)(nil)
	_ plugin.GRPCPlugin       = (*HostPlugin)(nil)
	_ proto.HostServiceServer = (*hostServer)(nil)
	_ component.Host          = (*hostClient)(nil)
)
