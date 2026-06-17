package keeper

import (
	"context"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Lord1Egypt/Maat/chain/token"
	"github.com/Lord1Egypt/Maat/x/maat/types"
)

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
}

type Keeper struct {
	storeKey      storetypes.StoreKey
	bankKeeper    BankKeeper
	accountKeeper AccountKeeper
	authority     string
	params        types.Params
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	bankKeeper BankKeeper,
	accountKeeper AccountKeeper,
	authority string,
) Keeper {
	return Keeper{
		storeKey:      storeKey,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		authority:     authority,
		params:        types.DefaultParams(),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetParams() types.Params {
	return k.params
}

func (k *Keeper) SetParams(p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}
	k.params = p
	return nil
}

func (k Keeper) GetCirculatingSupply(month int64) (int64, error) {
	allocs := token.Standard()
	return token.CirculatingSupply(allocs, int(month))
}

func (k Keeper) MintBlockReward(ctx sdk.Context) error {
	p := k.params
	supply := int64(token.TotalSupply)

	reward, err := token.BlockReward(supply, p.AnnualInflationBps, p.BlocksPerYear)
	if err != nil {
		return err
	}

	if reward == 0 {
		return nil
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(types.NativeDenom, reward))

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		authtypes.FeeCollectorName,
		coins,
	)
}
