package pureapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Client communicates with the Pure Storage FlashArray REST API 2.x.
type Client struct {
	baseURL      string
	apiToken     string
	sessionToken string
	apiVersion   string
	httpClient   *http.Client
}

// NewClient creates a new FlashArray API client.
// endpoint is the array IP or FQDN, apiToken is the API token for auth.
// If insecure is true, TLS certificate verification is skipped.
// If apiVersion is empty, the highest available 2.x version is negotiated.
func NewClient(endpoint, apiToken, apiVersion string, insecure bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
	return &Client{
		baseURL:    fmt.Sprintf("https://%s", strings.TrimRight(endpoint, "/")),
		apiToken:   apiToken,
		apiVersion: apiVersion,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// NegotiateAPIVersion queries /api/api_version and picks the highest 2.x version.
func (c *Client) NegotiateAPIVersion() error {
	if c.apiVersion != "" {
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/api_version", nil)
	if err != nil {
		return fmt.Errorf("create version request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("version negotiation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("version negotiation: HTTP %d", resp.StatusCode)
	}

	var vr APIVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		return fmt.Errorf("parse version response: %w", err)
	}

	var v2versions []string
	for _, v := range vr.Version {
		if strings.HasPrefix(v, "2.") {
			v2versions = append(v2versions, v)
		}
	}
	if len(v2versions) == 0 {
		return fmt.Errorf("no REST API 2.x version available on array")
	}

	sort.Strings(v2versions)
	c.apiVersion = v2versions[len(v2versions)-1]
	return nil
}

// Login authenticates via API token and stores the session token.
func (c *Client) Login() error {
	if err := c.NegotiateAPIVersion(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/%s/login", c.baseURL, c.apiVersion)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}
	req.Header.Set("api-token", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: HTTP %d: %s", resp.StatusCode, string(body))
	}

	token := resp.Header.Get("x-auth-token")
	if token == "" {
		return fmt.Errorf("login response missing x-auth-token header")
	}

	c.sessionToken = token
	return nil
}

// Logout invalidates the current session.
func (c *Client) Logout() {
	if c.sessionToken == "" {
		return
	}
	url := fmt.Sprintf("%s/api/%s/logout", c.baseURL, c.apiVersion)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("x-auth-token", c.sessionToken)
	resp, err := c.httpClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	c.sessionToken = ""
}

// doGet performs an authenticated GET request and decodes the JSON response.
func (c *Client) doGet(path string, result interface{}) error {
	url := fmt.Sprintf("%s/api/%s%s", c.baseURL, c.apiVersion, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-auth-token", c.sessionToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

// GetArraySpace returns capacity/space data for the array.
func (c *Client) GetArraySpace() (*ArraySpaceResponse, error) {
	var resp ArraySpaceResponse
	if err := c.doGet("/arrays/space", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetArrayPerformance returns performance metrics for the array.
func (c *Client) GetArrayPerformance() (*ArrayPerformanceResponse, error) {
	var resp ArrayPerformanceResponse
	if err := c.doGet("/arrays/performance", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHardware returns hardware component status for the array.
func (c *Client) GetHardware() (*HardwareResponse, error) {
	var resp HardwareResponse
	if err := c.doGet("/hardware", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVolumesSpace returns space data for all volumes.
// If names is non-empty, filters to only those volumes.
func (c *Client) GetVolumesSpace(names []string) (*VolumeSpaceResponse, error) {
	path := "/volumes/space"
	if len(names) > 0 {
		path += "?names=" + strings.Join(names, ",")
	}
	var resp VolumeSpaceResponse
	if err := c.doGet(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ArrayName returns the array name from a space response for use in messages.
func (c *Client) ArrayName() (string, error) {
	var resp ArraySpaceResponse
	if err := c.doGet("/arrays/space", &resp); err != nil {
		return "", err
	}
	if len(resp.Items) == 0 {
		return "unknown", nil
	}
	return resp.Items[0].Name, nil
}
