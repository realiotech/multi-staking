package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) PreDelegate(
	ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress,
	bondToken sdk.Coin,
) (sdk.Coin, error) {
	// check if bond denom match val's bond denom
	valBondDenom := k.GetValidatorBondDenom(ctx, valAcc)
	if bondToken.Denom != valBondDenom {
		return sdk.Coin{}, fmt.Errorf("mismatch bond token; expect %s got %s", valBondDenom, bondToken.Denom)
	}

	// calculate converted sdk bond token
	sdkBondAmount, err := k.CalSDKBondAmount(ctx, bondToken)
	if err != nil {
		return sdk.Coin{}, err
	}
	sdkBondToken := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdkBondAmount)

	intermediaryAcc := types.GetIntermediaryAccount(delAcc.String(), valAcc.String())

	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdkBondToken))

	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(sdkBondToken))

	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAcc) == nil {
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAcc, delAcc)
	}

	k.UpdateDVPairBondAmount(ctx, delAcc, valAcc, bondToken.Amount)

	k.UpdateDVPairSDKBondAmount(ctx, delAcc, valAcc, sdkBondAmount)
	return sdkBondToken, nil
}

func (k Keeper) UpdateDVPairBondAmount(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, updateAmount math.Int) {
	existingBondAmount := k.GetDVPairBondAmount(ctx, delAcc, valAcc)
	if existingBondAmount.IsZero() {
		k.SetDVPairBondAmount(ctx, delAcc, valAcc, updateAmount)
	} else {
		k.SetDVPairBondAmount(ctx, delAcc, valAcc, existingBondAmount.Add(updateAmount))
	}
}

func (k Keeper) UpdateDVPairSDKBondAmount(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, updateAmount math.Int) {
	existingSDKBondAmount := k.GetDVPairSDKBondAmount(ctx, delAcc, valAcc)
	if existingSDKBondAmount.IsZero() {
		k.SetDVPairSDKBondAmount(ctx, delAcc, valAcc, updateAmount)
	} else {
		k.SetDVPairSDKBondAmount(ctx, delAcc, valAcc, existingSDKBondAmount.Add(updateAmount))
	}
}

func (k Keeper) CalSDKBondAmount(ctx sdk.Context, bondToken sdk.Coin) (math.Int, error) {
	bondDenomWeight, isBondToken := k.GetBondTokenWeight(ctx, bondToken.Denom)
	if !isBondToken {
		return math.Int{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", bondToken.Denom,
		)
	}
	sdkBondAmount := bondDenomWeight.MulInt(bondToken.Amount).RoundInt()

	return sdkBondAmount, nil
}

// func (k Keeper) MintSDKBondTo
