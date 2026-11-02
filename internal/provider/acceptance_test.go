package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Acceptance tests run against a real server and change the tenant's
// configuration: destroy deletes its QualityGate. Point them at a test
// tenant only.
//
//	TF_ACC=1 VULNETIX_ENDPOINT=https://www.vulnetix.com/api/site/v1 \
//	VULNETIX_CONFIG_TOKEN=vmt_... go test ./internal/provider -run TestAcc
//
// resource.Test skips them without TF_ACC; accPreCheck skips them without
// the server settings, so TF_ACC alone never aims them at production.

func accPreCheck(t *testing.T) {
	t.Helper()
	for _, env := range []string{EnvToken, EnvEndpoint} {
		if os.Getenv(env) == "" {
			t.Skipf("%s is not set", env)
		}
	}
}

func accFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"vulnetix": providerserver.NewProtocol6WithError(New("acc")()),
	}
}

func TestAccKinds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { accPreCheck(t) },
		ProtoV6ProviderFactories: accFactories(),
		Steps: []resource.TestStep{{
			Config: `data "vulnetix_kinds" "all" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "api_version", "vulnetix.com/v1"),
				resource.TestCheckTypeSetElemNestedAttrs("data.vulnetix_kinds.all", "kinds.*", map[string]string{"kind": "QualityGate"}),
			),
		}},
	})
}

func TestAccQualityGate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { accPreCheck(t) },
		ProtoV6ProviderFactories: accFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
resource "vulnetix_quality_gate" "acc" {
  severity      = "critical"
  block_malware = true
  cooldown_days = 3
  eol = {
    block      = true
    severities = { retired = "critical" }
  }
  sca_autofix = { strategy = "safest", max_major_bump = 0 }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vulnetix_quality_gate.acc", "id", "default"),
					resource.TestCheckResourceAttr("vulnetix_quality_gate.acc", "severity", "critical"),
				),
			},
			{
				Config: `
resource "vulnetix_quality_gate" "acc" {
  severity      = "high"
  block_malware = true
  cooldown_days = 3
}
`,
				Check: resource.TestCheckResourceAttr("vulnetix_quality_gate.acc", "severity", "high"),
			},
			{
				ResourceName:  "vulnetix_quality_gate.acc",
				ImportState:   true,
				ImportStateId: "default",
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 || s[0].Attributes["severity"] != "high" || s[0].Attributes["block_malware"] != "true" {
						return fmt.Errorf("imported %v", s)
					}
					return nil
				},
			},
		},
	})
}
