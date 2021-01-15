package core

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
)

// MachinePlugin is just a GRPC client for a machine
type MachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Machine
}

func NewMachine(client *MachineClient, m core.Machine) *Machine {
	var machine *Machine
	mapstructure.Decode(m, &machine)
	machine.client = client
	return machine
}

// Machine implements core.Machine interface
type machineClient struct {
	Broker     *plugin.GRPCBroker
	Logger     hclog.Logger
	Mappers    []*argmapper.Func
	ResourceID string // NOTE(spox): This needs to be added (resource identifier)

	client proto.MachineServiceClient
}

func (m *machineClient) Communicate() (comm core.Communicator, err error) {
	// TODO
	return nil, nil
}

func (m *machineClient) Guest() (g core.Guest, err error) {
	// TODO
	return nil, nil
}

func (m *machineClient) State() (state *core.MachineState, err error) {
	// TODO
	return nil, nil
}

func (m *machineClient) IndexUUID() (id string, err error) {
	// TODO
	return "", nil
}

func (m *machineClient) Inspect() (printable string, err error) {
	// TODO
	return "", nil
}

func (m *machineClient) Reload() (err error) {
	// TODO
	return nil
}

func (m *machineClient) ConnectionInfo() (info *core.ConnectionInfo, err error) {
	// TODO
	return nil, nil
}

func (m *machineClient) UID() (user_id int, err error) {
	// TODO
	return 10, nil
}

func (m *machineClient) GetName() (name string, err error) {
	r, err := m.client.GetName(context.Background(), &pb.Machine_GetNameRequest{ResourceId: m.ResourceID})
	if err != nil {
		return "", err
	}

	return r.Name, nil
}

func (m *machineClient) SetName(name string) (err error) {
	_, err := m.client.SetName(
		context.Background(),
		&pb.Machine_SetNameRequest{
			ResourceId: m.ResourceID,
			Name:       name,
		},
	)
	return
}

func (m *machineClient) GetID() (id string, err error) {
	r, err := m.client.GetID(
		context.Background(),
		&pb.Machine_GetIDRequest{
			ResourceId: m.ResourceID,
		},
	)
	if err != nil {
		return
	}
	id = r.Id
	return
}

func (m *machineClient) SetID(id string) (err error) {
	_, err := m.client.SetID(
		context.Background(),
		&pb.Machine_SetIDRequest{
			ResourceId: m.ResourceID,
		},
	)
	return
}

func (m *machineClient) Box() (b Box, err error) {
	_, err := m.client.Box(
		context.Background(),
		&empty.Empty{},
	)
	if err != nil {
		return
	}
	// TODO(spox): this needs to be converted
	//	b = r.Box
	return
}

func (m *machineClient) SyncedFolders() (folders []core.SyncedFolder, err error) {
	// TODO
	return nil, nil
}

// Implements plugin.GRPCPlugin
func (p *MachinePlugin) GRPCClient(
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
func (p *MachinePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// Not implemented. The machine plugin server is in vagrant core
	return nil
}

type MachineClient struct {
	Broker  *plugin.GRPCBroker
	Logger  hclog.Logger
	Mappers []*argmapper.Func
	client  pb.MachineServiceClient
}

// Implements component.Machine
func (m *MachineClient) GetServerAddr() string {
	// TODO: I don't think this is needed on the client side
	return "nothing!"
}

// Implements component.Machine
func (m *MachineClient) GetMachine(id string) (core.Machine, error) {
	rawMachine, err := m.client.GetMachine(
		context.Background(),
		&pb.GetMachineRequest{Ref: &pb.Ref_Machine{Id: id}},
	)
	if err != nil {
		return nil, err
	}
	var machine *Machine
	mapstructure.Decode(rawMachine.Machine, &machine)
	return machine, nil
}

// Implements component.Machine
func (m *MachineClient) ListMachines() ([]core.Machine, error) {
	rawMachines, err := m.client.ListMachines(
		context.Background(),
		&pb.ListMachineRequest{})
	if err != nil {
		return nil, err
	}

	// TODO: test
	var machines []core.Machine
	mapstructure.Decode(rawMachines, &machines)
	return machines, nil
}

// Implements component.Machine
func (m *MachineClient) UpsertMachine(mach core.Machine) (core.Machine, error) {
	var machinepb *pb.Machine
	mapstructure.Decode(mach.(*Machine), &machinepb)

	resp, err := m.client.UpsertMachine(
		context.Background(),
		&pb.UpsertMachineRequest{Machine: machinepb},
	)
	if err != nil {
		return nil, err
	}
	var resultMachine *Machine
	mapstructure.Decode(resp.Machine, &resultMachine)
	return resultMachine, nil
}

var (
	_ plugin.Plugin     = (*MachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*MachinePlugin)(nil)
	_ component.Machine = (*MachineClient)(nil)
	_ core.Machine      = (*Machine)(nil)
)
