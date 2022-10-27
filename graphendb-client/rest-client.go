package graphendbclient

import resty "github.com/go-resty/resty/v2"

type RestApiClient struct {
	BaseUrl string
	ApiKey  string
	ApiClient *resty.Client
}

func NewApiClient(endpoint string, apiKey string) *RestApiClient {
	
	client := resty.New()
	client.SetHeader("Api_key", apiKey)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "*/*")
	client.SetBaseURL(endpoint)
	
	return &RestApiClient{
		BaseUrl: endpoint,
		ApiKey:  apiKey,
		ApiClient: client,
	}
}
