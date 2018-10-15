package main

import (
	"github.com/dbgeek/terraform-provider-oraclerdbms/oraclerdbms"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: oraclerdbms.Provider})
}
