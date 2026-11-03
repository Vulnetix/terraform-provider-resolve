package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vulnetix/terraform-provider-resolve/internal/client"
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
	for i := 0; i < 1200; i++ {
		f.vex = append(f.vex, map[string]any{"uuid": fmt.Sprintf("s%d", i), "documentUuid": "d1", "documentId": "urn:x",
			"author": "sec@acme.example", "vulnerabilityId": fmt.Sprintf("CVE-2026-%04d", i), "vulnerabilityName": "lib",
			"aliases": []string{"GHSA-x"}, "status": "not_affected", "justification": "component_not_present",
			"products": []map[string]any{{"@id": "pkg:npm/lib@1.0.0"}}, "purls": []string{"pkg:npm/lib@1.0.0"},
			"timestamp": 1767225600000 + i, "grounded": true, "sourceFormat": "openvex", "repositoryFullName": "acme/api"})
	}
	f.vex = append(f.vex, map[string]any{"uuid": "inv", "vulnerabilityId": "CVE-2026-9999", "status": "under_investigation",
		"aliases": []string{}, "purls": []string{}, "timestamp": 1})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
data "vulnetix_vex_statements" "default_limit" {}
data "vulnetix_vex_statements" "few" { limit = 5 }
data "vulnetix_vex_statements" "inv" { status = "under_investigation" }
data "vulnetix_vex_statements" "one" { vulnerability_id = "cve-2026-0007" }
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.default_limit", "statements.#", "1000"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.few", "statements.#", "5"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.inv", "statements.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.inv", "statements.0.products", "null"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.vulnerability_id", "CVE-2026-0007"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.purls.0", "pkg:npm/lib@1.0.0"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.aliases.0", "GHSA-x"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.products", `[{"@id":"pkg:npm/lib@1.0.0"}]`),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.repository", "acme/api"),
				resource.TestCheckResourceAttr("data.vulnetix_vex_statements.one", "statements.0.timestamp", "1767225600007"),
			),
		}},
	})
}

func TestGitHubInstallationReposDataSource(t *testing.T) {
	f, srv := newFakeAPI(t)
	for i := 0; i < 600; i++ {
		f.repos = append(f.repos, client.GitHubRepo{InstallationID: 1, AccountLogin: "acme", RepoID: int64(i),
			OwnerLogin: "acme", Name: fmt.Sprintf("svc-%03d", i), FullName: fmt.Sprintf("acme/svc-%03d", i),
			Private: true, Archived: i%100 == 0, Host: "github.com", UpdatedAt: 5})
	}
	f.repos = append(f.repos, client.GitHubRepo{InstallationID: 2, AccountLogin: "other", OwnerLogin: "other", Name: "web", FullName: "other/web", Host: "github.com"})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testFactories(),
		Steps: []resource.TestStep{{
			Config: providerBlock(srv.URL) + `
data "vulnetix_github_installation_repos" "all" {}
data "vulnetix_github_installation_repos" "archived" { include_archived = true }
data "vulnetix_github_installation_repos" "other" { installation_id = 2 }
data "vulnetix_github_installation_repos" "q" { query = "SVC-01" }
`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.all", "repos.#", "595"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.archived", "repos.#", "601"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.other", "full_names.#", "1"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.other", "full_names.0", "other/web"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.q", "repos.#", "10"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.q", "repos.0.private", "true"),
				resource.TestCheckResourceAttr("data.vulnetix_github_installation_repos.q", "repos.0.host", "github.com"),
			),
		}},
	})
}
