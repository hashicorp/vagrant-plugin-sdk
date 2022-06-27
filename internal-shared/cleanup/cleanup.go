package cleanup

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

type CleanupFn func() error

// Cleanup is an interface for registering cleanup functions
type Cleanup interface {
	Prepend(CleanupFn) // registers a cleanup function at start of current cleanup stack
	Do(CleanupFn)      // registers a cleanup function at the end of the current cleanup stack
	Close() error      // executes registered cleanup functions
}

// Create a new cleanup instance
func New() Cleanup {
	return &cleanup{
		fns: []CleanupFn{},
	}
}

type cleanup struct {
	fns []CleanupFn
	l   sync.Mutex
}

// Register a new cleanup task to be performed on close
// NOTE: cleanup tasks are only called once
func (c *cleanup) Do(fn CleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *cleanup) Prepend(fn CleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	fns := make([]CleanupFn, len(c.fns)+1)
	fns[0] = fn
	copy(fns[1:], c.fns)
	c.fns = fns
}

// Run all cleanup tasks
func (c *cleanup) Close() (err error) {
	c.l.Lock()
	defer c.l.Unlock()
	for _, fn := range c.fns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	// Remove cleanup tasks as they have been called
	c.fns = []CleanupFn{}

	return
}
