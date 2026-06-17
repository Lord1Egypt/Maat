#!/usr/bin/env python3
"""
Ma'at economic simulation — the Phase 0 STOP-GATE.

Proves (or disproves) that the v0.2 "oracle mid +/- spread, fixed per block"
market-making model lets the reserve GROW, while the abandoned v0.1
"fixed off-market price" model DRAINS the reserve (Terra-style).

Pure standard library. No external dependencies (numpy/pandas not required).
Charts are optional and only drawn if matplotlib is installed (--charts).

The economic core, per block (we step hourly for tractability):
    true_mid : geometric Brownian motion (the real market price)
    oracle   : an EMA of true_mid -> models oracle freshness/lag
    delta    : (true_mid - oracle) / oracle  -> the protocol's mispricing
    quote    : buy = oracle*(1-s), sell = oracle*(1+s)   (s = per-side spread)

    Revenue  : retail flow pays the spread          -> + retail_vol * s
    Cost     : when |delta| > s, arbitrageurs hit the stale side and extract
               (|delta| - s) on capped arb volume   -> - adverse selection
    Jumps    : occasional gaps (flash crash) widen delta sharply -> stress

PnL accrues to the reserve buffer (default 15% over-collateralization).
PASS if the reserve ends >= it started in >= 95% of Monte Carlo paths,
backing never breaks 100%, and net revenue beats infra cost.

Exit code 0 = gate passed, 1 = gate failed (so CI can enforce it).
Usage:
    python maat_sim.py                      # default 1000-path gate
    python maat_sim.py --paths 5000 --days 180
    python maat_sim.py --base-spread 0.0015 --oracle-alpha 0.8
    python maat_sim.py --charts             # also write PNGs (needs matplotlib)
"""

import argparse
import math
import random
import statistics
import sys

# annualised ETH vol ~80%; convert to per-hour
HOURS_PER_YEAR = 365 * 24


def hourly_sigma(annual_vol: float) -> float:
    return annual_vol / math.sqrt(HOURS_PER_YEAR)


def simulate_path(p, rng, mode="v2", flash=False):
    """Simulate one price path. Returns a dict of outcomes.

    mode: "v2" = oracle mid +/- spread (spread-capture)
          "v1" = fixed off-market price (constant discount) -- the broken model
    flash: if True, inject a flash-crash jump to stress backing.
    """
    steps = p.days * 24
    sig = hourly_sigma(p.annual_vol)

    true_mid = 1.0
    oracle = 1.0
    buffer = p.buffer_start          # fraction of supply value (e.g. 0.15)
    inventory = 0.0                  # net asset inventory (fraction of supply)
    cum_revenue = 0.0
    cum_adverse = 0.0
    min_buffer = buffer
    max_abs_inventory = 0.0

    # where to drop the flash crash (somewhere in the middle third)
    crash_at = rng.randint(steps // 3, 2 * steps // 3) if flash else -1

    for t in range(steps):
        # --- market mid evolves (GBM) ---
        shock = rng.gauss(0.0, sig)
        true_mid *= math.exp(shock - 0.5 * sig * sig)

        # --- occasional jump / flash crash ---
        if t == crash_at:
            true_mid *= (1.0 - p.crash_size)        # e.g. -35% gap
        elif rng.random() < p.jump_prob:
            jump = rng.gauss(0.0, p.jump_size)
            true_mid *= math.exp(jump)

        # --- oracle freshness: EMA of the true mid (alpha high = fresh) ---
        oracle += p.oracle_alpha * (true_mid - oracle)

        # --- volatility-scaled spread (v2) ---
        realized = abs(shock)
        spread = p.base_spread * (1.0 + p.vol_mult * realized / sig) if mode == "v2" else 0.0

        if mode == "v2":
            buy = oracle * (1.0 - spread)
            sell = oracle * (1.0 + spread)
            delta = (true_mid - oracle) / oracle           # mispricing vs reality
            edge = abs(delta) - spread                       # arb edge over the spread
        else:  # v1: fixed off-market price, constant discount, NOT tracking market
            sell = true_mid * (1.0 - p.v1_discount)          # protocol sells below market
            buy = true_mid * (1.0 - p.v1_discount)
            spread = 0.0
            delta = p.v1_discount
            edge = p.v1_discount                              # full gap is free to arbs

        # --- inventory-skew pricing (v2): tilt to pull inventory back ---
        skew = -p.skew_k * inventory if mode == "v2" else 0.0

        # --- retail flow pays the spread (both directions) ---
        retail_vol = max(0.0, rng.gauss(p.retail_vol, p.retail_vol * 0.4))
        revenue = retail_vol * spread

        # --- adverse selection: arbs hit the stale side when edge > 0 ---
        adverse = 0.0
        if edge > 0:
            arb_vol = min(p.arb_cap, p.arb_aggr * edge)      # capped per-block
            adverse = arb_vol * edge
            # arbs push inventory the opposite way of the mispricing
            inventory += math.copysign(arb_vol, -delta)

        # retail nudges inventory back toward target via skew incentive
        inventory += skew * retail_vol
        inventory *= (1.0 - p.inventory_decay)               # passive rebalancing

        cum_revenue += revenue
        cum_adverse += adverse
        buffer += revenue - adverse

        min_buffer = min(min_buffer, buffer)
        max_abs_inventory = max(max_abs_inventory, abs(inventory))

    return {
        "final_buffer": buffer,
        "pnl": buffer - p.buffer_start,
        "min_buffer": min_buffer,
        "backing_ok": min_buffer > -1e-9,          # buffer never went negative => >=100% backed
        "revenue": cum_revenue,
        "adverse": cum_adverse,
        "max_inventory": max_abs_inventory,
        "inventory_ok": max_abs_inventory <= p.inventory_cap + 1e-9,
    }


def monte_carlo(p, n, mode="v2", flash=False, seed=0):
    rng = random.Random(seed)
    results = [simulate_path(p, rng, mode=mode, flash=flash) for _ in range(n)]
    pnls = [r["pnl"] for r in results]
    grew = sum(1 for r in results if r["pnl"] >= 0)
    backing_ok = sum(1 for r in results if r["backing_ok"])
    inv_ok = sum(1 for r in results if r["inventory_ok"])
    return {
        "n": n,
        "grew_frac": grew / n,
        "backing_ok_frac": backing_ok / n,
        "inventory_ok_frac": inv_ok / n,
        "median_pnl": statistics.median(pnls),
        "mean_pnl": statistics.mean(pnls),
        "p05_pnl": sorted(pnls)[max(0, int(0.05 * n) - 1)],
        "mean_revenue": statistics.mean(r["revenue"] for r in results),
        "mean_adverse": statistics.mean(r["adverse"] for r in results),
        "results": results,
    }


def pct(x):
    return f"{100 * x:6.2f}%"


def main(argv=None):
    ap = argparse.ArgumentParser(description="Ma'at economic stop-gate simulation")
    ap.add_argument("--paths", type=int, default=1000)
    ap.add_argument("--days", type=int, default=90)
    ap.add_argument("--seed", type=int, default=42, help="fixed seed -> deterministic CI gate")
    ap.add_argument("--base-spread", type=float, default=0.0015, help="per-side spread (15 bps)")
    ap.add_argument("--oracle-alpha", type=float, default=0.8, help="EMA freshness (1=instant, low=laggy)")
    ap.add_argument("--annual-vol", type=float, default=0.80)
    ap.add_argument("--retail-vol", type=float, default=0.01, help="retail turnover per block (frac of supply)")
    ap.add_argument("--arb-aggr", type=float, default=2.0)
    ap.add_argument("--arb-cap", type=float, default=0.02, help="max arb volume per block (frac of supply)")
    ap.add_argument("--vol-mult", type=float, default=0.5, help="spread widening with realized vol")
    ap.add_argument("--buffer-start", type=float, default=0.15, help="over-collateralization buffer")
    ap.add_argument("--skew-k", type=float, default=0.5)
    ap.add_argument("--inventory-cap", type=float, default=0.20)
    ap.add_argument("--inventory-decay", type=float, default=0.02)
    ap.add_argument("--jump-prob", type=float, default=0.002)
    ap.add_argument("--jump-size", type=float, default=0.03)
    ap.add_argument("--crash-size", type=float, default=0.35, help="flash-crash gap for the stress test")
    ap.add_argument("--v1-discount", type=float, default=0.02, help="v0.1 fixed off-market discount")
    ap.add_argument("--infra-cost", type=float, default=0.001, help="infra cost as frac of supply over run")
    ap.add_argument("--pass-threshold", type=float, default=0.95)
    ap.add_argument("--charts", action="store_true", help="write PNG charts (needs matplotlib)")
    ap.add_argument("--output", default="charts/")
    p = ap.parse_args(argv)

    print("=" * 64)
    print(" Ma'at Economic Stop-Gate  —  v0.2 spread-capture vs v0.1 fixed")
    print("=" * 64)
    print(f" paths={p.paths}  days={p.days}  spread={p.base_spread*1e4:.0f}bps  "
          f"oracle_alpha={p.oracle_alpha}  seed={p.seed}")
    print("-" * 64)

    v2 = monte_carlo(p, p.paths, mode="v2", flash=False, seed=p.seed)
    v1 = monte_carlo(p, p.paths, mode="v1", flash=False, seed=p.seed)
    stress = monte_carlo(p, p.paths, mode="v2", flash=True, seed=p.seed + 1)

    print(" MODEL                     reserve grows   median PnL   5th-pct PnL")
    print(f" v0.2 spread-capture       {pct(v2['grew_frac'])}        "
          f"{pct(v2['median_pnl'])}      {pct(v2['p05_pnl'])}")
    print(f" v0.1 fixed off-market     {pct(v1['grew_frac'])}        "
          f"{pct(v1['median_pnl'])}      {pct(v1['p05_pnl'])}")
    print("-" * 64)
    print(f" v0.2 mean spread revenue : {pct(v2['mean_revenue'])} of supply")
    print(f" v0.2 mean adverse cost   : {pct(v2['mean_adverse'])} of supply")
    print(f" v0.2 backing >=100%      : {pct(v2['backing_ok_frac'])} of paths")
    print(f" v0.2 inventory under cap : {pct(v2['inventory_ok_frac'])} of paths")
    print(f" flash-crash backing OK   : {pct(stress['backing_ok_frac'])} of paths")
    print("=" * 64)

    # ---- the gate ----
    checks = {
        "v0.2 reserve grows >= threshold": v2["grew_frac"] >= p.pass_threshold,
        "v0.2 backing never breaks 100%": v2["backing_ok_frac"] >= 0.999,
        "v0.2 inventory stays under cap": v2["inventory_ok_frac"] >= 0.999,
        "v0.2 revenue beats infra cost": v2["mean_revenue"] > p.infra_cost,
        "flash-crash buffer absorbs shock": stress["backing_ok_frac"] >= p.pass_threshold,
        "v0.1 demonstrably drains (sanity)": v1["median_pnl"] < 0,
    }
    all_pass = all(checks.values())
    for name, ok in checks.items():
        print(f"  [{'PASS' if ok else 'FAIL'}] {name}")
    print("=" * 64)

    if p.charts:
        try:
            _draw_charts(p, v2, v1, stress)
            print(f" charts written to {p.output}")
        except ImportError:
            print(" (matplotlib not installed — skipping charts)")

    if all_pass:
        print(" STOP-GATE: PASSED — economics are sound, cleared to build the chain.")
        return 0
    print(" STOP-GATE: FAILED — DO NOT BUILD. Fix the model/params first.")
    return 1


def _draw_charts(p, v2, v1, stress):
    import os
    import matplotlib
    matplotlib.use("Agg")
    import matplotlib.pyplot as plt
    os.makedirs(p.output, exist_ok=True)

    fig, ax = plt.subplots(figsize=(8, 5))
    ax.hist([r["pnl"] for r in v2["results"]], bins=50, alpha=0.7, label="v0.2 spread-capture")
    ax.hist([r["pnl"] for r in v1["results"]], bins=50, alpha=0.7, label="v0.1 fixed off-market")
    ax.axvline(0, color="k", lw=1)
    ax.set_title("Reserve PnL distribution (fraction of supply)")
    ax.set_xlabel("final reserve PnL")
    ax.set_ylabel("paths")
    ax.legend()
    fig.tight_layout()
    fig.savefig(os.path.join(p.output, "pnl_distribution.png"), dpi=120)
    plt.close(fig)


if __name__ == "__main__":
    sys.exit(main())
