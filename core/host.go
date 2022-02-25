package core

import (
	"io"
)

type Host interface {
	CapabilityPlatform
	Seeder
	Named

	Detect(state StateBag) (bool, error)
	Parent() (string, error)

	io.Closer
}
