# Ma'at: An MEV-Free, Reserve-Backed Cross-Chain Exchange with Per-Block Fixed Pricing

**Mohamed Mounir (Lord1Egypt)** · v0.2 · 2026-06-17

> 𓁧 *"Bring order to the chaos."* This whitepaper synthesizes the corrected Ma'at
> design and the evidence behind it. It supersedes the original v0.1 concept.

---

## Abstract

Ma'at is a Cosmos-SDK Layer-1 that executes cross-chain asset swaps at **one
clearing price per block** — fixed *within* the block (eliminating slippage and
MEV) but tracking the market through a hardened oracle (preserving solvency). The
protocol quotes a two-sided price around the oracle mid and **earns the spread**,
so its reserve **grows** under normal trading rather than draining. Every wrapped
asset is backed 1:1+ by real reserves held in multi-sig, published live. We show
— in Monte-Carlo simulation on both synthetic and **real ETH data**, and in a
deterministic on-chain core with executable tests — that the model is solvent,
that its protections (oracle circuit breaker, bridge caps, leak-free fee
accounting) hold, and that the abandoned v0.1 "fixed off-market price" model
drains by construction.

---

## 1. The Problem

Crypto market microstructure prices **per trade**, which produces slippage, MEV
(~$400M+/yr), impermanent loss for LPs, and a heavy dependence on oracles (a
$1.6B+ attack surface in 2023). Users pay a "chaos tax" on every interaction.

## 2. What Doesn't Work: The v0.1 Trap (and why we abandoned it)

The original Ma'at proposed a **fixed off-market price** with "guaranteed
risk-free arbitrage" for users. This is structurally unsustainable: the user's
guaranteed profit is paid out of the reserve on every trade — the *same coins* —
so the reserve drains to zero by construction. It is the Terra/UST reflexive loop
with extra steps. Over-collateralization only changes how many days the drain
takes. (Full proof: `REDESIGN.md` §1.)

**Corollary.** The abandoned v0.1 claim that the protocol could *be* the oracle
(no external price needed) guarantees insolvency: a price untethered from the
market is a faucet. A market-tracking price *requires* an oracle. No third option.

## 3. The Model: Per-Block Fixed Price, Spread Captured

Each block, for each asset:

```
mid   = robust multi-source TWAP (the oracle)
spread= base × volatility-multiplier, tilted by inventory skew (governance-set)
buy   = mid × (1 − spread)     # protocol BUYS the asset from users
sell  = mid × (1 + spread)     # protocol SELLS the asset to users
```

- **Fixed within the block** ⇒ transaction ordering cannot change anyone's price
  ⇒ no slippage, no sandwich, no MEV.
- **Centered on the market mid** ⇒ the protocol never sells below / buys above the
  market ⇒ the reserve does not bleed.
- **The protocol keeps the spread** ⇒ the reserve compounds.

The honest value proposition is therefore **guaranteed best execution +
guaranteed no MEV + guaranteed transparent backing** — not free money (which is
always someone's loss, and in v0.1 was the protocol's).

## 4. Reserve & Solvency

Every wrapped token is a 1:1 claim on a real asset in multi-sig custody, plus an
over-collateralization buffer (default 15%) that absorbs oracle lag and one-sided
flow. Trading accrues spread to the buffer. A circuit breaker halts trading if
backing falls below 100%. Reserve addresses are published and verified monthly.

## 5. Hardened Oracle

`mid` is the **median** of ≥3 independent venues (resisting a single manipulated
feed), **time-weighted** (resisting flash spikes), with a **deviation guard**
(drop sources off by >X%) and a **circuit breaker** that halts on an implausible
jump. Implemented in integer fixed-point math for consensus determinism.

## 6. Protections (Risk Register, abridged)

| # | Risk | Status / mitigation |
|---|------|---------------------|
| 0 | Structural reserve drain (the v0.1 killer) | **Eliminated by design** (spread capture). |
| 1 | Bridge hack (#1 attack vector) | Battle-tested bridges, per-asset withdrawal caps, **delay queue with cancel**, insurance fund. |
| 2 | Oracle manipulation | Median of ≥3, TWAP, deviation guard, breaker. |
| 3 | Securities-law exposure | No "guaranteed profit" marketing; licensed ramp partners. |
| 4 | Inventory imbalance | Skew-pricing + spread widening + caps. |

## 7. Economic Validation (simulation)

`simulation/maat_sim.py` (pure stdlib) runs 1000 Monte-Carlo paths and, on the
same paths, the abandoned v0.1 model as a control. The gate **fails** the build
if the reserve does not grow.

| Model | Reserve grows | Median outcome |
|-------|---------------|----------------|
| **v0.2 spread-capture** | **100%** of paths | **+4.5%** |
| v0.1 fixed off-market (control) | 0% | **−86%** (drains) |

Calibrated against **real ETH** (90d hourly: 45.5% annual vol, ~1%/hr jumps,
−37.9% max drawdown), the real-data block-bootstrap gate still passes (+4.45%).
Backing held ≥100% in 100% of paths; the flash-crash buffer absorbed the shock.
The gate has teeth: a laggy oracle + near-zero spread correctly fails it.

## 8. On-Chain Core (executable)

`chain/` implements the risk-bearing logic in deterministic integer Go,
decoupled from the SDK so it is exhaustively unit-tested:

- `pricing/oracle.go` — median + deviation guard + EMA TWAP + breaker
- `pricing/market.go` — quote + buy/sell with spread accrual + backing checks
- `bridge/limits.go` — withdrawal caps + cancellable delay queue
- `treasury/fees.go` — fee splits with remainder routed to reserve (no leak)
- `scenario/` + `cmd/demo` — a 720-block end-to-end run; demo shows backing held
  at 100%, ~$21.6k spread captured and split, 60/240 bridge-outs throttled.

The Python sim proves the *policy*; the Go core proves the *on-chain
arithmetic*. They encode the same economics.

## 9. Roadmap

Phase 0 (concept + stop-gate) ✅ → Phase 2 SDK modules wrapping the verified
cores → audited Wormhole bridge → guarded mainnet (as a market maker) →
remittance corridor (Egypt's $30B) → multi-asset + transparent reserve index →
decentralize. Full plan: `BUILD_PLAN.md`.

## 10. Conclusion

Crypto does not need more chaos; it needs fair execution backed by discipline.
Ma'at delivers one fair price per block, a hardened oracle, real 1:1+ reserves,
and an honest spread that lets the protocol survive forever. The design is proven
solvent in simulation on real data and in an executable on-chain core. Order, not
chaos.

---

*References: `REDESIGN.md`, `ARCHITECTURE.md`, `ECONOMICS.md`, `RISKS.md`,
`SIMULATION.md`, `BUILD_PLAN.md`, `STATUS.md`, and the `chain/` + `simulation/`
sources.*
