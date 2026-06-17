# REDESIGN.md — Making Ma'at Actually Work

> An honest economic and security re-architecture of Ma'at.
> Read this **before** ECONOMICS.md / PLANNED_ECONOMY.md / STRATEGY.md — it corrects them.

---

## 0. The One-Sentence Problem

The current design pays the arbitrage gap **out of its own reserve**, which means
"guaranteed risk-free profit for users" and "fully-backed reserve" are claims on the
**same coins** — so the reserve drains to zero by design, exactly like Terra/UST.

## 1. Proof (cold accounting)

```
Market BTC = $100,000 | Ma'at wBTC = $90,000 (fixed) | MAAT = $1

User: buy 90,000 MAAT ($90k) → swap for 1 wBTC → bridge out 1 REAL BTC → sell $100k
User profit  = +$10,000
Reserve      = −1 real BTC ($100,000 hard asset) + 90,000 MAAT (own, reflexive token)
```

Arbitrage is zero-sum minus fees. The user's "guaranteed" profit is **literally** the
reserve's loss. Over-collateralization (115%) only changes how many days the drain takes,
not its direction. This is the LUNA↔UST reflexive loop.

### 1.1 The "Joseph Strategy" loses both ways
- Bull (Ma'at below market): users drain real assets cheap → reserve loses value.
- Bear (Ma'at above market): users dump assets on you above market → MAAT diluted.
Whoever quotes the off-market price is the counterparty to every free trade. The house
always loses.

### 1.2 "We ARE the oracle" is the root error
- Quote a price unrelated to market → you bleed.
- Quote a price near market → you must KNOW the market → you need an oracle anyway.
There is no third option. Drop this claim entirely.

### 1.3 Announced future prices = free options
Pre-announcing the next price writes a free call/put (struck at the old price) to the
entire market. Free optionality = guaranteed negative EV for the protocol. Remove
"exploit the schedule" as a user feature.

---

## 2. The Fix: flip who pockets the spread

Quote **bid/ask around the market mid and keep the spread** (like every FX bureau,
market maker, and central bank). Fix the price *within each block/batch* (kills MEV and
slippage) but *track the market* via a robust TWAP/oracle (reserve never bleeds).

```
Ma'at quote:   BUY wBTC @ mid − 0.15%      SELL wBTC @ mid + 0.15%
- Fixed per block/batch  → no MEV, no slippage, no sandwich
- Tracks market          → reserve never drains
- Protocol earns ~0.30%  → reserve GROWS, compounding
```

Every promised property is now TRUE and sustainable:
zero slippage ✅ · no MEV ✅ · predictable ✅ · no impermanent loss ✅ · transparent ✅ ·
**reserve compounds instead of draining.**

New honest pitch (replaces "guaranteed arb"):
> **Guaranteed best execution. Guaranteed no MEV. Guaranteed transparent backing.**

---

## 3. The three sustainable products (map to old phases)

| Product | What it is | Revenue | Replaces |
|---|---|---|---|
| **MEV-free batch settlement** | One clearing price per block, oracle-mid ± spread | Spread/fee | Weapon #1 (fixed) |
| **Crawling-peg remittance rail** | Cross-chain stable transfer; Egypt runs a real crawling peg | 0.3% fee | Weapon #6 (strongest idea) |
| **Scheduled-rebalance reserve index** | On-chain ETF, transparent NAV, announced rebalances | Mgmt fee | Phase 5 "world reserve" |

Pharaoh Bonds (old Weapon #5) become legitimate **only after** this fix — now there is
real spread/fee revenue to back the 15–20% payout. Under the old model they are a Ponzi
layer on top of a draining reserve.

---

## 4. Liquidity / cold-start (corrected)

- Old model: "liquidity = reserve," and arb *drains* it. Unsolvable.
- New model: LPs and bondholders earn **real yield from spread revenue** (not inflationary
  emissions). Bootstrapping is a normal market-maker inventory problem.
- Seed inventory + pay LPs the spread. No "give assets away at a discount" bootstrap.

---

## 5. Risk register (corrected priorities)

| # | Risk | Sev | Note |
|---|---|---|---|
| **0** | **Structural reserve drain by design** | **Existential** | NEW — the real killer; happens in normal operation, not a panic. Fixed by §2. |
| 1 | Bridge hack | Existential | Correctly the #1 real-world attack vector. Use battle-tested bridges, daily caps, insurance fund. |
| 2 | Oracle manipulation | High | Now unavoidable (you track market). Multi-source TWAP, deviation circuit breakers, quote freeze. |
| 3 | Securities-law exposure | High | Marketing "guaranteed risk-free profit" + yield bonds = a security in US/EU. Drop that framing regardless of mechanism. |
| 4 | Smart-contract bug | High | 3+ audits, fuzzing, bug bounty, upgrade governance. |
| 5 | Inventory imbalance | Medium | One-sided flow leaves you long/short an asset; widen spread dynamically + rebalance. |

---

## 6. Honest probability

- Current design surviving: **~0%** (drains by construction).
- Corrected design avoiding a Terra-style collapse: **high** — this is the "no risk" worth having.
- Corrected design becoming **#1 in crypto**: not guaranteeable by anyone; that's
  execution + timing + luck on a sound base. Sound base is necessary, not sufficient.

---

## 7. What to change in the existing docs

- **README / VISION:** replace "guaranteed user arb" with "guaranteed best execution / no MEV."
- **ECONOMICS §6, §7:** delete the arb-engine-as-profit and Joseph-wins-both-ways claims; replace with spread-capture model.
- **PLANNED_ECONOMY §1.1, §3:** drop "no oracle" and "exploit the schedule."
- **ARCHITECTURE:** x/price becomes oracle-mid ± spread per batch, not arbitrary fixed price.
- **STRATEGY:** keep brand/marketing engine; rebase value prop on execution quality + remittance, not free money.
- **RISKS:** add Risk #0 (above) as P0.
