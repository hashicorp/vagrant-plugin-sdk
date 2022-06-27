package core

type Downloader interface {
	// component.Configurable

	Download() error
}
