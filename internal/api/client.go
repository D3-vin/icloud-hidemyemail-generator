package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/yourusername/hidemyemail-generator/pkg/models"
)

// Client handles iCloud Hide My Email API requests with TLS fingerprinting
type Client struct {
	tlsClient tls_client.HttpClient
	baseURLV1 string
	baseURLV2 string
	cookies   string
	params    models.APIParams
}

// NewClient creates a new API client with Chrome 146 TLS fingerprinting
func NewClient(cookies string) (*Client, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(10),
		tls_client.WithClientProfile(profiles.Chrome_146),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}

	tlsClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize TLS client: %w", err)
	}

	client := &Client{
		tlsClient: tlsClient,
		baseURLV1: "https://p68-maildomainws.icloud.com/v1/hme",
		baseURLV2: "https://p68-maildomainws.icloud.com/v2/hme",
		cookies:   sanitizeCookies(cookies),
	}

	if err := client.extractParams(); err != nil {
		return nil, fmt.Errorf("failed to extract API params: %w", err)
	}

	return client, nil
}

// sanitizeCookies removes CR/LF characters to prevent header injection
func sanitizeCookies(cookies string) string {
	cookies = strings.ReplaceAll(cookies, "\r", "")
	cookies = strings.ReplaceAll(cookies, "\n", "")
	return strings.TrimSpace(cookies)
}

// extractParams parses cookies to extract DSID and clientId
func (c *Client) extractParams() error {
	cookieMap := make(map[string]string)
	for _, part := range strings.Split(c.cookies, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			cookieMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	dsid := cookieMap["X-APPLE-WEBAUTH-USER"]
	if dsid == "" {
		dsid = cookieMap["X-APPLE-DS-WEB-SESSION-TOKEN"]
		if dsid == "" {
			return fmt.Errorf("failed to extract DSID from cookies")
		}
	}

	clientID := cookieMap["X-APPLE-WEBAUTH-TOKEN"]
	if clientID == "" {
		decoded, _ := url.QueryUnescape(dsid)
		if len(decoded) > 32 {
			clientID = decoded[:32]
		} else {
			clientID = decoded
		}
	}

	c.params = models.APIParams{
		ClientBuildNumber:     "2536Project32",
		ClientMasteringNumber: "2536B20",
		ClientID:              clientID,
		DSID:                  dsid,
	}

	return nil
}

func (c *Client) buildQueryParams() string {
	params := url.Values{}
	params.Add("clientBuildNumber", c.params.ClientBuildNumber)
	params.Add("clientMasteringNumber", c.params.ClientMasteringNumber)
	params.Add("clientId", c.params.ClientID)
	params.Add("dsid", c.params.DSID)
	return params.Encode()
}

// doRequest executes an HTTP request with standardized headers and error handling
func (c *Client) doRequest(method, endpoint, baseURL string, body interface{}) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, c.buildQueryParams())

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Validate sanitized cookies
	if strings.ContainsAny(c.cookies, "\r\n") {
		return nil, fmt.Errorf("Invalid cookies: contains illegal newline characters")
	}

	// Set headers matching Chrome browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", "https://www.icloud.com")
	req.Header.Set("Referer", "https://www.icloud.com/")
	req.Header.Set("Cookie", c.cookies)
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	resp, err := c.tlsClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			return nil, fmt.Errorf("Request timeout: API did not respond within 10 seconds")
		}
		return nil, fmt.Errorf("%s", formatNetworkError(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("Authentication failed: invalid or expired cookies. Please refresh your cookies using the cookie extraction script.")
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("Server error: iCloud API returned status %d. Please try again later.", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response: %w", err)
	}

	return respBody, nil
}

// GenerateEmail generates a new Hide My Email address
func (c *Client) GenerateEmail() (*models.GenerateResponse, error) {
	body := map[string]string{"langCode": "en-us"}
	respBody, err := c.doRequest(http.MethodPost, "/generate", c.baseURLV1, body)
	
	var genResp models.GenerateResponse
	if err != nil {
		genResp.Success = false
		genResp.Error = err.Error()
		return &genResp, nil
	}

	if err := json.Unmarshal(respBody, &genResp); err != nil {
		return &models.GenerateResponse{
			Success: false,
			Error:   "Failed to parse API response: invalid JSON format",
		}, nil
	}

	if !genResp.Success {
		genResp.Error = formatAPIError(genResp.Error)
	}

	return &genResp, nil
}

// ReserveEmail reserves and activates a generated email address
func (c *Client) ReserveEmail(hme, label, note string) (*models.ReserveResponse, error) {
	body := map[string]string{"hme": hme, "label": label, "note": note}
	respBody, err := c.doRequest(http.MethodPost, "/reserve", c.baseURLV1, body)
	
	var reserveResp models.ReserveResponse
	if err != nil {
		reserveResp.Success = false
		reserveResp.Error = err.Error()
		return &reserveResp, nil
	}

	if err := json.Unmarshal(respBody, &reserveResp); err != nil {
		return &models.ReserveResponse{
			Success: false,
			Error:   "Failed to parse API response: invalid JSON format",
		}, nil
	}

	if !reserveResp.Success {
		reserveResp.Error = formatAPIError(reserveResp.Error)
	}

	return &reserveResp, nil
}

// ListEmails fetches all Hide My Email addresses
func (c *Client) ListEmails() (*models.ListResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, "/list", c.baseURLV2, nil)
	
	var listResp models.ListResponse
	if err != nil {
		listResp.Success = false
		listResp.Error = err.Error()
		return &listResp, nil
	}

	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return &models.ListResponse{
			Success: false,
			Error:   "Failed to parse API response: invalid JSON format",
		}, nil
	}

	if !listResp.Success {
		listResp.Error = formatAPIError(listResp.Error)
	}

	return &listResp, nil
}

// formatAPIError extracts error messages from API responses
func formatAPIError(errorField interface{}) string {
	if errorField == nil {
		return "API request failed"
	}

	if errorMap, ok := errorField.(map[string]interface{}); ok {
		if errorMessage, ok := errorMap["errorMessage"].(string); ok && errorMessage != "" {
			return errorMessage
		}
		if reason, ok := errorMap["reason"].(string); ok && reason != "" {
			return reason
		}
	}

	if errorString, ok := errorField.(string); ok && errorString != "" {
		return errorString
	}

	return "API request failed"
}

// formatNetworkError formats network errors
func formatNetworkError(err error) string {
	if err == nil {
		return "Unknown network error"
	}
	return fmt.Sprintf("Network error: %s", err.Error())
}
