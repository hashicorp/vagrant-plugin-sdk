package core

type BoxVersion struct {
	Version   string
	Providers []*BoxProvider
}

type BoxProvider struct {
	Name         string
	URL          string
	Checksum     string
	ChecksumType string
}

type BoxSummary struct {
	Name     string
	Version  string
	Provider string
}

type BoxCollection interface {
	Add(path, name, version string, metadataURL string, providers []*BoxProvider) (box Box, err error)
	All() (list *BoxSummary, err error)
	Find(name string, providers []string, version string) (box Box, err error)
	Clean(name string) (err error)
}
