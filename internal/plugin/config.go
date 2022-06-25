package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/protostructure"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// ConfigPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Config component type.
type ConfigPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Config // Impl is the concrete implementation
	*BasePlugin
}

func (p *ConfigPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterConfigServiceServer(s,
		&configServer{
			Impl:       p.Impl,
			BaseServer: p.NewServer(broker, p.Impl),
		},
	)
	return nil
}

func (p *ConfigPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewConfigServiceClient(c)
	return &configClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl),
	}, nil
}

// configClient is an implementation of component.Config over gRPC.
type configClient struct {
	*BaseClient

	client vagrant_plugin_sdk.ConfigServiceClient
}

func (c *configClient) Register() (*component.ConfigRegistration, error) {
	r, err := c.client.Register(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &component.ConfigRegistration{
		Identifier: r.Identifier,
		Scope:      r.Scope,
	}, nil
}

func (c *configClient) StructFunc() interface{} {
	spec, err := c.client.StructSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		ctx, _ = c.GenerateContext(ctx)
		resp, err := c.client.Struct(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}
		switch v := resp.GetValue().(type) {
		case *vagrant_plugin_sdk.Config_StructResponse_Raw:
			return true, nil
		case *vagrant_plugin_sdk.Config_StructResponse_Struct:
			return v.Struct.Struct, nil
		}

		return fmt.Errorf("unknown config struct response"), nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *configClient) Struct() (interface{}, error) {
	f := c.StructFunc()
	raw, err := c.CallDynamicFunc(f, false, argmapper.Typed(c.Ctx))
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (c *configClient) MergeFunc() interface{} {
	spec, err := c.client.MergeSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Args_ConfigData, error) {
		ctx, _ = c.GenerateContext(ctx)
		result, err := c.client.Merge(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *configClient) Merge(base, toMerge *component.ConfigData) (*component.ConfigData, error) {
	f := c.MergeFunc()
	m := &component.ConfigMerge{
		Base:    base,
		Overlay: toMerge,
	}

	raw, err := c.CallDynamicFunc(f, (**component.ConfigData)(nil),
		argmapper.Typed(c.Ctx, m),
	)
	if err != nil {
		return nil, err
	}

	return raw.(*component.ConfigData), nil
}

func (c *configClient) FinalizeFunc() interface{} {
	spec, err := c.client.FinalizeSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Args_ConfigData, error) {
		ctx, _ = c.GenerateContext(ctx)
		resp, err := c.client.Finalize(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *configClient) Finalize(data *component.ConfigData) (*component.ConfigData, error) {
	f := c.FinalizeFunc()
	raw, err := c.CallDynamicFunc(f, (**component.ConfigData)(nil),
		argmapper.Typed(&component.ConfigFinalize{Config: data}, c.Ctx))
	if err != nil {
		return nil, err
	}

	return raw.(*component.ConfigData), nil
}

// func (c *configClient) Configure(fn core.ConfigFn) error {
// 	r, err := c.client.ConfigStruct(c.Ctx, &emptypb.Empty{})
// 	if err != nil {
// 		return err
// 	}

// 	var s interface{}

// 	if r.Struct != nil {
// 		s, err = protostructure.New(r.Struct)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		s = r.Fields
// 	}

// 	result, err := fn(s)
// 	if err != nil {
// 		return err
// 	}

// 	rj, err := json.Marshal(result)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = c.client.Configure(c.Ctx, &vagrant_plugin_sdk.Config_ConfigureRequest{
// 		Json: rj,
// 	})

// 	return err
// }

// configServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type configServer struct {
	*BaseServer

	Impl component.Config
	vagrant_plugin_sdk.UnsafeConfigServiceServer
}

func (s *configServer) Register(
	ctx context.Context,
	_ *emptypb.Empty,
) (resp *vagrant_plugin_sdk.Config_RegisterResponse, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to receive config plugin registration",
				"error", err,
			)
		}
	}()
	n, err := s.Impl.Register()
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Config_RegisterResponse{
		Identifier: n.Identifier,
		Scope:      n.Scope,
	}, nil
}

func (s *configServer) StructSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (result *vagrant_plugin_sdk.FuncSpec, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to generate config struct spec",
				"error", err,
			)
		}
	}()
	if err = isImplemented(s, "config"); err != nil {
		return
	}
	result, err = s.GenerateSpec(s.Impl.StructFunc())

	return
}

func (s *configServer) Struct(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (result *vagrant_plugin_sdk.Config_StructResponse, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to generate config struct",
				"error", err,
			)
		}
	}()
	raw, err := s.CallDynamicFunc(s.Impl.StructFunc(),
		false, req.Args, argmapper.Typed(ctx))
	if err != nil {
		return
	}
	result = &vagrant_plugin_sdk.Config_StructResponse{}
	switch v := raw.(type) {
	case bool:
		result.Value = &vagrant_plugin_sdk.Config_StructResponse_Raw{
			Raw: true,
		}
	default:
		var val *protostructure.Struct
		val, err = protostructure.Encode(v)
		result.Value = &vagrant_plugin_sdk.Config_StructResponse_Struct{
			Struct: &vagrant_plugin_sdk.Config_Structure{
				Struct: val,
			},
		}
	}

	return
}

func (s *configServer) MergeSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (result *vagrant_plugin_sdk.FuncSpec, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to generate config merge spec",
				"error", err,
			)
		}
	}()
	if err = isImplemented(s, "config"); err != nil {
		return
	}
	result, err = s.GenerateSpec(s.Impl.MergeFunc())
	for _, i := range result.Args {
		s.Logger.Info("go merge spec argument",
			"name", i.Name,
			"type", i.Type,
		)
	}

	return
}

func (s *configServer) Merge(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Args_ConfigData, error) {
	for _, i := range req.Args {
		s.Logger.Info("config merge go plugin argument",
			"name", i.Name,
			"type", i.Type,
			"value", i.Value,
		)
	}
	s.Logger.Info("running config merge on go config plugin",
		"args", hclog.Fmt("%#v", req.Args),
	)
	raw, err := s.CallDynamicFunc(s.Impl.MergeFunc(),
		(**vagrant_plugin_sdk.Args_ConfigData)(nil), req.Args, argmapper.Typed(ctx))
	if err != nil {
		s.Logger.Error("failed to merge config",
			"error", err,
		)
		return nil, err
	}

	s.Logger.Info("result of config mrege on go config plugin",
		"config", hclog.Fmt("%#v", raw.(*vagrant_plugin_sdk.Args_ConfigData).Data),
	)

	return raw.(*vagrant_plugin_sdk.Args_ConfigData), nil
}

func (s *configServer) FinalizeSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (result *vagrant_plugin_sdk.FuncSpec, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to generate finalize spec",
				"error", err,
			)
		}
	}()
	if err = isImplemented(s, "config"); err != nil {
		return
	}
	result, err = s.GenerateSpec(s.Impl.FinalizeFunc())
	for _, i := range result.Args {
		s.Logger.Info("go finalize spec argument",
			"name", i.Name,
			"type", i.Type,
		)
	}

	s.Logger.Info("finalize spec generated for config",
		"spec", hclog.Fmt("%#v", result),
	)

	return
}

func (s *configServer) Finalize(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (result *vagrant_plugin_sdk.Args_ConfigData, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to finalize config",
				"error", err,
			)
		}
	}()
	for _, i := range req.Args {
		s.Logger.Info("config finalize go plugin argument",
			"name", i.Name,
			"type", i.Type,
			"value", i.Value,
		)
	}
	s.Logger.Info("running config finalize on go config plugin",
		"args", hclog.Fmt("%#v", req.Args),
	)

	raw, err := s.CallDynamicFunc(s.Impl.FinalizeFunc(),
		(**vagrant_plugin_sdk.Args_ConfigData)(nil), req.Args, argmapper.Typed(ctx))
	if err != nil {
		return
	}

	return raw.(*vagrant_plugin_sdk.Args_ConfigData), nil
}

var (
	_ plugin.Plugin                          = (*ConfigPlugin)(nil)
	_ plugin.GRPCPlugin                      = (*ConfigPlugin)(nil)
	_ vagrant_plugin_sdk.ConfigServiceServer = (*configServer)(nil)
	_ component.Config                       = (*configClient)(nil)
	_ core.Config                            = (*configClient)(nil)
	_ core.Seeder                            = (*configClient)(nil)
)
