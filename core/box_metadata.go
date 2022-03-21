package core

type BoxProvider struct {
	Name         string
	Url          string
	Checksum     string
	ChecksumType string

	Version *BoxVersion
}

type BoxVersion struct {
	Version     string
	Status      string
	Description string
}

type BoxMetadata interface {
	Name() string

	Version(version string, opts *BoxProvider) (*BoxVersion, error)
	ListVersions(opts ...*BoxProvider) ([]string, error)

	Provider(version string, name string) (*BoxProvider, error)
	ListProviders(version string) ([]string, error)

	Matches(version string, name string, provider *BoxProvider) (bool, error)
	MatchesAny(version string, name string, provider ...*BoxProvider) (bool, error)
}
