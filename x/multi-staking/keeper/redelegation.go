package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) SharesFromBondToken(ctx sdk.Context, bondAmt math.Int, delegation stakingtypes.Delegation) sdk.Dec {
	bondAmount := k.GetDVPairBondAmount(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr())

	totalShares := delegation.Shares

	shares := totalShares.MulInt(bondAmt).QuoInt(bondAmount)

	return shares
}

func (k Keeper) CalSDKUnbondAmount(ctx sdk.Context, delegation stakingtypes.Delegation, unbondAmount sdk.Coin) (math.Int, error) {

	shares := k.SharesFromBondToken(ctx, unbondAmount.Amount, delegation)

	validator, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
	if !found {
		return math.Int{}, stakingtypes.ErrNoValidatorFound
	}

	sdkBondToken := validator.TokensFromShares(shares)

	return sdkBondToken.RoundInt(), nil
}

func (k Keeper) GetSDKDelegation(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress) (stakingtypes.Delegation, bool) {
	intermediaryAccount := types.GetIntermediaryAccount(delAcc.String(), valAcc.String())

	return k.stakingKeeper.GetDelegation(ctx, intermediaryAccount, valAcc)
}
