# Setting Up the Ma'at Cosmos SDK Chain

## Prerequisites

- Go 1.22+
- `protoc` (Protocol Buffers compiler)
- `buf` (recommended) or `grpc-gateway` protoc plugins

## Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Generate protobuf code (if using buf)
buf generate

# 3. Build the chain binary
go build -o maatd ./cmd/maatd

# 4. Initialize a single-validator testnet
./maatd init test-node --chain-id maat-testnet-1
./maatd genesis add-genesis-account maat1... 1000000000umaat
./maatd genesis gentx val1 100000000umaat --chain-id maat-testnet-1
./maatd genesis collect-gentxs

# 5. Start the chain
./maatd start
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
