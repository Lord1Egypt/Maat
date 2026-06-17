package pricing

import "testing"

// FuzzReserveInvariants: drive a reserve with an arbitrary sequence of buy/sell
// operations and assert the safety invariants always hold:
//   - WrappedSupply never goes negative
//   - whenever a BuyWrapped succeeds, backing remains >= the minimum
//   - cumulative spread revenue is monotonic non-decreasing (never leaks out)
func FuzzReserveInvariants(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{0, 0, 0, 255, 128})

	f.Fuzz(func(t *testing.T, ops []byte) {
		r := &Reserve{MinBackingBps: BpsDenom}
		r.BridgeIn(1000 * MicroUSD) // start fully backed
		q := MakeQuote(2000*MicroUSD, 15, 0)

		lastSpread := r.SpreadUSD
		for _, b := range ops {
			size := (int64(b%16) + 1) * MicroUSD // 1..16 units
			if b%2 == 0 {
				_ = r.SellWrapped(size, q)
			} else {
				if err := r.BuyWrapped(size, q); err == nil && !r.Healthy() {
					t.Fatalf("buy succeeded but backing broke: %d bps", r.BackingBps())
				}
			}
			if r.WrappedSupply < 0 {
				t.Fatalf("negative supply: %d", r.WrappedSupply)
			}
			if r.SpreadUSD < lastSpread {
				t.Fatalf("spread revenue decreased: %d -> %d", lastSpread, r.SpreadUSD)
			}
			lastSpread = r.SpreadUSD
		}
	})
}
