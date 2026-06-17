// Package scenario wires the deterministic cores (oracle + market + treasury +
// bridge) into a multi-block end-to-end run. It is both an integration test
// fixture and the engine behind cmd/demo. No floats, no RNG — fully repeatable.
package scenario

import (
	"github.com/Lord1Egypt/Maat/chain/bridge"
	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/chain/treasury"
)

const M = pricing.MicroUSD

// Config controls a run. Zero value is not valid; use Default().
type Config struct {
	Blocks       int
	SeedAsset    int64 // initial real asset bridged in (micro-units)
	BasePrice    int64 // starting mid (micro-USD)
	RoundTrips   int   // buy+sell round-trips per block (size TradeSize each)
	TradeSize    int64 // micro-units per leg
	BridgeOutEvery int  // attempt a bridge-out every N blocks
	BridgeOutSize  int64
}

func Default() Config {
	return Config{
		Blocks: 720, SeedAsset: 1000 * M, BasePrice: 2000 * M,
		RoundTrips: 5, TradeSize: 1 * M,
		// heavy outflow: 8 units every 3 blocks = 64/window vs a 50/window cap,
		// so the cap throttles ~2 of every 8 attempts (protection demonstrated).
		BridgeOutEvery: 3, BridgeOutSize: 8 * M,
	}
}

// Result summarises a run.
type Result struct {
	Blocks          int
	FinalBackingBps int64
	BackingHeld     bool  // backing stayed >= 100% every block
	SpreadCaptured  int64 // total spread routed through treasury (micro-USD)
	ReserveFund     int64
	InsuranceFund   int64
	RewardsFund     int64
	TreasuryFund    int64
	BridgeAccepted  int64 // bridge-outs accepted (within cap)
	BridgeThrottled int64 // bridge-outs refused by the cap
	OracleHalts     int
}

func absI(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// price returns a deterministic oscillating mid for block b (a slow wave plus a
// sawtooth, so spreads see both calm and volatile blocks).
func price(base int64, b int) int64 {
	wave := int64((b%48)-24) * (base / 2000) // +/- ~1.2%
	saw := int64(b%6-3) * (base / 4000)      // small high-freq chop
	return base + wave + saw
}

// Run executes the scenario and returns a summary.
func Run(c Config) Result {
	o := pricing.NewOracle(8, 10, 500, 5000) // fresh, 5% outlier guard, 50% breaker
	r := &pricing.Reserve{MinBackingBps: pricing.BpsDenom}
	r.BridgeIn(c.SeedAsset) // 1:1 backed to start
	sp := pricing.SpreadParams{BaseBps: 15, VolMultBps: 15, MaxBps: 300}
	tr := &treasury.Treasury{}
	lim := bridge.NewLimiter(50*M, 24, 10*M, 12)

	res := Result{Blocks: c.Blocks, BackingHeld: true}
	prev := c.BasePrice

	for b := 0; b < c.Blocks; b++ {
		px := price(c.BasePrice, b)
		// three sources, two near px and one slightly off (still within guard)
		mid, err := o.Update([]int64{px, px + px/5000, px - px/8000})
		if err != nil {
			res.OracleHalts++
			o.Resume()
			prev = px
			continue
		}

		volBps := absI(px-prev) * pricing.BpsDenom / prev
		q := pricing.MakeQuote(mid, sp.EffectiveSpreadBps(volBps), 0)

		before := r.SpreadUSD
		for i := 0; i < c.RoundTrips; i++ {
			// round-trip keeps net supply flat -> backing preserved, spread earned twice
			if err := r.SellWrapped(c.TradeSize, q); err != nil {
				continue
			}
			if err := r.BuyWrapped(c.TradeSize, q); err != nil {
				// re-mint failed (backing) -> undo the sell by bridging the claim back
				r.WrappedSupply += c.TradeSize
			}
		}
		if gained := r.SpreadUSD - before; gained > 0 {
			a, _ := tr.Collect(gained, treasury.DefaultSpreadSplit)
			res.SpreadCaptured += a.Total()
		}

		// periodic bridge-out, throttled by the cap + delay queue
		if c.BridgeOutEvery > 0 && b%c.BridgeOutEvery == 0 {
			switch _, err := lim.RequestOut(c.BridgeOutSize, int64(b)); err {
			case nil:
				res.BridgeAccepted++
			case bridge.ErrCapExceeded:
				res.BridgeThrottled++
			}
		}

		if !r.Healthy() {
			res.BackingHeld = false
		}
		prev = px
	}

	res.FinalBackingBps = r.BackingBps()
	res.ReserveFund = tr.Reserve
	res.InsuranceFund = tr.Insurance
	res.RewardsFund = tr.Rewards
	res.TreasuryFund = tr.Treasury
	return res
}
