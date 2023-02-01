// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs" // TODO(spox): this is unmaintained - look to vendor internally
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type parser interface {
	ParseVagrantfile(string) (*vagrant_plugin_sdk.Args_Hash, error)
}

type VagrantfileFormat uint8

const (
	JSON VagrantfileFormat = iota
	HCL
	RUBY
)

var Mappers []*argmapper.Func

// Load Vagrant configuration using a Ruby based Vagrantfile
func LoadRubyVagrantfile(
	p path.Path, // path to Ruby Vagrantfile
	rubyRuntime parser, // ruby runtime plugin
) (*Vagrantfile, error) {
	return nil, fmt.Errorf("not implemented")
	// s, err := rubyRuntime.ParseVagrantfile(p.String())
	// if err != nil {
	// 	return nil, err
	// }

	// return LoadVagrantfile(s.Json, p.String(), JSON)
}

// Load Vagrant configuration using an HCL based Vagrantfile
func LoadHCLVagrantfile(
	p path.Path, // path to HCL Vagrantfile
) (*Vagrantfile, error) {
	f, err := os.Open(p.String())
	if err != nil {
		return nil, err
	}

	c, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return LoadVagrantfile(c, p.String(), HCL)
}

// Load a Vagrantfile
func LoadVagrantfile(
	content []byte, // Vagrantfile content
	loc string, // path of file imported
	kind VagrantfileFormat, // type of content (JSON or HCL file)
) (*Vagrantfile, error) {
	var f *hcl.File
	var d hcl.Diagnostics
	switch kind {
	case JSON:
		f, d = hcljson.Parse(content, loc)
	case HCL:
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("invalid Vagrantfile format defined")
	}
	if d.HasErrors() {
		return nil, d
	}

	v := &Vagrantfile{}
	d = gohcl.DecodeBody(f.Body, &hcl.EvalContext{}, v)
	if d.HasErrors() {
		// we never error!
		// return nil, d
	}

	return v, nil
}

// Decode a proto encoded Vagrantfile
func DecodeVagrantfile(
	data *vagrant_plugin_sdk.Args_Hash,
	args ...argmapper.Arg,
) (v *Vagrantfile, err error) {
	args = append(args, argmapper.ConverterFunc(Mappers...))
	mapped, err := dynamic.Map(data,
		(*map[string]interface{})(nil),
		args...,
	)
	if err != nil {
		return
	}
	v = &Vagrantfile{}
	err = mapstructure.Decode(mapped, v)

	return
}

// Encode the Vagrantfile for storage
func EncodeVagrantfile(
	v *Vagrantfile,
	args ...argmapper.Arg,
) (*vagrant_plugin_sdk.Args_Hash, error) {
	args = append(args, argmapper.ConverterFunc(Mappers...))
	mapped := structs.New(v)
	raw, err := dynamic.Map(mapped,
		(**vagrant_plugin_sdk.Args_Hash)(nil),
		args...,
	)
	if err != nil {
		return nil, err
	}

	return raw.(*vagrant_plugin_sdk.Args_Hash), nil
}

func DecodeConfiguration(
	data *vagrant_plugin_sdk.Args_Hash,
	i interface{},
	args ...argmapper.Arg,
) (err error) {
	args = append(args, argmapper.ConverterFunc(Mappers...))
	mapped, err := dynamic.Map(data,
		(*map[string]interface{})(nil),
		args...,
	)
	if err != nil {
		return
	}
	err = mapstructure.Decode(mapped, i)

	return
}

// Restore an encoded Vagrantfile
func RestoreVagrantfile(
	s []byte,
	args ...argmapper.Arg,
) (v *Vagrantfile, err error) {
	if s == nil {
		return &Vagrantfile{}, nil
	}

	v = &Vagrantfile{}
	err = RestoreConfiguration(s, v)
	return
}

// Encode a piece of the configuration
func EncodeConfiguration(
	c interface{}, // configuration to encode
	args ...argmapper.Arg,
) (*vagrant_plugin_sdk.Args_Hash, error) {
	args = append(args, argmapper.ConverterFunc(Mappers...))
	raw, err := dynamic.Map(c,
		(**vagrant_plugin_sdk.Args_Hash)(nil),
		args...,
	)
	if err != nil {
		return nil, err
	}

	return raw.(*vagrant_plugin_sdk.Args_Hash), nil
}

// Restore a piece of the confiugration
func RestoreConfiguration(
	s []byte, // encoded configuration
	i interface{}, // decoding destination
) (err error) {
	if s == nil {
		return nil
	}
	d := map[string]interface{}{}
	if err = json.Unmarshal(s, &d); err != nil {
		return
	}

	return mapstructure.Decode(d, i)
}
