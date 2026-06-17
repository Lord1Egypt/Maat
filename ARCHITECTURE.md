# ARCHITECTURE.md - Technical Design (v0.2)

> Cosmos SDK Modules, Oracle, Bridges, Consensus, and Smart Contracts

> ⚠️ Corrected from v0.1: `x/price` no longer sets an arbitrary fixed price, and there is
> now a required `x/oracle` module. Pricing = oracle mid ± spread, fixed per block.
> See [REDESIGN.md](REDESIGN.md).

---

## 1. Tech Stack

| Component | Technology |
|-----------|-----------|
| Blockchain Framework | Cosmos SDK v0.50+ |
| Consensus | CometBFT (Tendermint) v0.38+ |
| Smart Contracts | CosmWasm (Rust) |
| Cross-Chain (EVM) | Wormhole |
| Cross-Chain (BTC) | Threshold ECDSA (custom, audited) |
| Oracle | Multi-source TWAP (median of >=3 venues) |
| Storage | IAVL+ (validators) |
| RPC | gRPC + REST + WebSocket |
| Explorers | Mintscan, Ping.pub |
| Wallets | Keplr, Leap, Cosmostation |

## 2. High-Level Architecture

```
+-------------------------------------------------------+
|                   MA'AT CHAIN                         |
|                    (Cosmos SDK)                       |
|                                                       |
|  +-----------------------------------------------+  |
|  |     Consensus Layer (CometBFT)                 |  |
|  |     50-100 Validators, ~2s block time          |  |
|  +-----------------------------------------------+  |
|                     |                                |
|  +-----------------------------------------------+  |
|  |     Application Layer                          |  |
|  |  x/maat (core) | x/pegged (wraps) | x/bridge |  |
|  |  x/oracle (mid) | x/market (quote) | x/reserve|  |
|  +-----------------------------------------------+  |
|                     |                                |
|  +-----------------------------------------------+  |
|  |     Cross-Chain Layer                          |  |
|  |  Ethereum | Bitcoin | Solana | Cosmos (IBC)   |  |
|  +-----------------------------------------------+  |
+-------------------------------------------------------+
```

## 3. Core Modules

### 3.1 x/maat - Core Module
Handles native coin, validators, governance.
- Native denom: umaat (micro MAAT)
- Supply: 1,000,000,000,000,000 umaat (1B MAAT)
- Validator set: 50-100
- Inflation: 0-10% annual (adjustable via gov)

### 3.2 x/pegged - Wrapped Assets Module
Manages minting/burning of wrapped tokens against real backing.
- Each asset: denom, name, decimals, total supply
- Reserve address (multi-sig holding real assets)
- Reserve ratio: current vs required (>= 110%)

**Messages:** MsgMintWrapped / MsgBurnWrapped

### 3.3 x/oracle - Price Source Module (REQUIRED)
The hardened market-price feed. **This is what keeps the protocol solvent.**

**State per asset:**
- Sources (>=3 independent venues)
- Median + time-weighted mid (TWAP)
- Deviation guard (drop sources off by > X%)
- Circuit breaker flag (freeze if mid moves > 50% in 1h)

### 3.4 x/market - Quote & Settlement Module (THE INNOVATION)
Computes one clearing quote per block and settles swaps at it.

**Per asset per block:**
- mid = x/oracle TWAP
- spread_bps = base × vol_multiplier × inventory_skew
- buy = mid × (1 − spread); sell = mid × (1 + spread)

**Swap logic (no AMM, fixed-per-block):**
```
func (k Keeper) Swap(fromDenom, toDenom string, amount Int) (Int, error) {
    sellPx := k.SellPrice(fromDenom)   // protocol buys 'from'
    buyPx  := k.BuyPrice(toDenom)      // protocol sells 'to'
    value  := amount.Mul(sellPx)
    out    := value.Quo(buyPx)
    k.BurnCoins(fromDenom, amount)
    k.MintCoins(toDenom, out)          // backing-checked
    k.AccrueSpread(...)                // spread -> reserve/insurance
    return out, nil
}
```

Price fixed per block → no slippage, no MEV. Spread → reserve grows.

### 3.5 x/bridge - Cross-Chain Bridge Module
**Phased Bridge Strategy (battle-tested first, never roll-your-own bridge):**
1. Phase 1: Wormhole (EVM chains)
2. Phase 2: Native BTC (threshold-sig, audited)
3. Phase 3: IBC (Cosmos ecosystem)
4. Phase 4: Solana (Wormhole/SVM)

Controls: daily caps, large-tx delay/cancel window, per-bridge multisig, insurance fund.

### 3.6 x/reserve - Reserve Module
- ReserveProof: on-chain attestation of real assets
- ReserveAlert: health status (healthy/warning/critical)
- Auto-pause quotes if backing drops below 100%

## 4. Consensus & Security

| Parameter | Value |
|-----------|-------|
| Consensus | Tendermint BFT |
| Validators | 50-100 |
| Block time | ~2 seconds |
| Finality | Instant (1 block) |
| Safety | BFT (< 1/3 malicious) |

**Security:**
- Slashing (5-100%) for double-sign/duplicate-bridge
- Jailing for >24h downtime
- Circuit breaker if reserve < 100% or oracle deviates
- Daily bridge-out caps

## 5. Smart Contract Layer
CosmWasm (Rust) for advanced products: the reserve index, bonds backed by real spread
revenue, and remittance routing.
