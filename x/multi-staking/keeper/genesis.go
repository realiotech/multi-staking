package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// multi-staking state
	for _, multiStakingLock := range data.MultiStakingLocks {
		// set staking lock
		lockID := types.MultiStakingLockID(sdk.AccAddress(multiStakingLock.DelAddr), sdk.ValAddress(multiStakingLock.ValAddr))
		k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
		// set intermediaryAccount
		intermediaryAccount := types.IntermediaryAccount(sdk.AccAddress(multiStakingLock.DelAddr))
		k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}

	return k.stakingKeeper.InitGenesis(ctx, data.StakingGenesisState)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		StakingGenesisState: k.stakingKeeper.ExportGenesis(ctx),
	}
}
