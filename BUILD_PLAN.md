# BUILD_PLAN.md — How to Start & The Phased Plan

> The execution plan, rebuilt on the corrected model in [REDESIGN.md](REDESIGN.md).
> Old ROADMAP.md assumed the "fixed off-market price" model that drains the reserve — this supersedes it.

---

## Golden Rule

**Do not write a single line of chain code until the simulation proves the reserve GROWS
(not drains) under the spread-capture model.** The old roadmap had this backwards: it
planned to launch with "10% below market" discounts, which is the drain. Prove
sustainability first.

---

## Phase 0 — Validate the Economics (Weeks 1–4)  ← START HERE

**Goal:** A Python simulation proving the spread-capture model is sustainable.

- [ ] Build `simulation/maat_sim.py`: real BTC/ETH price history + bid/ask = mid ± spread.
- [ ] Model arb bots, retail flow, inventory imbalance, and oracle lag.
- [ ] Output charts: reserve value over time (must trend **up**), daily spread revenue,
      worst-case inventory skew, behavior under a 50% flash crash.
- [ ] Run 1,000 Monte Carlo paths. **Pass condition: reserve ends ≥ starting value in ≥95% of paths.**
- [ ] If a parameter set passes → lock it. If none pass → the model is wrong, iterate before coding.

**Deliverable:** `simulation/` + `whitepaper.md` with results. **Stop-gate:** no green sim, no build.

---

## Phase 1 — Single-Asset Testnet (Weeks 5–10)

**Goal:** A Cosmos SDK chain that quotes wETH at oracle-mid ± spread, per block.

- [ ] Init Cosmos SDK v0.50 + CometBFT, 4 local validators.
- [ ] `x/maat` (native coin, staking, gov), `x/pegged` (mint/burn wETH).
- [ ] `x/oracle` — multi-source TWAP feed (this is REQUIRED; "we are the oracle" is dropped).
- [ ] `x/market` — batch settlement: one clearing price per block = mid ± spread. No AMM.
- [ ] `x/reserve` — track real backing; circuit-break if backing < 100%.
- [ ] Unit + fuzz tests on every module. Verify spread revenue accrues in sim-driven tests.

**Deliverable:** Testnet where wETH round-trips and the reserve provably grows.

---

## Phase 2 — One Bridge, Done Right (Weeks 8–14, overlaps)

**Goal:** Real cross-chain ETH ↔ wETH via a battle-tested bridge (bridges are crypto's #1 hack vector).

- [ ] Integrate Wormhole (do NOT roll your own bridge).
- [ ] Sepolia: lock ETH → mint wETH; burn wETH → release ETH.
- [ ] Daily bridge cap + per-tx delay on large amounts + separate multisig.
- [ ] **Independent audit of the bridge integration before any mainnet value.**

**Deliverable:** Audited, capped wETH bridge on testnet → limited mainnet.

---

## Phase 3 — Guarded Mainnet Launch (Months 4–5)

**Goal:** First real users. **Launch as a market maker, not a faucet.**

- [ ] Seed reserve inventory (your own ETH + raise via legit bond, see §Funding).
- [ ] Quote wETH at mid ± 0.15%. Pitch: **"zero-MEV, zero-slippage execution + live reserve."**
- [ ] NOT "10% below market" — that is the drain. Marketing = execution quality, not free money.
- [ ] Public live reserve dashboard; publish quote curve + update cadence.
- [ ] Conservative caps; raise them as confidence grows.

**Success metrics:** reserve trending up · spread revenue > infra cost · backing ≥ 110% always · no bridge incident.

---

## Phase 4 — The Remittance Rail (Months 5–8)

**Goal:** Ship the strongest, most defensible product — Egypt's $30B corridor.

- [ ] Crawling-peg stable transfer (Egypt literally runs a crawling-peg FX regime — lean into it).
- [ ] Arabic app, 0.3% fee, instant settlement.
- [ ] Licensed fintech / VASP partners for fiat on/off ramp (regulatory-safe path).

**Deliverable:** Diaspora users moving real money. Real, recurring fee revenue.

---

## Phase 5 — Multi-Asset + Reserve Index (Months 8–14)

**Goal:** Expand assets; turn the "world reserve" vision into an on-chain ETF (NAV-priced), not a price-fixer.

- [ ] Add wBTC (threshold-sig or audited bridge), stables, wSOL, IBC assets.
- [ ] Dynamic spreads by volatility + inventory skew.
- [ ] Scheduled, announced **rebalances** of a transparent basket; charge a management fee.
- [ ] Legitimate Pharaoh Bonds — now backed by REAL spread/fee revenue.

**Success metrics:** $10M+ reserve · multi-asset live · revenue > $3K/day · backing healthy.

---

## Phase 6 — Decentralize & Harden (Months 12–18)

- [ ] DAO governance for spreads, caps, listings.
- [ ] 20+ validators, 50%+ staked.
- [ ] Insurance fund ≥ 5% of reserve; rotating security council.
- [ ] Continuous bug bounty; annual re-audits.

---

## Funding (corrected — no "guaranteed profit" framing)

| Stage | Source | Note |
|---|---|---|
| MVP (~$50K) | Self-fund | sim + testnet + 1 audit |
| Seed inventory | Bonds backed by real spread revenue | only AFTER Phase 0 green |
| Growth | Strategic/VC or revenue | profitable on fees once rail + spread live |

**Never** market "risk-free guaranteed returns" — that is a security and a legal landmine
(RISKS Risk #3). Sell predictability and execution quality.

---

## The Critical Path (one line)

```
Phase 0 sim proves reserve grows  →  testnet (oracle ± spread)  →  audited bridge
 →  guarded mainnet (market maker, not faucet)  →  remittance rail  →  multi-asset + index  →  decentralize
```

**You are here:** Phase 0. Next concrete step → build `simulation/maat_sim.py` and run the 1,000-path test.
Do not skip the stop-gate.
