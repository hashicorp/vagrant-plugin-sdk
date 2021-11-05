package core

// for boxes
type BoxMetadataMap map[string]interface{}

type Box interface {
	AutomaticUpdateCheckAllowed() (allowed bool, err error)
	Destroy() (err error)
	Directory() (path string, err error)
	HasUpdate(version string) (updateAvailable bool, err error)
	InUse(index TargetIndex) (inUse bool, err error)
	Metadata() (metadata BoxMetadataMap, err error)
	MetadataURL() (url string, err error)
	Name() (name string, err error)
	Provider() (name string, err error)
	Repackage() (err error)
	Version() (version string, err error)

	// TODO(spox): Needs comparison function for sorting
}

type BoxMetadata interface {
	Description() (description string, err error)
	Load(pathOrURL string) (err error)
	Name() (name string, err error)
	Version(version string, providers []string) (v BoxVersion, err error)
	Versions(providers []string) (versions []BoxVersion, err error)
}
