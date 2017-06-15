package main

import (
	gdb "github.com/ervinjohnson/terraform-provider-graphenedb/provider/graphenedb"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gdb.Provider,
	})
}
