terraform {
  required_version = ">= 1"
  required_providers {
    graphenedb = {
      version = "2.0.0"
      source  = "github.com/quelhasu/graphenedb"
    }
  }
}

provider "graphenedb" {
  endpoint = "https://api.graphenedb.com"
  api_key  = "api_key"
}

resource "graphenedb_vpc" "vpc" {
  label  = "vpc_label"
  region = "vpc_region"
  cidr   = "vpc_cidr"
}

resource "graphenedb_database" "db" {
  name    = "db_name"
  version = "db_version"
  region  = "db_region"
  plan    = "db_plan"
  vpc_id  = graphenedb_vpc.vpc.id

  plugins {
    name = "gds"
    kind = "extension"
    url  = "https://github.com/neo4j/graph-data-science/releases/download/2.1.5/neo4j-graph-data-science-2.1.5.zip"
  }

  plugins {
    name = "apoc"
    kind = "stored-procedure"
    url  = "https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/4.3.0.6/apoc-4.3.0.6-all.jar"
  }

  depends_on = [
    graphenedb_vpc.vpc
  ]
}
