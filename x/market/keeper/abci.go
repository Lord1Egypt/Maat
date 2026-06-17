package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) EndBlock(ctx sdk.Context) error {
	allHealthy := true
	denoms := k.AllDenoms()

	for _, denom := range denoms {
		r := k.GetReserve(denom)
		if r == nil {
			continue
		}
		if !r.Healthy() {
			allHealthy = false
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"market_backing_breach",
					sdk.NewAttribute("denom", denom),
					sdk.NewAttribute("backing_bps", fmt.Sprintf("%d", r.BackingBps())),
					sdk.NewAttribute("min_backing_bps", fmt.Sprintf("%d", r.MinBackingBps)),
				),
			)
		}
	}

	wasHalted := k.IsMarketHalted()
	if !allHealthy && !wasHalted {
		k.SetMarketHalted(true)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("market_halted"),
		)
	} else if allHealthy && wasHalted {
		k.SetMarketHalted(false)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("market_resumed"),
		)
	}

	return nil
}
