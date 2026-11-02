package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

func TestKindsDataSource(t *testing.T) {
	_, srv := newFakeAPI(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `data "vulnetix_kinds" "all" {}`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "api_version", "vulnetix.com/v1"),
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "schema_version", "v1.0.1"),
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "kinds.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "kinds.0.kind", "QualityGate"),
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "kinds.0.schema", `{"title":"QualityGate"}`),
				resource.TestCheckResourceAttr("data.vulnetix_kinds.all", "schema_kinds.0", "QualityGate"),
			),
		}},
	})
}

func seedFeeds(f *fakeAPI) {
	base := "https://feeds.example.test/feeds/org/f1"
	pub := int64(1700000000000)
	f.feeds = []client.Feed{
		{UUID: "f1", Name: "edge waf", Status: "PUBLISHED", FeedType: "MODSECURITY", SchemaVersion: 1, RuleCount: 3, PublicBaseURL: &base, PublishedAt: &pub, CreatedAt: 1},
		{UUID: "f2", Name: "dup", Status: "DRAFT", FeedType: "SIGMA", SchemaVersion: 1, CreatedAt: 2},
		{UUID: "f3", Name: "dup", Status: "DRAFT", FeedType: "YARA", SchemaVersion: 1, CreatedAt: 3},
	}
	f.entries["f1"] = []client.FeedEntry{
		{ID: 1, RuleTable: "virtual_patch", RuleType: "MODSECURITY", RuleUUID: "vp-1", RuleName: "CVE-2024-1 block", MatchedCve: "CVE-2024-1", AddedAt: 5},
		{ID: 2, RuleTable: "virtual_patch", RuleType: "MODSECURITY", RuleUUID: "vp-2", RuleName: "off", Disabled: true, AddedAt: 6},
		{ID: 3, RuleTable: "snort_rule", RuleType: "SNORT", RuleUUID: "s-1", RuleName: "catalogue", AddedAt: 7},
	}
}

func TestDetectionRuleFeedDataSource(t *testing.T) {
	f, srv := newFakeAPI(t)
	seedFeeds(f)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config:      providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "d" {}`,
				ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured`),
			},
			{
				Config:      providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "d" { name = "dup" }`,
				ExpectError: regexp.MustCompile(`2 detection rule feeds are named "dup"`),
			},
			{
				Config:      providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "d" { uuid = "nope" }`,
				ExpectError: regexp.MustCompile(`feed nope not found`),
			},
			{
				Config: providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "waf" { name = "edge waf" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vulnetix_detection_rule_feed.waf", "uuid", "f1"),
					resource.TestCheckResourceAttr("data.vulnetix_detection_rule_feed.waf", "feed_type", "MODSECURITY"),
					resource.TestCheckResourceAttr("data.vulnetix_detection_rule_feed.waf", "manifest_url", "https://feeds.example.test/feeds/org/f1/manifest.json"),
					resource.TestCheckResourceAttr("data.vulnetix_detection_rule_feed.waf", "entries.#", "3"),
					resource.TestCheckNoResourceAttr("data.vulnetix_detection_rule_feed.waf", "revoked_at"),
				),
			},
			{
				Config: providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "d" { uuid = "f2" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vulnetix_detection_rule_feed.d", "status", "DRAFT"),
					resource.TestCheckNoResourceAttr("data.vulnetix_detection_rule_feed.d", "manifest_url"),
				),
			},
		},
	})
}

func TestVirtualPatchesDataSource(t *testing.T) {
	f, srv := newFakeAPI(t)
	seedFeeds(f)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerBlock(srv.URL) + `data "vulnetix_virtual_patches" "vp" { feed_uuid = "f1" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.vulnetix_virtual_patches.vp", "feed_name", "edge waf"),
					resource.TestCheckResourceAttr("data.vulnetix_virtual_patches.vp", "patches.#", "1"),
					resource.TestCheckResourceAttr("data.vulnetix_virtual_patches.vp", "patches.0.matched_cve", "CVE-2024-1"),
				),
			},
			{
				Config: providerBlock(srv.URL) + `data "vulnetix_virtual_patches" "vp" {
  feed_name        = "edge waf"
  include_disabled = true
}`,
				Check: resource.TestCheckResourceAttr("data.vulnetix_virtual_patches.vp", "patches.#", "2"),
			},
		},
	})
}

func TestFeedsPlanLocked(t *testing.T) {
	f, srv := newFakeAPI(t)
	f.locked = true
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config:      providerBlock(srv.URL) + `data "vulnetix_detection_rule_feed" "d" { name = "x" }`,
			ExpectError: regexp.MustCompile(`need a Pro plan or higher \(requires pro\)`),
		}},
	})
}

func TestAiModelsDataSource(t *testing.T) {
	f, srv := newFakeAPI(t)
	for i := 0; i < 130; i++ {
		f.models = append(f.models, client.AiModel{UUID: fmt.Sprintf("m%d", i), Slug: fmt.Sprintf("model-%d", i),
			Name: "Model", ProviderUUID: "p1", ProviderSlug: "acme", IsActive: true})
	}
	f.models = append(f.models, client.AiModel{UUID: "c1", Slug: "claude-x", Name: "Other", ProviderUUID: "p2", ProviderSlug: "other"})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
data "vulnetix_ai_models" "all" {}
data "vulnetix_ai_models" "p2" { provider_uuid = "p2" }
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_ai_models.all", "models.#", "131"),
				resource.TestCheckResourceAttr("data.vulnetix_ai_models.p2", "models.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_ai_models.p2", "models.0.provider_slug", "other"),
				resource.TestCheckResourceAttr("data.vulnetix_ai_models.p2", "models.0.is_active", "false"),
			),
		}},
	})
}

func TestVexStatementsDataSource(t *testing.T) {
	f, srv := newFakeAPI(t)
	for i := 0; i < 120; i++ {
		f.vex = append(f.vex, map[string]any{"kind": "vex", "at": 1767225600000 + i, "subjectKey": fmt.Sprintf("s%d", i),
			"findingUuid": "", "subjectLabel": fmt.Sprintf("CVE-2026-%04d · lib", i), "outcome": "not_affected",
			"justification": "component_not_present", "detail": "", "repoName": "", "actor": "sec@acme.example", "automated": false})
	}
	f.vex = append(f.vex, map[string]any{"subjectKey": "fx", "at": 1, "subjectLabel": "GHSA-xxxx", "outcome": "fixed", "actor": "a"})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
data "vulnetix_vex_statements" "all" {}
data "vulnetix_vex_statements" "few" { limit = 5 }
data "vulnetix_vex_statements" "fixed" { status = "fixed" }
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.all", "statements.#", "121"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.few", "statements.#", "5"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.all", "statements.0.vulnerability_id", "CVE-2026-0000"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.all", "statements.0.vulnerability_name", "lib"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.all", "statements.0.timestamp", "1767225600000"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.fixed", "statements.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.fixed", "statements.0.vulnerability_id", "GHSA-xxxx"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.fixed", "statements.0.vulnerability_name", ""),
			),
		}},
	})
}
