package plugin

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Machine is just a GRCP client for a machine
type Machine struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *Machine) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &machineClient{
		baseClient: &baseClient{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

func (p *Machine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

type machineClient struct {
	*baseClient
	// client proto.MachineServiceClient
}

var (
	_ plugin.Plugin     = (*Machine)(nil)
	_ plugin.GRPCPlugin = (*Machine)(nil)
	// _ core.Machine      = (*machineClient)(nil)
)
