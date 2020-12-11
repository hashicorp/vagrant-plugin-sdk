package plugin

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
)

// Machine is just a GRCP client for a machine
type Machine struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Machine
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
	client proto.MachineServiceClient
}

func (m *machineClient) GetMachine() (*core.Machine, error) {
	return nil, nil
}

func (m *machineClient) ListMachines() ([]*core.Machine, error) {
	rawMachines, err := m.client.ListMachines(
		context.Background(),
		&proto.ListMachineRequest{})
	if err != nil {
		return nil, err
	}

	var machines []*core.Machine
	mapstructure.Decode(rawMachines, &machines)
	return machines, nil
}

func (m *machineClient) UpsertMachine(*core.Machine) error {
	return nil
}

var (
	_ plugin.Plugin     = (*Machine)(nil)
	_ plugin.GRPCPlugin = (*Machine)(nil)
	// _ core.Machine      = (*machineClient)(nil)
)
