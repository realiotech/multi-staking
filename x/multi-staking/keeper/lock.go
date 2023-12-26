package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetOrCreateMultiStakingLock(ctx sdk.Context, lockID types.LockID) types.MultiStakingLock {
	multiStakingLock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(&lockID, types.MultiStakingCoin{Amount: sdk.ZeroInt()})
	}
	return multiStakingLock
}

func (k Keeper) EscrowCoinFromAcc(ctx sdk.Context, fromAcc sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, sdk.NewCoins(coin))
}

func (k Keeper) WithdrawEscrowCoinTo(ctx sdk.Context, toAcc sdk.AccAddress, coin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcc, sdk.NewCoins(coin))
}

func (k Keeper) LockCoinAndMintBondCoin(
	ctx sdk.Context,
	lockID types.LockID,
	fromAcc sdk.AccAddress,
	mintedTo sdk.AccAddress,
	coin sdk.Coin,
) (mintedBondCoin sdk.Coin, err error) {
	// escrow coin
	err = k.EscrowCoinFromAcc(ctx, fromAcc, coin)
	if err != nil {
		return sdk.Coin{}, err
	}

	// get multistaking coin's bond weight
	bondWeight, isMultiStakingCoin := k.GetBondCoinWeight(ctx, coin.Denom)
	if !isMultiStakingCoin {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", coin.Denom,
		)
	}

	// update multistaking lock
	multiStakingCoin := types.NewMultiStakingCoin(coin.Denom, coin.Amount, bondWeight)
	lock := k.GetOrCreateMultiStakingLock(ctx, lockID)
	err = lock.AddCoinToMultiStakingLock(multiStakingCoin)
	if err != nil {
		return sdk.Coin{}, err
	}

	k.SetMultiStakingLock(ctx, lock)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking coin * bond coin weight
	mintedBondAmount := multiStakingCoin.BondAmount()
	mintedBondCoin = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond coin to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondCoin))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintedTo, sdk.NewCoins(mintedBondCoin))

	return mintedBondCoin, nil
}
