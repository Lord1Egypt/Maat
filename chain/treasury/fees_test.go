package treasury

import "testing"

func TestDistributeSumsExactly(t *testing.T) {
	a, err := Distribute(1000, DefaultSpreadSplit)
	if err != nil {
		t.Fatalf("distribute: %v", err)
	}
	if a.Reserve != 400 || a.Rewards != 250 || a.Insurance != 200 || a.Treasury != 150 {
		t.Fatalf("alloc=%+v, want 400/250/200/150", a)
	}
	if a.Total() != 1000 {
		t.Fatalf("total=%d, want 1000", a.Total())
	}
}

func TestRemainderGoesToReserveNoLeak(t *testing.T) {
	// 1001 doesn't divide cleanly; reserve must absorb the remainder.
	a, err := Distribute(1001, DefaultSpreadSplit)
	if err != nil {
		t.Fatalf("distribute: %v", err)
	}
	if a.Total() != 1001 {
		t.Fatalf("total=%d, want 1001 (no value leak)", a.Total())
	}
	if a.Reserve != 401 {
		t.Fatalf("reserve=%d, want 401 (absorbs remainder)", a.Reserve)
	}
}

func TestBadSplitRejected(t *testing.T) {
	bad := Split{ReserveBps: 5000, RewardsBps: 2500, InsuranceBps: 2000, TreasuryBps: 1000} // 10500
	if _, err := Distribute(100, bad); err != ErrBadSplit {
		t.Fatalf("err=%v, want ErrBadSplit", err)
	}
}

func TestDefaultSplitsValidate(t *testing.T) {
	if err := DefaultSpreadSplit.Validate(); err != nil {
		t.Fatalf("spread split invalid: %v", err)
	}
	if err := DefaultBridgeOutSplit.Validate(); err != nil {
		t.Fatalf("bridge-out split invalid: %v", err)
	}
}

func TestCollectAccumulates(t *testing.T) {
	var tr Treasury
	tr.Collect(1000, DefaultSpreadSplit)
	tr.Collect(1000, DefaultSpreadSplit)
	if tr.Reserve != 800 || tr.Rewards != 500 || tr.Insurance != 400 || tr.Treasury != 300 {
		t.Fatalf("treasury=%+v after two collects", tr)
	}
}

func TestNegativeAmountRejected(t *testing.T) {
	if _, err := Distribute(-1, DefaultSpreadSplit); err != ErrBadAmount {
		t.Fatalf("err=%v, want ErrBadAmount", err)
	}
}
