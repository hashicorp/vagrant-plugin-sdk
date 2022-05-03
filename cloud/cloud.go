package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Type is an enum of all the available http methods
type HTTPMethod int64

const (
	Undef HTTPMethod = iota
	DELETE
	GET
	HEAD
	POST
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
	}
	return "unknown"
}

const (
	DEFAULT_URL            = "https://vagrantcloud.com/api/v1"
	DEFAULT_RETRY_COUNT    = 3
	DEFAULT_RETRY_INTERVAL = 2 // in seconds
)

type VagrantCloudClient struct {
	accessToken   string
	retryCount    int
	retryInterval int
	url           string

	headers http.Header
}

func NewVagrantCloudClient(accessToken string, retryCount int, retryInterval int, url string) (*VagrantCloudClient, error) {
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
		accessToken:   accessToken,
		retryCount:    retryCount,
		retryInterval: retryInterval,
		url:           url,
		headers:       headers,
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
) (map[string]interface{}, error) {
	// Request with query parameters if the HTTPMethod is GET, HEAD or DELETE
	queryParamMethods := []string{DELETE.String(), GET.String(), HEAD.String()}
	if contains(queryParamMethods, method.String()) {
		stringParams := make(map[string]string)
		for k, v := range params {
			stringParams[k] = v.(string)
		}
		return vc.requestWithQueryParams(path, method, stringParams)
	} else {
		return vc.requestWithBody(path, method, params)
	}
}

func (vc *VagrantCloudClient) requestWithBody(
	path string, method HTTPMethod, params map[string]interface{},
) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", vc.url, path)

	// Create the request body
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(method.String(), url, bytes.NewBuffer(jsonBody))

	// Set headers
	req.Header = vc.headers

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func (vc *VagrantCloudClient) requestWithQueryParams(
	path string, method HTTPMethod, params map[string]string,
) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", vc.url, path)

	req, _ := http.NewRequest(method.String(), url, nil)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header = vc.headers

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
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
