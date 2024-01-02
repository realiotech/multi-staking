package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type StakingKeeper interface {
	DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair)
	BondDenom(ctx sdk.Context) (res string)
	InitGenesis(ctx sdk.Context, data *stakingtypes.GenesisState) (res []abci.ValidatorUpdate)
	ExportGenesis(ctx sdk.Context) *stakingtypes.GenesisState
}
