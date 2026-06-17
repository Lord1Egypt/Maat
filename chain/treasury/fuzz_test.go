package treasury

import "testing"

// FuzzDistributeNoLeak: for any non-negative amount within a realistic domain
// and any valid split, the four allocations must be non-negative and sum back
// to exactly the input — no value is ever created or destroyed.
//
// Domain note: amounts are bounded to <= $1e9 (1e15 micro-USD). Production uses
// the SDK's big.Int, so this fixed-point int64 core is only exercised in its
// real operating range; absurd values that would overflow int64 are out of scope.
func FuzzDistributeNoLeak(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(1001))
	f.Add(int64(999_999_999))

	splits := []Split{DefaultSpreadSplit, DefaultBridgeOutSplit, {ReserveBps: 10000}}

	f.Fuzz(func(t *testing.T, amount int64) {
		if amount < 0 || amount > 1_000_000_000_000_000 {
			return // out of the supported domain
		}
		for _, s := range splits {
			a, err := Distribute(amount, s)
			if err != nil {
				t.Fatalf("unexpected err for amount=%d split=%+v: %v", amount, s, err)
			}
			if a.Reserve < 0 || a.Rewards < 0 || a.Insurance < 0 || a.Treasury < 0 {
				t.Fatalf("negative allocation: %+v", a)
			}
			if a.Total() != amount {
				t.Fatalf("value leak: total=%d, want %d (alloc=%+v)", a.Total(), amount, a)
			}
		}
	})
}
