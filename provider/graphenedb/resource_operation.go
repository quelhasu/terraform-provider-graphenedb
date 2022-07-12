package graphenedb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOperationCreate,
		ReadContext:   resourceOperationRead,
		UpdateContext: resourceOperationUpdate,
		DeleteContext: resourceOperationDelete,

		Schema: map[string]*schema.Schema{
			"operation": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceOperationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	var diags diag.Diagnostics

	operation := d.Get("operation").(string)

	log.Printf("Fetch operation : ", operation)
	client := m.(*GrapheneDBClient).NewOperationsClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	operationDetail, err := client.FetchOperationDetail(operation)
	if err != nil {
		return diag.Errorf("Error fetching operation %s: %s", operation, err)
	}

	log.Printf("Operation detail : ", operationDetail)
	d.SetId(operationDetail.DatabaseId)
	return diags	
}

func resourceOperationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceOperationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceOperationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}
