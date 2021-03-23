package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/julian3xl/terraform-provider-appstream/appstream"
)

func main() {
        plugin.Serve(&plugin.ServeOpts{
               ProviderFunc: appstream.Provider})
}

