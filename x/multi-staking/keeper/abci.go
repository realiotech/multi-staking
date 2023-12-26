package keeper

import (
	"fmt"

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
	for _, unlock := range matureUnbondingDelegations {
		intermediaryAcc, valAcc, err := types.AccAddrAndValAddrFromStrings(unlock.DelegatorAddress, unlock.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		unbondingHeightsAndUnbondedAmounts := GetUnbondingHeightsAndUnbondedAmounts(ctx, unlock)
		for unbondingHeight, unbondedAmount := range unbondingHeightsAndUnbondedAmounts {
			k.BurnUnbondedCoinAndUnlockedMultiStakingCoin(ctx, intermediaryAcc, valAcc, unbondingHeight, unbondedAmount)
		}
	}
}

func (k Keeper) BurnUnbondedCoinAndUnlockedMultiStakingCoin(
	ctx sdk.Context,
	intermediaryAcc sdk.AccAddress,
	valAddr sdk.ValAddress,
	unbondingHeight int64,
	unbondAmount math.Int,
) (unlockedCoin sdk.Coin, err error) {
	// get multiStakerAddr
	multiStakerAddr := types.MultiStakerAddress(intermediaryAcc)

	// get unlock record
	unlockID := types.MultiStakingUnlockID(multiStakerAddr.String(), valAddr.String())
	unlockEntry, found := k.GetUnlockEntryAtCreationHeight(ctx, unlockID, unbondingHeight)
	if !found {
		return sdk.Coin{}, fmt.Errorf("unlock entry not found")
	}

	unlockDenom := unlockEntry.UnlockingCoin.Denom
	unlockedAmount := unlockEntry.UnbondAmountToUnlockAmount(unbondAmount)
	unlockedCoin = sdk.NewCoin(unlockDenom, unlockedAmount)

	// check amount
	if unlockedAmount.GT(unlockEntry.UnlockingCoin.Amount) {
		return sdk.Coin{}, fmt.Errorf("unlock amount greater than lock amount")
	}

	// burn bonded coin
	burnCoin := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondAmount))
	k.BurnCoin(ctx, intermediaryAcc, burnCoin)

	err = k.bankKeeper.SendCoins(ctx, intermediaryAcc, multiStakerAddr, sdk.NewCoins(unlockedCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	err = k.DeleteUnlockEntryAtCreationHeight(ctx, unlockID, unbondingHeight)
	if err != nil {
		return sdk.Coin{}, err
	}

	return unlockedCoin, nil
}
