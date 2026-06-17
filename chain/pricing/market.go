package pricing

import "errors"

var (
	ErrInsufficientBacking = errors.New("market: trade would break 100% backing")
	ErrMarketHalted        = errors.New("market: halted (backing below minimum)")
	ErrBadAmount           = errors.New("market: amount must be positive")
)

// SpreadParams controls how the per-block quote is built around the oracle mid.
// All values are basis points.
type SpreadParams struct {
	BaseBps    int64 // per-side base spread (e.g. 15 = 0.15%)
	VolMultBps int64 // extra spread added per 1x of realized-vol signal
	SkewMaxBps int64 // max inventory-skew tilt
	MaxBps     int64 // hard cap on effective per-side spread
}

// clamp keeps v within [lo, hi].
func clamp(v, lo, hi int64) int64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// EffectiveSpreadBps widens the base spread with volatility and tilts it with
// inventory skew. volSignalBps is a realized-vol proxy (e.g. |last move| in bps);
// invSkewBps is signed: positive when the protocol is long the asset (so it
// shades the SELL side cheaper to offload — handled by the caller via sign).
func (sp SpreadParams) EffectiveSpreadBps(volSignalBps int64) int64 {
	eff := sp.BaseBps + sp.VolMultBps*volSignalBps/BpsDenom
	cap := sp.MaxBps
	if cap == 0 {
		cap = sp.BaseBps * 20 // sane default cap if unset
	}
	return clamp(eff, sp.BaseBps, cap)
}

// Quote is the single clearing price for the current block: fixed within the
// block (no slippage, no MEV) but centered on the oracle mid (solvent).
type Quote struct {
	Mid  int64
	Buy  int64 // price at which the protocol BUYS the asset from a user
	Sell int64 // price at which the protocol SELLS the asset to a user
}

// MakeQuote builds the two-sided quote around mid given an effective spread and
// an inventory-skew tilt (skewBps, signed). Skew shifts the whole quote to pull
// inventory back toward target.
func MakeQuote(mid, effSpreadBps, skewBps int64) Quote {
	center := mid + mid*skewBps/BpsDenom
	return Quote{
		Mid:  mid,
		Buy:  center - center*effSpreadBps/BpsDenom,
		Sell: center + center*effSpreadBps/BpsDenom,
	}
}

// Reserve is the asset book for one wrapped asset. Amounts are in micro-units of
// the asset; value bookkeeping is in micro-USD.
type Reserve struct {
	AssetUnits    int64 // real asset held in custody (micro-units)
	WrappedSupply int64 // wrapped tokens in circulation (micro-units), a claim on AssetUnits
	SpreadUSD     int64 // cumulative spread revenue captured (micro-USD)
	MinBackingBps int64 // halt trading below this (e.g. 10000 = 100%)
}

// BackingBps = AssetUnits / WrappedSupply in basis points. Full (>=10000) means
// every wrapped token is backed 1:1 by a real asset unit.
func (r *Reserve) BackingBps() int64 {
	if r.WrappedSupply == 0 {
		return BpsDenom * 1000 // effectively infinite when nothing is circulating
	}
	return r.AssetUnits * BpsDenom / r.WrappedSupply
}

// Healthy reports whether backing is at/above the configured minimum.
func (r *Reserve) Healthy() bool { return r.BackingBps() >= r.MinBackingBps }

// BuyWrapped: a user buys `assetUnits` of wrapped asset at the quote's Sell
// price. The protocol mints wrapped tokens (a new claim) and books the spread
// vs mid as revenue. Fails if it would break the minimum backing.
func (r *Reserve) BuyWrapped(assetUnits int64, q Quote) error {
	if assetUnits <= 0 {
		return ErrBadAmount
	}
	if !r.Healthy() {
		return ErrMarketHalted
	}
	newSupply := r.WrappedSupply + assetUnits
	if r.AssetUnits*BpsDenom < newSupply*r.MinBackingBps {
		return ErrInsufficientBacking
	}
	r.WrappedSupply = newSupply
	// spread captured = (sell - mid) * units, in micro-USD
	r.SpreadUSD += (q.Sell - q.Mid) * assetUnits / MicroUSD
	return nil
}

// SellWrapped: a user sells `assetUnits` of wrapped asset to the protocol at the
// quote's Buy price. Wrapped tokens are burned; the protocol books the spread
// vs mid (it acquired the claim below mid).
func (r *Reserve) SellWrapped(assetUnits int64, q Quote) error {
	if assetUnits <= 0 {
		return ErrBadAmount
	}
	if assetUnits > r.WrappedSupply {
		return ErrBadAmount
	}
	r.WrappedSupply -= assetUnits
	r.SpreadUSD += (q.Mid - q.Buy) * assetUnits / MicroUSD
	return nil
}

// BridgeIn adds real custody (deposit) and mints matching wrapped supply 1:1.
func (r *Reserve) BridgeIn(assetUnits int64) error {
	if assetUnits <= 0 {
		return ErrBadAmount
	}
	r.AssetUnits += assetUnits
	r.WrappedSupply += assetUnits
	return nil
}

// BridgeOut burns wrapped supply and releases real custody 1:1, only if backing
// stays healthy afterward.
func (r *Reserve) BridgeOut(assetUnits int64) error {
	if assetUnits <= 0 || assetUnits > r.WrappedSupply || assetUnits > r.AssetUnits {
		return ErrBadAmount
	}
	r.AssetUnits -= assetUnits
	r.WrappedSupply -= assetUnits
	if !r.Healthy() {
		// revert
		r.AssetUnits += assetUnits
		r.WrappedSupply += assetUnits
		return ErrInsufficientBacking
	}
	return nil
}
