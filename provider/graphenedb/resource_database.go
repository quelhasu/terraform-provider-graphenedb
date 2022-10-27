package graphenedb

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)

const waitAfterCreateOrUpdate = 0
const waitForPlugins = 0

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"plugins": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
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
					},
				},
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	databaseId, err := meta.(*graphendbclient.RestApiClient).CreateDatabase(ctx, extractDatabaseInfoFromSchema(ctx, d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(databaseId)
	tflog.Debug(ctx, "CREATE DATABASE - DATABASE CREATED SUCCESSFULLY",  map[string]interface{}{
			"databaseId": databaseId,
	})
	time.Sleep(waitAfterCreateOrUpdate * time.Second)

	plugins := extractPluginInfoFromSchema(ctx, d)
	
	for _, plugin := range plugins {

		pluginCreateResult, err := meta.(*graphendbclient.RestApiClient).CreatePlugin(ctx, databaseId, plugin)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "CREATE DATABASE - PLUGIN ADDED SUCCESSFULLY",  map[string]interface{}{
			"Id": pluginCreateResult.Detail.Id,
			"Kind": pluginCreateResult.Detail.Kind,
			"Enabled": pluginCreateResult.Detail.Enabled,
			"Name": pluginCreateResult.Detail.Name,
		})
		time.Sleep(waitForPlugins * time.Second)
		err = meta.(*graphendbclient.RestApiClient).ChangePluginStatus(ctx, databaseId, pluginCreateResult.Detail.Id, graphendbclient.PluginEnabledStatus)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "CREATE DATABASE - PLUGIN ENABLED SUCCESSFULLY",  map[string]interface{}{
			"DatabaseId": databaseId,
			"PluginId": pluginCreateResult.Detail.Id,
			"Status":  graphendbclient.PluginEnabledStatus,
		})
	}
	if(len(plugins) > 0) {
		
		err = meta.(*graphendbclient.RestApiClient).RestartDatabase(ctx, databaseId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Cannot restart the database",
				Detail:   fmt.Sprintf("The database ID is %s", databaseId),
			})
		}
		tflog.Debug(ctx, "CREATE DATABASE - DATABASE RESTARTED SUCCESSFULLY",  map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}

	readDiags := resourceDatabaseRead(ctx, d, meta)
	if(len(readDiags) > 0){
		diags = append(diags, readDiags...)
	}
	return diags
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	databaseInfo, err := meta.(*graphendbclient.RestApiClient).GetUpstreamDatabaseInfo(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "READ DATABASE - DATABASE INFO",  map[string]interface{}{
		"DatabaseInfo": fmt.Sprintf("%+v", databaseInfo),
	})
	
	if(databaseInfo == nil){
		tflog.Debug(ctx, "READ DATABASE - DATABASE NOT FOUND",  map[string]interface{}{
			"DatabaseId":  d.Id(),
		})
		d.SetId("")
		return diags
	}

	flattenedData := flattenDatabase(ctx, databaseInfo, d)
	tflog.Debug(ctx, "READ DATABASE - FLATTENED DATABASE DATA", flattenedData)
	err = AttributesToResourceData(flattenedData, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	databaseId := d.Id()

	if d.HasChanges("plugins") {
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE PLUGIN CHANGED")
		upstreamDatabaseInfo, err := meta.(*graphendbclient.RestApiClient).GetUpstreamDatabaseInfo(ctx, databaseId)
		if err != nil {
			return diag.FromErr(err)
		}
		if(len(upstreamDatabaseInfo.Plugins) > 0){
			for _, plugin := range upstreamDatabaseInfo.Plugins {
				tflog.Debug(ctx, "UPDATE DATABASE - DATABASE PLUGIN IS DISABLING...", map[string]interface{}{
					"DatabaseId":  databaseId,
					"PluginId":  plugin.Id,
				})
				err = meta.(*graphendbclient.RestApiClient).ChangePluginStatus(ctx, databaseId, plugin.Id, graphendbclient.PluginDisabledStatus)
				if err != nil {
					return diag.FromErr(err)
				}
				time.Sleep(waitForPlugins * time.Second)
				tflog.Debug(ctx, "UPDATE DATABASE - DATABASE PLUGIN IS REMOVING...", map[string]interface{}{
					"DatabaseId":  databaseId,
					"PluginId":  plugin.Id,
				})
				err = meta.(*graphendbclient.RestApiClient).DeletePlugin(ctx, databaseId, plugin.Id)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		
		plugins := extractPluginInfoFromSchema(ctx, d)
		for _, plugin := range plugins {
			pluginCreateResult, err := meta.(*graphendbclient.RestApiClient).CreatePlugin(ctx, databaseId, plugin)
			if err != nil {
				return diag.FromErr(err)
			}
			tflog.Debug(ctx, "UPDATE DATABASE - PLUGIN ADDED SUCCESSFULLY",  map[string]interface{}{
				"Id": pluginCreateResult.Detail.Id,
				"Kind": pluginCreateResult.Detail.Kind,
				"Enabled": pluginCreateResult.Detail.Enabled,
				"Name": pluginCreateResult.Detail.Name,
			})
			time.Sleep(waitForPlugins * time.Second)
			err = meta.(*graphendbclient.RestApiClient).ChangePluginStatus(ctx, databaseId, pluginCreateResult.Detail.Id, graphendbclient.PluginEnabledStatus)
			if err != nil {
				return diag.FromErr(err)
			}
			tflog.Debug(ctx, "UPDATE DATABASE - PLUGIN ENABLED SUCCESSFULLY",  map[string]interface{}{
				"DatabaseId": databaseId,
				"PluginId": pluginCreateResult.Detail.Id,
				"Status":  graphendbclient.PluginEnabledStatus,
			})
		}
	}

	if d.HasChange("plan") {
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE MUST BE UPDATED",  map[string]interface{}{
			"DatabaseId": databaseId,
		})
		databaseId, err := meta.(*graphendbclient.RestApiClient).UpdateDatabase(ctx, databaseId, graphendbclient.DatabaseUpgradeInfo {
			Plan: d.Get("plan").(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(databaseId)
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE SUCCESSFULLY UPDATED",  map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}
	
	time.Sleep(waitAfterCreateOrUpdate * time.Second)

	if d.HasChanges("plugins") && !d.HasChange("plan") {
		err := meta.(*graphendbclient.RestApiClient).RestartDatabase(ctx, databaseId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Cannot restart the database",
				Detail:   fmt.Sprintf("The database ID is %s", databaseId),
			})
		}
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE RESTARTED SUCCESSFULLY",  map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}

	readDiags := resourceDatabaseRead(ctx, d, meta)
	if(len(readDiags) > 0){
		diags = append(diags, readDiags...)
	}
	return diags
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return diag.Errorf("Database delete is called")
}




func flattenDatabase(ctx context.Context, database *graphendbclient.UpstreamDatabaseInfo, d *schema.ResourceData ) map[string]interface{} {
	c := make(map[string]interface{})
  // c["id"] = database.Id
	c["plan"]= database.Plan.PlanType
	c["region"]= database.AwsRegion
	c["vpc_id"]= database.PrivateNetworkId
	c["version"]= database.VersionKey
	c["name"] = database.Name
	c["plugins"] = flattenPlugins(ctx , database, d)
	//tflog.Warn(ctx, "UPSTREAM DATABASE STATE",  c )
  return c
}

func flattenPlugins(ctx context.Context, database *graphendbclient.UpstreamDatabaseInfo, d *schema.ResourceData) []interface{} {
	countOfPlugins := len(database.Plugins)
	if(countOfPlugins >= 0){
		allPlugins := extractPluginInfoFromSchema(ctx, d)
		allPluginsMap := make([]interface{}, countOfPlugins)
		for i, plugin := range database.Plugins {
			pluginUrl := ""
			pluginKind := ""
			currentPluginFromSchema := getPluginByName(ctx, allPlugins, plugin.Name)
			if(currentPluginFromSchema != nil){
				pluginUrl = currentPluginFromSchema.Url
			}
			switch plugin.Type {
				case "extensions":
					pluginKind = "extension"
				case "storedprocedure":
					pluginKind = "stored-procedure"
			}
			currentPluginMap := make(map[string]interface{})
			currentPluginMap["plugin_id"] = plugin.Id
			currentPluginMap["name"] = plugin.Name
			currentPluginMap["kind"] = pluginKind
			currentPluginMap["url"] = pluginUrl
			allPluginsMap[i] = currentPluginMap
		}
		return allPluginsMap
	}
	return make([]interface{}, 0)
}

func extractDatabaseInfoFromSchema(ctx context.Context, d *schema.ResourceData) graphendbclient.DatabaseInfo {
	return graphendbclient.DatabaseInfo {
		Name: d.Get("name").(string),
		Version: d.Get("version").(string),
		AwsRegion: d.Get("region").(string),
		Plan: d.Get("plan").(string),
		Vpc: d.Get("vpc_id").(string),
	}
}

func extractPluginInfoFromSchema(ctx context.Context, d *schema.ResourceData) []graphendbclient.PluginInfo {
	if d.Get("plugins") != nil {
	inputPlugins := d.Get("plugins").([]interface{})
	databasePlugins := []graphendbclient.PluginInfo{}
		for _, inputPlugin := range inputPlugins {
			current := inputPlugin.(map[string]interface{})
			oi := graphendbclient.PluginInfo{
				Name: current["name"].(string),
				Kind: current["kind"].(string),
				Url:  current["url"].(string),
			}
			databasePlugins = append(databasePlugins, oi)
		}
		return databasePlugins
	}
	return nil
}

func getPluginByName(ctx context.Context, allPlugins []graphendbclient.PluginInfo, pluginName string) *graphendbclient.PluginInfo {
	if(len(allPlugins) > 0){
		for _, plugin := range allPlugins {
			if(plugin.Name == pluginName ){
				return &plugin
			}
		}
	}
	return nil
}


