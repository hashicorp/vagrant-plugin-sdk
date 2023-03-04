// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cleanup

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

type CleanupFn func() error

// Cleanup is an interface for registering cleanup functions
type Cleanup interface {
	Append(CleanupFn)  // registers a cleanup function to run after general stack
	Prepend(CleanupFn) // registers a cleanup function to run before general stack
	Do(CleanupFn)      // registers a cleanup function to run in the general stack
	Close() error      // executes registered cleanup functions
}

// Create a new cleanup instance
func New() Cleanup {
	return &cleanup{
		aFns: []CleanupFn{},
		fns:  []CleanupFn{},
		pFns: []CleanupFn{},
	}
}

type cleanup struct {
	aFns []CleanupFn
	fns  []CleanupFn
	pFns []CleanupFn
	l    sync.Mutex
}

// Register a new cleanup task to be performed on close
// NOTE: cleanup tasks are only called once
func (c *cleanup) Do(fn CleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *cleanup) Append(fn CleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.aFns = append(c.aFns, fn)
}

func (c *cleanup) Prepend(fn CleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.pFns = append(c.pFns, fn)
}

// Run all cleanup tasks
func (c *cleanup) Close() (err error) {
	c.l.Lock()
	defer c.l.Unlock()

	// First run all tasks in the prepend collection
	for _, fn := range c.pFns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	// Next run all tasks in regular collection
	for _, fn := range c.fns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	// Finally run all tasks in append collection
	for _, fn := range c.aFns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	// Remove cleanup tasks as they have been called
	c.aFns = []CleanupFn{}
	c.pFns = []CleanupFn{}
	c.fns = []CleanupFn{}

	return
}
