package core

import "github.com/hashicorp/vagrant-plugin-sdk/component"

type Downloader interface {
	component.Configurable

	Download() error
}
