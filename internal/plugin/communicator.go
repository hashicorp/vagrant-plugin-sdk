package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// CommunicatorPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Communicator component type.
type CommunicatorPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Communicator // Impl is the concrete implementation
	Mappers []*argmapper.Func      // Mappers
	Logger  hclog.Logger           // Logger
}

func (p *CommunicatorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterCommunicatorServiceServer(s, &communicatorServer{
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

func (p *CommunicatorPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &communicatorClient{
		client: vagrant_plugin_sdk.NewCommunicatorServiceClient(c),
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

// communicatorClient is an implementation of component.Communicator over gRPC.
type communicatorClient struct {
	*baseClient

	client vagrant_plugin_sdk.CommunicatorServiceClient
}

func (c *communicatorClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *communicatorClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *communicatorClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *communicatorClient) MatchFunc() interface{} {
	// Get our function specification from the server
	spec, err := c.client.MatchSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	// Mark that it's not a mapper
	spec.Result = nil

	// Create a callback to call the actual function on the server
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Match(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}

		return resp.Match, nil
	}

	// Generate a callable function from our specification
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Match(machine core.Machine) (bool, error) {
	f := c.MatchFunc()
	// Call the function and include our local machine argument
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.ctx),
		argmapper.Typed(machine),
	)

	if err != nil {
		return false, err
	}

	// and fin
	return raw.(bool), nil
}

func (c *communicatorClient) InitFunc() interface{} {
	spec, err := c.client.InitSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		_, err := c.client.Init(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Init(machine core.Machine) error {
	f := c.InitFunc()
	_, err := c.callDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Typed(c.ctx),
	)
	return err
}

func (c *communicatorClient) ReadyFunc() interface{} {
	spec, err := c.client.ReadySpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Ready(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Ready, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Ready(machine core.Machine) (bool, error) {
	f := c.ReadyFunc()
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) WaitForReadyFunc() interface{} {
	spec, err := c.client.WaitForReadySpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.WaitForReady(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Ready, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) WaitForReady(machine core.Machine, wait int) (bool, error) {
	f := c.WaitForReadyFunc()
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(wait),
		argmapper.Typed(c.ctx),
	)

	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) DownloadFunc() interface{} {
	spec, err := c.client.DownloadSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		_, err := c.client.Download(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Download(machine core.Machine, source, destination string) error {
	f := c.DownloadFunc()
	_, err := c.callDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Named("source", source),
		argmapper.Named("destination", destination),
		argmapper.Typed(c.ctx),
	)
	return err
}

func (c *communicatorClient) UploadFunc() interface{} {
	spec, err := c.client.UploadSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		_, err := c.client.Upload(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Upload(machine core.Machine, source, destination string) error {
	f := c.DownloadFunc()
	_, err := c.callDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Named("source", source),
		argmapper.Named("destination", destination),
		argmapper.Typed(c.ctx),
	)
	return err
}

func (c *communicatorClient) ExecuteFunc() interface{} {
	spec, err := c.client.ExecuteSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int32, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		result, err := c.client.Execute(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return -1, err
		}
		return result.ExitCode, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Execute(machine core.Machine, cmd []string, opts *core.CommunicatorOptions) (status int32, err error) {
	f := c.ExecuteFunc()
	raw, err := c.callDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(opts),
		argmapper.Typed(cmd),
		argmapper.Named("command", cmd),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return -1, err
	}

	return raw.(int32), nil
}

func (c *communicatorClient) PrivilegedExecuteFunc() interface{} {
	spec, err := c.client.PrivilegedExecuteSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int32, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		result, err := c.client.PrivilegedExecute(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return -1, err
		}
		return result.ExitCode, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) PrivilegedExecute(machine core.Machine, cmd []string, opts *core.CommunicatorOptions) (status int32, err error) {
	f := c.PrivilegedExecuteFunc()
	raw, err := c.callDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(opts),
		argmapper.Typed(cmd),
		argmapper.Named("command", cmd),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return -1, err
	}

	return raw.(int32), nil
}

func (c *communicatorClient) TestFunc() interface{} {
	spec, err := c.client.TestSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		result, err := c.client.Test(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return result.Valid, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Test(machine core.Machine, cmd []string, opts *core.CommunicatorOptions) (valid bool, err error) {
	f := c.TestFunc()
	raw, err := c.callDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(opts),
		argmapper.Typed(cmd),
		argmapper.Named("command", cmd),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) ResetFunc() interface{} {
	spec, err := c.client.ResetSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		_, err := c.client.Reset(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *communicatorClient) Reset(machine core.Machine) (err error) {
	f := c.ResetFunc()
	_, err = c.callDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return err
	}

	return nil
}

// communicatorServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type communicatorServer struct {
	*baseServer

	Impl component.Communicator
	vagrant_plugin_sdk.UnimplementedCommunicatorServiceServer
}

func (s *communicatorServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *communicatorServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *communicatorServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *communicatorServer) MatchSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.MatchFunc())
}

func (s *communicatorServer) Match(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_MatchResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.MatchFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_MatchResp{Match: raw.(bool)}, nil
}

func (s *communicatorServer) InitSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.InitFunc())
}

func (s *communicatorServer) Init(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_InitResp, error) {
	_, err := s.callDynamicFunc(s.Impl.InitFunc(), false, args.Args,
		argmapper.Typed(ctx),
	)
	return &vagrant_plugin_sdk.Communicator_InitResp{}, err
}

func (s *communicatorServer) ReadySpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ReadyFunc())
}

func (s *communicatorServer) Ready(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ReadyResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.ReadyFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ReadyResp{Ready: raw.(bool)}, nil
}

func (s *communicatorServer) WaitForReadySpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.WaitForReadyFunc())
}

func (s *communicatorServer) WaitForReady(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ReadyResp, error) {

	raw, err := s.callDynamicFunc(s.Impl.WaitForReadyFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ReadyResp{Ready: raw.(bool)}, nil
}

func (s *communicatorServer) DownloadSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.DownloadFunc())
}

func (s *communicatorServer) Download(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_FileTransferResp, error) {
	_, err := s.callDynamicFunc(s.Impl.DownloadFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_FileTransferResp{}, nil
}

func (s *communicatorServer) UploadSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.UploadFunc())
}

func (s *communicatorServer) Upload(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_FileTransferResp, error) {
	_, err := s.callDynamicFunc(s.Impl.UploadFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_FileTransferResp{}, nil
}

func (s *communicatorServer) ExecuteSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ExecuteFunc())
}

func (s *communicatorServer) Execute(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ExecuteResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.ExecuteFunc(), (*int32)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ExecuteResp{ExitCode: raw.(int32)}, nil
}

func (s *communicatorServer) PrivilegedExecuteSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.PrivilegedExecuteFunc())
}

func (s *communicatorServer) PrivilegedExecute(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ExecuteResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.PrivilegedExecuteFunc(), (*int32)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ExecuteResp{ExitCode: raw.(int32)}, nil
}

func (s *communicatorServer) TestSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.TestFunc())
}

func (s *communicatorServer) Test(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_TestResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.TestFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_TestResp{Valid: raw.(bool)}, nil
}

func (s *communicatorServer) ResetSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ResetFunc())
}

func (s *communicatorServer) Reset(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ResetResp, error) {
	_, err := s.callDynamicFunc(s.Impl.ResetFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ResetResp{}, nil
}

var (
	_ plugin.Plugin                                = (*CommunicatorPlugin)(nil)
	_ plugin.GRPCPlugin                            = (*CommunicatorPlugin)(nil)
	_ vagrant_plugin_sdk.CommunicatorServiceServer = (*communicatorServer)(nil)
	_ component.Communicator                       = (*communicatorClient)(nil)
)
