# 🧠 Brainstorm Swarm Report: Ma'at Codebase Security & Stability Audit
*A synthesis of specialized agent perspectives to establish a bulletproof execution plan.*

> [!IMPORTANT]
> **Executive Summary:** The brainstorm swarm has completed a comprehensive audit of the Ma'at codebase. While the scaffolding is functionally complete, we identified critical consensus violations, memory-based persistence flaws, swap rate math omissions, and data races. This report outlines these vulnerabilities and establishes a Master Plan to transition the codebase to a production-ready, fully persisted, and secure state machine.

---

## ⚔️ The Swarm Debate
Below is a matrix summarizing the core findings, tradeoffs, and consensus reached by the subagents.

| Specialist Role | Core Perspective / Stance | Key Recommendation / Constraint |
| :--- | :--- | :--- |
| **Principal Architect** | Standardize Cosmos SDK persistence layers and Context propagation | Eliminate all in-memory keeper variables; serialize all state directly to the KVStore using proto. |
| **CSO / Security Auditor** | Verify input bounds, prevent signature forgery, and enforce value transfers | Add cryptographic quorum validation on VAAs, replay prevention, and concrete bank transfer logic in `x/market/Swap`. |
| **QA Lead / Devil's Advocate** | Stress failure paths, halting conditions, and edge-case execution | Resolve non-deterministic map iteration in bridge lookup queries to prevent state-divergence halts. |
| **Performance Engineer** | Minimize lock contention and avoid pointer escapes | Encapsulate structs or implement internal Mutexes to avoid data races between consensus and query threads. |
| **Economics & Compliance** | Enforce USD-value backing and protect against economic drains | Normalize decimal calculations for 18-decimal assets like wETH, and implement global NAV-based backing checks. |

---

## 🔍 Detailed Perspectives

### 📐 Principal Architect
- **In-Memory State Non-Determinism**: Storing active states (`oracles`, `reserves`, `limiters`, etc.) in Go variables instead of Cosmos SDK's KVStore breaks node restarts, State Sync, and query endpoints. Map iteration order in Go is randomized, causing state divergence.
- **Pass-By-Value State Loss**: Passing `AppModule` keepers by value copies their fields, causing updates to pointer methods (like `EndBlockCheck`) to occur on transient copies and throw away state mutations.
- **API Context Propagation**: All keeper methods must accept `sdk.Context` (or `context.Context`) to read and write from the multi-store database.

### 🛡️ Security Auditor
- **Cryptographic Bridge Verification**: `BridgeIn` unmarshals the Wormhole VAA but lacks cryptographic signature check against a validator quorum, enabling anyone to forge deposits.
- **VAA Replay Attack**: There is no lookup table of processed VAA hashes, allowing a single valid VAA to be replayed to mint tokens indefinitely.
- **Swap Balance Checks**: The `Swap` message handler mutates reserve states but does not perform actual bank token burns/mints, leaving user balances unchanged.
- **Governance Bypass**: `ResumeOracle` ignores authority checking, allowing unauthorized actors to resume halted price feeds.

### 😈 Devil's Advocate
- **Denom Lookup Collisions**: `AllDenoms()` uses randomized map iteration. If multiple denoms have withdrawals with the same sequential ID, `ExecuteWithdrawal` will non-deterministically match different denoms across validators, causing an App Hash halt.
- **Stale Price Feed Lock**: Stale oracle prices (with zero votes) continue to be served without stale age limits, creating pricing exploits.

### ⚡ Performance Engineer
- **Pointer Escape Data Races**: Keepers release mutexes but return direct pointers to `*pricing.Reserve` and `*bridge.Limiter` which are mutated concurrently by message servers, creating data races against gRPC query routines.
- **Memory Leak**: Bounded bridge withdrawal maps are never pruned, causing memory leaks over time.

### 📈 Economics & Compliance Analyst
- **Decimal Mismatch**: The spread calculation uses `MicroUSD` scaling directly, causing scale errors for 18-decimal wETH that inflate revenues by $10^{12}$ and drain the reserves.
- **Asset Backing Deadlock**: Enforcing physical backing per-asset blocks swap minting on assets with low physical backing, even if the protocol is globally over-collateralized.

---

## 🗺️ The Master Plan (Incremental & Verifiable)

### Phase 1: Stateless Keepers & KVStore Persistence
- [ ] Migrate `reserves`, `oracles`, `limiters`, and `withdrawers` maps to Protobuf schemas (`proto/maat/`) and serialize them directly to the `KVStore` inside the respective keepers.
- [ ] Update all keeper method signatures to propagate `ctx sdk.Context` (or `context.Context`).
- [ ] Replace `AllDenoms()` map iteration with sorted denom lists or KVStore iterator loops to ensure determinism.

### Phase 2: Security, Token Transfers, and Swap Pricing
- [ ] Integrate bank keeper transfers (burns/mints) directly inside `x/market/keeper/msg_server.go`.
- [ ] Calculate exchange rates inside `Swap` based on the FromDenom/ToDenom oracle prices rather than a 1:1 ratio.
- [ ] Implement VAA replay protection by persisting processed sequence hashes to the KVStore in `x/bridge`.
- [ ] Add the authority check to `ResumeOracle` in `x/oracle`.

### Phase 3: Decimal Scaling & Concurrency Safeguards
- [ ] Add decimal normalization factors in `market.go` to support 18-decimal tokens (wETH).
- [ ] Clone all slices and values returned by queries (like `GetVotes`) to prevent memory leaks or pointer race conditions.
- [ ] Delete executed/cancelled withdrawals from the active bridge list to prevent memory growth.
