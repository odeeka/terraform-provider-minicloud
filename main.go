// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/odeeka/terraform-provider-minicloud/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

// // Default main() function
// func main() {
// 	var debug bool

// 	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
// 	flag.Parse()

// 	opts := providerserver.ServeOpts{
// 		// TODO: Update this string with the published name of your provider.
// 		// Also update the tfplugindocs generate command to either remove the
// 		// -provider-name flag or set its value to the updated provider name.
// 		Address: "registry.terraform.io/hashicorp/scaffolding",
// 		Debug:   debug,
// 	}

// 	err := providerserver.Serve(context.Background(), provider.New(version), opts)

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// }

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/minicloud. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/minicloud",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
