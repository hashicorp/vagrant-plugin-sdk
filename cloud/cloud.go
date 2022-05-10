package cloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

const (
	DEFAULT_SERVER_URL       = "https://vagrantcloud.com"
	DEFAULT_API_URL          = "https://vagrantcloud.com/api/v1"
	DEFAULT_RETRY_COUNT      = 3
	VAGRANT_LOGIN_TOKEN_FILE = "vagrant_login_token"
)

func ServerUrl() (url *url.URL, err error) {
	vagrantServerUrl := os.Getenv("VAGRANT_SERVER_URL")
	if vagrantServerUrl == "" {
		url, err = url.Parse(DEFAULT_SERVER_URL)
	} else {
		url, err = url.Parse(vagrantServerUrl)
	}
	return
}

func VagrantCloudToken(dataDir path.Path, ui terminal.UI) (token string, err error) {
	tokenPath := ""
	loginTokenFile := dataDir.Join(VAGRANT_LOGIN_TOKEN_FILE)
	if _, err := os.Stat(loginTokenFile.String()); err == nil {
		tokenPath = loginTokenFile.String()
	}

	vagrantCloudToken := os.Getenv("VAGRANT_CLOUD_TOKEN")
	if tokenPath != "" && vagrantCloudToken != "" {
		if ui != nil {
			// TODO: this output should go through the localization service
			ui.Output(`Vagrant detected both the VAGRANT_CLOUD_TOKEN environment variable and a Vagrant login
	token are present on this system. The VAGRANT_CLOUD_TOKEN environment variable takes
	precedence over the locally stored token. To remove this error, either unset
	the VAGRANT_CLOUD_TOKEN environment variable or remove the login token stored on disk:

			%s`, loginTokenFile)
		}
	}

	// Set the token
	if vagrantCloudToken != "" {
		token = vagrantCloudToken
	} else if tokenPath != "" {
		tokenData, err := os.ReadFile(tokenPath)
		if err == nil {
			token = string(tokenData)
		}
	} else if os.Getenv("ATLAS_TOKEN") != "" {
		token = os.Getenv("ATLAS_TOKEN")
	}
	return
}

type VagrantCloudClient struct {
	accessToken string
	retryCount  int
	url         string
}

type VagrantCloudClientOptions func(*VagrantCloudClient) error

func WithServerURL(url string) VagrantCloudClientOptions {
	return func(c *VagrantCloudClient) (err error) {
		c.url = url
		return
	}
}

func WithClientRetryCount(r int) VagrantCloudClientOptions {
	return func(c *VagrantCloudClient) (err error) {
		c.retryCount = r
		return
	}
}

func NewVagrantCloudClient(accessToken string, opts ...VagrantCloudClientOptions) (vcc *VagrantCloudClient, err error) {
	vcc = &VagrantCloudClient{
		accessToken: accessToken,
		retryCount:  DEFAULT_RETRY_COUNT,
		url:         DEFAULT_API_URL,
	}
	for _, opt := range opts {
		if oerr := opt(vcc); oerr != nil {
			err = multierror.Append(err, oerr)
		}
	}
	if err != nil {
		return nil, err
	}

	return
}

func (vc *VagrantCloudClient) request(
	path string, method HTTPMethod, params map[string]interface{},
) (jsonResp map[string]interface{}, err error) {
	var raw []byte
	var vcr *VagrantCloudRequest

	// Build url
	url := fmt.Sprintf("%s/%s", vc.url, path)

	// Request with query parameters if the HTTPMethod is GET, HEAD or DELETE
	queryParamMethods := []string{DELETE.String(), GET.String(), HEAD.String()}
	if contains(queryParamMethods, method.String()) {
		stringParams := make(map[string]string)
		for k, v := range params {
			stringParams[k] = v.(string)
		}
		vcr, err = NewVagrantCloudRequest(
			WithAuthTokenHeader(vc.accessToken),
			WithRetryCount(vc.retryCount),
			WithURL(url),
			ReplaceHosts(),
			WithMethod(method),
			WithURLQueryParams(stringParams),
		)
		if err != nil {
			return nil, err
		}
	} else {
		vcr, err = NewVagrantCloudRequest(
			WithAuthTokenHeader(vc.accessToken),
			WithRetryCount(vc.retryCount),
			WithURL(url),
			ReplaceHosts(),
			WithMethod(method),
			WithRequestJSONableData(params),
		)
		if err != nil {
			return nil, err
		}
	}
	// Execute request against Vagrant Cloud
	raw, err = vcr.Do()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &jsonResp); err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func (vc *VagrantCloudClient) AuthedRequest(url string, vagrantServerUrl *url.URL, method HTTPMethod, ui terminal.UI) (data []byte, err error) {
	opts := []VagrantCloudRequestOptions{
		WithRetryCount(vc.retryCount),
		WithURL(url),
		ReplaceHosts(),
		WithMethod(method),
		WarnDifferentTarget(vagrantServerUrl, ui),
	}

	var accessTokenByUrl bool
	accessTokenEnvVar := os.Getenv("VAGRANT_SERVER_ACCESS_TOKEN_BY_URL")
	if accessTokenEnvVar == "" {
		accessTokenByUrl = false
	} else {
		accessTokenByUrl, err = strconv.ParseBool(accessTokenEnvVar)
		if err != nil {
			return nil, err
		}
	}
	if accessTokenByUrl {
		// TODO: warn user
		opts = append(opts, WithAuthTokenURLParam(vc.accessToken))
	} else {
		opts = append(opts, WithAuthTokenHeader(vc.accessToken))
	}

	vcr, err := NewVagrantCloudRequest(opts...)
	if err != nil {
		return nil, err
	}
	return vcr.Do()
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
