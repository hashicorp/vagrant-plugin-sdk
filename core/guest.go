package core

import (
	"io"
)

type Guest interface {
	CapabilityPlatform
	Seeder
	Named

	Detect(Target) (bool, error)
	Parent() (string, error)

	io.Closer
}
