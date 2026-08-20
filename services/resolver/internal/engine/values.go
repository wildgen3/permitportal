package engine

import (
	"fmt"
	"sort"
	"strings"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types/ref"
	"gopkg.in/yaml.v3"

	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// Profile is a subject's fact set, keyed by scope unit and then by attribute.
//
// Scope is part of the key rather than a property of the value because regulations
// key to different units of analysis, and a fact set that flattens them cannot
// express "eight employees company-wide at a site with a different industry code" —
// which is the case OSHA's partial exemption turns on.
type Profile struct {
	Facts map[string]map[string]any
}

// LocalKey is the profile spelling of a registered attribute: the URI with its
// scope-name prefix removed when it has one, and dots replaced by underscores.
//
//	establishment.classification              (scope establishment) -> classification
//	company.peak_employment_prior_cy          (scope business)      -> company_peak_employment_prior_cy
func LocalKey(attr spec.Attribute) string {
	uri := attr.URI
	if p := attr.ScopeUnit + "."; strings.HasPrefix(uri, p) {
		uri = strings.TrimPrefix(uri, p)
	}
	return strings.ReplaceAll(uri, ".", "_")
}

func (p Profile) Lookup(attr spec.Attribute) (any, bool) {
	scoped, ok := p.Facts[attr.ScopeUnit]
	if !ok {
		return nil, false
	}
	v, ok := scoped[LocalKey(attr)]
	if !ok || v == nil {
		return nil, false
	}
	return v, true
}

// Validate rejects a profile carrying a key that resolves to no registered attribute
// at that scope.
//
// Without this, a typo — or an attribute written at the wrong scope — reads as
// "absent", the determination comes back INDETERMINATE, and the product asks a
// question the user already answered. Silence is the wrong failure here, so an
// unrecognised key is an error.
func (p Profile) Validate(c *spec.Corpus) error {
	index := map[string]map[string]struct{}{}
	for _, a := range c.Attributes {
		if index[a.ScopeUnit] == nil {
			index[a.ScopeUnit] = map[string]struct{}{}
		}
		index[a.ScopeUnit][LocalKey(a)] = struct{}{}
	}

	var problems []string
	for scope, facts := range p.Facts {
		known, ok := index[scope]
		if !ok {
			problems = append(problems, fmt.Sprintf("scope %q is not a scope unit used by any registered attribute", scope))
			continue
		}
		for key := range facts {
			if _, ok := known[key]; !ok {
				problems = append(problems, fmt.Sprintf("%s.%s is not a registered attribute at that scope", scope, key))
			}
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("profile does not match the attribute registry:\n  - %s",
			strings.Join(problems, "\n  - "))
	}
	return nil
}

// ---------------------------------------------------------------- coercion

func numericLiteral(n *spec.Node) (string, *cel.Type, error) {
	if n.Value.Kind == 0 {
		return "", nil, fmt.Errorf("predicate %q has no value", n.Predicate)
	}
	var i int64
	if err := n.Value.Decode(&i); err == nil {
		return fmt.Sprintf("%d", i), cel.IntType, nil
	}
	var f float64
	if err := n.Value.Decode(&f); err == nil {
		return fmt.Sprintf("%v", f), cel.DoubleType, nil
	}
	return "", nil, fmt.Errorf("predicate %q: value is neither an integer nor a decimal", n.Predicate)
}

func coerce(v any, t *cel.Type, uri string) (any, error) {
	switch t {
	case cel.IntType:
		switch n := v.(type) {
		case int:
			return int64(n), nil
		case int64:
			return n, nil
		case float64:
			if n != float64(int64(n)) {
				return nil, fmt.Errorf("%s: %v is not an integer", uri, n)
			}
			return int64(n), nil
		}
		return nil, fmt.Errorf("%s: expected an integer, got %T", uri, v)

	case cel.DoubleType:
		switch n := v.(type) {
		case int:
			return float64(n), nil
		case int64:
			return float64(n), nil
		case float64:
			return n, nil
		}
		return nil, fmt.Errorf("%s: expected a number, got %T", uri, v)

	case cel.BoolType:
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("%s: expected a boolean, got %T", uri, v)
		}
		return b, nil
	}

	switch t.String() {
	case cel.ListType(cel.StringType).String():
		items, ok := v.([]any)
		if !ok {
			if ss, ok := v.([]string); ok {
				return ss, nil
			}
			return nil, fmt.Errorf("%s: expected a list of strings, got %T", uri, v)
		}
		out := make([]string, 0, len(items))
		for _, it := range items {
			s, ok := it.(string)
			if !ok {
				return nil, fmt.Errorf("%s: list element %v is not a string", uri, it)
			}
			out = append(out, s)
		}
		return out, nil

	case cel.MapType(cel.StringType, cel.StringType).String():
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s: expected a mapping, got %T", uri, v)
		}
		out := make(map[string]string, len(m))
		for k, val := range m {
			s, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s: %s is %T, expected a string — codes and vintages "+
					"are strings, because YAML will happily turn 2022 into an integer and "+
					"a leading zero into nothing", uri, k, val)
			}
			out[k] = s
		}
		return out, nil
	}
	return nil, fmt.Errorf("%s: unsupported attribute type %s", uri, t.String())
}

func stringMap(v ref.Val) (map[string]string, error) {
	raw, err := v.ConvertToNative(mapStringStringType)
	if err != nil {
		return nil, err
	}
	m, ok := raw.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("expected map[string]string, got %T", raw)
	}
	return m, nil
}

func stringSlice(v ref.Val) ([]string, error) {
	raw, err := v.ConvertToNative(sliceStringType)
	if err != nil {
		return nil, err
	}
	s, ok := raw.([]string)
	if !ok {
		return nil, fmt.Errorf("expected []string, got %T", raw)
	}
	return s, nil
}

// UnmarshalProfile decodes a YAML fact set.
func UnmarshalProfile(data []byte) (Profile, error) {
	var facts map[string]map[string]any
	if err := yaml.Unmarshal(data, &facts); err != nil {
		return Profile{}, err
	}
	return Profile{Facts: facts}, nil
}
