package plugin

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
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

// Implements plugin.GRPCPlugin
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

// Implements plugin.GRPCPlugin
func (p *Machine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

type machineClient struct {
	*baseClient
	client pb.MachineServiceClient
}

func (m *machineClient) GetMachine(id string) (*core.Machine, error) {
	rawMachine, err := m.client.GetMachine(
		context.Background(),
		&pb.GetMachineRequest{Ref: &pb.Ref_Machine{Id: id}})
	if err != nil {
		return nil, err
	}

	var machine *core.Machine
	mapstructure.Decode(rawMachine, &machine)
	return machine, nil
}

func (m *machineClient) ListMachines() ([]*core.Machine, error) {
	rawMachines, err := m.client.ListMachines(
		context.Background(),
		&pb.ListMachineRequest{})
	if err != nil {
		return nil, err
	}

	var machines []*core.Machine
	mapstructure.Decode(rawMachines, &machines)
	return machines, nil
}

func (m *machineClient) UpsertMachine(machine *core.Machine) error {
	var machinepb *pb.Machine
	mapstructure.Decode(machine, &machinepb)
	_, err := m.client.UpsertMachine(
		context.Background(),
		&pb.UpsertMachineRequest{Machine: machinepb})
	if err != nil {
		return err
	}

	return nil
}

// TODO: (sophia) delete machine

var (
	_ plugin.Plugin     = (*Machine)(nil)
	_ plugin.GRPCPlugin = (*Machine)(nil)
	// _ core.Machine      = (*machineClient)(nil)
)
