package resolver_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wildgen3/permitportal/services/resolver/internal/engine"
	"github.com/wildgen3/permitportal/services/resolver/internal/spec"
)

// The golden files in packages/rules/tests/golden are the correctness specification:
// a profile in, an expected obligation set out. They are authored as data so that
// adding a case never means adding a branch to the evaluator.

type goldenCase struct {
	ID                 string                    `yaml:"id"`
	Mode               string                    `yaml:"mode"`
	Expect             string                    `yaml:"expect"`
	ExpectError        string                    `yaml:"expect_error"`
	ExpectReason       string                    `yaml:"expect_reason"`
	ExpectMissing      []string                  `yaml:"expect_missing"`
	AlsoExpectSurfaced []string                  `yaml:"also_expect_surfaced"`
	Profile            map[string]map[string]any `yaml:"profile"`
}

type goldenFile struct {
	Rule    string       `yaml:"rule"`
	AsOfLaw string       `yaml:"as_of_law"`
	Cases   []goldenCase `yaml:"cases"`
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "spec", "registry", "attributes.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate the repository root (no spec/registry/attributes.yaml above the test)")
		}
		dir = parent
	}
}

func loadCorpus(t *testing.T) (*spec.Corpus, string) {
	t.Helper()
	root := repoRoot(t)
	c, err := spec.Load(filepath.Join(root, "spec"))
	if err != nil {
		t.Fatalf("loading spec: %v", err)
	}
	if err := engine.ValidateOverlays(c); err != nil {
		t.Fatalf("overlay validation: %v", err)
	}
	return c, root
}

func TestGolden(t *testing.T) {
	corpus, root := loadCorpus(t)

	pattern := filepath.Join(root, "packages", "rules", "tests", "golden", "*.golden.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	// Fail closed. A test run that discovered no cases and reported PASS is
	// indistinguishable from one that verified something, and this repository has
	// already been bitten by exactly that.
	if len(files) == 0 {
		t.Fatalf("no golden files matched %s — refusing to report a pass over an empty set", pattern)
	}
	sort.Strings(files)

	total := 0
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var gf goldenFile
		if err := yaml.Unmarshal(data, &gf); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if len(gf.Cases) == 0 {
			t.Fatalf("%s declares no cases", path)
		}

		for _, tc := range gf.Cases {
			total++
			name := filepath.Base(path) + "/" + tc.ID
			t.Run(name, func(t *testing.T) {
				mode := engine.Fixture
				if tc.Mode != "" {
					mode = engine.Mode(tc.Mode)
				}
				eng, err := engine.New(corpus, mode)
				if err != nil {
					t.Fatal(err)
				}
				d, err := eng.Evaluate(gf.Rule, engine.Profile{Facts: tc.Profile}, gf.AsOfLaw)

				if tc.ExpectError != "" {
					// A case that asserts a refusal. Producing a determination here is
					// the failure: the engine would have answered from data the
					// specification says cannot support an answer.
					if err == nil {
						t.Fatalf("expected a refusal containing %q, got determination %s",
							tc.ExpectError, d.Outcome)
					}
					if !strings.Contains(err.Error(), tc.ExpectError) {
						t.Fatalf("refusal did not mention %q:\n%v", tc.ExpectError, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("evaluate: %v", err)
				}

				if got := string(d.Outcome); got != tc.Expect {
					t.Errorf("outcome = %s, want %s\n  reasons: %v\n  missing: %v",
						got, tc.Expect, d.RouteToAuthority, d.MissingAttributes)
				}
				if tc.ExpectReason != "" && d.Reason != tc.ExpectReason {
					t.Errorf("reason = %q, want %q", d.Reason, tc.ExpectReason)
				}
				if len(tc.ExpectMissing) > 0 {
					want := append([]string(nil), tc.ExpectMissing...)
					sort.Strings(want)
					got := append([]string(nil), d.MissingAttributes...)
					sort.Strings(got)
					if strings.Join(got, ",") != strings.Join(want, ",") {
						t.Errorf("missing_attributes = %v, want %v", got, want)
					}
				}
				for _, want := range tc.AlsoExpectSurfaced {
					found := false
					for _, s := range d.Surfaced {
						if s.ID == want {
							found = true
						}
					}
					if !found {
						t.Errorf("obligation %q was not surfaced; an answer that omits it is a harmful answer", want)
					}
				}
			})
		}
	}
	t.Logf("ran %d golden case(s) across %d file(s)", total, len(files))
}
