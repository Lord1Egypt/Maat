// Package pricing is the deterministic core of Ma'at's x/oracle and x/market
// modules. It uses integer fixed-point math only (no floats) so it is safe to
// run inside CometBFT consensus, where every validator must compute identical
// results.
//
// Units:
//   - prices are micro-USD: 1 USD = 1_000_000 (MicroUSD)
//   - spreads / deviations / ratios are in basis points: 1% = 100 bps (Bps)
package pricing

import (
	"errors"
	"sort"
)

const (
	// MicroUSD is the fixed-point scale for prices (6 decimals).
	MicroUSD int64 = 1_000_000
	// BpsDenom is the basis-points denominator (10000 bps = 100%).
	BpsDenom int64 = 10_000
)

var (
	ErrNoSources    = errors.New("oracle: no price sources")
	ErrAllFiltered  = errors.New("oracle: all sources filtered as outliers")
	ErrOracleHalted = errors.New("oracle: circuit breaker halted")
)

// Median returns the median of the prices (caller must pass len > 0).
// For an even count it averages the two middle values (rounding down).
func Median(prices []int64) int64 {
	c := append([]int64(nil), prices...)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	n := len(c)
	if n%2 == 1 {
		return c[n/2]
	}
	return (c[n/2-1] + c[n/2]) / 2
}

// absDiffBps returns |a-b| / b in basis points.
func absDiffBps(a, b int64) int64 {
	if b == 0 {
		return 0
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d * BpsDenom / b
}

// DeviationGuard drops sources that deviate from the median by more than
// maxDevBps. This resists a single manipulated feed. Returns the kept sources.
func DeviationGuard(sources []int64, maxDevBps int64) []int64 {
	if len(sources) == 0 {
		return nil
	}
	med := Median(sources)
	kept := make([]int64, 0, len(sources))
	for _, s := range sources {
		if absDiffBps(s, med) <= maxDevBps {
			kept = append(kept, s)
		}
	}
	return kept
}

// Oracle is the hardened price feed: median of de-outliered sources, smoothed
// into a TWAP via an integer EMA, with a circuit breaker on large moves.
type Oracle struct {
	Mid       int64 // current time-weighted mid (micro-USD); 0 = uninitialised
	AlphaNum  int64 // EMA smoothing numerator   (freshness)
	AlphaDen  int64 // EMA smoothing denominator; alpha = Num/Den, 0<alpha<=1
	MaxDevBps int64 // outlier threshold for DeviationGuard
	BreakBps  int64 // halt if accepted median moves more than this vs current Mid
	Halted    bool
}

// NewOracle builds an oracle. alphaNum/alphaDen sets freshness (e.g. 8/10 ~ fast).
func NewOracle(alphaNum, alphaDen, maxDevBps, breakBps int64) *Oracle {
	return &Oracle{AlphaNum: alphaNum, AlphaDen: alphaDen, MaxDevBps: maxDevBps, BreakBps: breakBps}
}

// Update ingests one round of source prices and advances the TWAP mid.
// Returns the new mid, or an error (and halts) on bad input / breaker trip.
func (o *Oracle) Update(sources []int64) (int64, error) {
	if o.Halted {
		return o.Mid, ErrOracleHalted
	}
	if len(sources) == 0 {
		return o.Mid, ErrNoSources
	}
	kept := DeviationGuard(sources, o.MaxDevBps)
	if len(kept) == 0 {
		return o.Mid, ErrAllFiltered
	}
	med := Median(kept)

	if o.Mid == 0 { // first observation seeds the mid
		o.Mid = med
		return o.Mid, nil
	}

	// circuit breaker: refuse + halt on an implausible jump
	if o.BreakBps > 0 && absDiffBps(med, o.Mid) > o.BreakBps {
		o.Halted = true
		return o.Mid, ErrOracleHalted
	}

	// integer EMA: mid += alpha*(med - mid)
	o.Mid += (med - o.Mid) * o.AlphaNum / o.AlphaDen
	return o.Mid, nil
}

// Resume clears the halt (governance / security-council action in production).
func (o *Oracle) Resume() { o.Halted = false }
