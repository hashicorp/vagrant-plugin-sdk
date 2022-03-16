package core

type BoxMetadataOpts struct {
	Name         string
	Url          string
	Checksum     string
	ChecksumType string
}

type BoxProviderData struct {
	Name         string
	Url          string
	Checksum     string
	ChecksumType string
	Metadata     BoxMetadata
}

type BoxVersionData struct {
	Version     string
	Status      string
	Description string
	Metadata    BoxMetadata
}

type BoxMetadata interface {
	Name() string

	Version(version string, opts *BoxMetadataOpts) (*BoxVersionData, error)
	ListVersions(opts ...*BoxMetadataOpts) ([]string, error)

	Provider(version string, name string) (*BoxProviderData, error)
	ListProviders(version string) ([]string, error)
}
