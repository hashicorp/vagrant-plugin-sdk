package core

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"google.golang.org/grpc"

// 	"github.com/hashicorp/go-argmapper"
// 	"github.com/hashicorp/go-hclog"
// 	"github.com/hashicorp/go-plugin"

// 	"github.com/hashicorp/vagrant-plugin-sdk/core"
// 	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
// 	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
// 	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
// 	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
// )

// type TargetPlugin struct {
// 	plugin.NetRPCUnsupportedPlugin

// 	Impl   core.Target
// 	Logger hclog.Logger
// }

// // Implements plugin.GRPCPlugin
// func (p *TargetPlugin) GRPCClient(
// 	ctx context.Context,
// 	broker *plugin.GRPCBroker,
// 	c *grpc.ClientConn,
// ) (interface{}, error) {
// 	return &targetClient{
// 		client: vagrant_plugin_sdk.NewTargetServiceClient(c),
// 		Logger: p.Logger,
// 		Broker: broker,
// 	}, nil
// }

// func (p *TargetPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
// 	vagrant_plugin_sdk.RegisterTargetServiceServer(s, &targetServer{
// 		Impl:   p.Impl,
// 		Logger: p.Logger,
// 		Broker: broker,
// 	})
// 	return nil
// }

// // Machine implements core.Machine interface
// type targetClient struct {
// 	Logger hclog.Logger
// 	Broker *plugin.GRPCBroker

// 	client vagrant_plugin_sdk.TargetServiceClient
// }

// func (t *targetClient) Communicate() (comm core.Communicator, err error) {

// 	// TODO
// 	return nil, nil
// }

// func (t *targetClient) Guest() (g core.Guest, err error) {
// 	// TODO
// 	return nil, nil
// }

// func (t *targetClient) State() (state *core.MachineState, err error) {
// 	// TODO
// 	return nil, nil
// }

// func (t *targetClient) IndexUUID() (id string, err error) {
// 	// TODO
// 	return "", nil
// }

// func (t *targetClient) Inspect() (printable string, err error) {
// 	// TODO
// 	return "", nil
// }

// func (t *targetClient) Reload() (err error) {
// 	// TODO
// 	return nil
// }

// func (t *targetClient) ConnectionInfo() (info *core.ConnectionInfo, err error) {
// 	// TODO
// 	return nil, nil
// }

// func (t *targetClient) UID() (user_id int, err error) {
// 	// TODO
// 	return 10, nil
// }

// func (t *targetClient) GetName() (name string, err error) {
// 	r, err := m.c.client.GetName(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_GetNameRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return "", err
// 	}

// 	return r.Name, nil
// }

// func (t *targetClient) SetName(name string) (err error) {
// 	_, err = m.c.client.SetName(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_SetNameRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 			Name: name,
// 		},
// 	)
// 	return
// }

// func (t *targetClient) GetID() (id string, err error) {
// 	r, err := m.c.client.GetID(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_GetIDRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	id = r.Id
// 	return
// }

// func (t *targetClient) SetID(id string) (err error) {
// 	_, err = m.c.client.SetID(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_SetIDRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 			Id: id,
// 		},
// 	)
// 	return
// }

// func (t *targetClient) Box() (b core.Box, err error) {
// 	_, err = m.c.client.Box(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_BoxRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	// TODO(spox): this needs to be converted
// 	//	b = r.Box
// 	return
// }

// func (t *targetClient) Datadir() (d *datadir.Machine, err error) {
// 	_, err = m.c.client.Datadir(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_DatadirRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	// TODO(spox): this needs to be converted
// 	// d = r.Datadir
// 	return
// }

// func (t *targetClient) LocalDataPath() (p path.Path, err error) {
// 	r, err := m.c.client.LocalDataPath(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_LocalDataPathRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	p = path.NewPath(r.Path)
// 	return
// }

// func (t *targetClient) Provider() (p core.Provider, err error) {
// 	_, err = m.c.client.Provider(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_ProviderRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	// TODO(spox): need to extract and convert provider
// 	return
// }

// func (t *targetClient) VagrantfileName() (name string, err error) {
// 	r, err := m.c.client.VagrantfileName(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_VagrantfileNameRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}

// 	name = r.Name
// 	return
// }

// func (t *targetClient) VagrantfilePath() (p path.Path, err error) {
// 	r, err := m.c.client.VagrantfilePath(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_VagrantfilePathRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		return
// 	}

// 	p = path.NewPath(r.Path)
// 	return
// }

// func (t *targetClient) UpdatedAt() (t *time.Time, err error) {
// 	_, err = m.c.client.UpdatedAt(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_UpdatedAtRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}

// 	// TODO(spox): need to figure out proto types
// 	return
// }

// func (t *targetClient) UI() (ui *terminal.UI, err error) {
// 	_, err = m.c.client.UI(
// 		context.Background(),
// 		&vagrant_plugin_sdk.Machine_UIRequest{
// 			Machine: &vagrant_plugin_sdk.Ref_Machine{
// 				ResourceId: m.ResourceID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}

// 	// TODO(spox): mapper to convert
// 	return
// }

// func (t *targetClient) SyncedFolders() (folders []core.SyncedFolder, err error) {
// 	// TODO
// 	return nil, nil
// }

// var (
// 	_ plugin.Plugin     = (*MachinePlugin)(nil)
// 	_ plugin.GRPCPlugin = (*MachinePlugin)(nil)
// 	_ core.Machine      = (*Machine)(nil)
// )
