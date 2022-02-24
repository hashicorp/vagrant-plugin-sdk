package core

import (
	"io"
)

type Guest interface {
	CapabilityPlatform
	Seeder

	Detect(Target) (bool, error)
	Parent() (string, error)

	io.Closer
}
