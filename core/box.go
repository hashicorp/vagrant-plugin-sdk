package core

type Box interface {
	// Check if a box is allowed
	AutomaticUpdateCheckAllowed() (allowed bool, err error)
	// Deletes the box
	Destroy() (err error)
	// The directory on disk where this box exists
	Directory() (path string, err error)
	// Checks if the box has an update
	HasUpdate(version string) (updateAvailable bool, err error)
	// Checks if this box is in use according to the given machine index
	InUse(index TargetIndex) (inUse bool, err error)
	// Returns the machines from the machine index that are using the box
	Machines(index TargetIndex) (machines []Machine, err error)
	// Returns the metadata associated with the box (metadata.json within the box file)
	BoxMetadata() (metadata map[string]interface{}, err error)
	// The metadata information for this box from the metadata url (given by the box repository)
	Metadata() (metadata BoxMetadata, err error)
	// The URL to the version info and other metadata for this box
	MetadataURL() (url string, err error)
	// Box name
	Name() (name string, err error)
	// Box provider
	Provider() (name string, err error)
	// This repackages this box and outputs it to the given path.
	Repackage(path string) (err error)
	// Box version
	Version() (version string, err error)
	// Compares a box to this box. Returns -1, 0, or 1 if this version is smaller, equal, or
	// larger than the other version, respectively.
	Compare(box Box) (int, error)
}
