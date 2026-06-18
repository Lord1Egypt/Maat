# Setting Up the Ma'at Cosmos SDK Chain

## Prerequisites

- Go 1.23+

> The module `types` are hand-written Go (not generated from the `.proto` files),
> so **no `protoc`/`buf` step is required to build or run the chain**. The `proto/`
> files document the intended wire schema for a future codegen pass.

## Quick Start (one command)

The fastest path — `scripts/testnet-up.sh` runs the full, verified bring-up
(init → key → genesis account → gentx → collect → validate → start):

```bash
go build -o maatd ./cmd/maatd
PATH="$PWD:$PATH" scripts/testnet-up.sh        # builds genesis, then starts the node
# or: START=0 scripts/testnet-up.sh            # build + validate genesis only
```

This was verified end-to-end: the node boots, the validator bonds, and blocks
finalize continuously (see the `app` package `TestAppInitChainPipeline` and the
CI `sdk-app` job for the automated guards).

## Quick Start (manual steps)

```bash
# 1. Build the chain binary
go build -o maatd ./cmd/maatd

# 2. Initialize a single-validator testnet
./maatd init test-node --chain-id maat-testnet-1 --default-denom umaat
./maatd keys add validator --keyring-backend test
./maatd genesis add-genesis-account "$(./maatd keys show validator -a --keyring-backend test)" 1000000000umaat
./maatd genesis gentx validator 500000000umaat --chain-id maat-testnet-1 --keyring-backend test
./maatd genesis collect-gentxs
./maatd genesis validate-genesis

# 3. Start the chain
./maatd start --minimum-gas-prices 0umaat
```

## Architecture

```
cmd/maatd/           chain binary entry point
app/                 SDK app wiring (module registration, keeper dependencies)
x/
  maat/              native MAAT token: vesting, inflation, block rewards
  oracle/            hardened price feed: median + TWAP + circuit breaker
  market/            per-block quote + settlement (the core innovation)
  reserve/           backing invariant: halt if < 100%
  bridge/            bridge-out safety: caps + delay queue
  treasury/          fee/spread distribution to 4 funds
chain/               verified cores (SDK-free, tested independently)
proto/               protobuf definitions for all modules
```

## Module Dependencies (wiring order)

```
x/maat      -> bankKeeper (for minting block rewards)
x/oracle    -> (standalone — price feeds)
x/treasury  -> (standalone — accumulates fees)
x/market    -> oracleKeeper, treasuryKeeper, reserveKeeper
x/reserve   -> marketKeeper (reads reserve state)
x/bridge    -> marketKeeper (calls BridgeIn/BridgeOut on reserves)
```

## Running the Verified Cores (no SDK deps needed)

```bash
cd chain && go test ./... -v      # all core tests
cd chain && go run ./cmd/demo     # 720-block end-to-end proof
```

## Running the Economic Simulation

```bash
python3 simulation/maat_sim.py --paths 1000 --seed 42
```
