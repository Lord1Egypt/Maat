// Package token implements the deterministic MAAT supply logic for x/maat:
// the vesting schedule (ECONOMICS.md distribution) and bounded inflation /
// staking-reward math. Integer-only, consensus-safe.
package token

import "errors"

const (
	BpsDenom int64 = 10_000
	// TotalSupply is the fixed MAAT supply: 1,000,000,000 MAAT.
	TotalSupply int64 = 1_000_000_000
	// MaxInflationBps caps annual inflation at 10% (ECONOMICS.md).
	MaxInflationBps int64 = 1_000
)

var (
	ErrBadInflation = errors.New("token: inflation exceeds 10% cap")
	ErrBadMonth     = errors.New("token: month must be >= 0")
)

// Allocation is one bucket of the MAAT distribution with its unlock terms.
// TGEBps unlocks at month 0; the remainder unlocks linearly over LinearMonths
// after CliffMonths have elapsed.
type Allocation struct {
	Name         string
	Total        int64
	TGEBps       int64
	CliffMonths  int
	LinearMonths int
}

// Standard is the ECONOMICS.md distribution. Sums to TotalSupply.
func Standard() []Allocation {
	return []Allocation{
		{"Community Sale", 300_000_000, 1000, 0, 12},  // 10% TGE, 12mo linear
		{"Treasury Reserve", 250_000_000, 10000, 0, 0}, // DAO-controlled, unlocked
		{"Validator Rewards", 200_000_000, 0, 0, 60},  // emitted over 5 years
		{"Team", 150_000_000, 0, 12, 24},              // 12mo cliff, 24mo linear
		{"Advisors/Partners", 50_000_000, 0, 6, 18},   // 6mo cliff, 18mo linear
		{"Airdrop", 50_000_000, 10000, 0, 0},          // at TGE
	}
}

// Unlocked returns how much of an allocation has vested by the given month.
func (a Allocation) Unlocked(month int) (int64, error) {
	if month < 0 {
		return 0, ErrBadMonth
	}
	tge := a.Total * a.TGEBps / BpsDenom
	rest := a.Total - tge
	if month <= a.CliffMonths {
		return tge, nil
	}
	elapsed := month - a.CliffMonths
	if a.LinearMonths == 0 || elapsed >= a.LinearMonths {
		return a.Total, nil
	}
	return tge + rest*int64(elapsed)/int64(a.LinearMonths), nil
}

// CirculatingSupply sums the unlocked amounts across all allocations at month.
func CirculatingSupply(allocs []Allocation, month int) (int64, error) {
	var total int64
	for _, a := range allocs {
		u, err := a.Unlocked(month)
		if err != nil {
			return 0, err
		}
		total += u
	}
	return total, nil
}

// SumAllocations returns the total of all allocation caps (should == TotalSupply).
func SumAllocations(allocs []Allocation) int64 {
	var s int64
	for _, a := range allocs {
		s += a.Total
	}
	return s
}

// BlockReward returns the per-block validator reward for a given circulating
// supply, annual inflation (bps, capped at 10%), and blocks per year.
func BlockReward(supply, annualBps, blocksPerYear int64) (int64, error) {
	if annualBps < 0 || annualBps > MaxInflationBps {
		return 0, ErrBadInflation
	}
	if blocksPerYear <= 0 {
		return 0, nil
	}
	return supply * annualBps / BpsDenom / blocksPerYear, nil
}
