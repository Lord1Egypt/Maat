package governance

import "testing"

func TestPasses(t *testing.T) {
	// 40% turnout (quorum 30% met), 70% yes (>60% approval)
	tally := Tally{Yes: 280, No: 120, Abstain: 0} // turnout 400 of 1000 = 40%
	if o := Decide(tally, 1000, DefaultParams(), false); o != OutcomePassed {
		t.Fatalf("outcome=%v, want passed", o)
	}
}

func TestFailsNoQuorum(t *testing.T) {
	// only 20% turnout < 30% quorum
	tally := Tally{Yes: 200, No: 0}
	if o := Decide(tally, 1000, DefaultParams(), false); o != OutcomeRejectedNoQuorum {
		t.Fatalf("outcome=%v, want no-quorum", o)
	}
}

func TestFailsApproval(t *testing.T) {
	// quorum met (50%) but only 50% yes < 60% approval
	tally := Tally{Yes: 250, No: 250}
	if o := Decide(tally, 1000, DefaultParams(), false); o != OutcomeRejected {
		t.Fatalf("outcome=%v, want rejected", o)
	}
}

func TestAbstainCountsForQuorumNotApproval(t *testing.T) {
	// 70 yes / 30 no = 70% approval; 300 abstain lifts turnout to 40% (quorum)
	tally := Tally{Yes: 70, No: 30, Abstain: 300}
	if o := Decide(tally, 1000, DefaultParams(), false); o != OutcomePassed {
		t.Fatalf("outcome=%v, want passed", o)
	}
}

func TestEmergencyNeeds80(t *testing.T) {
	// 70% yes passes normally but fails the 80% emergency bar
	tally := Tally{Yes: 280, No: 120} // 40% turnout, 70% yes
	if o := Decide(tally, 1000, DefaultParams(), false); o != OutcomePassed {
		t.Fatalf("normal outcome=%v, want passed", o)
	}
	if o := Decide(tally, 1000, DefaultParams(), true); o != OutcomeRejected {
		t.Fatalf("emergency outcome=%v, want rejected (<80%%)", o)
	}
	// 85% yes clears emergency
	strong := Tally{Yes: 340, No: 60} // 40% turnout, 85% yes
	if o := Decide(strong, 1000, DefaultParams(), true); o != OutcomePassed {
		t.Fatalf("strong emergency outcome=%v, want passed", o)
	}
}

func TestZeroStakedRejected(t *testing.T) {
	if o := Decide(Tally{Yes: 1}, 0, DefaultParams(), false); o != OutcomeRejectedNoQuorum {
		t.Fatalf("outcome=%v, want no-quorum", o)
	}
}
