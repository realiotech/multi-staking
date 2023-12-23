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

func (k Keeper) AddCoinToLock(
	ctx sdk.Context,
	fromAcc sdk.AccAddress,
	lockID []byte,
	weightedCoin types.WeightedCoin,
) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, sdk.NewCoins(weightedCoin.ToCoin()))
	if err != nil {
		return err
	}

	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(nil, weightedCoin)
	} else {
		multiStakingLock, err = multiStakingLock.AddWeightedCoinToMultiStakingLock(weightedCoin)
		if err != nil {
			return err
		}
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
	return nil
}

// removing coin from lock won't change the conversion ratio of the lock
func (k Keeper) WithdrawCoinFromLock(
	ctx sdk.Context,
	fromLockID []byte,
	toAcc sdk.AccAddress,
	withdrawalCoin sdk.Coin,
) error {
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, fromLockID)
	if !found {
		return fmt.Errorf("can't find multi staking lock")
	}

	// remove coin from lock on source val
	lock, err := lock.RemoveCoinFromMultiStakingLock(withdrawalCoin)
	if err != nil {
		return err
	}

	// update lock on source val
	k.SetMultiStakingLock(ctx, fromLockID, lock)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcc, sdk.NewCoins(withdrawalCoin))
	return nil
}

func (k Keeper) MoveLockedMultistakingCoin(
	ctx sdk.Context,
	fromLockID []byte,
	toLockID []byte,
	movedCoin sdk.Coin,
) (err error) {
	// withdraw coin from source lock
	fromLock, found := k.GetMultiStakingLock(ctx, fromLockID)
	if !found {
		return fmt.Errorf("can't find multi staking lock")
	}
	weightedCoin := fromLock.ToWeightedCoin(movedCoin)

	// remove coin from lock on source val
	fromLock, err = fromLock.RemoveCoinFromMultiStakingLock(movedCoin)
	if err != nil {
		return err
	}

	k.SetMultiStakingLock(ctx, fromLockID, fromLock)

	// add coin to destination lock
	toLock, found := k.GetMultiStakingLock(ctx, toLockID)
	if !found {
		toLock = types.NewMultiStakingLock(nil, weightedCoin)
	} else {
		toLock, err = toLock.AddWeightedCoinToMultiStakingLock(weightedCoin)
		if err != nil {
			return err
		}
	}
	k.SetMultiStakingLock(ctx, toLockID, toLock)
	return nil
}

func (k Keeper) LockMultiStakingCoinAndMintBondCoin(
	ctx sdk.Context,
	lockID []byte,
	fromAcc sdk.AccAddress,
	mintedTo sdk.AccAddress,
	multiStakingCoin sdk.Coin,
) (mintedBondCoin sdk.Coin, err error) {
	// get bond denom weight
	bondDenomWeight, isBondCoin := k.GetBondCoinWeight(ctx, multiStakingCoin.Denom)
	if !isBondCoin {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", multiStakingCoin.Denom,
		)
	}

	// update multistaking lock
	err = k.AddCoinToLock(ctx, fromAcc, lockID, types.NewWeightedCoin(multiStakingCoin.Denom, multiStakingCoin.Amount, bondDenomWeight))

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking coin * bond coin weight
	mintedBondAmount := bondDenomWeight.MulInt(multiStakingCoin.Amount).RoundInt()
	mintedBondCoin = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond coin to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondCoin))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintedTo, sdk.NewCoins(mintedBondCoin))

	return mintedBondCoin, nil
}

func (k Keeper) BurnBondCoinAndUnlockMultiStakingCoin(
	ctx sdk.Context,
	intermediaryAcc sdk.AccAddress,
	valAddr sdk.ValAddress,

	unbondCoinAmount sdk.Coin,
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
	// unlockMultiStakingAmount = unbondCoinAmount/multiStakingLock.ConversionRatio
	unlockDenom := k.GetValidatorAllowedCoin(ctx, valAddr)
	unlockMultiStakingAmount := sdk.NewDecFromInt(unbondCoinAmount.Amount).Quo(multiStakingLock.LockedCoin.Weight).RoundInt()

	// check amount
	if unlockMultiStakingAmount.GT(multiStakingLock.LockedCoin.Amount) {
		return unlockedAmount, fmt.Errorf("unlock amount greater than lock amount")
	}

	// burn bonded coin
	burnAmount := sdk.NewCoins(unbondCoinAmount)
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
