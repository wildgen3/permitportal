// Package credential resolves credential requirement trees into the order a subject
// must actually obtain things in.
//
// CTDL supplies the vocabulary but ships no OWL axioms and no SHACL, so ordering and
// boolean evaluation are not specified by it and are implemented here (ADR-0012).
// Cycle detection exists in scripts/check-rules.py as a gate on committed data; this
// is the runtime equivalent, because a corpus assembled from more than one file at
// runtime is not the corpus the linter saw.
package credential

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// Step is one credential in the resolved order.
type Step struct {
	Order         int      `json:"order"`
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	Type          string   `json:"type"`
	Authority     string   `json:"issuing_authority"`
	Citation      string   `json:"citation,omitempty"`
	SourceURL     string   `json:"source_url,omitempty"`
	Prerequisites []string `json:"prerequisites,omitempty"`

	// Conditions are the non-credential requirements — a bond, an examination, a
	// supervised-hours pathway. They are rendered as a requirement tree rather than a
	// flat list because `electrician-certification` has a genuine OR: an
	// experience-hours path or an education-plus-reduced-hours path. A flat foreign
	// key cannot express that.
	Conditions *Condition `json:"conditions,omitempty"`
}

// Condition mirrors the authored requirement tree, minus the credential edges that
// the ordering already accounts for.
type Condition struct {
	Kind            string       `json:"kind"` // AND_GROUP | OR_GROUP | LEAF
	LegalSource     string       `json:"legal_source,omitempty"`
	Predicate       string       `json:"predicate,omitempty"`
	YearsExperience int          `json:"years_experience,omitempty"`
	Children        []*Condition `json:"children,omitempty"`
}

// Chain is a resolved credential order.
type Chain struct {
	Targets []string `json:"targets"`
	Steps   []Step   `json:"steps"`
}

// Resolve returns the transitive closure of targets in an order where every
// credential appears after everything it depends on.
//
// Ties are broken by identifier, not by map or file order. Two runs over the same
// corpus must produce the same chain, or a determination that embeds one is not
// reproducible.
func Resolve(c *spec.Corpus, targets []string) (*Chain, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("credential: no target credentials")
	}

	// Transitive closure.
	needed := map[string]bool{}
	var visit func(id string, path []string) error
	visit = func(id string, path []string) error {
		for _, seen := range path {
			if seen == id {
				return fmt.Errorf("credential: cycle %s -> %s", strings.Join(path, " -> "), id)
			}
		}
		cred, ok := c.Credentials[id]
		if !ok {
			return fmt.Errorf("credential: %q is required but not declared in spec/credentials/", id)
		}
		if needed[id] {
			return nil
		}
		needed[id] = true
		deps := cred.Prerequisites()
		sort.Strings(deps)
		for _, d := range deps {
			if err := visit(d, append(path, id)); err != nil {
				return err
			}
		}
		return nil
	}
	sortedTargets := append([]string(nil), targets...)
	sort.Strings(sortedTargets)
	for _, t := range sortedTargets {
		if err := visit(t, nil); err != nil {
			return nil, err
		}
	}

	// Kahn's algorithm over the closure, with a sorted ready set.
	indegree := map[string]int{}
	dependents := map[string][]string{}
	for id := range needed {
		if _, ok := indegree[id]; !ok {
			indegree[id] = 0
		}
		for _, dep := range uniqueSorted(c.Credentials[id].Prerequisites()) {
			if !needed[dep] {
				continue
			}
			indegree[id]++
			dependents[dep] = append(dependents[dep], id)
		}
	}

	var ready []string
	for id, deg := range indegree {
		if deg == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)

	chain := &Chain{Targets: sortedTargets}
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		cred := c.Credentials[id]
		chain.Steps = append(chain.Steps, Step{
			Order:         len(chain.Steps) + 1,
			ID:            cred.ID,
			Label:         cred.Label,
			Type:          cred.Type,
			Authority:     cred.IssuingAuthority,
			Citation:      cred.Citation,
			SourceURL:     cred.SourceURL,
			Prerequisites: uniqueSorted(cred.Prerequisites()),
			Conditions:    conditionsOf(cred),
		})
		next := append([]string(nil), dependents[id]...)
		sort.Strings(next)
		for _, d := range next {
			indegree[d]--
			if indegree[d] == 0 {
				ready = append(ready, d)
			}
		}
		sort.Strings(ready)
	}

	if len(chain.Steps) != len(needed) {
		// Unreachable if the corpus passed the linter's cycle check, which is exactly
		// why it is asserted here: this runs against whatever was loaded, not against
		// whatever was linted.
		return nil, fmt.Errorf("credential: ordered %d of %d credentials — the graph contains a cycle",
			len(chain.Steps), len(needed))
	}
	return chain, nil
}

// conditionsOf strips credential edges from the requirement tree, leaving the
// conditions that are not themselves credentials. A group left with one child is
// collapsed; a group left with none disappears.
func conditionsOf(cred spec.Credential) *Condition {
	var conv func(spec.Requirement) *Condition
	conv = func(r spec.Requirement) *Condition {
		if r.NodeType == spec.Leaf {
			if r.TargetCredential != "" {
				return nil
			}
			return &Condition{
				Kind:            string(spec.Leaf),
				LegalSource:     r.LegalSource,
				Predicate:       r.TargetPredicate,
				YearsExperience: r.YearsExperience,
			}
		}
		var kids []*Condition
		for _, child := range r.Children {
			if c := conv(child); c != nil {
				kids = append(kids, c)
			}
		}
		switch len(kids) {
		case 0:
			return nil
		case 1:
			// An AND of one thing is that thing. An OR of one *surviving* alternative
			// is not: the discarded alternatives were credentials, and collapsing
			// would present a choice as a mandate.
			if r.NodeType == spec.AndGroup {
				return kids[0]
			}
		}
		return &Condition{Kind: string(r.NodeType), LegalSource: r.LegalSource, Children: kids}
	}

	var roots []*Condition
	for _, r := range cred.Requirements {
		if c := conv(r); c != nil {
			roots = append(roots, c)
		}
	}
	switch len(roots) {
	case 0:
		return nil
	case 1:
		return roots[0]
	default:
		return &Condition{Kind: string(spec.AndGroup), Children: roots}
	}
}

func uniqueSorted(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
