package reserve

import (
	"testing"
)

func TestReserveIndex(t *testing.T) {
	weights := []IndexWeight{
		{Denom: "weth", WeightBps: 5000},
		{Denom: "wbtc", WeightBps: 3000},
		{Denom: "wsol", WeightBps: 2000},
	}

	ri, err := NewReserveIndex(weights, 100)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}

	// Validate weights validation
	badWeights := []IndexWeight{
		{Denom: "weth", WeightBps: 5000},
	}
	_, err = NewReserveIndex(badWeights, 100)
	if err != ErrInvalidWeights {
		t.Fatalf("expected ErrInvalidWeights, got %v", err)
	}

	// Rebalance
	newWeights := []IndexWeight{
		{Denom: "weth", WeightBps: 4000},
		{Denom: "wbtc", WeightBps: 4000},
		{Denom: "wsol", WeightBps: 2000},
	}
	err = ri.Rebalance(newWeights)
	if err != nil {
		t.Fatalf("rebalance failed: %v", err)
	}

	// Accrue management fee
	// NAV = $1,000,000, fee = 200 bps (2%), elapsed = 1000 blocks, blocks/year = 100,000
	// fee = 1,000,000 * 200 * 1000 / (10000 * 100,000) = 200 USD
	fee, err := ri.AccrueManagementFee(1100, 1000000, 200, 100000)
	if err != nil {
		t.Fatalf("fee accrual failed: %v", err)
	}
	if fee != 200 {
		t.Fatalf("expected fee 200, got %d", fee)
	}
	if ri.FeeAccruedUSD != 200 {
		t.Fatalf("expected accrued fee 200, got %d", ri.FeeAccruedUSD)
	}
	if ri.LastFeeHeight != 1100 {
		t.Fatalf("expected LastFeeHeight 1100, got %d", ri.LastFeeHeight)
	}
}
