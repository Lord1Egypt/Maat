# chain/ — Ma'at on-chain core

The deterministic heart of Ma'at's pricing, written in plain Go with **integer
fixed-point math only** (no floats — safe for CometBFT consensus, where every
validator must compute byte-identical results).

This package is intentionally **decoupled from the Cosmos SDK** so the
risk-bearing logic can be exhaustively unit-tested on its own (`go test ./...`,
zero external dependencies). The Cosmos modules wrap it; they don't reimplement
it.

## What's here

| File | Maps to module | Responsibility |
|------|----------------|----------------|
| `token/token.go` | `x/maat` | MAAT vesting schedule (ECONOMICS distribution) + bounded inflation / block-reward math. |
| `pricing/oracle.go` | `x/oracle` | Median of de-outliered sources → integer-EMA TWAP `Mid`, with a deviation guard and a circuit breaker that halts on implausible jumps. |
| `pricing/market.go` | `x/market` + `x/reserve` | `MakeQuote` (mid ± vol/skew-adjusted spread, fixed per block), `BuyWrapped`/`SellWrapped` (spread accrual + backing checks), `BridgeIn`/`BridgeOut` (1:1 custody). |
| `bridge/limits.go` | `x/bridge` | Bridge-out safety (Risk #1): rolling per-asset withdrawal cap + delay queue (large withdrawals held `DelayBlocks`, cancellable by governance). |
| `treasury/fees.go` | `x/treasury` | Deterministic spread/fee distribution (ECONOMICS.md splits); rounding remainder → reserve so no micro-unit leaks. |
| `scenario/scenario.go` | (integration) | Wires all cores through a deterministic multi-block run; asserted by `scenario_test.go`. |
| `cmd/demo/main.go` | (demo) | `go run ./cmd/demo` — runnable end-to-end proof the cores work together. |

Units: prices in **micro-USD** (1 USD = 1_000_000), spreads/ratios in **bps**.

## Run the tests

```bash
cd chain
go vet ./...
go test ./... -v
go run ./cmd/demo      # human-readable end-to-end proof

# or from the repo root, run EVERYTHING (Go + economic gate):
make test
```

## How the Cosmos keepers will use it (integration sketch)

```go
// x/oracle keeper — called in EndBlock after collecting feeder votes
mid, err := oracle.Update(sourcePricesFromVotes)   // pure, deterministic
k.SetMid(ctx, denom, mid)                          // SDK state write

// x/market keeper — MsgSwap handler
mid := k.oracle.GetMid(ctx, denom)
q := pricing.MakeQuote(mid, sp.EffectiveSpreadBps(volSignal), invSkewBps)
if err := reserve.BuyWrapped(amt, q); err != nil { return err }
k.SaveReserve(ctx, denom, reserve)                 // spread accrued, backing checked
```

The float Python model in `../simulation/` and this integer Go core encode the
**same economics** (oracle mid ± spread, protocol earns the spread, backing
never breaks) — the simulation proves the policy; this proves the on-chain
arithmetic. See [../REDESIGN.md](../REDESIGN.md) and [../ARCHITECTURE.md](../ARCHITECTURE.md).

## Not yet here (next)

Full SDK plumbing: `MsgServer`/proto types, genesis, params, EndBlock wiring,
keeper state stores, and the `x/bridge` Wormhole integration. Those are
mechanical wrappers around this verified core — see **[MODULE_SPEC.md](MODULE_SPEC.md)**
for the exact per-module wiring guide.
