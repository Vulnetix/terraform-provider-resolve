package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// A Member is applied and then waits for the invitation to be accepted:
// the phase is Pending, the apply succeeds, and a refresh keeps saying so.
func TestStatusPending(t *testing.T) {
	_, srv := newFakeAPI(t)
	cfg := providerBlock(srv.URL) + `
resource "vulnetix_member" "dev" {
  name  = "dev"
  email = "dev@acme.com"
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vulnetix_member.dev", "status.phase", "Pending"),
					resource.TestCheckResourceAttr("vulnetix_member.dev", "status.extra.invitationExpires", "2026-12-01T00:00:00Z"),
					resource.TestCheckNoResourceAttr("vulnetix_member.dev", "status.authorization_url"),
				),
			},
			// The refresh before the second plan reads the object again and
			// finds the same phase; nothing to change.
			{Config: cfg, PlanOnly: true},
		},
	})
}

// A Slack connection is applied and then waits for a person to authorize it
// in a browser: the URL is in the status, and a refresh keeps it while the
// phase is unchanged.
func TestStatusAwaitingAuthorization(t *testing.T) {
	_, srv := newFakeAPI(t)
	ex := loadExample(t, mustKind(t, "MessagingConnection"), "minimal")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config:          providerBlock(srv.URL) + ex.hcl,
				ConfigVariables: exampleVars(ex),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ex.addrs[0], "status.phase", "AwaitingAuthorization"),
					resource.TestMatchResourceAttr(ex.addrs[0], "status.authorization_url", regexp.MustCompile(`^https://auth\.example\.test/slack\?c=`)),
				),
			},
			{
				Config:          providerBlock(srv.URL) + ex.hcl,
				ConfigVariables: exampleVars(ex),
				PlanOnly:        true,
				// the refresh must not lose the URL
			},
			{
				RefreshState: true,
				Check:        resource.TestMatchResourceAttr(ex.addrs[0], "status.authorization_url", regexp.MustCompile(`^https://`)),
			},
		},
	})
}

// A ServiceAccount's credentials are shown once, on create. They are kept in
// status.once, sensitive, across refreshes and updates, and an import does
// not have them.
func TestStatusShownOnce(t *testing.T) {
	f, srv := newFakeAPI(t)
	ex := loadExample(t, mustKind(t, "ServiceAccount"), "minimal")
	name := ex.names[0]
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock(srv.URL) + ex.hcl,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ex.addrs[0], "status.once.apiKey", "key-"+name),
					resource.TestCheckResourceAttr(ex.addrs[0], "status.once.secret", "shh-"+name),
					resource.TestCheckResourceAttr(ex.addrs[0], "status.extra.principal", "sa-"+name),
				),
			},
			{
				// A refresh reads an object whose answer no longer has them.
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ex.addrs[0], "status.once.apiKey", "key-"+name),
					func(*terraform.State) error {
						// The credentials were not asked for again: one PUT.
						if f.putCount() != 1 {
							return fmt.Errorf("%d PUTs", f.putCount())
						}
						return nil
					},
				),
			},
			{
				ResourceName:            ex.addrs[0],
				ImportState:             true,
				ImportStateId:           name,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

// Guardian refusing a SIEM write answers 502: the apply fails with the
// reason, says nothing was recorded, and leaves nothing in state.
func TestStatusFailedExternal(t *testing.T) {
	f, srv := newFakeAPI(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
resource "vulnetix_saved_search" "bad" {
  name = "refuse-this"
  text = "dns_events"
}
`,
			ExpectError: regexp.MustCompile(`(?s)external_failed.*Guardian\s+refused\s+the\s+SavedSearch.*nothing\s+was\s+recorded`),
		}},
	})
	if _, ok := f.objects["SavedSearch/refuse-this"]; ok {
		t.Error("a refused write was stored")
	}
}

// An owner-only kind refused for want of the owner role says what to do.
func TestOwnerOnlyForbidden(t *testing.T) {
	f, srv := newFakeAPI(t)
	f.forbidden["Organization"] = true
	var org kindSpec
	for _, k := range generatedKinds {
		if k.Kind == "Organization" {
			org = k
		}
	}
	if !org.OwnerOnly {
		t.Fatal("Organization is not marked owner-only")
	}
	ex := loadExample(t, org, "minimal")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config:      providerBlock(srv.URL) + ex.hcl,
			ExpectError: regexp.MustCompile(`(?s)HTTP\s+403.*can\s+only\s+be\s+applied\s+by\s+an\s+owner\s+or\s+admin`),
		}},
	})
}

// A member's credential is refused with 403 on apply, and the error says the
// principal must be an owner or admin. It holds for an ApiKey and a bearer token.
func TestMemberIsForbidden(t *testing.T) {
	for name, block := range map[string]func(url string) string{
		"apikey": func(url string) string {
			return `
provider "vulnetix" {
  endpoint = "` + url + `/v1"
  api_key  = "` + fakeMemberKey + `"
  org_id   = "` + fakeOrg + `"
}
`
		},
		"bearer": func(url string) string {
			return `
provider "vulnetix" {
  endpoint = "` + url + `/v1"
  token    = "` + fakeMemberToken + `"
}
`
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, srv := newFakeAPI(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				Steps: []resource.TestStep{{
					Config: block(srv.URL) + `
resource "vulnetix_quality_gate" "gate" {
  severity = "high"
}
`,
					ExpectError: regexp.MustCompile(`HTTP\s+403`),
				}},
			})
			if len(f.objects) != 0 {
				t.Errorf("a member wrote %d objects", len(f.objects))
			}
		})
	}
}

func mustKind(t *testing.T, kind string) kindSpec {
	t.Helper()
	for _, k := range generatedKinds {
		if k.Kind == kind {
			return k
		}
	}
	t.Fatalf("%s is not a served kind", kind)
	return kindSpec{}
}

// Only the four kinds that decide who can act for the tenant are owner-only.
func TestOwnerOnlyKinds(t *testing.T) {
	want := map[string]bool{"Organization": true, "Member": true, "ServiceAccount": true, "VaultEntry": true}
	for _, k := range generatedKinds {
		if k.OwnerOnly != want[k.Kind] {
			t.Errorf("%s: OwnerOnly = %v", k.Kind, k.OwnerOnly)
		}
	}
}
