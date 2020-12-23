package core

import (
	"context"
	"time"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/multistep"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
)

// Machine is just a GRCP client for a machine
type MachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Machine
}

// Machine implements sdkCore.Machine interface
type Machine struct {
	Box             *core.Box
	Datadir         string
	Environment     *core.Environment
	Id              string
	LocalDataPath   string
	Name            string
	Provider        *core.Provider
	VagrantfileName string
	VagrantfilePath string
	UpdatedAt       *time.Time
	UI              *terminal.UI
}

func (m *Machine) Action(name string, state multistep.StateBag) (err error) {
	return nil
}

func (m *Machine) Communicate() (comm core.Communicator, err error) {
	return nil, nil
}

func (m *Machine) Guest() (g core.Guest, err error) {
	return nil, nil
}

func (m *Machine) SetID(value string) (err error) {
	return nil
}

func (m *Machine) State() (state *core.MachineState, err error) {
	return nil, nil
}

func (m *Machine) IndexUUID() (id string, err error) {
	return "", nil
}

func (m *Machine) Inspect() (printable string, err error) {
	return "", nil
}

func (m *Machine) Reload() (err error) {
	return nil
}

func (m *Machine) ConnectionInfo() (info core.ConnectionInfo, err error) {
	return nil, nil
}

func (m *Machine) UID() (user_id int, err error) {
	return 10, nil
}
func (m *Machine) SyncedFolders() (folders []core.SyncedFolder, err error) {
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
	var machine *Machine
	mapstructure.Decode(rawMachine.Machine, &machine)
	return machine, nil
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
	_ plugin.Plugin     = (*MachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*MachinePlugin)(nil)
	_ component.Machine = (*machineClient)(nil)
	_ core.Machine      = (*Machine)(nil)
)
