# Ma'at Simulation — The Phase 0 Stop-Gate

This is the gate that must go **green before any blockchain code is written**
(see [../BUILD_PLAN.md](../BUILD_PLAN.md)). It proves the v0.2 spread-capture
economics let the reserve **grow**, and that the abandoned v0.1 fixed-off-market
model **drains** (Terra-style).

## Run it

```bash
# zero dependencies — pure Python standard library
python3 maat_sim.py                     # default 1000-path gate
python3 maat_sim.py --paths 5000 --days 180
python3 maat_sim.py --charts            # also write PNGs (needs matplotlib)
```

Exit code `0` = gate passed, `1` = gate failed. CI runs this on every push.

## What it models (per block, stepped hourly)

| Piece | Meaning |
|-------|---------|
| `true_mid` | real market price (geometric Brownian motion + jumps) |
| `oracle` | EMA of the true mid → models oracle freshness (`--oracle-alpha`) |
| `delta` | protocol mispricing = `(true_mid - oracle) / oracle` |
| `quote` | `buy = oracle*(1-s)`, `sell = oracle*(1+s)` |
| revenue | retail flow pays the spread `s` |
| adverse selection | when `|delta| > s`, arbs hit the stale side and extract the excess |
| flash crash | a `--crash-size` gap stress-tests whether the buffer holds 100% backing |

## The gate (all must pass)

1. v0.2 reserve grows in **≥95%** of paths
2. v0.2 backing never breaks 100%
3. v0.2 inventory stays under cap (skew-pricing works)
4. v0.2 spread revenue beats infra cost
5. flash-crash buffer absorbs the shock
6. v0.1 demonstrably drains (sanity check the model isn't rigged)

## Key knobs

`--base-spread` (per-side spread) · `--oracle-alpha` (freshness; low = laggy/risky) ·
`--annual-vol` · `--arb-cap` · `--retail-vol` · `--buffer-start` (over-collateralization) ·
`--crash-size` · `--v1-discount`. The gate is **real**: a laggy oracle + near-zero
spread makes it fail, as it should.
