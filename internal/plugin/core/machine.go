package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
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
		client:  pb.NewMachineServiceClient(c),
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	}, nil
}

// Implements plugin.GRPCPlugin
func (p *Machine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

type machineClient struct {
	Broker  *plugin.GRPCBroker
	Logger  hclog.Logger
	Mappers []*argmapper.Func
	client  pb.MachineServiceClient
}

// Implements component.Machine
func (m *machineClient) GetServerAddr() string {
	// TODO
	return "nothing!"
}

// Implements component.Machine
func (m *machineClient) GetMachine(id string) (core.Machine, error) {
	rawMachine, err := m.client.GetMachine(
		context.Background(),
		&pb.GetMachineRequest{Ref: &pb.Ref_Machine{Id: id}},
	)
	if err != nil {
		return nil, err
	}

	// TODO: I think this needs to have a GetMachineFunc with the mappers and
	//   everything. Then you can maybe decode the response into a more useful
	//   machine implementation?
	var machine *TestMachine
	mapstructure.Decode(rawMachine.Machine, &machine)
	return machine, nil
}

type ahherror struct {
}

func (m *ahherror) Error() string { return "ahh" }

// Machine implements sdkCore.Machine interface
type TestMachine struct {
	Datadir       string
	Id            string
	LocalDataPath string
	Name          string
}

func (m *TestMachine) UID() (user_id int, err error) {
	return 0, nil
}

// Implements component.Machine
func (m *machineClient) ListMachines() ([]core.Machine, error) {
	rawMachines, err := m.client.ListMachines(
		context.Background(),
		&pb.ListMachineRequest{})
	if err != nil {
		return nil, err
	}

	var machines []core.Machine
	mapstructure.Decode(rawMachines, &machines)
	return machines, nil
}

// Implements component.Machine
func (m *machineClient) UpsertMachine(machine core.Machine) error {
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

var (
	_ plugin.Plugin     = (*Machine)(nil)
	_ plugin.GRPCPlugin = (*Machine)(nil)
	_ component.Machine = (*machineClient)(nil)
)
