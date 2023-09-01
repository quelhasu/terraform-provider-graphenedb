# GrapheneDB Provider

This provider is used to interact with the new [GrapheneDB service](https://console.graphenedb.com).

## Example Usage

```hcl
terraform {
  required_version = ">= 1"
  required_providers {
    graphenedb = {
      version = "0.0.1"
      source  = "github.com/quelhasu/graphenedb"
    }
  }
}

provider "graphenedb" {
  environment_id = "graphenedb_environment_id"
  client_id      = "graphenedb_client_id"
  client_secret  = "graphenedb_client_secret"
}

resource "graphenedb_vpc_peering" "vpc" {
  label           = "vpc_name"
  aws_account_id  = "vpc_aws_account_id"
  vpc_id          = "vpc_id"
  peer_vpc_region = "vpc_peer_region"
}

resource "graphenedb_database" "db" {
  name    = "db_name"
  version = "db_version"
  plan    = "db_plan"
  edition = "db_edition"
  vendor  = "db_vendor"
  plugins {
    name = "gds"
    url  = "https://github.com/neo4j/graph-data-science/releases/download/2.1.5/neo4j-graph-data-science-2.1.5.zip"
  }

  plugins {
    name = "apoc"
    url  = "https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/4.3.0.6/apoc-4.3.0.6-all.jar"
  }

  configuration {
    key    = "dbms.transaction.timeout"
    value  = "45s"
    secret = false
  }

  depends_on = [
    graphenedb_vpc_peering.vpc
  ]
}
```

## Argument Reference

The following arguments are supported:

- `environment_id` - (Required) The ID of the GrapheneDB environment after creating it directly from the console.
- `client_id` - (Required) The ID of the GrapheneDB organization client.
- `client_secret` - (Required) The secret associated with the GrapheneDB organization client.

## Resources

- `graphenedb_vpc_peering` - Creates a VPC peering connection.
- `graphenedb_database` - Creates a new GrapheneDB database with plugins.
