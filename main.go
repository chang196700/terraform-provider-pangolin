package main

import (
	"context"
	"flag"
	"log"

	"github.com/chang196700/terraform-provider-pangolin/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	debug := flag.Bool("debug", false, "Run with debugger support")
	flag.Parse()
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/chang196700/pangolin", Debug: *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
