package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/x/oracle/keeper"
	"github.com/Lord1Egypt/Maat/x/oracle/types"
)

func mockContext() sdk.Context {
	return sdk.Context{}.
		WithLogger(log.NewNopLogger()).
		WithEventManager(sdk.NewEventManager())
}

func TestOracleKeeper(t *testing.T) {
	params := types.DefaultParams()
	feeders := []string{"feeder1", "feeder2"}
	k := keeper.NewKeeper(nil, params, feeders)

	require.Equal(t, params, k.GetParams())

	// Whitelisting
	require.True(t, k.IsFeeder("feeder1"))
	require.False(t, k.IsFeeder("nonfeeder"))

	// Params updating
	newParams := types.DefaultParams()
	newParams.MaxDevBps = 1000
	k.SetParams(newParams)
	require.Equal(t, newParams.MaxDevBps, k.GetParams().MaxDevBps)

	// Mid price getter/setter
	mid, found := k.GetMid("umaat")
	require.False(t, found)
	require.Zero(t, mid)

	k.SetMid("umaat", 12345)
	mid, found = k.GetMid("umaat")
	require.True(t, found)
	require.Equal(t, int64(12345), mid)

	// Votes append/clear/list
	k.AppendVote("umaat", 12000)
	k.AppendVote("umaat", 13000)
	require.Equal(t, []int64{12000, 13000}, k.GetVotes("umaat"))
	require.Contains(t, k.VoteDenoms(), "umaat")

	k.ClearVotes("umaat")
	require.Nil(t, k.GetVotes("umaat"))
}

func TestOracleMsgServer(t *testing.T) {
	params := types.DefaultParams()
	feeders := []string{"feeder1"}
	k := keeper.NewKeeper(nil, params, feeders)
	server := keeper.NewMsgServer(&k)
	ctx := mockContext()

	// Unauthorized feeder
	_, err := server.SubmitPrice(ctx, &types.MsgSubmitPrice{
		Feeder: "unauthorized",
		Denom:  "umaat",
		Price:  100,
	})
	require.Error(t, err)

	// Authorized feeder submit
	_, err = server.SubmitPrice(ctx, &types.MsgSubmitPrice{
		Feeder: "feeder1",
		Denom:  "umaat",
		Price:  1000000,
	})
	require.NoError(t, err)
	require.Equal(t, []int64{1000000}, k.GetVotes("umaat"))

	// EndBlock update
	err = k.EndBlock(ctx)
	require.NoError(t, err)
	require.Empty(t, k.GetVotes("umaat"))

	// Halted oracle handling
	k.SetMid("umaat", 1000000)
	k.AppendVote("umaat", 2000000) // jump of 100% will halt because breakbps default is 5000 (50%)
	err = k.EndBlock(ctx)
	require.NoError(t, err) // EndBlock handles halt gracefully (emits halt event, does not fail)
	require.True(t, k.IsHalted("umaat"))

	// Submitting price to halted oracle must fail
	_, err = server.SubmitPrice(ctx, &types.MsgSubmitPrice{
		Feeder: "feeder1",
		Denom:  "umaat",
		Price:  1000000,
	})
	require.Error(t, err)

	// Resume oracle
	_, err = server.ResumeOracle(ctx, &types.MsgResumeOracle{
		Denom: "umaat",
	})
	require.NoError(t, err)
	require.False(t, k.IsHalted("umaat"))
}
