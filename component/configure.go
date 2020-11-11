package component

import (
	"context"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
)

// Configurable can be optionally implemented by any compontent to
// accept user configuration.
type Configurable interface {
	// Config should return a pointer to an allocated configuration
	// structure. This structure will be written to directly with the
	// decoded configuration. If this returns nil, then it is as if
	// Configurable was not implemented.
	Config() (interface{}, error)
}

// Documented can be optionally implemented by any component to
// return documentation about the component.
type Documented interface {
	// Documentation() returns a completed docs.Documentation struct
	// describing the components configuration.
	Documentation() (*docs.Documentation, error)
}

// ConfigurableNotify is an optional interface that can be implemented
// by any component to receive a notification that the configuration
// was decoded.
type ConfigurableNotify interface {
	Configurable

	// ConfigSet is called with the value of the configuration after
	// decoding is complete successfully.
	ConfigSet(interface{}) error
}


// Configure configures c with the provided configuration.
//
// If c does not implement Configurable AND body is non-empty, then it is
// an error. If body is empty in that case, it is not an error.
func Configure(c interface{}, body string, ctx context.Context) error {
	return nil
}

// Documentation returns the documentation for the given component.
//
// If c does not implement Documented, nil is returned.
func Documentation(c interface{}) (*docs.Documentation, error) {
	if d, ok := c.(Documented); ok {
		return d.Documentation()
	}

	if c, ok := c.(Configurable); ok {
		// Get the configuration value
		v, err := c.Config()

		// If there is no configuration structure for this component,
		// then there is really no documentation, so just return an empty
		// docs structure.
		if err != nil || v == nil {
			return docs.New()
		}

		return docs.New(docs.FromConfig(v))
	}

	return nil, nil
}
