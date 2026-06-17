package types

import "fmt"

type DenomMid struct {
	Denom  string `json:"denom"`
	Mid    int64  `json:"mid"`
	Halted bool   `json:"halted"`
}

type GenesisState struct {
	Params Params     `json:"params"`
	Mids   []DenomMid `json:"mids"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Mids:   nil,
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	seen := make(map[string]bool)
	for _, dm := range gs.Mids {
		if dm.Denom == "" {
			return fmt.Errorf("denom mid has empty denom")
		}
		if seen[dm.Denom] {
			return fmt.Errorf("duplicate denom mid: %s", dm.Denom)
		}
		seen[dm.Denom] = true
		if dm.Mid < 0 {
			return fmt.Errorf("negative mid for denom %s: %d", dm.Denom, dm.Mid)
		}
	}
	return nil
}
