package provider

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vulnetix/terraform-provider-vulnetix/internal/client"
)

// The data source examples are written by hand: unlike the resources, no
// config-schemas example describes them. So every data source must have one
// (examples/data-sources/vulnetix_<name>/data-source.tf, which tfplugindocs
// renders), and each must plan and apply against the fake API, so an
// example cannot drift from the schema.
func TestDataSourceExamples(t *testing.T) {
	f, srv := newFakeAPI(t)
	seedFeeds(f)
	base := "https://feeds.example.test/feeds/org/f4"
	f.feeds = append(f.feeds,
		client.Feed{UUID: "f4", Name: "edge-ips", Status: "PUBLISHED", FeedType: "SURICATA", SchemaVersion: 1, PublicBaseURL: &base, CreatedAt: 4},
		client.Feed{UUID: "f5", Name: "edge-waf", Status: "PUBLISHED", FeedType: "MODSECURITY", SchemaVersion: 1, CreatedAt: 5},
	)
	f.entries["f5"] = f.entries["f1"]
	f.models = []client.AiModel{{UUID: "m1", Slug: "claude-sonnet-4", Name: "Claude Sonnet 4", ProviderSlug: "anthropic", IsActive: true}}
	f.vex = []map[string]any{{"uuid": "s1", "vulnerabilityId": "CVE-2026-1", "status": "not_affected", "justification": "component_not_present",
		"aliases": []string{}, "purls": []string{"pkg:npm/lodash@4.17.21"}, "products": []any{}, "timestamp": 1}}
	f.repos = []client.GitHubRepo{{InstallationID: 1, AccountLogin: "acme", OwnerLogin: "acme", Name: "api", FullName: "acme/api", Host: "github.com"}}

	for _, newDS := range New("test")().DataSources(context.Background()) {
		var md datasource.MetadataResponse
		newDS().Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "vulnetix"}, &md)
		t.Run(md.TypeName, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join("..", "..", "examples", "data-sources", md.TypeName, "data-source.tf"))
			if err != nil {
				t.Fatalf("no example for %s: %v", md.TypeName, err)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				Steps:                    []resource.TestStep{{Config: providerBlock(srv.URL) + string(b)}},
			})
		})
	}
}
