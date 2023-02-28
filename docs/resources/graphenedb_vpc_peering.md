# Resource `graphenedb_vpc_peering`

This resource creates a VPC peering connection between a GrapheneDB environment and a VPC in your AWS account. The peering connection allows communication between resources in the GrapheneDB environment and resources in your VPC, as if they were part of the same network.

## Example Usage

```hcl
resource "graphenedb_vpc_peering" "vpc" {
  label           = "vpc_name"
  aws_account_id  = "vpc_aws_account_id"
  vpc_id          = "vpc_id"
  peer_vpc_region = "vpc_peer_region"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) A label for the VPC peering connection.
- `aws_account_id` - (Required) The AWS account ID that owns the VPC to peer with.
- `vpc_id` - (Required) The ID of the VPC to peer with.
- `peer_vpc_region` - (Required) The region in which the peer VPC is located.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the VPC peering connection. (e.g. `pcx-XXXXXXX`)
- `vpc_id` - The ID of the VPC to peer with. (e.g. `vpc-XXXXXXX`)
