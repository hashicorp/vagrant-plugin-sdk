package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// CommunicatorPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Communicator component type.
type CommunicatorPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Communicator // Impl is the concrete implementation
	*BasePlugin
}

func (p *CommunicatorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterCommunicatorServiceServer(s, &communicatorServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, p.Impl),
	})
	return nil
}

func (p *CommunicatorPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewCommunicatorServiceClient(c)
	return &communicatorClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

// communicatorClient is an implementation of component.Communicator over gRPC.
type communicatorClient struct {
	*BaseClient

	client vagrant_plugin_sdk.CommunicatorServiceClient
}

func (c *communicatorClient) MatchFunc() interface{} {
	// Get our function specification from the server
	spec, err := c.client.MatchSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	// Mark that it's not a mapper
	spec.Result = nil

	// Create a callback to call the actual function on the server
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Match(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}

		return resp.Match, nil
	}

	// Generate a callable function from our specification
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Match(machine core.Machine) (bool, error) {
	f := c.MatchFunc()
	// Call the function and include our local machine argument
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
	)

	if err != nil {
		return false, err
	}

	// and fin
	return raw.(bool), nil
}

func (c *communicatorClient) InitFunc() interface{} {
	spec, err := c.client.InitSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Init(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Init(machine core.Machine) error {
	f := c.InitFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Typed(c.Ctx),
	)
	return err
}

func (c *communicatorClient) ReadyFunc() interface{} {
	spec, err := c.client.ReadySpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Ready(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Ready, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Ready(machine core.Machine) (bool, error) {
	f := c.ReadyFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) WaitForReadyFunc() interface{} {
	spec, err := c.client.WaitForReadySpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.WaitForReady(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Ready, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) WaitForReady(machine core.Machine, wait int) (bool, error) {
	f := c.WaitForReadyFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(wait),
		argmapper.Typed(c.Ctx),
	)

	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) DownloadFunc() interface{} {
	spec, err := c.client.DownloadSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Download(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Download(machine core.Machine, source, destination string) error {
	f := c.DownloadFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Named("source", source),
		argmapper.Named("destination", destination),
		argmapper.Typed(c.Ctx),
	)
	return err
}

func (c *communicatorClient) UploadFunc() interface{} {
	spec, err := c.client.UploadSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Upload(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err == nil, err
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Upload(machine core.Machine, source, destination string) error {
	f := c.DownloadFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Named("source", source),
		argmapper.Named("destination", destination),
		argmapper.Typed(c.Ctx),
	)
	return err
}

func (c *communicatorClient) ExecuteFunc() interface{} {
	spec, err := c.client.ExecuteSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*core.CommunicatorMessage, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		result, err := c.client.Execute(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		return &core.CommunicatorMessage{ExitCode: result.ExitCode, Stdout: result.Stdout, Stderr: result.Stderr}, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Execute(machine core.Machine, cmd []string, opts ...interface{}) (status int32, err error) {
	f := c.ExecuteFunc()
	raw, err := c.CallDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(opts),
		argmapper.Typed(cmd),
		argmapper.Named("command", cmd),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return -1, err
	}

	return raw.(int32), nil
}

func (c *communicatorClient) PrivilegedExecuteFunc() interface{} {
	spec, err := c.client.PrivilegedExecuteSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*core.CommunicatorMessage, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		result, err := c.client.PrivilegedExecute(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		return &core.CommunicatorMessage{ExitCode: result.ExitCode, Stdout: result.Stdout, Stderr: result.Stderr}, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) PrivilegedExecute(machine core.Machine, cmd []string, opts ...interface{}) (status int32, err error) {
	f := c.PrivilegedExecuteFunc()
	raw, err := c.CallDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(opts),
		argmapper.Typed(cmd),
		argmapper.Named("command", cmd),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return -1, err
	}

	return raw.(int32), nil
}

func (c *communicatorClient) TestFunc() interface{} {
	spec, err := c.client.TestSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		c.Logger.Debug("got args for test func in the communicatorClient %w", args)
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		result, err := c.client.Test(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return result.Valid, nil
	}
	c.Logger.Debug("generated func")
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Test(machine core.Machine, cmd []string, opts ...interface{}) (valid bool, err error) {
	c.Logger.Debug("Running Test from communicatorClient")
	f := c.TestFunc()
	c.Logger.Debug("got test func from communicator client")

	c.Logger.Debug("running with opts %w", opts)
	var optsArgs []interface{}
	if opts == nil {
		// is the opts are empty then pass in an empty args hash
		optsArgs = []interface{}{&vagrant_plugin_sdk.Args_Hash{}}
	} else {
		optsArgs = opts
	}

	c.Logger.Debug("running with optsArgs %w", optsArgs)
	c.Logger.Debug("here we go running")
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(optsArgs...),
		argmapper.Typed(cmd),
		argmapper.Typed(c.Ctx, c.Logger, c.Internal),
	)
	c.Logger.Debug("result of running the dynamic function %w", raw)
	if err != nil {
		c.Logger.Debug("oh boy, there was an error %w", err)
		return false, err
	}

	return raw.(bool), nil
}

func (c *communicatorClient) ResetFunc() interface{} {
	spec, err := c.client.ResetSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Reset(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *communicatorClient) Reset(machine core.Machine) (err error) {
	f := c.ResetFunc()
	_, err = c.CallDynamicFunc(f, false,
		argmapper.Typed(machine),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return err
	}

	return nil
}

// communicatorServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type communicatorServer struct {
	*BaseServer

	Impl component.Communicator
	vagrant_plugin_sdk.UnsafeCommunicatorServiceServer
}

func (s *communicatorServer) ConfigStruct(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *communicatorServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*emptypb.Empty, error) {
	return configure(s.Impl, req)
}

func (s *communicatorServer) Documentation(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *communicatorServer) MatchSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.MatchFunc())
}

func (s *communicatorServer) Match(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_MatchResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.MatchFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_MatchResp{Match: raw.(bool)}, nil
}

func (s *communicatorServer) InitSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.InitFunc())
}

func (s *communicatorServer) Init(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.InitFunc(), false, args.Args,
		argmapper.Typed(ctx),
	)
	return &emptypb.Empty{}, err
}

func (s *communicatorServer) ReadySpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ReadyFunc())
}

func (s *communicatorServer) Ready(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ReadyResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ReadyFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ReadyResp{Ready: raw.(bool)}, nil
}

func (s *communicatorServer) WaitForReadySpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.WaitForReadyFunc())
}

func (s *communicatorServer) WaitForReady(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ReadyResp, error) {

	raw, err := s.CallDynamicFunc(s.Impl.WaitForReadyFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_ReadyResp{Ready: raw.(bool)}, nil
}

func (s *communicatorServer) DownloadSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.DownloadFunc())
}

func (s *communicatorServer) Download(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.DownloadFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *communicatorServer) UploadSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.UploadFunc())
}

func (s *communicatorServer) Upload(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.UploadFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *communicatorServer) ExecuteSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ExecuteFunc())
}

func (s *communicatorServer) Execute(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ExecuteResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ExecuteFunc(), (**core.CommunicatorMessage)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}
	result := raw.(*core.CommunicatorMessage)
	return &vagrant_plugin_sdk.Communicator_ExecuteResp{
		ExitCode: result.ExitCode,
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
	}, nil
}

func (s *communicatorServer) PrivilegedExecuteSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.PrivilegedExecuteFunc())
}

func (s *communicatorServer) PrivilegedExecute(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_ExecuteResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.PrivilegedExecuteFunc(), (**core.CommunicatorMessage)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	result := raw.(*core.CommunicatorMessage)
	return &vagrant_plugin_sdk.Communicator_ExecuteResp{
		ExitCode: result.ExitCode,
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
	}, nil
}

func (s *communicatorServer) TestSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.TestFunc())
}

func (s *communicatorServer) Test(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Communicator_TestResp, error) {
	s.Logger.Debug("Runnin Test from communicator server with args %w", args)
	raw, err := s.CallDynamicFunc(s.Impl.TestFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Communicator_TestResp{Valid: raw.(bool)}, nil
}

func (s *communicatorServer) ResetSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "communicator"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ResetFunc())
}

func (s *communicatorServer) Reset(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.ResetFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

var (
	_ plugin.Plugin                                = (*CommunicatorPlugin)(nil)
	_ plugin.GRPCPlugin                            = (*CommunicatorPlugin)(nil)
	_ vagrant_plugin_sdk.CommunicatorServiceServer = (*communicatorServer)(nil)
	_ component.Communicator                       = (*communicatorClient)(nil)
	_ core.Seeder                                  = (*communicatorClient)(nil)
	_ core.Communicator                            = (*communicatorClient)(nil)
)
