package spec

// Credential requirement trees. CTDL supplies the vocabulary but ships no SHACL and
// no ordering semantics (ADR-0012), so the tree shape, the topological order, and
// cycle detection are all implemented here.

type NodeType string

const (
	AndGroup NodeType = "AND_GROUP"
	OrGroup  NodeType = "OR_GROUP"
	Leaf     NodeType = "LEAF"
)

type Requirement struct {
	ID               string        `yaml:"id"`
	NodeType         NodeType      `yaml:"node_type"`
	LegalSource      string        `yaml:"legal_source"`
	EdgeKind         string        `yaml:"edge_kind"`
	TargetCredential string        `yaml:"target_credential"`
	TargetPredicate  string        `yaml:"target_predicate"`
	YearsExperience  int           `yaml:"years_experience"`
	Children         []Requirement `yaml:"children"`
}

// Prerequisites returns every credential this requirement tree depends on, in
// document order. OR alternatives are included: an alternative that is not taken is
// still a real edge in the dependency graph, and omitting it would make the ordering
// wrong for anyone who does take it.
func (r Requirement) Prerequisites() []string {
	var out []string
	var walk func(Requirement)
	walk = func(n Requirement) {
		if n.TargetCredential != "" {
			out = append(out, n.TargetCredential)
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(r)
	return out
}

type Credential struct {
	ID               string        `yaml:"id"`
	Type             string        `yaml:"type"`
	Label            string        `yaml:"label"`
	IssuingAuthority string        `yaml:"issuing_authority"`
	SourceURL        string        `yaml:"source_url"`
	Citation         string        `yaml:"citation"`
	Requirements     []Requirement `yaml:"requirements"`
}

func (c Credential) Prerequisites() []string {
	var out []string
	for _, r := range c.Requirements {
		out = append(out, r.Prerequisites()...)
	}
	return out
}

type CredentialSet struct {
	Version      int          `yaml:"version"`
	Jurisdiction string       `yaml:"jurisdiction"`
	Credentials  []Credential `yaml:"credentials"`

	Path string `yaml:"-"`
}
