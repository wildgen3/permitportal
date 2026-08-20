// Command resolve evaluates one rule against one profile and prints the
// determination, its evidence tree, and — when the rule produces a credential — the
// order that credential must be obtained in.
//
// It is deliberately a CLI before it is a service. The determination contract is
// stable long before the HTTP contract is (spec/api/openapi.yaml is not written), and
// a CLI is the shortest path to a reviewer being able to run this themselves.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wildgen3/permitportal/services/resolver/internal/credential"
	"github.com/wildgen3/permitportal/services/resolver/internal/engine"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

type output struct {
	Determination *engine.Determination `json:"determination"`
	Chain         *credential.Chain     `json:"credential_chain,omitempty"`
}

func main() {
	var (
		specDir   = flag.String("spec", "", "path to the spec/ directory (default: found by walking up from the working directory)")
		ruleID    = flag.String("rule", "", "rule id to evaluate")
		profile   = flag.String("profile", "", "path to a YAML fact set")
		asOfLaw   = flag.String("as-of-law", "", "the date the law is read as of (YYYY-MM-DD)")
		mode      = flag.String("mode", string(engine.Production), "production | fixture")
		listRules = flag.Bool("list-rules", false, "print the rule ids in the corpus and exit")
	)
	flag.Parse()

	if err := run(*specDir, *ruleID, *profile, *asOfLaw, *mode, *listRules); err != nil {
		fmt.Fprintf(os.Stderr, "resolve: %v\n", err)
		os.Exit(1)
	}
}

func run(specDir, ruleID, profilePath, asOfLaw, mode string, listRules bool) error {
	if specDir == "" {
		found, err := findSpec()
		if err != nil {
			return err
		}
		specDir = found
	}
	corpus, err := spec.Load(specDir)
	if err != nil {
		return err
	}
	if err := engine.ValidateOverlays(corpus); err != nil {
		return err
	}

	if listRules {
		for _, id := range corpus.RuleIDs {
			fmt.Println(id)
		}
		return nil
	}
	if ruleID == "" || profilePath == "" || asOfLaw == "" {
		return fmt.Errorf("--rule, --profile and --as-of-law are all required (see --help)")
	}

	data, err := os.ReadFile(profilePath)
	if err != nil {
		return err
	}
	p, err := engine.UnmarshalProfile(data)
	if err != nil {
		return err
	}

	eng, err := engine.New(corpus, engine.Mode(mode))
	if err != nil {
		return err
	}
	d, err := eng.Evaluate(ruleID, p, asOfLaw)
	if err != nil {
		return err
	}

	out := output{Determination: d}
	if d.Outcome == engine.OutcomeTrue && d.Produces != nil && d.Produces.Credential != "" {
		chain, err := credential.Resolve(corpus, []string{d.Produces.Credential})
		if err != nil {
			return err
		}
		out.Chain = chain
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// findSpec walks up from the working directory looking for the specification. The
// binary is useful from anywhere in the tree, and guessing a relative path would make
// the determination depend on where the operator happened to be standing.
func findSpec() (string, error) {
	dir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "spec")
		if _, err := os.Stat(filepath.Join(candidate, "registry", "attributes.yaml")); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no spec/ directory found above the working directory; pass --spec")
		}
		dir = parent
	}
}
