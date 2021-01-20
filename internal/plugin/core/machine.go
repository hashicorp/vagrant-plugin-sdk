package core

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	proto "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Machine struct {
	c          *MachineClient
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

func (p *MachinePlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	return errors.New("Server plugin not provided")
}

func NewMachine(client *MachineClient, resourceID string) *Machine {
	return &Machine{
		c:          client,
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
	r, err := m.c.client.GetName(
		context.Background(),
		&pb.Machine_GetNameRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return "", err
	}

	return r.Name, nil
}

func (m *Machine) SetName(name string) (err error) {
	_, err = m.c.client.SetName(
		context.Background(),
		&pb.Machine_SetNameRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
			Name: name,
		},
	)
	return
}

func (m *Machine) GetID() (id string, err error) {
	r, err := m.c.client.GetID(
		context.Background(),
		&pb.Machine_GetIDRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}
	id = r.Id
	return
}

func (m *Machine) SetID(id string) (err error) {
	_, err = m.c.client.SetID(
		context.Background(),
		&pb.Machine_SetIDRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
			Id: id,
		},
	)
	return
}

func (m *Machine) Box() (b core.Box, err error) {
	_, err = m.c.client.Box(
		context.Background(),
		&pb.Machine_BoxRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}
	// TODO(spox): this needs to be converted
	//	b = r.Box
	return
}

func (m *Machine) Datadir() (d *datadir.Machine, err error) {
	_, err = m.c.client.Datadir(
		context.Background(),
		&pb.Machine_DatadirRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}
	// TODO(spox): this needs to be converted
	// d = r.Datadir
	return
}

func (m *Machine) LocalDataPath() (p path.Path, err error) {
	r, err := m.c.client.LocalDataPath(
		context.Background(),
		&pb.Machine_LocalDataPathRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}
	p = path.NewPath(r.Path)
	return
}

func (m *Machine) Provider() (p core.Provider, err error) {
	_, err = m.c.client.Provider(
		context.Background(),
		&pb.Machine_ProviderRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}
	// TODO(spox): need to extract and convert provider
	return
}

func (m *Machine) VagrantfileName() (name string, err error) {
	r, err := m.c.client.VagrantfileName(
		context.Background(),
		&pb.Machine_VagrantfileNameRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}

	name = r.Name
	return
}

func (m *Machine) VagrantfilePath() (p path.Path, err error) {
	r, err := m.c.client.VagrantfilePath(
		context.Background(),
		&pb.Machine_VagrantfilePathRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)

	if err != nil {
		return
	}

	p = path.NewPath(r.Path)
	return
}

func (m *Machine) UpdatedAt() (t *time.Time, err error) {
	_, err = m.c.client.UpdatedAt(
		context.Background(),
		&pb.Machine_UpdatedAtRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
	if err != nil {
		return
	}

	// TODO(spox): need to figure out proto types
	return
}

func (m *Machine) UI() (ui *terminal.UI, err error) {
	_, err = m.c.client.UI(
		context.Background(),
		&pb.Machine_UIRequest{
			Machine: &pb.Ref_Machine{
				ResourceId: m.ResourceID,
			},
		},
	)
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
