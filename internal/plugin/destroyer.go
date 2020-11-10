package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// destroyerClient implements component.Destroyer for a service that
// has the destroy methods implemented.
type destroyerClient struct {
	Client  destroyerProtoClient
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
}

func (c *destroyerClient) Implements(ctx context.Context) (bool, error) {
	if c == nil {
		return false, nil
	}

	resp, err := c.Client.IsDestroyer(ctx, &empty.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Implements, nil
}

func (c *destroyerClient) DestroyFunc() interface{} {
	impl, err := c.Implements(context.Background())
	if err != nil {
		return funcErr(err)
	}
	if !impl {
		return nil
	}

	// Get the spec
	spec, err := c.Client.DestroySpec(context.Background(), &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}

	return funcspec.Func(spec, c.destroy,
		argmapper.Logger(c.Logger),
		argmapper.Typed(&pluginargs.Internal{
			Broker:  c.Broker,
			Mappers: c.Mappers,
			Cleanup: &pluginargs.Cleanup{},
		}),
	)
}

func (c *destroyerClient) destroy(
	ctx context.Context,
	args funcspec.Args,
	internal *pluginargs.Internal,
) error {
	// Run the cleanup
	defer internal.Cleanup.Close()

	// Call our function
	_, err := c.Client.Destroy(ctx, &pb.FuncSpec_Args{Args: args})
	return err
}

// destroyerServer implements the common Destroyer-related RPC calls.
// This should be embedded into the service implementation.
type destroyerServer struct {
	*base
	Impl interface{}
}

func (s *destroyerServer) IsDestroyer(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.ImplementsResp, error) {
	d, ok := s.Impl.(component.Destroyer)
	return &pb.ImplementsResp{
		Implements: ok && d.DestroyFunc() != nil,
	}, nil
}

func (s *destroyerServer) DestroySpec(
	ctx context.Context,
	args *empty.Empty,
) (*pb.FuncSpec, error) {
	return funcspec.Spec(s.Impl.(component.Destroyer).DestroyFunc(),
		//argmapper.WithNoOutput(), // we only expect an error value so ignore the rest
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Logger(s.Logger),
		argmapper.Typed(s.internal()),
	)
}

func (s *destroyerServer) Destroy(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*empty.Empty, error) {
	internal := s.internal()
	defer internal.Cleanup.Close()

	_, err := callDynamicFunc2(s.Impl.(component.Destroyer).DestroyFunc(), args.Args,
		argmapper.ConverterFunc(s.Mappers...),
		argmapper.Typed(internal),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// destroyerProtoClient is the interface we expect any gRPC service that
// supports destroy to implement.
type destroyerProtoClient interface {
	IsDestroyer(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.ImplementsResp, error)
	DestroySpec(context.Context, *empty.Empty, ...grpc.CallOption) (*pb.FuncSpec, error)
	Destroy(context.Context, *pb.FuncSpec_Args, ...grpc.CallOption) (*empty.Empty, error)
}

var (
	_ component.Destroyer = (*destroyerClient)(nil)
)
