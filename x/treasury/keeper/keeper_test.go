package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/chain/treasury"
	"github.com/Lord1Egypt/Maat/x/treasury/keeper"
	"github.com/Lord1Egypt/Maat/x/treasury/types"
)

func TestTreasuryKeeper(t *testing.T) {
	k := keeper.NewKeeper()
	params := types.DefaultParams()
	k.SetParams(params)
	require.Equal(t, params, k.GetParams())

	// Set/Get balances
	balances := types.TreasuryBalances{
		Reserve:   100,
		Rewards:   200,
		Insurance: 300,
		Treasury:  400,
	}
	k.SetBalances(balances)
	require.Equal(t, balances, k.GetBalances())

	// Collect with splits
	split := treasury.Split{
		ReserveBps:   4000,
		RewardsBps:   3000,
		InsuranceBps: 2000,
		TreasuryBps:  1000,
	}
	alloc, err := k.Collect(1000, split)
	require.NoError(t, err)
	require.Equal(t, int64(400), alloc.Reserve)
	require.Equal(t, int64(300), alloc.Rewards)
	require.Equal(t, int64(200), alloc.Insurance)
	require.Equal(t, int64(100), alloc.Treasury)

	// Verify balances updated correctly
	newBalances := k.GetBalances()
	require.Equal(t, int64(500), newBalances.Reserve)
	require.Equal(t, int64(500), newBalances.Rewards)
	require.Equal(t, int64(500), newBalances.Insurance)
	require.Equal(t, int64(500), newBalances.Treasury)
}
