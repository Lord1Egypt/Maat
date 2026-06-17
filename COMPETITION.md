# COMPETITION.md - Market Landscape (v0.2)

> Where Ma'at sits, honestly. It competes on **fair execution + spread efficiency +
> transparent reserves + cross-chain reach** — not on "free money."

> ⚠️ v0.1 claimed "no oracle needed" and "zero competition." Both were wrong. Ma'at uses a
> hardened oracle, and it competes directly with RFQ/MEV-resistant venues. See [REDESIGN.md](REDESIGN.md).

---

## The Category

Ma'at is an **MEV-free, oracle-priced, reserve-backed cross-chain exchange.** Its closest
analogues are RFQ/market-maker venues and batch-auction DEXs — but on its own L1 with a
single fixed-per-block clearing price and a reserve that earns the spread.

---

## Direct Comparisons

### 1. Uniswap / Curve (AMMs)

| Feature | Uniswap | Ma'at |
|---------|---------|-------|
| Pricing | AMM curve (moves every trade) | Oracle mid ± spread, fixed per block |
| Slippage | Yes, grows with size | Zero |
| MEV | Severe (sandwiches) | Non-existent |
| Liquidity source | LP deposits (suffer IL) | Protocol reserve (no IL) |
| Who earns the spread | LPs (net of IL) | Reserve (compounds) |

**Edge:** no slippage, no MEV, no IL.

### 2. CoW Protocol / Hashflow / 0x RFQ (the real competitors)

| Feature | CoW / RFQ | Ma'at |
|---------|-----------|-------|
| MEV protection | Yes (batch/RFQ) | Yes (fixed per block) |
| Pricing | Solver/MM quotes | Native oracle ± spread |
| Settlement | On host chain | Native L1 + reserve |
| Cross-chain | Limited | Native multi-chain |
| Reserve transparency | n/a | Live, on-chain |

**Edge:** native cross-chain settlement + transparent reserve in one place. **Honest
weakness:** these are credible, funded competitors — execution and liquidity decide.

### 3. Binance / Coinbase (CeFi)

| Feature | CeFi | Ma'at |
|---------|------|-------|
| Pricing | Order book | Oracle ± spread |
| Custody | Centralized | Multi-sig + validators |
| Permissionless | KYC | No KYC (ramps via licensed partners) |
| Transparent | No | Fully on-chain |

**Edge:** transparency + permissionless. **Weakness:** liquidity depth.

### 4. Terra (LUNA/UST) — the cautionary tale we DON'T repeat

| Feature | Terra | Ma'at v0.2 |
|---------|-------|-----------|
| Backing | Algorithmic (own token) | Real assets 1:1+ |
| Price | Defended an off-market peg | Tracks market (oracle ± spread) |
| Drain vector | Arb against the peg | None — protocol earns the spread |

> Ma'at **v0.1 shared Terra's flaw** (off-market price drained by arb). v0.2 removes it.

### 5. Thorchain

| Feature | Thorchain | Ma'at |
|---------|-----------|-------|
| Model | Cross-chain AMM | Fixed-per-block oracle exchange |
| Slippage | Yes | Zero |
| MEV | Possible | Impossible |
| Liquidity | LP pools (IL) | Reserve (no IL) |

---

## Why This Is Hard (honest)

| Challenge | Reality | Ma'at's Answer |
|-----------|---------|----------------|
| Bridge risk | #1 hack vector | Battle-tested bridges, caps, insurance |
| Oracle robustness | Mandatory | Multi-source TWAP + breakers |
| Liquidity bootstrap | Cold start | Seed inventory + remittance demand |
| Funded competitors | Real | Compete on cross-chain + transparency + spread |

---

## Market Opportunity

| Metric | Current Market | Ma'at Target |
|--------|---------------|--------------|
| DEX monthly volume | ~$100B | Capture a slice on execution quality |
| Cross-chain bridge TVL | ~$20B | Reserve-backed settlement share |
| Egypt remittances | ~$30B/yr | Non-speculative, recurring volume |

**Key insight:** Ma'at wins by being the **best place for predictable, MEV-free cross-chain
settlement** — not by promising free profit.

---

## Competitive Moat

1. **Microstructure**: fixed-per-block clearing price (no MEV, no slippage, no IL)
2. **Reserve that compounds**: spread accrues to backing, not to LPs with IL
3. **Transparency**: live, audited reserves vs opaque CeFi
4. **Cross-chain native**: bridge + exchange + settlement in one
5. **Remittance beachhead**: real demand from the Egyptian corridor
6. **Brand**: Egyptian order-and-discipline philosophy
