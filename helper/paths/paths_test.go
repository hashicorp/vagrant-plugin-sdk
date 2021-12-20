package paths

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVagrantCwd(t *testing.T) {
	t.Run("uses actual cwd when env var is unset", func(t *testing.T) {
		require := require.New(t)

		oldVcwd, ok := os.LookupEnv("VAGRANT_CWD")
		if ok {
			os.Unsetenv("VAGRANT_CWD")
			defer os.Setenv("VAGRANT_CWD", oldVcwd)
		}

		dir, err := ioutil.TempDir("", "test")
		require.NoError(err)
		defer os.RemoveAll(dir)

		oldCwd, err := os.Getwd()
		require.NoError(err)
		os.Chdir(dir)
		defer os.Chdir(oldCwd)

		out, err := VagrantCwd()
		require.NoError(err)
		require.Equal(dir, out.String())
	})

	t.Run("honors VAGRANT_CWD if it's set and exists", func(t *testing.T) {
		require := require.New(t)

		dir, err := ioutil.TempDir("", "test")
		require.NoError(err)
		defer os.RemoveAll(dir)

		os.Setenv("VAGRANT_CWD", dir)
		defer os.Unsetenv("VAGRANT_CWD")

		out, err := VagrantCwd()
		require.NoError(err)
		require.Equal(dir, out.String())
	})

	t.Run("errors if VAGRANT_CWD is set and does not exist", func(t *testing.T) {
		require := require.New(t)

		os.Setenv("VAGRANT_CWD", filepath.Join(os.TempDir(), "idontexit"))
		defer os.Unsetenv("VAGRANT_CWD")

		_, err := VagrantCwd()
		require.Error(err)
		require.Contains(err.Error(), "does not exist")
	})
}
