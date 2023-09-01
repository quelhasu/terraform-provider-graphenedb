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
			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "enterprise" && v != "community" {
						errs = append(errs, fmt.Errorf("%q must be of type 'enterprise' or 'community', got: %s", key, v))
					}
					return
				},
			},
			"enabled_extras_kinds": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "graphneo" && v != "ongdb" {
						errs = append(errs, fmt.Errorf("%q must be of type 'graphneo' or 'ongdb', got: %s", key, v))
					}
					return
				},
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"private_domain_name": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"http_port": {
				Type:     schema.TypeInt,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"bolt_port": {
				Type:     schema.TypeInt,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"plugins": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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

	databaseId, err := meta.(*graphendbclient.RestApiClient).CreateDatabase(ctx, extractDatabaseInfoFromSchema(ctx, d), d.Get("vendor").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(databaseId)
	tflog.Debug(ctx, "CREATE DATABASE - DATABASE CREATED SUCCESSFULLY", map[string]interface{}{
		"databaseId": databaseId,
	})
	time.Sleep(waitAfterCreateOrUpdate * time.Second)

	plugins := extractPluginInfoFromSchema(ctx, d)

	for _, plugin := range plugins {

		pluginCreateResult, err := meta.(*graphendbclient.RestApiClient).CreatePlugin(ctx, databaseId, plugin, d.Get("vendor").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "CREATE DATABASE - PLUGIN ADDED SUCCESSFULLY", map[string]interface{}{
			"Id":        pluginCreateResult.Id,
			"Name":      pluginCreateResult.Name,
			"CreatedAt": pluginCreateResult.CreatedAt,
		})
		time.Sleep(waitForPlugins * time.Second)
	}

	if len(plugins) > 0 {
		err = meta.(*graphendbclient.RestApiClient).RestartDatabase(ctx, databaseId, d.Get("vendor").(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Cannot restart the database",
				Detail:   fmt.Sprintf("The database ID is %s", databaseId),
			})
		}
		tflog.Debug(ctx, "CREATE DATABASE - DATABASE RESTARTED SUCCESSFULLY", map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}

	readDiags := resourceDatabaseRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}
	return diags
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	databaseInfo, err := meta.(*graphendbclient.RestApiClient).GetUpstreamDatabaseInfo(ctx, d.Id(), d.Get("vendor").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "READ DATABASE - DATABASE INFO", map[string]interface{}{
		"DatabaseInfo": fmt.Sprintf("%+v", databaseInfo),
	})
	d.Set("domain_name", databaseInfo.DomainName)
	d.Set("private_domain_name", databaseInfo.PrivateDomainName)
	d.Set("http_port", databaseInfo.HTTPPort)
	d.Set("bolt_port", databaseInfo.BoltPort)

	if databaseInfo == nil {
		tflog.Debug(ctx, "READ DATABASE - DATABASE NOT FOUND", map[string]interface{}{
			"DatabaseId": d.Id(),
		})
		d.SetId("")
		return diags
	}

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	databaseId := d.Id()

	if d.HasChanges("plugins") {
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE PLUGIN CHANGED")
		DBPluginsInfo, err := meta.(*graphendbclient.RestApiClient).GetUpstreamDatabasePluginsInfo(ctx, databaseId, d.Get("vendor").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if len(DBPluginsInfo.Plugins) > 0 {
			for _, plugin := range DBPluginsInfo.Plugins {
				tflog.Debug(ctx, "UPDATE DATABASE - DATABASE PLUGIN IS REMOVING...", map[string]interface{}{
					"DatabaseId": databaseId,
					"PluginId":   plugin.ID,
				})
				err = meta.(*graphendbclient.RestApiClient).DeletePlugin(ctx, databaseId, d.Get("vendor").(string), plugin.ID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		plugins := extractPluginInfoFromSchema(ctx, d)
		for _, plugin := range plugins {
			pluginCreateResult, err := meta.(*graphendbclient.RestApiClient).CreatePlugin(ctx, databaseId, plugin, d.Get("vendor").(string))
			if err != nil {
				return diag.FromErr(err)
			}
			tflog.Debug(ctx, "UPDATE DATABASE - PLUGIN ADDED SUCCESSFULLY", map[string]interface{}{
				"Id":        pluginCreateResult.Id,
				"Name":      pluginCreateResult.Name,
				"CreatedAt": pluginCreateResult.CreatedAt,
			})
			time.Sleep(waitForPlugins * time.Second)
		}
	}

	if d.HasChange("plan") {
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE MUST BE UPDATED", map[string]interface{}{
			"DatabaseId": databaseId,
		})
		databaseId, err := meta.(*graphendbclient.RestApiClient).UpdateDatabase(ctx, databaseId, graphendbclient.DatabaseUpgradeInfo{
			Plan: d.Get("plan").(string),
		}, d.Get("vendor").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(databaseId)
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE SUCCESSFULLY UPDATED", map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}

	time.Sleep(waitAfterCreateOrUpdate * time.Second)

	if d.HasChanges("plugins") && !d.HasChange("plan") {
		err := meta.(*graphendbclient.RestApiClient).RestartDatabase(ctx, databaseId, d.Get("vendor").(string))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Cannot restart the database",
				Detail:   fmt.Sprintf("The database ID is %s", databaseId),
			})
		}
		tflog.Debug(ctx, "UPDATE DATABASE - DATABASE RESTARTED SUCCESSFULLY", map[string]interface{}{
			"DatabaseId": databaseId,
		})
	}

	readDiags := resourceDatabaseRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}
	return diags
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	var diags diag.Diagnostics
	err := meta.(*graphendbclient.RestApiClient).DeleteDatabase(ctx, d.Id(), d.Get("vendor").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func extractDatabaseInfoFromSchema(ctx context.Context, d *schema.ResourceData) graphendbclient.DatabaseInfo {
	return graphendbclient.DatabaseInfo{
		Name:               d.Get("name").(string),
		Plan:               d.Get("plan").(string),
		Version:            d.Get("version").(string),
		Edition:            d.Get("edition").(string),
		EnabledExtrasKinds: d.Get("enabled_extras_kinds").([]interface{}),
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
				Url:  current["url"].(string),
			}
			databasePlugins = append(databasePlugins, oi)
		}
		return databasePlugins
	}
	return nil
}

func getPluginByName(ctx context.Context, allPlugins []graphendbclient.PluginInfo, pluginName string) *graphendbclient.PluginInfo {
	if len(allPlugins) > 0 {
		for _, plugin := range allPlugins {
			if plugin.Name == pluginName {
				return &plugin
			}
		}
	}
	return nil
}
