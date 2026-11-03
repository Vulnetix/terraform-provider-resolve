package provider

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
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
//	VULNETIX_API_KEY=... VULNETIX_ORG_ID=... go test ./internal/provider -run TestAcc
//
// or, with a bearer token (an Authentik access token) instead of the ApiKey,
// VULNETIX_API_TOKEN=... The principal needs the owner or admin role: use the
// ApiKey of an admin or owner service account. Run the same suite under
// OpenTofu with TF_ACC_TERRAFORM_PATH=$(which tofu) and
// TF_ACC_PROVIDER_NAMESPACE=vulnetix.
//
// resource.Test skips them without TF_ACC; accPreCheck skips them without
// the server settings, so TF_ACC alone never aims them at production.

// accMissing is what keeps an acceptance run from starting: the endpoint, and
// either VULNETIX_API_KEY with VULNETIX_ORG_ID or VULNETIX_API_TOKEN.
func accMissing(getenv func(string) string) string {
	if getenv(EnvEndpoint) == "" {
		return EnvEndpoint
	}
	if getenv(EnvToken) != "" {
		return ""
	}
	for _, env := range []string{EnvAPIKey, EnvOrgID} {
		if getenv(env) == "" {
			return env + " (with " + EnvAPIKey + " and " + EnvOrgID + ") or " + EnvToken
		}
	}
	return ""
}

func accPreCheck(t *testing.T) {
	t.Helper()
	if m := accMissing(os.Getenv); m != "" {
		t.Skipf("%s is not set", m)
	}
}

func TestAccGate(t *testing.T) {
	env := func(kv ...string) func(string) string {
		m := map[string]string{}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return func(k string) string { return m[k] }
	}
	for name, tc := range map[string]struct {
		env     func(string) string
		missing bool
	}{
		"api key, org and endpoint": {env(EnvEndpoint, "e", EnvAPIKey, "k", EnvOrgID, "o"), false},
		"bearer token and endpoint": {env(EnvEndpoint, "e", EnvToken, "t"), false},
		"no endpoint":               {env(EnvAPIKey, "k", EnvOrgID, "o"), true},
		"api key without org":       {env(EnvEndpoint, "e", EnvAPIKey, "k"), true},
		"org without api key":       {env(EnvEndpoint, "e", EnvOrgID, "o"), true},
		"no credential":             {env(EnvEndpoint, "e"), true},
	} {
		if got := accMissing(tc.env) != ""; got != tc.missing {
			t.Errorf("%s: skipped = %v, want %v", name, got, tc.missing)
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
	// A kind that needs something outside the tenant (a GitHub App installation
	// to adopt, a Wiz account, a real mailbox for an invitation) can be left
	// out of a run: VULNETIX_ACC_SKIP_KINDS=GitHubInstallation,WizConnection.
	for _, s := range strings.Split(os.Getenv("VULNETIX_ACC_SKIP_KINDS"), ",") {
		if strings.TrimSpace(s) == kind {
			t.Skipf("%s is in VULNETIX_ACC_SKIP_KINDS", kind)
		}
	}
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

func TestAccAgentProfile(t *testing.T)          { testAccKind(t, "AgentProfile") }
func TestAccAiFirewallPolicy(t *testing.T)      { testAccKind(t, "AiFirewallPolicy") }
func TestAccAiGovernanceList(t *testing.T)      { testAccKind(t, "AiGovernanceList") }
func TestAccAiProvider(t *testing.T)            { testAccKind(t, "AiProvider") }
func TestAccAlertTrigger(t *testing.T)          { testAccKind(t, "AlertTrigger") }
func TestAccApiSpecSource(t *testing.T)         { testAccKind(t, "ApiSpecSource") }
func TestAccAwsAccount(t *testing.T)            { testAccKind(t, "AwsAccount") }
func TestAccAwsConnection(t *testing.T)         { testAccKind(t, "AwsConnection") }
func TestAccCampaign(t *testing.T)              { testAccKind(t, "Campaign") }
func TestAccCloudSensors(t *testing.T)          { testAccKind(t, "CloudSensors") }
func TestAccComplianceTarget(t *testing.T)      { testAccKind(t, "ComplianceTarget") }
func TestAccCrew(t *testing.T)                  { testAccKind(t, "Crew") }
func TestAccDetectionRule(t *testing.T)         { testAccKind(t, "DetectionRule") }
func TestAccDetectionRuleFeed(t *testing.T)     { testAccKind(t, "DetectionRuleFeed") }
func TestAccFactoryBudget(t *testing.T)         { testAccKind(t, "FactoryBudget") }
func TestAccFactoryManagedOrg(t *testing.T)     { testAccKind(t, "FactoryManagedOrg") }
func TestAccFactorySchedule(t *testing.T)       { testAccKind(t, "FactorySchedule") }
func TestAccFactorySettings(t *testing.T)       { testAccKind(t, "FactorySettings") }
func TestAccFactoryTrigger(t *testing.T)        { testAccKind(t, "FactoryTrigger") }
func TestAccFactoryWorkflow(t *testing.T)       { testAccKind(t, "FactoryWorkflow") }
func TestAccForgeConnection(t *testing.T)       { testAccKind(t, "ForgeConnection") }
func TestAccGitHubInstallation(t *testing.T)    { testAccKind(t, "GitHubInstallation") }
func TestAccJailExemption(t *testing.T)         { testAccKind(t, "JailExemption") }
func TestAccJailPolicy(t *testing.T)            { testAccKind(t, "JailPolicy") }
func TestAccLibraryItem(t *testing.T)           { testAccKind(t, "LibraryItem") }
func TestAccLogSource(t *testing.T)             { testAccKind(t, "LogSource") }
func TestAccMalwareIoc(t *testing.T)            { testAccKind(t, "MalwareIoc") }
func TestAccMember(t *testing.T)                { testAccKind(t, "Member") }
func TestAccMessagingConnection(t *testing.T)   { testAccKind(t, "MessagingConnection") }
func TestAccMessagingRoute(t *testing.T)        { testAccKind(t, "MessagingRoute") }
func TestAccMessagingTemplate(t *testing.T)     { testAccKind(t, "MessagingTemplate") }
func TestAccNotificationDefaults(t *testing.T)  { testAccKind(t, "NotificationDefaults") }
func TestAccOrganization(t *testing.T)          { testAccKind(t, "Organization") }
func TestAccPackageFirewallPolicy(t *testing.T) { testAccKind(t, "PackageFirewallPolicy") }
func TestAccPackageList(t *testing.T)           { testAccKind(t, "PackageList") }
func TestAccQualityGate(t *testing.T)           { testAccKind(t, "QualityGate") }
func TestAccRemediationPolicy(t *testing.T)     { testAccKind(t, "RemediationPolicy") }
func TestAccRepositorySettings(t *testing.T)    { testAccKind(t, "RepositorySettings") }
func TestAccRiskStrategy(t *testing.T)          { testAccKind(t, "RiskStrategy") }
func TestAccRuntimeSensorPolicy(t *testing.T)   { testAccKind(t, "RuntimeSensorPolicy") }
func TestAccSandboxEgressPolicy(t *testing.T)   { testAccKind(t, "SandboxEgressPolicy") }
func TestAccSandboxLaunchProfile(t *testing.T)  { testAccKind(t, "SandboxLaunchProfile") }
func TestAccSavedSearch(t *testing.T)           { testAccKind(t, "SavedSearch") }
func TestAccScannerProfile(t *testing.T)        { testAccKind(t, "ScannerProfile") }
func TestAccScheduledReport(t *testing.T)       { testAccKind(t, "ScheduledReport") }
func TestAccServiceAccount(t *testing.T)        { testAccKind(t, "ServiceAccount") }
func TestAccSuppression(t *testing.T)           { testAccKind(t, "Suppression") }
func TestAccRegistryConnection(t *testing.T)    { testAccKind(t, "RegistryConnection") }
func TestAccSsvcPolicy(t *testing.T)            { testAccKind(t, "SsvcPolicy") }
func TestAccTeaAccessPolicy(t *testing.T)       { testAccKind(t, "TeaAccessPolicy") }
func TestAccTeaCollection(t *testing.T)         { testAccKind(t, "TeaCollection") }
func TestAccTeaComponent(t *testing.T)          { testAccKind(t, "TeaComponent") }
func TestAccTeaProduct(t *testing.T)            { testAccKind(t, "TeaProduct") }
func TestAccTeaPublicationTarget(t *testing.T)  { testAccKind(t, "TeaPublicationTarget") }
func TestAccVaultEntry(t *testing.T)            { testAccKind(t, "VaultEntry") }
func TestAccVexSource(t *testing.T)             { testAccKind(t, "VexSource") }
func TestAccWizConnection(t *testing.T)         { testAccKind(t, "WizConnection") }

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
