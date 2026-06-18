package cmd

import (
	"io"
	"os"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"

	"github.com/Lord1Egypt/Maat/app"
)

// NewRootCmd builds the root maatd command with the full client + server command set.
func NewRootCmd() *cobra.Command {
	encodingConfig := app.MakeEncodingConfig()

	// Pre-instantiate the app to obtain a BasicModuleManager whose AppModuleBasics
	// carry real codecs (needed by staking/gov CLI command builders).
	tempApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, false, simtestutil.NewAppOptionsWithFlagHome(tempDir()))
	basicManager := tempApp.BasicModuleManager

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("") // empty prefix: no env-var prefix for Ma'at

	rootCmd := &cobra.Command{
		Use:   "maatd",
		Short: "Ma'at blockchain daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, initCometBFTConfig())
		},
	}

	initRootCmd(rootCmd, encodingConfig, basicManager)

	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfig, basicManager module.BasicManager) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		genutilcli.Commands(encodingConfig.TxConfig, basicManager, app.DefaultNodeHome),
		server.StatusCommand(),
		queryCommand(basicManager),
		txCommand(basicManager),
		keys.Commands(),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExporter, addModuleInitFlags)
}

func queryCommand(basicManager module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	basicManager.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand(basicManager module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	basicManager.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	// DefaultBaseappOptions wires pruning, min-gas-prices, snapshots, and crucially
	// SetChainID (with genesis fallback) so InitChain accepts the genesis chain-id.
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	return app.New(logger, db, traceStore, true, appOpts, baseappOptions...)
}

func appExporter(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	return servertypes.ExportedApp{}, nil
}

func addModuleInitFlags(startCmd *cobra.Command) {}

// initCometBFTConfig returns the default CometBFT config for the Ma'at node.
func initCometBFTConfig() *cmtcfg.Config {
	return cmtcfg.DefaultConfig()
}

// initAppConfig returns the default app.toml template and config, with a zero
// minimum gas price default suitable for a fresh local/testnet node.
func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0umaat"
	return serverconfig.DefaultConfigTemplate, srvCfg
}

// tempDir returns a throwaway home dir for the pre-instantiated app used only to
// derive the basic module manager.
func tempDir() string {
	dir, err := os.MkdirTemp("", "maatd")
	if err != nil {
		dir = app.DefaultNodeHome
	}
	defer os.RemoveAll(dir)
	return dir
}
