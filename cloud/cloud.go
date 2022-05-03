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
		return "delete"
	case GET:
		return "get"
	case HEAD:
		return "head"
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
	req, _ := http.NewRequest(method.String(), vc.url, nil)
	// set headers
	req.Header = vc.headers
	// add query parameters
	for k, v := range params {
		req.URL.Query().Add(k, v)
	}
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
	params := map[string]string{
		"q":        query,
		"provider": provider,
		"sort":     sort,
		"order":    order,
		"limit":    strconv.Itoa(limit),
		"page":     strconv.Itoa(page),
	}
	return vc.request("search", GET, params)
}
