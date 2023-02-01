// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// MapperPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Mapper plugin type.
type MapperPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	*BasePlugin
}

func (p *MapperPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterMapperServer(s, &mapperServer{
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

func (p *MapperPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &MapperClient{
		client:     vagrant_plugin_sdk.NewMapperClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

// MapperClient is an implementation of component.Mapper over gRPC.
type MapperClient struct {
	*BaseClient

	client vagrant_plugin_sdk.MapperClient
}

// Mappers returns the list of mappers that are supported by this plugin.
func (c *MapperClient) Mappers() ([]*argmapper.Func, error) {
	// Get our list of mapper FuncSpecs
	resp, err := c.client.ListMappers(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	c.Logger.Info("list of mappers for plugin", "mappers", resp.Funcs)

	// For each FuncSpec we turn that into a real mapper.Func which calls back
	// into our client to make an RPC call to generate the proper type.
	var funcs []*argmapper.Func
	for _, spec := range resp.Funcs {
		specCopy := spec

		if len(specCopy.Result) < 1 {
			c.Logger.Error("spec result is invalid length", "length", len(specCopy.Result))
			continue
		}
		// We use a closure here to capture spec so that we can provide
		// the correct result type. All we're doing is making our callback
		// call the Map RPC call and return the result/error.
		cb := func(ctx context.Context, args funcspec.Args) (*anypb.Any, error) {
			resp, err := c.client.Map(ctx, &vagrant_plugin_sdk.Map_Request{
				Args:   &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
				Result: specCopy.Result[0].Type,
			})
			if err != nil {
				return nil, err
			}

			return resp.Result, nil
		}

		// Build our funcspec function
		f := funcspec.Func(specCopy, cb)

		// Accumulate our functions
		funcs = append(funcs, f.Func)
	}

	return funcs, nil
}

// mapperServer is a gRPC server that implements the Mapper service.
type mapperServer struct {
	*BaseServer

	vagrant_plugin_sdk.UnimplementedMapperServer
}

func (s *mapperServer) ListMappers(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Map_ListResponse, error) {
	// Go through each mapper and build up our FuncSpecs for each of them.
	var result vagrant_plugin_sdk.Map_ListResponse
	for _, m := range s.Mappers {
		fn := m.Func()

		// Skip our built-in protomappers
		if _, ok := ProtomapperAllMap[reflect.ValueOf(fn).Type()]; ok {
			continue
		}

		spec, err := s.GenerateSpec(fn)
		if err != nil {
			s.Logger.Warn(
				"error converting mapper, will not notify plugin host",
				"func", m.String(),
				"err", err,
			)
			continue
		}

		result.Funcs = append(result.Funcs, spec)
	}

	return &result, nil
}

func (s *mapperServer) Map(
	ctx context.Context,
	args *vagrant_plugin_sdk.Map_Request,
) (*vagrant_plugin_sdk.Map_Response, error) {
	// Find the output type, which we should know about.
	protoType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(args.Result))
	if err != nil {
		return nil, err
	}
	if protoType == nil {
		return nil, status.Newf(
			codes.FailedPrecondition,
			"output type is not known: %s",
			args.Result,
		).Err()
	}

	goType := reflect.TypeOf(protoType.New())
	// Build our function that expects this type as an argument
	// so that we can return it. We do this dynamic function thing so
	// that we can just pretend that this is a function we have so that
	// callDynamicFunc just works.
	f := reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{goType}, []reflect.Type{goType}, false),
		func(args []reflect.Value) []reflect.Value {
			return args
		},
	).Interface()

	// Call it!
	result, err := s.CallDynamicFunc(f, (*proto.Message)(nil), args.Args.Args,
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}
	resultAny, err := dynamic.EncodeAny(result.(protoreflect.ProtoMessage))
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Map_Response{Result: resultAny}, nil
}

var (
	_ plugin.Plugin                   = (*MapperPlugin)(nil)
	_ plugin.GRPCPlugin               = (*MapperPlugin)(nil)
	_ vagrant_plugin_sdk.MapperServer = (*mapperServer)(nil)

	// ProtomapperAllMap is a set of all the protomapper mappers so
	// that we can easily filter them in ListMappers.
	ProtomapperAllMap = map[reflect.Type]struct{}{}
)
