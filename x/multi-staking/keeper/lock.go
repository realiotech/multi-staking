package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) LockedAmountToBondAmount(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	lockedAmount math.Int,
) (math.Int, error) {
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		return math.Int{}, fmt.Errorf("can't find multi staking lock")
	}

	return lock.LockedAmountToBondAmount(lockedAmount).RoundInt(), nil
}

func (k Keeper) AddTokenToLock(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	amountAdded math.Int,
	conversionRatio sdk.Dec,
) types.MultiStakingLock {
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(amountAdded, conversionRatio, delAddr, valAddr)
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(amountAdded, conversionRatio)
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
	return multiStakingLock
}

// removing token from lock won't change the conversion ratio of the lock
func (k Keeper) RemoveTokenFromLock(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	amountRemoved math.Int,
) (types.MultiStakingLock, error) {
	// get lock on source val
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	lock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		return types.MultiStakingLock{}, fmt.Errorf("can't find multi staking lock")
	}

	// remove token from lock on source val
	lock, err := lock.RemoveTokenFromMultiStakingLock(amountRemoved)
	if err != nil {
		return types.MultiStakingLock{}, err
	}

	// update lock on source val
	k.SetMultiStakingLock(ctx, lockID, lock)
	return lock, nil
}

func (k Keeper) MoveLockedMultistakingToken(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	srcValAddr sdk.ValAddress,
	dstValAddr sdk.ValAddress,
	lockedToken sdk.Coin,
) (err error) {
	// get lock on source val
	srcLock, err := k.RemoveTokenFromLock(ctx, delAddr, srcValAddr, lockedToken.Amount)
	if err != nil {
		return err
	}

	// get lock on destination val
	k.AddTokenToLock(ctx, delAddr, dstValAddr, lockedToken.Amount, srcLock.ConversionRatio)

	return err
}

func (k Keeper) LockMultiStakingTokenAndMintBondToken(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	multiStakingToken sdk.Coin,
) (mintedBondToken sdk.Coin, err error) {
	intermediaryAcc := types.IntermediaryAccount(delAddr)

	// get bond denom weight
	bondDenomWeight, isBondToken := k.GetBondTokenWeight(ctx, multiStakingToken.Denom)
	if !isBondToken {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", multiStakingToken.Denom,
		)
	}

	// lock coin in intermediary account
	err = k.bankKeeper.SendCoins(ctx, delAddr, intermediaryAcc, sdk.NewCoins(multiStakingToken))
	if err != nil {
		return sdk.Coin{}, err
	}

	// update multistaking lock
	k.AddTokenToLock(ctx, delAddr, valAddr, multiStakingToken.Amount, bondDenomWeight)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking token * bond token weight
	mintedBondAmount := bondDenomWeight.MulInt(multiStakingToken.Amount).RoundInt()
	mintedBondToken = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond token to intermediary account
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondToken))
	if err != nil {
		return sdk.Coin{}, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(mintedBondToken))
	if err != nil {
		return sdk.Coin{}, err
	}
	return mintedBondToken, nil
}

func (k Keeper) BurnBondTokenAndUnlockMultiStakingToken(
	ctx sdk.Context,
	intermediaryAcc sdk.AccAddress,
	valAddr sdk.ValAddress,
	unbondTokenAmount sdk.Coin,
) (unlockedAmount sdk.Coins, err error) {
	// get delAddr
	delAddr := types.DelegatorAccount(intermediaryAcc)

	// get Lock
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		return unlockedAmount, fmt.Errorf("StakingLock not exists")
	}

	// unlock amount
	// unlockMultiStakingAmount = unbondTokenAmount/multiStakingLock.ConversionRatio
	unlockDenom := k.GetValidatorAllowedToken(ctx, valAddr)
	unlockMultiStakingAmount := sdk.NewDecFromInt(unbondTokenAmount.Amount).Quo(multiStakingLock.ConversionRatio).RoundInt()

	// check amount
	if unlockMultiStakingAmount.GT(multiStakingLock.LockedAmount) {
		return unlockedAmount, fmt.Errorf("unlock amount greater than lock amount")
	}

	// burn bonded token
	burnAmount := sdk.NewCoins(unbondTokenAmount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, intermediaryAcc, types.ModuleName, burnAmount)
	if err != nil {
		return unlockedAmount, err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnAmount)
	if err != nil {
		return unlockedAmount, err
	}

	// unlock coin
	unlockedAmount = sdk.NewCoins(sdk.NewCoin(unlockDenom, unlockMultiStakingAmount))
	err = k.bankKeeper.SendCoins(ctx, intermediaryAcc, delAddr, unlockedAmount)
	if err != nil {
		return unlockedAmount, err
	}

	return unlockedAmount, nil
}
