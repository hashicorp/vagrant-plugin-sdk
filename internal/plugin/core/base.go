package core

import (
	"errors"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
)

type base struct {
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Cleanup *pluginargs.Cleanup
}

func (b *base) GRPCBroker() *plugin.GRPCBroker {
	return b.Broker
}

// This runs any registered cleanup items that may have been
// generated by the mappers (or anything else). This should
// be called when complete with the core interface in use
func (b *base) Close() error {
	return b.Cleanup.Close()
}

// Map a value to the expected type using registered mappers
// NOTE: The expected type must be a pointer, so an expected type
// of `*int` means an `int` is wanted. Expected type of `**int`
// means an `*int` is wanted, etc.
func (b *base) Map(resultValue, expectedType interface{}, args ...argmapper.Arg) (interface{}, error) {
	typPtr := reflect.TypeOf(expectedType)
	if typPtr.Kind() != reflect.Ptr {
		return nil, errors.New("expectedType must be nil pointer")
	}
	typ := typPtr.Elem()

	vIn := argmapper.Value{Type: typ}
	vOut := argmapper.Value{Type: typ}
	vsIn, err := argmapper.NewValueSet([]argmapper.Value{vIn})
	if err != nil {
		return nil, err
	}
	vsOut, err := argmapper.NewValueSet([]argmapper.Value{vOut})
	if err != nil {
		return nil, err
	}

	cb := func(in, out *argmapper.ValueSet) error {
		val := in.Typed(typ).Value.Interface()
		out.Typed(typ).Value = reflect.ValueOf(val)
		return nil
	}

	callFn, err := argmapper.BuildFunc(vsIn, vsOut, cb)
	if err != nil {
		return nil, err
	}

	args = append(args,
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.Typed(resultValue),
		argmapper.Typed(b.internal()),
		argmapper.Typed(b.Logger),
		argmapper.Logger(b.Logger),
	)

	if err = vsOut.FromResult(callFn.Call(args...)); err != nil {
		return nil, err
	}

	return vsOut.Typed(typ).Value.Interface(), nil
}

func (b *base) internal() *pluginargs.Internal {
	return &pluginargs.Internal{
		Broker:  b.Broker,
		Mappers: b.Mappers,
		Cleanup: b.Cleanup,
	}
}
