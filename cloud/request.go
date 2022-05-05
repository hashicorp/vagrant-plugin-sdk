package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	TARGET_HOST = "vagrantcloud.com"
)

// Type is an enum of all the available http methods
type HTTPMethod int64

const (
	Undef HTTPMethod = iota
	DELETE
	GET
	HEAD
	POST
	PUT
)

func (m HTTPMethod) String() string {
	switch m {
	case DELETE:
		return "DELETE"
	case GET:
		return "GET"
	case HEAD:
		return "HEAD"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	}
	return "unknown"
}

type VagrantCloudRequest struct {
	headers        http.Header
	method         HTTPMethod
	retryCount     int
	requestBody    []byte
	url            string
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
	return
}

func (vcr *VagrantCloudRequest) Do() (raw []byte, err error) {
	client := retryablehttp.NewClient()
	client.RetryMax = vcr.retryCount
	var req *retryablehttp.Request

	// Create request with request body if one is provided
	if vcr.requestBody != nil {
		req, err = retryablehttp.NewRequest(
			vcr.method.String(), vcr.url, bytes.NewBuffer(vcr.requestBody),
		)
		if err != nil {
			return nil, err
		}
	} else {
		// If no request body is provided then create an empty request
		req, err = retryablehttp.NewRequest(vcr.method.String(), vcr.url, nil)
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
	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type VagrantCloudRequestOptions func(*VagrantCloudRequest) error

func WithAuthTokenHeader(t string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.headers.Set("Authorization", fmt.Sprintf("Bearer %s", t))
		return
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
		r.url = u
		return
	}
}

func WithURLQueryParams(u map[string]string) VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		r.urlQueryParams = u
		return
	}
}

func ReplaceHosts() VagrantCloudRequestOptions {
	return func(r *VagrantCloudRequest) (err error) {
		newUrl, err := replaceHostUrl(r.url)
		if err != nil {
			return err
		}
		r.url = newUrl
		return
	}
}

func replaceHostUrl(urlIn string) (string, error) {
	replacementHosts := []string{"app.vagrantup.com", "atlas.hashicorp.com"}
	parsedUrl, err := url.Parse(urlIn)
	if err != nil {
		return "", err
	}
	// Replace the url host name with the TARGET_HOST if it is one of the hosts
	// that should be replaced (eg. in the replacementHosts list).
	if parsedUrl.Host != TARGET_HOST && contains(replacementHosts, parsedUrl.Host) {
		parsedUrl.Host = TARGET_HOST
	}
	return parsedUrl.String(), nil
}
