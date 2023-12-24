package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) AddCoinToLock(
	ctx sdk.Context,
	fromAcc sdk.AccAddress,
	lockID []byte,
	addedCoin types.MultiStakingCoin,
) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, sdk.NewCoins(addedCoin.ToCoin()))
	if err != nil {
		return err
	}

	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(nil, addedCoin)
	} else {
		multiStakingLock, err = multiStakingLock.AddCoinToMultiStakingLock(addedCoin)
		if err != nil {
			return err
		}
	}
	k.SetMultiStakingLock(ctx, lockID, multiStakingLock)
	return nil
}

// removing coin from lock won't change the conversion ratio of the lock
func (k Keeper) RemoveCoinFromLock(
	ctx sdk.Context,
	fromLockID []byte,
	removedCoin sdk.Coin,
) error {
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, fromLockID)
	if !found {
		return fmt.Errorf("can't find multi staking lock")
	}
	multiStakingCoin := lock.ToMultiStakingCoin(removedCoin)

	// remove coin from lock on source val
	lock, err := lock.RemoveCoinFromMultiStakingLock(multiStakingCoin)
	if err != nil {
		return err
	}

	// update lock on source val
	k.SetMultiStakingLock(ctx, fromLockID, lock)

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
	multiStakingCoin := fromLock.ToMultiStakingCoin(movedCoin)

	// remove coin from lock on source val
	fromLock, err = fromLock.RemoveCoinFromMultiStakingLock(multiStakingCoin)
	if err != nil {
		return err
	}

	k.SetMultiStakingLock(ctx, fromLockID, fromLock)

	// add coin to destination lock
	toLock, found := k.GetMultiStakingLock(ctx, toLockID)
	if !found {
		toLock = types.NewMultiStakingLock(nil, multiStakingCoin)
	} else {
		toLock, err = toLock.AddCoinToMultiStakingLock(multiStakingCoin)
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
	bondWeight, isBondCoin := k.GetBondCoinWeight(ctx, multiStakingCoin.Denom)
	if !isBondCoin {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", multiStakingCoin.Denom,
		)
	}

	// update multistaking lock
	err = k.AddCoinToLock(ctx, fromAcc, lockID, types.NewWeightedCoin(multiStakingCoin.Denom, multiStakingCoin.Amount, bondWeight))
	if err != nil {
		return sdk.Coin{}, err
	}

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking coin * bond coin weight
	mintedBondAmount := bondWeight.MulInt(multiStakingCoin.Amount).RoundInt()
	mintedBondCoin = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond coin to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondCoin))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintedTo, sdk.NewCoins(mintedBondCoin))

	return mintedBondCoin, nil
}
