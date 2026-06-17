package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) BeginBlock(ctx sdk.Context) error {
	if err := k.MintBlockReward(ctx); err != nil {
		k.Logger(ctx).Error("failed to mint block reward", "error", err)
		return err
	}
	return nil
}
