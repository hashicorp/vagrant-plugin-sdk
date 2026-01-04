package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	vcc := TestUnauthedClient(t)
	res, err := vcc.Search("hashicorp", "", "", "", 10, 1)
	require.NoError(t, err)
	require.Len(t, res["boxes"], 10)
}
