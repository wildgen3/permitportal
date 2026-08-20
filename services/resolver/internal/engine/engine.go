// Package engine evaluates authored rules against a subject's facts and returns a
// determination with the evidence tree that produced it.
//
// ADR-0001: this is the system of record. It has, and must keep, zero dependencies on
// any model client — scripts/check-engine-purity.py asserts that against the import
// graph on every pull request. If the AI layer ever becomes load-bearing in a
// determination, the build breaks rather than the determination becoming
// unreproducible.
package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/wildgen3/permitportal/services/resolver/internal/kleene"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// Version identifies the evaluator. It is pinned into every determination, because a
// determination is only reproducible if the thing that produced it is identified.
const Version = "0.1.0"

// Mode separates real determinations from specification fixtures.
//
// The distinction is load-bearing, not hygiene. Two of this repository's own
// invariants — that an incomplete code list cannot support a FALSE, and that a code
// may not cross vintages without a crosswalk — would each block the worked examples
// that demonstrate the engine. Rather than weakening the invariants, Fixture mode
// suspends them explicitly, refuses to touch any rule that has not declared itself a
// fixture, and stamps itself on every determination it produces.
type Mode string

const (
	Production Mode = "production"
	Fixture    Mode = "fixture"
)

type Outcome string

const (
	OutcomeTrue                       Outcome = "TRUE"
	OutcomeFalse                      Outcome = "FALSE"
	OutcomeUnknownMissingInput        Outcome = "UNKNOWN_MISSING_INPUT"
	OutcomeUnknownNotSelfDeterminable Outcome = "UNKNOWN_NOT_SELF_DETERMINABLE"
)

type Engine struct {
	corpus    *spec.Corpus
	mode      Mode
	crosswalk Crosswalk
}

func New(c *spec.Corpus, mode Mode) (*Engine, error) {
	if c == nil {
		return nil, fmt.Errorf("engine: nil corpus")
	}
	cw, err := crosswalkFor(mode)
	if err != nil {
		return nil, err
	}
	return &Engine{corpus: c, mode: mode, crosswalk: cw}, nil
}

func (e *Engine) Mode() Mode { return e.mode }

// SurfacedObligation is an obligation attached to the determination rather than
// derived from the tree.
type SurfacedObligation struct {
	ID          string `json:"id"`
	Citation    string `json:"citation,omitempty"`
	Note        string `json:"note,omitempty"`
	NonWaivable bool   `json:"non_waivable"`
}

// Determination is the output of record.
type Determination struct {
	RuleID           string `json:"rule_id"`
	RuleVersion      int    `json:"rule_version"`
	EngineVersion    string `json:"engine_version"`
	Mode             Mode   `json:"mode"`
	AsOfLaw          string `json:"as_of_law"`
	InputSnapshotSHA string `json:"input_snapshot_sha256"`

	Outcome Outcome `json:"outcome"`
	// Reason is the label of the node that settled the determination, when the rule
	// labels its branches. It is what a plain-language explanation leads with.
	Reason string `json:"reason,omitempty"`

	// MissingAttributes are the questions to ask next. Populated only for
	// UNKNOWN_MISSING_INPUT — this is part of the contract, not an error payload.
	MissingAttributes []string `json:"missing_attributes,omitempty"`

	// RouteToAuthority carries the reasons the subject cannot resolve themselves.
	RouteToAuthority []kleene.Reason `json:"route_to_authority,omitempty"`

	Evidence *Evidence                `json:"evidence"`
	Surfaced []SurfacedObligation     `json:"surfaced_obligations,omitempty"`
	Produces *spec.ProducesObligation `json:"produces_obligation,omitempty"`
}

// Evaluate runs one rule against one profile as the law stood on asOfLaw.
func (e *Engine) Evaluate(ruleID string, p Profile, asOfLaw string) (*Determination, error) {
	rule, ok := e.corpus.Rules[ruleID]
	if !ok {
		return nil, fmt.Errorf("engine: no rule %q in the corpus", ruleID)
	}
	if asOfLaw == "" {
		return nil, fmt.Errorf("engine: as_of_law is required; a determination with no law date " +
			"cannot be reproduced (ADR-0009)")
	}
	if rule.FixtureOnly && e.mode != Fixture {
		return nil, fmt.Errorf("engine: rule %s declares fixture_only: true and cannot be evaluated "+
			"in %s mode.\n  It references a code list reproduced here as a subset, so a FALSE from it "+
			"would be an assertion the source data does not support (linter L-06).", ruleID, e.mode)
	}
	if err := p.Validate(e.corpus); err != nil {
		return nil, err
	}
	if err := inForce(rule, asOfLaw); err != nil {
		return nil, err
	}

	ev, err := e.eval(rule, &rule.Expr, p, map[string]bool{rule.ID: true})
	if err != nil {
		return nil, err
	}

	val := kleene.Value{T: truthFrom(ev.Truth), Reasons: ev.Reasons}
	val = e.applyOverlays(rule, val, p)

	hash, err := snapshotSHA(p)
	if err != nil {
		return nil, err
	}

	d := &Determination{
		RuleID:           rule.ID,
		RuleVersion:      rule.Version,
		EngineVersion:    Version,
		Mode:             e.mode,
		AsOfLaw:          asOfLaw,
		InputSnapshotSHA: hash,
		Evidence:         ev,
		Produces:         rule.ProducesObligation,
	}

	switch val.T {
	case kleene.True:
		d.Outcome = OutcomeTrue
		d.Reason = ev.decisive()
	case kleene.False:
		d.Outcome = OutcomeFalse
		d.Reason = ev.decisive()
	default:
		// Route-to-authority wins over ask-the-applicant when both are present: you
		// cannot close a gap in the system's own knowledge by asking the subject a
		// question, so offering them a form field would be a dead end.
		var missing []string
		var route []kleene.Reason
		for _, r := range val.Reasons {
			if r.Kind.AskableByApplicant() {
				missing = append(missing, r.Ref)
			} else {
				route = append(route, r)
			}
		}
		sort.Strings(missing)
		if len(route) > 0 {
			d.Outcome = OutcomeUnknownNotSelfDeterminable
			d.RouteToAuthority = route
			d.MissingAttributes = missing
		} else {
			d.Outcome = OutcomeUnknownMissingInput
			d.MissingAttributes = missing
		}
		// The tree's own annotation must agree with the determination.
		ev.Truth = kleene.Unknown.String()
		ev.Reasons = val.Reasons
	}

	// Non-waivable obligations survive every outcome, which is the whole point of
	// attaching them to the rule. An "exempt" answer that omits severe-injury
	// reporting is a harmful answer.
	for _, nw := range rule.NonWaivable {
		d.Surfaced = append(d.Surfaced, SurfacedObligation{
			ID:          nw.Obligation,
			Citation:    nw.Cite,
			Note:        nw.Note,
			NonWaivable: true,
		})
	}
	return d, nil
}

// eval walks the authored tree, producing the evidence tree as it goes.
func (e *Engine) eval(rule *spec.Rule, n *spec.Node, p Profile, visiting map[string]bool) (*Evidence, error) {
	ev := &Evidence{Kind: n.Kind, Label: n.Label}

	switch n.Kind {
	case spec.KindAll, spec.KindAny, spec.KindNone:
		vals := make([]kleene.Value, 0, len(n.Children))
		for i := range n.Children {
			child, err := e.eval(rule, &n.Children[i], p, visiting)
			if err != nil {
				return nil, err
			}
			ev.Children = append(ev.Children, child)
			vals = append(vals, kleene.Value{T: truthFrom(child.Truth), Reasons: child.Reasons})
		}
		var v kleene.Value
		switch n.Kind {
		case spec.KindAll:
			v = kleene.All(vals)
		case spec.KindAny:
			v = kleene.Any(vals)
		default:
			v = kleene.None(vals)
		}
		ev.Truth, ev.Reasons = v.T.String(), v.Reasons
		return ev, nil

	case spec.KindNot:
		child, err := e.eval(rule, n.Child, p, visiting)
		if err != nil {
			return nil, err
		}
		ev.Children = []*Evidence{child}
		v := kleene.Not(kleene.Value{T: truthFrom(child.Truth), Reasons: child.Reasons})
		ev.Truth, ev.Reasons = v.T.String(), v.Reasons
		return ev, nil

	case spec.KindRef:
		target, ok := e.corpus.Rules[n.RefRule]
		if !ok {
			return nil, fmt.Errorf("rule %s: ref %q does not resolve", rule.ID, n.RefRule)
		}
		if visiting[target.ID] {
			return nil, fmt.Errorf("rule %s: cycle through ref %q", rule.ID, n.RefRule)
		}
		visiting[target.ID] = true
		defer delete(visiting, target.ID)
		child, err := e.eval(target, &target.Expr, p, visiting)
		if err != nil {
			return nil, err
		}
		ev.Children = []*Evidence{child}
		ev.Truth, ev.Reasons = child.Truth, child.Reasons
		return ev, nil

	case spec.KindPredicate:
		out, err := e.evalLeaf(rule, n, p)
		if err != nil {
			return nil, err
		}
		ev.Predicate = n.Predicate
		ev.Scope = n.Scope
		ev.Attribute = n.Attribute
		ev.ListRef = n.ListRef
		ev.Truth = out.Value.T.String()
		ev.Reasons = out.Value.Reasons
		ev.Translation = out.Translation
		ev.ListCompletenessWaived = out.Waived
		for _, id := range n.Cites {
			if c, ok := rule.Cite(id); ok {
				ev.Citations = append(ev.Citations, c)
			} else {
				return nil, fmt.Errorf("rule %s: predicate cites %q, which is not declared", rule.ID, id)
			}
		}
		return ev, nil
	}
	return nil, fmt.Errorf("rule %s: unhandled node kind %q", rule.ID, n.Kind)
}

// applyOverlays downgrades a result an overlay says cannot yet be asserted.
//
// The federal recordkeeping appendix is a floor: a state plan may require records
// from establishments the federal rule exempts. A rule carrying
// effect: MAY_NARROW / resolution: REQUIRE_JURISDICTION_RULE therefore may not answer
// "exempt" until the jurisdiction fact is actually loaded — regardless of which key
// produced the exemption. Without this the size key would answer TRUE while the
// industry key correctly answered INDETERMINATE on the same profile.
func (e *Engine) applyOverlays(rule *spec.Rule, v kleene.Value, p Profile) kleene.Value {
	if v.T != kleene.True || rule.ResultSemantics != "exemption" {
		return v
	}
	var reasons []kleene.Reason
	for _, o := range rule.Overlays {
		if o.Effect != "MAY_NARROW" || o.Resolution != "REQUIRE_JURISDICTION_RULE" {
			continue
		}
		for _, uri := range jurisdictionRuleAttrs(e.corpus, rule) {
			attr := e.corpus.Attributes[uri]
			if _, present := p.Lookup(attr); !present {
				reasons = append(reasons, kleene.Reason{
					Kind:   kleene.JurisdictionRuleMissing,
					Ref:    uri,
					Detail: fmt.Sprintf("overlay %q may narrow this exemption; it is INDETERMINATE, never exempt, until the jurisdiction rule is loaded", o.Scope),
				})
			}
		}
	}
	if len(reasons) == 0 {
		return v
	}
	return kleene.Unk(reasons...)
}

// jurisdictionRuleAttrs returns the attributes this rule reads that are supplied by a
// jurisdiction rule rather than by the applicant.
func jurisdictionRuleAttrs(c *spec.Corpus, rule *spec.Rule) []string {
	seen := map[string]bool{}
	var out []string
	var walk func(*spec.Node)
	walk = func(n *spec.Node) {
		if n == nil {
			return
		}
		if n.Kind == spec.KindPredicate {
			if a, ok := c.Attributes[n.Attribute]; ok &&
				a.CollectionMethod == "jurisdiction_rule" && !seen[a.URI] {
				seen[a.URI] = true
				out = append(out, a.URI)
			}
		}
		for i := range n.Children {
			walk(&n.Children[i])
		}
		walk(n.Child)
	}
	walk(&rule.Expr)
	sort.Strings(out)
	return out
}

// ValidateOverlays checks, at load time, that a rule whose overlay demands a
// jurisdiction rule actually reads one. A rule that demands a fact it never consults
// would pass every test while enforcing nothing.
func ValidateOverlays(c *spec.Corpus) error {
	var problems []string
	for _, id := range c.RuleIDs {
		rule := c.Rules[id]
		for _, o := range rule.Overlays {
			if o.Resolution != "REQUIRE_JURISDICTION_RULE" {
				continue
			}
			if len(jurisdictionRuleAttrs(c, rule)) == 0 {
				problems = append(problems, fmt.Sprintf(
					"%s: overlay %q requires a jurisdiction rule, but the rule reads no attribute "+
						"whose collection_method is jurisdiction_rule, so the requirement is inert",
					rule.ID, o.Scope))
			}
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("overlay validation failed:\n  - %s", joinLines(problems))
	}
	return nil
}

func joinLines(s []string) string {
	out := ""
	for i, v := range s {
		if i > 0 {
			out += "\n  - "
		}
		out += v
	}
	return out
}

func inForce(rule *spec.Rule, asOfLaw string) error {
	if f := rule.Effective.LawFrom; f != "" && asOfLaw < f {
		return fmt.Errorf("rule %s took effect %s; as_of_law is %s", rule.ID, f, asOfLaw)
	}
	if t := rule.Effective.LawTo; t != "" && asOfLaw > t {
		return fmt.Errorf("rule %s ceased to have effect %s; as_of_law is %s", rule.ID, t, asOfLaw)
	}
	return nil
}

func truthFrom(s string) kleene.Truth {
	switch s {
	case "TRUE":
		return kleene.True
	case "FALSE":
		return kleene.False
	default:
		return kleene.Unknown
	}
}

// snapshotSHA hashes the fact set the determination was made against.
//
// encoding/json sorts map keys, so the digest does not depend on Go's randomised map
// iteration order. Replay needs this to be stable across processes, not merely within
// one.
func snapshotSHA(p Profile) (string, error) {
	b, err := json.Marshal(p.Facts)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
