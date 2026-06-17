package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/chain/token"
	"github.com/Lord1Egypt/Maat/x/maat/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) error {
	if err := gs.Validate(); err != nil {
		return err
	}

	if err := k.SetParams(gs.Params); err != nil {
		return err
	}

	allocs := token.Standard()
	for _, alloc := range allocs {
		tgeAmount := (alloc.Total * int64(alloc.TGEBps)) / 10_000

		if tgeAmount == 0 {
			continue
		}

		coins := sdk.NewCoins(sdk.NewInt64Coin(types.NativeDenom, tgeAmount))

		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return fmt.Errorf("mint for allocation %s: %w", alloc.Name, err)
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(_ sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(),
	}
}
