package types

import (
	"fmt"
	"github.com/Lord1Egypt/Maat/chain/bridge"
)

type DenomLimiter struct {
	Denom        string              `json:"denom"`
	WindowStart  int64               `json:"window_start"`
	UsedInWindow int64               `json:"used_in_window"`
	NextID       uint64              `json:"next_id"`
	Withdrawals  []bridge.Withdrawal `json:"withdrawals"`
}

type GenesisState struct {
	Params   Params         `json:"params"`
	Limiters []DenomLimiter `json:"limiters"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Limiters: nil,
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	seen := make(map[string]bool)
	for _, l := range gs.Limiters {
		if l.Denom == "" {
			return fmt.Errorf("limiter has empty denom")
		}
		if seen[l.Denom] {
			return fmt.Errorf("duplicate limiter denom: %s", l.Denom)
		}
		seen[l.Denom] = true
		if l.WindowStart < 0 {
			return fmt.Errorf("negative window start for denom %s: %d", l.Denom, l.WindowStart)
		}
		if l.UsedInWindow < 0 {
			return fmt.Errorf("negative used in window for denom %s: %d", l.Denom, l.UsedInWindow)
		}
	}
	return nil
}
