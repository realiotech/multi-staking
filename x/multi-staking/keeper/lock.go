package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) LockedAmountToBondAmount(ctx sdk.Context, multiStakingLockID []byte, lockedAmount math.Int) (math.Int, error) {
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, multiStakingLockID)
	if !found {
		return math.Int{}, fmt.Errorf("can't find multi staking lock")
	}

	return lock.LockedAmountToBondAmount(lockedAmount).RoundInt(), nil
}

func (k Keeper) AddTokenToLock(ctx sdk.Context, lockID []byte, amountAdded math.Int, conversionRatio sdk.Dec) types.MultiStakingLock {
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(amountAdded, conversionRatio)
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(amountAdded, conversionRatio)
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
	return multiStakingLock
}

// removing token from lock won't change the conversion ratio of the lock
func (k Keeper) RemoveTokenFromLock(ctx sdk.Context, lockID []byte, amountRemoved math.Int) (types.MultiStakingLock, error) {
	// get lock on source val
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

func (k Keeper) MoveLockedMultistakingToken(ctx sdk.Context, srcLockID []byte, dstLockID []byte, lockedToken sdk.Coin) (err error) {
	// get lock on source val
	srcLock, err := k.RemoveTokenFromLock(ctx, srcLockID, lockedToken.Amount)
	if err != nil {
		return err
	}

	// get lock on destination val
	k.AddTokenToLock(ctx, dstLockID, lockedToken.Amount, srcLock.ConversionRatio)

	return err
}

func (k Keeper) LockMultiStakingTokenAndMintBondToken(
	ctx sdk.Context, delAcc sdk.AccAddress, lockID []byte,
	multiStakingToken sdk.Coin,
) (mintedBondToken sdk.Coin, err error) {
	intermediaryAcc := k.GetIntermediaryAccountDelegator(ctx, delAcc)

	// get bond denom weight
	bondDenomWeight, isBondToken := k.GetBondTokenWeight(ctx, multiStakingToken.Denom)
	if !isBondToken {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", multiStakingToken.Denom,
		)
	}

	// lock coin in intermediary account
	err = k.bankKeeper.SendCoins(ctx, delAcc, intermediaryAcc, sdk.NewCoins(multiStakingToken))
	if err != nil {
		return sdk.Coin{}, err
	}

	// update multistaking lock
	k.AddTokenToLock(ctx, lockID, multiStakingToken.Amount, bondDenomWeight)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking token * bond token weight
	mintedBondAmount := bondDenomWeight.MulInt(multiStakingToken.Amount).RoundInt()
	mintedBondToken = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond token to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondToken))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(mintedBondToken))

	return mintedBondToken, nil
}
