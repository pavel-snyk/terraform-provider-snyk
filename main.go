package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/pavel-snyk/terraform-provider-snyk/internal/provider"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with debug support")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pavel-snyk/snyk",
		Debug:   debugMode,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
