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
  endpoint = "api_endpoint"
}
resource "graphenedb_pvc" "creation" {
  label  = "pvc_label"
  region = "pvc_region"
  cidr   = "pvc_cidr"
}

resource "graphenedb_database" "creation" {
  name    = "db_db"
  version = "db_version"
  region  = "db_region"
  plan    = "db_plan"
  cidr    = "db_cidr"
}
