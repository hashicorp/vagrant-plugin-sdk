// Package core has the interfaces for all core types
// used by Vagrant. These are implemented within Vagrant
// and provided to plugins as needed/required. This is
// currently a high level mapping of the interface
// provided by Vagrant in its current state. By doing
// a best effort match (and mapping where things are
// different) it should be easier to integrated with
// Vagrant to keep existing plugins working as well as
// making porting plugins less time consuming.

package core

type PluginType interface {
	PluginName() (string, error)
}
