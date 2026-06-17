#!/usr/bin/env python3
"""
Calibrate Ma'at simulation parameters from a real price CSV.

Pure standard library. Reads data/eth_hourly.csv (timestamp_ms,price_usd) and
estimates the parameters maat_sim.py needs, separating ordinary diffusion from
jumps so the flash-crash stress is grounded in real tail behaviour:

  - annual_vol      : annualised volatility of hourly log returns (ex-jumps)
  - jump_prob       : hourly probability of a |return| > 4*sigma jump
  - jump_size       : stdev of those jump returns
  - max_drawdown    : worst peak-to-trough over the window (sanity vs crash-size)
  - worst_hour      : largest single-hour drop (informs oracle freshness need)

Outputs a ready-to-paste maat_sim.py command and writes calibrated_params.json.

Usage:
    python3 calibrate.py                       # uses data/eth_hourly.csv
    python3 calibrate.py --csv data/btc_hourly.csv --emit-json out.json
"""

import argparse
import csv
import json
import math
import statistics
import sys

HOURS_PER_YEAR = 365 * 24


def load_prices(path):
    prices = []
    with open(path) as f:
        for row in csv.DictReader(f):
            prices.append(float(row["price_usd"]))
    return prices


def log_returns(prices):
    return [math.log(prices[i] / prices[i - 1]) for i in range(1, len(prices))
            if prices[i] > 0 and prices[i - 1] > 0]


def max_drawdown(prices):
    peak = prices[0]
    mdd = 0.0
    for px in prices:
        peak = max(peak, px)
        mdd = min(mdd, px / peak - 1.0)
    return mdd


def calibrate(prices):
    rets = log_returns(prices)
    sd_all = statistics.pstdev(rets)
    # separate jumps (|r| > 4 sigma) from ordinary diffusion
    jumps = [r for r in rets if abs(r) > 4 * sd_all]
    body = [r for r in rets if abs(r) <= 4 * sd_all]
    sd_body = statistics.pstdev(body) if len(body) > 1 else sd_all

    annual_vol = sd_body * math.sqrt(HOURS_PER_YEAR)
    jump_prob = len(jumps) / len(rets) if rets else 0.0
    jump_size = statistics.pstdev(jumps) if len(jumps) > 1 else (abs(jumps[0]) if jumps else 0.0)
    worst_hour = min(rets) if rets else 0.0
    mdd = max_drawdown(prices)

    return {
        "n_points": len(prices),
        "n_returns": len(rets),
        "annual_vol": round(annual_vol, 4),
        "hourly_sigma_body": round(sd_body, 6),
        "jump_prob": round(jump_prob, 5),
        "jump_size": round(jump_size, 4),
        "n_jumps": len(jumps),
        "worst_hour_return": round(worst_hour, 4),
        "max_drawdown": round(mdd, 4),
    }


def main(argv=None):
    ap = argparse.ArgumentParser(description="Calibrate Ma'at sim params from real data")
    ap.add_argument("--csv", default="data/eth_hourly.csv")
    ap.add_argument("--emit-json", default="calibrated_params.json")
    p = ap.parse_args(argv)

    try:
        prices = load_prices(p.csv)
    except FileNotFoundError:
        print(f"CSV not found: {p.csv} (run fetch_data.py first)", file=sys.stderr)
        return 1
    if len(prices) < 50:
        print("not enough data points", file=sys.stderr)
        return 1

    c = calibrate(prices)

    print("=" * 60)
    print(f" Calibration from {p.csv}")
    print("=" * 60)
    print(f" data points          : {c['n_points']}")
    print(f" annual vol (ex-jump) : {c['annual_vol']*100:.1f}%")
    print(f" jumps detected       : {c['n_jumps']}  (hourly prob {c['jump_prob']})")
    print(f" jump size (stdev)    : {c['jump_size']*100:.1f}%")
    print(f" worst single hour    : {c['worst_hour_return']*100:.1f}%")
    print(f" max drawdown         : {c['max_drawdown']*100:.1f}%")
    print("-" * 60)
    crash = max(0.10, round(-c["max_drawdown"], 2))
    cmd = (f" python3 maat_sim.py --paths 1000 "
           f"--annual-vol {c['annual_vol']} "
           f"--jump-prob {c['jump_prob']} "
           f"--jump-size {c['jump_size']} "
           f"--crash-size {crash}")
    print(" Calibrated stop-gate command:")
    print(cmd)
    print("=" * 60)

    with open(p.emit_json, "w") as f:
        json.dump({**c, "suggested_command": cmd.strip()}, f, indent=2)
    print(f" wrote {p.emit_json}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
