#!/bin/bash
set -e

# Initialize home folder
maatd init localnode --chain-id maat-testnet-1 --home ./localnode --overwrite

# Build genesis modifications using python
python3 -c '
import json

genesis_path = "./localnode/config/genesis.json"
with open(genesis_path, "r") as f:
    gen = json.load(f)

# Update x/oracle params
gen["app_state"]["oracle"] = {
    "params": {
        "alpha_num": 8,
        "alpha_den": 10,
        "max_dev_bps": 500,
        "break_bps": 5000,
        "vote_period": 1
    },
    "feeders": ["maat1valoper1", "maat1valoper2"]
}

# Update x/market params
gen["app_state"]["market"] = {
    "params": {
        "base_bps": 15,
        "vol_mult_bps": 50,
        "skew_max_bps": 100,
        "max_bps": 1000
    }
}

# Update x/bridge params
gen["app_state"]["bridge"] = {
    "params": {
        "daily_cap_units": 1000000000,
        "window_blocks": 14400,
        "large_tx_units": 50000000,
        "delay_blocks": 14400
    }
}

# Update x/maat params
gen["app_state"]["maat"] = {
    "params": {
        "annual_inflation_bps": 800,
        "blocks_per_year": 5256000,
        "vesting_start": 100000
    }
}

with open(genesis_path, "w") as f:
    json.dump(gen, f, indent=2)
'

echo "Local genesis configured successfully with custom Ma'at parameters."
