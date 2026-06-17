package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/x/market/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetParams(gs.Params)
	for _, dr := range gs.Reserves {
		r := &pricing.Reserve{
			AssetUnits:    dr.AssetUnits,
			WrappedSupply: dr.WrappedSupply,
			SpreadUSD:     dr.SpreadUSD,
			MinBackingBps: dr.MinBackingBps,
		}
		k.SetReserve(dr.Denom, r)
	}
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	denoms := k.AllDenoms()
	reserves := make([]types.DenomReserve, 0, len(denoms))
	for _, d := range denoms {
		r := k.GetReserve(d)
		if r == nil {
			continue
		}
		reserves = append(reserves, types.DenomReserve{
			Denom:         d,
			AssetUnits:    r.AssetUnits,
			WrappedSupply: r.WrappedSupply,
			SpreadUSD:     r.SpreadUSD,
			MinBackingBps: r.MinBackingBps,
		})
	}
	return &types.GenesisState{
		Params:   k.GetParams(),
		Reserves: reserves,
	}
}
