package core

import (
	"io"
)

type Guest interface {
	CapabilityPlatform

	Detect(Target) (bool, error)
	Parent() (string, error)

	io.Closer
}
