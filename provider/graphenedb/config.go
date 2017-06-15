package graphenedb

import (
	"fmt"
	"net/url"

	cli "github.com/ervinjohnson/terraform-provider-graphenedb/client"
)

type GrapheneDBClient struct {
	*cli.AuthenticatedClient
}

type Config struct {
	User     string
	Password string
	Endpoint string
}

func (c *Config) Client() (*GrapheneDBClient, error) {
	uri, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid endpoint URI: %s", err)
	}

	client := cli.NewClient(c.User, c.Password, uri)
	authenticatedClient, err := client.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("Authentication failed: %s", err)
	}

	grapheneDBClient := &GrapheneDBClient{
		AuthenticatedClient: authenticatedClient,
	}

	return grapheneDBClient, nil
}
