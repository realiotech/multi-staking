package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) CalUnbondedTokens(ctx sdk.Context, dvPair stakingtypes.DVPair) math.Int {
	addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

	// k.stakingKeeper(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, found bool)

	cachedCtx, _ := ctx.CacheContext()
	balances, err := k.CompleteUnbonding(cachedCtx, delegatorAddress, addr)
	if err != nil {

	} else {

	}

}
