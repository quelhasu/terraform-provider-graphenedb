package graphendbclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	resty "github.com/go-resty/resty/v2"
)

type RestApiClient struct {
	BaseUrl        string
	ClientId       string
	ClientSecret   string
	EnvironementId string
	BearerToken    string
	ApiClient      *resty.Client
}

func NewApiClient(endpoint string, client_id string, client_secret string, environement_id string) (*RestApiClient, error) {

	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "*/*")
	client.SetBaseURL(endpoint)

	apiClient := &RestApiClient{
		BaseUrl:        endpoint,
		EnvironementId: environement_id,
		ClientId:       client_id,
		ClientSecret:   client_secret,
		BearerToken:    "",
		ApiClient:      client,
	}

	if err := apiClient.GetBearerToken(); err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}

	var token strings.Builder
	token.WriteString("Bearer ")
	token.WriteString(apiClient.BearerToken)

	client.SetHeader("Authorization", token.String())

	return apiClient, nil
}

// Update the TokenResponse struct according the new API
func (c *RestApiClient) GetBearerToken() error {
	response, err := c.ApiClient.R().SetFormData(map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.ClientSecret,
		"grant_type":    "client_credentials",
	}).Post("/organizations/oauth/token")
	if err != nil {
		return fmt.Errorf("failed to retrieve OAuth2 token: %v", err)
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to retrieve OAuth2 token: %v", response.Status())
	}

	type TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	tokenResponse := &TokenResponse{}
	if err := json.Unmarshal(response.Body(), tokenResponse); err != nil {
		return fmt.Errorf("failed to unmarshal OAuth2 token response: %v", err)
	}

	c.BearerToken = tokenResponse.AccessToken

	return nil
}
