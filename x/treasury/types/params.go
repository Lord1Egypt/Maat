package types

import (
	"fmt"
)

type SplitParams struct {
	ReserveBps   int64 `json:"reserve_bps"`
	RewardsBps   int64 `json:"rewards_bps"`
	InsuranceBps int64 `json:"insurance_bps"`
	TreasuryBps  int64 `json:"treasury_bps"`
}

func (s SplitParams) Validate() error {
	sum := s.ReserveBps + s.RewardsBps + s.InsuranceBps + s.TreasuryBps
	if sum != 10000 {
		return fmt.Errorf("split bps must sum to 10000, got: %d", sum)
	}
	if s.ReserveBps < 0 || s.RewardsBps < 0 || s.InsuranceBps < 0 || s.TreasuryBps < 0 {
		return fmt.Errorf("split bps cannot be negative")
	}
	return nil
}

type Params struct {
	SpreadSplit    SplitParams `json:"spread_split"`
	BridgeOutSplit SplitParams `json:"bridge_out_split"`
}

func DefaultParams() Params {
	return Params{
		SpreadSplit: SplitParams{
			ReserveBps:   4000,
			RewardsBps:   2500,
			InsuranceBps: 2000,
			TreasuryBps:  1500,
		},
		BridgeOutSplit: SplitParams{
			ReserveBps:   6000,
			RewardsBps:   0,
			InsuranceBps: 2000,
			TreasuryBps:  2000,
		},
	}
}

func (p Params) Validate() error {
	if err := p.SpreadSplit.Validate(); err != nil {
		return fmt.Errorf("invalid spread split: %w", err)
	}
	if err := p.BridgeOutSplit.Validate(); err != nil {
		return fmt.Errorf("invalid bridge out split: %w", err)
	}
	return nil
}
