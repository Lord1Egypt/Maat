package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/x/bridge/keeper"
	"github.com/Lord1Egypt/Maat/x/bridge/types"
)

type mockBankKeeper struct {
	minted   sdk.Coins
	burned   sdk.Coins
	sentTo   sdk.Coins
	sentFrom sdk.Coins
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	m.minted = m.minted.Add(amt...)
	return nil
}

func (m *mockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	m.burned = m.burned.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	m.sentTo = m.sentTo.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	m.sentFrom = m.sentFrom.Add(amt...)
	return nil
}

type mockMarketKeeper struct {
	reserves map[string]*pricing.Reserve
}

func (m *mockMarketKeeper) GetReserve(denom string) *pricing.Reserve {
	return m.reserves[denom]
}

func (m *mockMarketKeeper) SetReserve(denom string, r *pricing.Reserve) {
	m.reserves[denom] = r
}

func TestBridgeKeeper(t *testing.T) {
	bk := &mockBankKeeper{}
	mk := &mockMarketKeeper{
		reserves: map[string]*pricing.Reserve{
			"weth": {
				AssetUnits:    100 * pricing.MicroUSD,
				WrappedSupply: 100 * pricing.MicroUSD,
				MinBackingBps: 10000,
			},
		},
	}
	k := keeper.NewKeeper(bk, mk, "authority")
	require.Equal(t, "authority", k.GetAuthority())

	lim := k.GetLimiter("weth")
	require.NotNil(t, lim)

	k.SetWithdrawer(123, "withdrawer1")
	addr, found := k.GetWithdrawer(123)
	require.True(t, found)
	require.Equal(t, "withdrawer1", addr)
}

func TestBridgeMsgServer(t *testing.T) {
	bk := &mockBankKeeper{}
	mk := &mockMarketKeeper{
		reserves: map[string]*pricing.Reserve{
			"weth": {
				AssetUnits:    100 * pricing.MicroUSD,
				WrappedSupply: 100 * pricing.MicroUSD,
				MinBackingBps: 10000,
			},
		},
	}
	k := keeper.NewKeeper(bk, mk, "authority")
	server := keeper.NewMsgServer(&k)
	ctx := sdk.Context{}.
		WithBlockHeight(100).
		WithEventManager(sdk.NewEventManager())

	// BridgeIn (verify Wormhole VAA mock, mints coins)
	depositor := sdk.AccAddress([]byte("addr1")).String()
	_, err := server.BridgeIn(ctx, &types.MsgBridgeIn{
		Depositor: depositor,
		Denom:     "weth",
		Amount:    10 * pricing.MicroUSD,
	})
	require.NoError(t, err)
	require.Equal(t, int64(10*pricing.MicroUSD), bk.minted.AmountOf("weth").Int64())

	// BridgeOut (Immediate because amount is under large tx delay threshold)
	withdrawer := sdk.AccAddress([]byte("addr2")).String()
	res, err := server.BridgeOut(ctx, &types.MsgBridgeOut{
		Withdrawer: withdrawer,
		Denom:      "weth",
		Amount:     1 * pricing.MicroUSD,
	})
	require.NoError(t, err)
	require.Equal(t, int32(2), res.Status)
}
