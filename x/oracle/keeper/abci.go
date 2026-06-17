package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) EndBlock(ctx sdk.Context) error {
	denoms := k.VoteDenoms()
	for _, denom := range denoms {
		votes := k.GetVotes(denom)
		if len(votes) == 0 {
			continue
		}

		o := k.getOrCreateOracle(denom)
		newMid, err := o.Update(votes)
		k.ClearVotes(denom)

		if err != nil {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"oracle_halt",
					sdk.NewAttribute("denom", denom),
					sdk.NewAttribute("reason", err.Error()),
					sdk.NewAttribute("last_mid", fmt.Sprintf("%d", o.Mid)),
				),
			)
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"oracle_update",
				sdk.NewAttribute("denom", denom),
				sdk.NewAttribute("mid", fmt.Sprintf("%d", newMid)),
				sdk.NewAttribute("votes", fmt.Sprintf("%d", len(votes))),
			),
		)
	}
	return nil
}
