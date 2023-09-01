# Resource `graphenedb_database`

This resource manages a GrapheneDB database instance.

## Example Usage

```hcl
resource "graphenedb_database" "db" {
  name    = "my-db"
  version = "4.3.0"
  plan    = "standard"
  edition = "enterprise"
  vendor  = "neo4j"

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

- `name` (Required) - The name of the database instance.
- `version` (Required) - The version of Neo4j to use for the database instance. Available Neo4j versions can be find by fetching the following endpoint: https://api.db.graphenedb.com/deployments/versions.
- `plan` (Required) - The plan of the database instance. Available database plan can be find by fetching the following endpoint: https://api.db.graphenedb.com/deployments/plans.
- `edition` (Required) - The edition of Neo4j to use for the database instance.
- `vendor` (Required) - The vendor of Neo4j to use for the database instance. The available vendors are _graphneo_ and _ongdb_.
- `plugins` (Optional) - A list of plugin objects to install in the database instance. Each object must have the name and url attributes.
- `configuration` (Optional) - A list of configuration objects to set in the database instance. Each object must have the `key`, `value` and `secret` attributes.

## Attribute Reference

In addition to the arguments listed above, the following attributes are exported:

- `id` - The unique ID of the database instance.
- `domain_name` - The domain name to access the database instance.
- `http_port` - The HTTP port of the database instance.
- `bolt_port` - The Bolt port of the database instance.
- `vendor` - The vendor of the Neo4j instance.
