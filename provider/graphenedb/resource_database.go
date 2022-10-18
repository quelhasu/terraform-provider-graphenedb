package graphenedb

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/quelhasu/terraform-provider-graphenedb/client"
)

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
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name, version, region, plan, vpc := d.Get("name").(string), d.Get("version").(string), d.Get("region").(string), d.Get("plan").(string), d.Get("vpc_id").(string)

	log.Printf("Creating db with : name: %s, version: %s, region: %s, plan: %s, vpc: %s", name, version, region, plan, vpc)
	client := meta.(*GrapheneDBClient).NewDatabasesClient()

	database, err := client.CreateDatabase(name, version, region, plan, vpc)
	if err != nil {
		return diag.Errorf("Error creating database %s: %s", name, err)
	}
	operationClient := meta.(*GrapheneDBClient).NewOperationsClient()
	databaseId, err := FetchDatabaseAsyncOperationStatus(database.OperationID, operationClient)
	if err != nil {
		return diag.Errorf("Error creating database %s: %s", name, err)
	}
	d.SetId(databaseId)

	return nil
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceData, err := getDatabaseResourceAttributesById(ctx,  meta.(*GrapheneDBClient).NewDatabasesClient(), d.Id())
	if(err != nil){
		return diag.FromErr(err)
	}
	if resourceData == nil {
		tflog.Warn(ctx, "Database not exists upstream and new one must be created")
		d.SetId("")
		return nil
	}
	err = AttributesToResourceData(resourceData, d)
	if  resourceData == nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name, version, region, plan, vpc := d.Get("name").(string), d.Get("version").(string), d.Get("region").(string), d.Get("plan").(string), d.Get("vpc_id").(string)
	tflog.Warn(ctx, "NEW STATE", map[string]interface{}{
		"name": name, 
		"version": version, 
		"region": region, 
		"plan": plan,
		"vpc_id": vpc,
	})


	databaseClient := meta.(*GrapheneDBClient).NewDatabasesClient()
	database, err := databaseClient.GetDatabaseInfo(ctx, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	if(database != nil) {
		name, version, region, _, vpc := d.Get("name").(string), d.Get("version").(string), d.Get("region").(string), d.Get("plan").(string), d.Get("vpc_id").(string)
		if(name != database.Name){
			return diag.FromErr(errors.New("name of database cannot be updated"))
		}
		if(version != database.VersionKey){
			return diag.FromErr(errors.New("version of database cannot be updated"))
		}
		if(region != database.AwsRegion){
			return diag.FromErr(errors.New("region of database cannot be updated"))
		}
		if(vpc != database.PrivateNetworkId){
			return diag.FromErr(errors.New("network id of database cannot be updated"))
		}
	} else {
		return diag.FromErr(errors.New("database now found"))
	}
	updateResult, err := databaseClient.UpgradeDatabaseInfo(ctx, plan, d.Id())
	if(err != nil){
		return diag.FromErr(err)
	}


	operationClient := meta.(*GrapheneDBClient).NewOperationsClient()
	databaseId, err := FetchDatabaseAsyncOperationStatus(updateResult.OperationID, operationClient)
	if err != nil {
		return diag.Errorf("Error updating database %s: %s", name, err)
	}
	d.SetId(databaseId)

	return nil
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return diag.Errorf("Database delete is called")
	return nil
}

func getDatabaseResourceAttributesById(ctx context.Context, client *client.DatabasesClient, databaseId string) (map[string]interface{}, error) {
	database, err := client.GetDatabaseInfo(ctx, databaseId)
	if err != nil {
		return nil, err
	}
	if(database == nil){
		return nil, nil
	}
	c := make(map[string]interface{})
  c["id"] = database.Id
	c["plan"]= database.Plan.PlanType
	c["region"]= database.AwsRegion
	c["vpc_id"]= database.PrivateNetworkId
	c["version"]= database.VersionKey
	c["name"] = database.Name
	tflog.Warn(ctx, "UPSTREAM DATABASE STATE",  c )
  return c, nil
}
