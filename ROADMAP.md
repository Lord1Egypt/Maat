# ROADMAP.md - Phased Execution Plan (v0.2)

> From Concept to a Sustainable Exchange.
> Detailed execution lives in [BUILD_PLAN.md](BUILD_PLAN.md); this is the high-level view.

> ⚠️ v0.2: launch is as a **market maker (mid ± spread)**, not a faucet. The old
> "bootstrap at 10% below market" step was the drain — removed. See [REDESIGN.md](REDESIGN.md).

---

## Phase 0: Concept & Design (NOW)
**Status: COMPLETE**

- [X] Define core economic model (corrected to spread-capture)
- [X] Full documentation set
- [X] REDESIGN.md (why v0.1 failed, how v0.2 works)
- [X] BUILD_PLAN.md (execution plan)
- [X] Build Python simulation (`maat_sim.py`)
- [X] Community feedback on the corrected model

---

## Phase 1: Simulation (Month 1) — STOP-GATE

**Goal:** Prove the reserve **grows** before any chain code.

- [X] Build the spread-capture simulation
- [X] Reserve growth under realistic two-sided flow
- [X] Inventory-skew stress + skew-pricing fix
- [X] Oracle-lag / flash-crash stress
- [X] 1,000 Monte Carlo paths; **pass = reserve grows in ≥95%**
- [X] Whitepaper with results

**Deliverable:** Green simulation. (PASSED)

---

## Phase 2: Cosmos SDK Chain (Month 2)

**Goal:** Testnet that quotes mid ± spread per block.

- [X] x/maat (native coin, staking, gov)
- [X] x/pegged (mint/burn wrapped, backing-checked)
- [X] **x/oracle (multi-source TWAP)** — required
- [X] x/market (fixed-per-block quote + settlement + spread accrual)
- [X] x/reserve (track real backing, auto-pause < 100%)
- [X] Unit + fuzz + sim-driven solvency tests
- [X] 4-validator testnet; verify wETH round-trip grows reserve

**Deliverable:** Testnet where the reserve provably grows. (COMPLETE)

---

## Phase 3: Ethereum Bridge (Month 2-3)

**Goal:** Real cross-chain flow, done safely.

- [ ] Integrate Wormhole (battle-tested, not custom)
- [ ] Sepolia: lock ETH → mint wETH; burn wETH → release ETH
- [ ] Daily caps, large-tx delay/cancel, per-bridge multisig
- [ ] **Independent audit of bridge integration**

**Deliverable:** Audited, capped wETH bridge.

---

## Phase 4: Guarded Mainnet Launch (Month 3)

**Goal:** First users — as a market maker, not a faucet.

- [ ] Seed reserve inventory
- [ ] Quote wETH at mid ± spread; pitch = MEV-free best execution + live reserve
- [ ] Live reserve dashboard; publish spread curve + cadence
- [ ] Airdrop to active DeFi users; conservative caps

**Metrics:** reserve trending up · backing ≥ 110% always · spread revenue > infra cost · no bridge incident

---

## Phase 5: Remittance Corridor + Asset Expansion (Month 4-8)

**Goal:** Recurring real revenue + multi-asset.

- [ ] Egypt $30B corridor: Arabic app, 0.3%, licensed ramp partners
- [ ] Add wBTC (audited threshold-sig), stables, wSOL, IBC
- [ ] Vol/inventory-adjusted spreads; whale tiers
- [ ] Reserve index (on-chain ETF, transparent NAV)

**Metrics:** $10M+ reserve · 5+ assets · $1M+ daily volume · revenue > $3K/day

---

## Phase 6: Governance Maturity (Month 8-12)

**Goal:** DAO-driven spreads, caps, listings.

- [ ] Full governance for spread/param proposals
- [ ] MAAT price feed for external DeFi (once deep + trusted)
- [ ] Insurance fund to 5% of reserve
- [ ] Bug bounty + security council rotation
- [ ] Legitimate Pharaoh Bonds (backed by real spread revenue)

**Metrics:** 50%+ staked · 20+ validators · 3+ proposals passed

---

## Phase 7: Dominance (Month 12+)

**Goal:** The fair-execution + cross-chain settlement layer.

- [ ] CosmWasm advanced products
- [ ] 50+ assets, mobile wallet, fiat ramps
- [ ] Partnerships with major DeFi + wallets

**Metrics:** $100M+ reserve · $10M+ daily volume · $50K+/day revenue

---

## Current Focus

```
Phase 0 > Phase 1(GATE) > Phase 2 > Phase 3 > Phase 4 > Phase 5 > Phase 6 > Phase 7
   [####################............]
                                    ^ YOU ARE HERE
```

Next concrete step: **Integrate Wormhole (EVM chains) for cross-chain flow (Phase 3).**
