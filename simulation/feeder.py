import json
import sys
import urllib.request
import time
from typing import Dict, List, Optional

COINGECKO_ETH_URL = "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"

def fetch_eth_price_live() -> Optional[float]:
    try:
        req = urllib.request.Request(
            COINGECKO_ETH_URL,
            headers={'User-Agent': 'Mozilla/5.0'}
        )
        with urllib.request.urlopen(req, timeout=5) as response:
            data = json.loads(response.read().decode())
            return float(data["ethereum"]["usd"])
    except Exception:
        return None

def simulate_feeders(base_price: float, num_feeders: int = 5, max_dev: float = 0.02) -> List[int]:
    import random
    random.seed(int(time.time()))
    prices: List[int] = []
    for _ in range(num_feeders):
        dev = random.uniform(-max_dev, max_dev)
        feeder_price = base_price * (1.0 + dev)
        prices.append(int(feeder_price * 1_000_000))
    return prices

def main() -> None:
    print("Starting Ma'at Price Feeder Daemon...")
    live_eth = fetch_eth_price_live()
    if live_eth:
        print(f"Fetched live ETH price from CoinGecko: ${live_eth:.2f}")
        base = live_eth
    else:
        print("Failed to fetch live price. Replaying fallback base price: $3000.00")
        base = 3000.0

    votes = simulate_feeders(base)
    
    print("\nSimulated validator consensus votes (micro-USD):")
    for i, vote in enumerate(votes):
        print(f"  Validator #{i+1} : {vote} (equivalent to ${vote/1_000_000:.4f})")

    txs: List[Dict[str, any]] = []
    for i, vote in enumerate(votes):
        tx = {
            "type": "cosmos-sdk/MsgSubmitPrice",
            "feeder": f"maatvaloper{i+1}",
            "denom": "weth",
            "price": vote,
            "timestamp": int(time.time())
        }
        txs.append(tx)

    print("\nPrepared transactions ready for consensus broadcast:")
    print(json.dumps(txs, indent=2))

if __name__ == "__main__":
    main()
