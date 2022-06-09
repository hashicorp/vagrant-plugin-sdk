package plugin

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/mitchellh/protostructure"
// 	"google.golang.org/grpc"
// 	"google.golang.org/protobuf/types/known/emptypb"

// 	"github.com/hashicorp/vagrant-plugin-sdk/component"
// 	"github.com/hashicorp/vagrant-plugin-sdk/docs"
// 	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
// )

// // configStruct is the shared helper to implement the ConfigStruct RPC call
// // for components. The logic is the same regardless of component so this can
// // be called instead.
// func configStruct(impl interface{}) (*vagrant_plugin_sdk.Config_StructResp, error) {
// 	c, ok := impl.(component.Configurable)

// 	// If Configurable isn't implemented, we just return an empty response.
// 	// The nil struct signals to the receiving side that this component
// 	// is not configurable.
// 	if !ok {
// 		return &vagrant_plugin_sdk.Config_StructResp{}, nil
// 	}

// 	v, err := c.Config()
// 	if err != nil {
// 		return nil, err
// 	}

// 	s, err := protostructure.Encode(v)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &vagrant_plugin_sdk.Config_StructResp{Struct: s}, nil
// }

// // configStructCall is the shared helper to call the ConfigStruct RPC call
// // and return the proper struct value for decoding configuration.
// func configStructCall(ctx context.Context, c configurableClient) (interface{}, error) {
// 	resp, err := c.ConfigStruct(ctx, &emptypb.Empty{})

// 	// If we had a failure receiving the configuration struct, then
// 	// panic because this should never happen. In the future maybe we can
// 	// support an error return value.
// 	if err != nil {
// 		return nil, err
// 	}

// 	// If we have no struct, then we have no value so return nil
// 	if resp.Struct == nil {
// 		return nil, nil
// 	}

// 	result, err := protostructure.New(resp.Struct)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// // configure is the shared helper to implement the Configure RPC call.
// func configure(impl interface{}, req *vagrant_plugin_sdk.Config_ConfigureRequest) (*emptypb.Empty, error) {
// 	c, ok := impl.(component.Configurable)

// 	// This should never happen but if it does just do nothing. This
// 	// should never happen because prior to this ever being called, our core
// 	// calls ConfigStruct and if we return nil then we don't configure anything.
// 	if !ok {
// 		return &emptypb.Empty{}, nil
// 	}

// 	// Get our value that we can decode into
// 	v, err := c.Config()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Decode our JSON value directly into our structure.
// 	if err := json.Unmarshal(req.Json, v); err != nil {
// 		return nil, err
// 	}

// 	// If our client also implements the notify interface, call that.
// 	if cn, ok := c.(component.ConfigurableNotify); ok {
// 		if err := cn.ConfigSet(v); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &emptypb.Empty{}, nil
// }

// // configureCall calls the Configure RPC endpoint.
// func configureCall(ctx context.Context, c configurableClient, v interface{}) error {
// 	jsonv, err := json.Marshal(v)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = c.Configure(ctx, &vagrant_plugin_sdk.Config_ConfigureRequest{
// 		Json: jsonv,
// 	})
// 	return err
// }

// // documentation is the shared helper to implement the Documentation RPC call
// // for components. The logic is the same regardless of component so this can
// // be called instead.
// func documentation(impl interface{}) (*vagrant_plugin_sdk.Config_Documentation, error) {
// 	d, err := component.Documentation(impl)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dets := d.Details()

// 	v := &vagrant_plugin_sdk.Config_Documentation{
// 		Description: dets.Description,
// 		Example:     dets.Example,
// 		Input:       dets.Input,
// 		Output:      dets.Output,
// 		Fields:      make(map[string]*vagrant_plugin_sdk.Config_FieldDocumentation),
// 	}

// 	for _, f := range d.Fields() {
// 		v.Fields[f.Field] = &vagrant_plugin_sdk.Config_FieldDocumentation{
// 			Name:     f.Field,
// 			Type:     f.Type,
// 			Default:  f.Default,
// 			Synopsis: f.Synopsis,
// 			Summary:  f.Summary,
// 			EnvVar:   f.EnvVar,
// 			Optional: f.Optional,
// 		}
// 	}

// 	for _, m := range dets.Mappers {
// 		v.Mappers = append(v.Mappers, &vagrant_plugin_sdk.Config_MapperDocumentation{
// 			Input:       m.Input,
// 			Output:      m.Output,
// 			Description: m.Description,
// 		})
// 	}

// 	return v, nil
// }

// // configStructCall is the shared helper to call the ConfigStruct RPC call
// // and return the proper struct value for decoding configuration.
// func documentationCall(ctx context.Context, c configurableClient) (*docs.Documentation, error) {
// 	resp, err := c.Documentation(ctx, &emptypb.Empty{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	d, err := docs.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	d.Example(resp.Example)
// 	d.Description(resp.Description)
// 	d.Input(resp.Input)
// 	d.Output(resp.Output)

// 	for _, f := range resp.Fields {
// 		d.OverrideField(&docs.FieldDocs{
// 			Field:    f.Name,
// 			Type:     f.Type,
// 			Default:  f.Default,
// 			Synopsis: f.Synopsis,
// 			Summary:  f.Summary,
// 			Optional: f.Optional,
// 			EnvVar:   f.EnvVar,
// 		})
// 	}

// 	for _, m := range resp.Mappers {
// 		d.AddMapper(m.Input, m.Output, m.Description)
// 	}

// 	return d, nil
// }

// // configurableClient is the interface implemented by all gRPC services that
// // have the configuration RPC methods. We use this with the helpers above
// // to extract shared logic for component configuration.
// type configurableClient interface {
// 	ConfigStruct(context.Context, *emptypb.Empty, ...grpc.CallOption) (*vagrant_plugin_sdk.Config_StructResp, error)
// 	Configure(context.Context, *vagrant_plugin_sdk.Config_ConfigureRequest, ...grpc.CallOption) (*emptypb.Empty, error)
// 	Documentation(context.Context, *emptypb.Empty, ...grpc.CallOption) (*vagrant_plugin_sdk.Config_Documentation, error)
// }
