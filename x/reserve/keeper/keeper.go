package keeper

import (
	"github.com/Lord1Egypt/Maat/chain/pricing"
)

type MarketKeeper interface {
	GetReserve(denom string) *pricing.Reserve
	AllDenoms() []string
}

type Keeper struct {
	storeKey     string
	marketKeeper MarketKeeper
	halted       bool
}

func NewKeeper(storeKey string, mk MarketKeeper) Keeper {
	return Keeper{
		storeKey:     storeKey,
		marketKeeper: mk,
	}
}

func (k Keeper) GetBacking(denom string) (assetUnits, wrappedSupply, backingBps int64, healthy bool, found bool) {
	r := k.marketKeeper.GetReserve(denom)
	if r == nil {
		return 0, 0, 0, false, false
	}
	return r.AssetUnits, r.WrappedSupply, r.BackingBps(), r.Healthy(), true
}

func (k Keeper) Healthy(denom string) bool {
	r := k.marketKeeper.GetReserve(denom)
	if r == nil {
		return true
	}
	return r.Healthy()
}

func (k Keeper) AllHealthy() bool {
	for _, denom := range k.marketKeeper.AllDenoms() {
		if !k.Healthy(denom) {
			return false
		}
	}
	return true
}

func (k *Keeper) IsMarketHalted() bool { return k.halted }

func (k *Keeper) SetMarketHalted(halted bool) { k.halted = halted }

func (k *Keeper) EndBlockCheck() bool {
	allHealthy := k.AllHealthy()
	k.halted = !allHealthy
	return allHealthy
}
