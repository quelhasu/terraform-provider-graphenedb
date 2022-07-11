package graphenedb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePvc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePvcCreate,
		ReadContext:   resourcePvcRead,
		UpdateContext: resourcePvcUpdate,
		DeleteContext: resourcePvcDelete,

		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourcePvcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())

	label, cidr, region := d.Get("label").(string), d.Get("cidr").(string), d.Get("region").(string)

	client := m.(*GrapheneDBClient).NewPvcClient()
	// cli.Debug.Printf("Resource state: %s %s %s %s %#v", name, version, awsRegion, plan, client)

	pvc, err := client.CreatePvc(label, cidr, region)
	if err != nil {
		return diag.Errorf("Error creating private network %s: %s", label, err)
	}

	d.SetId(pvc.ID)
	return nil
}

func resourcePvcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourcePvcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourcePvcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}
