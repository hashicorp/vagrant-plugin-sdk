package cloud

import (
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestUnauthedClient returns a fully in-memory and side-effect free VagrantCloudClient
// that does not have any auth tokens.
func TestUnauthedClient(t testing.T) *VagrantCloudClient {
	vcc, err := NewVagrantCloudClient(
		"", DEFAULT_RETRY_COUNT, DEFAULT_URL,
	)
	require.NoError(t, err)
	return vcc
}
