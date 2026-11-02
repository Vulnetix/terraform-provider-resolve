package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The generated examples (examples/resources/vulnetix_<kind>) drive a full
// lifecycle for every served kind, against the fake API in unit tests and
// against a real tenant in the acceptance tests (acceptance_test.go):
//
//  1. create from minimal.tf;
//  2. update to advanced.tf (every argument, secrets included);
//  3. import every object advanced.tf declares and compare with state;
//  4. destroy.
//
// Throughout, no write-only secret may reach the state.

var (
	reResource = regexp.MustCompile(`(?m)^resource "([a-z0-9_]+)" "([a-z0-9_]+)" \{`)
	reVariable = regexp.MustCompile(`(?m)^variable "([a-z0-9_]+)" \{`)
	reName     = regexp.MustCompile(`(?m)^  name\s+= "([^"]+)"`)
	reWO       = regexp.MustCompile(`(?m)^\s+([a-z0-9_]+_wo) `)
)

// example is one generated example file.
type example struct {
	hcl string
	// addrs are the resource addresses, names the object names, in order.
	addrs, names []string
	vars         []string
	// writeOnly are the *_wo arguments it sets.
	writeOnly []string
}

func loadExample(t *testing.T, k kindSpec, variant string) example {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "examples", "resources", "vulnetix_"+k.TypeName, variant+".tf"))
	if err != nil {
		t.Fatalf("%s: %v (run go generate ./...)", k.Kind, err)
	}
	e := example{hcl: string(b)}
	blocks := reResource.FindAllStringSubmatchIndex(e.hcl, -1)
	for i, m := range blocks {
		e.addrs = append(e.addrs, e.hcl[m[2]:m[3]]+"."+e.hcl[m[4]:m[5]])
		end := len(e.hcl)
		if i+1 < len(blocks) {
			end = blocks[i+1][0]
		}
		name := "default"
		if n := reName.FindStringSubmatch(e.hcl[m[0]:end]); n != nil {
			name = n[1]
		}
		e.names = append(e.names, name)
	}
	for _, m := range reVariable.FindAllStringSubmatch(e.hcl, -1) {
		e.vars = append(e.vars, m[1])
	}
	for _, m := range reWO.FindAllStringSubmatch(e.hcl, -1) {
		e.writeOnly = append(e.writeOnly, m[1])
	}
	if len(e.addrs) == 0 {
		t.Fatalf("%s %s.tf declares no resource", k.Kind, variant)
	}
	return e
}

// secretValue is what a test feeds an example's secret variable.
func secretValue(v string) string { return "s3cret-" + v }

func exampleVars(es ...example) config.Variables {
	out := config.Variables{}
	for _, e := range es {
		for _, v := range e.vars {
			out[v] = config.StringVariable(secretValue(v))
		}
	}
	return out
}

// noSecretsInState fails if any secret value, or any *_wo attribute, is in
// the state.
func noSecretsInState(es ...example) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, m := range s.Modules {
			for addr, r := range m.Resources {
				for k, v := range r.Primary.Attributes {
					if strings.HasSuffix(k, "_wo") && v != "" {
						return fmt.Errorf("%s: write-only %s is in state", addr, k)
					}
					for _, e := range es {
						for _, sv := range e.vars {
							if strings.Contains(v, secretValue(sv)) {
								return fmt.Errorf("%s: %s holds the secret %s", addr, k, sv)
							}
						}
					}
				}
			}
		}
		return nil
	}
}

// exampleSteps is the lifecycle for one kind. providerHCL configures the
// provider; check, when set, runs after the advanced apply.
func exampleSteps(t *testing.T, k kindSpec, providerHCL string, check resource.TestCheckFunc) []resource.TestStep {
	minimal, advanced := loadExample(t, k, "minimal"), loadExample(t, k, "advanced")
	advChecks := []resource.TestCheckFunc{noSecretsInState(advanced)}
	for i, addr := range advanced.addrs {
		advChecks = append(advChecks, resource.TestCheckResourceAttr(addr, "id", advanced.names[i]))
	}
	if check != nil {
		advChecks = append(advChecks, check)
	}
	steps := []resource.TestStep{
		{
			Config:          providerHCL + minimal.hcl,
			ConfigVariables: exampleVars(minimal),
			Check: resource.ComposeAggregateTestCheckFunc(
				noSecretsInState(minimal),
				resource.TestCheckResourceAttr(minimal.addrs[0], "id", minimal.names[0]),
			),
		},
		{
			Config:          providerHCL + advanced.hcl,
			ConfigVariables: exampleVars(advanced),
			Check:           resource.ComposeAggregateTestCheckFunc(advChecks...),
		},
	}
	for i, addr := range advanced.addrs {
		// status is what the server said at apply time (the once map, the
		// authorization URL); an import cannot know it.
		ignore := []string{"status"}
		// A write-only argument and its version are not in the server's answer, so
		// an import does not know them, nested or not.
		ignore = append(ignore, writeOnlyPaths(k.Fields, "")...)
		steps = append(steps, resource.TestStep{
			Config:                  providerHCL + advanced.hcl,
			ConfigVariables:         exampleVars(advanced),
			ResourceName:            addr,
			ImportState:             true,
			ImportStateId:           advanced.names[i],
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: ignore,
		})
	}
	return steps
}

// Every served kind's examples apply, update, import and destroy against
// the fake API, and the secrets they set are sent but never stored.
func TestExamplesLifecycle(t *testing.T) {
	for _, k := range generatedKinds {
		t.Run(k.Kind, func(t *testing.T) {
			f, srv := newFakeAPI(t)
			advanced := loadExample(t, k, "advanced")
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testFactories(),
				CheckDestroy: func(*terraform.State) error {
					for key := range f.objects {
						if strings.HasPrefix(key, k.Kind+"/") {
							return fmt.Errorf("%s left after destroy", key)
						}
					}
					return nil
				},
				Steps: exampleSteps(t, k, providerBlock(srv.URL), func(*terraform.State) error {
					// Every secret the advanced example sets reached the
					// server once.
					for _, v := range advanced.vars {
						sent := false
						for _, put := range f.puts {
							if strings.Contains(fmt.Sprint(put), secretValue(v)) {
								sent = true
							}
						}
						if !sent {
							return fmt.Errorf("secret %s was never sent", v)
						}
					}
					return nil
				}),
			})
		})
	}
}

// writeOnlyPaths lists the state paths of every write-only argument and its
// version, in the nested attributes too (oauth.client_secret_wo_version).
func writeOnlyPaths(fields []field, prefix string) []string {
	var out []string
	for _, f := range fields {
		switch {
		case f.WriteOnly || f.Local:
			out = append(out, prefix+f.Attr)
		case len(f.Children) > 0 && f.Kind == fieldObject:
			out = append(out, writeOnlyPaths(f.Children, prefix+f.Attr+".")...)
		}
	}
	return out
}
