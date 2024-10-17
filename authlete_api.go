package gauthlete

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	contentTypeJSON = "application/json;charset=UTF-8"
	acceptJSON      = "application/json"
)

type AuthleteError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e AuthleteError) Error() string {
	return e.Message
}

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
	u := c.baseURL + "/api/auth/introspection"

	payload := map[string]string{"token": token}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Accept", acceptJSON)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var authleteErr AuthleteError
		if err := json.NewDecoder(resp.Body).Decode(&authleteErr); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("Authlete error: %s (code: %d)", authleteErr.Message, authleteErr.Code)
	}

	var introspectionResp IntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspectionResp); err != nil {
		return nil, err
	}

	return &introspectionResp, nil
}

// IntrospectionResponse represents the response from token introspection
type IntrospectionResponse struct {
	Active bool `json:"active"`
	// TODO: Add other fields as per Authlete API response
}

// TODO: Add other Authlete API methods...
// Authorization Service API

func (c *ServiceClient) Authorization(parameters string) (*AuthorizationResponse, error) {
	url := c.baseURL + "/api/auth/authorization"
	payload := map[string]string{"parameters": parameters}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Accept", acceptJSON)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var authleteErr AuthleteError
		if err := json.NewDecoder(resp.Body).Decode(&authleteErr); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, authleteErr
	}

	var authResp AuthorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &authResp, nil
}

type AuthorizationResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`
	ResultMessage string `json:"resultMessage,omitempty"`

	Action string `json:"action,omitempty"`
	Ticket string `json:"ticket,omitempty"`
}

func (c *ServiceClient) AuthorizationFail(ticket string) (*AuthorizationFailResponse, error) {
	u := c.baseURL + "/api/auth/authorization/fail"
	payload := map[string]string{"ticket": ticket}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	responseContainer := &AuthorizationFailResponse{}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Accept", acceptJSON)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var authleteErr AuthleteError
		if err := json.NewDecoder(res.Body).Decode(&authleteErr); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, authleteErr
	}

	// Convert the JSON (body) to an instance.
	if err := json.NewDecoder(res.Body).Decode(responseContainer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return responseContainer, nil
}

func (c *ServiceClient) AuthorizationIssue(ticket, subject string) (*AuthorizationIssueResponse, error) {
	u := c.baseURL + "/api/auth/authorization/issue"
	payload := map[string]string{"ticket": ticket, "subject": subject}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	responseContainer := &AuthorizationIssueResponse{}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Accept", acceptJSON)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var authleteErr AuthleteError
		if err := json.NewDecoder(res.Body).Decode(&authleteErr); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}
		return nil, authleteErr
	}

	// Convert the JSON (body) to an instance.
	if err := json.NewDecoder(res.Body).Decode(responseContainer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return responseContainer, nil
}

type AuthorizationFailResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`
	ResultMessage string `json:"resultMessage,omitempty"`

	Action          string `json:"action,omitempty"`
	ResponseContent string `json:"responseContent,omitempty"`
}

type AuthorizationIssueResponse struct {
	ResultCode    string `json:"resultCode,omitempty"`
	ResultMessage string `json:"resultMessage,omitempty"`

	Action               string `json:"action,omitempty"`
	ResponseContent      string `json:"responseContent,omitempty"`
	AccessToken          string `json:"accessToken"`
	AccessTokenExpiresAt int64  `json:"accessTokenExpiresAt"`
	AccessTokenDuration  int64  `json:"accessTokenDuration"`
	IdToken              string `json:"idToken"`
	AuthorizationCode    string `json:"authorizationCode"`
	JwtAccessToken       string `json:"jwtAccessToken"`
	// there are more in the docs
}

// TokenResponse represents the response from Authlete's /api/auth/token API
type TokenResponse struct {
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`

	Action          string `json:"action"`
	ResponseContent string `json:"responseContent"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       int64  `json:"expiresIn"`
	RefreshToken    string `json:"refreshToken"`
}

// Token sends a token request to Authlete's /api/auth/token API
func (c *ServiceClient) Token(parameters, clientID, clientSecret string) (*TokenResponse, error) {
	url := c.baseURL + "/api/auth/token"

	payload := map[string]string{
		"parameters": parameters,
		"clientId":   clientID,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var authleteErr AuthleteError
		if err := json.NewDecoder(resp.Body).Decode(&authleteErr); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return nil, authleteErr
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}
