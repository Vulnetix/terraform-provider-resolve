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
	UniqueItems bool             `json:"uniqueItems"`
	Properties  map[string]*node `json:"properties"`
	Required    []string         `json:"required"`
	AddlProps   json.RawMessage  `json:"additionalProperties"`
	Items       *node            `json:"items"`
	AllOf       []*node          `json:"allOf"`
	OneOf       []*node          `json:"oneOf"`
	AnyOf       []*node          `json:"anyOf"`
	Defs        map[string]*node `json:"$defs"`
	// Conditional constraints, read only to find properties a document may not
	// always carry (the reference example leaves those out).
	If   *node `json:"if"`
	Then *node `json:"then"`
	Else *node `json:"else"`
	Not  *node `json:"not"`
	// X is the root x-vulnetix annotation (config-schemas v1.3.0 and later).
	X *annotations `json:"x-vulnetix"`

	// file is the schema file this node was read from, for resolving its
	// local "#/..." refs.
	file string
	// boolean is set for a boolean schema (true or false). The kind schemas
	// use false inside oneOf branches to forbid a property, which only
	// constrains values the server checks anyway, so the generator ignores it.
	boolean *bool
}

// UnmarshalJSON accepts the boolean schemas 2020-12 allows wherever a schema
// may appear, such as "properties": {"provider": false}.
func (n *node) UnmarshalJSON(b []byte) error {
	var v bool
	if json.Unmarshal(b, &v) == nil {
		*n = node{boolean: &v}
		return nil
	}
	type plain node
	return json.Unmarshal(b, (*plain)(n))
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
	for _, c := range []*node{n.If, n.Then, n.Else, n.Not} {
		setFile(c, file)
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
		// the kind's own description wins over the shared one, and so do a
		// narrower pattern and a const.
		merged := *target
		if n.Description != "" {
			merged.Description = n.Description
		}
		if n.Default != nil {
			merged.Default = n.Default
		}
		if n.Pattern != "" {
			merged.Pattern = n.Pattern
		}
		if n.Const != nil {
			merged.Const = n.Const
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
		u, err := s.union(r)
		return u, ref, err
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
				// Keep the stricter, later layer but not lose a ref target: a
				// const name (singletonMetadata) or a narrower pattern and its
				// description (AiProvider's slug).
				c := *cur
				if v.Const != nil {
					c.Const = v.Const
				}
				if v.Pattern != "" {
					c.Pattern = v.Pattern
				}
				if v.Description != "" {
					c.Description = v.Description
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
	u, err := s.union(&out)
	return u, ref, err
}

// union folds an untyped oneOf or anyOf whose branches share a shape into one
// node, so it can be typed instead of carried as JSON:
//
//   - branches that are all strings (RiskStrategy metric: three enum refs; a
//     PackageList purl: two patterns) become one string, its enum the union
//     when every branch is an enum;
//   - branches that are all objects with properties (RemediationPolicy spec:
//     a standard, or custom windows) become one object with every branch's
//     properties, none required, since which are required depends on the
//     branch. The server checks the branch.
//
// A node with its own type or properties keeps them; its oneOf only
// constrains them (VexSource spec: name exactly one source).
func (s *schemaSet) union(n *node) (*node, error) {
	branches := n.OneOf
	if len(branches) == 0 {
		branches = n.AnyOf
	}
	if len(branches) == 0 || len(n.types()) > 0 || len(n.Properties) > 0 || len(n.Enum) > 0 || n.Const != nil {
		return n, nil
	}
	var flat []*node
	strs, objs := true, true
	for _, b := range branches {
		f, _, err := s.flatten(b)
		if err != nil {
			return nil, err
		}
		flat = append(flat, f)
		isStr := f.hasType("string") && len(f.types()) == 1 || (len(f.types()) == 0 && len(f.Properties) == 0 && (len(f.Enum) > 0 || f.Const != nil) && constType(append([]any{f.Const}, f.Enum...)) == "string")
		strs = strs && isStr
		objs = objs && len(f.Properties) > 0 && (len(f.types()) == 0 || f.hasType("object"))
	}
	out := *n
	out.OneOf, out.AnyOf = nil, nil
	switch {
	case strs:
		out.Type = json.RawMessage(`"string"`)
		allEnum := true
		var enum []any
		for _, f := range flat {
			vals := append([]any{}, f.Enum...)
			if f.Const != nil {
				vals = append(vals, f.Const)
			}
			if len(vals) == 0 {
				allEnum = false
			}
			enum = append(enum, vals...)
		}
		if allEnum {
			out.Enum = enum
		}
		return &out, nil
	case objs:
		out.Type = json.RawMessage(`"object"`)
		out.Properties = map[string]*node{}
		for _, f := range flat {
			for k, v := range f.Properties {
				if _, dup := out.Properties[k]; !dup {
					out.Properties[k] = v
				}
			}
			if len(out.AddlProps) == 0 {
				out.AddlProps = f.AddlProps
			}
		}
		return &out, nil
	}
	return n, nil
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

// annotations is a kind schema's x-vulnetix block: what the server does with
// the kind, beyond the shape of its documents.
type annotations struct {
	// ValidateOnly kinds are checked but not stored; plan answers
	// not_stored. They get no resource.
	ValidateOnly bool `json:"validateOnly"`
	// Singleton kinds have one object per tenant, named default.
	Singleton bool `json:"singleton"`
	// DependsOn are the kinds an object of this kind can reference, applied
	// first.
	DependsOn []string `json:"dependsOn"`
	// Secrets are the dotted paths (spec.apiKey) of valueFrom fields.
	Secrets []string `json:"secrets"`
	// CelFields are the CEL expressions: a dotted path ([] for list
	// elements) and the CEL environment it is evaluated in.
	CelFields []celField `json:"celFields"`
}

type celField struct {
	Path string `json:"path"`
	Env  string `json:"env"`
}

// annotationsOf returns a kind's x-vulnetix block, or nil for a release
// that predates them.
func (s *schemaSet) annotationsOf(kind string) *annotations {
	if n, ok := s.files[kind+".schema.json"]; ok {
		return n.X
	}
	return nil
}

// isValidateOnly reads x-vulnetix.validateOnly. A schema without the
// annotation (the test schemas) is stored.
func (s *schemaSet) isValidateOnly(kind string) bool {
	a := s.annotationsOf(kind)
	return a != nil && a.ValidateOnly
}
