package core

type DownloaderOptions func(*Downloader) error

type Downloader interface {
	Source() (string, error)
	Destination() (string, error)
	Download() error
}
