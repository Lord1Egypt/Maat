#!/bin/bash
#
# testnet-up.sh — bring up a single-node Ma'at localnet from scratch and start it.
#
# This is the full, verified flow: init -> key -> genesis account -> gentx ->
# collect-gentxs -> validate-genesis -> (optional custom params) -> start.
# Without the validator (gentx) steps a chain cannot produce blocks.
#
# Usage:
#   scripts/testnet-up.sh            # build genesis, then `maatd start`
#   START=0 scripts/testnet-up.sh    # build + validate genesis only, don't start
#
set -euo pipefail

BIN="${BIN:-maatd}"
HOME_DIR="${HOME_DIR:-$(pwd)/localnode}"
CHAIN_ID="${CHAIN_ID:-maat-testnet-1}"
DENOM="${DENOM:-umaat}"
KEYRING="${KEYRING:-test}"
MONIKER="${MONIKER:-localnode}"
START="${START:-1}"

echo ">>> Resetting $HOME_DIR"
rm -rf "$HOME_DIR"

echo ">>> init"
"$BIN" init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR" --default-denom "$DENOM" --overwrite >/dev/null

echo ">>> create validator key"
"$BIN" keys add validator --keyring-backend "$KEYRING" --home "$HOME_DIR" >/dev/null
VAL_ADDR="$("$BIN" keys show validator -a --keyring-backend "$KEYRING" --home "$HOME_DIR")"
echo "    validator: $VAL_ADDR"

echo ">>> fund genesis account"
"$BIN" genesis add-genesis-account "$VAL_ADDR" "1000000000${DENOM}" --home "$HOME_DIR"

echo ">>> apply custom Ma'at module params"
# NOTE: oracle/market/bridge serialize Go field names (no json tags); maat uses
# snake_case json tags. Keys below match the actual genesis schema produced by init.
python3 - "$HOME_DIR/config/genesis.json" <<'PY'
import json, sys
path = sys.argv[1]
gen = json.load(open(path))
st = gen["app_state"]
st["oracle"]["params"] = {"AlphaNum": 8, "AlphaDen": 10, "MaxDevBps": 500, "BreakBps": 5000, "VotePeriod": 1}
st["market"]["params"] = {"BaseBps": 15, "VolMultBps": 50, "SkewMaxBps": 100, "MaxBps": 1000}
st["bridge"]["params"] = {"DailyCapUnits": 1000000000, "WindowBlocks": 14400, "LargeTxUnits": 50000000, "DelayBlocks": 14400}
st["maat"]["params"] = {"annual_inflation_bps": 800, "blocks_per_year": 5256000}
json.dump(gen, open(path, "w"), indent=2)
print("    custom params written")
PY

echo ">>> gentx (self-delegate 500 MAAT)"
"$BIN" genesis gentx validator "500000000${DENOM}" --chain-id "$CHAIN_ID" --keyring-backend "$KEYRING" --home "$HOME_DIR" >/dev/null

echo ">>> collect-gentxs"
"$BIN" genesis collect-gentxs --home "$HOME_DIR" >/dev/null

echo ">>> validate-genesis"
"$BIN" genesis validate-genesis --home "$HOME_DIR"

echo ""
echo "Local genesis ready at $HOME_DIR with custom Ma'at parameters."

if [ "$START" = "1" ]; then
  echo ">>> starting node (Ctrl-C to stop)"
  exec "$BIN" start --home "$HOME_DIR" --minimum-gas-prices "0${DENOM}"
else
  echo "To start: $BIN start --home $HOME_DIR --minimum-gas-prices 0${DENOM}"
fi
