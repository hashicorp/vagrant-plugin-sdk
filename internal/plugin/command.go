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
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// CommandPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Command component type.
type CommandPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Command // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *CommandPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCommandServiceServer(s, &commandServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *CommandPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &commandClient{
		client: proto.NewCommandServiceClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// commandClient is an implementation of component.Command over gRPC.
type commandClient struct {
	client  proto.CommandServiceClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *commandClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *commandClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *commandClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *commandClient) CommandFunc() interface{} {
	//TODO
	return nil
}

// commandServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type commandServer struct {
	Impl    component.Command
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *commandServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *commandServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *commandServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin              = (*CommandPlugin)(nil)
	_ plugin.GRPCPlugin          = (*CommandPlugin)(nil)
	_ proto.CommandServiceServer = (*commandServer)(nil)
	_ component.Command          = (*commandClient)(nil)
)
