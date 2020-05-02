package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/plugin"
        "github.com/hashicorp/terraform-plugin-sdk/terraform"
        "github.com/arnvid/terraform-provider-appstream/appstream"
)

func main() {
        plugin.Serve(&plugin.ServeOpts{
               ProviderFunc: appstream.Provider})
}

