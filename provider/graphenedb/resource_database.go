package graphenedb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())

	var diags diag.Diagnostics

	name, version, region, plan, vpc := d.Get("name").(string), d.Get("version").(string), d.Get("region").(string), d.Get("plan").(string), d.Get("vpc_id").(string)

	log.Printf("Creating db with : name: %s, version: %s, region: %s, plan: %s, vpc: %s", name, version, region, plan, vpc)
	client := m.(*GrapheneDBClient).NewDatabasesClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	database, err := client.CreateDatabase(name, version, region, plan, vpc)
	if err != nil {
		return diag.Errorf("Error creating database %s: %s", name, err)
	}

	operation := database.OperationID
	log.Printf("Fetch operation: %s", operation)
	client_op := m.(*GrapheneDBClient).NewOperationsClient()

	operationDetail, err_op := client_op.FetchOperationDetail(operation)

	for !operationDetail.Stopped && err_op == nil  {
		log.Println("Still fetching operation....")

		if err_op != nil {
			return diag.Errorf("Error fetching operation %s: %s", operation, err_op)
		}

		// wait 10 sec to avoid overloading the API.
		time.Sleep(10 * time.Second)

		operationDetail, err_op = client_op.FetchOperationDetail(operation)
	}

	if operationDetail.CurrentState != "finished"{
		return diag.Errorf("Failed creating the database for operation %s", operation)
	}

	d.SetId(operationDetail.DatabaseId)
	return diags	
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*GrapheneDBClient).NewDatabasesClient()
	database, err := client.GetDatabaseInfo(ctx, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Warn(ctx, "DATABASE INFO",  map[string]interface{}{
    "database": fmt.Sprintf("%+v", database),
	})
	if(database != nil) {
		tflog.Warn(ctx, "DATABASE PLAN TYPE", map[string]interface{}{
			"PlanType":database.Plan.PlanType,
		})
		d.Set("plan", database.Plan.PlanType)
	} else{
		d.SetId("")
	}
	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}
