#!/usr/bin/env python3
"""
Fetch real historical hourly price data for calibration / replay.

Pure standard library (urllib). Saves a simple 2-column CSV:
    timestamp_ms,price_usd

Source: CoinGecko public API (no key, no geo-block). For days in [2,90] the
free endpoint returns ~hourly granularity automatically.

Usage:
    python3 fetch_data.py                      # ETH, 90 days -> data/eth_hourly.csv
    python3 fetch_data.py --coin bitcoin --days 90 --out data/btc_hourly.csv

If the network is unavailable, the existing CSV in data/ is kept (the sim and
calibrate scripts work offline against whatever CSV is already committed).
"""

import argparse
import csv
import json
import os
import sys
import urllib.request

API = "https://api.coingecko.com/api/v3/coins/{coin}/market_chart"


def fetch(coin: str, days: int, vs: str = "usd"):
    url = f"{API.format(coin=coin)}?vs_currency={vs}&days={days}&interval=hourly"
    req = urllib.request.Request(url, headers={"User-Agent": "maat-sim/1.0"})
    with urllib.request.urlopen(req, timeout=30) as r:
        data = json.loads(r.read().decode())
    return data.get("prices", [])


def main(argv=None):
    ap = argparse.ArgumentParser(description="Fetch hourly price data for Ma'at calibration")
    ap.add_argument("--coin", default="ethereum")
    ap.add_argument("--days", type=int, default=90)
    ap.add_argument("--vs", default="usd")
    ap.add_argument("--out", default="data/eth_hourly.csv")
    p = ap.parse_args(argv)

    try:
        prices = fetch(p.coin, p.days, p.vs)
    except Exception as e:  # offline / rate-limited -> keep existing CSV
        print(f"fetch failed ({e}); keeping any existing CSV at {p.out}", file=sys.stderr)
        return 0 if os.path.exists(p.out) else 1

    if not prices:
        print("no prices returned", file=sys.stderr)
        return 1

    os.makedirs(os.path.dirname(p.out) or ".", exist_ok=True)
    with open(p.out, "w", newline="") as f:
        w = csv.writer(f)
        w.writerow(["timestamp_ms", "price_usd"])
        for ts, px in prices:
            w.writerow([int(ts), f"{px:.6f}"])

    print(f"wrote {len(prices)} points to {p.out} "
          f"({p.coin}/{p.vs}, {p.days}d): {prices[0][1]:.2f} -> {prices[-1][1]:.2f}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
