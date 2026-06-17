# SIMULATION.md - Economic Model (v0.2)

> Python simulation to prove the economics **before** writing blockchain code.
> The pass condition changed: the reserve must **grow**, not merely "survive."

> ⚠️ v0.1 asked "can the reserve survive a bank run?" — the honest v0.1 answer was "no,
> it drains by design." v0.2 asks the right question: does the spread-capture reserve
> compound under realistic flow? See [REDESIGN.md](REDESIGN.md).

---

## Purpose

Before writing a single line of Cosmos SDK, we simulate the Ma'at economy to answer:

1. Does the spread engine make the reserve **grow** under realistic two-sided flow?
2. How bad does **inventory imbalance** get under one-sided flow, and does skew-pricing fix it?
3. What is the worst-case loss from **oracle lag** during a flash crash?
4. What spread maximizes (volume × spread) revenue without killing volume?

---

## Simulation Model

The simulation models:
- **Market mid**: real historical BTC/ETH price data + synthetic shocks
- **Ma'at quote**: mid ± spread, fixed per block, with vol/inventory adjustment
- **Oracle lag**: quote follows mid with realistic delay
- **Flow**: arb bots (close gap to mid), retail (random), one-sided stress
- **Reserve & inventory**: real backing + spread accrual + skew

### Key Variables

| Variable | Range | Default |
|----------|-------|---------|
| Base spread (each side) | 5-50 bps | 15 bps |
| Oracle lag | 1-30 blocks | 5 blocks |
| Vol multiplier | 1-5x | dynamic |
| Inventory cap | 5-50% of reserve | 20% |
| Reserve buffer target | 100-150% | 115% |

---

## How to Run

```bash
cd simulation
pip install numpy pandas matplotlib
python maat_sim.py --days 90 --base-spread 0.0015 --paths 1000 --output charts/
```

---

## Planned Charts

1. Reserve value over time (must trend **up**)
2. Daily spread revenue vs base spread
3. Inventory skew over time (with vs without skew-pricing)
4. Worst-case loss during a 50% flash crash (oracle lag stress)
5. Revenue surface: spread × volume

---

## Pass / Fail Gate (the stop-gate)

| Test | Pass condition |
|------|----------------|
| 1,000 Monte Carlo paths | Reserve ends ≥ start in **≥95%** of paths |
| Flash-crash stress | Backing stays ≥ 100% with breakers on |
| One-sided flow | Skew-pricing keeps inventory under cap |
| Revenue | Net spread revenue > simulated infra cost |

**If no parameter set passes → the model is wrong. Do not build. Iterate first.**

---

## Next Steps

1. Implement `maat_sim.py`
2. Run 1,000 Monte Carlo paths
3. Find the optimal spread/lag/skew parameter set
4. Document in whitepaper
5. Only then build the real chain (see [BUILD_PLAN.md](BUILD_PLAN.md))
