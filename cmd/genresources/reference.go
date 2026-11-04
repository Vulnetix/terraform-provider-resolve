package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// The reference example (examples/resources/vulnetix_<kind>/reference.tf)
// declares every argument that has a schema default, set to that default,
// plus the arguments the kind requires, taken from the minimal example. A
// plan of it against a tenant only sets values the tenant left unset, and
// changes none it chose: since 1.3.3 an import reads a value equal to its
// default as unset, so declaring the default is the only diff it can cause.
//
// The Vulnetix/platform-iac-examples repository publishes these files.

// specDefaults returns the default of every property of an object node that
// has one, recursing into objects. An object is only included when its
// defaults satisfy its own required properties, so the result is a valid
// document fragment. Secrets have no defaults and are skipped.
func (g *generator) specDefaults(raw *node, present map[string]any) (map[string]any, error) {
	out := map[string]any{}
	n, _, err := g.set.flatten(raw)
	if err != nil || n == nil {
		return out, err
	}
	cond, err := g.conditional(raw)
	if err != nil {
		return nil, err
	}
	for name, p := range n.Properties {
		if p == nil || p.boolean != nil {
			continue
		}
		if _, have := present[name]; cond[name] && !have {
			// Allowed only in some documents (a oneOf branch, or a property an
			// if/then/else forbids): only the minimal example decides.
			continue
		}
		secret, err := g.set.isValueFrom(p)
		if err != nil {
			return nil, err
		}
		if secret {
			continue
		}
		r, _, err := g.set.flatten(p)
		if err != nil {
			return nil, err
		}
		if r == nil {
			continue
		}
		if r.Default != nil {
			out[name] = r.Default
			continue
		}
		if isObject(r) && len(r.Properties) > 0 {
			subPresent, _ := present[name].(map[string]any)
			sub, err := g.specDefaults(p, subPresent)
			if err != nil {
				return nil, err
			}
			if len(sub) == 0 {
				continue
			}
			complete := true
			for _, req := range r.Required {
				if _, ok := sub[req]; !ok {
					complete = false
					break
				}
			}
			if complete {
				out[name] = sub
			}
		}
	}
	return out, nil
}

// overlayDefaults adds the defaults to a spec without changing a value it
// already has: the minimal example's required values win.
func overlayDefaults(spec, defaults map[string]any) map[string]any {
	if spec == nil {
		spec = map[string]any{}
	}
	for k, d := range defaults {
		cur, ok := spec[k]
		if !ok {
			spec[k] = d
			continue
		}
		cm, cok := cur.(map[string]any)
		dm, dok := d.(map[string]any)
		if cok && dok {
			spec[k] = overlayDefaults(cm, dm)
		}
	}
	return spec
}

// referenceYAML builds the reference document from the minimal example's
// first document and the kind's defaults.
func referenceYAML(k *kindOut, minimal []byte) ([]byte, error) {
	var doc map[string]any
	if err := yaml.NewDecoder(bytes.NewReader(minimal)).Decode(&doc); err != nil {
		return nil, fmt.Errorf("minimal example: %w", err)
	}
	spec, _ := doc["spec"].(map[string]any)
	doc["spec"] = overlayDefaults(spec, k.Defaults)
	if md, ok := doc["metadata"].(map[string]any); ok {
		delete(md, "labels")
		delete(md, "annotations")
	}
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return b.Bytes(), enc.Close()
}

// writeReference writes one kind's reference.tf.
func writeReference(out, src string, k *kindOut) error {
	minimal, err := os.ReadFile(filepath.Join(src, k.Kind, "minimal.yaml"))
	if err != nil {
		return fmt.Errorf("%s: minimal example: %w", k.Kind, err)
	}
	raw, err := referenceYAML(k, minimal)
	if err != nil {
		return fmt.Errorf("%s: %w", k.Kind, err)
	}
	dir := filepath.Join(out, "vulnetix_"+k.TypeName)
	head := fmt.Sprintf("# Every field with a schema default, set to that default, plus the required\n# fields from config-schemas examples/%s/minimal.yaml. reference.tf is the same\n# document as Terraform.\n", k.Kind)
	if err := os.WriteFile(filepath.Join(dir, "reference.yaml"), append([]byte(head), raw...), 0o644); err != nil {
		return err
	}
	hcl, _, err := exampleHCL(k, "reference", raw)
	if err != nil {
		return fmt.Errorf("%s: reference: %w", k.Kind, err)
	}
	return os.WriteFile(filepath.Join(dir, "reference.tf"), hcl, 0o644)
}

// conditional returns the properties of an object schema that a document may
// not always carry: those that exist only in a oneOf or anyOf branch, and
// those a not/required (in an if/then/else or allOf) forbids in some case.
func (g *generator) conditional(raw *node) (map[string]bool, error) {
	out := map[string]bool{}
	r, _, err := g.set.resolve(raw)
	if err != nil || r == nil {
		return out, err
	}
	base := map[string]bool{}
	var walk func(n *node, depth int) error
	walk = func(n *node, depth int) error {
		if n == nil || depth > 8 {
			return nil
		}
		// A ref this walk cannot follow (a pointer into a $defs map) only
		// means fewer names are known to be conditional: skip it.
		n, _, err := g.set.resolve(n)
		if err != nil || n == nil {
			return nil
		}
		if n.Not != nil {
			nn, _, err := g.set.resolve(n.Not)
			if err != nil || nn == nil {
				return nil
			}
			for _, req := range nn.Required {
				out[req] = true
			}
			for k := range nn.Properties {
				out[k] = true
			}
		}
		for _, c := range []*node{n.If, n.Then, n.Else} {
			if err := walk(c, depth+1); err != nil {
				return err
			}
		}
		for _, c := range n.AllOf {
			if err := walk(c, depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(r, 0); err != nil {
		return nil, err
	}
	for k := range r.Properties {
		base[k] = true
	}
	for _, c := range r.AllOf {
		cr, _, err := g.set.flatten(c)
		if err != nil {
			return nil, err
		}
		if cr != nil {
			for k := range cr.Properties {
				base[k] = true
			}
		}
	}
	for _, list := range [][]*node{r.OneOf, r.AnyOf} {
		for _, b := range list {
			br, _, err := g.set.flatten(b)
			if err != nil {
				return nil, err
			}
			if br == nil {
				continue
			}
			for k := range br.Properties {
				if !base[k] {
					out[k] = true
				}
			}
			for _, req := range br.Required {
				if !base[req] {
					out[req] = true
				}
			}
		}
	}
	return out, nil
}
