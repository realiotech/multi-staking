package multistaking

import (
	"context"
	"encoding/json"
	"fmt"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/client/cli"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/simulation"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sdksimulation "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

var (
	_ module.HasABCIGenesis      = AppModule{}
	_ module.HasServices         = AppModule{}
	_ module.HasABCIEndBlock     = AppModule{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasInvariants       = AppModule{}
	_ appmodule.HasBeginBlocker  = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/staking AppModuleBasic.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return multistakingtypes.ModuleName
}

// RegisterLegacyAminoCodec register module codec
func (am AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	multistakingtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module interface
func (am AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	stakingtypes.RegisterInterfaces(reg)
	multistakingtypes.RegisterInterfaces(reg)
}

// DefaultGenesis returns multi-staking module default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(multistakingtypes.DefaultGenesis())
}

// ValidateGenesis validate genesis state for multi-staking module
func (am AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState multistakingtypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", multistakingtypes.ModuleName, err)
	}

	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := multistakingtypes.RegisterQueryHandlerClient(context.Background(), mux, multistakingtypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the staking module's root tx command.
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd(amb.cdc.InterfaceRegistry().SigningContext().ValidatorAddressCodec(), amb.cdc.InterfaceRegistry().SigningContext().AddressCodec())
}

// GetQueryCmd returns the multi-staking and staking module's root query command.
func (AppModuleBasic) GetQueryCmd() (queryCmd *cobra.Command) {
	return cli.GetQueryCmd()
}

// AppModule embeds the Cosmos SDK's x/staking AppModule where we only override
// specific methods.
type AppModule struct {
	AppModuleBasic
	// embed the Cosmos SDK's x/staking AppModule
	skAppModule staking.AppModule

	keeper keeper.Keeper
	sk     *stakingkeeper.Keeper
	ak     stakingtypes.AccountKeeper
	bk     stakingtypes.BankKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace exported.Subspace
}

// NewAppModule creates a new AppModule object using the native x/staking module
// AppModule constructor.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, sk *stakingkeeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, ss exported.Subspace) AppModule {
	stakingAppMod := staking.NewAppModule(cdc, sk, ak, bk, ss)
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		skAppModule:    stakingAppMod,
		keeper:         keeper,
		sk:             sk,
		ak:             ak,
		bk:             bk,
		legacySubspace: ss,
	}
}

// Name returns the staking module's name.
func (AppModule) Name() string {
	return multistakingtypes.ModuleName
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.skAppModule.RegisterInvariants(ir)
	keeper.RegisterInvariants(ir, am.keeper)
}

// QuerierRoute returns the staking module's querier route name.
func (AppModule) QuerierRoute() string {
	return multistakingtypes.QuerierRoute
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	multistakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))

	multistakingtypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
	querier := stakingkeeper.Querier{Keeper: am.sk}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	// migrate staking module
	m := keeper.NewMigrator(am.sk, am.legacySubspace)
	if err := cfg.RegisterMigration(multistakingtypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", stakingtypes.ModuleName, err))
	}
}

// InitGenesis initial genesis state for multi-staking module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState multistakingtypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	return am.keeper.InitGenesis(ctx, genesisState)
}

// ExportGenesis export multi-staking state as raw message for multi-staking module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the multi-staking module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.skAppModule.BeginBlock(ctx)
}

// EndBlock returns the end blocker for the multi-staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	// calculate the amount of coin
	matureUnbondingDelegations, err := am.keeper.GetMatureUnbondingDelegations(ctx)
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}
	// staking endblock
	valUpdates, err := am.skAppModule.EndBlock(ctx)
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}
	// update endblock multi-staking
	am.keeper.EndBlocker(ctx, matureUnbondingDelegations)

	return valUpdates, nil
}

func (am AppModule) StakingAppModule() staking.AppModule {
	return am.skAppModule
}

// ConsensusVersion return module consensus version
func (AppModule) ConsensusVersion() uint64 { return 2 }

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// GenerateGenesisState creates a randomized GenState of the staking module.
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for staking module's types
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	// sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the staking module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return sdksimulation.WeightedOperations{}
}
