package client

import (
	"net/http"
	"net/url"
	"time"
)

type credentials struct {
	userName string
	password string
}

type Client struct {
	credentials
	apiClient
}

func NewClient(userName, password string, apiEndpoint *url.URL) *Client {
	return &Client{
		credentials: credentials{
			userName: userName,
			password: password,
		},
		apiClient: apiClient{
			apiEndpoint: apiEndpoint,
			httpClient: &http.Client{
				Transport: &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					TLSHandshakeTimeout: 120 * time.Second},
			},
		},
	}
}

func (c *Client) unqualify(names ...*string) {
	for _, name := range names {
		*name = *name
	}
}
