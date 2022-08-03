package graphenedb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseRestartCreate,
		ReadContext:   resourceDatabaseRestartRead,
		UpdateContext: resourceDatabaseRestartUpdate,
		DeleteContext: resourceDatabaseRestartDelete,

		Schema: map[string]*schema.Schema{
			"database_id": &schema.Schema{
				Type: schema.TypeString, 
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceDatabaseRestartCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())

	var diags diag.Diagnostics

	database_id := d.Get("database_id").(string)

	log.Printf("Restart the database identified by : ", database_id)
	client := m.(*GrapheneDBClient).NewDatabasesClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	database, err := client.RestartDatabase(database_id)
	if err != nil {
		return diag.Errorf("Error restarting database %s: %s", database_id, err)
	}

	operation := database.OperationID
	log.Printf("Fetch operation : ", operation)
	client_op := m.(*GrapheneDBClient).NewOperationsClient()

	operationDetail, err_op := client_op.FetchOperationDetail(operation)

	for operationDetail.Stopped == false && err_op == nil  {
		log.Println("Still fetching operation....")

		if err_op != nil {
			return diag.Errorf("Error fetching operation %s: %s", operation, err_op)
		}

		// wait 10 sec to avoid overloading the API.
		time.Sleep(10 * time.Second)

		operationDetail, err_op = client_op.FetchOperationDetail(operation)
	}

	if operationDetail.CurrentState != "finished"{
		return diag.Errorf("Failed restarting the database for operation %s", operation)
	}

	d.SetId(operationDetail.DatabaseId)
	
	return diags	
}

func resourceDatabaseRestartRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceDatabaseRestartUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceDatabaseRestartDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}
