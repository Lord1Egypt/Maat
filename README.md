# 𓁧 Ma'at — The Order Blockchain

> **"Bring Order to the Chaos"**
> *MEV-free, fixed-per-block cross-chain asset exchange — predictable execution, transparent reserves, sustainable spread.*

<p align="center">
  <a href="https://github.com/Lord1Egypt/Maat/actions/workflows/ci.yml"><img src="https://github.com/Lord1Egypt/Maat/actions/workflows/ci.yml/badge.svg"></a>
  <img src="https://img.shields.io/badge/Concept-v0.2-blue?style=flat-square">
  <img src="https://img.shields.io/badge/Ecosystem-Cosmos%20SDK-2CA5E0?style=flat-square&logo=cosmos&logoColor=white">
  <img src="https://img.shields.io/badge/Layer-L1%20App--Chain-success?style=flat-square">
  <img src="https://img.shields.io/badge/Phase-Concept%20%26%20Design-yellow?style=flat-square">
  <img src="https://img.shields.io/badge/🇪🇬-Made%20in%20Egypt-red?style=flat-square">
</p>

---

## 🏛️ What is Ma'at?

**Ma'at** is an L1 blockchain named after the ancient Egyptian goddess of truth, balance, order, and justice. It is a **cross-chain asset exchange that quotes a single price per block** — fixed *within* the block so there is no slippage and no MEV, but tracking the market so the protocol stays solvent and **earns the spread** on every trade.

Think of it as the **central bank's FX window for crypto**: you always know the price, the price is fair, nobody front-runs you — and the house keeps a tiny, honest spread instead of giving money away.

> ⚠️ **v0.2 note:** This replaces the original "fixed off-market price + guaranteed user arbitrage" model, which was structurally unsustainable (it drained the reserve by design — see [REDESIGN.md](REDESIGN.md) for the full proof). Same brand, same philosophy, sound economics.

---

## 🔥 The Core Innovation

### Traditional Crypto Markets (Chaos)
```
Per-block price moves every trade → slippage, sandwiches, MEV, oracle attacks
```

### Ma'at (Order)
```
One clearing price per block = oracle mid ± spread → no slippage, no MEV, reserve grows
```

| Component | Market (AMM/DEX) | Ma'at |
|-----------|------------------|-------|
| **Native Coin (MAAT)** | Free market | Free market (traded externally) |
| **Wrapped BTC/ETH/USDC** | Price moves every trade | **One fixed price per block (mid ± spread)** |
| **Price source** | AMM curve / order book | **Robust multi-source TWAP oracle** |
| **Slippage** | Yes, grows with size | **Zero — fixed per block** |
| **MEV** | Sandwiches, frontrunning | **Impossible — no intra-block price to exploit** |
| **Who keeps the spread?** | LPs (with impermanent loss) | **The protocol reserve (no IL, it grows)** |

---

## 💡 The Spread Engine (how the reserve *grows*)

Ma'at quotes a two-sided price around the market mid and **keeps the spread**, exactly like a real FX bureau or market maker:

```
Market BTC mid = $100,000
Ma'at BUY  wBTC @ $99,850  (mid − 0.15%)
Ma'at SELL wBTC @ $100,150 (mid + 0.15%)

Every round-trip → ~0.30% to the reserve.
Trade flow is two-sided → the reserve compounds.
```

This is the opposite of the old model. Instead of "guaranteed profit for users" (which was guaranteed *loss* for the reserve), Ma'at offers **guaranteed best execution** and keeps an honest spread — so it can run forever.

**The real value proposition:**
> **Guaranteed best execution. Guaranteed zero MEV. Guaranteed transparent 1:1 backing.**

---

## 🎯 What Success Looks Like

| Phase | What Happens | Result |
|-------|-------------|--------|
| ✅ Prove | Simulation shows reserve grows in ≥95% of paths | Green light to build |
| ✅ Bootstrap | Launch wETH at mid ± spread, MEV-free | Traders come for execution quality |
| ✅ Remit | Egypt's $30B diaspora corridor at 0.3% | Real recurring revenue |
| ✅ Scale | BTC, USDC, SOL, IBC + a transparent reserve index | Multi-asset MEV-free hub |
| ✅ Maturity | DAO-governed spreads, deep reserves | **The fair-execution layer of crypto** |

---

## 📜 Documentation

| Document | What It Covers |
|----------|---------------|
| [WHITEPAPER.md](WHITEPAPER.md) | The synthesis — model, proof, validation, on-chain core |
| [REDESIGN.md](REDESIGN.md) | **Read first** — why v0.1 was broken and how v0.2 fixes it |
| [STATUS.md](STATUS.md) | Current project state & continuation handoff |
| [BUILD_PLAN.md](BUILD_PLAN.md) | How to start + the phased execution plan |
| [VISION.md](VISION.md) | The full philosophy — why this matters |
| [ECONOMICS.md](ECONOMICS.md) | Tokenomics, reserve model, spread mechanics |
| [PLANNED_ECONOMY.md](PLANNED_ECONOMY.md) | How pricing works in detail |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Technical design: Cosmos SDK modules, oracle, bridges |
| [RISKS.md](RISKS.md) | Every risk with mitigations (honest) |
| [ROADMAP.md](ROADMAP.md) | Phased execution plan |
| [COMPETITION.md](COMPETITION.md) | Market landscape and differentiation |
| [SIMULATION.md](SIMULATION.md) | Economic model simulation |

---

## ✅ Verification (this is proven, not just claimed)

Every economic and safety claim is enforced by code in CI:

```bash
make test          # runs everything below
```

| What | Where | Status |
|------|-------|--------|
| Reserve **grows** (1000-path Monte Carlo, synthetic + real ETH) | `simulation/` | gate fails CI if not |
| Deterministic on-chain cores: `x/maat`, `x/oracle`, `x/market`, `x/reserve`, `x/bridge`, `x/treasury`, `x/gov` | `chain/` | Go tests + fuzzing |
| End-to-end run (backing held 100%, spread captured, bridge throttled) | `chain/scenario`, `chain/cmd/demo` | `go run ./cmd/demo` |
| No "guaranteed risk-free profit" claims resurface | docs-integrity job | grep-enforced |

The Python sim proves the **policy**; the Go core proves the **on-chain arithmetic**.
See [WHITEPAPER.md](WHITEPAPER.md) and [chain/README.md](chain/README.md).

---

## ⚠️ Honest Warning

**This project is NOT a typical DeFi protocol.** It is an MEV-free, oracle-priced exchange whose reserve grows on spread. The hard parts are real and we won't hide them:
- **Bridge security** is crypto's #1 attack vector — we use battle-tested bridges with caps and insurance.
- **Oracle robustness** is mandatory — multi-source TWAP with deviation circuit breakers.
- **Solvency** is proven by simulation *before* any chain code is written.

We document everything — including what could still go wrong.

---

<p align="center">
  Made with 𓁧 by <b>Mohamed Mounir (Lord1Egypt)</b><br>
  <i>"السوق مش دايماً صح — أحياناً النظام أحسن"</i>
</p>
