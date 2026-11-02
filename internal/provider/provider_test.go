package provider

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
