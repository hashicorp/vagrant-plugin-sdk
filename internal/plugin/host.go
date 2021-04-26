package plugin

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
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
	vagrant_plugin_sdk.RegisterHostServiceServer(s, &hostServer{
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
		client: vagrant_plugin_sdk.NewHostServiceClient(c),
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

	client vagrant_plugin_sdk.HostServiceClient
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
		resp, err := c.client.Detect(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
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

func (c *hostClient) HasCapability(capName string) bool {
	resp, err := c.client.HasCapability(
		c.ctx,
		&vagrant_plugin_sdk.Host_Capability_NamedRequest{
			Name: capName,
		},
	)
	if err != nil {
		return false
	}
	return resp.HasCapability
}

func (c *hostClient) Capability(capName string, args ...argmapper.Arg) (interface{}, error) {
	resp, err := c.client.Capability(
		c.ctx,
		&vagrant_plugin_sdk.Host_Capability_NamedRequest{
			Name: capName,
			// TODO: Insert args here
			// FuncArgs: args,
		},
	)

	// TODO: do something to result here?
	return resp.Result, err
}

func (c *hostClient) InitializeCapabilities() (err error) {
	_, err = c.client.InitializeCapabilities(c.ctx, &empty.Empty{})
	return
}

// hostServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type hostServer struct {
	*baseServer

	Impl component.Host
	vagrant_plugin_sdk.UnimplementedHostServiceServer
}

func (s *hostServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *hostServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *hostServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *hostServer) DetectSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.DetectFunc())
}

func (s *hostServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Host_DetectResp, error) {
	s.Logger.Debug("running the detect function on the server to call real implementation")
	raw, err := s.callLocalDynamicFunc(s.Impl.DetectFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Host_DetectResp{Detected: raw.(bool)}, nil
}

func (s *hostServer) HasCapability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Host_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Host_Capability_CheckResp, error) {
	result := s.Impl.HasCapability(args.Name)
	return &vagrant_plugin_sdk.Host_Capability_CheckResp{HasCapability: result}, nil
}

func (s *hostServer) Capability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Host_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Host_Capability_Resp, error) {
	// TODO: pass this args
	hasCap := s.Impl.HasCapability(args.Name)
	if hasCap == false {
		return nil, errors.New("Capability " + args.Name + " not found")
	}
	result, err := s.Impl.Capability(args.Name)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Host_Capability_Resp{Result: result.(*anypb.Any)}, nil
}

func (s *hostServer) InitializeCapabilities(
	ctx context.Context,
	_ *empty.Empty,
) (*empty.Empty, error) {
	err := s.Impl.InitializeCapabilities()
	return &empty.Empty{}, err
}

var (
	_ plugin.Plugin                        = (*HostPlugin)(nil)
	_ plugin.GRPCPlugin                    = (*HostPlugin)(nil)
	_ vagrant_plugin_sdk.HostServiceServer = (*hostServer)(nil)
	_ component.Host                       = (*hostClient)(nil)
)
