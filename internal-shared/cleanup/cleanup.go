package cleanup

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Cleanup is an interface for registering cleanup functions
type Cleanup interface {
	Do(cleanupFn) // registers a cleanup function
	Close() error // executes registered cleanup functions
}

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

func (c *cleanup) Do(fn cleanupFn) {
	c.l.Lock()
	defer c.l.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *cleanup) Close() (err error) {
	c.l.Lock()
	defer c.l.Unlock()
	if c.mark {
		return fmt.Errorf("Cleanup has already been closed")
	}

	for _, fn := range c.fns {
		e := fn()
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	c.mark = true

	return
}
