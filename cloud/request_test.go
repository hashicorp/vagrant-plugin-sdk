package cloud

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
	"github.com/stretchr/testify/require"
)

func simpleServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
}

func validateRequestServer(t *testing.T, validate func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validate(t, r)
	}))
}

type TestingUI struct {
	terminal.UI
	buf bytes.Buffer
}

func NewTestingUI(buf bytes.Buffer) *TestingUI {
	return &TestingUI{UI: terminal.NonInteractiveUI(context.Background()), buf: buf}
}

func (ui *TestingUI) Output(msg string, raw ...interface{}) {
	ui.UI.Output(msg, terminal.WithWriter(&ui.buf), raw)
}

func TestNewVagrantCloudRequest(t *testing.T) {
	_, err := NewVagrantCloudRequest()
	require.Error(t, err)

	ts := simpleServer()
	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
	)
	require.NoError(t, err)
	raw, err := vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
	ts.Close()

	ts = simpleServer()
	vcr, err = NewVagrantCloudRequest(
		WithURL(ts.URL),
		ReplaceHosts(),
	)
	require.NoError(t, err)
	raw, err = vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
	ts.Close()

	ts = simpleServer()
	vcr, err = NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithMethod(PUT),
	)
	require.NoError(t, err)
	raw, err = vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
	ts.Close()
}

func TestWarnDifferentURL(t *testing.T) {
	ts := simpleServer()
	testServerURL, err := url.Parse(ts.URL)
	require.NoError(t, err)
	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WarnDifferentTarget(testServerURL, nil),
	)
	require.NoError(t, err)
	raw, err := vcr.Do()
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
	ts.Close()

	ts = simpleServer()
	testServerURL, err = url.Parse(ts.URL)
	require.NoError(t, err)
	var buf bytes.Buffer
	ui := NewTestingUI(buf)
	vcr, err = NewVagrantCloudRequest(
		WithURL(ts.URL),
		WarnDifferentTarget(testServerURL, ui),
	)
	require.NoError(t, err)
	raw, err = vcr.Do()
	require.Contains(t, ui.buf.String(), "Vagrant has detected a custom Vagrant server in use for downloading")
	require.NoError(t, err)
	require.NotNil(t, raw)
	require.Contains(t, string(raw), "Hello, client")
	ts.Close()
}

func TestAuthTokenRequest(t *testing.T) {
	ts := validateRequestServer(t, func(t *testing.T, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		require.NotEmpty(t, authHeader)
		require.Equal(t, authHeader, "Bearer mytoken")
	})

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithAuthTokenHeader("mytoken"),
	)
	require.NoError(t, err)
	_, err = vcr.Do()
	require.NoError(t, err)
	ts.Close()
}

func TestAuthURLParameterRequest(t *testing.T) {
	ts := validateRequestServer(t, func(t *testing.T, r *http.Request) {
		queryParams := r.URL.Query()
		require.NotEmpty(t, queryParams)
		require.Equal(t, queryParams["access_token"], []string{"mytoken"})
	})

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithAuthTokenURLParam("mytoken"),
	)
	require.NoError(t, err)
	_, err = vcr.Do()
	require.NoError(t, err)
	ts.Close()
}

func TestRequestWithBody(t *testing.T) {
	ts := validateRequestServer(t, func(t *testing.T, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.Equal(t, string(body), "mybody")
	})

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithRequestBody([]byte("mybody")),
	)
	require.NoError(t, err)
	_, err = vcr.Do()
	require.NoError(t, err)
	ts.Close()
}

func TestRequestWithJSONBody(t *testing.T) {
	ts := validateRequestServer(t, func(t *testing.T, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.Equal(t, string(body), "{\"key\":\"myval\"}")
	})

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithRequestJSONableData(map[string]interface{}{
			"key": "myval",
		}),
	)
	require.NoError(t, err)
	_, err = vcr.Do()
	require.NoError(t, err)
	ts.Close()
}

func TestQueryURLParameterRequest(t *testing.T) {
	ts := validateRequestServer(t, func(t *testing.T, r *http.Request) {
		queryParams := r.URL.Query()
		require.NotEmpty(t, queryParams)
		require.Equal(t, queryParams["test"], []string{"val"})
		require.Equal(t, queryParams["othertest"], []string{"otherval"})
	})

	vcr, err := NewVagrantCloudRequest(
		WithURL(ts.URL),
		WithURLQueryParams(map[string]string{
			"test":      "val",
			"othertest": "otherval",
		}),
	)
	require.NoError(t, err)
	_, err = vcr.Do()
	require.NoError(t, err)
	ts.Close()
}
