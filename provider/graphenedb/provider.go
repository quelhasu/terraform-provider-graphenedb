package graphenedb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)

func Provider() *schema.Provider {
	// cli.ClientLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"graphenedb_vpc_peering": resourceVpcPeering(),
			"graphenedb_database":    resourceDatabase(),
			"graphenedb_environment": resourceEnvironment(),
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_CLIENT_ID", nil),
				Description: "The Client Id for GrapeheneDB API operations.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_CLIENT_SECRET", nil),
				Description: "The Client Secret for GrapheneDB API operations.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_ENDPOINT", nil),
				Description: "The Environment ID for GrapheneDB API operations.",
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	environmentId := d.Get("environment_id").(string)
	return graphendbclient.NewApiClient(environmentId, clientId, clientSecret)
}
