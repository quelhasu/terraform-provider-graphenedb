package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ervinjohnson/terraform-provider-graphenedb/provider/gdb"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gdb.Provider
	})
}