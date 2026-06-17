package keeper

import (
	"sync"

	storetypes "cosmossdk.io/store/types"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/chain/treasury"
	"github.com/Lord1Egypt/Maat/x/market/types"
)

type OracleKeeper interface {
	GetMid(denom string) (int64, bool)
	IsHalted(denom string) bool
}

type TreasuryKeeper interface {
	Collect(amount int64, split treasury.Split) (treasury.Allocation, error)
}

type Keeper struct {
	storeKey       storetypes.StoreKey
	oracleKeeper   OracleKeeper
	treasuryKeeper TreasuryKeeper
	params         types.Params
	marketHalted   bool

	mu       sync.RWMutex
	reserves map[string]*pricing.Reserve
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	oracleKeeper OracleKeeper,
	treasuryKeeper TreasuryKeeper,
	params types.Params,
) Keeper {
	return Keeper{
		storeKey:       storeKey,
		oracleKeeper:   oracleKeeper,
		treasuryKeeper: treasuryKeeper,
		params:         params,
		reserves:       make(map[string]*pricing.Reserve),
	}
}

func (k *Keeper) SetParams(p types.Params) {
	k.params = p
}

func (k *Keeper) GetParams() types.Params {
	return k.params
}

func (k *Keeper) IsMarketHalted() bool {
	return k.marketHalted
}

func (k *Keeper) SetMarketHalted(halted bool) {
	k.marketHalted = halted
}

func (k *Keeper) GetReserve(denom string) *pricing.Reserve {
	k.mu.RLock()
	r := k.reserves[denom]
	k.mu.RUnlock()
	return r
}

func (k *Keeper) SetReserve(denom string, r *pricing.Reserve) {
	k.mu.Lock()
	k.reserves[denom] = r
	k.mu.Unlock()
}

func (k *Keeper) SpreadParams() pricing.SpreadParams {
	return pricing.SpreadParams{
		BaseBps:    k.params.BaseBps,
		VolMultBps: k.params.VolMultBps,
		SkewMaxBps: k.params.SkewMaxBps,
		MaxBps:     k.params.MaxBps,
	}
}

func (k *Keeper) GetQuote(denom string, volSignalBps, skewBps int64) (pricing.Quote, error) {
	mid, ok := k.oracleKeeper.GetMid(denom)
	if !ok {
		return pricing.Quote{}, types.ErrNoOraclePrice.Wrapf("denom %s", denom)
	}
	sp := k.SpreadParams()
	effSpread := sp.EffectiveSpreadBps(volSignalBps)
	return pricing.MakeQuote(mid, effSpread, skewBps), nil
}

func (k *Keeper) Healthy(denom string) bool {
	r := k.GetReserve(denom)
	if r == nil {
		return false
	}
	return r.Healthy()
}

func (k *Keeper) AllDenoms() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	denoms := make([]string, 0, len(k.reserves))
	for d := range k.reserves {
		denoms = append(denoms, d)
	}
	return denoms
}
