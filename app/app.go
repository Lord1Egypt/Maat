package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	serverapi "github.com/cosmos/cosmos-sdk/server/api"
	apiConfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	maatmodule "github.com/Lord1Egypt/Maat/x/maat"
	maatkeeper "github.com/Lord1Egypt/Maat/x/maat/keeper"
	maattypes "github.com/Lord1Egypt/Maat/x/maat/types"

	oraclemodule "github.com/Lord1Egypt/Maat/x/oracle"
	oraclekeeper "github.com/Lord1Egypt/Maat/x/oracle/keeper"
	oracletypes "github.com/Lord1Egypt/Maat/x/oracle/types"

	marketmodule "github.com/Lord1Egypt/Maat/x/market"
	marketkeeper "github.com/Lord1Egypt/Maat/x/market/keeper"
	markettypes "github.com/Lord1Egypt/Maat/x/market/types"

	reservemodule "github.com/Lord1Egypt/Maat/x/reserve"
	reservekeeper "github.com/Lord1Egypt/Maat/x/reserve/keeper"
	reservetypes "github.com/Lord1Egypt/Maat/x/reserve/types"

	bridgemodule "github.com/Lord1Egypt/Maat/x/bridge"
	bridgekeeper "github.com/Lord1Egypt/Maat/x/bridge/keeper"
	bridgetypes "github.com/Lord1Egypt/Maat/x/bridge/types"

	treasurymodule "github.com/Lord1Egypt/Maat/x/treasury"
	treasurykeeper "github.com/Lord1Egypt/Maat/x/treasury/keeper"
	treasurytypes "github.com/Lord1Egypt/Maat/x/treasury/types"
)

const AppName = "MaatApp"

func init() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultNodeHome = filepath.Join(userHome, ".maat")
}

var (
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		gov.NewAppModuleBasic(nil),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		evidence.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		consensus.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		distr.AppModuleBasic{},
		maatmodule.AppModule{},
		oraclemodule.AppModule{},
		marketmodule.AppModule{},
		reservemodule.AppModule{},
		bridgemodule.AppModule{},
		treasurymodule.AppModule{},
	)

	// maccPerms defines module account permissions for token minting/burning.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		maattypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		bridgetypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
	}
)

type MaatApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	keys  map[string]*storetypes.KVStoreKey
	tkeys map[string]*storetypes.TransientStoreKey

	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.BaseKeeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        *govkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	ConsensusKeeper  consensuskeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper

	MaatKeeper     maatkeeper.Keeper
	OracleKeeper   oraclekeeper.Keeper
	MarketKeeper   marketkeeper.Keeper
	ReserveKeeper  reservekeeper.Keeper
	BridgeKeeper   bridgekeeper.Keeper
	TreasuryKeeper treasurykeeper.Keeper

	mm                 *module.Manager
	BasicModuleManager module.BasicManager
}

func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *MaatApp {
	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		slashingtypes.StoreKey,
		distrtypes.StoreKey,
		govtypes.StoreKey,
		crisistypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		paramstypes.StoreKey,
		consensustypes.StoreKey,
		capabilitytypes.StoreKey,
		maattypes.StoreKey,
		oracletypes.StoreKey,
		markettypes.StoreKey,
		reservetypes.StoreKey,
		bridgetypes.StoreKey,
		treasurytypes.StoreKey,
	)

	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)

	app := &MaatApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// x/consensus manages CometBFT consensus params; baseapp reads them via the param store.
	app.ConsensusKeeper = consensuskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensustypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.ProvideEventService(),
	)
	bApp.SetParamStore(app.ConsensusKeeper.ParamsStore)

	addrPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	addressCodec := address.NewBech32Codec(addrPrefix)
	validatorAddressCodec := address.NewBech32Codec(addrPrefix + "valoper")
	consensusAddressCodec := address.NewBech32Codec(addrPrefix + "valcons")

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addressCodec,
		addrPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		validatorAddressCodec,
		consensusAddressCodec,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Wire staking hooks so distribution + slashing react to validator/delegation
	// changes. Without this, slashing never records validator signing info and the
	// node hits "no validator signing info found" at FinalizeBlock.
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		5,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper.AddressCodec(),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		nil,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		DefaultNodeHome,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.StakingKeeper,
		app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	// Ma'at custom module keepers.
	app.TreasuryKeeper = treasurykeeper.NewKeeper()

	app.OracleKeeper = oraclekeeper.NewKeeper(
		keys[oracletypes.StoreKey],
		oracletypes.DefaultParams(),
		nil,
	)

	app.MarketKeeper = marketkeeper.NewKeeper(
		keys[markettypes.StoreKey],
		&app.OracleKeeper,
		&app.TreasuryKeeper,
		markettypes.DefaultParams(),
	)

	app.ReserveKeeper = reservekeeper.NewKeeper(
		reservetypes.StoreKey,
		&app.MarketKeeper,
	)

	app.BridgeKeeper = bridgekeeper.NewKeeper(
		app.BankKeeper,
		&app.MarketKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// x/maat: native token with BeginBlock minting
	app.MaatKeeper = maatkeeper.NewKeeper(
		keys[maattypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		bApp.MsgServiceRouter(),
		govtypes.DefaultConfig(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, nil),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, nil, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, nil),
		crisis.NewAppModule(app.CrisisKeeper, false, nil),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		maatmodule.NewAppModule(app.MaatKeeper),
		oraclemodule.NewAppModule(&app.OracleKeeper),
		marketmodule.NewAppModule(&app.MarketKeeper),
		reservemodule.NewAppModule(app.ReserveKeeper),
		bridgemodule.NewAppModule(&app.BridgeKeeper),
		treasurymodule.NewAppModule(&app.TreasuryKeeper),
	)

	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		distrtypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		consensustypes.ModuleName,
		maattypes.ModuleName,
		oracletypes.ModuleName,
		markettypes.ModuleName,
		reservetypes.ModuleName,
		bridgetypes.ModuleName,
		treasurytypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		distrtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		upgradetypes.ModuleName,
		paramstypes.ModuleName,
		consensustypes.ModuleName,
		maattypes.ModuleName,
		oracletypes.ModuleName,
		markettypes.ModuleName,
		reservetypes.ModuleName,
		bridgetypes.ModuleName,
		treasurytypes.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		consensustypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		genutiltypes.ModuleName,
		maattypes.ModuleName,
		oracletypes.ModuleName,
		markettypes.ModuleName,
		reservetypes.ModuleName,
		bridgetypes.ModuleName,
		treasurytypes.ModuleName,
	)

	// Basic manager derived from the constructed modules — its AppModuleBasics carry
	// real codecs, so CLI tx/query command builders (e.g. staking) don't nil-panic.
	// NOTE: interfaces/amino are already registered on the shared registry by
	// MakeEncodingConfig (via the global ModuleBasics); re-registering here would
	// double-register Msg types (e.g. MsgCreateValidator) and panic. This manager
	// is used only to build CLI tx/query commands with codec-bearing module basics.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName:     gov.NewAppModuleBasic(nil),
		},
	)

	app.mm.RegisterInvariants(app.CrisisKeeper)

	// Register all module Msg/Query services on the routers — without this, txs like
	// MsgCreateValidator (gentx) have "no message handler found" and InitChain panics.
	if err := app.mm.RegisterServices(module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())); err != nil {
		panic(err)
	}

	// PreBlock runs the upgrade module's scheduled-upgrade check before each block.
	app.mm.SetOrderPreBlockers(upgradetypes.ModuleName)

	if err := app.RegisterStores(keys, tkeys); err != nil {
		panic(err)
	}

	// ABCI lifecycle wiring: without these the node cannot init genesis or produce blocks.
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	anteHandler, err := ante.NewAnteHandler(ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		FeegrantKeeper:  nil,
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	})
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}
	}

	return app
}

// InitChainer initializes module genesis state and records module versions.
func (app *MaatApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}
	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		return nil, err
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// PreBlocker runs module PreBlock hooks (upgrade scheduling).
func (app *MaatApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

// BeginBlocker runs all module BeginBlock hooks (incl. x/maat block-reward mint).
func (app *MaatApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker runs all module EndBlock hooks (incl. x/oracle aggregation, x/market health).
func (app *MaatApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

func (app *MaatApp) RegisterStores(
	keys map[string]*storetypes.KVStoreKey,
	tkeys map[string]*storetypes.TransientStoreKey,
) error {
	// Mount the exact key objects the keepers hold — re-creating keys with the same
	// name would mount stores under different pointers, breaking ctx.KVStore() lookups.
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	return nil
}

func (app *MaatApp) Name() string { return AppName }

func (app *MaatApp) LegacyAmino() *codec.LegacyAmino { return app.legacyAmino }

func (app *MaatApp) AppCodec() codec.Codec { return app.appCodec }

func (app *MaatApp) ModuleManager() *module.Manager { return app.mm }

func BlockedAddresses() map[string]bool {
	blocked := make(map[string]bool)
	for acc := range maccPerms {
		blocked[authtypes.NewModuleAddress(acc).String()] = true
	}
	return blocked
}

func initParamsKeeper(
	appCodec codec.Codec,
	legacyAmino *codec.LegacyAmino,
	key *storetypes.KVStoreKey,
	tkey *storetypes.TransientStoreKey,
) paramskeeper.Keeper {
	pk := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	pk.Subspace(authtypes.ModuleName)
	pk.Subspace(banktypes.ModuleName)
	pk.Subspace(stakingtypes.ModuleName)
	pk.Subspace(distrtypes.ModuleName)
	pk.Subspace(slashingtypes.ModuleName)
	pk.Subspace(govtypes.ModuleName)
	pk.Subspace(crisistypes.ModuleName)

	return pk
}

func (app *MaatApp) RegisterAPIRoutes(apiSrv *serverapi.Server, apiCfg apiConfig.APIConfig) {}

func (app *MaatApp) RegisterNodeService(clientCtx client.Context, cfg apiConfig.Config) {}

func (app *MaatApp) RegisterTendermintService(clientCtx client.Context) {}

func (app *MaatApp) RegisterTxService(clientCtx client.Context) {}
