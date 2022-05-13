package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

const (
	TARGET_HOST             = "vagrantcloud.com"
	CUSTOM_HOST_NOTIFY_WAIT = 5 // time in seconds
)

// Type is an enum of all the available http methods
type HTTPMethod int64

const (
	GET HTTPMethod = iota
	DELETE
	HEAD
	POST
	PUT
)

type VagrantCloudRequest struct {
	headers        http.Header
	method         HTTPMethod
	retryCount     int
	requestBody    []byte
	url            *url.URL
	urlQueryParams map[string]string
}

func NewVagrantCloudRequest(opts ...VagrantCloudRequestOptions) (r *VagrantCloudRequest, err error) {
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	r = &VagrantCloudRequest{
		headers: headers,
	}
	for _, opt := range opts {
		if oerr := opt(r); oerr != nil {
			err = multierror.Append(err, oerr)
		}
	}
	if r.url == nil {
		return nil, errors.New("no URL provided")
	}
	return
}

func (vcr *VagrantCloudRequest) Do() (raw []byte, err error) {
	client := retryablehttp.NewClient()
	client.RetryMax = vcr.retryCount
	var req *retryablehttp.Request

	// Create request with request body if one is provided
	if vcr.requestBody != nil {
		req, err = retryablehttp.NewRequest(
			vcr.method.String(), vcr.url.String(), bytes.NewBuffer(vcr.requestBody),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// If no request body is provided then create an empty request
		req, err = retryablehttp.NewRequest(vcr.method.String(), vcr.url.String(), nil)
	}

	// Add query params if provided
	if vcr.urlQueryParams != nil {
		q := req.URL.Query()
		for k, v := range vcr.urlQueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	// Set headers
	req.Header = vcr.headers
	// Add headers to redirects
	client.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type VagrantCloudRequestOptions func(*VagrantCloudRequest) error

func ReplaceHosts() VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.url = replaceHostUrl(r.url)
		return
	}
}

func WarnDifferentTarget(serverUrl *url.URL, ui terminal.UI) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		if serverUrl.Host == r.url.Host {
			if serverUrl.Host != TARGET_HOST {
				if ui != nil {
					// TODO: This output should go through the localization module
					ui.Output(fmt.Sprintf(`Vagrant has detected a custom Vagrant server in use for downloading
				box files. An authentication token is currently set which will be
				added to the box request. If the custom Vagrant server should not
				be receiving the authentication token, please unset it.

					Known Vagrant server:  %s
					Custom Vagrant server: %s`, TARGET_HOST, serverUrl.Host))
					time.Sleep(CUSTOM_HOST_NOTIFY_WAIT * time.Second)
				}
			}
		}
		return
	}
}

func WithAuthTokenHeader(t string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		if t == "" {
			return
		}
		r.headers.Set("Authorization", fmt.Sprintf("Bearer %s", t))
		return
	}
}

func WithAuthTokenURLParam(t string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		if t == "" {
			return
		}
		if r.urlQueryParams == nil {
			r.urlQueryParams = make(map[string]string)
		}
		r.urlQueryParams["access_token"] = t
		return nil
	}
}

func WithMethod(m HTTPMethod) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.method = m
		return
	}
}

func WithRetryCount(c int) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.retryCount = c
		return
	}
}

func WithRequestBody(b []byte) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.requestBody = b
		return
	}
}

func WithRequestJSONableData(d map[string]interface{}) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		jsonBody, err := json.Marshal(d)
		if err != nil {
			return err
		}
		r.requestBody = jsonBody
		return
	}
}

func WithURL(u string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		requestedUrl, err := url.Parse(u)
		if err != nil {
			return err
		}
		r.url = requestedUrl
		return
	}
}

func WithURLQueryParams(u map[string]string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.urlQueryParams = u
		return
	}
}

func replaceHostUrl(url *url.URL) *url.URL {
	replacementHosts := []string{"app.vagrantup.com", "atlas.hashicorp.com"}
	// Replace the url host name with the TARGET_HOST if it is one of the hosts
	// that should be replaced (eg. in the replacementHosts list).
	if url.Host != TARGET_HOST && contains(replacementHosts, url.Host) {
		url.Host = TARGET_HOST
	}
	return url
}
