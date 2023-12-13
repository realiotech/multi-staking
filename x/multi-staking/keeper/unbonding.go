package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) CompleteUnbonding(
	ctx sdk.Context,
	intermediaryAcc sdk.AccAddress,
	valAddr sdk.ValAddress,
	unbondingHeight int64,
	balance math.Int,
) (unlockedAmount sdk.Coins, err error) {
	// get delAddr
	delAddr := types.DelegatorAccount(intermediaryAcc)

	// get unbonded record
	ubd, found := k.GetMultiStakingUnlock(ctx, delAddr, valAddr)
	if !found {
		return unlockedAmount, types.ErrRecordNotExists
	}
	var (
		unbondEntry      types.UnlockEntry
		unbondEntryIndex int64 = -1
	)

	for i, entry := range ubd.Entries {
		if entry.CreationHeight == unbondingHeight {
			unbondEntry = entry
			unbondEntryIndex = int64(i)
			break
		}
	}
	if unbondEntryIndex == -1 {
		return nil, sdkerrors.ErrNotFound.Wrapf("unbonding delegation entry is not found at block height %d", unbondingHeight)
	}

	unlockDenom := k.GetValidatorAllowedToken(ctx, valAddr)
	unlockMultiStakingAmount := sdk.NewDecFromInt(balance).Quo(unbondEntry.ConversionRatio).RoundInt()

	// check amount
	if unlockMultiStakingAmount.GT(unbondEntry.Balance) {
		return unlockedAmount, types.ErrCheckInsufficientAmount
	}

	// burn bonded token
	burnAmount := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), balance))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, intermediaryAcc, types.ModuleName, burnAmount)
	if err != nil {
		return unlockedAmount, err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnAmount)
	if err != nil {
		return unlockedAmount, err
	}

	// check unbond amount has been slashed or not
	if unbondEntry.Balance.GT(unlockMultiStakingAmount) {
		unlockedAmount = sdk.NewCoins(sdk.NewCoin(unlockDenom, unlockMultiStakingAmount))

		// Slash user amount
		burnUserAmount := sdk.NewCoins(sdk.NewCoin(unlockDenom, unbondEntry.Balance.Sub(unlockMultiStakingAmount)))
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, intermediaryAcc, types.ModuleName, burnUserAmount)
		if err != nil {
			return unlockedAmount, err
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnUserAmount)
		if err != nil {
			return unlockedAmount, err
		}
	} else {
		unlockedAmount = sdk.NewCoins(sdk.NewCoin(unlockDenom, unbondEntry.Balance))
	}

	err = k.bankKeeper.SendCoins(ctx, intermediaryAcc, delAddr, unlockedAmount)
	if err != nil {
		return unlockedAmount, err
	}

	return unlockedAmount, nil
}
