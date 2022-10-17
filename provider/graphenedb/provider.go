package graphenedb

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cli "github.com/quelhasu/terraform-provider-graphenedb/client"
)

type GrapheneDBClient struct {
	*cli.AuthenticatedClient
}

type Config struct {
	User     string
	Password string
	Endpoint string
}

func Provider() *schema.Provider {
	// cli.ClientLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"graphenedb_database": resourceDatabase(),
			"graphenedb_database_restart": resourceDatabaseRestart(),
			"graphenedb_vpc": resourceVpc(),
			"graphenedb_plugin": resourcePlugin(),
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_USER", nil),
				Description: "The user name for GrapheneDB API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_PASSWORD", nil),
				Description: "The user password for GrapheneDB API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_ENDPOINT", nil),
				Description: "The HTTP endpoint for GrapheneDB API operations.",
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func (c *Config) Client() (*GrapheneDBClient, error) {
	uri, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Invalid endpoint URI", err)
	}

	client := cli.NewClient(c.User, c.Password, uri)
	authenticatedClient, err := client.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Authentication failed", err)
	}

	grapheneDBClient := &GrapheneDBClient{
		AuthenticatedClient: authenticatedClient,
	}

	return grapheneDBClient, err
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	user := d.Get("user").(string)
	password := d.Get("password").(string)
	endpoint := d.Get("endpoint").(string)
	// Warning or errors can be collected in a slice type

	config := &Config{
		User:     user,
		Password: password,
		Endpoint: endpoint,
	}

	return config.Client()
}
