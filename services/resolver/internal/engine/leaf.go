package engine

import (
	"fmt"
	"sort"
	"strings"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"cel.dev/cel-go/common/types/ref"

	"github.com/wildgen3/permitportal/services/resolver/internal/kleene"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// Leaf predicates compile to CEL and are evaluated with cel-go's PartialActivation
// (ADR-0007). An attribute that is absent from the profile is registered as an
// unknown attribute pattern rather than being bound to a zero value, so cel-go
// returns an unknown and the leaf becomes UNKNOWN instead of FALSE.
//
// The Kleene combinators stay in Go rather than in CEL, for the reason ADR-0007 gives:
// only the leaves compile. CEL has no three-valued semantics of its own and no place
// to hang a per-node citation, and collapsing the tree into one expression would
// destroy the evidence tree.

// Tri-state results returned by the custom list functions. A string rather than a
// bool because "the code is not on the list" and "the code could not be compared to
// the list" are different facts, and only the first of them may become FALSE.
const (
	triTrue          = "TRUE"
	triFalse         = "FALSE"
	triNotComposable = "NOT_COMPOSABLE"
)

type leafOutcome struct {
	Value       kleene.Value
	Translation *TranslationEvidence
	Waived      bool
	CELSource   string
}

func celVar(uri string) string {
	return strings.NewReplacer(".", "_", "-", "_", "/", "_").Replace(uri)
}

// evalLeaf compiles and evaluates one predicate node.
func (e *Engine) evalLeaf(rule *spec.Rule, n *spec.Node, p Profile) (leafOutcome, error) {
	attr, ok := e.corpus.Attributes[n.Attribute]
	if !ok {
		return leafOutcome{}, fmt.Errorf("rule %s: predicate references attribute %q, "+
			"which is not in the registry (linter L-04)", rule.ID, n.Attribute)
	}
	if n.Scope != attr.ScopeUnit {
		return leafOutcome{}, fmt.Errorf("rule %s: predicate on %q declares scope %q but the "+
			"registry declares scope_unit %q — the two-keys-two-granularities bug (linter L-04)",
			rule.ID, n.Attribute, n.Scope, attr.ScopeUnit)
	}

	value, present := p.Lookup(attr)
	varName := celVar(n.Attribute)

	var (
		src     string
		varType *cel.Type
		fns     []cel.EnvOption
		tev     *TranslationEvidence
		list    *spec.CodeList
	)

	switch n.Predicate {
	case "attr_lte", "attr_lt", "attr_gte", "attr_gt", "attr_eq":
		lit, t, err := numericLiteral(n)
		if err != nil {
			return leafOutcome{}, fmt.Errorf("rule %s: %w", rule.ID, err)
		}
		varType = t
		src = fmt.Sprintf("%s %s %s", varName, comparator(n.Predicate), lit)

	case "is_true":
		varType = cel.BoolType
		src = varName

	case "is_false":
		varType = cel.BoolType
		src = "!" + varName

	case "code_in":
		var err error
		list, err = e.requireList(rule, n)
		if err != nil {
			return leafOutcome{}, err
		}
		varType = cel.MapType(cel.StringType, cel.StringType)
		src = fmt.Sprintf("pp_code_in(%s)", varName)
		tev = &TranslationEvidence{Crosswalk: e.crosswalk.Name(), ToVintage: list.ListVintage}
		fns = append(fns, e.codeInFn(n, list, tev))

	case "activity_in", "any_of":
		var err error
		list, err = e.requireList(rule, n)
		if err != nil {
			return leafOutcome{}, err
		}
		varType = cel.ListType(cel.StringType)
		src = fmt.Sprintf("pp_any_in(%s)", varName)
		fns = append(fns, e.anyInFn(list))

	case "location_within":
		if n.JurisdictionProfRef != "self" {
			return leafOutcome{}, fmt.Errorf("rule %s: location_within supports only "+
				"jurisdiction_profile_ref: self, got %q", rule.ID, n.JurisdictionProfRef)
		}
		if rule.JurisdictionProfile == nil {
			return leafOutcome{}, fmt.Errorf("rule %s: location_within references the rule's "+
				"own jurisdiction_profile, which is absent", rule.ID)
		}
		varType = cel.ListType(cel.StringType)
		src = fmt.Sprintf("pp_location_within(%s)", varName)
		fns = append(fns, locationWithinFn(rule.JurisdictionProfile))

	default:
		return leafOutcome{}, fmt.Errorf("rule %s: unsupported predicate %q", rule.ID, n.Predicate)
	}

	opts := append([]cel.EnvOption{cel.Variable(varName, varType)}, fns...)
	env, err := cel.NewEnv(opts...)
	if err != nil {
		return leafOutcome{}, fmt.Errorf("rule %s: building CEL env for %q: %w", rule.ID, n.Attribute, err)
	}
	ast, iss := env.Compile(src)
	if iss != nil && iss.Err() != nil {
		return leafOutcome{}, fmt.Errorf("rule %s: compiling %q: %w", rule.ID, src, iss.Err())
	}
	prg, err := env.Program(ast, cel.EvalOptions(cel.OptPartialEval))
	if err != nil {
		return leafOutcome{}, fmt.Errorf("rule %s: program for %q: %w", rule.ID, src, err)
	}

	vars := map[string]any{}
	var patterns []*cel.AttributePatternType
	if present {
		coerced, err := coerce(value, varType, n.Attribute)
		if err != nil {
			return leafOutcome{}, fmt.Errorf("rule %s: %w", rule.ID, err)
		}
		vars[varName] = coerced
		if tev != nil {
			if m, ok := coerced.(map[string]string); ok {
				tev.FromVintage = m["vintage"]
				tev.FromCode = m["code"]
			}
		}
	} else {
		patterns = append(patterns, cel.AttributePattern(varName))
	}
	act, err := cel.PartialVars(vars, patterns...)
	if err != nil {
		return leafOutcome{}, err
	}

	out, _, err := prg.Eval(act)
	if err != nil && !types.IsUnknown(out) {
		return leafOutcome{}, fmt.Errorf("rule %s: evaluating %q: %w", rule.ID, src, err)
	}

	res := leafOutcome{Translation: tev, CELSource: src}

	if types.IsUnknown(out) {
		// The only way to get here is an attribute the profile does not carry, which
		// is exactly the question to ask next.
		res.Value = kleene.Unk(kleene.Reason{
			Kind: kleene.MissingAttribute,
			Ref:  n.Attribute,
		})
		return res, nil
	}

	switch v := out.Value().(type) {
	case bool:
		res.Value = kleene.Known(v)
	case string:
		val, waived, err := e.applyPolarity(rule, n, list, v)
		if err != nil {
			return leafOutcome{}, err
		}
		res.Value = val
		res.Waived = waived
	default:
		return leafOutcome{}, fmt.Errorf("rule %s: predicate %q produced %T, expected bool or tri-state",
			rule.ID, n.Predicate, out.Value())
	}
	return res, nil
}

// applyPolarity turns a raw membership answer into a truth value, honouring list
// polarity (ADR-0008) and list completeness (linter L-06).
//
// This is the single most consequential function in the engine. "The code is not on
// the list" becomes FALSE only when the list is both semantically exhaustive AND
// completely reproduced here. Everything else is UNKNOWN.
func (e *Engine) applyPolarity(rule *spec.Rule, n *spec.Node, list *spec.CodeList, tri string) (kleene.Value, bool, error) {
	switch tri {
	case triTrue:
		// A hit is dispositive for inclusion under both polarities.
		return kleene.Known(true), false, nil

	case triNotComposable:
		return kleene.Unk(kleene.Reason{
			Kind:   kleene.CrosswalkUnavailable,
			Ref:    n.Attribute,
			Detail: fmt.Sprintf("no crosswalk from the code of record to %s vintage %s", list.ID, list.ListVintage),
		}), false, nil

	case triFalse:
		if list.ListSemantics == spec.SemanticsOpen {
			// A code miss proves nothing. Only 5 of the 11 industrial-activity
			// categories at 40 CFR 122.26(b)(14) reference an industry code at all.
			return kleene.Unk(kleene.Reason{
				Kind:   kleene.ListIllustrativeOpen,
				Ref:    list.ID,
				Detail: "the list is illustrative; a miss does not establish exclusion",
			}), false, nil
		}
		if !list.IsComplete {
			if e.mode == Fixture && rule.FixtureOnly {
				return kleene.Known(false), true, nil
			}
			return kleene.Unk(kleene.Reason{
				Kind:   kleene.ListIncomplete,
				Ref:    list.ID,
				Detail: "the list is semantically closed but is reproduced here as a subset (is_complete: false)",
			}), false, nil
		}
		return kleene.Known(false), false, nil

	default:
		return kleene.Value{}, false, fmt.Errorf("rule %s: list predicate returned %q", rule.ID, tri)
	}
}

func (e *Engine) requireList(rule *spec.Rule, n *spec.Node) (*spec.CodeList, error) {
	if n.ListRef == "" {
		return nil, fmt.Errorf("rule %s: predicate %q has no list_ref (linter L-03)", rule.ID, n.Predicate)
	}
	list, ok := e.corpus.Lists[n.ListRef]
	if !ok {
		return nil, fmt.Errorf("rule %s: list_ref %q does not resolve (linter L-03)", rule.ID, n.ListRef)
	}
	if n.ListSemantics != "" && n.ListSemantics != list.ListSemantics {
		return nil, fmt.Errorf("rule %s: predicate declares list_semantics %q but %s declares %q "+
			"(linter L-03)", rule.ID, n.ListSemantics, list.ID, list.ListSemantics)
	}
	if n.ListVintage != "" && n.ListVintage != list.ListVintage {
		return nil, fmt.Errorf("rule %s: predicate pins list_vintage %q but %s is vintage %q "+
			"(linter L-03)", rule.ID, n.ListVintage, list.ID, list.ListVintage)
	}
	return list, nil
}

// ---------------------------------------------------------------- custom functions

func (e *Engine) codeInFn(n *spec.Node, list *spec.CodeList, tev *TranslationEvidence) cel.EnvOption {
	return cel.Function("pp_code_in",
		cel.Overload("pp_code_in_map",
			[]*cel.Type{cel.MapType(cel.StringType, cel.StringType)},
			cel.StringType,
			cel.UnaryBinding(func(arg ref.Val) ref.Val {
				m, err := stringMap(arg)
				if err != nil {
					return types.NewErr("classification: %v", err)
				}
				scheme, code, vintage := m["scheme"], m["code"], m["vintage"]
				if code == "" {
					return types.NewErr("classification has no code")
				}
				if n.Scheme != "" && scheme != n.Scheme {
					// A cross-scheme hop needs a concordance this engine does not have.
					tev.Composable = false
					return types.String(triNotComposable)
				}
				if scheme != list.Scheme {
					tev.Composable = false
					return types.String(triNotComposable)
				}
				codes, composable := e.crosswalk.Translate(scheme, code, vintage, list.ListVintage)
				tev.Composable = composable
				tev.ToCodes = codes
				if !composable {
					return types.String(triNotComposable)
				}
				for _, c := range codes {
					if listContainsHierarchical(list, c) {
						return types.String(triTrue)
					}
				}
				return types.String(triFalse)
			}),
		),
	)
}

func (e *Engine) anyInFn(list *spec.CodeList) cel.EnvOption {
	return cel.Function("pp_any_in",
		cel.Overload("pp_any_in_list",
			[]*cel.Type{cel.ListType(cel.StringType)},
			cel.StringType,
			cel.UnaryBinding(func(arg ref.Val) ref.Val {
				vals, err := stringSlice(arg)
				if err != nil {
					return types.NewErr("%v", err)
				}
				for _, v := range vals {
					for _, entry := range list.Codes {
						if entry.Code == v {
							return types.String(triTrue)
						}
					}
				}
				return types.String(triFalse)
			}),
		),
	)
}

func locationWithinFn(jp *spec.JurisdictionProfile) cel.EnvOption {
	main := append([]string(nil), jp.Main...)
	excepts := append([]string(nil), jp.Exceptions...)
	sort.Strings(main)
	sort.Strings(excepts)
	return cel.Function("pp_location_within",
		cel.Overload("pp_location_within_list",
			[]*cel.Type{cel.ListType(cel.StringType)},
			cel.BoolType,
			cel.UnaryBinding(func(arg ref.Val) ref.Val {
				path, err := stringSlice(arg)
				if err != nil {
					return types.NewErr("%v", err)
				}
				in := make(map[string]struct{}, len(path))
				for _, p := range path {
					in[p] = struct{}{}
				}
				for _, ex := range excepts {
					if _, hit := in[ex]; hit {
						return types.Bool(false)
					}
				}
				for _, m := range main {
					if _, hit := in[m]; !hit {
						return types.Bool(false)
					}
				}
				return types.Bool(true)
			}),
		),
	)
}

// listContainsHierarchical reports membership under the hierarchy of the scheme: a
// list entry at a shorter level covers every code beneath it. NAICS 5412 on OSHA's
// Appendix A covers 541211, and a flat string comparison would answer "not exempt"
// for every six-digit accountancy establishment in the country.
func listContainsHierarchical(list *spec.CodeList, code string) bool {
	for _, entry := range list.Codes {
		if entry.Code == code {
			return true
		}
		if len(entry.Code) < len(code) && strings.HasPrefix(code, entry.Code) {
			return true
		}
	}
	return false
}

func comparator(pred string) string {
	switch pred {
	case "attr_lte":
		return "<="
	case "attr_lt":
		return "<"
	case "attr_gte":
		return ">="
	case "attr_gt":
		return ">"
	default:
		return "=="
	}
}
