# RISKS.md - Every Risk, Honest Mitigation (v0.2)

> Updated for the spread-capture model. The biggest risk in v0.1 was the model itself —
> now Risk #0, fixed by design. See [REDESIGN.md](REDESIGN.md).

---

## RISK #0: Structural Reserve Drain (the v0.1 killer — now FIXED)

**What:** v0.1 quoted a fixed *off-market* price and called the gap "guaranteed user arb."
That gap was paid out of the reserve on every trade — the reserve drained by construction,
exactly like Terra/UST. This happened in **normal operation**, not a panic.

**Severity:** WAS existential.

**Fix:** Quote **oracle mid ± spread, fixed per block.** The protocol now *earns* the
spread instead of giving it away; the reserve grows under normal trading. "Guaranteed user
profit" is replaced by "guaranteed best execution + no MEV." (Proof: REDESIGN.md §1–2.)

---

## RISK #1: Bridge Attack

**What:** Cross-chain bridge gets hacked (#1 attack vector in crypto, by far).

**Severity:** CRITICAL

**Mitigations:**
1. Battle-tested bridges (Wormhole), never custom for EVM
2. Threshold-sig for BTC (no single point of failure), audited
3. Gradual bridge limits (max X per day)
4. Separate multi-sig per bridge
5. Large-tx delay + cancel window
6. Insurance fund for bridge losses
7. Validator slashing for bridge fraud

---

## RISK #2: Oracle Manipulation

**What:** Now that pricing tracks the market, a manipulated feed could mis-quote.

**Severity:** HIGH

**Mitigations:**
1. Median of >=3 independent venues (CEX + on-chain TWAP)
2. Time-weighting to resist flash spikes
3. Deviation guard: drop any source off by > X%
4. Circuit breaker: freeze quotes if mid moves > 50% in 1h
5. Over-collateral buffer absorbs short oracle lag
6. Spread auto-widens with volatility

---

## RISK #3: Regulatory / Securities

**What:** Marketing "guaranteed risk-free profit" + yield-bearing bonds = a security in
US/EU. Fixed pricing could also be framed as manipulation.

**Severity:** HIGH

**Mitigations:**
1. **Drop all "guaranteed/risk-free profit" language** — sell execution quality, not returns
2. Full transparency: spreads, oracle sources, reserves published
3. Governance-driven parameters, not admin-decided
4. Licensed VASP/fintech partners for fiat ramps (remittance)
5. Legal structure in a crypto-friendly jurisdiction; legal defense fund

---

## RISK #4: Inventory Imbalance

**What:** Sustained one-sided flow leaves the protocol long or short an asset at risk of
adverse price moves.

**Severity:** MEDIUM

**Mitigations:**
1. Inventory-skewed quotes pull inventory back to target
2. Spread widens as skew grows (deters toxic flow)
3. Periodic rebalancing across reserve assets
4. Per-block, per-asset exposure caps

---

## RISK #5: MEV on the Swap Path

**Severity:** LOW

**Mitigations:**
1. Price fixed per block = zero slippage = no sandwich
2. Transaction ordering cannot change the clearing price
3. Fair block building (Tendermint)

---

## RISK #6: Liquidity / Cold Start

**Severity:** MEDIUM

**Mitigations:**
1. Seed reserve inventory at launch (own funds + bonds backed by real spread revenue)
2. Start with 1 asset (ETH), prove the spread engine, then expand
3. LPs/bondholders earn real yield from spread (not inflationary emissions)
4. No AMM pools = no fragmentation

---

## RISK #7: The Model Just Doesn't Get Traction

**What:** Sound economics, but nobody trades.

**Severity:** MEDIUM

**Mitigations:**
1. Lead with the remittance corridor (real $30B demand) for non-speculative volume
2. Compete on execution quality + transparency vs AMMs
3. Airdrop to active DeFi users; integrations with wallets
4. If it fails -> documented lesson, reserves return to users (fully backed)

---

## RISK #8: Technical Bug

**Severity:** HIGH

**Mitigations:**
1. Extensive testing (unit + integration + fuzz), sim-driven solvency tests
2. Bug bounty program
3. Circuit breakers for price/oracle anomalies
4. Emergency upgrade via governance
5. Audits by 3+ independent firms before mainnet

---

## Risk Matrix

| Risk | Likelihood | Impact | Priority |
|------|-----------|--------|----------|
| Structural Drain (#0) | **Eliminated by design** | (was Catastrophic) | Resolved |
| Bridge Attack | Low | Catastrophic | P0 |
| Oracle Manipulation | Medium | High | P0 |
| Regulatory/Securities | Medium | High | P1 |
| Technical Bug | Low | High | P1 |
| Inventory Imbalance | Medium | Medium | P1 |
| No Traction | Medium | Medium | P2 |
| Liquidity Cold Start | Medium | Medium | P2 |
| MEV | Very Low | Low | P3 |
