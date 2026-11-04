package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"gopkg.in/yaml.v3"
)

// The examples are written from config-schemas examples/<Kind>/minimal.yaml
// and advanced.yaml, so the HCL a user copies is the same document the CLI
// and the console's Download as code use, in the provider's attribute names:
//
//	examples/resources/vulnetix_<kind>/minimal.tf
//	examples/resources/vulnetix_<kind>/advanced.tf
//	examples/resources/vulnetix_<kind>/import.sh
//
// tfplugindocs (templates/resources.md.tmpl) renders them under Minimal,
// Advanced and Import.

// Attribute shapes, for walking a document beside its attributes.
const (
	shapeScalar     = "scalar"
	shapeObject     = "object"
	shapeMap        = "map"
	shapeMapObject  = "mapObject"
	shapeList       = "list"
	shapeListObject = "listObject"
	shapeJSON       = "json"
	shapeSecret     = "secret"
)

// hfield is one attribute as the example writer sees it.
type hfield struct {
	Attr, JSON string
	Shape      string
	// Elem is the field kind of a scalar, or of a list or map element
	// (fieldString and so on).
	Elem     string
	Children []hfield
}

// exampleVariants are the example files every kind has, in the order the
// docs show them.
var exampleVariants = []string{"minimal", "advanced"}

// examplesDir finds config-schemas' examples: the -examples flag, else the
// module directory of the config-schemas version this module requires.
func examplesDir(flagDir string) (string, error) {
	if flagDir != "" {
		return flagDir, nil
	}
	out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", schemasModule).Output()
	if err != nil {
		return "", fmt.Errorf("locate %s: %w", schemasModule, err)
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("locate %s: not downloaded; run go mod download", schemasModule)
	}
	return filepath.Join(dir, "examples"), nil
}

// writeExamples writes one kind's example directory.
func writeExamples(out, src string, k *kindOut) error {
	dir := filepath.Join(out, "vulnetix_"+k.TypeName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	var firstName string
	for _, v := range exampleVariants {
		raw, err := os.ReadFile(filepath.Join(src, k.Kind, v+".yaml"))
		if err != nil {
			return fmt.Errorf("%s: %s example: %w", k.Kind, v, err)
		}
		hcl, name, err := exampleHCL(k, v, raw)
		if err != nil {
			return fmt.Errorf("%s: %s.yaml: %w", k.Kind, v, err)
		}
		if firstName == "" {
			firstName = name
		}
		if err := os.WriteFile(filepath.Join(dir, v+".tf"), hcl, 0o644); err != nil {
			return err
		}
	}
	if err := writeReference(out, src, k); err != nil {
		return err
	}
	imp := fmt.Sprintf("# Import an existing %s by its name (metadata.name).\nterraform import vulnetix_%s.%s %s\n",
		k.Kind, k.TypeName, resourceLabel(firstName), firstName)
	return os.WriteFile(filepath.Join(dir, "import.sh"), []byte(imp), 0o644)
}

// removeStaleExamples deletes the generated files of every example
// directory, so a kind that is no longer served leaves nothing behind.
func removeStaleExamples(out string) error {
	for _, name := range append(append([]string{}, exampleVariants...), "reference", "import") {
		ext := ".tf"
		if name == "import" {
			ext = ".sh"
		}
		if name == "reference" {
			_ = removeGlob(filepath.Join(out, "vulnetix_*", "reference.yaml"))
		}
		matches, err := filepath.Glob(filepath.Join(out, "vulnetix_*", name+ext))
		if err != nil {
			return err
		}
		for _, m := range matches {
			if err := os.Remove(m); err != nil {
				return err
			}
		}
	}
	dirs, _ := filepath.Glob(filepath.Join(out, "vulnetix_*"))
	for _, d := range dirs {
		if entries, err := os.ReadDir(d); err == nil && len(entries) == 0 {
			_ = os.Remove(d)
		}
	}
	return nil
}

// exampleDoc is one YAML document of an example file.
type exampleDoc struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name        string            `yaml:"name"`
		Labels      map[string]string `yaml:"labels"`
		Annotations map[string]string `yaml:"annotations"`
	} `yaml:"metadata"`
	Spec map[string]any `yaml:"-"`
}

// exampleHCL converts an example file (one or more documents) to HCL, one
// resource per document. It returns the first document's name.
func exampleHCL(k *kindOut, variant string, raw []byte) ([]byte, string, error) {
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	var docs []exampleDoc
	for {
		var n yaml.Node
		err := dec.Decode(&n)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, "", err
		}
		var d exampleDoc
		if err := n.Decode(&d); err != nil {
			return nil, "", err
		}
		if sn := mappingValue(&n, "spec"); sn != nil {
			v, err := nodeValue(sn)
			if err != nil {
				return nil, "", err
			}
			d.Spec, _ = v.(map[string]any)
		}
		if d.Kind == "" && d.Spec == nil && d.Metadata.Name == "" {
			continue
		}
		if d.Kind != k.Kind {
			return nil, "", fmt.Errorf("document of kind %q in the %s examples", d.Kind, k.Kind)
		}
		docs = append(docs, d)
	}
	if len(docs) == 0 {
		return nil, "", errors.New("no documents")
	}

	var vars, body, alts bytes.Buffer
	seen := map[string]bool{}
	for i, d := range docs {
		// A kind with one object per tenant can only be declared once. Its
		// further documents are alternatives, written commented out.
		alt := k.Singleton && i > 0
		var doc, altVars bytes.Buffer
		dv := &vars
		if alt {
			dv = &altVars
		}
		label := resourceLabel(d.Metadata.Name)
		for n := 2; seen[label]; n++ {
			label = resourceLabel(d.Metadata.Name) + "_" + strconv.Itoa(n)
		}
		seen[label] = true
		w := &hclWriter{vars: dv, prefix: label}
		fmt.Fprintf(&doc, "resource \"vulnetix_%s\" %q {\n", k.TypeName, label)
		if !k.Singleton || d.Metadata.Name != "default" {
			fmt.Fprintf(&doc, "name = %s\n", hclString(d.Metadata.Name))
		}
		if err := w.object(&doc, d.Spec, k.Tree, "spec"); err != nil {
			return nil, "", err
		}
		for _, meta := range []struct {
			attr string
			m    map[string]string
		}{{"labels", d.Metadata.Labels}, {"annotations", d.Metadata.Annotations}} {
			if len(meta.m) == 0 {
				continue
			}
			fmt.Fprintf(&doc, "%s = ", meta.attr)
			anyMap := map[string]any{}
			for key, v := range meta.m {
				anyMap[key] = v
			}
			w.freeMap(&doc, anyMap)
			doc.WriteString("\n")
		}
		doc.WriteString("}\n")
		if alt {
			alts.WriteString("\n# Or, instead (one per tenant, so declare only one):\n#\n")
			for _, line := range strings.Split(strings.TrimRight(string(hclwrite.Format(append(altVars.Bytes(), doc.Bytes()...))), "\n"), "\n") {
				alts.WriteString(strings.TrimRight("# "+line, " ") + "\n")
			}
			continue
		}
		if body.Len() > 0 {
			body.WriteString("\n")
		}
		body.Write(doc.Bytes())
	}

	var b bytes.Buffer
	switch variant {
	case "reference":
		fmt.Fprintf(&b, "# Every argument with a schema default, set to that default, plus the required\n# arguments from config-schemas examples/%s/minimal.yaml. A plan only sets values\n# a tenant left unset. See advanced.tf for the arguments that have no default.\n\n", k.Kind)
	case "advanced":
		fmt.Fprintf(&b, "# Every argument the kind supports. From config-schemas examples/%s/%s.yaml.\n\n", k.Kind, variant)
	default:
		fmt.Fprintf(&b, "# Only the arguments the kind requires. From config-schemas examples/%s/%s.yaml.\n\n", k.Kind, variant)
	}
	if vars.Len() > 0 {
		b.Write(vars.Bytes())
		b.WriteString("\n")
	}
	b.Write(body.Bytes())
	out := hclwrite.Format(b.Bytes())
	return append(out, alts.Bytes()...), docs[0].Metadata.Name, nil
}

// resourceLabel turns an object name into a Terraform resource label.
func resourceLabel(name string) string {
	l := strings.NewReplacer("-", "_", ".", "_").Replace(name)
	if l == "" || (l[0] >= '0' && l[0] <= '9') {
		l = "r_" + l
	}
	return l
}

type hclWriter struct {
	vars   *bytes.Buffer
	prefix string
}

// object writes the attributes of one object, in attribute order.
func (w *hclWriter) object(b *bytes.Buffer, m map[string]any, fields []hfield, where string) error {
	byJSON := map[string]hfield{}
	for _, f := range fields {
		byJSON[f.JSON] = f
	}
	keys := make([]string, 0, len(m))
	for key := range m {
		if _, ok := byJSON[key]; !ok {
			return fmt.Errorf("%s.%s: no such argument", where, key)
		}
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return byJSON[keys[i]].Attr < byJSON[keys[j]].Attr })
	for _, key := range keys {
		f := byJSON[key]
		if f.Shape == shapeSecret {
			v := w.prefix + "_" + strings.ReplaceAll(strings.TrimPrefix(snakePath(where+"."+key), "spec_"), ".", "_")
			fmt.Fprintf(w.vars, "variable %q {\ntype = string\nsensitive = true\nephemeral = true\n}\n", v)
			fmt.Fprintf(b, "%s_wo = var.%s\n%s_wo_version = 1\n", f.Attr, v, f.Attr)
			continue
		}
		fmt.Fprintf(b, "%s = ", f.Attr)
		if err := w.value(b, m[key], f, where+"."+key); err != nil {
			return err
		}
		b.WriteString("\n")
	}
	return nil
}

func snakePath(p string) string {
	parts := strings.Split(p, ".")
	for i, s := range parts {
		parts[i] = snake(s)
	}
	return strings.Join(parts, "_")
}

// value writes one attribute's value.
func (w *hclWriter) value(b *bytes.Buffer, v any, f hfield, where string) error {
	if v == nil {
		b.WriteString("null")
		return nil
	}
	switch f.Shape {
	case shapeScalar:
		return scalarHCL(b, typed(v, f.Elem), where)
	case shapeJSON:
		b.WriteString("jsonencode(")
		w.free(b, v)
		b.WriteString(")")
		return nil
	case shapeObject:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: expected an object, got %T", where, v)
		}
		b.WriteString("{\n")
		if err := w.object(b, m, f.Children, where); err != nil {
			return err
		}
		b.WriteString("}")
		return nil
	case shapeMap, shapeMapObject:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: expected a map, got %T", where, v)
		}
		keys := sortedKeys(m)
		b.WriteString("{\n")
		for _, key := range keys {
			fmt.Fprintf(b, "%s = ", hclKey(key))
			if f.Shape == shapeMap {
				if err := scalarHCL(b, typed(m[key], f.Elem), where+"."+key); err != nil {
					return err
				}
			} else {
				el, ok := m[key].(map[string]any)
				if !ok {
					return fmt.Errorf("%s.%s: expected an object, got %T", where, key, m[key])
				}
				b.WriteString("{\n")
				if err := w.object(b, el, f.Children, where+"."+key); err != nil {
					return err
				}
				b.WriteString("}")
			}
			b.WriteString("\n")
		}
		b.WriteString("}")
		return nil
	case shapeList, shapeListObject:
		arr, ok := v.([]any)
		if !ok {
			return fmt.Errorf("%s: expected a list, got %T", where, v)
		}
		if f.Shape == shapeList {
			b.WriteString("[")
			for i, e := range arr {
				if i > 0 {
					b.WriteString(", ")
				}
				if err := scalarHCL(b, typed(e, f.Elem), fmt.Sprintf("%s[%d]", where, i)); err != nil {
					return err
				}
			}
			b.WriteString("]")
			return nil
		}
		b.WriteString("[\n")
		for i, e := range arr {
			el, ok := e.(map[string]any)
			if !ok {
				return fmt.Errorf("%s[%d]: expected an object, got %T", where, i, e)
			}
			b.WriteString("{\n")
			if err := w.object(b, el, f.Children, fmt.Sprintf("%s[%d]", where, i)); err != nil {
				return err
			}
			b.WriteString("},\n")
		}
		b.WriteString("]")
		return nil
	}
	return fmt.Errorf("%s: unknown shape %s", where, f.Shape)
}

// scalarHCL writes a string, number or boolean.
func scalarHCL(b *bytes.Buffer, v any, where string) error {
	switch x := v.(type) {
	case string:
		b.WriteString(hclString(x))
	case bool:
		b.WriteString(strconv.FormatBool(x))
	case int:
		b.WriteString(strconv.Itoa(x))
	case int64:
		b.WriteString(strconv.FormatInt(x, 10))
	case uint64:
		b.WriteString(strconv.FormatUint(x, 10))
	case float64:
		b.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	default:
		return fmt.Errorf("%s: expected a scalar, got %T", where, v)
	}
	return nil
}

// free writes any YAML value as an HCL expression with quoted keys: the
// inside of jsonencode(), where keys keep their JSON spelling.
func (w *hclWriter) free(b *bytes.Buffer, v any) {
	switch x := v.(type) {
	case map[string]any:
		w.freeMap(b, x)
	case []any:
		b.WriteString("[")
		for i, e := range x {
			if i > 0 {
				b.WriteString(", ")
			}
			w.free(b, e)
		}
		b.WriteString("]")
	case nil:
		b.WriteString("null")
	default:
		if err := scalarHCL(b, x, ""); err != nil {
			j, _ := json.Marshal(x)
			b.WriteString(hclString(string(j)))
		}
	}
}

func (w *hclWriter) freeMap(b *bytes.Buffer, m map[string]any) {
	b.WriteString("{\n")
	for _, key := range sortedKeys(m) {
		fmt.Fprintf(b, "%s = ", hclKey(key))
		w.free(b, m[key])
		b.WriteString("\n")
	}
	b.WriteString("}")
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// hclKey is a map key: bare when it is an identifier, quoted otherwise.
func hclKey(k string) string {
	if k == "" {
		return `""`
	}
	for i, r := range k {
		ok := r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (i > 0 && (r == '-' || (r >= '0' && r <= '9')))
		if !ok {
			return hclString(k)
		}
	}
	return k
}

// hclString quotes a string, escaping HCL's template sequences so a CEL
// expression or a glob is written literally.
func hclString(s string) string {
	q := strconv.Quote(s)
	q = strings.ReplaceAll(q, "${", "$${")
	return strings.ReplaceAll(q, "%{", "%%{")
}

// mappingValue is the value of key in a document or mapping node.
func mappingValue(n *yaml.Node, key string) *yaml.Node {
	if n.Kind == yaml.DocumentNode && len(n.Content) == 1 {
		n = n.Content[0]
	}
	if n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

// nodeValue converts YAML to plain values, keeping a scalar as written: an
// unquoted date stays the string it is in the document, as the server reads
// it, rather than becoming a time.
func nodeValue(n *yaml.Node) (any, error) {
	switch n.Kind {
	case yaml.AliasNode:
		return nodeValue(n.Alias)
	case yaml.MappingNode:
		m := map[string]any{}
		for i := 0; i+1 < len(n.Content); i += 2 {
			v, err := nodeValue(n.Content[i+1])
			if err != nil {
				return nil, err
			}
			m[n.Content[i].Value] = v
		}
		return m, nil
	case yaml.SequenceNode:
		out := []any{}
		for _, c := range n.Content {
			v, err := nodeValue(c)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!null":
			return nil, nil
		case "!!bool", "!!int", "!!float":
			var v any
			if err := n.Decode(&v); err != nil {
				return nil, err
			}
			return v, nil
		}
		return n.Value, nil
	}
	return nil, fmt.Errorf("line %d: unsupported YAML node", n.Line)
}

// typed writes a number as a string where the attribute is a string (a
// guardrail pattern may be 200 or "200"), so the example reads as the
// attribute's type.
func typed(v any, elem string) any {
	if elem != "fieldString" {
		return v
	}
	switch x := v.(type) {
	case int, int64, uint64, float64, bool:
		var b bytes.Buffer
		_ = scalarHCL(&b, x, "")
		return b.String()
	}
	return v
}

// removeGlob deletes every file matching pattern.
func removeGlob(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil {
			return err
		}
	}
	return nil
}
