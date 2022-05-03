package cloud

import (
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
)

func (m HTTPMethod) String() string {
	switch m {
	case DELETE:
		return "DELETE"
	case GET:
		return "GET"
	case HEAD:
		return "HEAD"
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

func (vc *VagrantCloudClient) request(
	path string, method HTTPMethod, params map[string]string,
) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", vc.url, path)
	req, _ := http.NewRequest(method.String(), url, nil)

	// Set headers
	req.Header = vc.headers

	// Add query parameters
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

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

func (vc *VagrantCloudClient) Seach(
	query string, provider string, sort string, order string, limit int, page int,
) (map[string]interface{}, error) {
	params := make(map[string]string)
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
