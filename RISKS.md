# RISKS.md - Every Risk, Honest Mitigation

---

## RISK #1: Bank Run (Reserve Depletion)

**What:** If too many people arb out at once, real asset reserve drains.

**Severity:** CRITICAL (existential)

**Mitigations:**
1. Over-collateralization (110-120%)
2. Daily withdrawal caps per asset
3. Gradual price changes (no sudden large gaps)
4. Circuit breaker: trading pauses if reserve < 100%
5. Insurance fund (5% of all fees)
6. Emergency MAAT minting as last resort (gov vote)

**Example:**
```
BTC in reserve: 100 BTC
wBTC circulating: 87 BTC (ratio = 115%)
Daily withdrawal cap: 10 BTC
Worst case: reserve drops to 90 BTC in 1 day
Circuit breaker at 87 BTC (100% ratio)
Market recovers, ratio restored
```

---

## RISK #2: Oracle Manipulation from Outside

**What:** Whales manipulate external market price to exploit spread.

**Severity:** HIGH

**Mitigations:**
1. TWAP pricing from multiple sources
2. Price change delay (announced N days before)
3. Circuit breaker if price deviates >50% in 1 hour
4. Emergency freeze (multi-sig)

---

## RISK #3: Bridge Attack

**What:** Cross-chain bridge gets hacked (#1 attack vector in crypto).

**Severity:** CRITICAL

**Mitigations:**
1. Battle-tested bridges (Wormhole, not custom)
2. Threshold-sig for BTC (no single point of failure)
3. Gradual bridge limits (max X per day)
4. Separate multi-sig for each bridge
5. Insurance fund for bridge losses
6. Validator slashing for bridge fraud

---

## RISK #4: Regulatory

**What:** Fixed prices look like market manipulation.

**Severity:** MEDIUM

**Mitigations:**
1. Full transparency: all prices, schedules published
2. Governance-driven, not admin-decided
3. Legal structure in crypto-friendly jurisdiction
4. No KYC, but trackable
5. Legal defense fund in treasury

---

## RISK #5: MEV on the Swap Path

**What:** Could MEV bots frontrun even with fixed prices?

**Severity:** LOW

**Mitigations:**
1. Fixed prices = zero slippage = no sandwich attack
2. Transaction ordering doesnt change price
3. Fair block building (Tendermint)

---

## RISK #6: Liquidity Fragmentation

**Severity:** MEDIUM

**Mitigations:**
1. Start with 1 asset (ETH) - prove the model
2. Add assets gradually
3. No AMM pools = no fragmentation
4. All liquidity is in the reserve

---

## RISK #7: The Model Just Doesnt Work

**What:** Nobody comes. No TVL. Dead chain.

**Severity:** MEDIUM

**Mitigations:**
1. Bootstrap with aggressive discounts (10% below market)
2. Simulate economics before mainnet
3. Airdrop to active DeFi users
4. If fails -> documented lesson

---

## RISK #8: Technical Bug

**Severity:** HIGH

**Mitigations:**
1. Extensive testing (unit + integration + fuzz)
2. Bug bounty program
3. Circuit breaker for price anomalies
4. Emergency upgrade via governance
5. Audits by 3+ independent firms before mainnet

---

## Risk Matrix

| Risk | Likelihood | Impact | Priority |
|------|-----------|--------|----------|
| Bank Run | Medium | Catastrophic | P0 |
| Bridge Attack | Low | Catastrophic | P0 |
| Oracle Manipulation | Medium | High | P1 |
| Regulatory | Medium | Medium | P1 |
| Technical Bug | Low | High | P1 |
| Model Fails | Medium | Medium | P2 |
| Liquidity Frag | Low | Medium | P2 |
| MEV | Very Low | Low | P3 |
