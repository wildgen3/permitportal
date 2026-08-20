package engine

import "fmt"

// Crosswalk translates a code between vintages of the same scheme.
//
// This exists because ADR-0002 makes the primary key `(scheme, vintage, code)`: a
// rule keyed to the 2007 edition of a list cannot be tested against a subject whose
// code of record is 2022 without an explicit, auditable translation step. Seven codes
// have been reused across NAICS revisions for different activities, so "the code
// matches" is not the same fact as "the activity matches".
type Crosswalk interface {
	// Name appears in the evidence tree. A determination must say which translation
	// authority produced it.
	Name() string

	// Translate maps code from fromVintage into toVintage. composable is false when
	// no mapping data supports the hop; callers must treat that as UNKNOWN and must
	// never fall back to comparing the untranslated code.
	Translate(scheme, code, fromVintage, toVintage string) (codes []string, composable bool)
}

// StrictCrosswalk is the production implementation: identity within a vintage, and
// nothing at all across vintages.
//
// data/crosswalks/ is empty. The honest behaviour for a cross-vintage hop is
// therefore "not composable", which propagates as UNKNOWN. The alternative —
// comparing a 2022 code against a 2007 list and calling a miss a FALSE — is the exact
// defect this repository exists to argue against, and it would be invisible in the
// output.
type StrictCrosswalk struct{}

func (StrictCrosswalk) Name() string { return "strict" }

func (StrictCrosswalk) Translate(_, code, from, to string) ([]string, bool) {
	if from == to {
		return []string{code}, true
	}
	return nil, false
}

// FixtureVintageStableCrosswalk asserts that codes are stable across vintages.
//
// That assertion is FALSE in general. It exists so specification fixtures can
// exercise the rest of the engine without shipping crosswalk data this repository has
// not actually retrieved, and it is usable only from Fixture mode. It names itself in
// every evidence node it touches so a fixture determination can never be mistaken for
// a production one.
type FixtureVintageStableCrosswalk struct{}

func (FixtureVintageStableCrosswalk) Name() string { return "fixture_vintage_stable" }

func (FixtureVintageStableCrosswalk) Translate(_, code, _, _ string) ([]string, bool) {
	return []string{code}, true
}

func crosswalkFor(m Mode) (Crosswalk, error) {
	switch m {
	case Production:
		return StrictCrosswalk{}, nil
	case Fixture:
		return FixtureVintageStableCrosswalk{}, nil
	default:
		return nil, fmt.Errorf("engine: unknown mode %q", m)
	}
}
