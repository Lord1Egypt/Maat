package keeper

import (
	"sync"

	"github.com/Lord1Egypt/Maat/chain/treasury"
	"github.com/Lord1Egypt/Maat/x/treasury/types"
)

type Keeper struct {
	params types.Params

	mu   sync.RWMutex
	core treasury.Treasury
}

func NewKeeper() Keeper {
	return Keeper{
		params: types.DefaultParams(),
	}
}

func (k *Keeper) SetParams(p types.Params) {
	k.params = p
}

func (k *Keeper) GetParams() types.Params {
	return k.params
}

func (k *Keeper) Collect(amount int64, split treasury.Split) (treasury.Allocation, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.core.Collect(amount, split)
}

func (k *Keeper) GetBalances() types.TreasuryBalances {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return types.TreasuryBalances{
		Reserve:   k.core.Reserve,
		Rewards:   k.core.Rewards,
		Insurance: k.core.Insurance,
		Treasury:  k.core.Treasury,
	}
}

func (k *Keeper) SetBalances(b types.TreasuryBalances) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.core.Reserve = b.Reserve
	k.core.Rewards = b.Rewards
	k.core.Insurance = b.Insurance
	k.core.Treasury = b.Treasury
}

func ToCoreSplit(s types.SplitParams) treasury.Split {
	return treasury.Split{
		ReserveBps:   s.ReserveBps,
		RewardsBps:   s.RewardsBps,
		InsuranceBps: s.InsuranceBps,
		TreasuryBps:  s.TreasuryBps,
	}
}
