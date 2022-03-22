package core

// Type is an enum of all the types of core components supported.
type Type uint

const (
	InvalidType Type = iota // Invalid
	BasisType
	BoxCollectionType
	BoxMetadataType
	BoxType
	MachineType
	PluginManagerType
	ProjectType
	StateBagType
	TargetIndexType
	TargetType
	maxType
)

// TypeMap is a mapping of Type to the nil pointer to the interface of that
// type. This can be used with libraries such as mapper.
var TypeMap = map[Type]interface{}{
	BasisType:         (*Basis)(nil),
	BoxCollectionType: (*BoxCollection)(nil),
	BoxType:           (*Box)(nil),
	MachineType:       (*Machine)(nil),
	PluginManagerType: (*PluginManager)(nil),
	ProjectType:       (*Project)(nil),
	StateBagType:      (*StateBag)(nil),
	TargetIndexType:   (*TargetIndex)(nil),
	TargetType:        (*Target)(nil),
}

type CorePluginManager interface {
	GetPlugin(pluginType Type) (interface{}, error)
}
