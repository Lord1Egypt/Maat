package types

import (
	"fmt"
)

type TreasuryBalances struct {
	Reserve   int64 `json:"reserve"`
	Rewards   int64 `json:"rewards"`
	Insurance int64 `json:"insurance"`
	Treasury  int64 `json:"treasury"`
}

type GenesisState struct {
	Params   Params           `json:"params"`
	Balances TreasuryBalances `json:"balances"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Balances: TreasuryBalances{},
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	if gs.Balances.Reserve < 0 || gs.Balances.Rewards < 0 || gs.Balances.Insurance < 0 || gs.Balances.Treasury < 0 {
		return fmt.Errorf("balances cannot be negative")
	}
	return nil
}
