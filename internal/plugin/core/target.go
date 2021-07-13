package core

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type TargetPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    core.Target
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

// Implements plugin.GRPCPlugin
func (p *TargetPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &targetClient{
		client: vagrant_plugin_sdk.NewTargetServiceClient(c),
		ctx:    ctx,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	}, nil
}

func (p *TargetPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterTargetServiceServer(s, &targetServer{
		Impl: p.Impl,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	})
	return nil
}

// Target implements core.Target interface
type targetClient struct {
	*base

	ctx    context.Context
	client vagrant_plugin_sdk.TargetServiceClient
}

type targetServer struct {
	*base

	Impl core.Target
	vagrant_plugin_sdk.UnsafeTargetServiceServer
	//vagrant_plugin_sdk.UnimplementedTargetServiceServer
}

func (t *targetClient) Communicate() (comm core.Communicator, err error) {

	// TODO:
	return nil, errNotImplemented
}

func (t *targetClient) SetName(name string) (err error) {
	_, err = t.client.SetName(t.ctx, &vagrant_plugin_sdk.Target_SetNameRequest{
		Name: name})
	return
}

func (t *targetClient) Provider() (p core.Provider, err error) {
	return nil, errNotImplemented
}

func (t *targetClient) VagrantfileName() (name string, err error) {
	r, err := t.client.VagrantfileName(t.ctx, &empty.Empty{})
	if err == nil {
		name = r.Name
	}

	return
}

func (t *targetClient) VagrantfilePath() (p path.Path, err error) {
	r, err := t.client.VagrantfilePath(t.ctx, &empty.Empty{})
	if err == nil {
		p = path.NewPath(r.Path)
	}

	return
}

func (t *targetClient) UpdatedAt() (utime *time.Time, err error) {
	r, err := t.client.UpdatedAt(t.ctx, &empty.Empty{})
	if err == nil {
		ut := r.UpdatedAt.AsTime()
		utime = &ut
	}

	return
}

func (c *targetClient) Name() (name string, err error) {
	r, err := c.client.Name(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	name = r.Name
	return
}

func (c *targetClient) ResourceId() (rid string, err error) {
	r, err := c.client.ResourceId(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	rid = r.ResourceId

	return
}

func (c *targetClient) Project() (project core.Project, err error) {
	r, err := c.client.Project(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	result, err := c.Map(r, (*core.Project)(nil))
	if err != nil {
		return
	}
	project = result.(core.Project)
	return
}

func (c *targetClient) Metadata() (mdata map[string]string, err error) {
	r, err := c.client.Metadata(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	result, err := c.Map(r, (*map[string]string)(nil))
	if err != nil {
		return
	}
	mdata = result.(map[string]string)
	return
}

func (c *targetClient) DataDir() (dir *datadir.Target, err error) {
	r, err := c.client.DataDir(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	result, err := c.Map(r, (*datadir.Target)(nil))
	if err != nil {
		return
	}
	dir = result.(*datadir.Target)
	return
}

func (c *targetClient) State() (state core.State, err error) {
	r, err := c.client.State(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	result, err := c.Map(r, (*core.State)(nil))
	if err != nil {
		return
	}
	state = result.(core.State)
	return
}

func (c *targetClient) Record() (record *anypb.Any, err error) {
	r, err := c.client.Record(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	record = r.Record
	return
}

func (c *targetClient) Specialize(kind interface{}) (specialized core.Machine, err error) {
	a, err := anypb.New(&empty.Empty{})
	if err != nil {
		return
	}
	r, err := c.client.Specialize(c.ctx, a)

	if err != nil {
		return
	}

	m := &vagrant_plugin_sdk.Args_Target_Machine{}
	if err = r.UnmarshalTo(m); err != nil {
		return
	}

	s, err := c.Map(m, (*core.Machine)(nil),
		argmapper.Typed(c.ctx))
	return s.(core.Machine), err
}

func (c *targetClient) UI() (ui terminal.UI, err error) {
	r, err := c.client.UI(c.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := c.Map(r, (*terminal.UI)(nil),
		argmapper.Typed(c.ctx))
	if err == nil {
		ui = result.(terminal.UI)
	}

	return
}

// Target Server

func (s *targetServer) Name(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Target_NameResponse, err error) {
	n, err := s.Impl.Name()
	if err == nil {
		r = &vagrant_plugin_sdk.Target_NameResponse{Name: n}
	}

	return
}

func (t *targetServer) SetName(
	ctx context.Context,
	in *vagrant_plugin_sdk.Target_SetNameRequest,
) (*empty.Empty, error) {
	if err := t.Impl.SetName(in.Name); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (t *targetServer) Provider(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Provider, err error) {
	p, err := t.Impl.Provider()
	if err != nil {
		return
	}

	result, err := t.Map(p, (**vagrant_plugin_sdk.Args_Provider)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Provider)
	}

	return
}

func (t *targetServer) VagrantfileName(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_VagrantfileNameResponse, error) {
	n, err := t.Impl.VagrantfileName()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_VagrantfileNameResponse{
		Name: n}, nil
}

func (t *targetServer) VagrantfilePath(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_VagrantfilePathResponse, error) {
	n, err := t.Impl.VagrantfilePath()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_VagrantfilePathResponse{
		Path: n.String()}, nil
}

func (t *targetServer) UpdatedAt(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Target_UpdatedAtResponse, error) {
	u, err := t.Impl.UpdatedAt()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Target_UpdatedAtResponse{
		UpdatedAt: timestamppb.New(*u)}, nil
}

func (s *targetServer) ResourceId(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Target_ResourceIdResponse, err error) {
	rid, err := s.Impl.ResourceId()
	if err == nil {
		r = &vagrant_plugin_sdk.Target_ResourceIdResponse{ResourceId: rid}
	}

	return
}

func (s *targetServer) Project(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Args_Project, error) {
	p, err := s.Impl.Project()
	if err != nil {
		return nil, err
	}

	result, err := s.Map(p, (**vagrant_plugin_sdk.Args_Project)(nil))
	if err != nil {
		return nil, err
	}

	return result.(*vagrant_plugin_sdk.Args_Project), nil
}

func (s *targetServer) Metadata(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_MetadataSet, err error) {
	m, err := s.Impl.Metadata()
	if err != nil {
		return
	}
	result, err := s.Map(m, (**vagrant_plugin_sdk.Args_MetadataSet)(nil))
	if err != nil {
		return
	}
	r = result.(*vagrant_plugin_sdk.Args_MetadataSet)

	return
}

func (s *targetServer) DataDir(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_DataDir_Target, err error) {
	d, err := s.Impl.DataDir()
	if err != nil {
		return
	}
	result, err := s.Map(d, (**vagrant_plugin_sdk.Args_DataDir_Target)(nil))
	if err != nil {
		return
	}
	r = result.(*vagrant_plugin_sdk.Args_DataDir_Target)

	return
}

func (s *targetServer) State(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_Target_State, err error) {
	st, err := s.Impl.State()
	if err != nil {
		return
	}
	result, err := s.Map(st, (**vagrant_plugin_sdk.Args_Target_State)(nil))
	if err != nil {
		return
	}
	r = result.(*vagrant_plugin_sdk.Args_Target_State)

	return
}

func (s *targetServer) Record(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Target_RecordResponse, err error) {
	record, err := s.Impl.Record()
	if err == nil {
		r = &vagrant_plugin_sdk.Target_RecordResponse{Record: record}
	}

	return
}

func (t *targetServer) UI(
	ctx context.Context,
	_ *empty.Empty,
) (r *vagrant_plugin_sdk.Args_TerminalUI, err error) {
	d, err := t.Impl.UI()
	if err != nil {
		return
	}

	result, err := t.Map(d, (**vagrant_plugin_sdk.Args_TerminalUI)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TerminalUI)
	}

	return
}

func (t *targetServer) Specialize(
	ctx context.Context,
	in *any.Any,
) (r *any.Any, err error) {
	mc, ok := t.Impl.(interface{ Machine() core.Machine })
	if !ok {
		return nil, errors.New("could not specialize to machine")
	}

	result, err := t.Map(mc.Machine(), (**vagrant_plugin_sdk.Args_Target_Machine)(nil),
		argmapper.Typed(ctx))
	if err != nil {
		return
	}
	return anypb.New(result.(*vagrant_plugin_sdk.Args_Target_Machine))
}

var (
	_ plugin.Plugin                          = (*TargetPlugin)(nil)
	_ plugin.GRPCPlugin                      = (*TargetPlugin)(nil)
	_ core.Target                            = (*targetClient)(nil)
	_ vagrant_plugin_sdk.TargetServiceServer = (*targetServer)(nil)
)
