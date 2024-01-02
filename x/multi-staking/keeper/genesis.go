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
		k.SetMultiStakingLock(ctx, multiStakingLock)
		// set intermediaryDelegator
		// intermediaryDelegator := types.IntermediaryDelegator(sdk.AccAddress(multiStakingLock.DelAddr))
		// k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}
	// for _, multiStakingUnlock := range data.MultiStakingUnlocks {

	// }

	for _, valMultiStakingCoin := range data.ValidatorMultiStakingCoins {
		valAddr, err := sdk.ValAddressFromBech32(valMultiStakingCoin.ValAddr)
		if err != nil {
			panic("error validator address")
		}
		k.SetValidatorMultiStakingCoin(ctx, valAddr, valMultiStakingCoin.CoinDenom)
	}

	return k.stakingKeeper.InitGenesis(ctx, &data.StakingGenesisState)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// get multiStakingLock
	var multiStakingLocks []types.MultiStakingLock
	k.MultiStakingLockIterator(ctx, func(stakingLock types.MultiStakingLock) bool {
		multiStakingLocks = append(multiStakingLocks, stakingLock)
		return false
	})

	// get validator allowed coin
	var ValidatorMultiStakingCoinLists []types.ValidatorMultiStakingCoin
	k.ValidatorMultiStakingCoinIterator(ctx, func(valAddr string, denom string) (stop bool) {
		ValidatorMultiStakingCoin := types.ValidatorMultiStakingCoin{
			ValAddr:   valAddr,
			CoinDenom: denom,
		}
		ValidatorMultiStakingCoinLists = append(ValidatorMultiStakingCoinLists, ValidatorMultiStakingCoin)
		return false
	})

	return &types.GenesisState{
		MultiStakingLocks:          multiStakingLocks,
		ValidatorMultiStakingCoins: ValidatorMultiStakingCoinLists,
		StakingGenesisState:        *k.stakingKeeper.ExportGenesis(ctx),
	}
}
