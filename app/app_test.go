package app_test

import (
	"encoding/json"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/app"
)

const testChainID = "maat-test-1"

func TestAppInitialization(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("application initialization panicked: %v", r)
		}
	}()

	maatApp := app.New(logger, db, nil, true, nil)
	require.NotNil(t, maatApp)
}

// TestAppInitChainPipeline drives InitChain with the default genesis and asserts that
// the entire genesis-init pipeline runs. This guards against regressions in store
// mounting, the InitChainer wiring, and the consensus param store — the bugs that
// "compile + construct fine" but break only when the node actually boots.
//
// With a default genesis (no gentx) the pipeline legitimately stops at staking's
// "validator set is empty" requirement. Reaching THAT specific error proves every
// module's InitGenesis executed against real stores; any wiring regression would
// instead surface as a store-not-found panic or a consensus-param error earlier.
// Full block production with a real validator is covered end-to-end by
// scripts/localnet.sh (see SETUP.md).
func TestAppInitChainPipeline(t *testing.T) {
	db := dbm.NewMemDB()
	maatApp := app.New(log.NewNopLogger(), db, nil, true, nil, baseapp.SetChainID(testChainID))
	require.NotNil(t, maatApp)

	genState := app.ModuleBasics.DefaultGenesis(maatApp.AppCodec())
	stateBytes, err := json.Marshal(genState)
	require.NoError(t, err)

	consensusParams := cmttypes.DefaultConsensusParams().ToProto()

	_, err = maatApp.InitChain(&abci.RequestInitChain{
		ChainId:         testChainID,
		ConsensusParams: &consensusParams,
		AppStateBytes:   stateBytes,
		InitialHeight:   1,
	})
	require.Error(t, err, "default genesis has no validators, so InitChain must reach the staking validator check")
	require.Contains(t, err.Error(), "validator set is empty",
		"InitChain should run the full pipeline and stop only at the staking validator requirement")
}
