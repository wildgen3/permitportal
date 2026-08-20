package spec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Corpus is the loaded specification: every rule, code list, credential set, and the
// attribute registry, indexed for lookup.
type Corpus struct {
	Root        string
	Rules       map[string]*Rule
	Lists       map[string]*CodeList
	Attributes  map[string]Attribute
	Credentials map[string]Credential

	// RuleIDs and friends are sorted, so anything that iterates the corpus produces
	// stable output. Go map order is randomised and determinations must be
	// byte-reproducible.
	RuleIDs       []string
	CredentialIDs []string
}

// Load reads the specification rooted at specDir (the repository's spec/ directory).
//
// Fails closed on every ambiguity. A corpus that loaded "most of" the rules would
// silently change determinations, so a single unreadable or duplicate-id file is an
// error, not a warning.
func Load(specDir string) (*Corpus, error) {
	abs, err := filepath.Abs(specDir)
	if err != nil {
		return nil, err
	}
	c := &Corpus{
		Root:        abs,
		Rules:       map[string]*Rule{},
		Lists:       map[string]*CodeList{},
		Attributes:  map[string]Attribute{},
		Credentials: map[string]Credential{},
	}

	if err := c.loadRegistry(filepath.Join(abs, "registry", "attributes.yaml")); err != nil {
		return nil, err
	}
	if err := c.loadLists(filepath.Join(abs, "lists")); err != nil {
		return nil, err
	}
	if err := c.loadRules(filepath.Join(abs, "rules")); err != nil {
		return nil, err
	}
	if err := c.loadCredentials(filepath.Join(abs, "credentials")); err != nil {
		return nil, err
	}

	if len(c.Rules) == 0 {
		return nil, fmt.Errorf("spec: loaded 0 rules from %s — refusing to return an "+
			"empty corpus, because every determination made against it would be vacuous", abs)
	}

	c.RuleIDs = sortedKeys(c.Rules)
	c.CredentialIDs = make([]string, 0, len(c.Credentials))
	for id := range c.Credentials {
		c.CredentialIDs = append(c.CredentialIDs, id)
	}
	sort.Strings(c.CredentialIDs)
	return c, nil
}

func sortedKeys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func decodeFile(path string, into any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(false)
	if err := dec.Decode(into); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func (c *Corpus) loadRegistry(path string) error {
	var reg Registry
	if err := decodeFile(path, &reg); err != nil {
		return err
	}
	for _, a := range reg.Attributes {
		if a.URI == "" {
			return fmt.Errorf("%s: attribute with no uri", path)
		}
		if _, dup := c.Attributes[a.URI]; dup {
			return fmt.Errorf("%s: duplicate attribute uri %q", path, a.URI)
		}
		c.Attributes[a.URI] = a
	}
	if len(c.Attributes) == 0 {
		return fmt.Errorf("%s: registry is empty", path)
	}
	return nil
}

func (c *Corpus) loadLists(dir string) error {
	return walkYAML(dir, func(path string) error {
		var l CodeList
		if err := decodeFile(path, &l); err != nil {
			return err
		}
		if l.ID == "" {
			return fmt.Errorf("%s: code list has no id", path)
		}
		if _, dup := c.Lists[l.ID]; dup {
			return fmt.Errorf("%s: duplicate list id %q", path, l.ID)
		}
		switch l.ListSemantics {
		case SemanticsClosed, SemanticsOpen:
		default:
			return fmt.Errorf("%s: list %q declares unknown list_semantics %q",
				path, l.ID, l.ListSemantics)
		}
		l.Path = path
		c.Lists[l.ID] = &l
		return nil
	})
}

func (c *Corpus) loadRules(dir string) error {
	return walkYAML(dir, func(path string) error {
		var r Rule
		if err := decodeFile(path, &r); err != nil {
			return err
		}
		if r.ID == "" {
			return fmt.Errorf("%s: rule has no id", path)
		}
		if _, dup := c.Rules[r.ID]; dup {
			return fmt.Errorf("%s: duplicate rule id %q", path, r.ID)
		}
		r.Path = path
		c.Rules[r.ID] = &r
		return nil
	})
}

func (c *Corpus) loadCredentials(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return walkYAML(dir, func(path string) error {
		var set CredentialSet
		if err := decodeFile(path, &set); err != nil {
			return err
		}
		for _, cred := range set.Credentials {
			if cred.ID == "" {
				return fmt.Errorf("%s: credential with no id", path)
			}
			if _, dup := c.Credentials[cred.ID]; dup {
				return fmt.Errorf("%s: duplicate credential id %q", path, cred.ID)
			}
			c.Credentials[cred.ID] = cred
		}
		return nil
	})
}

func walkYAML(dir string, fn func(path string) error) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("spec: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("spec: %s is not a directory", dir)
	}
	var paths []string
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext == ".yaml" || ext == ".yml" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Sorted, so load order — and therefore any order-dependent diagnostic — does not
	// depend on the filesystem.
	sort.Strings(paths)
	for _, p := range paths {
		if err := fn(p); err != nil {
			return err
		}
	}
	return nil
}
