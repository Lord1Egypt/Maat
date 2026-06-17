# PLANNED_ECONOMY.md

> How Pricing Works — The Detailed Mechanics (v0.2)

> ⚠️ Corrected from v0.1. Ma'at does **not** quote an arbitrary fixed price, and it does
> **not** run without an oracle. It quotes **oracle mid ± spread, fixed per block**.
> See [REDESIGN.md](REDESIGN.md) for why the old "we ARE the oracle" model self-destructs.

---

## 1. The Core Mechanism

### 1.1 The Protocol Reads a Robust Oracle, Then Fixes the Price Per Block

Ma'at uses a hardened multi-source TWAP oracle (the price *source*), then exposes a single
clearing price per block (the *execution* guarantee). Every block stores:

```
state.asset_quote["wBTC"] = {
    mid:        100000,     // robust multi-source TWAP
    spread_bps: 15,         // each side, governance-set, vol/inventory-adjusted
    buy:        99850,      // mid * (1 - spread)   protocol buys from user
    sell:       100150,     // mid * (1 + spread)   protocol sells to user
    block:      123456,
}
```

The price is **fixed for the whole block** → no intra-block reordering can change it →
no slippage, no sandwich, no MEV. It **tracks the market** → the reserve never bleeds.

### 1.2 Mint/Burn at the Quote

```
User sells 1 wBTC (protocol buys):
  1. Read buy price (99,850)
  2. Burn wBTC, credit user 99,850 (MAAT-equivalent)
  3. Reserve keeps the asset acquired below mid -> +spread

User buys 1 wBTC (protocol sells):
  1. Read sell price (100,150)
  2. Check backing exists
  3. Debit 100,150, release/mint 1 wBTC -> +spread
```

No slippage. No spread *given away*. No MEV. The protocol earns the spread.

## 2. The Oracle (mandatory, not optional)

```
mid = TWAP across >= 3 independent venues (e.g., Binance, Coinbase, on-chain DEX TWAP)
- Median of sources to resist a single manipulated feed
- Time-weighted to resist flash spikes
- Deviation guard: if a source deviates > X% it is dropped
- Circuit breaker: if mid moves > 50% in 1h, freeze quotes for review
```

"Being the oracle" is dropped entirely — that claim is what guaranteed insolvency.

## 3. Predictability: The Spread Schedule

What Ma'at commits to in advance is **how it prices**, not a future off-market price:

```
Published:
- Base spread per asset
- Volatility-scaling formula
- Inventory-skew formula
- Update cadence (per block) and circuit-breaker thresholds
```

This gives users total clarity about execution **without** writing free options to the
market (pre-announcing a future fixed price is a free call/put — see REDESIGN §1.3).

## 4. Quote Adjustment Types

| Type | Description | Use Case |
|------|-------------|---------|
| Flat spread | Fixed bps each side | Calm markets |
| Volatility-scaled | Spread grows with realized vol | Turbulent markets |
| Inventory-skewed | Quotes tilt to pull inventory back to target | One-sided flow |
| Emergency widen | Spread jumps to deter toxic flow | Stress |
| Emergency Freeze | Pause all quotes 48h | Circuit breaker |

## 5. Governance

### 5.1 Who Can Propose

Any MAAT staker with >= 100,000 MAAT delegated.

### 5.2 Proposal Format

```json
{
  "title": "wBTC base spread: 15bps -> 12bps",
  "asset": "wBTC",
  "param": "base_spread_bps",
  "current": 15,
  "new": 12,
  "rationale": "Deep reserve, high competition, capture more volume"
}
```

### 5.3 Parameters

| Parameter | Value |
|-----------|-------|
| Voting period | 7 days |
| Quorum | 30% of staked MAAT |
| Approval | >= 60% of votes |
| Emergency fast-track | >= 80%, 1-day vote |

### 5.4 Emergency Powers

Security Council (7-of-12 multi-sig) can:
- Pause all quotes (circuit breaker)
- Force-widen spreads
- Trigger emergency reserve rebalancing

Powers are audited, time-limited, and publicly logged.
