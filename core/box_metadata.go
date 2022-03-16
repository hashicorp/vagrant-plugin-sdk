package core

type BoxMetadataOpts struct {
	Name         string
	Url          string
	Checksum     string
	ChecksumType string
}

type BoxVersionProviderData struct {
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
	Version(version string, opts BoxMetadataOpts) (BoxVersionData, error)
	ListVersions(opts ...BoxMetadataOpts) ([]string, error)

	Provider(version string, name string) (BoxVersionProviderData, error)
	ListProviders(version string) ([]string, error)
}
