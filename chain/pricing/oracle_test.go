package pricing

import "testing"

const M = MicroUSD

func TestMedian(t *testing.T) {
	if got := Median([]int64{3, 1, 2}); got != 2 {
		t.Fatalf("odd median = %d, want 2", got)
	}
	if got := Median([]int64{1, 2, 3, 4}); got != 2 { // (2+3)/2 = 2 (floor)
		t.Fatalf("even median = %d, want 2", got)
	}
}

func TestDeviationGuardDropsOutlier(t *testing.T) {
	sources := []int64{100 * M, 101 * M, 99 * M, 200 * M}
	kept := DeviationGuard(sources, 500) // 5% threshold
	if len(kept) != 3 {
		t.Fatalf("kept %d sources, want 3 (200 should be dropped): %v", len(kept), kept)
	}
	for _, s := range kept {
		if s == 200*M {
			t.Fatal("outlier 200 was not dropped")
		}
	}
}

func TestOracleSeedsThenEMA(t *testing.T) {
	o := NewOracle(8, 10, 500, 5000) // alpha=0.8
	if mid, err := o.Update([]int64{100 * M}); err != nil || mid != 100*M {
		t.Fatalf("seed mid=%d err=%v, want 100M nil", mid, err)
	}
	// move to 110: EMA mid += 0.8*(110-100) = 108
	mid, err := o.Update([]int64{110 * M})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if mid != 108*M {
		t.Fatalf("EMA mid=%d, want 108M", mid)
	}
}

func TestOracleBreakerHalts(t *testing.T) {
	o := NewOracle(8, 10, 5000, 500) // breaker at 5%
	o.Update([]int64{100 * M})
	mid, err := o.Update([]int64{110 * M}) // +10% > 5% breaker
	if err != ErrOracleHalted {
		t.Fatalf("err=%v, want ErrOracleHalted", err)
	}
	if !o.Halted {
		t.Fatal("oracle should be halted")
	}
	if mid != 100*M {
		t.Fatalf("mid should stay 100M, got %d", mid)
	}
	o.Resume()
	if o.Halted {
		t.Fatal("Resume should clear halt")
	}
}

func TestOracleRejectsEmpty(t *testing.T) {
	o := NewOracle(8, 10, 500, 5000)
	if _, err := o.Update(nil); err != ErrNoSources {
		t.Fatalf("err=%v, want ErrNoSources", err)
	}
}
