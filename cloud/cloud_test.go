package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	vcc := TestUnauthedClient(t)
	res, err := vcc.Seach("hashicorp", "", "", "", 10, 1)
	require.NoError(t, err)
	require.Len(t, res["boxes"], 10)

	res2, err := vcc.Seach("hashicorp", "", "", "", 10, 2)
	require.NoError(t, err)
	require.Len(t, res2["boxes"], 10)
	require.NotEqual(t, res, res2)
}

func TestAuthTokenCreate(t *testing.T) {
	vcc := TestUnauthedClient(t)
	resp, err := vcc.AuthTokenCreate("username", "password", "testtoken", "")
	require.NoError(t, err)
	require.False(t, resp["success"].(bool))
}
