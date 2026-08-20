// Package kleene implements the three-valued logic the rules DSL is specified in
// (ADR-0006, docs/06-rules-dsl.md).
//
// The single property that matters: MISSING INPUT NEVER BECOMES FALSE. In a
// compliance product the false negative is the harm — a rule that silently answers
// "does not apply" because nobody asked about chemicals is the worst output this
// system can produce. Every combinator below is written so that the only way to
// reach False is for a child to have actually been evaluated to False.
package kleene

import "sort"

// Truth is a Kleene truth value.
type Truth uint8

const (
	False Truth = iota
	True
	Unknown
)

func (t Truth) String() string {
	switch t {
	case True:
		return "TRUE"
	case False:
		return "FALSE"
	default:
		return "UNKNOWN"
	}
}

// ReasonKind classifies *why* a value is Unknown. The distinction is not cosmetic:
// it decides which of the two UNKNOWN outcomes the determination surfaces, and
// therefore whether the product asks the applicant a question or routes them to the
// issuing authority.
type ReasonKind string

const (
	// The applicant can answer this. Drives missing_attributes[] -> next question.
	MissingAttribute ReasonKind = "missing_attribute"

	// The applicant cannot answer any of these. They are gaps in the system's own
	// knowledge or in the regulation's self-determinability, and they all route to
	// the issuing authority rather than to a form field.
	NotSelfDeterminable     ReasonKind = "not_self_determinable"
	ListIncomplete          ReasonKind = "list_incomplete"
	ListIllustrativeOpen    ReasonKind = "list_illustrative_open"
	CrosswalkUnavailable    ReasonKind = "crosswalk_unavailable"
	JurisdictionRuleMissing ReasonKind = "jurisdiction_rule_missing"
)

// AskableByApplicant reports whether a reason can be resolved by asking the subject
// a question. Only one kind can.
func (k ReasonKind) AskableByApplicant() bool { return k == MissingAttribute }

// Reason explains one contribution to an Unknown.
type Reason struct {
	Kind   ReasonKind `json:"kind" yaml:"kind"`
	Ref    string     `json:"ref" yaml:"ref"`
	Detail string     `json:"detail,omitempty" yaml:"detail,omitempty"`
}

// Value is a truth value plus, when Unknown, the reasons it is Unknown.
type Value struct {
	T       Truth    `json:"truth" yaml:"truth"`
	Reasons []Reason `json:"reasons,omitempty" yaml:"reasons,omitempty"`
}

func Known(b bool) Value {
	if b {
		return Value{T: True}
	}
	return Value{T: False}
}

func Unk(reasons ...Reason) Value { return Value{T: Unknown, Reasons: dedupe(reasons)} }

// All is Kleene AND. FALSE dominates: one evaluated-false child settles the node even
// if siblings are unknown, because no amount of later information can rescue it.
func All(vs []Value) Value {
	var unknown []Reason
	sawUnknown := false
	for _, v := range vs {
		switch v.T {
		case False:
			return Value{T: False}
		case Unknown:
			sawUnknown = true
			unknown = append(unknown, v.Reasons...)
		}
	}
	if sawUnknown {
		return Unk(unknown...)
	}
	return Value{T: True}
}

// Any is Kleene OR. TRUE dominates, symmetrically.
func Any(vs []Value) Value {
	var unknown []Reason
	sawUnknown := false
	for _, v := range vs {
		switch v.T {
		case True:
			return Value{T: True}
		case Unknown:
			sawUnknown = true
			unknown = append(unknown, v.Reasons...)
		}
	}
	if sawUnknown {
		return Unk(unknown...)
	}
	return Value{T: False}
}

// Not flips a known value and passes an Unknown through with its reasons intact.
//
// This is the line that rejects negation-as-failure (ADR-0007). Under Rego's
// semantics `not p` succeeds when p is merely absent; here it stays Unknown.
func Not(v Value) Value {
	switch v.T {
	case True:
		return Value{T: False}
	case False:
		return Value{T: True}
	default:
		return Value{T: Unknown, Reasons: v.Reasons}
	}
}

// None is NOR: true only when every child is false.
func None(vs []Value) Value { return Not(Any(vs)) }

// dedupe returns reasons in a deterministic order with duplicates removed.
// Determinations must be byte-identically reproducible, so reason ordering cannot
// depend on evaluation order or on Go map iteration.
func dedupe(rs []Reason) []Reason {
	if len(rs) == 0 {
		return nil
	}
	seen := make(map[Reason]struct{}, len(rs))
	out := make([]Reason, 0, len(rs))
	for _, r := range rs {
		if _, dup := seen[r]; dup {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		if out[i].Ref != out[j].Ref {
			return out[i].Ref < out[j].Ref
		}
		return out[i].Detail < out[j].Detail
	})
	return out
}
