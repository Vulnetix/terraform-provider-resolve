package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// node is the subset of JSON Schema 2020-12 the kind schemas use.
type node struct {
	Ref         string           `json:"$ref"`
	Type        json.RawMessage  `json:"type"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Enum        []any            `json:"enum"`
	Const       any              `json:"const"`
	Default     any              `json:"default"`
	Pattern     string           `json:"pattern"`
	Format      string           `json:"format"`
	MinLength   *int64           `json:"minLength"`
	MaxLength   *int64           `json:"maxLength"`
	Minimum     *json.Number     `json:"minimum"`
	Maximum     *json.Number     `json:"maximum"`
	MinItems    *int64           `json:"minItems"`
	MaxItems    *int64           `json:"maxItems"`
	Properties  map[string]*node `json:"properties"`
	Required    []string         `json:"required"`
	AddlProps   json.RawMessage  `json:"additionalProperties"`
	Items       *node            `json:"items"`
	AllOf       []*node          `json:"allOf"`
	OneOf       []*node          `json:"oneOf"`
	AnyOf       []*node          `json:"anyOf"`
	Defs        map[string]*node `json:"$defs"`

	// file is the schema file this node was read from, for resolving its
	// local "#/..." refs.
	file string
}

// types lists the node's declared types.
func (n *node) types() []string {
	if len(n.Type) == 0 {
		return nil
	}
	var one string
	if json.Unmarshal(n.Type, &one) == nil {
		return []string{one}
	}
	var many []string
	_ = json.Unmarshal(n.Type, &many)
	return many
}

// hasType reports whether the node is (only, ignoring null) of type t.
func (n *node) hasType(t string) bool {
	for _, x := range n.types() {
		if x == t {
			return true
		}
	}
	return false
}

// additional returns additionalProperties as a schema, or nil when it is
// absent or a boolean.
func (n *node) additional() *node {
	if len(n.AddlProps) == 0 || n.AddlProps[0] != '{' {
		return nil
	}
	var a node
	if json.Unmarshal(n.AddlProps, &a) != nil {
		return nil
	}
	a.file = n.file
	setFile(&a, n.file)
	return &a
}

func (n *node) requires(name string) bool {
	for _, r := range n.Required {
		if r == name {
			return true
		}
	}
	return false
}

func setFile(n *node, file string) {
	if n == nil {
		return
	}
	n.file = file
	for _, c := range n.Properties {
		setFile(c, file)
	}
	for _, c := range n.Defs {
		setFile(c, file)
	}
	for _, list := range [][]*node{n.AllOf, n.OneOf, n.AnyOf} {
		for _, c := range list {
			setFile(c, file)
		}
	}
	setFile(n.Items, file)
}

// schemaSet is every schema file in a directory, keyed by file name.
type schemaSet struct {
	files map[string]*node
}

func loadSchemas(fsys fs.FS, dir string) (*schemaSet, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	s := &schemaSet{files: map[string]*node{}}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".schema.json") {
			continue
		}
		b, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var n node
		if err := json.Unmarshal(b, &n); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		setFile(&n, e.Name())
		s.files[e.Name()] = &n
	}
	return s, nil
}

// kinds are the files that describe a document: they have a const kind.
func (s *schemaSet) kinds() []string {
	var out []string
	for name, n := range s.files {
		if k, ok := n.Properties["kind"]; ok && k.Const != nil {
			out = append(out, strings.TrimSuffix(name, ".schema.json"))
		}
	}
	sort.Strings(out)
	return out
}

// resolve follows $ref until it reaches a node without one. It returns the
// last ref followed, so callers can recognise shared definitions such as
// valueFrom.
func (s *schemaSet) resolve(n *node) (*node, string, error) {
	last := ""
	for i := 0; n != nil && n.Ref != ""; i++ {
		if i > 32 {
			return nil, "", fmt.Errorf("$ref loop at %s", n.Ref)
		}
		file, frag, _ := strings.Cut(n.Ref, "#")
		if file == "" {
			file = n.file
		}
		root, ok := s.files[file]
		if !ok {
			return nil, "", fmt.Errorf("$ref %s: no schema file %s", n.Ref, file)
		}
		target := root
		for _, part := range strings.Split(strings.TrimPrefix(frag, "/"), "/") {
			if part == "" {
				continue
			}
			switch part {
			case "$defs":
				continue
			case "properties":
				continue
			}
			next := target.Defs[part]
			if next == nil {
				next = target.Properties[part]
			}
			if next == nil {
				return nil, "", fmt.Errorf("$ref %s: %s not found", n.Ref, part)
			}
			target = next
		}
		last = file + "#" + frag
		// A $ref beside other keywords (description, default) keeps them:
		// the kind's own description wins over the shared one.
		merged := *target
		if n.Description != "" {
			merged.Description = n.Description
		}
		if n.Default != nil {
			merged.Default = n.Default
		}
		n = &merged
	}
	return n, last, nil
}

// flatten resolves n and folds its allOf branches into one node, the way the
// singletonMetadata definition layers a const name over metadata.
func (s *schemaSet) flatten(n *node) (*node, string, error) {
	r, ref, err := s.resolve(n)
	if err != nil {
		return nil, "", err
	}
	if len(r.AllOf) == 0 {
		return r, ref, nil
	}
	out := *r
	out.AllOf = nil
	out.Properties = map[string]*node{}
	for k, v := range r.Properties {
		out.Properties[k] = v
	}
	for _, branch := range r.AllOf {
		b, _, err := s.flatten(branch)
		if err != nil {
			return nil, "", err
		}
		if len(out.Type) == 0 {
			out.Type = b.Type
		}
		if out.Description == "" {
			out.Description = b.Description
		}
		out.Required = append(out.Required, b.Required...)
		for k, v := range b.Properties {
			if cur, ok := out.Properties[k]; ok {
				// Keep the stricter, later layer but not lose a ref target.
				c := *cur
				if v.Const != nil {
					c.Const = v.Const
				}
				out.Properties[k] = &c
				continue
			}
			out.Properties[k] = v
		}
		if len(out.AddlProps) == 0 {
			out.AddlProps = b.AddlProps
		}
	}
	return &out, ref, nil
}

// isValueFrom reports whether a property is a secret: a valueFrom
// reference, an object whose one property is valueFrom (the shape export
// writes), or a oneOf/anyOf with either as a branch.
func (s *schemaSet) isValueFrom(n *node) (bool, error) {
	r, ref, err := s.resolve(n)
	if err != nil {
		return false, err
	}
	if strings.HasSuffix(ref, "/$defs/valueFrom") {
		return true, nil
	}
	if len(r.Properties) == 1 {
		if v, ok := r.Properties["valueFrom"]; ok {
			_, vref, err := s.resolve(v)
			if err != nil {
				return false, err
			}
			if strings.HasSuffix(vref, "/$defs/valueFrom") {
				return true, nil
			}
		}
	}
	for _, list := range [][]*node{r.OneOf, r.AnyOf} {
		for _, b := range list {
			ok, err := s.isValueFrom(b)
			if err != nil || ok {
				return ok, err
			}
		}
	}
	return false, nil
}
