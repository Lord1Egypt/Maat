package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/chain/treasury"
	"github.com/Lord1Egypt/Maat/x/market/keeper"
	"github.com/Lord1Egypt/Maat/x/market/types"
)

type mockOracleKeeper struct {
	mids   map[string]int64
	halted map[string]bool
}

func (m *mockOracleKeeper) GetMid(denom string) (int64, bool) {
	v, ok := m.mids[denom]
	return v, ok
}

func (m *mockOracleKeeper) IsHalted(denom string) bool {
	return m.halted[denom]
}

type mockTreasuryKeeper struct {
	collected int64
	splits    []treasury.Split
}

func (m *mockTreasuryKeeper) Collect(amount int64, split treasury.Split) (treasury.Allocation, error) {
	m.collected += amount
	m.splits = append(m.splits, split)
	return treasury.Allocation{}, nil
}

func TestMarketKeeper(t *testing.T) {
	ok := &mockOracleKeeper{
		mids:   map[string]int64{"weth": 3000 * pricing.MicroUSD},
		halted: map[string]bool{"weth": false},
	}
	tk := &mockTreasuryKeeper{}
	params := types.DefaultParams()
	k := keeper.NewKeeper(nil, ok, tk, params)

	require.Equal(t, params, k.GetParams())

	// Set/get reserve
	r := &pricing.Reserve{
		AssetUnits:    10 * pricing.MicroUSD,
		WrappedSupply: 10 * pricing.MicroUSD,
		MinBackingBps: 10000,
	}
	k.SetReserve("weth", r)
	require.Equal(t, r, k.GetReserve("weth"))
	require.True(t, k.Healthy("weth"))

	// GetQuote
	q, err := k.GetQuote("weth", 0, 0)
	require.NoError(t, err)
	require.Equal(t, int64(3000*pricing.MicroUSD), q.Mid)
	// Base spread is 15 bps (0.15%)
	// Sell = mid + mid*0.0015 = 3004.5
	require.Equal(t, int64(3004500000), q.Sell)

	// Market Halted
	require.False(t, k.IsMarketHalted())
	k.SetMarketHalted(true)
	require.True(t, k.IsMarketHalted())
}

func TestMarketMsgServer(t *testing.T) {
	ok := &mockOracleKeeper{
		mids:   map[string]int64{"weth": 3000 * pricing.MicroUSD},
		halted: map[string]bool{"weth": false},
	}
	tk := &mockTreasuryKeeper{}
	params := types.DefaultParams()
	k := keeper.NewKeeper(nil, ok, tk, params)
	k.SetReserve("weth", &pricing.Reserve{
		AssetUnits:    1000 * pricing.MicroUSD,
		WrappedSupply: 100 * pricing.MicroUSD,
		MinBackingBps: 10000,
	})
	server := keeper.NewMsgServer(&k)
	ctx := context.Background()

	// Swap buy (user buys weth, supply increases)
	res, err := server.Swap(ctx, &types.MsgSwap{
		FromDenom: "umaat",
		ToDenom:   "weth",
		Amount:    10 * pricing.MicroUSD,
	})
	require.NoError(t, err)
	require.Equal(t, int64(10*pricing.MicroUSD), res.Executed)
	require.Positive(t, res.SpreadUSD)

	// Swap sell (user sells weth, supply decreases)
	res, err = server.Swap(ctx, &types.MsgSwap{
		FromDenom: "weth",
		ToDenom:   "umaat",
		Amount:    5 * pricing.MicroUSD,
	})
	require.NoError(t, err)
	require.Equal(t, int64(5*pricing.MicroUSD), res.Executed)
	require.Positive(t, res.SpreadUSD)

	// Oracle Halted
	ok.halted["weth"] = true
	_, err = server.Swap(ctx, &types.MsgSwap{
		FromDenom: "umaat",
		ToDenom:   "weth",
		Amount:    10 * pricing.MicroUSD,
	})
	require.Error(t, err)
}
