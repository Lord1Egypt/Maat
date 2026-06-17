package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/x/reserve/keeper"
)

type mockMarketKeeper struct {
	reserves map[string]*pricing.Reserve
}

func (m *mockMarketKeeper) GetReserve(denom string) *pricing.Reserve {
	return m.reserves[denom]
}

func (m *mockMarketKeeper) AllDenoms() []string {
	denoms := make([]string, 0, len(m.reserves))
	for d := range m.reserves {
		denoms = append(denoms, d)
	}
	return denoms
}

func TestReserveKeeper(t *testing.T) {
	mk := &mockMarketKeeper{
		reserves: map[string]*pricing.Reserve{
			"weth": {
				AssetUnits:    100 * pricing.MicroUSD,
				WrappedSupply: 100 * pricing.MicroUSD,
				MinBackingBps: 10000,
			},
		},
	}
	k := keeper.NewKeeper("storeKey", mk)

	// GetBacking
	asset, supply, backing, healthy, found := k.GetBacking("weth")
	require.True(t, found)
	require.True(t, healthy)
	require.Equal(t, int64(100*pricing.MicroUSD), asset)
	require.Equal(t, int64(100*pricing.MicroUSD), supply)
	require.Equal(t, int64(10000), backing)

	// GetBacking for unknown denom
	_, _, _, _, found = k.GetBacking("unknown")
	require.False(t, found)

	// Market halted
	require.False(t, k.IsMarketHalted())
	k.SetMarketHalted(true)
	require.True(t, k.IsMarketHalted())

	// EndBlockCheck healthy
	k.SetMarketHalted(false)
	require.True(t, k.EndBlockCheck())
	require.False(t, k.IsMarketHalted())

	// Unhealthy reserve
	mk.reserves["weth"].AssetUnits = 50 * pricing.MicroUSD // backing drops to 50%
	require.False(t, k.EndBlockCheck())
	require.True(t, k.IsMarketHalted())
}
