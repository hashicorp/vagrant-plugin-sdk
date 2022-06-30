package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/stretchr/testify/require"
)

func TestFindPathEmptyDir(t *testing.T) {
	t.Run("returns nil for empty dir", func(t *testing.T) {
		require := require.New(t)

		dir, err := os.MkdirTemp("", "configtest")
		require.NoError(err)

		filenames := []string{"Vagrantfile"}

		p, err := FindPath(path.NewPath(dir), filenames)
		require.NoError(err)
		require.Nil(p)
	})

	t.Run("returns path to Vagrantfile when it exists", func(t *testing.T) {
		require := require.New(t)

		dir, err := os.MkdirTemp("", "configtest")
		require.NoError(err)

		vagrantfilePath := filepath.Join(dir, "Vagrantfile")
		f, err := os.Create(vagrantfilePath)
		require.NoError(err)
		require.NoError(f.Close())

		filenames := []string{"Vagrantfile"}

		p, err := FindPath(path.NewPath(dir), filenames)
		require.NoError(err)
		require.Equal(vagrantfilePath, p.String())
	})

	t.Run("can find a vagrantfile in an ancestor dir", func(t *testing.T) {
		require := require.New(t)

		dir, err := os.MkdirTemp("", "configtest")
		require.NoError(err)

		vagrantfilePath := filepath.Join(dir, "Vagrantfile")
		f, err := os.Create(vagrantfilePath)
		require.NoError(err)
		require.NoError(f.Close())

		subdir := filepath.Join(dir, "a", "b")
		require.NoError(os.MkdirAll(subdir, 0750))

		filenames := []string{"Vagrantfile"}

		p, err := FindPath(path.NewPath(subdir), filenames)
		require.NoError(err)
		require.Equal(vagrantfilePath, p.String())
	})
}
