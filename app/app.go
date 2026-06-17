package app

import (
	"io"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	serverapi "github.com/cosmos/cosmos-sdk/server/api"
	apiConfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
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
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
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
		bridgetypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
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
	EvidenceKeeper   evidencekeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper

	MaatKeeper     maatkeeper.Keeper
	OracleKeeper   oraclekeeper.Keeper
	MarketKeeper   marketkeeper.Keeper
	ReserveKeeper  reservekeeper.Keeper
	BridgeKeeper   bridgekeeper.Keeper
	TreasuryKeeper treasurykeeper.Keeper

	mm *module.Manager
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
		maattypes.ModuleName,
		oracletypes.ModuleName,
		markettypes.ModuleName,
		reservetypes.ModuleName,
		bridgetypes.ModuleName,
		treasurytypes.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
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

	app.mm.RegisterInvariants(app.CrisisKeeper)

	if err := app.RegisterStores(keys, tkeys); err != nil {
		panic(err)
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}
	}

	return app
}

func (app *MaatApp) RegisterStores(
	keys map[string]*storetypes.KVStoreKey,
	tkeys map[string]*storetypes.TransientStoreKey,
) error {
	for _, key := range keys {
		app.MountKVStores(storetypes.NewKVStoreKeys(key.Name()))
	}
	for _, tkey := range tkeys {
		app.MountTransientStores(storetypes.NewTransientStoreKeys(tkey.Name()))
	}
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
