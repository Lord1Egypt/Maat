package pricing

import "testing"

func TestEffectiveSpreadWidensWithVol(t *testing.T) {
	sp := SpreadParams{BaseBps: 15, VolMultBps: 15, MaxBps: 300}
	if got := sp.EffectiveSpreadBps(0); got != 15 {
		t.Fatalf("calm spread = %d, want 15", got)
	}
	// vol signal of 1x (10000 bps) adds VolMultBps -> 30
	if got := sp.EffectiveSpreadBps(BpsDenom); got != 30 {
		t.Fatalf("vol spread = %d, want 30", got)
	}
	// huge vol is capped
	if got := sp.EffectiveSpreadBps(1000 * BpsDenom); got != 300 {
		t.Fatalf("capped spread = %d, want 300", got)
	}
}

func TestMakeQuoteIsTwoSidedAroundMid(t *testing.T) {
	q := MakeQuote(100*M, 15, 0)
	if !(q.Buy < q.Mid && q.Mid < q.Sell) {
		t.Fatalf("quote not two-sided: %+v", q)
	}
	if q.Buy != 99_850_000 || q.Sell != 100_150_000 {
		t.Fatalf("buy=%d sell=%d, want 99.85M / 100.15M", q.Buy, q.Sell)
	}
}

func TestBuyAccruesSpreadAndKeepsBacking(t *testing.T) {
	r := &Reserve{AssetUnits: 100 * M, WrappedSupply: 80 * M, MinBackingBps: BpsDenom}
	q := MakeQuote(100*M, 15, 0)
	if err := r.BuyWrapped(10*M, q); err != nil {
		t.Fatalf("buy err: %v", err)
	}
	if r.WrappedSupply != 90*M {
		t.Fatalf("supply=%d, want 90M", r.WrappedSupply)
	}
	// 10 units * $0.15 spread = $1.50 = 1_500_000 micro-USD
	if r.SpreadUSD != 1_500_000 {
		t.Fatalf("spread=%d, want 1_500_000", r.SpreadUSD)
	}
	if !r.Healthy() {
		t.Fatal("should still be healthy")
	}
}

func TestBuyFailsWhenUndercollateralized(t *testing.T) {
	r := &Reserve{AssetUnits: 100 * M, WrappedSupply: 100 * M, MinBackingBps: BpsDenom}
	q := MakeQuote(100*M, 15, 0)
	if err := r.BuyWrapped(1*M, q); err != ErrInsufficientBacking {
		t.Fatalf("err=%v, want ErrInsufficientBacking", err)
	}
	if r.WrappedSupply != 100*M {
		t.Fatal("supply must be unchanged on failed buy")
	}
}

func TestSellWrappedAccruesSpread(t *testing.T) {
	r := &Reserve{AssetUnits: 100 * M, WrappedSupply: 80 * M, MinBackingBps: BpsDenom}
	q := MakeQuote(100*M, 15, 0)
	if err := r.SellWrapped(10*M, q); err != nil {
		t.Fatalf("sell err: %v", err)
	}
	if r.WrappedSupply != 70*M {
		t.Fatalf("supply=%d, want 70M", r.WrappedSupply)
	}
	if r.SpreadUSD != 1_500_000 { // (mid-buy)=0.15 * 10 units = $1.50
		t.Fatalf("spread=%d, want 1_500_000", r.SpreadUSD)
	}
}

func TestBridgeInOutKeepsOneToOne(t *testing.T) {
	r := &Reserve{AssetUnits: 0, WrappedSupply: 0, MinBackingBps: BpsDenom}
	if err := r.BridgeIn(50 * M); err != nil {
		t.Fatalf("bridge in: %v", err)
	}
	if r.AssetUnits != 50*M || r.WrappedSupply != 50*M {
		t.Fatalf("after bridge-in asset=%d supply=%d", r.AssetUnits, r.WrappedSupply)
	}
	if err := r.BridgeOut(20 * M); err != nil {
		t.Fatalf("bridge out: %v", err)
	}
	if r.AssetUnits != 30*M || r.WrappedSupply != 30*M {
		t.Fatalf("after bridge-out asset=%d supply=%d", r.AssetUnits, r.WrappedSupply)
	}
	if err := r.BridgeOut(999 * M); err != ErrBadAmount {
		t.Fatalf("over-bridge err=%v, want ErrBadAmount", err)
	}
}

// A full round-trip of mixed flow should leave backing >= 100% and positive
// spread revenue — the core thesis, asserted in code.
func TestRoundTripGrowsSpreadAndStaysBacked(t *testing.T) {
	r := &Reserve{MinBackingBps: BpsDenom}
	r.BridgeIn(1000 * M) // 1000 units in, fully backed
	q := MakeQuote(2000*M, 15, 0)
	for i := 0; i < 100; i++ {
		if err := r.SellWrapped(1*M, q); err != nil {
			t.Fatalf("sell %d: %v", i, err)
		}
		if err := r.BuyWrapped(1*M, q); err != nil {
			t.Fatalf("buy %d: %v", i, err)
		}
	}
	if !r.Healthy() {
		t.Fatalf("backing dropped: %d bps", r.BackingBps())
	}
	if r.SpreadUSD <= 0 {
		t.Fatalf("expected positive spread revenue, got %d", r.SpreadUSD)
	}
}
