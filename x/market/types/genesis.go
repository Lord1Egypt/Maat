package types

import "fmt"

type DenomReserve struct {
	Denom         string `json:"denom"`
	AssetUnits    int64  `json:"asset_units"`
	WrappedSupply int64  `json:"wrapped_supply"`
	SpreadUSD     int64  `json:"spread_usd"`
	MinBackingBps int64  `json:"min_backing_bps"`
}

type GenesisState struct {
	Params   Params         `json:"params"`
	Reserves []DenomReserve `json:"reserves"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Reserves: nil,
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	seen := make(map[string]bool)
	for _, r := range gs.Reserves {
		if r.Denom == "" {
			return fmt.Errorf("reserve has empty denom")
		}
		if seen[r.Denom] {
			return fmt.Errorf("duplicate reserve denom: %s", r.Denom)
		}
		seen[r.Denom] = true
		if r.AssetUnits < 0 {
			return fmt.Errorf("negative asset units for denom %s: %d", r.Denom, r.AssetUnits)
		}
		if r.WrappedSupply < 0 {
			return fmt.Errorf("negative wrapped supply for denom %s: %d", r.Denom, r.WrappedSupply)
		}
		if r.MinBackingBps <= 0 {
			return fmt.Errorf("min backing bps must be positive for denom %s: %d", r.Denom, r.MinBackingBps)
		}
	}
	return nil
}
