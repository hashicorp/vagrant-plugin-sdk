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
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// VagrantConfiglugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Config component type.
type VagrantConfigPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.VagrantConfig // Impl is the concrete implementation
	Mappers []*argmapper.Func       // Mappers
	Logger  hclog.Logger            // Logger
}

func (p *VagrantConfigPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterVagrantConfigServer(s, &vagrantConfigServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *VagrantConfigPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &vagrantConfigClient{
		client: proto.NewVagrantConfigClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// configClient is an implementation of component.Config over gRPC.
type vagrantConfigClient struct {
	client  proto.VagrantConfigClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *vagrantConfigClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *vagrantConfigClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *vagrantConfigClient) ConfigFunc() interface{} {
	//TODO
	return nil
}

// logPlatformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type vagrantConfigServer struct {
	Impl    component.VagrantConfig
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *vagrantConfigServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *vagrantConfigServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *vagrantConfigServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin             = (*VagrantConfigPlugin)(nil)
	_ plugin.GRPCPlugin         = (*VagrantConfigPlugin)(nil)
	_ proto.VagrantConfigServer = (*vagrantConfigServer)(nil)
	_ component.VagrantConfig   = (*vagrantConfigClient)(nil)
)
