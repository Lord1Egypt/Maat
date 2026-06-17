package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/x/treasury/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetParams(gs.Params)
	k.SetBalances(gs.Balances)
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(),
		Balances: k.GetBalances(),
	}
}
