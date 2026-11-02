package provider

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests run against a real server and change the tenant's
// configuration: each creates, updates, imports and destroys the objects in
// its kind's examples (examples/resources/vulnetix_<kind>), including
// singletons such as the QualityGate. Point them at a test tenant only.
//
//	TF_ACC=1 VULNETIX_ENDPOINT=https://www.vulnetix.com/api/site/v1 \
//	VULNETIX_CONFIG_TOKEN=vmt_... go test ./internal/provider -run TestAcc
//
// The token needs the owner or admin role and every served kind. Run the
// same suite under OpenTofu with TF_ACC_TERRAFORM_PATH=$(which tofu) and
// TF_ACC_PROVIDER_NAMESPACE=vulnetix.
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

// testAccKind runs one served kind's example lifecycle against the server
// named by the environment: create, plan-empty, update, import, destroy.
func testAccKind(t *testing.T, kind string) {
	t.Helper()
	var spec *kindSpec
	for i := range generatedKinds {
		if generatedKinds[i].Kind == kind {
			spec = &generatedKinds[i]
		}
	}
	if spec == nil {
		t.Fatalf("%s is not a served kind", kind)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { accPreCheck(t) },
		ProtoV6ProviderFactories: accFactories(),
		Steps:                    exampleSteps(t, *spec, "provider \"vulnetix\" {}\n", nil),
	})
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

func TestAccAiFirewallPolicy(t *testing.T)      { testAccKind(t, "AiFirewallPolicy") }
func TestAccAiGovernanceList(t *testing.T)      { testAccKind(t, "AiGovernanceList") }
func TestAccAiProvider(t *testing.T)            { testAccKind(t, "AiProvider") }
func TestAccApiSpecSource(t *testing.T)         { testAccKind(t, "ApiSpecSource") }
func TestAccJailExemption(t *testing.T)         { testAccKind(t, "JailExemption") }
func TestAccJailPolicy(t *testing.T)            { testAccKind(t, "JailPolicy") }
func TestAccMalwareIoc(t *testing.T)            { testAccKind(t, "MalwareIoc") }
func TestAccPackageFirewallPolicy(t *testing.T) { testAccKind(t, "PackageFirewallPolicy") }
func TestAccPackageList(t *testing.T)           { testAccKind(t, "PackageList") }
func TestAccQualityGate(t *testing.T)           { testAccKind(t, "QualityGate") }
func TestAccVexSource(t *testing.T)             { testAccKind(t, "VexSource") }

// Every served kind has its TestAcc<Kind>, so a newly served kind cannot
// ship without an acceptance test.
func TestEveryServedKindHasAcceptanceTest(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "acceptance_test.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	have := map[string]bool{}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok {
			have[fn.Name.Name] = true
		}
	}
	for _, k := range generatedKinds {
		if !have["TestAcc"+k.Kind] {
			t.Errorf("served kind %s has no TestAcc%s in acceptance_test.go", k.Kind, k.Kind)
		}
	}
}
