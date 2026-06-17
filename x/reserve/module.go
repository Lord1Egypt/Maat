package reserve

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/Lord1Egypt/Maat/x/reserve/keeper"
	"github.com/Lord1Egypt/Maat/x/reserve/types"
)

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasGenesis     = AppModule{}
)

type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{keeper: k}
}

func (AppModule) Name() string { return types.ModuleName }

func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

func (AppModule) RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}

func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := keeper.GenesisState{MarketHalted: false}
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs keeper.GenesisState
	if err := json.Unmarshal(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err)
	}
	return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs keeper.GenesisState
	if err := json.Unmarshal(bz, &gs); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	am.keeper.InitGenesis(gs)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis()
	bz, _ := json.Marshal(gs)
	return bz
}

func (am AppModule) EndBlock(ctx sdk.Context) error {
	am.keeper.EndBlockCheck()
	return nil
}

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}

func (am AppModule) ConsensusVersion() uint64 { return 1 }
