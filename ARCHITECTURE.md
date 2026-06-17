# ARCHITECTURE.md - Technical Design

> Cosmos SDK Modules, Bridges, Consensus, and Smart Contracts

---

## 1. Tech Stack

| Component | Technology |
|-----------|-----------|
| Blockchain Framework | Cosmos SDK v0.50+ |
| Consensus | CometBFT (Tendermint) v0.38+ |
| Smart Contracts | CosmWasm (Rust) |
| Cross-Chain (EVM) | Wormhole |
| Cross-Chain (BTC) | Threshold ECDSA (custom) |
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
|  |  x/price (prices) | x/reserve | x/treasury   |  |
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
Manages minting/burning of wrapped tokens.
- Each asset: denom, name, decimals, total supply
- Reserve address (multi-sig holding real assets)
- Reserve ratio: current vs required

**Messages:** MsgMintWrapped / MsgBurnWrapped

### 3.3 x/price - Fixed Price Module (THE INNOVATION)
Manages scheduled price changes.

**State per asset:**
- CurrentPrice + CurrentPriceBlock
- NextPrice + NextPriceBlock
- LastAnnouncedBlock
- PriceHistory (last 100 snapshots)

**Swap logic (no AMM):**
```
func (k Keeper) Swap(fromDenom, toDenom string, amount Int) error {
    priceFrom = k.GetCurrentPrice(fromDenom)
    priceTo = k.GetCurrentPrice(toDenom)
    fromValue = amount * priceFrom
    toAmount = fromValue / priceTo
    k.BurnCoins(fromDenom, amount)
    k.MintCoins(toDenom, toAmount)
    return nil
}
```

No AMM. No slippage. No MEV.

### 3.4 x/bridge - Cross-Chain Bridge Module
**Phased Bridge Strategy:**
1. Phase 1: Wormhole (EVM chains)
2. Phase 2: Native BTC (threshold-sig)
3. Phase 3: IBC (Cosmos ecosystem)
4. Phase 4: Solana (Wormhole/SVM)

### 3.5 x/reserve - Reserve Module
- ReserveProof: on-chain attestation of real assets
- ReserveAlert: health status (healthy/warning/critical)
- Auto-pause if reserve drops below 100%

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
- Circuit breaker if reserve < 100%
- Daily bridge-out caps

## 5. Smart Contract Layer
CosmWasm (Rust smart contracts) for complex DeFi and automation.
