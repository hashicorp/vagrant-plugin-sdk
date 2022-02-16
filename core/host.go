package core

import (
	"io"
)

type Host interface {
	CapabilityPlatform

	Detect(state StateBag) (bool, error)
	Parent() (string, error)

	io.Closer
}
