package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func spec(t *testing.T, put map[string]any) map[string]any {
	t.Helper()
	doc, _ := put["document"].(map[string]any)
	s, _ := doc["spec"].(map[string]any)
	if s == nil {
		t.Fatalf("PUT body has no document.spec: %v", put)
	}
	return s
}

func checkPut(t *testing.T, f *fakeAPI, check func(body, spec map[string]any) error) resource.TestCheckFunc {
	return func(*terraform.State) error {
		body := f.lastPut()
		if body == nil {
			return fmt.Errorf("no PUT was made")
		}
		return check(body, spec(t, body))
	}
}

func TestQualityGateLifecycle(t *testing.T) {
	f, srv := newFakeAPI(t)
	// Fields the tenant's row holds that the configuration never sets. A
	// read that copied them into state would show a diff on every plan.
	f.defaults["QualityGate"] = map[string]any{
		"versionLag": json.Number("3"),
		"eol":        map[string]any{"severities": map[string]any{"retired": "critical"}},
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		CheckDestroy: func(*terraform.State) error {
			if len(f.objects) != 0 {
				return fmt.Errorf("objects left after destroy: %v", f.objects)
			}
			if len(f.deletes) == 0 || f.deletes[len(f.deletes)-1] != "QualityGate/default?manager=terraform" {
				return fmt.Errorf("deletes = %v", f.deletes)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: providerBlock(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity      = "high"
  block_malware = true
  cooldown_days = 7
  eol = {
    block = true
  }
  labels = { team = "platform" }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "id", "default"),
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "name", "default"),
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "severity", "high"),
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "eol.block", "true"),
					resource.TestCheckNoResourceAttr("vulnetix_quality_gate.gate", "version_lag"),
					resource.TestCheckNoResourceAttr("vulnetix_quality_gate.gate", "eol.severities"),
					checkPut(t, f, func(body, s map[string]any) error {
						if body["manager"] != "terraform" {
							return fmt.Errorf("manager = %v", body["manager"])
						}
						doc := body["document"].(map[string]any)
						if doc["apiVersion"] != "vulnetix.com/v1" || doc["kind"] != "QualityGate" {
							return fmt.Errorf("document header = %v %v", doc["apiVersion"], doc["kind"])
						}
						if doc["metadata"].(map[string]any)["labels"].(map[string]any)["team"] != "platform" {
							return fmt.Errorf("labels not sent: %v", doc["metadata"])
						}
						if s["severity"] != "high" || s["blockMalware"] != true || s["cooldownDays"] != json.Number("7") {
							return fmt.Errorf("spec = %v", s)
						}
						if eol, _ := s["eol"].(map[string]any); eol["block"] != true {
							return fmt.Errorf("spec.eol = %v", s["eol"])
						}
						if _, ok := s["versionLag"]; ok {
							return fmt.Errorf("unset field sent: %v", s)
						}
						return nil
					}),
				),
			},
			{
				Config: providerBlock(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity = "critical"
  eol = {
    block = true
    severities = { within_30_days = "high" }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "severity", "critical"),
					resource.TestCheckNoResourceAttr("vulnetix_quality_gate.gate", "block_malware"),
					resource.TestCheckResourceAttr("vulnetix_quality_gate.gate", "eol.severities.within_30_days", "high"),
					checkPut(t, f, func(_, s map[string]any) error {
						if _, ok := s["blockMalware"]; ok {
							return fmt.Errorf("removed argument still sent: %v", s)
						}
						sev := s["eol"].(map[string]any)["severities"].(map[string]any)
						if sev["within30Days"] != "high" {
							return fmt.Errorf("spec.eol.severities = %v", sev)
						}
						return nil
					}),
				),
			},
			{
				// Import reads every field the server holds, including the
				// ones no configuration set.
				ResourceName:  "vulnetix_quality_gate.gate",
				ImportState:   true,
				ImportStateId: "default",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("imported %d states", len(states))
					}
					a := states[0].Attributes
					for k, want := range map[string]string{
						"id": "default", "name": "default", "severity": "critical",
						"version_lag": "3", "block_malware": "true", "cooldown_days": "7",
						"eol.severities.retired": "critical", "eol.severities.within_30_days": "high",
					} {
						if a[k] != want {
							return fmt.Errorf("imported %s = %q, want %q (all: %v)", k, a[k], want, a)
						}
					}
					return nil
				},
			},
		},
	})
}

func TestQualityGateNameIsDefault(t *testing.T) {
	_, srv := newFakeAPI(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  name     = "other"
  severity = "high"
}
`,
			ExpectError: regexp.MustCompile(`value must be one of`),
		}, {
			Config: providerBlock(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity = "urgent"
}
`,
			ExpectError: regexp.MustCompile(`value must be one of`),
		}},
	})
}

func TestQualityGateDriftAndRemoval(t *testing.T) {
	f, srv := newFakeAPI(t)
	cfg := providerBlock(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				// Someone changes the gate in the console: the next plan
				// puts it back.
				PreConfig: func() {
					f.mu.Lock()
					f.objects["QualityGate/default"].Spec["severity"] = "low"
					f.mu.Unlock()
				},
				Config: cfg,
				Check: checkPut(t, f, func(_, s map[string]any) error {
					if s["severity"] != "high" {
						return fmt.Errorf("drift not corrected: %v", s)
					}
					return nil
				}),
			},
			{
				// Deleted outside Terraform: recreated.
				PreConfig: func() {
					f.mu.Lock()
					delete(f.objects, "QualityGate/default")
					f.mu.Unlock()
				},
				Config: cfg,
				Check: func(*terraform.State) error {
					if _, ok := f.objects["QualityGate/default"]; !ok {
						return fmt.Errorf("not recreated")
					}
					return nil
				},
			},
		},
	})
}

func TestWidgetShapesAndWriteOnly(t *testing.T) {
	f, srv := newFakeAPI(t)
	widget := func(color string, version int, token string) string {
		return providerBlock(srv.URL) + fmt.Sprintf(`
resource "vulnetix_widget" "w" {
  name             = "w1"
  color            = %q
  size             = 3
  ratio            = 0.5
  enabled          = true
  tags             = ["a", "b"]
  env              = { REGION = "ap-southeast-2" }
  extra            = jsonencode({ b = 1, a = [true] })
  token_wo         = %q
  token_wo_version = %d
  rules = [
    { purl = "pkg:npm/left-pad", reason = "component_not_present" },
  ]
  nested = {
    deep_flag      = true
    within_30_days = "soon"
  }
}
`, color, token, version)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config: widget("red", 1, "s3cret"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vulnetix_widget.w", "id", "w1"),
					resource.TestCheckNoResourceAttr("vulnetix_widget.w", "token_wo"),
					resource.TestCheckResourceAttr("vulnetix_widget.w", "token_wo_version", "1"),
					resource.TestCheckResourceAttr("vulnetix_widget.w", "rules.0.purl", "pkg:npm/left-pad"),
					checkPut(t, f, func(_, s map[string]any) error {
						if s["token"] != "s3cret" {
							return fmt.Errorf("secret not sent on create: %v", s)
						}
						if s["ratio"] != json.Number("0.5") || s["size"] != json.Number("3") {
							return fmt.Errorf("numbers = %v %v", s["ratio"], s["size"])
						}
						if x, _ := s["extra"].(map[string]any); x["b"] != json.Number("1") {
							return fmt.Errorf("JSON attribute not decoded: %v", s["extra"])
						}
						n := s["nested"].(map[string]any)
						if n["deepFlag"] != true || n["within30Days"] != "soon" {
							return fmt.Errorf("nested = %v", n)
						}
						r := s["rules"].([]any)[0].(map[string]any)
						if r["purl"] != "pkg:npm/left-pad" || r["reason"] != "component_not_present" {
							return fmt.Errorf("rules = %v", s["rules"])
						}
						return nil
					}),
				),
			},
			{
				// Same version: the secret is not sent again.
				Config: widget("blue", 1, "s3cret"),
				Check: checkPut(t, f, func(_, s map[string]any) error {
					if _, ok := s["token"]; ok || s["color"] != "blue" {
						return fmt.Errorf("spec = %v", s)
					}
					return nil
				}),
			},
			{
				Config: widget("blue", 2, "n3w"),
				Check: checkPut(t, f, func(_, s map[string]any) error {
					if s["token"] != "n3w" {
						return fmt.Errorf("new secret not sent: %v", s)
					}
					return nil
				}),
			},
			{
				// A new name replaces the object.
				Config: regexp.MustCompile(`"w1"`).ReplaceAllString(widget("blue", 2, "n3w"), `"w2"`),
				Check: func(*terraform.State) error {
					if _, ok := f.objects["Widget/w1"]; ok {
						return fmt.Errorf("old object not deleted")
					}
					if d, ok := f.objects["Widget/w2"]; !ok || d.Spec["token"] != "n3w" {
						return fmt.Errorf("replacement missing or lost its secret: %v", d)
					}
					return nil
				},
			},
		},
	})
}

func TestUnauthorizedIsReported(t *testing.T) {
	for name, creds := range map[string]string{
		"wrong api key": `api_key  = "wrong"
  org_id   = "` + fakeOrg + `"`,
		"wrong org": `api_key  = "` + fakeAPIKey + `"
  org_id   = "other-org"`,
		"wrong bearer": `token    = "wrong"`,
	} {
		t.Run(name, func(t *testing.T) {
			_, srv := newFakeAPI(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				Steps: []resource.TestStep{{
					Config: `
provider "vulnetix" {
  endpoint = "` + srv.URL + `/v1"
  ` + creds + `
}
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`,
					ExpectError: regexp.MustCompile(`HTTP 401:\s+not authenticated`),
				}},
			})
		})
	}
}

// Both credential forms apply, and the provider sends the header the
// Vulnetix CLI sends: ApiKey <org>:<key>, or Bearer <token>.
func TestCredentialForms(t *testing.T) {
	for name, tc := range map[string]struct {
		block func(url string) string
		auth  string
	}{
		"api key": {providerBlock, "ApiKey " + fakeOrg + ":" + fakeAPIKey},
		"api key with the org prefix": {func(url string) string {
			return `
provider "vulnetix" {
  endpoint = "` + url + `/v1"
  api_key  = "` + fakeOrg + ":" + fakeAPIKey + `"
  org_id   = "` + fakeOrg + `"
}
`
		}, "ApiKey " + fakeOrg + ":" + fakeAPIKey},
		"bearer": {providerBlockBearer, "Bearer " + fakeBearer},
	} {
		t.Run(name, func(t *testing.T) {
			f, srv := newFakeAPI(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				Steps: []resource.TestStep{{
					Config: tc.block(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`,
					Check: func(*terraform.State) error {
						f.mu.Lock()
						defer f.mu.Unlock()
						if f.lastAuth != tc.auth {
							return fmt.Errorf("Authorization = %q, want %q", f.lastAuth, tc.auth)
						}
						return nil
					},
				}},
			})
		})
	}
}

// The credentials may come from the environment instead of the configuration.
func TestCredentialsFromEnvironment(t *testing.T) {
	f, srv := newFakeAPI(t)
	t.Setenv(EnvAPIKey, fakeAPIKey)
	t.Setenv(EnvOrgID, fakeOrg)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: `
provider "vulnetix" {
  endpoint = "` + srv.URL + `/v1"
}
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`,
			Check: func(*terraform.State) error {
				f.mu.Lock()
				defer f.mu.Unlock()
				if len(f.objects) == 0 {
					return fmt.Errorf("nothing was applied")
				}
				return nil
			},
		}},
	})
}

// Configuring both an ApiKey and a bearer token, or an ApiKey without its
// organisation, is an error before any request is made.
func TestCredentialConflicts(t *testing.T) {
	for name, tc := range map[string]struct{ creds, want string }{
		"both": {`api_key = "k"
  org_id  = "o"
  token   = "t"`, `Two Vulnetix credentials`},
		"key without org": {`api_key = "k"`, `api_key needs org_id`},
		"org without key": {`org_id = "o"`, `org_id needs api_key`},
	} {
		t.Run(name, func(t *testing.T) {
			f, srv := newFakeAPI(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				Steps: []resource.TestStep{{
					Config: `
provider "vulnetix" {
  endpoint = "` + srv.URL + `/v1"
  ` + tc.creds + `
}
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`,
					ExpectError: regexp.MustCompile(tc.want),
				}},
			})
			if f.putCount() != 0 {
				t.Error("a request was made with conflicting credentials")
			}
		})
	}
}
