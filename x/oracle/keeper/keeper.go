package keeper

import (
	"sync"

	storetypes "cosmossdk.io/store/types"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/x/oracle/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	params   types.Params
	feeders  map[string]bool

	mu      sync.RWMutex
	oracles map[string]*pricing.Oracle
	votes   map[string][]int64
}

func NewKeeper(storeKey storetypes.StoreKey, params types.Params, feeders []string) Keeper {
	fm := make(map[string]bool, len(feeders))
	for _, f := range feeders {
		fm[f] = true
	}
	return Keeper{
		storeKey: storeKey,
		params:   params,
		feeders:  fm,
		oracles:  make(map[string]*pricing.Oracle),
		votes:    make(map[string][]int64),
	}
}

func (k *Keeper) SetParams(p types.Params) {
	k.params = p
}

func (k *Keeper) GetParams() types.Params {
	return k.params
}

func (k *Keeper) IsFeeder(addr string) bool {
	return k.feeders[addr]
}

func (k *Keeper) GetOracle(denom string) *pricing.Oracle {
	k.mu.RLock()
	o := k.oracles[denom]
	k.mu.RUnlock()
	return o
}

func (k *Keeper) getOrCreateOracle(denom string) *pricing.Oracle {
	k.mu.Lock()
	defer k.mu.Unlock()
	o, ok := k.oracles[denom]
	if !ok {
		o = pricing.NewOracle(k.params.AlphaNum, k.params.AlphaDen, k.params.MaxDevBps, k.params.BreakBps)
		k.oracles[denom] = o
	}
	return o
}

func (k *Keeper) GetMid(denom string) (int64, bool) {
	o := k.GetOracle(denom)
	if o == nil || o.Mid == 0 {
		return 0, false
	}
	return o.Mid, true
}

func (k *Keeper) SetMid(denom string, mid int64) {
	o := k.getOrCreateOracle(denom)
	o.Mid = mid
}

func (k *Keeper) IsHalted(denom string) bool {
	o := k.GetOracle(denom)
	if o == nil {
		return false
	}
	return o.Halted
}

func (k *Keeper) AppendVote(denom string, price int64) {
	k.mu.Lock()
	k.votes[denom] = append(k.votes[denom], price)
	k.mu.Unlock()
}

func (k *Keeper) GetVotes(denom string) []int64 {
	k.mu.RLock()
	v := k.votes[denom]
	k.mu.RUnlock()
	return v
}

func (k *Keeper) ClearVotes(denom string) {
	k.mu.Lock()
	delete(k.votes, denom)
	k.mu.Unlock()
}

func (k *Keeper) AllDenoms() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	denoms := make([]string, 0, len(k.oracles))
	for d := range k.oracles {
		denoms = append(denoms, d)
	}
	return denoms
}

func (k *Keeper) VoteDenoms() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	denoms := make([]string, 0, len(k.votes))
	for d := range k.votes {
		denoms = append(denoms, d)
	}
	return denoms
}
