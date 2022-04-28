package cleanup

import (
	//	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Cleanup is an interface for registering cleanup functions
type Cleanup interface {
	Do(cleanupFn) // registers a cleanup function
	Close() error // executes registered cleanup functions
}

// Create a new cleanup instance
func New() Cleanup {
	return &cleanup{
		fns: []cleanupFn{},
	}
}

type cleanupFn func() error

type cleanup struct {
	fns  []cleanupFn
	l    sync.Mutex
	mark bool
}

// Register a new cleanup task to be performed on close
// NOTE: cleanup tasks are only called once
func (c *cleanup) Do(fn cleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.fns = append(c.fns, fn)
}

// Run all cleanup tasks
func (c *cleanup) Close() (err error) {
	c.l.Lock()
	defer c.l.Unlock()
	// TODO(spox): Uncomment once closers are properly
	//             setup to only be called once.
	// if c.mark {
	// 	return fmt.Errorf("Cleanup has already been closed")
	// }

	for _, fn := range c.fns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	// Remove cleanup tasks as they have been called
	c.fns = []cleanupFn{}

	// Set the mark to show close has been called
	c.mark = true

	return
}
