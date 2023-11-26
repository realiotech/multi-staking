package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsAllowedToken(ctx sdk.Context, valAcc sdk.ValAddress, lockedToken sdk.Coin) bool {
	return lockedToken.Denom == k.GetValidatorBondDenom(ctx, valAcc)
}
