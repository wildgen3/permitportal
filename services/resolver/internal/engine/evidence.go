package engine

import (
	"github.com/wildgen3/permitportal/services/resolver/internal/kleene"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// Evidence is the authored rule tree annotated with truth values.
//
// It is not a rendering of the decision — it *is* the decision. docs/06-rules-dsl.md
// is explicit that if the tree is not the representation, the explanation becomes a
// separate artifact that drifts from what it purports to explain.
type Evidence struct {
	Kind  spec.NodeKind `json:"kind"`
	Label string        `json:"label,omitempty"`

	// Leaf detail
	Predicate string `json:"predicate,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ListRef   string `json:"list_ref,omitempty"`

	Truth   string          `json:"truth"`
	Reasons []kleene.Reason `json:"reasons,omitempty"`

	// Citations resolved from the node's cites:, carried per node because that is
	// where they attach in the authored form.
	Citations []spec.Citation `json:"citations,omitempty"`

	// Translation records the vintage hop actually performed on a code predicate.
	// Present only when a translation was attempted.
	Translation *TranslationEvidence `json:"translation,omitempty"`

	// ListCompletenessWaived is set when a FALSE was permitted from a list marked
	// is_complete: false. Only Fixture mode can set it. A production determination
	// with this flag set is a bug, and it is recorded so that it is visible.
	ListCompletenessWaived bool `json:"list_completeness_waived,omitempty"`

	Children []*Evidence `json:"children,omitempty"`
}

type TranslationEvidence struct {
	Crosswalk   string   `json:"crosswalk"`
	FromVintage string   `json:"from_vintage"`
	ToVintage   string   `json:"to_vintage"`
	FromCode    string   `json:"from_code"`
	ToCodes     []string `json:"to_codes,omitempty"`
	Composable  bool     `json:"composable"`
}

// decisive returns the label of the child that settled this node, if any. For an
// `any` that returned TRUE it is the first TRUE child; for an `all` that returned
// FALSE it is the first FALSE child. This is what `expect_reason` in the golden files
// names, and it is what a user-facing explanation leads with.
func (e *Evidence) decisive() string {
	if e == nil {
		return ""
	}
	switch e.Kind {
	case spec.KindAny:
		if e.Truth == kleene.True.String() {
			return firstChildWith(e, kleene.True.String())
		}
	case spec.KindAll:
		if e.Truth == kleene.False.String() {
			return firstChildWith(e, kleene.False.String())
		}
	}
	return e.Label
}

func firstChildWith(e *Evidence, truth string) string {
	for _, c := range e.Children {
		if c.Truth == truth {
			if c.Label != "" {
				return c.Label
			}
			return c.decisive()
		}
	}
	return e.Label
}
