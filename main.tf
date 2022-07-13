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
  user     = "api_user"
  password = "api_key"
  endpoint = "https://api.graphenedb.com"
}

resource "graphenedb_vpc" "creation" {
  label  = "vpc_label"
  region = "vpc_region"
  cidr   = "vpc_cidr"
}


resource "graphenedb_database" "creation" {
  name    = "db_name"
  version = "db_version"
  region  = "db_region"
  plan    = "db_plan"
  vpc_id  = graphenedb_vpc.creation.id

  depends_on = [
    graphenedb_vpc.creation
  ]
}

resource "graphenedb_plugin" "add_gds" {
  name        = "gds"
  kind        = "extension"
  url         = "https://github.com/neo4j/graph-data-science/releases/download/2.1.5/neo4j-graph-data-science-2.1.5.zip"
  database_id = graphenedb_database.creation.id

  depends_on = [
    graphenedb_database.creation
  ]
}
resource "graphenedb_plugin" "add_apoc" {
  name        = "apoc"
  kind        = "stored-procedure"
  url         = "https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/4.3.0.6/apoc-4.3.0.6-all.jar"
  database_id = graphenedb_database.creation.id

  depends_on = [
    graphenedb_database.creation
  ]
}


