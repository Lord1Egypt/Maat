package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Lord1Egypt/Maat/x/maat/keeper"
	"github.com/Lord1Egypt/Maat/x/maat/types"
)

type mockBankKeeper struct {
	minted sdk.Coins
	sent   map[string]sdk.Coins
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error {
	m.minted = m.minted.Add(amounts...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	if m.sent == nil {
		m.sent = make(map[string]sdk.Coins)
	}
	m.sent[senderModule+"->"+recipientModule] = m.sent[senderModule+"->"+recipientModule].Add(amt...)
	return nil
}

func (m *mockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewInt64Coin(denom, 0)
}

type mockAccountKeeper struct{}

func (m *mockAccountKeeper) GetModuleAddress(name string) sdk.AccAddress {
	return sdk.AccAddress([]byte(name))
}

func (m *mockAccountKeeper) GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI {
	return authtypes.NewEmptyModuleAccount(name)
}

func (m *mockAccountKeeper) SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI) {}

func TestMaatKeeper(t *testing.T) {
	bk := &mockBankKeeper{}
	ak := &mockAccountKeeper{}
	k := keeper.NewKeeper(nil, bk, ak, "authority")
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())

	// Params
	require.Equal(t, types.DefaultParams(), k.GetParams())
	params := types.DefaultParams()
	params.AnnualInflationBps = 500
	err := k.SetParams(params)
	require.NoError(t, err)
	require.Equal(t, params, k.GetParams())

	// Circulating Supply
	supply, err := k.GetCirculatingSupply(12)
	require.NoError(t, err)
	require.Positive(t, supply)

	// Block reward minting
	err = k.MintBlockReward(ctx)
	require.NoError(t, err)
	require.Positive(t, bk.minted.AmountOf(types.NativeDenom).Int64())
	key := types.ModuleName + "->" + authtypes.FeeCollectorName
	require.Positive(t, bk.sent[key].AmountOf(types.NativeDenom).Int64())
}
