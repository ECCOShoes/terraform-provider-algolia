package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories provides the provider to the acceptance test
// framework under the "algolia" name.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"algolia": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck verifies the environment is configured before running
// acceptance tests. Acceptance tests talk to a real Algolia application and
// only run when TF_ACC is set.
func testAccPreCheck(t *testing.T) {
	required := []string{
		"ALGOLIA_APP_ID",
		"ALGOLIA_API_KEY",
	}
	for _, env := range required {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}
