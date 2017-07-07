package graphenedb

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	// cli.ClientLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"graphenedb_database": resourceDatabase(),
		},
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_USER", nil),
				Description: "The user name for GrapheneDB API operations.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAPHENEDB_PASSWORD", nil),
				Description: "The user password for GrapheneDB API operations.",
			},

			"endpoint": &schema.Schema{
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
	config := Config{
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
		Endpoint: d.Get("endpoint").(string),
	}

	return config.Client()
}
