# Cosmos SDK Module Spec — wrapping the verified cores

> How to turn `chain/`'s tested, SDK-free cores into real Cosmos SDK v0.50
> modules. Each module is a thin keeper around a core in this repo; the **money
> logic is already written and tested** — this is the plumbing.
>
> Build/verify this in a Go + `ignite` environment (it needs the cosmos-sdk
> dependency tree, which the design sandbox can't compile). Scaffold with:
> `ignite scaffold chain maat --address-prefix maat`.

Common conventions:
- Amounts on-chain use `cosmossdk.io/math.Int`; the cores use `int64` fixed-point.
  Wrap at the keeper boundary (`Int.Int64()` within validated ranges, or port the
  core arithmetic to `math.Int` — the algorithms are identical).
- All params live in each module's `Params` (governance-updatable).
- Every state mutation that the cores guard (backing, caps, breakers) must call
  the core function and propagate its error as an SDK error.

---

## x/maat — native coin, staking, governance  (core: `token/`)

**State**
- `Params{ AnnualInflationBps, BlocksPerYear, VestingStart time }`
- Vesting accounts seeded at genesis from `token.Standard()`.

**Genesis**
- Assert `token.SumAllocations(token.Standard()) == token.TotalSupply` (1B).
- Create a vesting account per allocation with its TGE/cliff/linear terms.

**BeginBlock**
- `reward, err := token.BlockReward(circulating, p.AnnualInflationBps, p.BlocksPerYear)`
- Mint `reward` umaat to the fee collector for validator distribution.

**Msgs**: standard `x/staking`, `x/gov`, `x/distribution` (SDK built-ins).

---

## x/oracle — hardened price feed  (core: `pricing/oracle.go`)

**State**
- `Params{ AlphaNum, AlphaDen, MaxDevBps, BreakBps, VotePeriod }`
- `Mid[denom] int64`, `Halted[denom] bool`
- `Votes[denom][]int64` — feeder submissions for the current period.

**Msgs**
- `MsgSubmitPrice{ feeder, denom, price }` — only whitelisted feeders; appended to `Votes`.

**EndBlock** (every `VotePeriod` blocks)
- For each denom: `mid, err := oracle.Update(votes)` using a per-denom `Oracle`
  built from params + stored `Mid`. Persist new `Mid`; on `ErrOracleHalted` set
  `Halted=true` and emit an event (security council can `Resume`).

**Keeper API** (consumed by x/market): `GetMid(ctx, denom) (int64, halted bool)`.

---

## x/market — per-block quote + settlement  (core: `pricing/market.go`)

**State**
- `SpreadParams[denom]{ BaseBps, VolMultBps, SkewMaxBps, MaxBps }`
- `LastMid[denom]` (for the vol signal), `Inventory[denom]`.

**Msgs**
- `MsgSwap{ trader, fromDenom, toDenom, amount }`

**Handler**
1. `mid, halted := oracleKeeper.GetMid(ctx, denom)`; reject if `halted` or backing
   unhealthy (`reserveKeeper`).
2. `volBps := |mid - LastMid| * 10000 / LastMid`
3. `q := pricing.MakeQuote(mid, sp.EffectiveSpreadBps(volBps), invSkewBps)`
4. `reserve.BuyWrapped` / `reserve.SellWrapped` (returns error on backing breach).
5. Route captured spread to x/treasury (`treasury.Collect`).

Price is read once per block → identical for every tx in the block → no MEV.

---

## x/reserve — backing + custody  (core: `pricing/market.go` Reserve)

**State**: `Reserve[denom]{ AssetUnits, WrappedSupply, SpreadUSD, MinBackingBps }`,
plus published custody addresses for proofs.

**Invariant (EndBlock)**: if any `reserve.BackingBps() < MinBackingBps`, set a
chain-wide `MarketHalted` flag (x/market rejects swaps) until restored. This is
the on-chain circuit breaker.

**Keeper API**: `Healthy(ctx, denom) bool`, `Get(ctx, denom) Reserve`.

---

## x/bridge — cross-chain transport + safety  (core: `bridge/limits.go`)

**State**: `Limiter[denom]{ DailyCapUnits, WindowBlocks, LargeTxUnits, DelayBlocks, ... }`,
`Pending[denom][id]Withdrawal`.

**Msgs**
- `MsgBridgeIn{ proof }` → verify Wormhole VAA → `reserve.BridgeIn(units)`.
- `MsgBridgeOut{ to, amount }` → `lim.RequestOut(amount, height)`:
  - immediate → execute now (`reserve.BridgeOut` + release VAA);
  - pending → store, emit event with `UnlockAt`.
- `MsgExecuteWithdrawal{ id }` → `lim.Execute(id, height)` then `reserve.BridgeOut`.
- `MsgCancelWithdrawal{ id }` (authority/security council) → `lim.Cancel(id)`.

Transport is Wormhole (battle-tested); this module is the **safety layer** around it.

---

## x/treasury — fee/spread distribution  (core: `treasury/fees.go`)

**State**: running `Treasury{ Reserve, Rewards, Insurance, Treasury }`,
`Split` params (`DefaultSpreadSplit`, `DefaultBridgeOutSplit`).

**Keeper API**: `Collect(ctx, amount, split)` → `treasury.Collect` → credit each
fund's module account. Rounding remainder routes to reserve (no value leak).

---

## Build order

1. `ignite scaffold chain` → wire the six modules into `app.go`.
2. Port each core into its keeper (or call it directly with `math.Int` wrapping).
3. Proto: define the Msgs above; generate with `make proto-gen`.
4. Genesis: seed `token.Standard()` vesting + default params.
5. `make test` here stays green (cores unchanged); add module-level keeper tests.
6. Spin a 4-validator testnet; verify `MsgSwap` grows the reserve end-to-end.

The hard, money-handling logic is done and tested. This file is the wiring guide.
