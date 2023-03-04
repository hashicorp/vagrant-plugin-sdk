// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package pluginargs
package pluginargs

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cacher"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cleanup"
)

// Internal is an interface that is available to mappers.
type Internal interface {
	Broker() *plugin.GRPCBroker
	Cache() cacher.Cache
	Cleanup() cleanup.Cleanup
	Logger() hclog.Logger
	Mappers() []*argmapper.Func
}

// Create a new internal instance
func New(
	broker *plugin.GRPCBroker,
	cache cacher.Cache,
	cleanup cleanup.Cleanup,
	logger hclog.Logger,
	mappers []*argmapper.Func,
) Internal {
	return &internal{
		broker:  broker,
		cache:   cache,
		cleanup: cleanup,
		logger:  logger,
		mappers: mappers,
	}
}

type internal struct {
	broker  *plugin.GRPCBroker
	cache   cacher.Cache
	cleanup cleanup.Cleanup
	logger  hclog.Logger
	mappers []*argmapper.Func
}

// Broker implements Internal
func (i *internal) Broker() *plugin.GRPCBroker {
	return i.broker
}

// Cache implements Internal
func (i *internal) Cache() cacher.Cache {
	return i.cache
}

// Cleanup implements Internal
func (i *internal) Cleanup() cleanup.Cleanup {
	return i.cleanup
}

// Logger implements Internal
func (i *internal) Logger() hclog.Logger {
	return i.logger
}

// Mappers implements Internal
func (i *internal) Mappers() []*argmapper.Func {
	return i.mappers
}
