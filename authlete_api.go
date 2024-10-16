package gauthlete

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// ServiceClient represents the API client
type ServiceClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	apiSecret  string
}

// NewServiceClient creates a new API client
func NewServiceClient() *ServiceClient {
	baseURL, ok := os.LookupEnv("AUTHLETE_BASE_URL")
	if !ok {
		baseURL = "https://api.authlete.com"
	}

	return &ServiceClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:   baseURL,
		apiKey:    os.Getenv("AUTHLETE_SERVICE_APIKEY"),
		apiSecret: os.Getenv("AUTHLETE_SERVICE_APISECRET"),
	}
}

// TokenIntrospect introspects a token
func (c *ServiceClient) TokenIntrospect(token string) (*IntrospectionResponse, error) {
	payload := map[string]string{"token": token}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/auth/introspection", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var introspectionResp IntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspectionResp); err != nil {
		return nil, err
	}

	return &introspectionResp, nil
}

// IntrospectionResponse represents the response from token introspection
type IntrospectionResponse struct {
	Active bool `json:"active"`
	// Add other fields as per Authlete's API response
}

// TODO: Add other Authlete API methods...
// Authorization Service API
