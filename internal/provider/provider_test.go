package provider

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// The provider offers exactly the kinds the engine applies: an unserved
// kind's resource would fail every apply with unknown kind.
func TestResourcesAreServedKindsOnly(t *testing.T) {
	p := New("test")()
	got := map[string]bool{}
	for _, f := range p.Resources(context.Background()) {
		var resp resource.MetadataResponse
		f().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "vulnetix"}, &resp)
		got[resp.TypeName] = true
	}
	for _, k := range generatedKinds {
		if !got["vulnetix_"+k.TypeName] {
			t.Errorf("served kind %s has no resource", k.Kind)
		}
	}
	for _, k := range generatedKindsUnserved {
		if got["vulnetix_"+k.TypeName] {
			t.Errorf("unserved kind %s is offered as a resource", k.Kind)
		}
	}
	if len(got) != len(generatedKinds) {
		t.Errorf("%d resources for %d served kinds", len(got), len(generatedKinds))
	}
}

// The README and the registry's index page list every served resource and
// every data source.
func TestReadmeAndIndexListEverything(t *testing.T) {
	names := []string{}
	for _, k := range generatedKinds {
		names = append(names, "vulnetix_"+k.TypeName)
	}
	for _, f := range New("test")().DataSources(context.Background()) {
		var resp datasource.MetadataResponse
		f().Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "vulnetix"}, &resp)
		names = append(names, resp.TypeName)
	}
	for _, file := range []string{"../../README.md", "../../templates/index.md.tmpl"} {
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, n := range names {
			if !strings.Contains(string(b), "`"+n+"`") {
				t.Errorf("%s does not list %s", file, n)
			}
		}
	}
}

// The credential rules: an ApiKey with its organisation, or a bearer token,
// never both; none is a warning, not an error.
func TestCheckCredentials(t *testing.T) {
	for name, tc := range map[string]struct {
		apiKey, orgID, token string
		errors, warnings     int
	}{
		"api key and org":    {apiKey: "k", orgID: "o"},
		"bearer token":       {token: "t"},
		"bearer with org id": {token: "t", orgID: "o"},
		"api key and token":  {apiKey: "k", orgID: "o", token: "t", errors: 1},
		"api key only":       {apiKey: "k", errors: 1},
		"api key with token": {apiKey: "k", token: "t", errors: 1},
		"org id only":        {orgID: "o", errors: 1},
		"nothing configured": {warnings: 1},
	} {
		t.Run(name, func(t *testing.T) {
			d := checkCredentials(tc.apiKey, tc.orgID, tc.token)
			if d.ErrorsCount() != tc.errors || d.WarningsCount() != tc.warnings {
				t.Errorf("%d errors and %d warnings, want %d and %d: %v", d.ErrorsCount(), d.WarningsCount(), tc.errors, tc.warnings, d)
			}
		})
	}
}

// The credential arguments say who may use them.
func TestCredentialDescriptions(t *testing.T) {
	var resp provider.SchemaResponse
	(&vulnetixProvider{}).Schema(context.Background(), provider.SchemaRequest{}, &resp)
	for _, name := range []string{"api_key", "org_id", "token"} {
		if resp.Schema.Attributes[name] == nil {
			t.Fatalf("no %s argument", name)
		}
	}
	for _, name := range []string{"api_key", "token"} {
		a := resp.Schema.Attributes[name]
		if !a.IsSensitive() {
			t.Errorf("%s is not sensitive", name)
		}
		d := a.GetMarkdownDescription()
		for _, want := range []string{"owner or admin", "service account", "member", "403"} {
			if !strings.Contains(d, want) {
				t.Errorf("%s description does not mention %q: %s", name, want, d)
			}
		}
	}
}
