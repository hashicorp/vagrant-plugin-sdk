// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/protostructure"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/errors"
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

func (c *configClient) InitFunc() interface{} {
	spec, err := c.client.InitSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		ctx, _ = c.GenerateContext(ctx)
		resp, err := c.client.Init(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *configClient) Init(in *component.ConfigData) (*component.ConfigData, error) {
	f := c.InitFunc()

	// NOTE: Need to map result directly to prevent invalid result
	raw1, err := c.CallDynamicFunc(f, false, argmapper.Typed(c.Ctx, in))
	if err != nil {
		return nil, errors.Wrap("config init failed", err)
	}

	raw, err := c.Map(raw1, (**component.ConfigData)(nil), argmapper.Typed(c.Ctx))
	if err != nil {
		return nil, errors.Wrap("config init map failed", err)
	}

	return raw.(*component.ConfigData), nil
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
	cb := func(ctx context.Context, args funcspec.Args) (*component.ConfigData, error) {
		ctx, _ = c.GenerateContext(ctx)
		result, err := c.client.Merge(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}
		raw, err := c.Map(result, (**component.ConfigData)(nil), argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		return raw.(*component.ConfigData), err
	}

	return c.GenerateFunc(spec, cb)
}

func (c *configClient) Merge(base, overlay *component.ConfigData) (*component.ConfigData, error) {
	f := c.MergeFunc()
	baseProto, err := c.Map(
		base,
		(**vagrant_plugin_sdk.Args_ConfigData)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		c.Logger.Error("failed to convert base to proto",
			"error", err,
		)
	}

	overlayProto, err := c.Map(
		overlay,
		(**vagrant_plugin_sdk.Args_ConfigData)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		c.Logger.Error("failed to convert base to proto",
			"error", err,
		)
	}

	raw, err := c.CallDynamicFunc(f,
		(**component.ConfigData)(nil),
		argmapper.Typed(c.Ctx),
		argmapper.Named("Base", baseProto),
		argmapper.Named("Overlay", overlayProto),
		argmapper.Typed(&vagrant_plugin_sdk.Config_Merge{
			Base:    baseProto.(*vagrant_plugin_sdk.Args_ConfigData),
			Overlay: overlayProto.(*vagrant_plugin_sdk.Args_ConfigData),
		}),
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
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Config_FinalizeResponse, error) {
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
	c.Logger.Warn("Running config finalize from the config client")
	r, err := c.CallDynamicFunc(f, (**vagrant_plugin_sdk.Config_FinalizeResponse)(nil),
		argmapper.Typed(data, c.Ctx))
	if err != nil {
		return nil, err
	}

	raw, err := c.Map(r.(*vagrant_plugin_sdk.Config_FinalizeResponse).Data, (**component.ConfigData)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw.(*component.ConfigData), nil
}

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
	log := s.Logger.With("function", "MergeSpec")
	defer func() {
		if err != nil {
			log.Error("failed to generate config merge spec",
				"error", err,
			)
		}
	}()
	if err = isImplemented(s, "config"); err != nil {
		return
	}

	result, err = s.generateConfigSpec(ctx, s.Impl.FinalizeFunc(), log)
	if err != nil {
		log.Error("spec generation", "error", err)
		return nil, err
	}

	// Add the finalize proto to the request arguments
	result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
		Type: "hashicorp.vagrant.sdk.Config.Merge",
	})

	return
}

func (s *configServer) Merge(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Args_ConfigData, error) {
	log := s.Logger.With("function", "Merge")
	// Grab configuration structures for base and overlay
	base, err := s.CallDynamicFunc(
		s.Impl.StructFunc(),
		false,
		funcspec.Args{},
		argmapper.Typed(ctx),
	)
	if err != nil {
		log.Error("configuration structure retrieval", "error", err)
		return nil, err
	}
	overlay, err := s.CallDynamicFunc(
		s.Impl.StructFunc(),
		false,
		funcspec.Args{},
		argmapper.Typed(ctx),
	)
	if err != nil {
		log.Error("configuration structure retrieval", "error", err)
		return nil, err
	}

	// Extract the merge data
	data, err := s.CallDynamicFunc(
		func(in *component.ConfigMerge) *component.ConfigMerge { return in },
		false, req.Args, argmapper.Typed(ctx),
	)
	if err != nil {
		log.Error("fetch merge data", "error", err)
		return nil, err
	}
	mergeData := data.(*component.ConfigMerge)

	// Decode into custom configuration types
	if mergeData.Base != nil {
		if err = mapstructure.Decode(mergeData.Base.Data, base); err != nil {
			log.Error("base config decode", "error", err)
			return nil, err
		}
	}
	if mergeData.Overlay != nil {
		if err = mapstructure.Decode(mergeData.Overlay.Data, overlay); err != nil {
			log.Error("overlay config decode", "error", err)
			return nil, err
		}
	}

	// Include the custom types as named arguments
	// when calling the merge function
	raw, err := s.CallDynamicFunc(
		s.Impl.MergeFunc(),
		false,
		req.Args,
		argmapper.Typed(ctx),
		argmapper.Named("Base", base),
		argmapper.Named("Overlay", overlay),
	)
	if err != nil {
		log.Error("merge execution", "error", err)
		return nil, err
	}

	// The value returned from finalize can be a *component.ConfigData _or_ it can be
	// the actual configuration structure. For the latter, it will need to be converted
	// to a *component.ConfigData.
	confData, ok := raw.(*component.ConfigData)
	if !ok {
		mapData := structs.Map(raw)
		confData = &component.ConfigData{
			Source: structs.Name(raw),
			Data:   mapData,
		}
	}

	// Map the value into a proto
	confProto, err := s.Map(confData,
		(**vagrant_plugin_sdk.Args_ConfigData)(nil),
		argmapper.Typed(ctx),
	)
	if err != nil {
		log.Error("configuration to proto mapping", "error", err)
		return nil, err
	}

	log.Trace("merged configuration", "proto", confProto)

	return confProto.(*vagrant_plugin_sdk.Args_ConfigData), nil
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

	// Now generate the spec with the customized function
	result, err = s.generateConfigSpec(
		ctx, s.Impl.FinalizeFunc(), s.Logger.With("function", "FinalizeSpec"))
	if err != nil {
		s.Logger.Error("finalize spec generation", "error", err)
		return
	}

	// Check spec arguments to verify the inclusion of the ConfigData
	// proto which is the actual configuration data. If it is not included,
	// manually add it.
	addDataArg := true
	for _, i := range result.Args {
		if i.Type == "hashicorp.vagrant.sdk.Args.ConfigData" {
			addDataArg = false
		}
	}

	if addDataArg {
		result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
			Type: "hashicorp.vagrant.sdk.Args.ConfigData",
		})
	}

	return
}

func (s *configServer) Finalize(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (result *vagrant_plugin_sdk.Config_FinalizeResponse, err error) {
	defer func() {
		if err != nil {
			s.Logger.Error("failed to finalize config",
				"error", err,
			)
		}
	}()
	// Start with grabbing the configuration structure from the
	// plugin so it can be populated with the content
	sConfig, err := s.CallDynamicFunc(
		s.Impl.StructFunc(),
		false,
		req.Args,
		argmapper.Typed(ctx),
	)
	if err != nil {
		s.Logger.Error("unable to retrieve configuration structure", "error", err)
		return nil, err
	}

	s.Logger.Trace("configuration structure", "type", hclog.Fmt("%T", sConfig))

	// Extract the configuration data from the request arguments
	data, err := s.CallDynamicFunc(
		func(in *component.ConfigData) *component.ConfigData { return in },
		false, req.Args, argmapper.Typed(ctx),
	)
	if err != nil {
		s.Logger.Error("unable to extract configuration data", "error", err)
		return nil, err
	}

	s.Logger.Trace("configuration data extracted", "data", data)

	// Decode the configuration data into the configuration structure
	if err = mapstructure.Decode(data.(*component.ConfigData).Data, sConfig); err != nil {
		s.Logger.Error("configuration decoding to struct", "error", err)
		return nil, err
	}

	s.Logger.Trace("decoded configuration into struct", "decoded", sConfig)

	// Now execute the finalize function including the native configuration type
	raw, err := s.CallDynamicFunc(s.Impl.FinalizeFunc(),
		false, req.Args, argmapper.Typed(ctx), argmapper.Typed(sConfig))
	if err != nil {
		s.Logger.Error("finalization function execution", "error", err)
		return nil, err
	}

	s.Logger.Trace("configuration now finalized", "data", hclog.Fmt("%#v", raw))

	// The value returned from finalize can be a *component.ConfigData _or_ it can be
	// the actual configuration structure. For the latter, it will need to be converted
	// to a *component.ConfigData.
	confData, ok := raw.(*component.ConfigData)
	if !ok {
		mapData := structs.Map(raw)
		confData = &component.ConfigData{
			Source: structs.Name(raw),
			Data:   mapData,
		}
	}

	// Map the value into a proto
	confProto, err := s.Map(confData,
		(**vagrant_plugin_sdk.Args_ConfigData)(nil),
		argmapper.Typed(ctx),
	)
	if err != nil {
		s.Logger.Error("configuration to proto mapping", "error", err)
		return nil, err
	}

	s.Logger.Trace("finalized configuration", "proto", confProto)

	return &vagrant_plugin_sdk.Config_FinalizeResponse{
		Data: confProto.(*vagrant_plugin_sdk.Args_ConfigData),
	}, nil
}

func (s *configServer) InitSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "config"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.InitFunc())
}

func (s *configServer) Init(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Config_InitResponse, error) {
	// Now execute the finalize function including the native configuration type
	raw, err := s.CallDynamicFunc(s.Impl.InitFunc(),
		false, req.Args, argmapper.Typed(ctx))
	if err != nil {
		return nil, errors.Wrap("config init failure", err)
	}

	configData, err := s.generateConfigData(ctx, raw)
	if err != nil {
		return nil, errors.Wrap("could not convert to proto", err)
	}

	return &vagrant_plugin_sdk.Config_InitResponse{
		Data: configData,
	}, nil
}

func (s *configServer) generateConfigData(ctx context.Context, in any) (*vagrant_plugin_sdk.Args_ConfigData, error) {
	confData, ok := in.(*component.ConfigData)
	if !ok {
		mapData := structs.Map(in)
		confData = &component.ConfigData{
			Source: structs.Name(in),
			Data:   mapData,
		}
	}

	confProto, err := s.Map(confData,
		(**vagrant_plugin_sdk.Args_ConfigData)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return confProto.(*vagrant_plugin_sdk.Args_ConfigData), nil
}

func (s *configServer) generateConfigSpec(
	ctx context.Context,
	fn interface{},
	log hclog.Logger,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	// Fetch the configuration structure from the plugin as this
	// will be the type that is scrubbed from the function arguments
	sConfig, err := s.CallDynamicFunc(
		s.Impl.StructFunc(),
		false,
		funcspec.Args{},
		argmapper.Typed(ctx),
	)
	if err != nil {
		log.Error("configuration structure retrieval", "error", err)
		return nil, err
	}

	log.Trace("original plugin function", "func", hclog.Fmt("%T", fn))

	// Build a custom function so the native configuration structure can
	// be filtered from the arguments
	ins := []reflect.Type{}
	fnType := reflect.TypeOf(fn)
	checkType := reflect.TypeOf(sConfig)
	for i := 0; i < fnType.NumIn(); i++ {
		if fnType.In(i) != checkType {
			ins = append(ins, fnType.In(i))
		}
	}
	// NOTE: Since this function is used to generate the spec which is used
	//       for determining required input arguments only, the out values
	//       are irrelevant so they are ignored.
	newFnType := reflect.FuncOf(ins, []reflect.Type{}, fnType.IsVariadic())
	newFnValue := reflect.New(newFnType)
	newFn := newFnValue.Elem().Interface()

	log.Trace("modified plugin function", "func", hclog.Fmt("%T", newFn))

	// Now generate the spec with the customized function
	result, err := s.GenerateSpec(newFn)
	if err != nil {
		log.Error("spec generation", "error", err)
		return nil, err
	}

	return result, err
}

var (
	_ plugin.Plugin                          = (*ConfigPlugin)(nil)
	_ plugin.GRPCPlugin                      = (*ConfigPlugin)(nil)
	_ vagrant_plugin_sdk.ConfigServiceServer = (*configServer)(nil)
	_ component.Config                       = (*configClient)(nil)
	_ core.Config                            = (*configClient)(nil)
	_ core.Seeder                            = (*configClient)(nil)
)
