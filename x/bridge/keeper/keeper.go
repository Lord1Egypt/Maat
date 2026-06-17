package keeper

import (
	"context"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/chain/bridge"
	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/x/bridge/types"
)

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type MarketKeeper interface {
	GetReserve(denom string) *pricing.Reserve
	SetReserve(denom string, r *pricing.Reserve)
}

type Keeper struct {
	bankKeeper   BankKeeper
	marketKeeper MarketKeeper
	authority    string
	params       types.Params

	mu          sync.RWMutex
	limiters    map[string]*bridge.Limiter
	withdrawers map[uint64]string
}

func NewKeeper(
	bankKeeper BankKeeper,
	marketKeeper MarketKeeper,
	authority string,
) Keeper {
	return Keeper{
		bankKeeper:   bankKeeper,
		marketKeeper: marketKeeper,
		authority:    authority,
		params:       types.DefaultParams(),
		limiters:     make(map[string]*bridge.Limiter),
		withdrawers:  make(map[uint64]string),
	}
}

func (k *Keeper) SetParams(p types.Params) {
	k.params = p
}

func (k *Keeper) GetParams() types.Params {
	return k.params
}

func (k *Keeper) GetAuthority() string {
	return k.authority
}

func (k *Keeper) GetLimiter(denom string) *bridge.Limiter {
	k.mu.Lock()
	defer k.mu.Unlock()
	l, ok := k.limiters[denom]
	if !ok {
		l = bridge.NewLimiter(k.params.DailyCapUnits, k.params.WindowBlocks, k.params.LargeTxUnits, k.params.DelayBlocks)
		k.limiters[denom] = l
	}
	return l
}

func (k *Keeper) AllDenoms() []string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	denoms := make([]string, 0, len(k.limiters))
	for d := range k.limiters {
		denoms = append(denoms, d)
	}
	return denoms
}

func (k *Keeper) SetWithdrawer(id uint64, addr string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.withdrawers == nil {
		k.withdrawers = make(map[uint64]string)
	}
	k.withdrawers[id] = addr
}

func (k *Keeper) GetWithdrawer(id uint64) (string, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	if k.withdrawers == nil {
		return "", false
	}
	addr, ok := k.withdrawers[id]
	return addr, ok
}
