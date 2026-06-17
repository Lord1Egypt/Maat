# ECONOMICS.md

> The Ma'at Economic Model (v0.2 — spread-capture)

Tokenomics, Reserve Model, Supply Mechanics, and Revenue

> ⚠️ Corrected from v0.1. The old "fixed off-market price + guaranteed user arbitrage"
> model drained the reserve by construction — see [REDESIGN.md](REDESIGN.md) §1 for the proof.
> Ma'at now quotes **oracle mid ± spread per block** and **earns** the spread.

---

## 1. Token Overview

| Token | Type | Supply | Price Model |
|-------|------|--------|-------------|
| MAAT | Native coin (gas + gov) | Fixed: 1B | Free market |
| wBTC/wETH/wUSDC | Wrapped assets | Demand-driven | **Oracle mid ± spread, fixed per block** |

## 2. MAAT Tokenomics

### 2.1 Distribution

| Allocation | % | Amount | Vesting |
|------------|---|--------|---------|
| Community Sale | 30% | 300M | 10% TGE, 90% over 12mo |
| Treasury Reserve | 25% | 250M | Multi-sig, DAO-controlled |
| Validator Rewards | 20% | 200M | Inflated over 5 years |
| Team | 15% | 150M | 12mo cliff, 24mo linear |
| Advisors/Partners | 5% | 50M | 6mo cliff, 18mo linear |
| Airdrop | 5% | 50M | Initial users + cross-chain communities |

### 2.2 Utility

- **Gas**: Every transaction on Ma'at costs MAAT
- **Governance**: Vote on spreads, reserve ratios, listings, bridge params
- **Staking**: Validators stake MAAT for consensus
- **Fee Discount**: Stakers get reduced swap/bridge fees
- **Spread Dividend**: A share of spread revenue can flow to stakers (gov-set)

> Note: MAAT is **never** the backing for wrapped assets. Backing is always real assets
> (1:1+). Using your own token as collateral is the Terra/UST reflexive trap.

## 3. The Reserve Model

### 3.1 Full-Backed Reserve (1:1+)

Every wrapped asset is backed 1:1 by real assets in multi-sig:

```
wBTC minted = 100 --> 100 BTC held in reserve
wETH minted = 500 --> 500 ETH held in reserve
```

Because Ma'at trades at **mid ± spread** (not an off-market price), normal trading does
not drain backing — it adds spread revenue to it.

### 3.2 Over-Collateralization Buffer

A buffer absorbs oracle lag and one-sided flow:

```
wBTC: 1 BTC in circ -> 1.15 BTC in reserve (15% buffer)
wETH: 1 ETH in circ -> 1.12 ETH in reserve (12% buffer)
stables: 1 USDC -> 1.02 USDC in reserve (2% buffer)
```

### 3.3 Reserve Management

1. On-chain multi-sig (7-of-12 from validator set)
2. Gradual rebalancing - no large moves without announcement
3. Transparency - reserve addresses published + verified monthly (live dashboard)
4. Insurance fund - 5% of all fees go to a separate reserve

## 4. Revenue Model

| Source | Fee | Who Pays |
|--------|-----|----------|
| **Spread (every swap)** | ~0.30% round-trip (mid ± 0.15%) | Trader |
| Bridge-out fee (withdraw real asset) | 0.5% | Withdrawer |
| Bridge-in fee (deposit real asset) | 0.1% | Depositor |
| Remittance fee | 0.3% | Sender |
| Validator commission | Variable | Delegators |

### 4.1 Fee/Spread Distribution

```
Spread + swap fees:
  40% -> Reserve buffer (compounds backing)
  25% -> Validator / staker rewards
  20% -> Insurance fund
  15% -> Treasury

Bridge-out fee (0.5%):
  60% -> Reserve buffer
  20% -> Insurance fund
  20% -> Treasury
```

## 5. Price Setting Mechanics

### 5.1 The Quote

Each block, for each asset:

```
mid   = robust multi-source TWAP (the oracle)
buy   = mid * (1 - spread)     # protocol buys the asset from users
sell  = mid * (1 + spread)     # protocol sells the asset to users
```

`spread` is governance-set per asset (default 0.15% each side), and **widens
automatically** with volatility and inventory skew.

### 5.2 Spread Schedule (predictability without giving money away)

The protocol publishes its **spread curve and update cadence** in advance — users know
exactly how pricing behaves. What is *not* pre-committed is a future off-market price
(that would be writing free options to the market — see [REDESIGN.md](REDESIGN.md) §1.3).

### 5.3 Spread Adjustment Types

| Type | Description | Use Case |
|------|-------------|---------|
| Flat spread | Fixed bps each side | Calm markets |
| Volatility-scaled | Spread = base × volatility multiplier | Turbulent markets |
| Inventory-skewed | Quote tilts to rebalance inventory | One-sided flow |
| Emergency widen / freeze | Spread → large, or quotes paused 48h | Circuit breaker |

## 6. How Trading Affects the Reserve

### 6.1 User buys wBTC (protocol sells)

```
Market mid = $100,000 | Ma'at sell = $100,150
User pays 100,150 (in MAAT-equivalent) -> receives 1 wBTC
Protocol releases 1 wBTC against backing, books +$150 spread
```

### 6.2 User sells wBTC (protocol buys)

```
Market mid = $100,000 | Ma'at buy = $99,850
User delivers 1 wBTC -> receives 99,850
Protocol acquires asset $150 below mid, books +$150 spread
```

Effect: **both directions add to the reserve.** The spread is the revenue; the oracle
keeps it solvent.

## 7. Behavior Across the Cycle (the honest "Joseph")

| Phase | Market | Ma'at behavior |
|-------|--------|----------------|
| Bull | Prices rising | Oracle tracks up; spread earned on heavy two-sided volume |
| Bear | Prices crashing | Spread widens with volatility; reserve still earns, backing held |
| Calm | Sideways | Tight spread, high volume, steady compounding |

The grain-store discipline of Joseph applies to **reserves and insurance fund growth in
good times to survive shocks** — not to quoting off-market prices.

## 8. Key Metrics

| Metric | Target |
|--------|--------|
| Reserve Ratio | >= 110% (and growing) |
| Reserve trend | **Up** (net of all flows) |
| TVL | $100M+ at maturity |
| Daily Volume | $10M+ at maturity |
| Effective spread | 0.10–0.50% (vol-dependent) |
| MAAT MC | $50M+ at maturity |
| Daily Revenue | $50K+ at maturity |
