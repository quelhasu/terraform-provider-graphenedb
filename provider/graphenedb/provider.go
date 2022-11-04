package graphenedb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)


type Config struct {
	ApiKey string
	Endpoint string
}

func Provider() *schema.Provider {
	// cli.ClientLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"graphenedb_vpc": resourceVpc(),
			"graphenedb_database": resourceDatabase(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_API_KEY", nil),
				Description: "The API key for GrapheneDB API operations.",
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

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiKey := d.Get("api_key").(string)
	endpoint := d.Get("endpoint").(string)
	return graphendbclient.NewApiClient(endpoint, apiKey), nil
}
