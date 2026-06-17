package pricing

import "testing"

// FuzzQuoteIsOrdered: for any realistic mid and non-negative spread (no skew),
// the quote must satisfy 0 < buy <= mid <= sell. This guarantees the protocol
// never quotes a negative or inverted price.
//
// Domain: mid in (0, $10M], spread in [0, 5000] bps. Larger values are out of
// the supported fixed-point range (production uses the SDK big.Int).
func FuzzQuoteIsOrdered(f *testing.F) {
	f.Add(int64(100*MicroUSD), int64(15))
	f.Add(int64(1), int64(0))
	f.Add(int64(2000*MicroUSD), int64(300))

	f.Fuzz(func(t *testing.T, mid, spreadBps int64) {
		if mid <= 0 || mid > 10_000_000*MicroUSD {
			return
		}
		if spreadBps < 0 || spreadBps > 5000 {
			return
		}
		q := MakeQuote(mid, spreadBps, 0)
		if q.Buy <= 0 {
			t.Fatalf("non-positive buy: %+v", q)
		}
		if !(q.Buy <= q.Mid && q.Mid <= q.Sell) {
			t.Fatalf("quote not ordered buy<=mid<=sell: %+v", q)
		}
	})
}

// FuzzMedianWithinBounds: the median is always within [min, max] of its inputs.
func FuzzMedianWithinBounds(f *testing.F) {
	f.Add(int64(1), int64(2), int64(3))
	f.Add(int64(-5), int64(0), int64(100))

	f.Fuzz(func(t *testing.T, a, b, c int64) {
		in := []int64{a, b, c}
		m := Median(in)
		lo, hi := a, a
		for _, v := range in {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}
		if m < lo || m > hi {
			t.Fatalf("median %d outside [%d,%d] for %v", m, lo, hi, in)
		}
	})
}
