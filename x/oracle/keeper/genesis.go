package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/x/oracle/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetParams(gs.Params)
	for _, dm := range gs.Mids {
		o := k.getOrCreateOracle(dm.Denom)
		o.Mid = dm.Mid
		o.Halted = dm.Halted
	}
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	denoms := k.AllDenoms()
	mids := make([]types.DenomMid, 0, len(denoms))
	for _, d := range denoms {
		o := k.GetOracle(d)
		if o == nil {
			continue
		}
		mids = append(mids, types.DenomMid{
			Denom:  d,
			Mid:    o.Mid,
			Halted: o.Halted,
		})
	}
	return &types.GenesisState{
		Params: k.GetParams(),
		Mids:   mids,
	}
}
