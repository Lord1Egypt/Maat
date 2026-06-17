package scenario

import "testing"

// The end-to-end integration assertion: across a full multi-block run wiring
// oracle + market + treasury + bridge, backing must hold >=100%, the protocol
// must capture positive spread, the treasury splits must populate, and the
// bridge cap must actually throttle some withdrawals.
func TestEndToEndReserveGrowsAndProtectionsHold(t *testing.T) {
	r := Run(Default())

	if !r.BackingHeld {
		t.Fatalf("backing dropped below 100%% during run")
	}
	if r.FinalBackingBps < 10000 {
		t.Fatalf("final backing=%d bps, want >= 10000", r.FinalBackingBps)
	}
	if r.SpreadCaptured <= 0 {
		t.Fatalf("expected positive spread captured, got %d", r.SpreadCaptured)
	}
	// treasury split must have populated every fund
	if r.ReserveFund <= 0 || r.InsuranceFund <= 0 || r.RewardsFund <= 0 || r.TreasuryFund <= 0 {
		t.Fatalf("treasury funds not all populated: %+v", r)
	}
	// no value leak: the four funds sum to the captured spread
	if r.ReserveFund+r.InsuranceFund+r.RewardsFund+r.TreasuryFund != r.SpreadCaptured {
		t.Fatalf("fund sum != captured (value leak): %+v", r)
	}
	// the bridge cap must bite at least once (proves the throttle works)
	if r.BridgeThrottled == 0 {
		t.Fatalf("expected bridge cap to throttle some withdrawals, got 0")
	}
}

func TestDeterministic(t *testing.T) {
	a, b := Run(Default()), Run(Default())
	if a != b {
		t.Fatalf("scenario not deterministic:\n%+v\n%+v", a, b)
	}
}
