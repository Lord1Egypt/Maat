# STATUS.md — Project State & Continuation Handoff

> **Read this first if you're picking up Ma'at cold** (new session, after a context
> reset, or a weekly-limit interruption). It captures where we are, what was decided
> and why, and exactly what to do next.

**Last updated:** 2026-06-19
**Phase:** 2 (Cosmos SDK Chain) — **COMPLETE + the node now actually BOOTS and PRODUCES BLOCKS.** All modules/keepers scaffolded, app boot wiring fixed, verified end-to-end (init → validator → sustained finalized blocks; 33+ blocks in a live run).

---

## 1. TL;DR — the corrected thesis

Ma'at's original v0.1 idea (**fixed off-market price + "guaranteed risk-free arbitrage"
for users**) was structurally unsustainable: the user's "guaranteed profit" was paid out
of the reserve on every trade — the same coins — so the reserve drained to zero by
construction, exactly like Terra/UST. This was abandoned.

**v0.2 (current model):** quote **oracle mid ± spread, fixed per block**.
- Fixed *within* the block → no slippage, no MEV (the good part, kept).
- Tracks the market via a hardened multi-source TWAP oracle → always solvent.
- The protocol **earns** the spread → the reserve **grows** instead of draining.
- Value prop flipped from "free money" → **guaranteed best execution + no MEV + transparent backing.**

Full proof and reasoning: [REDESIGN.md](REDESIGN.md).

---

## 2. What has been done (merged to `main`)

| PR | What | Status |
|----|------|--------|
| #1 | `fix:` replace unsustainable fixed-price model with spread-capture (v0.2). Rewrote all 10 docs + added REDESIGN.md & BUILD_PLAN.md | ✅ Merged |
| #2 | `feat:` economic stop-gate simulation + first CI/CD pipeline | ✅ Merged |
| #3 | `docs:` STATUS.md continuation/handoff doc | ✅ Merged |
| #4 | `feat:` calibration vs real ETH data + real-data CI gate | ✅ Merged |
| #5 | `feat:` deterministic Go pricing core (x/oracle + x/market) + Go CI | ✅ Merged |
| #6 | `feat:` bridge withdrawal caps/delay queue + treasury fee splits | ✅ Merged |
| #7 | `feat:` end-to-end integration scenario + demo + `make test` runner | ✅ Merged |
| #8 | `docs:` WHITEPAPER.md (synthesis + Phase 1 deliverable) | ✅ Merged |
| #9 | `test:` Go fuzz tests for safety invariants | ✅ Merged |
| #10 | `feat:` MAAT token core (x/maat) — vesting + inflation | ✅ Merged |
| #11 | `docs:` x/* module integration spec (SDK wiring guide) | ✅ Merged |
| #12 | `feat:` governance tally core (x/gov) | ✅ Merged |
| #13 | `docs:` README CI badge + verification section | ✅ Merged |
| #14 | `feat:` scaffold Cosmos SDK v0.50 modules wrapping the verified cores | ✅ Merged |
| #15 | `feat:` Wormhole VAA parser + x/bridge BridgeIn integration | ✅ Merged |
| #16 | `feat:` Phase 4–7 — dashboard, feeder daemon, reserve index, Pharaoh Bonds, testnet, Arabic remittance | ✅ Merged |
| #17 | `fix:` **make the node actually boot & produce blocks** (app ABCI wiring + CLI) | ⏳ This PR |

**Workflow in use:** branch → PR → merge (never direct push to `main`). CI must be green
before merge.

### PR #17 — node boot wiring (the gap between "compiles" and "runs")

The SDK scaffold (#14–#16) compiled and passed keeper unit tests, but the binary
could not boot a chain. Fixed and **verified by running** (`scripts/testnet-up.sh`
brings up a single-validator localnet that finalizes blocks continuously):

1. **`maatd` panicked on every command** — `client.Context` built without `.WithViper("")`
   (nil Viper in `ReadFromClientConfig`). Also added `initAppConfig`/`initCometBFTConfig`
   so `InterceptConfigsPreRunHandler` gets non-nil defaults.
2. **No client/genesis commands** — `root.go` only had `server.AddCommands`. Added
   `init`, `genesis`, `keys`, `tx`, `query`, `status`, derived from a pre-instantiated
   app's `BasicModuleManager` (so staking/gov CLI builders get real codecs).
3. **No ABCI lifecycle in `app.go`** — added `InitChainer`, `PreBlocker`, `BeginBlocker`,
   `EndBlocker`, and the `AnteHandler`, plus `SetOrderPreBlockers`. Without these the
   module BeginBlock/EndBlock (x/maat mint, x/oracle, x/market) never ran and `InitChain`
   panicked on a nil chainer.
4. **Missing `x/consensus`** — required in SDK 0.50; added the module + `SetParamStore`,
   else `InitChain` panics storing consensus params.
5. **`RegisterStores` mounted fresh key objects** instead of the keepers' own keys →
   runtime `ctx.KVStore()` lookups would fail. Now mounts the real maps.
6. **`MakeEncodingConfig` had no address codecs** in the interface registry's signing
   options → `gentx` failed ("requires a proper address codec"). Rebuilt with
   `NewInterfaceRegistryWithOptions`.
7. **`DefaultNodeHome` was empty** (no `init()`) → CometBFT `mkdir ""`. Set to `~/.maat`.
8. **`mm.RegisterServices` was never called** → `MsgCreateValidator` had "no message
   handler" and InitChain panicked delivering the gentx. Now registered.
9. **Staking hooks never wired** → slashing had "no validator signing info found" and the
   node hit a CONSENSUS FAILURE at height 2. Added `StakingKeeper.SetHooks(distr, slashing)`.
10. **`newApp` didn't pass `server.DefaultBaseappOptions`** → baseapp chainID empty →
    "invalid chain-id on InitChain". Now wired (pruning/snapshots/min-gas/SetChainID).

Guards added: `app.TestAppInitChainPipeline` (drives InitChain through every module's
genesis), a new CI `sdk-app` job (builds the whole chain + runs boot/keeper tests), and
`scripts/testnet-up.sh` rewritten to the full verified flow.

**Known follow-up (not blocking):** module `query`/`tx` CLI subcommands need AutoCLI
(`cosmossdk.io/client/v2`), which would add a new dependency. The gRPC/REST query
services ARE registered, so querying works over those today.

---

## 3. Simulation result (the Phase 0 stop-gate)

`simulation/maat_sim.py` (pure stdlib, zero deps). 1000 Monte Carlo paths, same paths for
both models:

| Model | Reserve grows | Median outcome |
|-------|---------------|----------------|
| **v0.2 spread-capture** | **100%** of paths | **+4.5%** |
| v0.1 fixed off-market (abandoned) | 0% of paths | **−86%** (drains) |

- Backing stayed ≥100% in 100% of paths; flash-crash absorbed by the buffer.
- The gate has **teeth**: a laggy oracle + near-zero spread correctly **fails** it (exit 1).
- **Now also validated on REAL ETH data** (90d hourly, CoinGecko): calibrated to 45.5%
  annual vol / ~1%-per-hour jumps / −37.9% max drawdown; block-bootstrap gate still PASSES
  (+4.45% median, v0.1 −86%). Tools: `fetch_data.py`, `calibrate.py`; replay via `--price-csv`.
- Run it: `python3 simulation/maat_sim.py --paths 1000` (exit 0 = pass, 1 = fail).

---

## 4. CI/CD (enforced on every push/PR)

`.github/workflows/ci.yml` — two jobs:
1. **economics-gate** — compiles + runs the 1000-path stop-gate (~20s in CI). Red build if economics break.
2. **docs-integrity** — fails CI if abandoned v0.1 claims ("guaranteed risk-free arb",
   "we are the oracle") reappear as live copy. (Citing/warning about them is allowed.)

---

## 5. NEXT STEPS

- **B — Calibrate vs real ETH data** ✅ DONE (PR #4). Params locked: annual_vol 0.4554,
  jump 0.0291, crash 0.38; gate passes on real-data bootstrap.
- **A — Cosmos SDK Chain (Phase 2)** ✅ DONE.
  - ✅ Deterministic pricing **core** in `chain/pricing/` (Go, integer fixed-point, 12
    unit tests passing, Go CI job added).
  - ✅ Protection cores (Go, tested): `bridge/limits.go` (per-asset withdrawal cap + delay queue).
  - ✅ Token core (Go, tested): `token/token.go` (`x/maat`) — vesting schedule and inflation math.
  - ✅ Integration: `chain/scenario` wires oracle+market+treasury+bridge through 720 blocks.
  - ✅ Governance tally core (`governance/tally.go`, `x/gov`) rules from PLANNED_ECONOMY, tested.
  - ✅ SDK wiring guide written: `chain/MODULE_SPEC.md`.
  - ✅ Scaffolded and implemented all 6 Cosmos SDK modules under `x/` (maat, oracle, market, reserve, bridge, treasury), wired in `app/app.go`, entry binary in `cmd/maatd`, all fully verified by keeper-level tests.
- **B — Ethereum Bridge (Phase 3)** ⏭️ NEXT.
  - [ ] Integrate Wormhole EVM transport layer (Sepolia testnet deployment).
  - [ ] Configure daily bridge cap, delays, and multisig security.
  - [ ] Schedule independent security audit of the bridge integration.

(See [BUILD_PLAN.md](BUILD_PLAN.md) for the full phased plan and stop-gate rules.)

---

## 6. File map (where everything lives)

| File | Purpose |
|------|---------|
| `REDESIGN.md` | Why v0.1 fails + the corrected v0.2 model (the core reasoning) |
| `BUILD_PLAN.md` | Phased execution plan + the simulation stop-gate rule |
| `STATUS.md` | **This file** — current state & continuation handoff |
| `ECONOMICS.md` / `PLANNED_ECONOMY.md` | Spread mechanics, reserve, pricing detail |
| `ARCHITECTURE.md` | Cosmos SDK modules (incl. required `x/oracle`) |
| `RISKS.md` | Risk register (Risk #0 = structural drain, now fixed) |
| `VISION.md` / `STRATEGY.md` / `COMPETITION.md` | Philosophy, growth plan, landscape |
| `SIMULATION.md` / `simulation/` | The stop-gate sim + how to run it |
| `.github/workflows/ci.yml` | CI/CD: economics-gate + docs-integrity |

---

## 7. Hard rules to honor when resuming

1. **No chain code until the stop-gate is green** (it is — keep it that way).
2. **Never reintroduce the v0.1 model** (off-market fixed price / "we are the oracle" /
   guaranteed user arb). The oracle is mandatory; the protocol earns the spread.
3. **branch → PR → merge**, CI green before merge.
4. **Zero/standard-library dependencies** where possible (sim is pure stdlib — keep it).
5. Don't market "guaranteed/risk-free profit" — it's a security and a legal landmine.
