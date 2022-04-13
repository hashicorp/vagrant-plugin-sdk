package core

import (
	"context"
	"errors"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type BoxCollectionPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.BoxCollection
	*vplugin.BasePlugin
}

func (p *BoxCollectionPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &boxCollectionClient{
		client:     vagrant_plugin_sdk.NewBoxCollectionServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *BoxCollectionPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBoxCollectionServiceServer(s, &boxCollectionServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type boxCollectionClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.BoxCollectionServiceClient
}

func (b *boxCollectionClient) Add(path, name, version, metadataURL string, force bool, providers ...string) (box core.Box, err error) {
	r, err := b.client.Add(b.Ctx, &vagrant_plugin_sdk.BoxCollection_AddRequest{
		Path: path, Name: name, Version: version, MetadataUrl: metadataURL, Force: force, Providers: providers,
	})
	if err != nil {
		b.Logger.Error("failed to add box",
			"error", err,
		)
		return
	}
	result, err := b.Map(
		r, (*core.Box)(nil), argmapper.Typed(b.Ctx),
	)
	if err != nil {
		b.Logger.Error("failed to add box",
			"error", err,
		)

		return nil, err
	}
	return result.(core.Box), nil
}

func (b *boxCollectionClient) All() (boxes []core.Box, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box list",
				"error", err,
			)
		}
	}()

	r, err := b.client.All(b.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	boxes = []core.Box{}
	for _, box := range r.Boxes {
		var mappedBox interface{}
		mappedBox, err = b.Map(
			box, (*core.Box)(nil), argmapper.Typed(b.Ctx),
		)
		if err != nil {
			return nil, err
		}
		boxes = append(boxes, mappedBox.(core.Box))
	}
	return
}

func (b *boxCollectionClient) Clean(name string) (err error) {
	_, err = b.client.Clean(b.Ctx, &vagrant_plugin_sdk.BoxCollection_CleanRequest{
		Name: name,
	})
	if err != nil {
		b.Logger.Error("failed to clean box",
			"error", err,
		)
	}
	return
}

func (b *boxCollectionClient) Find(name string, version string, providers ...string) (box core.Box, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to find box",
				"error", err,
			)
		}
	}()

	r, err := b.client.Find(b.Ctx, &vagrant_plugin_sdk.BoxCollection_FindRequest{
		Name: name, Version: version, Providers: providers,
	})
	if err != nil {
		return
	}
	result, err := b.Map(
		r, (*core.Box)(nil), argmapper.Typed(b.Ctx),
	)
	if err != nil {
		return nil, err
	}
	return result.(core.Box), nil
}

type boxCollectionServer struct {
	*vplugin.BaseServer

	Impl core.BoxCollection
	vagrant_plugin_sdk.UnimplementedBoxCollectionServiceServer
}

func (b *boxCollectionServer) Add(
	ctx context.Context, in *vagrant_plugin_sdk.BoxCollection_AddRequest,
) (r *vagrant_plugin_sdk.Args_Box, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to add box",
				"error", err,
			)
		}
	}()

	box, err := b.Impl.Add(
		in.Path, in.Name, in.Version, in.MetadataUrl, in.Force, in.Providers...,
	)
	if err != nil {
		return
	}
	boxProto, err := b.Map(
		box, (**vagrant_plugin_sdk.Args_Box)(nil), argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}
	return boxProto.(*vagrant_plugin_sdk.Args_Box), nil
}

func (b *boxCollectionServer) All(
	ctx context.Context, in *emptypb.Empty,
) (r *vagrant_plugin_sdk.BoxCollection_AllResponse, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get box list",
				"error", err,
			)
		}
	}()

	boxes, err := b.Impl.All()
	if err != nil {
		return
	}
	boxesProto := []*vagrant_plugin_sdk.Args_Box{}
	for _, box := range boxes {
		var boxProto interface{}
		boxProto, err = b.Map(
			box, (**vagrant_plugin_sdk.Args_Box)(nil), argmapper.Typed(ctx),
		)
		if err != nil {
			return nil, err
		}
		boxesProto = append(boxesProto, boxProto.(*vagrant_plugin_sdk.Args_Box))
	}
	return &vagrant_plugin_sdk.BoxCollection_AllResponse{Boxes: boxesProto}, nil
}

func (b *boxCollectionServer) Clean(
	ctx context.Context, in *vagrant_plugin_sdk.BoxCollection_CleanRequest,
) (r *emptypb.Empty, err error) {
	err = b.Impl.Clean(in.Name)
	if err != nil {
		b.Logger.Error("failed to clean box",
			"error", err,
		)

		return
	}
	return &emptypb.Empty{}, nil
}

func (b *boxCollectionServer) Find(
	ctx context.Context, in *vagrant_plugin_sdk.BoxCollection_FindRequest,
) (r *vagrant_plugin_sdk.Args_Box, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to find box",
				"error", err,
			)
		}
	}()

	box, err := b.Impl.Find(
		in.Name, in.Version, in.Providers...,
	)
	if err != nil {
		return
	}
	if box == nil {
		return nil, errors.New("no box found")
	}
	boxProto, err := b.Map(
		box, (**vagrant_plugin_sdk.Args_Box)(nil), argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}
	return boxProto.(*vagrant_plugin_sdk.Args_Box), nil
}

var (
	_ plugin.Plugin      = (*BoxCollectionPlugin)(nil)
	_ plugin.GRPCPlugin  = (*BoxCollectionPlugin)(nil)
	_ core.BoxCollection = (*boxCollectionClient)(nil)
)
