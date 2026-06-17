package keeper

type GenesisState struct {
	MarketHalted bool
}

func (k *Keeper) InitGenesis(state GenesisState) {
	k.halted = state.MarketHalted
}

func (k Keeper) ExportGenesis() GenesisState {
	return GenesisState{MarketHalted: k.halted}
}
