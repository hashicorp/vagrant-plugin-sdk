package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
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

const (
	DEFAULT_URL            = "https://vagrantcloud.com/api/v1"
	DEFAULT_RETRY_COUNT    = 3
	DEFAULT_RETRY_INTERVAL = 2 // in seconds
)

type CloudClient interface {
	AuthTokenCreate(username string, password string, description string, code string) (map[string]interface{}, error)
	AuthTokenDelete() (map[string]interface{}, error)
	AuthRequest2faCode(username string, password string, delivery_method string) (map[string]interface{}, error)
	AuthTokenValidate() (map[string]interface{}, error)
	BoxCreate(username string, name string, shortDescription string, description string, isPrivate bool) (map[string]interface{}, error)
	BoxDelete(username string, name string) (map[string]interface{}, error)
	BoxGet(username string, name string) (map[string]interface{}, error)
	BoxUpdate(username string, name string, shortDescription string, description string, isPrivate bool) (map[string]interface{}, error)
	BoxVersionGet(username string, name string, version string) (map[string]interface{}, error)
	BoxVersionCreate(username string, name string, version string, description string) (map[string]interface{}, error)
	BoxVersionUpdate(username string, name string, version string, description string) (map[string]interface{}, error)
	BoxVersionDelete(username string, name string, version string) (map[string]interface{}, error)
	BoxVersionRelease(username string, name string, version string) (map[string]interface{}, error)
	BoxVersionRevoke(username string, name string, version string) (map[string]interface{}, error)
	BoxVersionProviderCreate(username string, name string, version string, provider string, url string, checksum string, checksumType string) (map[string]interface{}, error)
	BoxVersionProviderDelete(username string, name string, version string, provider string) (map[string]interface{}, error)
	BoxVersionProviderGet(username string, name string, version string, provider string) (map[string]interface{}, error)
	BoxVersionProviderUpdate(username string, name string, version string, provider string, url string, checksum string, checksumType string) (map[string]interface{}, error)
	BoxVersionProviderUpload(username string, name string, version string, provider string) (map[string]interface{}, error)
	BoxVersionProviderUploadDirect(username string, name string, version string, provider string) (map[string]interface{}, error)
	OrganizationGet(name string) (map[string]interface{}, error)
	Seach(query string, provider string, sort string, order string, limit int, page int) (map[string]interface{}, error)
}

type VagrantCloudClient struct {
	accessToken string
	headers     http.Header
	retryCount  int
	url         string
}

func NewVagrantCloudClient(accessToken string, retryCount int, url string) (*VagrantCloudClient, error) {
	// Set default url if none is provided
	if url == "" {
		url = DEFAULT_URL
	}
	// Set default headers
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	if accessToken != "" {
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	client := &VagrantCloudClient{
		accessToken: accessToken,
		retryCount:  retryCount,
		url:         url,
		headers:     headers,
	}
	return client, nil
}

func contains(one []string, two string) bool {
	for _, v := range one {
		if v == two {
			return true
		}
	}

	return false
}

func (vc *VagrantCloudClient) request(
	path string, method HTTPMethod, params map[string]interface{},
) (jsonResp map[string]interface{}, err error) {
	var resp *http.Response
	// Request with query parameters if the HTTPMethod is GET, HEAD or DELETE
	queryParamMethods := []string{DELETE.String(), GET.String(), HEAD.String()}
	if contains(queryParamMethods, method.String()) {
		stringParams := make(map[string]string)
		for k, v := range params {
			stringParams[k] = v.(string)
		}
		resp, err = vc.requestWithQueryParams(path, method, stringParams)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = vc.requestWithBody(path, method, params)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func (vc *VagrantCloudClient) requestWithBody(
	path string, method HTTPMethod, params map[string]interface{},
) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = vc.retryCount
	url := fmt.Sprintf("%s/%s", vc.url, path)

	// Create the request body
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, _ := retryablehttp.NewRequest(
		method.String(), url, bytes.NewBuffer(jsonBody),
	)

	// Set headers
	req.Header = vc.headers

	// Execute request
	return client.Do(req)
}

func (vc *VagrantCloudClient) requestWithQueryParams(
	path string, method HTTPMethod, params map[string]string,
) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = vc.retryCount
	url := fmt.Sprintf("%s/%s", vc.url, path)

	// Add query parameters to request
	req, _ := retryablehttp.NewRequest(method.String(), url, nil)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header = vc.headers

	// Execute request
	return client.Do(req)
}

func (vc *VagrantCloudClient) AuthTokenCreate(
	username string, password string, description string, code string,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	params["user"] = map[string]string{
		"login":    username,
		"password": password,
	}
	if description != "" {
		params["token"] = map[string]string{
			"description": description,
		}
	}
	if code != "" {
		params["two_factor"] = map[string]string{
			"code": code,
		}
	}
	return vc.request("authenticate", POST, params)
}

func (vc *VagrantCloudClient) AuthTokenDelete() (map[string]interface{}, error) {
	return vc.request("authenticate", DELETE, nil)
}

func (vc *VagrantCloudClient) AuthRequest2faCode(
	username string, password string, delivery_method string,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	params["user"] = map[string]string{
		"login":    username,
		"password": password,
	}
	params["two_factor"] = map[string]string{
		"delivery_method": delivery_method,
	}
	return vc.request("two-factor/request-code", POST, params)
}

func (vc *VagrantCloudClient) AuthTokenValidate() (map[string]interface{}, error) {
	return vc.request("authenticate", GET, nil)
}

func (vc *VagrantCloudClient) BoxCreate(
	username string, name string, shortDescription string, description string, isPrivate bool,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	params["username"] = username
	params["name"] = name
	if shortDescription != "" {
		params["short_description"] = shortDescription
	}
	if description != "" {
		params["description"] = description
	}
	params["is_private"] = strconv.FormatBool(isPrivate)
	return vc.request("boxes", POST, params)
}

func (vc *VagrantCloudClient) BoxDelete(
	username string, name string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s", username, name)
	return vc.request(path, DELETE, nil)
}

func (vc *VagrantCloudClient) BoxGet(
	username string, name string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s", username, name)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) BoxUpdate(
	username string, name string, shortDescription string, description string, isPrivate bool,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if shortDescription != "" {
		params["short_description"] = shortDescription
	}
	if description != "" {
		params["description"] = description
	}
	params["is_private"] = strconv.FormatBool(isPrivate)
	path := fmt.Sprintf("box/%s/%s", username, name)
	return vc.request(path, PUT, params)
}

func (vc *VagrantCloudClient) BoxVersionGet(
	username string, name string, version string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s", username, name, version)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) BoxVersionCreate(
	username string, name string, version string, description string,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	versionHash := map[string]string{
		"version": version,
	}
	if description != "" {
		versionHash["description"] = description
	}
	params["version"] = versionHash
	path := fmt.Sprintf("box/%s/%s/version/%s", username, name, version)
	return vc.request(path, POST, nil)
}

func (vc *VagrantCloudClient) BoxVersionUpdate(
	username string, name string, version string, description string,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	versionHash := map[string]string{
		"version": version,
	}
	if description != "" {
		versionHash["description"] = description
	}
	params["version"] = versionHash
	path := fmt.Sprintf("box/%s/%s/version/%s", username, name, version)
	return vc.request(path, PUT, nil)
}

func (vc *VagrantCloudClient) BoxVersionDelete(
	username string, name string, version string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s", username, name, version)
	return vc.request(path, DELETE, nil)
}

func (vc *VagrantCloudClient) BoxVersionRelease(
	username string, name string, version string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/release", username, name, version)
	return vc.request(path, PUT, nil)
}

func (vc *VagrantCloudClient) BoxVersionRevoke(
	username string, name string, version string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/revoke", username, name, version)
	return vc.request(path, PUT, nil)
}

func (vc *VagrantCloudClient) BoxVersionProviderCreate(
	username string, name string, version string, provider string, url string, checksum string, checksumType string,
) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"provider": map[string]string{
			"name":          provider,
			"url":           url,
			"checksum":      checksum,
			"checksum_type": checksumType,
		},
	}

	path := fmt.Sprintf("box/%s/%s/version/%s/providers", username, name, version)
	return vc.request(path, POST, params)
}

func (vc *VagrantCloudClient) BoxVersionProviderDelete(
	username string, name string, version string, provider string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/provider/%s", username, name, version, provider)
	return vc.request(path, DELETE, nil)
}

func (vc *VagrantCloudClient) BoxVersionProviderGet(
	username string, name string, version string, provider string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/provider/%s", username, name, version, provider)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) BoxVersionProviderUpdate(
	username string, name string, version string, provider string, url string, checksum string, checksumType string,
) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"provider": map[string]string{
			"name":          provider,
			"url":           url,
			"checksum":      checksum,
			"checksum_type": checksumType,
		},
	}

	path := fmt.Sprintf("box/%s/%s/version/%s/provider/%s", username, name, version, provider)
	return vc.request(path, PUT, params)
}

func (vc *VagrantCloudClient) BoxVersionProviderUpload(
	username string, name string, version string, provider string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/provider/%s/upload", username, name, version, provider)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) BoxVersionProviderUploadDirect(
	username string, name string, version string, provider string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("box/%s/%s/version/%s/provider/%s/upload/direct", username, name, version, provider)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) OrganizationGet(
	name string,
) (map[string]interface{}, error) {
	path := fmt.Sprintf("user/%s", name)
	return vc.request(path, GET, nil)
}

func (vc *VagrantCloudClient) Seach(
	query string, provider string, sort string, order string, limit int, page int,
) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if query != "" {
		params["q"] = query
	}
	if provider != "" {
		params["provider"] = provider
	}
	if sort != "" {
		params["sort"] = sort
	}
	if order != "" {
		params["order"] = order
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	return vc.request("search", GET, params)
}

var (
	_ CloudClient = (*VagrantCloudClient)(nil)
)
