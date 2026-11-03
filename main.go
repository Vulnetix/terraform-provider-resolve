// Command terraform-provider-vulnetix is the vulnetix provider for Terraform
// and OpenTofu.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/provider"
)

// Generate the resources from the config kind schemas, then the docs from
// the provider schema. Run go generate ./... after bumping config-schemas.
//go:generate go run ./cmd/genresources -out internal/provider -examples-out examples/resources
//go:generate go run ./cmd/genresources -schemas cmd/genresources/testdata/schemas -out internal/provider -suffix _test -var testKinds
//go:generate go tool tfplugindocs generate --provider-name vulnetix

// version is set by goreleaser.
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers such as delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		// The registry address, which is also the source users write in
		// required_providers.
		Address: "registry.terraform.io/vulnetix/resolve",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
