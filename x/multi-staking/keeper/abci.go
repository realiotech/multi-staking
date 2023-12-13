package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func GetUnbondingHeightsAndUnbondedAmounts(ctx sdk.Context, unbondingDelegation stakingtypes.UnbondingDelegation) map[int64]math.Int {
	ctxTime := ctx.BlockHeader().Time

	unbondingHeightsAndUnbondedAmounts := map[int64]math.Int{}
	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(unbondingDelegation.Entries); i++ {
		entry := unbondingDelegation.Entries[i]
		if entry.IsMature(ctxTime) && !entry.Balance.IsZero() {
			unbondedAmount, found := unbondingHeightsAndUnbondedAmounts[entry.CreationHeight]
			if !found {
				unbondingHeightsAndUnbondedAmounts[entry.CreationHeight] = entry.Balance
			} else {
				unbondedAmount = unbondedAmount.Add(entry.Balance)
				unbondingHeightsAndUnbondedAmounts[entry.CreationHeight] = unbondedAmount
			}
		}
	}
	return unbondingHeightsAndUnbondedAmounts
}

func (k Keeper) EndBlocker(ctx sdk.Context, matureUnbondingDelegations []stakingtypes.UnbondingDelegation) {
	for _, ubd := range matureUnbondingDelegations {
		delAcc, valAcc, err := types.DelAccAndValAccFromStrings(ubd.DelegatorAddress, ubd.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		unbondingHeightsAndUnbondedAmounts := GetUnbondingHeightsAndUnbondedAmounts(ctx, ubd)
		for unbondingHeight, unbondedAmount := range unbondingHeightsAndUnbondedAmounts {
			_, err = k.CompleteUnbonding(ctx, delAcc, valAcc, unbondingHeight, unbondedAmount)
			if err != nil {
				panic(err)
			}
		}
	}
}
