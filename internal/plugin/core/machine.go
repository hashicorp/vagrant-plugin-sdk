package core

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Machine struct {
	client     *MachineClient
	ResourceID string
	ServerAddr string
}

// MachinePlugin is just a GRPC client for a machine
type MachinePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Machine
}

// Implements plugin.GRPCPlugin
func (p *MachinePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &MachineClient{
		client:       pb.NewMachineServiceClient(c),
		ServerTarget: c.Target(),
		Mappers:      p.Mappers,
		Logger:       p.Logger,
		Broker:       broker,
	}, nil
}

func NewMachine(client *MachineClient, resourceID string) *Machine {
	return &Machine{
		client:     client,
		ResourceID: resourceID,
		ServerAddr: client.ServerTarget,
	}
}

// Machine implements core.Machine interface
type MachineClient struct {
	Broker       *plugin.GRPCBroker
	Logger       hclog.Logger
	Mappers      []*argmapper.Func
	ResourceID   string // NOTE(spox): This needs to be added (resource identifier)
	ServerTarget string

	client proto.MachineServiceClient
}

func (m *Machine) Communicate() (comm core.Communicator, err error) {

	// TODO
	return nil, nil
}

func (m *Machine) Guest() (g core.Guest, err error) {
	// TODO
	return nil, nil
}

func (m *Machine) State() (state *core.MachineState, err error) {
	// TODO
	return nil, nil
}

func (m *Machine) IndexUUID() (id string, err error) {
	// TODO
	return "", nil
}

func (m *Machine) Inspect() (printable string, err error) {
	// TODO
	return "", nil
}

func (m *Machine) Reload() (err error) {
	// TODO
	return nil
}

func (m *Machine) ConnectionInfo() (info *core.ConnectionInfo, err error) {
	// TODO
	return nil, nil
}

func (m *Machine) UID() (user_id int, err error) {
	// TODO
	return 10, nil
}

func (m *Machine) GetName() (name string, err error) {
	r, err := m.client.client.GetName(
		context.Background(),
		&pb.Machine_GetNameRequest{ResourceId: m.ResourceID},
	)
	if err != nil {
		return "", err
	}

	return r.Name, nil
}

func (m *Machine) SetName(name string) (err error) {
	_, err := m.client.client.SetName(
		context.Background(),
		&pb.Machine_SetNameRequest{
			ResourceId: m.ResourceID,
			Name:       name,
		},
	)
	return
}

func (m *Machine) GetID() (id string, err error) {
	r, err := m.client.client.GetID(
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

func (m *Machine) SetID(id string) (err error) {
	_, err := m.client.client.SetID(
		context.Background(),
		&pb.Machine_SetIDRequest{
			ResourceId: m.ResourceID,
		},
	)
	return
}

func (m *Machine) Box() (b Box, err error) {
	_, err := m.client.client.Box(
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

func (m *Machine) Datadir() (d datadir.Machine, err error) {
	_, err := m.client.client.Datadir(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	// TODO(spox): this needs to be converted
	// d = r.Datadir
	return
}

func (m *Machine) LocalDataPath() (p path.Path, err error) {
	r, err := m.client.client.LocalDataPath(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	p = path.NewPath(r.Path)
	return
}

func (m *Machine) Provider() (p core.Provider, err error) {
	_, err := m.client.client.Provider(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	// TODO(spox): need to extract and convert provider
	return
}

func (m *Machine) VagrantfileName() (name string, err error) {
	r, err := m.client.client.VagrantfileName(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}

	name = r.Name
	return
}

func (m *Machine) VagrantfilePath() (p path.Path, err error) {
	r, err := m.client.client.VagrantfilePath(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}

	p = path.NewPath(r.Path)
	return
}

func (m *Machine) UpdatedAt() (t *time.Time, err error) {
	_, err := m.client.client.UpdatedAt(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}

	// TODO(spox): need to figure out proto types
	return
}

func (m *Machine) UI() (ui *terminal.UI, err error) {
	_, err := m.client.client.UI(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}

	// TODO(spox): mapper to convert
	return
}

func (m *Machine) SyncedFolders() (folders []core.SyncedFolder, err error) {
	// TODO
	return nil, nil
}

var (
	_ plugin.Plugin     = (*MachinePlugin)(nil)
	_ plugin.GRPCPlugin = (*MachinePlugin)(nil)
	_ core.Machine      = (*Machine)(nil)
)
