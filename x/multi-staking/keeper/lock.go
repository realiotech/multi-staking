package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) AddTokenToLock(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	amountAdded math.Int,
	conversionRatio sdk.Dec,
) {
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(amountAdded, conversionRatio, delAddr, valAddr)
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(amountAdded, conversionRatio)
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
}

func (k Keeper) LockMultiStakingTokenAndMintBondToken(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	multiStakingToken sdk.Coin,
) (mintedBondToken sdk.Coin, err error) {
	intermediaryAcc := k.GetIntermediaryAccountDelegator(ctx, delAddr)

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
	lockID := types.MultiStakingLockID(delAddr, valAddr)
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(multiStakingToken.Amount, multiStakingLock.ConversionRatio, delAddr, valAddr)
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(multiStakingToken.Amount, bondDenomWeight)
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking token * bond token weight
	mintedBondAmount := bondDenomWeight.MulInt(multiStakingToken.Amount).RoundInt()
	mintedBondToken = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond token to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondToken))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(mintedBondToken))

	return mintedBondToken, nil
}

func (k Keeper) MoveLockedMultistakingToken(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valSrcAddr sdk.ValAddress,
	valDstAddr sdk.ValAddress,
	lockedToken sdk.Coin,
) (err error) {
	// get lock on source val
	srcLockID := types.MultiStakingLockID(delAddr, valSrcAddr)
	srcLock, found := k.GetMultiStakingLock(ctx, srcLockID)
	if !found {
		return fmt.Errorf("can't find multi staking lock")
	}

	// remove token from lock on source val
	srcLock, err = srcLock.RemoveTokenFromMultiStakingLock(lockedToken.Amount)
	if err != nil {
		return err
	}

	// update lock on source val
	k.SetMultiStakingLock(ctx, srcLockID, srcLock)

	// get lock on destination val
	dstLockID := types.MultiStakingLockID(delAddr, valDstAddr)
	dstLock, found := k.GetMultiStakingLock(ctx, dstLockID)
	if !found {
		dstLock = types.NewMultiStakingLock(lockedToken.Amount, srcLock.ConversionRatio, delAddr, valDstAddr)
		k.SetMultiStakingLock(ctx, dstLockID, dstLock)
	} else {
		dstLock = dstLock.AddTokenToMultiStakingLock(lockedToken.Amount, srcLock.ConversionRatio)
	}

	// update lock on destination val
	k.SetMultiStakingLock(ctx, dstLockID, dstLock)

	return err
}

func (k Keeper) LockedAmountToBondAmount(ctx sdk.Context, multiStakingLockID []byte, lockedAmount math.Int) (math.Int, error) {
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, multiStakingLockID)
	if !found {
		return math.Int{}, fmt.Errorf("can't find multi staking lock")
	}

	return lock.LockedAmountToBondAmount(lockedAmount).RoundInt(), nil
}
