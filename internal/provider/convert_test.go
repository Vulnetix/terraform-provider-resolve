package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// A server copy that repeats a set element exactly (PackageFirewallPolicy
// mirrors held one homebrew-bottle mirror twice) must read as one element:
// Terraform refuses a set value with duplicates, which made the object
// unimportable.
func TestSetReadDropsExactDuplicates(t *testing.T) {
	mirror := field{Attr: "mirrors", JSON: "mirrors", Kind: fieldSetObject, Children: []field{
		{Attr: "ecosystem", JSON: "ecosystem", Kind: fieldString},
		{Attr: "url", JSON: "url", Kind: fieldString},
		{Attr: "priority", JSON: "priority", Kind: fieldInt},
	}}
	ot := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"ecosystem": tftypes.String, "url": tftypes.String, "priority": tftypes.Number,
	}}
	st := tftypes.Set{ElementType: ot}
	dup := map[string]any{"ecosystem": "homebrew-bottle", "url": "https://ghcr.io", "priority": float64(0)}
	other := map[string]any{"ecosystem": "npm", "url": "https://registry.npmjs.org", "priority": float64(0)}

	v, err := jsonToValue([]any{dup, other, dup}, st, mirror, tftypes.NewValue(st, nil), true)
	if err != nil {
		t.Fatal(err)
	}
	var elems []tftypes.Value
	if err := v.As(&elems); err != nil {
		t.Fatal(err)
	}
	if len(elems) != 2 {
		t.Fatalf("got %d elements, want 2 (the repeat dropped)", len(elems))
	}

	// Elements that differ in any field are distinct and all kept.
	near := map[string]any{"ecosystem": "homebrew-bottle", "url": "https://ghcr.io", "priority": float64(1)}
	v, err = jsonToValue([]any{dup, near}, st, mirror, tftypes.NewValue(st, nil), true)
	if err != nil {
		t.Fatal(err)
	}
	if err := v.As(&elems); err != nil {
		t.Fatal(err)
	}
	if len(elems) != 2 {
		t.Fatalf("got %d elements, want 2 (different priority is a different element)", len(elems))
	}
}

// A list keeps repeats: its order and length are the document's.
func TestListReadKeepsRepeats(t *testing.T) {
	tags := field{Attr: "tags", JSON: "tags", Kind: fieldList, Elem: fieldString}
	lt := tftypes.List{ElementType: tftypes.String}
	v, err := jsonToValue([]any{"a", "a"}, lt, tags, tftypes.NewValue(lt, nil), true)
	if err != nil {
		t.Fatal(err)
	}
	var elems []tftypes.Value
	if err := v.As(&elems); err != nil {
		t.Fatal(err)
	}
	if len(elems) != 2 {
		t.Fatalf("got %d elements, want 2", len(elems))
	}
}

// An import reads a value the server filled in from the schema default as
// unset, so the first plan after it does not change every default the
// configuration left out (PackageList entries came back with scope "full").
// A value that is not the default is kept, and a normal read, which has a
// prior to go by, is not affected.
func TestImportReadsDefaultsAsUnset(t *testing.T) {
	fields := []field{
		{Attr: "scope", JSON: "scope", Kind: fieldString, Default: `"full"`},
		{Attr: "range_ms", JSON: "rangeMs", Kind: fieldInt, Default: "86400000"},
		{Attr: "enabled", JSON: "enabled", Kind: fieldBool, Default: "true"},
		{Attr: "purl", JSON: "purl", Kind: fieldString},
	}
	ot := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"scope": tftypes.String, "range_ms": tftypes.Number, "enabled": tftypes.Bool, "purl": tftypes.String,
	}}
	server := map[string]any{"scope": "full", "rangeMs": float64(86400000), "enabled": true, "purl": "pkg:npm/a"}

	vals, err := fromSpec(server, ot, tftypes.NewValue(ot, nil), fields, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range []string{"scope", "range_ms", "enabled"} {
		if !vals[a].IsNull() {
			t.Errorf("import kept the default %s = %v", a, vals[a])
		}
	}
	if vals["purl"].IsNull() {
		t.Error("import dropped purl, which has no default")
	}

	// Not the default: kept.
	server["scope"], server["enabled"] = "bad_actor", false
	if vals, err = fromSpec(server, ot, tftypes.NewValue(ot, nil), fields, true); err != nil {
		t.Fatal(err)
	}
	if vals["scope"].IsNull() || vals["enabled"].IsNull() {
		t.Errorf("import dropped a value that is not the default: %v %v", vals["scope"], vals["enabled"])
	}

	// A normal read keeps what the configuration set, default or not.
	prior := tftypes.NewValue(ot, map[string]tftypes.Value{
		"scope": tftypes.NewValue(tftypes.String, "full"), "range_ms": tftypes.NewValue(tftypes.Number, nil),
		"enabled": tftypes.NewValue(tftypes.Bool, nil), "purl": tftypes.NewValue(tftypes.String, "pkg:npm/a"),
	})
	server["scope"] = "full"
	if vals, err = fromSpec(server, ot, prior, fields, false); err != nil {
		t.Fatal(err)
	}
	if vals["scope"].IsNull() {
		t.Error("a read dropped a default the configuration wrote out")
	}
}

func TestIsDefault(t *testing.T) {
	for _, tc := range []struct {
		raw  any
		def  string
		want bool
	}{
		{"full", `"full"`, true},
		{"bad_actor", `"full"`, false},
		{float64(86400000), "86400000", true},
		{json.Number("8.64e7"), "86400000", true},
		{json.Number("86400001"), "86400000", false},
		{true, "true", true},
		{false, "true", false},
		{"true", "true", false},
		{"x", "", false},
	} {
		if got := isDefault(tc.raw, tc.def); got != tc.want {
			t.Errorf("isDefault(%#v, %s) = %v, want %v", tc.raw, tc.def, got, tc.want)
		}
	}
}
