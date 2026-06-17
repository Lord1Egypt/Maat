package types

import (
	"fmt"

	"github.com/Lord1Egypt/Maat/chain/token"
)

type GenesisState struct {
	Params Params `json:"params"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	allocs := token.Standard()
	sum := token.SumAllocations(allocs)
	if sum != token.TotalSupply {
		return fmt.Errorf("%w: sum %d != total %d", ErrSupplyMismatch, sum, token.TotalSupply)
	}

	return nil
}
