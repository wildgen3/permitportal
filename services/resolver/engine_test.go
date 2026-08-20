package resolver_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/wildgen3/permitportal/services/resolver/internal/credential"
	"github.com/wildgen3/permitportal/services/resolver/internal/engine"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

const oshaRule = "osha.1904.partial_exemption"

func exemptBySizeProfile() engine.Profile {
	return engine.Profile{Facts: map[string]map[string]any{
		"business": {"company_peak_employment_prior_cy": 8},
		"establishment": {
			"classification": map[string]any{"scheme": "NAICS", "vintage": "2022", "code": "238210"},
			"osha_bls_state_written_exemption_override": false,
		},
	}}
}

func fixtureEngine(t *testing.T) (*engine.Engine, *spec.Corpus) {
	t.Helper()
	corpus, _ := loadCorpus(t)
	eng, err := engine.New(corpus, engine.Fixture)
	if err != nil {
		t.Fatal(err)
	}
	return eng, corpus
}

// Determinations must replay byte-identically from
// (rule_version, engine_version, as_of_law, input_snapshot_hash). Anything that
// varies between runs — map iteration order, a timestamp, a hash seed — breaks the
// audit story, and this repository has already been bitten by a generator that was
// non-deterministic in exactly that way.
func TestDeterminismWithinAProcess(t *testing.T) {
	eng, _ := fixtureEngine(t)
	var first string
	for i := 0; i < 200; i++ {
		d, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19")
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			first = string(b)
			continue
		}
		if string(b) != first {
			t.Fatalf("determination %d differs from the first:\n  %s\n  %s", i, first, string(b))
		}
	}
	if first == "" {
		t.Fatal("no determination was produced")
	}
}

// And across independently loaded corpora, because replay happens in a different
// process than the original determination.
func TestDeterminismAcrossCorpusLoads(t *testing.T) {
	var seen []string
	for i := 0; i < 5; i++ {
		eng, _ := fixtureEngine(t)
		d, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19")
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(d)
		seen = append(seen, string(b))
	}
	for i := 1; i < len(seen); i++ {
		if seen[i] != seen[0] {
			t.Fatalf("load %d produced a different determination", i)
		}
	}
}

func TestSnapshotHashChangesWithTheFacts(t *testing.T) {
	eng, _ := fixtureEngine(t)
	a, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19")
	if err != nil {
		t.Fatal(err)
	}
	p := exemptBySizeProfile()
	p.Facts["business"]["company_peak_employment_prior_cy"] = 9
	b, err := eng.Evaluate(oshaRule, p, "2026-08-19")
	if err != nil {
		t.Fatal(err)
	}
	if a.InputSnapshotSHA == b.InputSnapshotSHA {
		t.Fatal("the snapshot hash did not change when the facts changed — replay would " +
			"reproduce a determination against the wrong inputs")
	}
}

// The two-keys-two-granularities trap, asserted rather than described. Swapping the
// scope on the size predicate must be refused, not silently evaluated against
// whatever happens to be at that scope.
func TestScopeMismatchIsRefused(t *testing.T) {
	eng, corpus := fixtureEngine(t)
	rule := corpus.Rules[oshaRule]
	original := rule.Expr.Children[0].Scope
	rule.Expr.Children[0].Scope = "establishment"
	defer func() { rule.Expr.Children[0].Scope = original }()

	_, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19")
	if err == nil {
		t.Fatal("a company-scoped attribute evaluated at establishment scope was accepted")
	}
	for _, want := range []string{"company.peak_employment_prior_cy", "establishment", "business"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error does not name %q, so it does not tell the author what to fix:\n%v", want, err)
		}
	}
}

// A profile key that resolves to no registered attribute must be an error. If it read
// as "absent" instead, a typo would come back as INDETERMINATE and the product would
// ask a question the user already answered.
func TestUnregisteredProfileKeyIsRefused(t *testing.T) {
	eng, _ := fixtureEngine(t)
	p := exemptBySizeProfile()
	p.Facts["business"]["company_peak_employmnet_prior_cy"] = 8 // transposed
	_, err := eng.Evaluate(oshaRule, p, "2026-08-19")
	if err == nil {
		t.Fatal("an unregistered profile key was accepted")
	}
	if !strings.Contains(err.Error(), "company_peak_employmnet_prior_cy") {
		t.Errorf("the error does not name the offending key:\n%v", err)
	}
}

// A code miss against an illustrative list may never be a FALSE. This mutates the
// loaded list rather than authoring a second rule, so the assertion is about the
// evaluator rather than about a fixture.
func TestOpenListNeverProducesFalse(t *testing.T) {
	eng, corpus := fixtureEngine(t)
	list := corpus.Lists["lists/osha_1904_app_a"]
	list.ListSemantics = spec.SemanticsOpen
	list.IsComplete = true // isolate polarity from completeness
	rule := corpus.Rules[oshaRule]
	node := &rule.Expr.Children[1].Children[0]
	originalSemantics := node.ListSemantics
	node.ListSemantics = spec.SemanticsOpen
	defer func() {
		list.ListSemantics = spec.SemanticsClosed
		list.IsComplete = false
		node.ListSemantics = originalSemantics
	}()

	// 240 employees defeats the size key, so the industry key alone decides. 238210 is
	// not on the list.
	p := exemptBySizeProfile()
	p.Facts["business"]["company_peak_employment_prior_cy"] = 240

	d, err := eng.Evaluate(oshaRule, p, "2026-08-19")
	if err != nil {
		t.Fatal(err)
	}
	if d.Outcome == engine.OutcomeFalse {
		t.Fatal("a miss against an illustrative_open list produced FALSE — the code miss " +
			"was treated as proof of exclusion")
	}
	if d.Outcome != engine.OutcomeUnknownNotSelfDeterminable {
		t.Fatalf("outcome = %s, want UNKNOWN_NOT_SELF_DETERMINABLE", d.Outcome)
	}
}

// Production mode has no crosswalk data, so a rule that translates a code of record
// into an older list vintage cannot be evaluated. The refusal must be explicit.
func TestProductionRefusesAFixtureOnlyRule(t *testing.T) {
	corpus, _ := loadCorpus(t)
	eng, err := engine.New(corpus, engine.Production)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19"); err == nil {
		t.Fatal("production mode produced a determination from a fixture_only rule")
	}
}

func TestRuleOutsideItsEffectiveWindowIsRefused(t *testing.T) {
	eng, _ := fixtureEngine(t)
	// The NAICS-based appendix was adopted 2014-09-18; before that the appendix keyed
	// to SIC, so a 2010 determination from this rule would be an anachronism.
	if _, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2010-01-01"); err == nil {
		t.Fatal("a rule was evaluated before it took effect")
	}
}

func TestEvidenceTreeCarriesCitations(t *testing.T) {
	eng, _ := fixtureEngine(t)
	d, err := eng.Evaluate(oshaRule, exemptBySizeProfile(), "2026-08-19")
	if err != nil {
		t.Fatal(err)
	}
	var cited, leaves int
	var walk func(*engine.Evidence)
	walk = func(e *engine.Evidence) {
		if e == nil {
			return
		}
		if e.Predicate != "" {
			leaves++
			if len(e.Citations) > 0 {
				cited++
			}
		}
		for _, c := range e.Children {
			walk(c)
		}
	}
	walk(d.Evidence)
	if leaves == 0 {
		t.Fatal("the evidence tree has no leaves")
	}
	if cited != leaves {
		t.Fatalf("%d of %d predicate nodes carry a citation; a determination without one on "+
			"every contributing node is announced rather than arguable", cited, leaves)
	}
}

func TestCredentialChainIsOrderedAndStable(t *testing.T) {
	_, corpus := fixtureEngine(t)
	chain, err := credential.Resolve(corpus, []string{"cred/us-wa/electrical-contractor"})
	if err != nil {
		t.Fatal(err)
	}
	pos := map[string]int{}
	for _, s := range chain.Steps {
		pos[s.ID] = s.Order
	}
	for _, s := range chain.Steps {
		for _, dep := range s.Prerequisites {
			if pos[dep] == 0 {
				t.Fatalf("%s depends on %s, which is not in the chain", s.ID, dep)
			}
			if pos[dep] >= s.Order {
				t.Fatalf("%s (order %d) depends on %s (order %d) — the order is not topological",
					s.ID, s.Order, dep, pos[dep])
			}
		}
	}
	// The AND-of-ORs must survive. electrician-certification has a genuine choice
	// between an experience path and an education path, and flattening it would
	// present a choice as a mandate.
	var found bool
	for _, s := range chain.Steps {
		if s.ID != "cred/us-wa/electrician-certification" {
			continue
		}
		found = true
		if s.Conditions == nil {
			t.Fatal("electrician-certification lost its requirement tree")
		}
		if !containsOrGroup(s.Conditions) {
			t.Fatalf("the OR pathway was flattened: %+v", s.Conditions)
		}
	}
	if !found {
		t.Fatal("electrician-certification is not in the chain")
	}

	// Stable across runs.
	first, _ := json.Marshal(chain)
	for i := 0; i < 20; i++ {
		again, err := credential.Resolve(corpus, []string{"cred/us-wa/electrical-contractor"})
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(again)
		if string(b) != string(first) {
			t.Fatal("the credential chain is not stable across runs")
		}
	}
}

func containsOrGroup(c *credential.Condition) bool {
	if c == nil {
		return false
	}
	if c.Kind == string(spec.OrGroup) {
		return true
	}
	for _, k := range c.Children {
		if containsOrGroup(k) {
			return true
		}
	}
	return false
}

func TestCredentialCycleIsDetectedAtRuntime(t *testing.T) {
	_, corpus := fixtureEngine(t)
	// Make the root of the chain depend on its own dependent.
	base := corpus.Credentials["cred/us-wa/business-license"]
	base.Requirements = []spec.Requirement{{
		ID:       "req/injected/cycle",
		NodeType: spec.AndGroup,
		Children: []spec.Requirement{{
			ID:               "req/injected/cycle-leaf",
			NodeType:         spec.Leaf,
			TargetCredential: "cred/us-wa/electrical-contractor",
		}},
	}}
	corpus.Credentials["cred/us-wa/business-license"] = base

	if _, err := credential.Resolve(corpus, []string{"cred/us-wa/electrical-contractor"}); err == nil {
		t.Fatal("a cycle in the credential graph was ordered without complaint")
	}
}

func TestUndeclaredCredentialIsRefused(t *testing.T) {
	_, corpus := fixtureEngine(t)
	if _, err := credential.Resolve(corpus, []string{"cred/us-wa/does-not-exist"}); err == nil {
		t.Fatal("an undeclared credential resolved")
	}
}
