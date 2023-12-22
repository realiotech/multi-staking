package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// multi-staking state
	for _, multiStakingLock := range data.MultiStakingLocks {
		delAcc, valAcc, err := types.DelAccAndValAccFromStrings(multiStakingLock.LockID.DelAddr, multiStakingLock.LockID.ValAddr)
		if err != nil {

		}

		// set staking lock
		k.SetMultiStakingLock(ctx, types.MultiStakingLockID(delAcc, valAcc), multiStakingLock)
		// set intermediaryAccount
		// intermediaryAccount := types.IntermediaryAccount(sdk.AccAddress(multiStakingLock.DelAddr))
		// k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}
	// for _, multiStakingUnlock := range data.MultiStakingUnlocks {

	// }

	for _, valAllowedCoin := range data.ValidatorAllowedCoin {
		valAddr, err := sdk.ValAddressFromBech32(valAllowedCoin.ValAddr)
		if err != nil {
			panic("error validator address")
		}
		k.SetValidatorAllowedCoin(ctx, valAddr, valAllowedCoin.CoinDenom)
	}

	return k.stakingKeeper.InitGenesis(ctx, data.StakingGenesisState)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// get multiStakingLock
	var multiStakingLocks []types.MultiStakingLock
	k.MultiStakingLockIterator(ctx, func(stakingLock types.MultiStakingLock) bool {
		multiStakingLocks = append(multiStakingLocks, stakingLock)
		return false
	})

	// get validator allowed coin
	var validatorAllowedCoinLists []types.ValidatorAllowedCoin
	k.ValidatorAllowedCoinIterator(ctx, func(valAddr string, denom string) (stop bool) {
		validatorAllowedCoin := types.ValidatorAllowedCoin{
			ValAddr:   valAddr,
			CoinDenom: denom,
		}
		validatorAllowedCoinLists = append(validatorAllowedCoinLists, validatorAllowedCoin)
		return false
	})

	return &types.GenesisState{
		MultiStakingLocks:    multiStakingLocks,
		ValidatorAllowedCoin: validatorAllowedCoinLists,
		StakingGenesisState:  k.stakingKeeper.ExportGenesis(ctx),
	}
}
