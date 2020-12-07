package core

// for boxes
type Metadata map[string]interface{}

type Box interface {
	// accessors
	Name() (name string, err error)
	Provider() (name string, err error)
	Version() (version string, err error)
	Directory() (path string, err error)
	Metadata() (metadata Metadata, err error)
	MetadataURL() (url string, err error)

	// action functions
	Destroy() (err error)
	InUse(index MachineIndex) (inUse bool, err error)
	LoadMetadata() (metadata BoxMetadata, err error)
	HasUpdate(version string) (updateAvailable bool, err error)
	AutomaticUpdateCheckAllowed() (allowed bool, err error)
	Repackage() (err error)

	// TODO(spox): Needs comparison function for sorting
}

type BoxMetadata interface {
	// accessors
	Name() (name string, err error)
	Description() (description string, err error)

	// action
	Load(pathOrURL string) (err error)
	Version(version string, providers []string) (v BoxVersion, err error)
	Versions(providers []string) (versions []BoxVersion, err error)
}
