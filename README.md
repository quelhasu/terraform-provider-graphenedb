# terraform-provider-graphenedb

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Usage

```
terraform {
	required_providers {
		graphenedb =  {
			version = "2.0.0"
			source = "github.com/quelhasu/graphenedb"
		}
	}
}

provider  "graphenedb" {
	endpoint =  "https://api.graphenedb.com"
	api_key =  "your_api_key"
}
```

## Building The Provider

Clone the repository.

- **Windows:** Run `windows-build.sh` file.Change the version on the changes and run bash file. It buidl the project and copy new version in local terraform registry. not you can use the local registry provder like this

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

- **Other systems:** Use make file for build and release

## Using the provider

## graphenedb_vpc

```
resource  "graphenedb_vpc"  "vpc" {
	label =  "vpc_label"
	region =  "vpc_region"
	cidr =  "vpc_cidr"
}
```

## graphenedb_database

```
resource  "graphenedb_database"  "db" {
	name =  "db_name"
	version =  "db_version"
	region =  "db_region"
	plan =  "db_plan"
	vpc_id =  graphenedb_vpc.vpc.id

	plugins {
		name =  "gds"
		kind =  "extension"
		url =  "https://github.com/neo4j/graph-data-science/releases/download/2.1.5/neo4j-graph-data-science-2.1.5.zip"
	}

	plugins {
		name =  "apoc"
		kind =  "stored-procedure"
		url =  "https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/4.3.0.6/apoc-4.3.0.6-all.jar"
	}

	depends_on  =  [
		graphenedb_vpc.vpc
	]
}
```

## Release and publish on the registry

```sh
$ export GITHUB_TOKEN=...
$ git tag vX.Y.Z
$ git push origin vX.Y.Z
$ goreleaser release --rm-dist --skip-sign
# sign manually the SHA file
$ gpg --detach-sign dist/terraform-provider-graphenedb_X.Y.Z_SHA256SUMS
```
