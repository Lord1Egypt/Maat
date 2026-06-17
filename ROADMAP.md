# ROADMAP.md - Phased Execution Plan

> From Concept to Dominance

---

## Phase 0: Concept & Design (NOW)
**Status: IN PROGRESS**

- [X] Define core economic model
- [X] Create repository with full documentation
- [X] Write VISION.md (philosophy)
- [X] Write ECONOMICS.md (tokenomics)
- [X] Write PLANNED_ECONOMY.md (price mechanics)
- [X] Write ARCHITECTURE.md (technical design)
- [X] Write RISKS.md (honest disclosure)
- [X] Write COMPETITION.md (market comparison)
- [ ] Build Python simulation model (arb engine simulation)
- [ ] Get community feedback on concept

---

## Phase 1: Simulation (Month 1)

**Goal:** Prove the economics work before writing any blockchain code.

- [ ] Build Python simulation of arb engine
- [ ] Simulate: 30 days of arb with various spreads
- [ ] Simulate: Bank run scenario (reserve depletion)
- [ ] Simulate: Bull market vs bear market behavior
- [ ] Simulate: MAAT price discovery model
- [ ] Write whitepaper with simulation results
- [ ] Publish simulation results for community review
- [ ] Iterate economic model based on simulation findings

**Deliverable:** whitepaper.md + simulation/ folder with reproducible notebooks

---

## Phase 2: Cosmos SDK Chain (Month 2)

**Goal:** Running blockchain with basic modules.

- [ ] Initialize Cosmos SDK project
- [ ] Implement x/maat module (native coin, staking, governance)
- [ ] Implement x/pegged module (mint/burn wrapped assets)
- [ ] Implement x/price module (fixed prices + schedule)
- [ ] Implement x/reserve module (track real reserves)
- [ ] Implement x/treasury module (fee collection)
- [ ] Write unit tests for all modules
- [ ] Deploy testnet with 4 validators
- [ ] Verify basic swap flow: MAAT <-> wETH

**Deliverable:** Testnet with swap functionality

---

## Phase 3: Ethereum Bridge (Month 2-3)

**Goal:** Real cross-chain flow with Ethereum.

- [ ] Deploy Wormhole contract on Ethereum testnet
- [ ] Integrate Wormhole VAAs into x/bridge module
- [ ] Test: Lock ETH on Sepolia -> mint wETH on Ma'at testnet
- [ ] Test: Burn wETH on Ma'at -> release ETH on Sepolia
- [ ] End-to-end arb test: Buy wETH cheap -> bridge out -> sell on Ethereum
- [ ] Security audit of bridge integration
- [ ] Deploy to Ethereum mainnet (limited, audited)

**Deliverable:** wETH flows cross-chain

---

## Phase 4: Bootstrap Launch (Month 3)

**Goal:** First real users and TVL.

- [ ] Launch Ma'at mainnet
- [ ] Bootstrap with ETH (10% below market price)
- [ ] 1-week announcement: "wETH at 10% discount starting Day 8"
- [ ] Marketing push: "Guaranteed arb on Day 8"
- [ ] Deploy MAAT token on Ethereum (ERC-20 representation)
- [ ] List MAAT on Uniswap for initial price discovery
- [ ] Airdrop to active DeFi users (Snapshot of top 10K addresses)
- [ ] Monitor first arb flows
- [ ] Verify reserve ratio stays healthy

**Metrics:**
- [ ] $1M+ TVL within 30 days
- [ ] >100 unique swappers per day
- [ ] Reserve ratio > 110% at all times
- [ ] No bank run events

---

## Phase 5: Asset Expansion (Month 4-6)

**Goal:** Multi-asset cross-chain exchange.

- [ ] Add wBTC (threshold-sig bridge)
- [ ] Add USDC + USDT (stablecoins at 1:1)
- [ ] Add wSOL (Solana via Wormhole)
- [ ] IBC integration for Cosmos ecosystem
- [ ] Add 10+ additional EVM assets
- [ ] Price tiering for whales (volume discounts)
- [ ] Gas optimization for high-frequency trading

**Metrics:**
- [ ] $10M+ TVL
- [ ] 5+ assets trading
- [ ] $1M+ daily volume
- [ ] Revenue: $3K+/day in fees

---

## Phase 6: Governance Maturity (Month 6-9)

**Goal:** DAO-driven, decentralized price setting.

- [ ] Full governance for price proposals
- [ ] Price oracle module (MAAT price feeds for DeFi)
- [ ] Staking rewards optimization
- [ ] Validator incentives aligned with reserve health
- [ ] Insurance fund growth to 5% of TVL
- [ ] Bug bounty program launch
- [ ] Security council rotation

**Metrics:**
- [ ] 50%+ MAAT staked
- [ ] 20+ active validators
- [ ] 3+ governance proposals passed

---

## Phase 7: Dominance (Month 9-12+)

**Goal:** Ma'at becomes the primary cross-chain exchange.

- [ ] Smart contract layer (CosmWasm) for advanced DeFi
- [ ] Leverage trading on fixed prices
- [ ] Institutional-grade KYC/AML option
- [ ] Mobile wallet
- [ ] Fiat on-ramp integration
- [ ] 50+ supported assets
- [ ] Referral program
- [ ] Partnership with major DeFi protocols

**Metrics:**
- [ ] $100M+ TVL
- [ ] $10M+ daily volume
- [ ] $50K+/day revenue
- [ ] Top 100 crypto project by market cap
- [ ] Referenced as oracle by other protocols

---

## Phase 8: What Comes After

- Layer 2 for Ma'at (faster, cheaper swaps)
- Cross-chain lending/borrowing with fixed prices
- NFT exchange with fixed prices
- Real-world asset tokenization at protocol-set prices

---

## Current Focus

```
Phase 0 > Phase 1 > Phase 2 > Phase 3 > Phase 4 > Phase 5 > Phase 6 > Phase 7
   [##.......]  [........]  [........]  [........]  [........]  [........]  [........]
   ^ YOU ARE HERE
```

Next concrete step: **Build the Python simulation model** to prove the economics.
