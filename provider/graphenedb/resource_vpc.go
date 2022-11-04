package graphenedb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)

func resourceVpc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcCreate,
		ReadContext:   resourceVpcRead,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	var diags diag.Diagnostics
	
	result, err := m.(*graphendbclient.RestApiClient).CreateVPC(ctx, graphendbclient.VpcInfo{
		Label: d.Get("label").(string),
		AwsRegion: d.Get("region").(string),
		CidrBlock: d.Get("cidr").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result.Id)
	
	return diags
}

func resourceVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := meta.(*graphendbclient.RestApiClient).DeleteVPC(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
