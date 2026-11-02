package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	configschemas "github.com/vulnetix/config-schemas"
)

func TestSnake(t *testing.T) {
	for in, want := range map[string]string{
		"blockMalware":  "block_malware",
		"within30Days":  "within_30_days",
		"QualityGate":   "quality_gate",
		"URLPath":       "url_path",
		"maxMajorBump":  "max_major_bump",
		"scaAutofix":    "sca_autofix",
		"severity":      "severity",
		"AiFirewallKey": "ai_firewall_key",
	} {
		if got := snake(in); got != want {
			t.Errorf("snake(%q) = %q, want %q", in, got, want)
		}
	}
}

// The embedded release must generate, and every kind it holds must come out.
func TestGenerateEmbeddedKinds(t *testing.T) {
	set, err := loadSchemas(configschemas.FS(), "schemas")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(set.kinds(), ","), strings.Join(configschemas.Kinds(), ","); got != want {
		t.Fatalf("kinds = %s, want %s", got, want)
	}
	g := &generator{set: set, version: "test"}
	for _, k := range set.kinds() {
		if _, err := g.generate(k); err != nil {
			t.Fatalf("%s: %v", k, err)
		}
	}
}

func TestGenerateQualityGate(t *testing.T) {
	set, err := loadSchemas(configschemas.FS(), "schemas")
	if err != nil {
		t.Fatal(err)
	}
	o, err := (&generator{set: set, version: "test"}).generate("QualityGate")
	if err != nil {
		t.Fatal(err)
	}
	src := string(o.Source)
	for _, want := range []string{
		`Singleton:   true`,
		`"within_30_days": schema.StringAttribute{`,
		`{Attr: "within_30_days", JSON: "within30Days", Kind: fieldString}`,
		`{Attr: "cooldown_days", JSON: "cooldownDays", Kind: fieldInt}`,
		`stringvalidator.OneOf("poc", "active", "weaponized")`,
		`int64validator.AtLeast(0)`,
		`Schema default: ` + "`critical`",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated QualityGate lacks %s", want)
		}
	}
	if strings.Contains(src, "WriteOnly") {
		t.Error("QualityGate has no secrets but got a write-only attribute")
	}
}

func TestGenerateWidgetShapes(t *testing.T) {
	set, err := loadSchemas(os.DirFS("testdata/schemas"), ".")
	if err != nil {
		t.Fatal(err)
	}
	o, err := (&generator{set: set, version: "test"}).generate("Widget")
	if err != nil {
		t.Fatal(err)
	}
	src := string(o.Source)
	for _, want := range []string{
		`NamePattern: "^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$"`,
		`"color": schema.StringAttribute{`,
		`Required:            true`,
		`"token_wo": schema.StringAttribute{`,
		`WriteOnly:           true`,
		`Sensitive:           true`,
		`"token_wo_version": schema.Int64Attribute{`,
		`{Attr: "token_wo", JSON: "token", WriteOnly: true, Version: "token_wo_version"}`,
		`{Attr: "webhook_secret_wo", JSON: "webhookSecret", WriteOnly: true, Version: "webhook_secret_wo_version"}`,
		`"rules": schema.ListNestedAttribute{`,
		`"env": schema.MapAttribute{`,
		`"tags": schema.ListAttribute{`,
		`listvalidator.SizeAtMost(5)`,
		`CustomType:          jsontypes.NormalizedType{}`,
		`stringvalidator.RegexMatches(regexp.MustCompile("^pkg:[a-z][a-z0-9.+-]*/.+")`,
		`float64validator.AtLeast(0)`,
		`int64validator.AtMost(10)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated Widget lacks %s", want)
		}
	}
}

const commonOnly = `{"$defs": {
  "metadata": {"type": "object", "properties": {"name": {"type": "string"}}},
  "valueFrom": {"type": "object", "properties": {"env": {"type": "string"}}}
}}`

func kindSchema(spec string) string {
	return `{"properties": {"kind": {"const": "Bad"}, "metadata": {"$ref": "common.schema.json#/$defs/metadata"}, "spec": ` + spec + `}}`
}

func TestGenerateRefusals(t *testing.T) {
	for name, tc := range map[string]struct{ spec, want string }{
		"reserved name":    {`{"properties": {"name": {"type": "string"}}}`, "reserves"},
		"collision":        {`{"properties": {"fooBar": {"type": "string"}, "foo_bar": {"type": "string"}}}`, "both map to foo_bar"},
		"secret in a list": {`{"properties": {"keys": {"type": "array", "items": {"type": "object", "properties": {"k": {"$ref": "common.schema.json#/$defs/valueFrom"}}}}}}`, "cannot be write-only"},
		"dangling ref":     {`{"properties": {"x": {"$ref": "#/$defs/nope"}}}`, "not found"},
	} {
		t.Run(name, func(t *testing.T) {
			fsys := fstest.MapFS{
				"common.schema.json": {Data: []byte(commonOnly)},
				"Bad.schema.json":    {Data: []byte(kindSchema(tc.spec))},
			}
			set, err := loadSchemas(fsys, ".")
			if err != nil {
				t.Fatal(err)
			}
			_, err = (&generator{set: set}).generate("Bad")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want it to mention %q", err, tc.want)
			}
		})
	}
}

// The checked-in generated files must be what the generator writes now, so
// a schema bump without go generate fails CI.
func TestGeneratedFilesUpToDate(t *testing.T) {
	for _, tc := range []struct {
		args []string
	}{
		{[]string{"-version", "check"}},
		{[]string{"-version", "check", "-schemas", "testdata/schemas", "-suffix", "_test", "-var", "testKinds"}},
	} {
		dir := t.TempDir()
		if err := run(append(tc.args, "-out", dir)); err != nil {
			t.Fatal(err)
		}
		fresh, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		for _, f := range fresh {
			want, _ := os.ReadFile(f)
			got, err := os.ReadFile(filepath.Join("..", "..", "internal", "provider", filepath.Base(f)))
			if err != nil {
				t.Errorf("%s is not checked in: run go generate ./...", filepath.Base(f))
				continue
			}
			if !bytes.Equal(stripHeader(got), stripHeader(want)) {
				t.Errorf("%s is stale: run go generate ./...", filepath.Base(f))
			}
		}
	}
}

// stripHeader drops the first line, which names the schema version.
func stripHeader(b []byte) []byte {
	if i := bytes.IndexByte(b, '\n'); i >= 0 {
		return b[i+1:]
	}
	return b
}

// Every array of a served kind must be decided in collections.go from the
// engine adapter, and every entry there must still name an array.
func TestEveryArrayDecided(t *testing.T) {
	set, err := loadSchemas(configschemas.FS(), "schemas")
	if err != nil {
		t.Fatal(err)
	}
	g := &generator{set: set, version: "test"}
	for _, k := range set.kinds() {
		if _, err := g.generate(k); err != nil {
			t.Fatalf("%s: %v", k, err)
		}
	}
	for _, where := range g.undecided {
		// A kind with no engine adapter has no evidence to decide from; it
		// falls back to uniqueItems until its adapter is served.
		if kind, _, _ := strings.Cut(where, "."); !served[kind] || set.isValidateOnly(kind) {
			continue
		}
		t.Errorf("%s: no collections entry; decide list or set from the engine adapter", where)
	}
	for where := range collections {
		if !g.decided[where] {
			t.Errorf("collections entry %s names no array in the released schemas", where)
		}
	}
}

// Every served kind must have a released schema.
func TestServedKindsHaveSchemas(t *testing.T) {
	have := map[string]bool{}
	for _, k := range configschemas.Kinds() {
		have[k] = true
	}
	for k := range served {
		if !have[k] {
			t.Errorf("served kind %s has no schema in this config-schemas release", k)
		}
	}
}

func TestExampleHCL(t *testing.T) {
	set, err := loadSchemas(os.DirFS("testdata/schemas"), ".")
	if err != nil {
		t.Fatal(err)
	}
	o, err := (&generator{set: set, version: "test"}).generate("Widget")
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := writeExamples(out, "testdata/examples", o); err != nil {
		t.Fatal(err)
	}
	read := func(name string) string {
		b, err := os.ReadFile(filepath.Join(out, "vulnetix_widget", name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	minimal := read("minimal.tf")
	if !strings.Contains(minimal, `resource "vulnetix_widget" "w1" {`) || !strings.Contains(minimal, `color = "red"`) {
		t.Errorf("minimal.tf:\n%s", minimal)
	}
	adv := read("advanced.tf")
	for _, want := range []string{
		`resource "vulnetix_widget" "r_1st_widget" {`,
		`resource "vulnetix_widget" "r_1st_widget_2" {`,
		`variable "r_1st_widget_token" {`,
		`variable "r_1st_widget_nested_webhook_secret" {`,
		`ephemeral = true`,
		`token_wo         = var.r_1st_widget_token`,
		`token_wo_version = 1`,
		`webhook_secret_wo         = var.r_1st_widget_nested_webhook_secret`,
		`within_30_days            = "2026-12-31"`,
		`match = "$${literal} %%{ also }"`,
		`extra = jsonencode({`,
		`= ["a", "7"]`,
		`REGION = "ap-southeast-2"`,
		`"ticket/url" = "https://tickets.example/1"`,
		`purl   = "pkg:npm/left-pad"`,
		`ratio = 0.5`,
	} {
		if !strings.Contains(adv, want) {
			t.Errorf("advanced.tf lacks %s:\n%s", want, adv)
		}
	}
	if imp := read("import.sh"); !strings.Contains(imp, "terraform import vulnetix_widget.w1 w1") {
		t.Errorf("import.sh: %s", imp)
	}

	bad := t.TempDir()
	if err := os.MkdirAll(filepath.Join(bad, "Widget"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, v := range exampleVariants {
		_ = os.WriteFile(filepath.Join(bad, "Widget", v+".yaml"), []byte("kind: Widget\nmetadata: {name: x}\nspec: {nope: 1}\n"), 0o644)
	}
	if err := writeExamples(t.TempDir(), bad, o); err == nil || !strings.Contains(err.Error(), "spec.nope: no such argument") {
		t.Errorf("unknown key: %v", err)
	}
}

// The checked-in examples must be what the generator writes from the
// required config-schemas release, and nothing generated may be left over.
func TestGeneratedExamplesUpToDate(t *testing.T) {
	dir, ex := t.TempDir(), t.TempDir()
	if err := run([]string{"-version", "check", "-out", dir, "-examples-out", ex}); err != nil {
		t.Fatal(err)
	}
	checked := filepath.Join("..", "..", "examples", "resources")
	list := func(root string) map[string]bool {
		out := map[string]bool{}
		for _, pat := range []string{"*/minimal.tf", "*/advanced.tf", "*/import.sh"} {
			m, _ := filepath.Glob(filepath.Join(root, pat))
			for _, f := range m {
				rel, _ := filepath.Rel(root, f)
				out[rel] = true
			}
		}
		return out
	}
	fresh, have := list(ex), list(checked)
	if len(fresh) == 0 {
		t.Fatal("the generator wrote no examples")
	}
	for f := range fresh {
		want, _ := os.ReadFile(filepath.Join(ex, f))
		got, err := os.ReadFile(filepath.Join(checked, f))
		if err != nil {
			t.Errorf("examples/resources/%s is not checked in: run go generate ./...", f)
			continue
		}
		if !bytes.Equal(got, want) {
			t.Errorf("examples/resources/%s is stale: run go generate ./...", f)
		}
	}
	for f := range have {
		if !fresh[f] {
			t.Errorf("examples/resources/%s is no longer generated: run go generate ./...", f)
		}
	}
}

// A one-per-tenant kind can be declared once: further documents in an
// example become a commented-out alternative.
func TestExampleHCLSingletonAlternatives(t *testing.T) {
	set, err := loadSchemas(os.DirFS("testdata/schemas"), ".")
	if err != nil {
		t.Fatal(err)
	}
	o, err := (&generator{set: set, version: "test"}).generate("Widget")
	if err != nil {
		t.Fatal(err)
	}
	o.Singleton = true
	raw := []byte("kind: Widget\nmetadata: {name: default}\nspec: {color: red}\n---\nkind: Widget\nmetadata: {name: default}\nspec: {color: blue}\n")
	hcl, _, err := exampleHCL(o, "advanced", raw)
	if err != nil {
		t.Fatal(err)
	}
	s := string(hcl)
	if strings.Count(s, "\nresource ") != 1 || !strings.Contains(s, "# Or, instead") ||
		!strings.Contains(s, `# resource "vulnetix_widget" "default_2" {`) || !strings.Contains(s, `#   color = "blue"`) {
		t.Errorf("advanced.tf:\n%s", s)
	}
}
