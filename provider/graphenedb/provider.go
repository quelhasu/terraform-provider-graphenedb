package graphenedb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)

func Provider() *schema.Provider {
	// cli.ClientLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"graphenedb_vpc":      resourceVpc(),
			"graphenedb_database": resourceDatabase(),
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
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_ENDPOINT", nil),
				Description: "The HTTP endpoint for GrapheneDB API operations.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_ENV_ID", nil),
				Description: "The Environement ID of the DB for GrapheneDB API operations.",
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	environementId := d.Get("environement_id").(string)
	endpoint := d.Get("endpoint").(string)
	return graphendbclient.NewApiClient(endpoint, clientId, clientSecret, environementId)
}
