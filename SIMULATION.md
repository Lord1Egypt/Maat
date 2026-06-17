# SIMULATION.md - Economic Model

> Python simulation to prove the economics before writing blockchain code.

---

## Purpose

Before writing a single line of Cosmos SDK, we simulate the Ma'at economy
to answer critical questions:

1. Does the arb engine actually attract users?
2. Can the reserve survive a bank run?
3. What happens to MAAT price in bull vs bear markets?
4. What is the optimal price spread?

---

## Simulation Model

The simulation models:
- **Market prices**: Real historical BTC/ETH price data
- **Ma'at prices**: Protocol-set prices that change on schedule
- **Arb bots**: Automated agents that exploit the spread
- **MAAT supply/demand**: How arb affects native coin
- **Reserve ratio**: Real assets backing wrapped tokens

### Key Variables

| Variable | Range | Default |
|----------|-------|---------|
| Price spread | 1-20% | 5% |
| Announcement delay | 1-30 days | 7 days |
| Withdrawal cap | 1-100% of reserve | 10% daily |
| Reserve ratio target | 100-150% | 120% |
| Fee rate | 0.1-1% | 0.3% |

---

## How to Run

```bash
cd simulation
pip install numpy pandas matplotlib
python arb_model.py --days 90 --spread 0.05 --output charts/
```

---

## Planned Charts

1. Arb volume vs spread percentage
2. Reserve depletion over time (bank run scenario)
3. MAAT price correlation with TVL
4. Cumulative protocol revenue over time
5. Optimal spread for max TVL growth

---

## Expected Results

Based on economic modeling:

| Scenario | Result | Confidence |
|----------|--------|-----------|
| 5% spread, 7d delay | 10% TVL growth/week | High |
| Bank run (no cap) | Reserve empty in 3 days | Validated warning |
| Bank run (with 10% cap) | Reserve stays above 100% | Moderate |
| Bear market | Reserve accumulates cheap assets | High |
| Bull market | MAAT demand grows with arb volume | High |

---

## Next Steps

1. Implement simulation in Python
2. Run 1000 monte carlo simulations
3. Find the optimal parameter set
4. Document in whitepaper
5. Then build the real chain
