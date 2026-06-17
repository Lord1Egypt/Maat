# PLANNED_ECONOMY.md

> How Price Setting Works - The Detailed Mechanics

---

## 1. The Core Mechanism

### 1.1 The Protocol IS the Oracle

Ma'at does not use external oracles (Chainlink, etc.). Every block stores:

```
state.asset_prices["wBTC"] = {
    price: 90000,
    effective_block: 123456,
    announced_block: 123400,
    next_price: 95000,
    next_effective_block: 130000,
}
```

### 1.2 Mint/Burn at Protocol Price

```
User sends 1 wBTC:
  1. Read fixed price ($90,000)
  2. Burn wBTC
  3. Credit user 90,000 MAAT

User sends 90,000 MAAT:
  1. Read fixed price ($90,000)
  2. Check reserve has real BTC
  3. Debit 90,000 MAAT
  4. Mint 1 wBTC
```

No slippage. No spread. No MEV.

## 2. The Announcement Cycle

```
Default cycle: 30 days

Day 0:   Price = $90,000 (current)
Day 1:   Announce NEXT price = $95,000 (effective Day 30)
Day 1-29: Market adjusts, users plan, arb flows
Day 30:  Price changes to $95,000
Day 31:  Announce next price...
```

## 3. How Users Exploit the Schedule

### 3.1 Next Price HIGHER than Current

```
Current = $90,000, Next = $95,000 (in 29 days)
- Buy wBTC NOW at $90,000
- Hold 29 days
- Price becomes $95,000
- Sell wBTC for $95,000 MAAT
- Profit = $5,000/BTC
```

### 3.2 Next Price LOWER than Current

```
Current = $95,000, Next = $90,000 (in 29 days)
- SELL wBTC NOW at $95,000
- Get 95,000 MAAT
- Wait for drop
- Buy back wBTC at $90,000
- Profit = $5,000/BTC
```

## 4. Price Change Types

| Type | Description | Use Case |
|------|-------------|---------|
| Fixed Reprice | Move to new fixed value | Normal adjustments |
| Percentage | Move by X% | Following market trends |
| Market-Relative | Peg: Market - 2% | Safety mode, constant arb |
| Staircase | Volume-tiered pricing | Whale discounts |
| Emergency Freeze | Pause all trading 48h | Circuit breaker |

## 5. Price Change Governance

### 5.1 Who Can Propose

Any MAAT staker with >= 100,000 MAAT delegated.

### 5.2 Proposal Format

```json
{
  "title": "wBTC: $90K -> $95K",
  "asset": "wBTC",
  "current_price": 90000,
  "new_price": 95000,
  "effective_delay_blocks": 300000,
  "rationale": "Market recovering, BTC up 15%"
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
- Pause all trading (circuit breaker)
- Force price freeze for up to 48 hours
- Trigger emergency reserve rebalancing

Powers are audited, time-limited, and publicly logged.
