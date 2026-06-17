# ECONOMICS.md

> The Ma'at Economic Model

Tokenomics, Reserve Model, Supply Mechanics, and Revenue

---

## 1. Token Overview

| Token | Type | Supply | Price Model |
|-------|------|--------|-------------|
| MAAT | Native coin (gas + gov) | Fixed: 1B | Free market |
| wBTC/wETH/wUSDC | Wrapped assets | Demand-driven | Fixed protocol price |

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
- **Governance**: Vote on price changes, reserve ratios, bridge params
- **Staking**: Validators stake MAAT for consensus
- **Arb Fee Discount**: Stakers get reduced swap/bridge fees
- **Reserve Backstop**: In emergencies, MAAT backs wrapped assets (gov vote)

## 3. The Reserve Model

### 3.1 Full-Backed Reserve (1:1+)

Every wrapped asset is backed 1:1 by real assets in multi-sig:

```
wBTC minted = 100 --> 100 BTC held in reserve
wETH minted = 500 --> 500 ETH held in reserve
```

### 3.2 Over-Collateralization Buffer

10-20% over-collateralization for volatile assets:

```
wBTC: 1 BTC in circ -> 1.15 BTC in reserve (15% buffer)
wETH: 1 ETH in circ -> 1.12 ETH in reserve (12% buffer)
stables: 1 USDC -> 1.02 USDC in reserve (2% buffer)
```

### 3.3 Reserve Management

1. On-chain multi-sig (7-of-12 from validator set)
2. Gradual rebalancing - no large moves without announcement
3. Transparency - reserve addresses published + verified monthly
4. Insurance fund - 5% of all fees go to separate reserve

## 4. Revenue Model

| Source | Fee | Who Pays |
|--------|-----|----------|
| Swap fee (asset <-> MAAT) | 0.3% | Swapper |
| Bridge-out fee (withdraw real asset) | 0.5% | Withdrawer |
| Bridge-in fee (deposit real asset) | 0.1% | Depositor |
| Validator commission | Variable | Delegators |
| Treasury staking yield | From inflation | MAAT holders |

### 4.1 Fee Distribution

```
Swap fee (0.3%):
  40% -> Liquidity reserve
  30% -> Validator rewards
  20% -> Insurance fund
  10% -> Treasury

Bridge-out fee (0.5%):
  60% -> Reserve buffer
  20% -> Insurance fund
  20% -> Treasury
```

## 5. Price Setting Mechanics

### 5.1 Initial Price

On launch, first prices are set:
- wETH = Market ETH * 0.98 (2% below market)

### 5.2 Price Schedule

Price changes follow a fixed schedule:
1. Governance proposes new price schedule
2. 7-day discussion period
3. 7-day voting period
4. If passed -> price changes at announced block height

### 5.3 Price Change Types

| Type | Description | Example |
|------|-------------|---------|
| Flat reprice | Move to new fixed price | wBTC $90K -> $95K |
| Percentage shift | Move by X% | wBTC +5% |
| Market-relative | Peg to market + offset | wBTC = Market - 2% |
| Tiered | Different prices per volume | <1 BTC: $90K, >=1 BTC: $88K |

## 6. The Arb Engine

### 6.1 When Ma'at Price > Market Price

```
Market BTC = $100,000 | Ma'at wBTC = $110,000 (premium)
  1. User buys 1 BTC on Binance ($100,000)
  2. Bridges to Ma'at
  3. Sells wBTC for 110,000 MAAT
  4. Sells MAAT on market
  Profit = ~$9,000
```
Effect: Real BTC flows INTO Ma'at. Reserve grows.

### 6.2 When Ma'at Price < Market Price

```
Market BTC = $100,000 | Ma'at wBTC = $90,000 (discount)
  1. User buys MAAT on market
  2. Sends 90,000 MAAT to Ma'at chain
  3. Buys 1 wBTC at $90,000
  4. Bridges out real BTC, sells at $100,000
  Profit = ~$9,000
```
Effect: Real BTC flows OUT. Reserve shrinks. MAAT burned.

## 7. The Joseph Strategy (7 Plenty, 7 Famine)

| Phase | Market | Ma'at Strategy |
|-------|--------|----------------|
| Bull | Prices rising | Ma'at rises SLOWER -> cheap buys -> MAAT demand soars |
| Bear | Prices crashing | Ma'at falls SLOWER -> sell higher -> reserve accumulates |

Result: In bull, Ma'at is best to BUY. In bear, best to SELL. Always the best venue.

## 8. Key Metrics

| Metric | Target |
|--------|--------|
| Reserve Ratio | >= 110% |
| TVL | $100M+ at maturity |
| Daily Volume | $10M+ at maturity |
| Arb Spread | 1-10% (intentional) |
| MAAT MC | $50M+ at maturity |
| Daily Revenue | $50K+ at maturity |
