package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The resource works on raw tftypes values rather than per-kind structs, so a
// kind is nothing more than its generated schema and field list.

// objectAttrs splits an object value into its attributes. A null or unknown
// value has none.
func objectAttrs(v tftypes.Value) (map[string]tftypes.Value, error) {
	if v.IsNull() || !v.IsKnown() {
		return map[string]tftypes.Value{}, nil
	}
	m := map[string]tftypes.Value{}
	if err := v.As(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func isSet(v tftypes.Value) bool {
	return v.Type() != nil && v.IsKnown() && !v.IsNull()
}

// specInput is what toSpec reads: the plan, plus the config (the only place a
// write-only value exists) and the prior state (to tell whether a secret's
// version moved). config and prior may be null.
type specInput struct {
	plan, config, prior tftypes.Value
	create              bool
}

// toSpec builds a document spec from the attributes fields describe.
func toSpec(in specInput, fields []field) (map[string]any, error) {
	plan, err := objectAttrs(in.plan)
	if err != nil {
		return nil, err
	}
	config, err := objectAttrs(in.config)
	if err != nil {
		return nil, err
	}
	prior, err := objectAttrs(in.prior)
	if err != nil {
		return nil, err
	}
	spec := map[string]any{}
	for _, f := range fields {
		switch {
		case f.Local:
			continue
		case f.WriteOnly:
			cv, ok := config[f.Attr]
			if !ok || !isSet(cv) || !secretChanged(f, in.create, plan, prior) {
				continue
			}
			var s string
			if err := cv.As(&s); err != nil {
				return nil, fmt.Errorf("%s: %w", f.Attr, err)
			}
			spec[f.JSON] = s
			continue
		}
		v, ok := plan[f.Attr]
		if !ok || !isSet(v) {
			continue
		}
		child := specInput{plan: v, config: nullOr(config[f.Attr]), prior: nullOr(prior[f.Attr]), create: in.create || !isSet(nullOr(prior[f.Attr]))}
		out, err := valueToJSON(child, f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Attr, err)
		}
		spec[f.JSON] = out
	}
	return spec, nil
}

// secretChanged reports whether a write-only value must be sent. Terraform
// never shows a write-only value in a plan, so the only signal that it
// changed is its _wo_version sibling. Without one it is sent every time.
func secretChanged(f field, create bool, plan, prior map[string]tftypes.Value) bool {
	if create || f.Version == "" || len(prior) == 0 {
		return true
	}
	pv, cv := nullOr(prior[f.Version]), nullOr(plan[f.Version])
	return !pv.Equal(cv)
}

var nullValue = tftypes.NewValue(tftypes.DynamicPseudoType, nil)

func nullOr(v tftypes.Value) tftypes.Value {
	if v.Type() == nil {
		return nullValue
	}
	return v
}

func valueToJSON(in specInput, f field) (any, error) {
	v := in.plan
	switch f.Kind {
	case fieldObject:
		return toSpec(in, f.Children)
	case fieldList, fieldSet:
		var elems []tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make([]any, 0, len(elems))
		for _, e := range elems {
			j, err := scalarToJSON(e, f.Elem)
			if err != nil {
				return nil, err
			}
			out = append(out, j)
		}
		return out, nil
	case fieldMap:
		var elems map[string]tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make(map[string]any, len(elems))
		for k, e := range elems {
			j, err := scalarToJSON(e, f.Elem)
			if err != nil {
				return nil, err
			}
			out[k] = j
		}
		return out, nil
	case fieldListObject, fieldSetObject:
		var elems []tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make([]any, 0, len(elems))
		for _, e := range elems {
			j, err := toSpec(specInput{plan: e, config: nullValue, prior: nullValue, create: true}, f.Children)
			if err != nil {
				return nil, err
			}
			out = append(out, j)
		}
		return out, nil
	case fieldMapObject:
		var elems map[string]tftypes.Value
		if err := v.As(&elems); err != nil {
			return nil, err
		}
		out := make(map[string]any, len(elems))
		for k, e := range elems {
			j, err := toSpec(specInput{plan: e, config: nullValue, prior: nullValue, create: true}, f.Children)
			if err != nil {
				return nil, err
			}
			out[k] = j
		}
		return out, nil
	}
	return scalarToJSON(v, f.Kind)
}

func scalarToJSON(v tftypes.Value, k fieldKind) (any, error) {
	if !isSet(v) {
		return nil, nil
	}
	switch k {
	case fieldString:
		var s string
		err := v.As(&s)
		return s, err
	case fieldJSON:
		var s string
		if err := v.As(&s); err != nil {
			return nil, err
		}
		var out any
		dec := json.NewDecoder(bytes.NewReader([]byte(s)))
		dec.UseNumber()
		if err := dec.Decode(&out); err != nil {
			return nil, fmt.Errorf("not valid JSON: %w", err)
		}
		return out, nil
	case fieldBool:
		var b bool
		err := v.As(&b)
		return b, err
	case fieldInt:
		var n big.Float
		if err := v.As(&n); err != nil {
			return nil, err
		}
		i, acc := n.Int64()
		if acc != big.Exact {
			return nil, fmt.Errorf("%s is not a whole number", n.String())
		}
		return i, nil
	case fieldFloat:
		var n big.Float
		if err := v.As(&n); err != nil {
			return nil, err
		}
		f, _ := n.Float64()
		return f, nil
	}
	return nil, fmt.Errorf("field kind %d is not a scalar", k)
}

// fromSpec builds the attribute values fields describe from a document spec.
//
// With keepAll false only attributes already set in prior are filled in: the
// engine merges, and the tenant's row carries server defaults the
// configuration never mentioned, so reading those back would show a diff on
// every plan. keepAll is true on import, where there is no prior to go by.
func fromSpec(spec map[string]any, typ tftypes.Object, prior tftypes.Value, fields []field, keepAll bool) (map[string]tftypes.Value, error) {
	pm, err := objectAttrs(prior)
	if err != nil {
		return nil, err
	}
	out := map[string]tftypes.Value{}
	for _, f := range fields {
		at, ok := typ.AttributeTypes[f.Attr]
		if !ok {
			return nil, fmt.Errorf("schema has no attribute %s", f.Attr)
		}
		pv := pm[f.Attr]
		switch {
		case f.WriteOnly:
			out[f.Attr] = tftypes.NewValue(at, nil)
			continue
		case f.Local:
			if pv.Type() == nil {
				pv = tftypes.NewValue(at, nil)
			}
			out[f.Attr] = pv
			continue
		}
		raw, present := spec[f.JSON]
		if !present || raw == nil || (!keepAll && !isSet(nullOr(pv))) {
			out[f.Attr] = tftypes.NewValue(at, nil)
			continue
		}
		v, err := jsonToValue(raw, at, f, nullOr(pv), keepAll)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Attr, err)
		}
		out[f.Attr] = v
	}
	return out, nil
}

func jsonToValue(raw any, at tftypes.Type, f field, prior tftypes.Value, keepAll bool) (tftypes.Value, error) {
	switch f.Kind {
	case fieldObject:
		m, ok := raw.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected an object, got %T", raw)
		}
		ot := at.(tftypes.Object)
		vals, err := fromSpec(m, ot, prior, f.Children, keepAll)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(ot, vals), nil
	case fieldList, fieldListObject, fieldSet, fieldSetObject:
		arr, ok := raw.([]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected a list, got %T", raw)
		}
		var et tftypes.Type
		if lt, ok := at.(tftypes.List); ok {
			et = lt.ElementType
		} else {
			et = at.(tftypes.Set).ElementType
		}
		var pe []tftypes.Value
		if isSet(prior) {
			if err := prior.As(&pe); err != nil {
				return tftypes.Value{}, err
			}
		}
		used := make([]bool, len(pe))
		vals := make([]tftypes.Value, 0, len(arr))
		for i, e := range arr {
			p := nullValue
			switch {
			case keepAll:
			case f.Kind == fieldListObject && i < len(pe):
				p = pe[i]
			case f.Kind == fieldSetObject:
				p = matchElem(e, et, f, pe, used)
			}
			v, err := elemToValue(e, et, f, p)
			if err != nil {
				return tftypes.Value{}, err
			}
			if (f.Kind == fieldSet || f.Kind == fieldSetObject) && containsValue(vals, v) {
				// A set cannot hold the same element twice, and Terraform
				// refuses a value that does ("Duplicate Set Element"), which
				// made an object whose server copy repeats an element
				// unreadable. The server matches these elements by their
				// fields, so an exact repeat is the same element: keep one.
				continue
			}
			vals = append(vals, v)
		}
		return tftypes.NewValue(at, vals), nil
	case fieldMap, fieldMapObject:
		m, ok := raw.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected a map, got %T", raw)
		}
		et := at.(tftypes.Map).ElementType
		pm := map[string]tftypes.Value{}
		if isSet(prior) && !keepAll {
			if err := prior.As(&pm); err != nil {
				return tftypes.Value{}, err
			}
		}
		vals := make(map[string]tftypes.Value, len(m))
		for k, e := range m {
			v, err := elemToValue(e, et, f, nullOr(pm[k]))
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[k] = v
		}
		return tftypes.NewValue(at, vals), nil
	}
	return scalarToValue(raw, at, f.Kind)
}

// elemToValue converts one list, set or map element. An object element with
// a prior element is read managed-only against it, like a top-level object,
// so the defaults the server fills into each element (a rule's enabled, a
// mirror's priority) never show as a diff. One with no prior (a new element,
// an import) is taken whole.
func elemToValue(e any, et tftypes.Type, f field, prior tftypes.Value) (tftypes.Value, error) {
	if f.Kind == fieldListObject || f.Kind == fieldMapObject || f.Kind == fieldSetObject {
		m, ok := e.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected an object element, got %T", e)
		}
		ot := et.(tftypes.Object)
		vals, err := fromSpec(m, ot, prior, f.Children, !isSet(prior))
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(ot, vals), nil
	}
	return scalarToValue(e, et, f.Elem)
}

// matchElem finds the prior set element a server element is: the first
// unused one that reading the server element managed-only against gives back
// unchanged. A list element is matched by position instead; a set has none,
// and the server returns its elements in its own order. No match (the
// element is new, or a managed field drifted) reads the element whole, so
// the drift shows.
func matchElem(e any, et tftypes.Type, f field, prior []tftypes.Value, used []bool) tftypes.Value {
	for i, p := range prior {
		if used[i] || !isSet(p) {
			continue
		}
		v, err := elemToValue(e, et, f, p)
		if err == nil && v.Equal(p) {
			used[i] = true
			return p
		}
	}
	return nullValue
}

func scalarToValue(raw any, at tftypes.Type, k fieldKind) (tftypes.Value, error) {
	if raw == nil {
		return tftypes.NewValue(at, nil), nil
	}
	switch k {
	case fieldString:
		switch s := raw.(type) {
		case string:
			return tftypes.NewValue(at, s), nil
		case json.Number:
			// A property that is a string or a number (a guardrail pattern)
			// is a string attribute; the server may answer with the number.
			return tftypes.NewValue(at, s.String()), nil
		}
		return tftypes.Value{}, fmt.Errorf("expected a string, got %T", raw)
	case fieldJSON:
		b, err := json.Marshal(raw)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(at, string(b)), nil
	case fieldBool:
		b, ok := raw.(bool)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected a boolean, got %T", raw)
		}
		return tftypes.NewValue(at, b), nil
	case fieldInt, fieldFloat:
		n, err := toBigFloat(raw)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(at, n), nil
	}
	return tftypes.Value{}, fmt.Errorf("field kind %d is not a scalar", k)
}

func toBigFloat(raw any) (*big.Float, error) {
	switch n := raw.(type) {
	case json.Number:
		f, _, err := big.ParseFloat(n.String(), 10, 512, big.ToNearestEven)
		return f, err
	case float64:
		return big.NewFloat(n), nil
	case int64:
		return new(big.Float).SetInt64(n), nil
	case int:
		return new(big.Float).SetInt64(int64(n)), nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, fmt.Errorf("expected a number, got %q", n)
		}
		return big.NewFloat(f), nil
	}
	return nil, fmt.Errorf("expected a number, got %T", raw)
}

// containsValue reports whether vals already holds a value equal to v.
func containsValue(vals []tftypes.Value, v tftypes.Value) bool {
	for _, x := range vals {
		if x.Equal(v) {
			return true
		}
	}
	return false
}
