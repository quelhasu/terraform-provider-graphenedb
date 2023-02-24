package graphenedb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	graphendbclient "github.com/quelhasu/terraform-provider-graphenedb/graphendb-client"
)

func resourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPeeringCreate,
		ReadContext:   resourceVpcPeeringRead,
		UpdateContext: resourceVpcPeeringUpdate,
		DeleteContext: resourceVpcPeeringDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"aws_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"peer_vpc_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"peering_id": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
		},
	}
}

func resourceVpcPeeringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	var diags diag.Diagnostics

	result, err := m.(*graphendbclient.RestApiClient).CreateVpcPeering(ctx, graphendbclient.VpcPeeringInfo{
		Label:         d.Get("label").(string),
		AwsAccountId:  d.Get("aws_account_id").(string),
		VpcId:         d.Get("vpc_id").(string),
		PeerVpcRegion: d.Get("peer_vpc_region").(string),
	})

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result.PeeringConnectionID)
	d.Set("peering_id", result.ID)

	return diags
}

func resourceVpcPeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceVpcPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// cli.Debug.Printf("Resource state: %#v", d.State())
	return nil
}

func resourceVpcPeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := meta.(*graphendbclient.RestApiClient).DeleteVPCPeering(ctx, d.Get("peering_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
