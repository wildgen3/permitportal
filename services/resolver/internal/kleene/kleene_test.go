package kleene

import "testing"

var all3 = []Truth{True, False, Unknown}

func v(t Truth) Value {
	if t == Unknown {
		return Unk(Reason{Kind: MissingAttribute, Ref: "x"})
	}
	return Value{T: t}
}

// The truth tables, written out rather than derived, because the whole design rests
// on them and a derivation would share any mistake with the implementation.
func TestTruthTables(t *testing.T) {
	and := map[[2]Truth]Truth{
		{True, True}: True, {True, False}: False, {True, Unknown}: Unknown,
		{False, True}: False, {False, False}: False, {False, Unknown}: False,
		{Unknown, True}: Unknown, {Unknown, False}: False, {Unknown, Unknown}: Unknown,
	}
	or := map[[2]Truth]Truth{
		{True, True}: True, {True, False}: True, {True, Unknown}: True,
		{False, True}: True, {False, False}: False, {False, Unknown}: Unknown,
		{Unknown, True}: True, {Unknown, False}: Unknown, {Unknown, Unknown}: Unknown,
	}
	not := map[Truth]Truth{True: False, False: True, Unknown: Unknown}

	for _, a := range all3 {
		for _, b := range all3 {
			if got := All([]Value{v(a), v(b)}).T; got != and[[2]Truth{a, b}] {
				t.Errorf("All(%v, %v) = %v, want %v", a, b, got, and[[2]Truth{a, b}])
			}
			if got := Any([]Value{v(a), v(b)}).T; got != or[[2]Truth{a, b}] {
				t.Errorf("Any(%v, %v) = %v, want %v", a, b, got, or[[2]Truth{a, b}])
			}
		}
		if got := Not(v(a)).T; got != not[a] {
			t.Errorf("Not(%v) = %v, want %v", a, got, not[a])
		}
	}
}

// The one property the product depends on: nothing turns an absent fact into FALSE.
// Rego's negation-as-failure fails exactly this test, which is why it was rejected.
func TestUnknownNeverBecomesFalseWithoutAnEvaluatedFalse(t *testing.T) {
	u := v(Unknown)
	if got := Not(u).T; got == False {
		t.Fatal("Not(UNKNOWN) produced FALSE — negation-as-failure")
	}
	if got := All([]Value{v(True), u}).T; got == False {
		t.Fatal("All(TRUE, UNKNOWN) produced FALSE")
	}
	if got := Any([]Value{u}).T; got == False {
		t.Fatal("Any(UNKNOWN) produced FALSE")
	}
	if got := None([]Value{u}).T; got == False {
		t.Fatal("None(UNKNOWN) produced FALSE")
	}
	// And the converse: a real FALSE must still be reachable, or the test above is
	// satisfied by an evaluator that never says FALSE at all.
	if got := All([]Value{v(True), v(False)}).T; got != False {
		t.Fatalf("All(TRUE, FALSE) = %v, want FALSE", got)
	}
}

func TestNotPreservesReasons(t *testing.T) {
	in := Unk(Reason{Kind: MissingAttribute, Ref: "establishment.exemption_claims"})
	out := Not(in)
	if len(out.Reasons) != 1 || out.Reasons[0].Ref != "establishment.exemption_claims" {
		t.Fatalf("Not() dropped the reason: %+v", out)
	}
}

func TestReasonsAreDeterministic(t *testing.T) {
	a := Reason{Kind: ListIncomplete, Ref: "lists/b"}
	b := Reason{Kind: MissingAttribute, Ref: "a"}
	one := Any([]Value{Unk(a, b), Unk(b, a), Unk(a)})
	two := Any([]Value{Unk(b), Unk(a, b), Unk(b, a)})
	if len(one.Reasons) != 2 || len(two.Reasons) != 2 {
		t.Fatalf("dedupe failed: %+v %+v", one.Reasons, two.Reasons)
	}
	for i := range one.Reasons {
		if one.Reasons[i] != two.Reasons[i] {
			t.Fatalf("reason order depends on input order: %+v vs %+v", one.Reasons, two.Reasons)
		}
	}
}

func TestOnlyMissingAttributeIsAskable(t *testing.T) {
	askable := []ReasonKind{MissingAttribute}
	notAskable := []ReasonKind{
		NotSelfDeterminable, ListIncomplete, ListIllustrativeOpen,
		CrosswalkUnavailable, JurisdictionRuleMissing,
	}
	for _, k := range askable {
		if !k.AskableByApplicant() {
			t.Errorf("%s should be askable", k)
		}
	}
	for _, k := range notAskable {
		if k.AskableByApplicant() {
			t.Errorf("%s must not be presented to the applicant as a question", k)
		}
	}
}
