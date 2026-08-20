// Package spec loads the authored specification — rules, code lists, the attribute
// registry, and credential chains — from spec/ into Go values.
//
// Nothing here interprets law. It is a faithful reader for the YAML that
// scripts/check-rules.py already validates, and it deliberately re-checks the
// invariants the evaluator depends on rather than trusting that the linter ran.
package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------- rules

// Citation is serialised into the evidence tree, so it carries JSON tags: the wire
// shape of a determination is part of the contract and must not be Go field names.
type Citation struct {
	ID     string `yaml:"id"     json:"id"`
	Cite   string `yaml:"cite"   json:"cite"`
	Source string `yaml:"source" json:"source"`
	AsOf   string `yaml:"as_of"  json:"as_of"`
	Note   string `yaml:"note"   json:"note,omitempty"`
}

type Effective struct {
	LawFrom string `yaml:"law_from"`
	LawTo   string `yaml:"law_to"`
}

type JurisdictionProfile struct {
	Main       []string `yaml:"main"`
	Exceptions []string `yaml:"exceptions"`
}

type Overlay struct {
	Scope      string `yaml:"scope"`
	Effect     string `yaml:"effect"`
	Resolution string `yaml:"resolution"`
	Note       string `yaml:"note"`
}

type NonWaivable struct {
	Obligation string `yaml:"obligation"`
	Cite       string `yaml:"cite"`
	Note       string `yaml:"note"`
}

type ProducesObligation struct {
	Type       string `yaml:"type"       json:"type"`
	Credential string `yaml:"credential" json:"credential,omitempty"`
}

type Fixture struct {
	ID   string `yaml:"id"`
	Note string `yaml:"note"`
}

type Fixtures struct {
	Positive []Fixture `yaml:"positive"`
	Negative []Fixture `yaml:"negative"`
}

type Rule struct {
	ID                  string               `yaml:"rule"`
	Version             int                  `yaml:"version"`
	Scope               string               `yaml:"scope"`
	ResultSemantics     string               `yaml:"result_semantics"`
	FixtureOnly         bool                 `yaml:"fixture_only"`
	Effective           Effective            `yaml:"effective"`
	SourceURL           string               `yaml:"source_url"`
	Citations           []Citation           `yaml:"citations"`
	Expr                Node                 `yaml:"expr"`
	Overlays            []Overlay            `yaml:"overlays"`
	NonWaivable         []NonWaivable        `yaml:"non_waivable"`
	JurisdictionProfile *JurisdictionProfile `yaml:"jurisdiction_profile"`
	ProducesObligation  *ProducesObligation  `yaml:"produces_obligation"`
	SoleKeyJustify      string               `yaml:"sole_key_justification"`
	Fixtures            Fixtures             `yaml:"fixtures"`

	// Path is where the rule was loaded from. Diagnostics only.
	Path string `yaml:"-"`
}

// Cite resolves a citation id declared on this rule.
func (r *Rule) Cite(id string) (Citation, bool) {
	for _, c := range r.Citations {
		if c.ID == id {
			return c, true
		}
	}
	return Citation{}, false
}

// ---------------------------------------------------------------- expression tree

type NodeKind string

const (
	KindAll       NodeKind = "all"
	KindAny       NodeKind = "any"
	KindNone      NodeKind = "none"
	KindNot       NodeKind = "not"
	KindRef       NodeKind = "ref"
	KindPredicate NodeKind = "predicate"
)

type Translation struct {
	FromCodeOfRecord  bool   `yaml:"from_code_of_record"`
	RequireComposable bool   `yaml:"require_composable"`
	OnCloseMatch      string `yaml:"on_close_match"`
}

// Node is one node of the authored predicate tree. The tree is the representation —
// not a compiled expression string — because citations, effective dates, list
// vintages and scope hops attach per node, and because the evidence tree returned to
// the user *is* this tree annotated with truth values.
type Node struct {
	Kind  NodeKind
	Label string

	// Interior
	Children []Node
	Child    *Node  // `not`
	RefRule  string // `ref`

	// Leaf
	Predicate           string
	Scope               string
	Attribute           string
	Value               yaml.Node
	Scheme              string
	ListRef             string
	ListVintage         string
	ListSemantics       string
	Translation         *Translation
	JurisdictionProfRef string
	Cites               []string
}

type rawNode struct {
	Label string `yaml:"label"`

	All  []Node `yaml:"all"`
	Any  []Node `yaml:"any"`
	None []Node `yaml:"none"`
	Not  *Node  `yaml:"not"`
	Ref  string `yaml:"ref"`

	Predicate           string       `yaml:"predicate"`
	Scope               string       `yaml:"scope"`
	Attribute           string       `yaml:"attribute"`
	Value               yaml.Node    `yaml:"value"`
	Scheme              string       `yaml:"scheme"`
	ListRef             string       `yaml:"list_ref"`
	ListVintage         string       `yaml:"list_vintage"`
	ListSemantics       string       `yaml:"list_semantics"`
	Translation         *Translation `yaml:"translation"`
	JurisdictionProfRef string       `yaml:"jurisdiction_profile_ref"`
	Cites               []string     `yaml:"cites"`
}

// UnmarshalYAML decodes the tagged-union node shape. Exactly one structural key may
// be present; a node carrying two is a malformed rule and is rejected here rather
// than being silently resolved by field order.
func (n *Node) UnmarshalYAML(value *yaml.Node) error {
	var raw rawNode
	if err := value.Decode(&raw); err != nil {
		return err
	}

	present := []NodeKind{}
	if raw.All != nil {
		present = append(present, KindAll)
	}
	if raw.Any != nil {
		present = append(present, KindAny)
	}
	if raw.None != nil {
		present = append(present, KindNone)
	}
	if raw.Not != nil {
		present = append(present, KindNot)
	}
	if raw.Ref != "" {
		present = append(present, KindRef)
	}
	if raw.Predicate != "" {
		present = append(present, KindPredicate)
	}

	switch len(present) {
	case 0:
		return fmt.Errorf("line %d: expression node has no structural key "+
			"(expected one of all, any, none, not, ref, predicate)", value.Line)
	case 1:
	default:
		return fmt.Errorf("line %d: expression node declares %v; exactly one is allowed",
			value.Line, present)
	}

	n.Kind = present[0]
	n.Label = raw.Label
	switch n.Kind {
	case KindAll:
		n.Children = raw.All
	case KindAny:
		n.Children = raw.Any
	case KindNone:
		n.Children = raw.None
	case KindNot:
		n.Child = raw.Not
	case KindRef:
		n.RefRule = raw.Ref
	case KindPredicate:
		n.Predicate = raw.Predicate
		n.Scope = raw.Scope
		n.Attribute = raw.Attribute
		n.Value = raw.Value
		n.Scheme = raw.Scheme
		n.ListRef = raw.ListRef
		n.ListVintage = raw.ListVintage
		n.ListSemantics = raw.ListSemantics
		n.Translation = raw.Translation
		n.JurisdictionProfRef = raw.JurisdictionProfRef
		n.Cites = raw.Cites
	}
	return nil
}

// ---------------------------------------------------------------- code lists

type ListEntry struct {
	Code             string `yaml:"code"`
	Title            string `yaml:"title"`
	RevisionStatus   string `yaml:"revision_status"`
	SuccessorCode    string `yaml:"successor_code"`
	SuccessorVintage string `yaml:"successor_vintage"`
}

type CodeList struct {
	ID            string      `yaml:"id"`
	Label         string      `yaml:"label"`
	Scheme        string      `yaml:"scheme"`
	ListVintage   string      `yaml:"list_vintage"`
	ListSemantics string      `yaml:"list_semantics"`
	IsComplete    bool        `yaml:"is_complete"`
	Edition       string      `yaml:"edition"`
	SourceURL     string      `yaml:"source_url"`
	Citation      string      `yaml:"citation"`
	RetrievedAt   string      `yaml:"retrieved_at"`
	Codes         []ListEntry `yaml:"codes"`

	Path string `yaml:"-"`
}

const (
	SemanticsClosed = "enumerative_closed"
	SemanticsOpen   = "illustrative_open"
)

// ---------------------------------------------------------------- attribute registry

type Attribute struct {
	URI              string `yaml:"uri"`
	Label            string `yaml:"label"`
	ScopeUnit        string `yaml:"scope_unit"`
	Datatype         string `yaml:"datatype"`
	CollectionMethod string `yaml:"collection_method"`
	DataClass        string `yaml:"data_class"`
	LLMEgressAllowed bool   `yaml:"llm_egress_allowed"`
	Note             string `yaml:"note"`
}

type Registry struct {
	Version    int         `yaml:"version"`
	Attributes []Attribute `yaml:"attributes"`
}
