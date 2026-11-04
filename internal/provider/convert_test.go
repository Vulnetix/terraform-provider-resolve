package provider

import (
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
