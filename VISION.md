# 𓁧 Ma'at — VISION.md

> **The Full Philosophy: Why MEV-free, predictable, reserve-backed execution is the next evolution of crypto**

> ⚠️ v0.2. The philosophy is unchanged — *order over chaos*. The mechanism is corrected:
> Ma'at earns an honest spread around a market-tracking oracle, instead of giving the
> reserve away at off-market prices. See [REDESIGN.md](REDESIGN.md).

---

## 1. The Problem: Crypto Markets Are Broken

### 1.1 Volatility is a Feature, Not a Bug — But It's Also a Tax

Every crypto user pays the "volatility tax":
- **Traders**: Lose money to MEV bots (sandwich attacks cost $1B+/year)
- **Liquidity Providers**: Impermanent loss on AMMs
- **Normal Users**: Can't predict execution quality block to block
- **Businesses**: Can't accept crypto payments without immediate conversion

The market says: "This is the price of decentralization."

We say: **"No — this is the price of poor microstructure design."**

### 1.2 The Oracle Problem (we harden it, we don't pretend it away)

Every DeFi protocol needs price data. Oracles can be:
- Expensive, slow, and an attack vector ($1.6B lost to oracle manipulation in 2023)

The original Ma'at draft tried to "be the oracle" and skip external prices. That is exactly
what made it insolvent — a price untethered from the market is a faucet. **Ma'at instead
builds a hardened, multi-source, time-weighted oracle and treats it as critical
infrastructure** (median of ≥3 venues, deviation guards, circuit breakers). We respect the
oracle problem and engineer around it.

### 1.3 The MEV Crisis

MEV is a $400M+/year tax on users (sandwiches, frontrunning, backrunning, liquidations).

In Ma'at, the clearing price is **fixed for the entire block**, so reordering transactions
cannot change anyone's price. **No intra-block price → no sandwich → no MEV.** The price is
the price for that block.

---

## 2. The Solution: One Fair Price Per Block, Honest Spread

### 2.1 The FX-Window Model for Crypto

A central bank's FX window and a professional market maker both do the same thing: quote a
**two-sided price around the market mid and keep a small spread.** Users get certainty and
fair execution; the institution stays solvent because it never sells below or buys above the
market.

**Ma'at is the first chain to make that the native execution model:**
- One clearing price per block (no MEV, no slippage)
- Price tracks the market via a robust oracle (always solvent)
- The protocol keeps a small spread (the reserve compounds)
- Every wrapped asset is backed 1:1+ by real reserves (transparent, audited live)

### 2.2 Why This Beats AMMs and Order Books

```
AMM:        price moves every trade -> slippage + MEV + impermanent loss
Order book: needs deep liquidity -> thin books = bad fills
Ma'at:      fixed-per-block quote around oracle mid -> no slippage, no MEV,
            no IL, and the spread accrues to the reserve instead of bleeding to LPs
```

---

## 3. The Philosophy: "سوق مظبوط"

In Egyptian Arabic: **"سوق مظبوط"** — a market that's *correct*, *fair*, *proper*.

ما معنى "سوق مظبوط"؟ يعني:
- السعر معروف وأنت بتشتري (ثابت داخل البلوك)
- مفيش حد يغشك (no MEV)
- مفيش frontrunning
- السعر بيتبع السوق (مش رقم من خيالنا)
- البروتوكول بياخد فرق بسيط عادل (spread) عشان يفضل صامد للأبد

### 3.1 Why This Is Egyptian

Egypt has a long history of **managed, disciplined economies**:
- Ancient Egypt: Pharaoh stored grain in years of plenty for years of famine (Joseph) —
  *reserve discipline*, the heart of Ma'at's insurance fund.
- Modern Egypt: a crawling-peg FX regime — predictable, managed, market-tracking.
- Egyptian character: "النظام أحسن من الفوضى."

Ma'at is the **digital extension of this discipline** — order and fairness, funded by an
honest spread, never by giving the treasury away.

---

## 4. The Comparison

| Feature | Free Market (AMM) | Ma'at (v0.2) |
|---------|-------------------|--------------|
| Price within a block | Moves every trade | Fixed |
| MEV | Endemic | Non-existent |
| Oracle | Often weak | Hardened, multi-source, critical |
| Slippage | Yes | Zero |
| Impermanent loss | Yes (LPs) | None |
| Who earns the spread | LPs (with IL) | Reserve (compounds) |
| Reserve under normal trade | n/a | **Grows** |
| Risk type | MEV + IL + market | Bridge + oracle + inventory |

---

## 5. The Thesis

**We believe that:**

1. Crypto's microstructure (per-trade pricing) causes MEV, slippage, and IL.
2. A chain with one fair, market-tracking price per block is more efficient and fairer.
3. The protocol should **earn** the spread, not give it away — that's what makes it eternal.
4. A hardened oracle plus real 1:1+ reserves plus circuit breakers can be robust.
5. The Egyptian discipline — order, fairness, reserve buffers — wins in a world tired of
   MEV and oracle attacks.

**If we are right:** Ma'at becomes the fair-execution and cross-chain settlement layer of crypto.

**If we are wrong:** we document exactly why, so the next person doesn't repeat it.

---

## 6. The Name: Ma'at 𓁧

**Ma'at** was the goddess of Truth, Balance, Order, Justice, Harmony. The feather of Ma'at
weighs the heart. We weigh every quote against the market with the same honesty: fair price,
fair spread, real backing.

---

<p align="center">
  <i>"مفيش حاجة اسمها سوق حر — في سوق مظبوط وسوق مكسور"</i><br>
  — <b>Ma'at Genesis</b>
</p>
