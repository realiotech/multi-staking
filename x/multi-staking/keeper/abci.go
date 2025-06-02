package keeper

import (
	"context"
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// need a way to better name this func
func GetUnbondingHeightsAndUnbondedAmounts(ctx context.Context, unbondingDelegation stakingtypes.UnbondingDelegation) map[int64]math.Int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ctxTime := sdkCtx.BlockHeader().Time

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

func (k Keeper) EndBlocker(ctx context.Context, matureUnbondingDelegations []stakingtypes.UnbondingDelegation) {
	for _, unbond := range matureUnbondingDelegations {
		multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(unbond.DelegatorAddress, unbond.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		unbondingHeightsAndUnbondedAmounts := GetUnbondingHeightsAndUnbondedAmounts(ctx, unbond)
		for unbondingHeight, unbondedAmount := range unbondingHeightsAndUnbondedAmounts {
			_, err := k.BurnUnbondedCoinAndUnlockedMultiStakingCoin(ctx, multiStakerAddr, valAcc, unbondingHeight, unbondedAmount)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (k Keeper) BurnUnbondedCoinAndUnlockedMultiStakingCoin(
	ctx context.Context,
	multiStakerAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	unbondingHeight int64,
	unbondAmount math.Int,
) (unlockedCoin sdk.Coin, err error) {
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
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	burnCoin := sdk.NewCoin(bondDenom, unbondAmount)
	err = k.BurnCoin(ctx, multiStakerAddr, burnCoin)
	if err != nil {
		return sdk.Coin{}, err
	}
	// burn remaining coin in unlock
	remaningCoin := unlockEntry.UnlockingCoin.ToCoin().Sub(unlockedCoin)
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(remaningCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	err = k.UnescrowCoinTo(ctx, multiStakerAddr, unlockedCoin)
	if err != nil {
		return sdk.Coin{}, err
	}

	err = k.DeleteUnlockEntryAtCreationHeight(ctx, unlockID, unbondingHeight)
	if err != nil {
		return sdk.Coin{}, err
	}

	return unlockedCoin, nil
}
