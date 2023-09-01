# terraform-provider-graphenedb

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Usage

```terraform
terraform {
	required_providers {
		graphenedb =  {
			version = "2.0.0"
			source = "github.com/quelhasu/graphenedb"
		}
	}
}

 provider "graphenedb" {
   environment_id = "graphenedb_environment_id"
   client_id      = "graphenedb_client_id"
   client_secret  = "graphenedb_client_secret"
 }
```

## Building The Provider

Clone the repository.

### Windows

Run `windows-build.sh` file.Change the version on the changes and run bash file. It buidl the project and copy new version in local terraform registry. not you can use the local registry provder like this

```sh
terraform {
	required_providers {
		graphenedb =  {
			version = "x.x.x" # build version
			source = "terraform.localhost.com/quelhasu/graphenedb"
		}
	}
}
```

### Other systems

Use make file for build and release

## Using the provider

## graphenedb_vpc_peering

Create a peering request to the provided vpc.

```tf
resource "graphenedb_vpc_peering" "vpc" {
	label           = "vpc_name"
	aws_account_id  = "vpc_aws_account_id"
	vpc_id          = "vpc_id"
	peer_vpc_region = "vpc_peer_region"
}
```

## graphenedb_database

Create a database with plugins url

```tf
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

## Release and publish on the registry

```sh
$ export GITHUB_TOKEN=...
$ git tag vX.Y.Z
$ git push origin vX.Y.Z
$ goreleaser release --clean --skip-sign
# sign manually the SHA file
$ gpg --detach-sign dist/terraform-provider-graphenedb_X.Y.Z_SHA256SUMS
```
