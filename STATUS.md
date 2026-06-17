# STATUS.md — Project State & Continuation Handoff

> **Read this first if you're picking up Ma'at cold** (new session, after a context
> reset, or a weekly-limit interruption). It captures where we are, what was decided
> and why, and exactly what to do next.

**Last updated:** 2026-06-17
**Phase:** 0 (Concept & Design) — **stop-gate PASSED (synthetic + real ETH data)**, cleared to build chain code.

---

## 1. TL;DR — the corrected thesis

Ma'at's original v0.1 idea (**fixed off-market price + "guaranteed risk-free arbitrage"
for users**) was structurally unsustainable: the user's "guaranteed profit" was paid out
of the reserve on every trade — the same coins — so the reserve drained to zero by
construction, exactly like Terra/UST. This was abandoned.

**v0.2 (current model):** quote **oracle mid ± spread, fixed per block**.
- Fixed *within* the block → no slippage, no MEV (the good part, kept).
- Tracks the market via a hardened multi-source TWAP oracle → always solvent.
- The protocol **earns** the spread → the reserve **grows** instead of draining.
- Value prop flipped from "free money" → **guaranteed best execution + no MEV + transparent backing.**

Full proof and reasoning: [REDESIGN.md](REDESIGN.md).

---

## 2. What has been done (merged to `main`)

| PR | What | Status |
|----|------|--------|
| #1 | `fix:` replace unsustainable fixed-price model with spread-capture (v0.2). Rewrote all 10 docs + added REDESIGN.md & BUILD_PLAN.md | ✅ Merged |
| #2 | `feat:` economic stop-gate simulation + first CI/CD pipeline | ✅ Merged |

**Workflow in use:** branch → PR → merge (never direct push to `main`). CI must be green
before merge.

---

## 3. Simulation result (the Phase 0 stop-gate)

`simulation/maat_sim.py` (pure stdlib, zero deps). 1000 Monte Carlo paths, same paths for
both models:

| Model | Reserve grows | Median outcome |
|-------|---------------|----------------|
| **v0.2 spread-capture** | **100%** of paths | **+4.5%** |
| v0.1 fixed off-market (abandoned) | 0% of paths | **−86%** (drains) |

- Backing stayed ≥100% in 100% of paths; flash-crash absorbed by the buffer.
- The gate has **teeth**: a laggy oracle + near-zero spread correctly **fails** it (exit 1).
- **Now also validated on REAL ETH data** (90d hourly, CoinGecko): calibrated to 45.5%
  annual vol / ~1%-per-hour jumps / −37.9% max drawdown; block-bootstrap gate still PASSES
  (+4.45% median, v0.1 −86%). Tools: `fetch_data.py`, `calibrate.py`; replay via `--price-csv`.
- Run it: `python3 simulation/maat_sim.py --paths 1000` (exit 0 = pass, 1 = fail).

---

## 4. CI/CD (enforced on every push/PR)

`.github/workflows/ci.yml` — two jobs:
1. **economics-gate** — compiles + runs the 1000-path stop-gate (~20s in CI). Red build if economics break.
2. **docs-integrity** — fails CI if abandoned v0.1 claims ("guaranteed risk-free arb",
   "we are the oracle") reappear as live copy. (Citing/warning about them is allowed.)

---

## 5. NEXT STEPS

- **B — Calibrate vs real ETH data** ✅ DONE (PR #4). Params locked: annual_vol 0.4554,
  jump 0.0291, crash 0.38; gate passes on real-data bootstrap.
- **A — Scaffold the chain (Phase 2)** ← IN PROGRESS / NEXT. Start with the core pricing
  math: `x/oracle` (hardened multi-source TWAP) and `x/market` (fixed-per-block quote +
  settlement + spread accrual), since every other module depends on them. Then `x/pegged`,
  `x/reserve`, `x/bridge`. Go 1.22 is available in the environment.

(See [BUILD_PLAN.md](BUILD_PLAN.md) for the full phased plan and stop-gate rules.)

---

## 6. File map (where everything lives)

| File | Purpose |
|------|---------|
| `REDESIGN.md` | Why v0.1 fails + the corrected v0.2 model (the core reasoning) |
| `BUILD_PLAN.md` | Phased execution plan + the simulation stop-gate rule |
| `STATUS.md` | **This file** — current state & continuation handoff |
| `ECONOMICS.md` / `PLANNED_ECONOMY.md` | Spread mechanics, reserve, pricing detail |
| `ARCHITECTURE.md` | Cosmos SDK modules (incl. required `x/oracle`) |
| `RISKS.md` | Risk register (Risk #0 = structural drain, now fixed) |
| `VISION.md` / `STRATEGY.md` / `COMPETITION.md` | Philosophy, growth plan, landscape |
| `SIMULATION.md` / `simulation/` | The stop-gate sim + how to run it |
| `.github/workflows/ci.yml` | CI/CD: economics-gate + docs-integrity |

---

## 7. Hard rules to honor when resuming

1. **No chain code until the stop-gate is green** (it is — keep it that way).
2. **Never reintroduce the v0.1 model** (off-market fixed price / "we are the oracle" /
   guaranteed user arb). The oracle is mandatory; the protocol earns the spread.
3. **branch → PR → merge**, CI green before merge.
4. **Zero/standard-library dependencies** where possible (sim is pure stdlib — keep it).
5. Don't market "guaranteed/risk-free profit" — it's a security and a legal landmine.
