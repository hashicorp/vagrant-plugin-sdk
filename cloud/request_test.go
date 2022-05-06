package cloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVagrantCloudRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	_, err := NewVagrantCloudRequest()
	require.Error(t, err)

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
	)
	require.NoError(t, err)
	raw, err := vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")

	vcr, err = NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithMethod(PUT),
	)
	require.NoError(t, err)
	raw, err = vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
}
