package datadir

import (
	"os"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestDir returns a Dir for testing.
func TestDir(t testing.T) (Dir, func()) {
	t.Helper()

	dir, err := newDir("datadir-test")
	require.NoError(t, err)

	return dir, func() {
		dirs := []path.Path{dir.CacheDir(), dir.ConfigDir(),
			dir.DataDir(), dir.TempDir()}
		for _, d := range dirs {
			os.RemoveAll(d.String())
		}
	}
}
