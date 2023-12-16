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
		// set intermediaryAccount
		intermediaryAccount := types.IntermediaryAccount(sdk.AccAddress(multiStakingLock.DelAddr))
		k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}

	for _, valAllowedToken := range data.ValidatorAllowedToken {
		valAddr, err := sdk.ValAddressFromBech32(valAllowedToken.ValAddr)
		if err != nil {
			panic("error validator address")
		}
		k.SetValidatorAllowedToken(ctx, valAddr, valAllowedToken.TokenDenom)
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

	// get validator allowed token
	var validatorAllowedTokenLists []types.ValidatorAllowedToken
	k.ValidatorAllowedTokenIterator(ctx, func(valAddr string, denom string) (stop bool) {
		validatorAllowedToken := types.ValidatorAllowedToken{
			ValAddr:    valAddr,
			TokenDenom: denom,
		}
		validatorAllowedTokenLists = append(validatorAllowedTokenLists, validatorAllowedToken)
		return false
	})

	return &types.GenesisState{
		MultiStakingLocks:     multiStakingLocks,
		ValidatorAllowedToken: validatorAllowedTokenLists,
		StakingGenesisState:   k.stakingKeeper.ExportGenesis(ctx),
	}
}
