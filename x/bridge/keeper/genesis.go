package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/chain/bridge"
	"github.com/Lord1Egypt/Maat/x/bridge/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetParams(gs.Params)
	for _, dl := range gs.Limiters {
		l := k.GetLimiter(dl.Denom)
		l.SetWindowStart(dl.WindowStart)
		l.SetUsedInWindow(dl.UsedInWindow)
		l.SetNextID(dl.NextID)

		pending := make(map[uint64]*bridge.Withdrawal)
		for _, w := range dl.Withdrawals {
			wCopy := w
			pending[w.ID] = &wCopy
		}
		l.SetPending(pending)
	}
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	denoms := k.AllDenoms()
	limiters := make([]types.DenomLimiter, 0, len(denoms))
	for _, d := range denoms {
		l := k.GetLimiter(d)
		if l == nil {
			continue
		}
		var withdrawals []bridge.Withdrawal
		for _, w := range l.GetPending() {
			withdrawals = append(withdrawals, *w)
		}
		limiters = append(limiters, types.DenomLimiter{
			Denom:        d,
			WindowStart:  l.GetWindowStart(),
			UsedInWindow: l.GetUsedInWindow(),
			NextID:       l.GetNextID(),
			Withdrawals:  withdrawals,
		})
	}
	return &types.GenesisState{
		Params:   k.GetParams(),
		Limiters: limiters,
	}
}
