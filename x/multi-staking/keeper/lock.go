package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) LockMultiStakingTokenAndMintBondToken(
	ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, intermediaryAcc sdk.AccAddress,
	bondToken sdk.Coin,
) (sdk.Coin, error) {
	bondDenomWeight, isBondToken := k.GetBondTokenWeight(ctx, bondToken.Denom)
	if !isBondToken {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", bondToken.Denom,
		)
	}
	sdkBondAmount := bondDenomWeight.MulInt(bondToken.Amount).RoundInt()

	sdkBondToken := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdkBondAmount)

	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdkBondToken))

	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(sdkBondToken))

	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAcc) == nil {
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAcc, delAcc)
	}



	return sdkBondToken, nil
}

func (k Keeper) LockMultiStakingToken(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, intermediaryAcc sdk.AccAddress, lockedCoin sdk.Coin, currentConversionRatio sdk.Dec) error {
	err := k.bankKeeper.SendCoins(ctx, delAcc, intermediaryAcc, sdk.NewCoins(lockedCoin))
	if err != nil {
		return err
	}


	multiStakingLock, found := k.GetMultiStakingLock(ctx, delAcc, valAcc)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(lockedCoin.Amount, multiStakingLock.ConversionRatio, intermediaryAcc.String())
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(lockedCoin.Amount, currentConversionRatio)
	}
	k.SetMultiStakingLock(ctx, delAcc, valAcc, multiStakingLock)





}

informal-system
binary-builder
decentr-ware

block-gang

