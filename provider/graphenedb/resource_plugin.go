package graphenedb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/quelhasu/terraform-provider-graphenedb/client"
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

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	
	name, kind, url, databaseId := d.Get("name").(string), d.Get("kind").(string), d.Get("url").(string), d.Get("database_id").(string)
	client := meta.(*GrapheneDBClient).NewDatabasesClient()
	pluginId, err := createPlugin(ctx, client, name, kind, url, databaseId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pluginId)
	err = RestartDatabase(databaseId, meta.(*GrapheneDBClient))
	if(err != nil){
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to restart database",
			Detail:   fmt.Sprintf("Error details is %#v", err),
		})
		return diags
	}
	return nil	
}

func resourcePluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourcePluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	var diags diag.Diagnostics
	pluginId := d.Id()
	name, kind, url, databaseId := d.Get("name").(string), d.Get("kind").(string), d.Get("url").(string), d.Get("database_id").(string)
	tflog.Warn(ctx, "NEW PLUGIN INFO", map[string]interface{}{
		"name": name,
		"kind": kind,
		"url": url,
		"databaseId": databaseId,
	})

	client := meta.(*GrapheneDBClient).NewDatabasesClient()
	err := removePlugin(ctx, client, databaseId, pluginId)
	if err != nil {
		return diag.FromErr(err)
	}
	newPluginId, err := createPlugin(ctx, client, name, kind, url, databaseId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newPluginId)
	_ = RestartDatabase(databaseId, meta.(*GrapheneDBClient))
	err = RestartDatabase(databaseId, meta.(*GrapheneDBClient))
	if(err != nil){
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to restart database",
			Detail:   fmt.Sprintf("Error details is %#v", err),
		})
		return diags
	}
	return nil
}

func resourcePluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	databaseId := d.Get("database_id").(string)
	pluginId := d.Id()
	client := meta.(*GrapheneDBClient).NewDatabasesClient()
	err := removePlugin(ctx, client, databaseId, pluginId)
	if(err != nil){
		return diag.FromErr(err)
	}
	_ = RestartDatabase(databaseId, meta.(*GrapheneDBClient))
	if(err != nil){
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to restart database",
			Detail:   fmt.Sprintf("Error details is %#v", err),
		})
		return diags
	}
	return nil
}

func removePlugin(ctx context.Context, client *client.DatabasesClient, databaseId string, pluginId string) error {
	_, err := client.ChangePluginStatus(databaseId, pluginId, "disabled")
	if(err != nil){
		return err
	}

	err = client.DeletePlugin(ctx, databaseId, pluginId)
	if(err != nil){
		return err
	}
	return nil
}

func createPlugin(ctx context.Context, client *client.DatabasesClient, name, kind, url, databaseId string) (string, error) {
	
	plugin, err := client.AddPlugin(name, kind, url, databaseId)
	if err != nil {
		return "", err
	}
	pluginId := plugin.Detail.Id
	_, err = client.ChangePluginStatus(databaseId, pluginId, "enabled")

	if err != nil {
		return "", err
	}
	return pluginId, nil

}

