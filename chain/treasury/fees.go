// Package treasury implements Ma'at's deterministic fee/spread distribution,
// matching ECONOMICS.md. Integer-only; the remainder from rounding is assigned
// to the reserve buffer so no micro-unit is ever lost or created.
package treasury

import "errors"

const BpsDenom int64 = 10_000

var (
	ErrBadAmount = errors.New("treasury: amount must be non-negative")
	ErrBadSplit  = errors.New("treasury: split must sum to 10000 bps")
)

// Split is the destination weights for spread/fee revenue, in basis points.
// Per ECONOMICS.md default: 40 reserve / 25 rewards / 20 insurance / 15 treasury.
type Split struct {
	ReserveBps   int64
	RewardsBps   int64
	InsuranceBps int64
	TreasuryBps  int64
}

// DefaultSpreadSplit is the documented spread/swap-fee distribution.
var DefaultSpreadSplit = Split{ReserveBps: 4000, RewardsBps: 2500, InsuranceBps: 2000, TreasuryBps: 1500}

// DefaultBridgeOutSplit is the documented bridge-out fee distribution
// (60 reserve / 20 insurance / 20 treasury).
var DefaultBridgeOutSplit = Split{ReserveBps: 6000, RewardsBps: 0, InsuranceBps: 2000, TreasuryBps: 2000}

func (s Split) Validate() error {
	if s.ReserveBps+s.RewardsBps+s.InsuranceBps+s.TreasuryBps != BpsDenom {
		return ErrBadSplit
	}
	return nil
}

// Allocation is the result of distributing an amount.
type Allocation struct {
	Reserve   int64
	Rewards   int64
	Insurance int64
	Treasury  int64
}

func (a Allocation) Total() int64 { return a.Reserve + a.Rewards + a.Insurance + a.Treasury }

// Distribute splits amount by the weights. Rounding remainder goes to Reserve so
// Total() always equals amount exactly (no value leaks).
func Distribute(amount int64, s Split) (Allocation, error) {
	if amount < 0 {
		return Allocation{}, ErrBadAmount
	}
	if err := s.Validate(); err != nil {
		return Allocation{}, err
	}
	a := Allocation{
		Rewards:   amount * s.RewardsBps / BpsDenom,
		Insurance: amount * s.InsuranceBps / BpsDenom,
		Treasury:  amount * s.TreasuryBps / BpsDenom,
	}
	a.Reserve = amount - a.Rewards - a.Insurance - a.Treasury // remainder-safe
	return a, nil
}

// Treasury accumulates the running balances of each fund.
type Treasury struct {
	Reserve   int64
	Rewards   int64
	Insurance int64
	Treasury  int64
}

// Collect distributes amount and adds it to the running balances.
func (t *Treasury) Collect(amount int64, s Split) (Allocation, error) {
	a, err := Distribute(amount, s)
	if err != nil {
		return Allocation{}, err
	}
	t.Reserve += a.Reserve
	t.Rewards += a.Rewards
	t.Insurance += a.Insurance
	t.Treasury += a.Treasury
	return a, nil
}
