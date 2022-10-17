package graphenedb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCreate,
		ReadContext:   resourcePluginRead,
		UpdateContext: resourcePluginUpdate,
		DeleteContext: resourcePluginDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"database_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())

	var diags diag.Diagnostics

	name, kind, url, database_id := d.Get("name").(string), d.Get("kind").(string), d.Get("url").(string), d.Get("database_id").(string)

	log.Printf("Adding plugin with : name: %s, kind: %s, url: %s, database id: %s", name, kind, url, database_id)
	client := m.(*GrapheneDBClient).NewDatabasesClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	plugin, err := client.AddPlugin(name, kind, url, database_id)
	if err != nil {
		return diag.Errorf("Error creating database %s: %s", name, err)
	}

	plugin_id := plugin.Detail.Id
	log.Printf("Enable plugin: %s", plugin_id)

	plugin_status, err_en := client.ChangePluginStatus(database_id, plugin_id, "enabled")

	if err_en != nil {
		return diag.Errorf("Error enabling plugin %s(%s) for database %s: %s", plugin_id, plugin_status, database_id, err_en)
	}

	d.SetId(plugin_id)
	
	return diags	
}

func resourcePluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourcePluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourcePluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}
