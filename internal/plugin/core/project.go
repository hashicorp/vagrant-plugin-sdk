package core

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// ProjectPlugin is just a GRPC client for a project
type ProjectPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.Project
	*vplugin.BasePlugin
}

// Implements plugin.GRPCPlugin
func (p *ProjectPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &projectClient{
		client:     vagrant_plugin_sdk.NewProjectServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *ProjectPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	vagrant_plugin_sdk.RegisterProjectServiceServer(s, &projectServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type projectClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.ProjectServiceClient
}

type projectServer struct {
	*vplugin.BaseServer

	Impl core.Project
	vagrant_plugin_sdk.UnimplementedProjectServiceServer
}

func (p *projectClient) ActiveTargets() (targets []core.Target, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get active targets",
				"error", err,
			)
		}
	}()
	resp, err := p.client.ActiveTargets(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	targets = []core.Target{}
	for _, t := range resp.Targets {
		coreTarget, err := p.Map(t, (*core.Target)(nil),
			argmapper.Typed(p.Ctx),
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, coreTarget.(core.Target))
	}

	return
}

func (p *projectClient) Boxes() (b core.BoxCollection, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get boxes",
				"error", err,
			)
		}
	}()
	r, err := p.client.Boxes(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.BoxCollection)(nil),
		argmapper.Typed(p.Ctx),
	)
	if err == nil {
		b = result.(core.BoxCollection)
	}

	return
}

func (p *projectClient) CWD() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get cwd",
				"error", err,
			)
		}
	}()
	r, err := p.client.CWD(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *projectClient) Config() (v *vagrant_plugin_sdk.Vagrantfile_Vagrantfile, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get config",
				"error", err,
			)
		}
	}()
	r, err := p.client.Config(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return r.Vagrantfile, nil
}

func (p *projectClient) DataDir() (dir *datadir.Project, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get datadir",
				"error", err,
			)
		}
	}()
	r, err := p.client.DataDir(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (**datadir.Project)(nil))
	if err == nil {
		dir = result.(*datadir.Project)
	}

	return
}

func (p *projectClient) DefaultPrivateKey() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get default private key",
				"error", err,
			)
		}
	}()
	r, err := p.client.DefaultPrivateKey(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Key
	}

	return
}

func (p *projectClient) DefaultProvider() (name string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get default provider",
				"error", err,
			)
		}
	}()
	d, err := p.client.DefaultProvider(p.Ctx, &emptypb.Empty{})
	if err == nil {
		name = d.ProviderName
	}

	return
}

func (p *projectClient) Home() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get home path",
				"error", err,
			)
		}
	}()
	r, err := p.client.Home(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *projectClient) Host() (h core.Host, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get host",
				"error", err,
			)
		}
	}()
	r, err := p.client.Host(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.Host)(nil),
		argmapper.Typed(p.Ctx),
	)
	if err == nil {
		h = result.(core.Host)
	}

	return
}

func (p *projectClient) LocalData() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get local data path",
				"error", err,
			)
		}
	}()
	r, err := p.client.LocalData(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *projectClient) PrimaryTargetName() (t string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get primary target name",
				"error", err,
			)
		}
	}()
	resp, err := p.client.PrimaryTargetName(p.Ctx, &emptypb.Empty{})
	if err == nil {
		t = resp.Name
	}

	return
}

func (p *projectClient) Push(name string) (err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to push",
				"error", err,
			)
		}
	}()
	_, err = p.client.Push(p.Ctx,
		&vagrant_plugin_sdk.Project_PushRequest{Name: name},
	)
	return
}

func (p *projectClient) ResourceId() (rid string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get resource id",
				"error", err,
			)
		}
	}()
	r, err := p.client.ResourceId(p.Ctx, &emptypb.Empty{})
	if err == nil {
		rid = r.ResourceId
	}

	return
}

func (p *projectClient) RootPath() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get root path",
				"error", err,
			)
		}
	}()
	r, err := p.client.RootPath(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *projectClient) Target(name string) (t core.Target, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target",
				name,
				"error", err,
			)
		}
	}()
	r, err := p.client.Target(p.Ctx, &vagrant_plugin_sdk.Project_TargetRequest{
		Name: name,
	})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.Target)(nil),
		argmapper.Typed(p.Ctx))
	if err == nil {
		t = result.(core.Target)
	}
	return
}

func (p *projectClient) TargetIds() (ids []string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target ids",
				"error", err,
			)
		}
	}()
	r, err := p.client.TargetIds(p.Ctx, &emptypb.Empty{})
	if err == nil {
		ids = r.Ids
	}

	return
}

func (p *projectClient) TargetIndex() (index core.TargetIndex, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target index",
				"error", err,
			)
		}
	}()
	r, err := p.client.TargetIndex(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.TargetIndex)(nil),
		argmapper.Typed(p.Ctx))
	if err == nil {
		index = result.(core.TargetIndex)
	}
	return
}

func (p *projectClient) TargetNames() (names []string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target names",
				"error", err,
			)
		}
	}()
	r, err := p.client.TargetNames(p.Ctx, &emptypb.Empty{})
	if err == nil {
		names = r.Names
	}

	return
}

func (p *projectClient) Tmp() (path string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get temp path",
				"error", err,
			)
		}
	}()
	r, err := p.client.Tmp(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *projectClient) UI() (ui terminal.UI, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get ui",
				"error", err,
			)
		}
	}()
	r, err := p.client.UI(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*terminal.UI)(nil),
		argmapper.Typed(p.Ctx))
	if err == nil {
		ui = result.(terminal.UI)
	}

	return
}

func (p *projectClient) VagrantfileName() (name string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get Vagrantfile name",
				"error", err,
			)
		}
	}()
	r, err := p.client.VagrantfileName(p.Ctx, &emptypb.Empty{})
	if err == nil {
		name = r.Name
	}

	return
}

func (p *projectClient) VagrantfilePath() (pp path.Path, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get Vagrantfile path",
				"error", err,
			)
		}
	}()
	r, err := p.client.VagrantfilePath(p.Ctx, &emptypb.Empty{})
	if err == nil {
		pp = path.NewPath(r.Path)
	}
	return
}

// Project server

func (p *projectServer) ActiveTargets(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Project_ActiveTargetsResponse, err error) {
	targets, err := p.Impl.ActiveTargets()
	if err != nil {
		return
	}

	targetProtos := []*vagrant_plugin_sdk.Args_Target{}
	for _, t := range targets {
		tp, err := p.Map(t, (**vagrant_plugin_sdk.Args_Target)(nil))
		if err != nil {
			return nil, err
		}
		targetProtos = append(targetProtos, tp.(*vagrant_plugin_sdk.Args_Target))
	}
	r = &vagrant_plugin_sdk.Project_ActiveTargetsResponse{
		Targets: targetProtos,
	}
	return
}

func (p *projectServer) Boxes(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_BoxCollection, err error) {
	boxCollection, err := p.Impl.Boxes()
	if err != nil {
		return
	}

	result, err := p.Map(boxCollection, (**vagrant_plugin_sdk.Args_BoxCollection)(nil))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_BoxCollection)
	}

	return
}

func (p *projectServer) Config(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_ConfigResponse, error) {
	v, err := p.Impl.Config()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_ConfigResponse{
		Vagrantfile: v,
	}, nil
}

func (p *projectServer) CWD(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_CwdResponse, error) {
	c, err := p.Impl.CWD()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_CwdResponse{
		Path: c,
	}, nil
}

func (p *projectServer) DataDir(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_DataDir_Project, err error) {
	d, err := p.Impl.DataDir()
	if err != nil {
		return
	}
	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_DataDir_Project)(nil))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_DataDir_Project)
	}

	return
}

func (p *projectServer) DefaultPrivateKey(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_DefaultPrivateKeyResponse, error) {
	key, err := p.Impl.DefaultPrivateKey()
	p.Logger.Warn("private key on project server", "key", key)
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_DefaultPrivateKeyResponse{
		Key: key,
	}, nil
}

func (p *projectServer) DefaultProvider(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_DefaultProviderResponse, error) {
	provider, err := p.Impl.DefaultProvider()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_DefaultProviderResponse{
		ProviderName: provider,
	}, nil
}

func (p *projectServer) Home(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_HomeResponse, error) {
	path, err := p.Impl.Home()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_HomeResponse{
		Path: path,
	}, nil
}

func (p *projectServer) Host(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_Host, err error) {
	d, err := p.Impl.Host()
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_Host)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Host)
	}

	return
}

func (p *projectServer) LocalData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_LocalDataResponse, error) {
	path, err := p.Impl.LocalData()
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Project_LocalDataResponse{
		Path: path,
	}, nil
}

func (p *projectServer) PrimaryTargetName(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_PrimaryTargetNameResponse, error) {
	name, err := p.Impl.PrimaryTargetName()
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Project_PrimaryTargetNameResponse{
		Name: name,
	}, nil
}

func (p *projectServer) Push(
	ctx context.Context,
	in *vagrant_plugin_sdk.Project_PushRequest,
) (*emptypb.Empty, error) {
	err := p.Impl.Push(in.Name)
	return &emptypb.Empty{}, err
}

func (p *projectServer) ResourceId(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_ResourceIdResponse, error) {
	rid, err := p.Impl.ResourceId()

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_ResourceIdResponse{
		ResourceId: rid,
	}, nil
}

func (p *projectServer) RootPath(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_RootPathResponse, error) {
	path, err := p.Impl.RootPath()

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_RootPathResponse{
		Path: path,
	}, nil
}

func (p *projectServer) Target(
	ctx context.Context,
	in *vagrant_plugin_sdk.Project_TargetRequest,
) (r *vagrant_plugin_sdk.Args_Target, err error) {
	d, err := p.Impl.Target(in.Name)
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_Target)(nil))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Target)
	}

	return
}

func (p *projectServer) TargetIds(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_TargetIdsResponse, error) {
	ids, err := p.Impl.TargetIds()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_TargetIdsResponse{
		Ids: ids}, nil
}

func (p *projectServer) TargetIndex(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TargetIndex, err error) {
	idx, err := p.Impl.TargetIndex()
	if err != nil {
		return nil, err
	}

	result, err := p.Map(idx, (**vagrant_plugin_sdk.Args_TargetIndex)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TargetIndex)
	}
	return
}

func (p *projectServer) TargetNames(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_TargetNamesResponse, error) {
	n, err := p.Impl.TargetNames()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_TargetNamesResponse{
		Names: n}, nil
}

func (p *projectServer) Tmp(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_TmpResponse, error) {
	path, err := p.Impl.Tmp()
	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Project_TmpResponse{
		Path: path,
	}, nil
}

func (p *projectServer) UI(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TerminalUI, err error) {
	d, err := p.Impl.UI()
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_TerminalUI)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TerminalUI)
	}

	return
}

func (p *projectServer) VagrantfileName(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_VagrantfileNameResponse, error) {
	name, err := p.Impl.VagrantfileName()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_VagrantfileNameResponse{
		Name: name,
	}, nil
}

func (p *projectServer) VagrantfilePath(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Project_VagrantfilePathResponse, error) {
	path, err := p.Impl.VagrantfilePath()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Project_VagrantfilePathResponse{
		Path: path.String(),
	}, nil
}

var (
	_ plugin.Plugin     = (*ProjectPlugin)(nil)
	_ plugin.GRPCPlugin = (*ProjectPlugin)(nil)
	_ core.Project      = (*projectClient)(nil)
)
