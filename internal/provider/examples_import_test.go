package provider

import (
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// An import reads a value equal to its schema default as unset (fromSpec),
// while advanced.tf writes some defaults out. The import steps therefore
// compare state and import with defaults allowed to be absent on either side,
// and with list and set positions ignored, so a set element that differs only
// in a default does not shift every index after it.

// defaultPath is a scalar with a schema default, as flatmap path segments in
// which "*" is any list or set position or map key.
type defaultPath struct {
	segs []string
	def  string
}

func defaultPaths(fields []field, prefix []string) []defaultPath {
	var out []defaultPath
	for _, f := range fields {
		at := append(slices.Clone(prefix), f.Attr)
		switch f.Kind {
		case fieldObject:
			out = append(out, defaultPaths(f.Children, at)...)
		case fieldListObject, fieldSetObject, fieldMapObject:
			out = append(out, defaultPaths(f.Children, append(at, "*"))...)
		default:
			if f.Default != "" {
				out = append(out, defaultPath{segs: at, def: f.Default})
			}
		}
	}
	return out
}

var reIndex = regexp.MustCompile(`^[0-9]+$`)

// flatForImport is a flatmap state as path (positions as "*") to its sorted
// values, leaving out the ignored paths and any value equal to its schema
// default.
func flatForImport(attrs map[string]string, ignore []string, defaults []defaultPath) map[string][]string {
	out := map[string][]string{}
	for k, v := range attrs {
		if slices.ContainsFunc(ignore, func(p string) bool { return strings.HasPrefix(k, p) }) {
			continue
		}
		segs := strings.Split(k, ".")
		for i, s := range segs {
			if reIndex.MatchString(s) {
				segs[i] = "*"
			}
		}
		if flatIsDefault(segs, v, defaults) {
			continue
		}
		key := strings.Join(segs, ".")
		out[key] = append(out[key], v)
	}
	for k := range out {
		slices.Sort(out[k])
	}
	return out
}

func flatIsDefault(segs []string, v string, defaults []defaultPath) bool {
	for _, d := range defaults {
		if len(d.segs) != len(segs) {
			continue
		}
		match := true
		for i := range segs {
			if d.segs[i] != "*" && d.segs[i] != segs[i] {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		var want any
		dec := json.NewDecoder(strings.NewReader(d.def))
		dec.UseNumber()
		if dec.Decode(&want) != nil {
			continue
		}
		switch w := want.(type) {
		case string:
			return v == w
		case bool:
			return v == fmt.Sprint(w)
		case json.Number:
			a, _, errA := big.ParseFloat(v, 10, 512, big.ToNearestEven)
			b, _, errB := big.ParseFloat(w.String(), 10, 512, big.ToNearestEven)
			return errA == nil && errB == nil && a.Cmp(b) == 0
		}
	}
	return false
}

func flatDiff(want, got map[string][]string) string {
	var b strings.Builder
	keys := slices.Sorted(maps.Keys(want))
	for k := range got {
		if _, ok := want[k]; !ok {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		if !slices.Equal(want[k], got[k]) {
			fmt.Fprintf(&b, "  %s: state %q, import %q\n", k, want[k], got[k])
		}
	}
	return b.String()
}

// The import comparison drops defaults wherever a set element puts them, and
// still sees a real difference.
func TestFlatForImport(t *testing.T) {
	defaults := []defaultPath{{segs: []string{"entries", "*", "scope"}, def: `"full"`}, {segs: []string{"range_ms"}, def: "86400000"}}
	state := map[string]string{
		"entries.#": "2", "entries.0.purl": "a", "entries.0.scope": "full", "entries.1.purl": "b", "entries.1.scope": "bad_actor",
		"range_ms": "86400000", "status.%": "1",
	}
	imported := map[string]string{
		"entries.#": "2", "entries.0.purl": "b", "entries.0.scope": "bad_actor", "entries.1.purl": "a",
		"range_ms": "8.64e+07",
	}
	ignore := []string{"status"}
	if !maps.EqualFunc(flatForImport(state, ignore, defaults), flatForImport(imported, ignore, defaults), slices.Equal[[]string]) {
		t.Fatalf("a default and a set's order made a difference:\n%s", flatDiff(flatForImport(state, ignore, defaults), flatForImport(imported, ignore, defaults)))
	}
	imported["entries.0.scope"] = "full"
	if maps.EqualFunc(flatForImport(state, ignore, defaults), flatForImport(imported, ignore, defaults), slices.Equal[[]string]) {
		t.Fatal("a changed value that is not a default went unnoticed")
	}
}
