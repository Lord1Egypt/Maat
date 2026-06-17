package maat

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/Lord1Egypt/Maat/x/maat/keeper"
	"github.com/Lord1Egypt/Maat/x/maat/types"
)

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasABCIGenesis = AppModule{}
	_ appmodule.AppModule   = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{keeper: k}
}

func (AppModule) Name() string { return types.ModuleName }

func (AppModule) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (AppModule) RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}

func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := types.DefaultGenesis()
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := json.Unmarshal(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return gs.Validate()
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		panic(fmt.Sprintf("failed to unmarshal genesis for %s: %v", types.ModuleName, err))
	}
	if err := am.keeper.InitGenesis(ctx, gs); err != nil {
		panic(fmt.Sprintf("failed to init genesis for %s: %v", types.ModuleName, err))
	}
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	bz, _ := json.Marshal(gs)
	return bz
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return am.keeper.BeginBlock(sdkCtx)
}

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}

func (AppModule) ConsensusVersion() uint64 { return 1 }
